package api

import (
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/catalog/item-types", h.GetCatalog)
	mux.HandleFunc("/v1/rent-actions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateRentAction(w, r)
		case http.MethodGet:
			// List rent actions (internal handler needed)
		}
	})
	mux.HandleFunc("/v1/rent-actions/", h.GetRentAction)

	return mux
}
