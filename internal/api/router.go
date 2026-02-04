package api

import (
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	// Auth
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
		if r.Method == http.MethodPost {
			h.CreateInspectionTemplate(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
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
	mux.HandleFunc("/v1/fleet/item-types/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/recall") {
			if r.Method == http.MethodPost {
				h.RecallItemTypeAssets(w, r)
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

	// Inventory (Assets)
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
		if strings.HasPrefix(r.URL.Path, "/v1/fleet/assets/") {
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
		}
		if strings.HasSuffix(r.URL.Path, "/required-inspections") {
			if r.Method == http.MethodGet {
				h.GetRequiredInspections(w, r)
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

	// Rent Actions
	mux.HandleFunc("/v1/rent-actions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateRentAction(w, r)
		case http.MethodGet:
			// List rent actions (internal handler needed)
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

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	return mux
}
