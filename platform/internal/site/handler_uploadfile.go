package site

import (
	"encoding/json"
	"log/slog"
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
	ev := log.EventFromContext(ctx)

	if r.Method != http.MethodPost {

		ev.AddStep(slog.LevelWarn, "site.UploadFile.method_not_allowed")

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadSessionId := r.Header.Get("H-UploadSession-Id")
	if uploadSessionId == "" {

		ev.AddStep(slog.LevelWarn, "site.UploadFile.missing_upload_session_id")

		http.Error(w, "Missing upload session id", http.StatusBadRequest)
		return
	}

	ev.AddAttrs(map[string]any{
		"uploadsession.id": uploadSessionId,
	})

	if err := r.ParseMultipartForm(20 << 20); err != nil { // 20 MB

		ev.AddStep(slog.LevelWarn, "site.UploadFile.invalid_multipart")

		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {

		ev.AddStep(slog.LevelWarn, "site.UploadFile.missing_file")

		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ev.AddAttrs(map[string]any{
		"file.name": handler.Filename,
		"file.size": handler.Size,
	})

	err = h.service.UploadFile(ctx, uploadSessionId, handler.Filename, file, handler.Size)
	if err != nil {
		switch err.(type) {
		case handlerError:
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(err)
			if err != nil {
				ev.AddError(err)
			}
			return
		default:
			log.WarnContext(ctx, "error with no error code", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(handlerError{ErrorCode: errorCodeOther, ErrorMessage: err.Error()})
			if err != nil {
				ev.AddError(err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}
