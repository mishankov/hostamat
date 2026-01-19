package site

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

type listSiteService interface {
	GetList(context.Context) ([]*Site, error)
}

type ListSiteHandler struct {
	service listSiteService
}

type listReponse struct {
	Sites    []*Site `json:"sites"`
	MaxAmoun int     `json:"maxAmount"`
}

func NewListSiteHandler(service listSiteService) *ListSiteHandler {
	return &ListSiteHandler{
		service: service,
	}
}

func (h *ListSiteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sites, err := h.service.GetList(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при получении списка сайтов",
		})
		return
	}

	resp := listReponse{sites, maxSites}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.ErrorContext(ctx, "failed to send response", "error", err, "resp", resp)
	}
}
