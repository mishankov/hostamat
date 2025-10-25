package site

import "net/http"

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
		http.Error(w, "Missing upload session id", http.StatusBadRequest)
		return
	}

	err := h.service.CloseUploadSession(r.Context(), uploadSessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
