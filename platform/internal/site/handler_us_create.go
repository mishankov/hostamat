package site

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

type CreateUploadSessionHandler struct {
	service *Service
}

func NewCreateUploadSessionHandler(service *Service) *CreateUploadSessionHandler {
	return &CreateUploadSessionHandler{
		service: service,
	}
}

func (h *CreateUploadSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ev := log.EventFromContext(r.Context())

	if r.Method != http.MethodPost {

		ev.AddStep(slog.LevelWarn, "site.CreateUploadSession.method_not_allowed")

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteId := r.Header.Get("H-Site-Id")
	if siteId == "" {

		ev.AddStep(slog.LevelWarn, "site.CreateUploadSession.missing_site_id")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingSiteId,
			ErrorMessage: "Отсутствует ID сайта",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	ev.AddAttrs(map[string]any{
		"site.id": siteId,
	})

	us, err := h.service.CreateUploadSession(r.Context(), siteId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при создании сессии загрузки",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)

	ev.AddAttrs(map[string]any{
		"uploadsession.id": us.Id,
	})

	_, err = w.Write([]byte(us.Id))
	if err != nil {
		ev.AddError(err)
	}
}
