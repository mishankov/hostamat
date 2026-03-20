package site

import (
	"context"
	"encoding/json"
	"log/slog"
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
	ev := log.EventFromContext(ctx)

	if r.Method != http.MethodGet {

		ev.AddStep(slog.LevelWarn, "site.List.method_not_allowed")

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sites, err := h.service.GetList(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		err := json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при получении списка сайтов",
		})
		if err != nil {
			ev.AddError(err)
		}

		return
	}

	resp := listReponse{sites, maxSites}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ev.AddError(err)
	}
}
