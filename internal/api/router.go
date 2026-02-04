package api

import (
	"net/http"
	"strings"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

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
	mux.HandleFunc("/v1/catalog/item-types/", func(w http.ResponseWriter, r *http.Request) {
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

	return mux
}
