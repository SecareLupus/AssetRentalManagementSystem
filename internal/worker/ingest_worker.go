package worker

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/desmond/rental-management-system/internal/db"
	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/oliveagle/jsonpath"
)

type IngestWorker struct {
	repo       db.Repository
	httpClient *http.Client
}

func NewIngestWorker(repo db.Repository) *IngestWorker {
	return &IngestWorker{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (w *IngestWorker) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial run
	w.ProcessPendingSources(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.ProcessPendingSources(ctx)
		}
	}
}

func (w *IngestWorker) ProcessPendingSources(ctx context.Context) {
	sources, err := w.repo.GetPendingIngestSources(ctx)
	if err != nil {
		log.Printf("[IngestWorker] Failed to fetch pending sources: %v", err)
		return
	}

	for _, src := range sources {
		if err := w.SyncSource(ctx, &src); err != nil {
			log.Printf("[IngestWorker] Failed to sync source %s (%d): %v", src.Name, src.ID, err)
		}
	}
}

func (w *IngestWorker) SyncSource(ctx context.Context, src *domain.IngestSource) error {
	log.Printf("[IngestWorker] Syncing source: %s", src.Name)

	// 1. Check/Refresh Authentication
	if src.AuthType == domain.IngestAuthBearer {
		if src.LastToken == "" || (src.TokenExpiry != nil && time.Now().After(src.TokenExpiry.Add(-5*time.Minute))) {
			if err := w.refreshAuth(ctx, src); err != nil {
				return fmt.Errorf("auth refresh failed: %w", err)
			}
		}
	}

	// 2. Iterate Endpoints
	for i := range src.Endpoints {
		ep := &src.Endpoints[i]
		if !ep.IsActive {
			continue
		}

		if err := w.SyncEndpoint(ctx, src, ep); err != nil {
			log.Printf("[IngestWorker] Failed to sync endpoint %s for source %s: %v", ep.Path, src.Name, err)
			continue
		}
	}

	src.UpdatedAt = time.Now()
	return w.repo.UpdateIngestSource(ctx, src)
}

func (w *IngestWorker) refreshAuth(ctx context.Context, src *domain.IngestSource) error {
	// Priority 1: Use Refresh Token if RefreshEndpoint is configured and we have a token
	if src.RefreshEndpoint != "" && src.RefreshToken != "" {
		log.Printf("[IngestWorker] Attempting token refresh for source: %s", src.Name)

		refreshURL := src.RefreshEndpoint
		// If it doesn't look like a full URL, prefix with BaseURL
		if !bytes.HasPrefix([]byte(refreshURL), []byte("http")) {
			refreshURL = src.BaseURL + refreshURL
		}

		refreshPayload := map[string]string{
			"refresh_token": src.RefreshToken,
		}
		payloadBytes, _ := json.Marshal(refreshPayload)

		req, err := http.NewRequestWithContext(ctx, "POST", refreshURL, bytes.NewReader(payloadBytes))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := w.httpClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				var data map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
					accessToken, refreshToken, expiresIn := domain.DiscoverTokens(data)
					if accessToken != "" {
						src.LastToken = accessToken
						if refreshToken != "" {
							src.RefreshToken = refreshToken
						}
						if expiresIn > 0 {
							expiry := time.Now().Add(time.Duration(expiresIn) * time.Second)
							src.TokenExpiry = &expiry
						}
						log.Printf("[IngestWorker] Token refreshed via Refresh Token for %s", src.Name)
						return nil
					}
				}
			}
			log.Printf("[IngestWorker] Refresh token attempt failed for %s (Status: %d), falling back to full login", src.Name, resp.StatusCode)
		} else {
			log.Printf("[IngestWorker] Refresh token request failed for %s: %v, falling back to full login", src.Name, err)
		}
	}

	// Priority 2: Standard Login (Existing Logic)
	log.Printf("[IngestWorker] Performing full login for source: %s", src.Name)

	authURL := src.AuthEndpoint
	if !strings.HasPrefix(authURL, "http") {
		authURL = src.BaseURL + authURL
	}

	creds := domain.UnwrapJSON(src.AuthCredentials)
	resp, err := w.httpClient.Post(authURL, "application/json", bytes.NewReader(creds))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth endpoint returned %s: %s", resp.Status, string(body))
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	// Step 1.5: Verify (Optional Triple-Auth)
	if src.VerifyEndpoint != "" {
		verifyURL := src.VerifyEndpoint
		if !strings.HasPrefix(verifyURL, "http") {
			verifyURL = src.BaseURL + verifyURL
		}
		log.Printf("[IngestWorker] Performing verification step for %s: %s", src.Name, verifyURL)

		// Attempt to find a token to use for verification if needed
		loginToken, _, _ := domain.DiscoverTokens(data)

		// Re-encode login response to send to verify
		verifyBody, _ := json.Marshal(data)
		req, err := http.NewRequest("POST", verifyURL, bytes.NewReader(verifyBody))
		if err != nil {
			return fmt.Errorf("failed to create verify request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if loginToken != "" {
			req.Header.Set("Authorization", "Bearer "+loginToken)
		}

		vResp, err := w.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("verify request failed: %w", err)
		}
		defer vResp.Body.Close()

		if vResp.StatusCode < 200 || vResp.StatusCode >= 300 {
			vBody, _ := io.ReadAll(vResp.Body)
			return fmt.Errorf("verify endpoint returned %s: %s", vResp.Status, string(vBody))
		}

		// Update data with verify response
		if err := json.NewDecoder(vResp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode verify response: %w", err)
		}
	}

	accessToken, refreshToken, expiresIn := domain.DiscoverTokens(data)
	if accessToken == "" {
		return fmt.Errorf("no access token found in auth response")
	}

	src.LastToken = accessToken
	if refreshToken != "" {
		src.RefreshToken = refreshToken
	}
	if expiresIn > 0 {
		expiry := time.Now().Add(time.Duration(expiresIn) * time.Second)
		src.TokenExpiry = &expiry
	}
	return nil
}

