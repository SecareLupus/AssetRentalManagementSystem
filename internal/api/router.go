package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:embed web/dist/*
var uiFS embed.FS

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("[%s] %s -> %d (%s)", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	// UI (Embedded)
	// We handle the root and static files here.
	distFS, err := fs.Sub(uiFS, "web/dist")
	if err != nil {
		log.Printf("Warning: failed to sub uiFS: %v. UI will be unavailable.", err)
	}
	uiServer := http.FileServer(http.FS(distFS))

	// System & UI
	mux.HandleFunc("/v1/health", h.Health)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Let API and Swagger through
		if strings.HasPrefix(path, "/v1/") || strings.HasPrefix(path, "/swagger/") {
			http.NotFound(w, r)
			return
		}

		// Serve static files if they exist in the embedded FS
		_, err := distFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil || path == "/" {
			uiServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routing
		r.URL.Path = "/"
		uiServer.ServeHTTP(w, r)
	})

	// Auth (Public)
	mux.HandleFunc("/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Login(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Register(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Catalog (ItemTypes)
	mux.HandleFunc("/v1/catalog/item-types", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateItemType(w, r)
		case http.MethodGet:
			h.GetCatalog(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/catalog/inspection-templates", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateInspectionTemplate(w, r)
		case http.MethodGet:
			h.ListInspectionTemplates(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/fleet/build-specs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateBuildSpec(w, r)
		case http.MethodGet:
			h.ListBuildSpecs(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Dashboard
	mux.HandleFunc("/v1/dashboard/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetDashboardStats(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/v1/catalog/item-types/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/inspections") {
			if r.Method == http.MethodPost {
				h.SetItemTypeInspections(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetItemType(w, r)
		case http.MethodPut:
			h.UpdateItemType(w, r)
		case http.MethodDelete:
			h.DeleteItemType(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/catalog/inspection-templates/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetInspectionTemplate(w, r)
		case http.MethodPut:
			h.UpdateInspectionTemplate(w, r)
		case http.MethodDelete:
			h.DeleteInspectionTemplate(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/fleet/item-types/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/recall") {
			if r.Method == http.MethodPost {
				h.RecallItemTypeAssets(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	mux.HandleFunc("/v1/inventory/assets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateAsset(w, r)
		case http.MethodGet:
			h.ListAssets(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/inventory/assets/bulk-recall", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.BulkRecallAssets(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/v1/inventory/reconcile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.VerifyInventory(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/v1/inventory/assets/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/status") {
			if r.Method == http.MethodPatch {
				h.UpdateAssetStatus(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/provision") {
			if r.Method == http.MethodPost {
				h.StartProvisioning(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/complete-provisioning") {
			if r.Method == http.MethodPost {
				h.CompleteProvisioning(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/repair") {
			if r.Method == http.MethodPost {
				h.RepairAsset(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/refurbish") {
			if r.Method == http.MethodPost {
				h.RefurbishAsset(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/required-inspections") {
			if r.Method == http.MethodGet {
				h.GetRequiredInspections(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/maintenance-logs") {
			if r.Method == http.MethodGet {
				h.ListMaintenanceLogs(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/inspections") {
			if r.Method == http.MethodPost {
				h.SubmitInspection(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetAsset(w, r)
		case http.MethodPut:
			h.UpdateAsset(w, r)
		case http.MethodDelete:
			h.DeleteAsset(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/fleet/assets/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/remote-status") {
			if r.Method == http.MethodGet {
				h.GetAssetRemoteStatus(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/remote-power") {
			if r.Method == http.MethodPost {
				h.ApplyAssetRemotePower(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	// Rent Actions
	mux.HandleFunc("/v1/rent-actions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateRentAction(w, r)
		case http.MethodGet:
			h.ListRentActions(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/rent-actions/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/submit") {
			h.SubmitRentAction(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/approve") {
			h.ApproveRentAction(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/reject") {
			h.RejectRentAction(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/cancel") {
			h.CancelRentAction(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetRentAction(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Intelligence
	mux.HandleFunc("/v1/intelligence/availability", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetAvailabilityTimeline(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/v1/intelligence/shortage-alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetShortageAlerts(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/v1/intelligence/maintenance-forecast", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetMaintenanceForecast(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Entities (Phase 24 Convergence)
	mux.HandleFunc("/v1/entities/companies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateCompany(w, r)
		case http.MethodGet:
			h.ListCompanies(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/entities/companies/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetCompany(w, r)
		case http.MethodPut:
			h.UpdateCompany(w, r)
		case http.MethodDelete:
			h.DeleteCompany(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/entities/people", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreatePerson(w, r)
		case http.MethodGet:
			h.ListPeople(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/entities/people/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetPerson(w, r)
		case http.MethodPut:
			h.UpdatePerson(w, r)
		case http.MethodDelete:
			h.DeletePerson(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/entities/roles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateOrganizationRole(w, r)
		case http.MethodGet:
			h.ListOrganizationRoles(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/entities/places", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreatePlace(w, r)
		case http.MethodGet:
			h.ListPlaces(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/entities/places/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetPlace(w, r)
		case http.MethodPut:
			h.UpdatePlace(w, r)
		case http.MethodDelete:
			h.DeletePlace(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/entities/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateEvent(w, r)
		case http.MethodGet:
			h.ListEvents(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/v1/entities/events/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/needs") {
			switch r.Method {
			case http.MethodPost:
				h.UpdateEventAssetNeeds(w, r)
			case http.MethodGet:
				h.ListEventAssetNeeds(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.GetEvent(w, r)
		case http.MethodPut:
			h.UpdateEvent(w, r)
		case http.MethodDelete:
			h.DeleteEvent(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Swagger UI (Public)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Wrap entire mux in LoggingMiddleware
	handler := LoggingMiddleware(mux)

	// Apply AuthMiddleware to all /v1 routes EXCEPT public ones
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Skip auth for health, login, register, and swagger
		if path == "/v1/health" || strings.HasPrefix(path, "/v1/auth/") || strings.HasPrefix(path, "/swagger/") || path == "/" {
			handler.ServeHTTP(w, r)
			return
		}

		// Require auth for all other /v1 routes
		if strings.HasPrefix(path, "/v1/") {
			h.AuthMiddleware(handler).ServeHTTP(w, r)
			return
		}

		// Fallback for any other routes
		handler.ServeHTTP(w, r)
	})
}
