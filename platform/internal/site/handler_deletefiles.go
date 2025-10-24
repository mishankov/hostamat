package site

import "net/http"

type DeleteFilesHandler struct {
	service *Service
}

func NewDeleteFileHandler(service *Service) *DeleteFilesHandler {
	return &DeleteFilesHandler{
		service: service,
	}
}
func (h *DeleteFilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	domain := r.Host

	err := h.service.DeleteFiles(ctx, domain)
	if err != nil {
		http.Error(w, "Failed to delete files", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