func (w *IngestWorker) SyncEndpoint(ctx context.Context, src *domain.IngestSource, ep *domain.IngestEndpoint) error {
	log.Printf("[IngestWorker] Syncing endpoint: %s%s", src.BaseURL, ep.Path)

	fullURL := src.BaseURL + ep.Path
	body := domain.UnwrapJSON(ep.RequestBody)

	// Refined body handling
	var bodyReader io.Reader
	sendBody := false
	if len(body) > 0 && string(body) != "null" {
		method := strings.ToUpper(ep.Method)
		if method == "POST" || method == "PUT" || method == "PATCH" {
			sendBody = true
			bodyReader = bytes.NewReader(body)
		}
	}

	req, err := http.NewRequestWithContext(ctx, ep.Method, fullURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	if sendBody {
		req.Header.Set("Content-Type", "application/json")
	}

	if src.AuthType == domain.IngestAuthBearer && src.LastToken != "" {
		req.Header.Set("Authorization", "Bearer "+src.LastToken)
	}

	if ep.LastETag != "" {
		req.Header.Set("If-None-Match", ep.LastETag)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return err
	}

	// Retry once on 401 or 403 if using Bearer auth
	if (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) && src.AuthType == domain.IngestAuthBearer {
		resp.Body.Close()
		log.Printf("[IngestWorker] Received %d for %s, attempting token refresh", resp.StatusCode, ep.Path)

		if err := w.refreshAuth(ctx, src); err != nil {
			return fmt.Errorf("auth refresh failed after %d: %w", resp.StatusCode, err)
		}

		// Re-create request
		req, err = http.NewRequestWithContext(ctx, ep.Method, fullURL, bodyReader)
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/json")
		if sendBody {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Authorization", "Bearer "+src.LastToken)
		if ep.LastETag != "" {
			req.Header.Set("If-None-Match", ep.LastETag)
		}

		resp, err = w.httpClient.Do(req)
		if err != nil {
			return err
		}
	}
	defer resp.Body.Close()

	now := time.Now()
	ep.LastSyncAt = &now

	if resp.StatusCode == http.StatusNotModified {
		log.Printf("[IngestWorker] Endpoint %s unchanged (304)", ep.Path)
		return w.repo.UpdateIngestEndpoint(ctx, ep)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Deduplication
	h := sha256.New()
	h.Write(body)
	hash := hex.EncodeToString(h.Sum(nil))
	if hash == ep.LastPayloadHash {
		log.Printf("[IngestWorker] Endpoint %s unchanged (Hash match)", ep.Path)
		ep.LastETag = resp.Header.Get("ETag")
		return w.repo.UpdateIngestEndpoint(ctx, ep)
	}

	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return err
	}

	// Strategy-based extraction
	items := w.extractItems(jsonData, ep.RespStrategy, ep.ItemsPath)

	count := 0
	for _, item := range items {
		if err := w.ingestItem(ctx, ep.Mappings, item); err == nil {
			count++
		}
	}

	log.Printf("[IngestWorker] Endpoint %s finished. Ingested %d items.", ep.Path, count)
	ep.LastPayloadHash = hash
	ep.LastETag = resp.Header.Get("ETag")
	ep.LastSuccessAt = &now
	return w.repo.UpdateIngestEndpoint(ctx, ep)
}

func (w *IngestWorker) extractItems(data interface{}, strategy string, itemsPath string) []interface{} {
	if itemsPath != "" && itemsPath != "$" {
		res, err := jsonpath.JsonPathLookup(data, itemsPath)
		if err == nil {
			if list, ok := res.([]interface{}); ok {
				return list
			}
			return []interface{}{res}
		}
		log.Printf("[IngestWorker] JSONPath lookup failed for %s: %v", itemsPath, err)
	}

	switch strategy {
	case "list":
		if list, ok := data.([]interface{}); ok {
			return list
		}
	case "single":
		return []interface{}{data}
	case "auto":
		if list, ok := data.([]interface{}); ok {
			return list
		}
		return []interface{}{data}
	}
	return nil
}

func (w *IngestWorker) ingestItem(ctx context.Context, mappings []domain.IngestMapping, item interface{}) error {
	// Group mappings by target model
	models := make(map[domain.IngestTargetModel]map[string]interface{})
	identities := make(map[domain.IngestTargetModel]interface{})

	for _, m := range mappings {
		res, err := jsonpath.JsonPathLookup(item, m.JSONPath)
		if err != nil {
			continue
		}

		if _, ok := models[m.TargetModel]; !ok {
			models[m.TargetModel] = make(map[string]interface{})
		}
		models[m.TargetModel][m.TargetField] = res
		if m.IsIdentity {
			identities[m.TargetModel] = res
		}
	}

	// Upsert each model
	for model, data := range models {
		identity := identities[model]
		if identity == nil {
			continue // Identity is required
		}

		var err error
		switch model {
		case domain.IngestTargetItemType:
			err = w.upsertItemType(ctx, data, identity)
		case domain.IngestTargetAsset:
			err = w.upsertAsset(ctx, data, identity)
		case domain.IngestTargetCompany:
			err = w.upsertCompany(ctx, data, identity)
		case domain.IngestTargetPerson:
			err = w.upsertPerson(ctx, data, identity)
		case domain.IngestTargetPlace:
			err = w.upsertPlace(ctx, data, identity)
		}
		if err != nil {
			log.Printf("[IngestWorker] Upsert failed for %s: %v", model, err)
		}
	}

	return nil
}

func (w *IngestWorker) upsertItemType(ctx context.Context, data map[string]interface{}, identity interface{}) error {
	it := &domain.ItemType{
		Code: fmt.Sprintf("%v", identity),
	}
	if name, ok := data["name"].(string); ok {
		it.Name = name
	}
	if kind, ok := data["kind"].(string); ok {
		it.Kind = domain.ItemKind(kind)
	}
	if isActive, ok := data["is_active"].(bool); ok {
		it.IsActive = isActive
	} else {
		it.IsActive = true
	}
	return w.repo.UpsertItemType(ctx, it)
}

func (w *IngestWorker) upsertAsset(ctx context.Context, data map[string]interface{}, identity interface{}) error {
	a := &domain.Asset{}
	idenStr := fmt.Sprintf("%v", identity)

	if tag, ok := data["asset_tag"].(string); ok {
		a.AssetTag = &tag
	} else {
		a.AssetTag = &idenStr
	}

	if sn, ok := data["serial_number"].(string); ok {
		a.SerialNumber = &sn
	}
	if status, ok := data["status"].(string); ok {
		a.Status = domain.AssetStatus(status)
	} else {
		a.Status = domain.AssetStatusAvailable
	}

	// Support ItemType lookup or code
	if itID, ok := data["item_type_id"].(float64); ok {
		a.ItemTypeID = int64(itID)
	}

	if a.ItemTypeID == 0 {
		return fmt.Errorf("item_type_id missing")
	}

	return w.repo.UpsertAsset(ctx, a)
}

func (w *IngestWorker) upsertCompany(ctx context.Context, data map[string]interface{}, identity interface{}) error {
	c := &domain.Company{
		Name: fmt.Sprintf("%v", identity),
	}
	if legalName, ok := data["legal_name"].(string); ok {
		c.LegalName = &legalName
	}
	if desc, ok := data["description"].(string); ok {
		c.Description = &desc
	}
	return w.repo.UpsertCompany(ctx, c)
}

func (w *IngestWorker) upsertPerson(ctx context.Context, data map[string]interface{}, identity interface{}) error {
	p := &domain.Person{}
	// Identity for Person is tricky, usually email or employee ID
	// For now assume given_name as identity if no employee_id
	if gn, ok := data["given_name"].(string); ok {
		p.GivenName = gn
	}
	if fn, ok := data["family_name"].(string); ok {
		p.FamilyName = fn
	}
	if cid, ok := data["company_id"].(float64); ok {
		id := int64(cid)
		p.CompanyID = &id
	}

	if p.GivenName == "" || p.FamilyName == "" {
		return fmt.Errorf("incomplete person data")
	}

	return w.repo.UpsertPerson(ctx, p)
}

func (w *IngestWorker) upsertPlace(ctx context.Context, data map[string]interface{}, identity interface{}) error {
	p := &domain.Place{
		Name: fmt.Sprintf("%v", identity),
	}
	if desc, ok := data["description"].(string); ok {
		p.Description = &desc
	}
	if cat, ok := data["category"].(string); ok {
		p.Category = &cat
	}
	if isInt, ok := data["is_internal"].(bool); ok {
		p.IsInternal = isInt
	}
	return w.repo.UpsertPlace(ctx, p)
}
