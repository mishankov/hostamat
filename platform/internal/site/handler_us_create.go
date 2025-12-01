package site

import "net/http"

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
		http.Error(w, "Missing site id", http.StatusBadRequest)
		return
	}

	us, err := h.service.CreateUploadSession(r.Context(), siteId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(us.Id))
}
