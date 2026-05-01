package publicapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"hostamat/internal/apikey"
	"hostamat/internal/file"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorWriter struct{}

func (ErrorWriter) WritePublicError(w http.ResponseWriter, err error) {
	RespondError(w, statusForError(err), codeForError(err), messageForError(err))
}

func RespondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func statusForError(err error) int {
	switch {
	case isAPIKeyError(err):
		return apikey.PublicStatusForError(err)
	case errors.Is(err, file.ErrInvalidFileType):
		return http.StatusBadRequest
	case strings.Contains(err.Error(), "not found"), strings.Contains(err.Error(), "не найден"):
		return http.StatusNotFound
	case strings.Contains(err.Error(), "access denied"), strings.Contains(err.Error(), "Доступ запрещен"):
		return http.StatusForbidden
	case strings.Contains(err.Error(), "max sites limit"):
		return http.StatusConflict
	case strings.Contains(err.Error(), "maximum"), strings.Contains(err.Error(), "максим"):
		return http.StatusRequestEntityTooLarge
	case strings.Contains(err.Error(), "too long"), strings.Contains(err.Error(), "не более"), strings.Contains(err.Error(), "корректный"), strings.Contains(err.Error(), "обязат"), strings.Contains(err.Error(), "Отсутствует"), strings.Contains(err.Error(), "поддомен"):
		return http.StatusBadRequest
	case strings.Contains(err.Error(), "уже подключен"):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func codeForError(err error) string {
	if isAPIKeyError(err) {
		return string(apikey.CodeForError(err))
	}

	switch {
	case errors.Is(err, file.ErrInvalidFileType):
		return "invalid_file_type"
	case strings.Contains(err.Error(), "max sites limit"):
		return "site_limit_reached"
	case strings.Contains(err.Error(), "not found"), strings.Contains(err.Error(), "не найден"):
		return "not_found"
	case strings.Contains(err.Error(), "access denied"), strings.Contains(err.Error(), "Доступ запрещен"):
		return "forbidden"
	case strings.Contains(err.Error(), "maximum"), strings.Contains(err.Error(), "максим"):
		return "upload_too_large"
	case strings.Contains(err.Error(), "уже подключен"):
		return "hostname_taken"
	case strings.Contains(err.Error(), "не более"):
		return "name_too_long"
	case strings.Contains(err.Error(), "не менее"), strings.Contains(err.Error(), "обязат"), strings.Contains(err.Error(), "Отсутствует"):
		return "missing_required_field"
	case strings.Contains(err.Error(), "корректный"), strings.Contains(err.Error(), "поддомен"):
		return "invalid_hostname"
	default:
		return "internal_error"
	}
}

func messageForError(err error) string {
	if isAPIKeyError(err) {
		return err.Error()
	}

	switch codeForError(err) {
	case "site_limit_reached":
		return "Site limit reached."
	case "not_found":
		return "Resource not found."
	case "forbidden":
		return "Access denied."
	case "upload_too_large":
		return "Upload exceeds the maximum allowed size."
	case "hostname_taken":
		return "Hostname is already connected to another site."
	case "name_too_long":
		return "Name is too long."
	case "missing_required_field":
		return "A required field is missing."
	case "invalid_hostname":
		return "Hostname is invalid."
	case "invalid_file_type":
		return "Invalid file type."
	default:
		return "Internal server error."
	}
}

func isAPIKeyError(err error) bool {
	return errors.Is(err, apikey.ErrMissingBearer) ||
		errors.Is(err, apikey.ErrMalformedBearer) ||
		errors.Is(err, apikey.ErrUnknownKey) ||
		errors.Is(err, apikey.ErrExpiredKey) ||
		errors.Is(err, apikey.ErrRevokedKey) ||
		errors.Is(err, apikey.ErrOwnerInactive) ||
		errors.Is(err, apikey.ErrInsufficientScope)
}
