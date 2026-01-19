package site

import (
	"encoding/json"
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

type UploadFileHandler struct {
	service *Service
}

func NewUploadFileHandler(service *Service) *UploadFileHandler {
	return &UploadFileHandler{
		service: service,
	}
}
func (h *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadSessionId := r.Header.Get("H-UploadSession-Id")
	if uploadSessionId == "" {
		http.Error(w, "Missing upload session id", http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(20 << 20) // 20 MB
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.service.UploadFile(ctx, uploadSessionId, handler.Filename, file, handler.Size)
	if err != nil {
		switch err.(type) {
		case handlerError:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		default:
			log.WarnContext(ctx, "error with no error code", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(handlerError{ErrorCode: errorCodeOther, ErrorMessage: err.Error()})
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}
