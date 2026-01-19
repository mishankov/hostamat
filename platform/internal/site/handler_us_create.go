package site

import (
	"encoding/json"
	"net/http"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteId := r.Header.Get("H-Site-Id")
	if siteId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeMissingSiteId,
			ErrorMessage: "Отсутствует ID сайта",
		})
		return
	}

	us, err := h.service.CreateUploadSession(r.Context(), siteId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlerError{
			ErrorCode:    errorCodeOther,
			ErrorMessage: "Ошибка сервера при создании сессии загрузки",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(us.Id))
}
