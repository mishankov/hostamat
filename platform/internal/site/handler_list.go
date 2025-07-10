package site

import (
	"encoding/json"
	"net/http"
)

type listSiteService interface {
	GetList() ([]*Site, error)
}

type ListSiteHandler struct {
	service listSiteService
}

func NewListSiteHandler(service listSiteService) *ListSiteHandler {
	return &ListSiteHandler{
		service: service,
	}
}

func (h *ListSiteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sites, err := h.service.GetList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}
