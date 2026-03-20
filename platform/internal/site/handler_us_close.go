package site

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

type CloseUploadSessionHandler struct {
	service *Service
}

func NewCloseUploadSessionHandler(service *Service) *CloseUploadSessionHandler {
	return &CloseUploadSessionHandler{
		service: service,
	}
}

func (h *CloseUploadSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ev := log.EventFromContext(r.Context())

	if r.Method != http.MethodDelete {

		ev.AddStep(slog.LevelWarn, "site.CloseUploadSession.method_not_allowed")

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadSessionId := r.Header.Get("H-UploadSession-Id")
	if uploadSessionId == "" {

		ev.AddStep(slog.LevelWarn, "site.CloseUploadSession.missing_upload_session_id")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingSessionId,
			ErrorMessage: "Отсутствует ID сессии загрузки",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSessionId,
	})

	files, err := h.service.CloseUploadSession(r.Context(), uploadSessionId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при завершении загрузки",
		})
		if err != nil {
			ev.AddError(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(files); err != nil {

		ev.AddError(err)

		log.ErrorContext(r.Context(), "failed to encode body", "error", err.Error())
		return
	}
}
