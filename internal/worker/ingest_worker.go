package worker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
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
		if err := w.SyncSource(ctx, src); err != nil {
			log.Printf("[IngestWorker] Failed to sync source %s (%d): %v", src.Name, src.ID, err)
		}
	}
}

func (w *IngestWorker) SyncSource(ctx context.Context, src domain.IngestSource) error {
	log.Printf("[IngestWorker] Syncing source: %s", src.Name)

	now := time.Now()
	src.LastSyncAt = &now

	// Prepare request
	req, err := http.NewRequestWithContext(ctx, "GET", src.APIURL, nil)
	if err != nil {
		return w.failSource(&src, "invalid request: "+err.Error())
	}

	if src.AuthType == domain.IngestAuthBearer {
		var credentials string
		json.Unmarshal(src.AuthCredentials, &credentials)
		if credentials != "" {
			req.Header.Set("Authorization", "Bearer "+credentials)
		}
	}

	if src.LastETag != "" {
		req.Header.Set("If-None-Match", src.LastETag)
	}

	// Perform request
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return w.failSource(&src, "request failed: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		log.Printf("[IngestWorker] Source %s unchanged (304)", src.Name)
		return w.finishSource(&src, "Not Modified", "")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return w.failSource(&src, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))
	}

	// Process body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return w.failSource(&src, "failed to read body: "+err.Error())
	}

	// Delta detection
	h := sha256.New()
	h.Write(body)
	hash := hex.EncodeToString(h.Sum(nil))
	if hash == src.LastPayloadHash {
		log.Printf("[IngestWorker] Source %s unchanged (Hash match)", src.Name)
		return w.finishSource(&src, "Success (No changes)", resp.Header.Get("ETag"))
	}

	// Parse JSON
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return w.failSource(&src, "failed to parse JSON: "+err.Error())
	}

	// Apply mappings and ingest
	count, err := w.ingestData(ctx, src, jsonData)
	if err != nil {
		return w.failSource(&src, "ingestion failed: "+err.Error())
	}

	log.Printf("[IngestWorker] Source %s finished. Ingested %d items.", src.Name, count)
	src.LastPayloadHash = hash
	return w.finishSource(&src, fmt.Sprintf("Success (%d items)", count), resp.Header.Get("ETag"))
}

func (w *IngestWorker) ingestData(ctx context.Context, src domain.IngestSource, data interface{}) (int, error) {
	// Root JSONPath mapping - if the API returns an array, we might want to iterate.
	// For simplicity, we assume if it's an array, we iterate over items.
	// In the future, we could add a "Root Path" to the IngestSource.

	items, ok := data.([]interface{})
	if !ok {
		// Single object
		items = []interface{}{data}
	}

	count := 0
	for _, item := range items {
		if err := w.ingestItem(ctx, src, item); err != nil {
			log.Printf("[IngestWorker] Skipping item in %s: %v", src.Name, err)
			continue
		}
		count++
	}

	return count, nil
}

func (w *IngestWorker) ingestItem(ctx context.Context, src domain.IngestSource, item interface{}) error {
	mapped := make(map[string]interface{})
	var identityValue interface{}

	for _, m := range src.Mappings {
		res, err := jsonpath.JsonPathLookup(item, m.JSONPath)
		if err != nil {
			if m.IsIdentity {
				return fmt.Errorf("identity field %s not found: %w", m.JSONPath, err)
			}
			continue
		}
		mapped[m.TargetField] = res
		if m.IsIdentity {
			identityValue = res
		}
	}

	if identityValue == nil {
		return fmt.Errorf("no identity value found for item")
	}

	// Ingest into target model
	switch src.TargetModel {
	case domain.IngestTargetItemType:
		return w.upsertItemType(ctx, mapped, identityValue)
	case domain.IngestTargetAsset:
		return w.upsertAsset(ctx, mapped, identityValue)
	case domain.IngestTargetCompany:
		return w.upsertCompany(ctx, mapped, identityValue)
	case domain.IngestTargetPerson:
		return w.upsertPerson(ctx, mapped, identityValue)
	case domain.IngestTargetPlace:
		return w.upsertPlace(ctx, mapped, identityValue)
	}

	return fmt.Errorf("unsupported target model: %s", src.TargetModel)
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
	}
	// For item types, default Active if not specified
	if _, ok := data["is_active"]; !ok {
		it.IsActive = true
	}

	return w.repo.UpsertItemType(ctx, it)
}

func (w *IngestWorker) upsertAsset(ctx context.Context, data map[string]interface{}, identity interface{}) error {
	a := &domain.Asset{}
	idenStr := fmt.Sprintf("%v", identity)

	// Determine if identity is tag or serial
	// For now, we'll try to match target fields
	if tag, ok := data["asset_tag"].(string); ok {
		a.AssetTag = &tag
	} else if _, ok := data["asset_tag"]; !ok {
		// If identity was mapped to asset_tag implicitly
		a.AssetTag = &idenStr
	}

	if sn, ok := data["serial_number"].(string); ok {
		a.SerialNumber = &sn
	}

	if status, ok := data["status"].(string); ok {
		a.Status = domain.AssetStatus(status)
	} else {
		a.Status = domain.AssetStatusAvailable // Default
	}

	// item_type_id is required. If not provided, we might need to look it up by code
	if itID, ok := data["item_type_id"].(float64); ok {
		a.ItemTypeID = int64(itID)
	} else if itIDStr, ok := data["item_type_id"].(string); ok {
		id, _ := strconv.ParseInt(itIDStr, 10, 64)
		a.ItemTypeID = id
	}

	if a.ItemTypeID == 0 {
		return fmt.Errorf("item_type_id is missing for asset %s", idenStr)
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
	if gn, ok := data["given_name"].(string); ok {
		p.GivenName = gn
	}
	if fn, ok := data["family_name"].(string); ok {
		p.FamilyName = fn
	}
	if cid, ok := data["company_id"].(float64); ok {
		p.CompanyID = ptrInt64(int64(cid))
	}

	if p.GivenName == "" || p.FamilyName == "" {
		return fmt.Errorf("given_name and family_name are required for person")
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

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrString(v string) *string {
	return &v
}

func (w *IngestWorker) failSource(src *domain.IngestSource, msg string) error {
	src.LastStatus = "Error"
	src.LastError = msg
	next := time.Now().Add(time.Duration(src.SyncIntervalSeconds) * time.Second)
	src.NextSyncAt = &next
	return w.repo.UpdateIngestSource(context.Background(), src)
}

func (w *IngestWorker) finishSource(src *domain.IngestSource, status, etag string) error {
	src.LastStatus = status
	src.LastError = ""
	now := time.Now()
	src.LastSuccessAt = &now
	if etag != "" {
		src.LastETag = etag
	}
	next := now.Add(time.Duration(src.SyncIntervalSeconds) * time.Second)
	src.NextSyncAt = &next
	return w.repo.UpdateIngestSource(context.Background(), src)
}
