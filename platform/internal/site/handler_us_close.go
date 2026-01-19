package site

import (
	"encoding/json"
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
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadSessionId := r.Header.Get("H-UploadSession-Id")
	if uploadSessionId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingSessionId,
			ErrorMessage: "Отсутствует ID сессии загрузки",
		})
		return
	}

	files, err := h.service.CloseUploadSession(r.Context(), uploadSessionId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при завершении загрузки",
		})
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(files); err != nil {
		log.ErrorContext(r.Context(), "failed to encode body", "error", err.Error())
		return
	}
}
