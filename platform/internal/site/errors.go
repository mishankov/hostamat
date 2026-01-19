package site

import "errors"

var (
	errLongName   = errors.New("Название сайта должно быть не более 25 символов")
	errShortName  = errors.New("Название сайта должно быть не менее 1 символа")
	errUploadSize = errors.New("Достигнут максимальный размер загрузки: 20 мегабайт")
)

type errorCode = string

const (
	// Validation errors
	errorCodeInvalidExt       = "INVALID_EXT"
	errorCodeMissingSiteId    = "MISSING_SITE_ID"
	errorCodeMissingSessionId = "MISSING_SESSION_ID"
	errorCodeMissingName      = "MISSING_NAME"
	errorCodeInvalidBody      = "INVALID_BODY"
	errorCodeNameTooLong      = "NAME_TOO_LONG"
	errorCodeNameTooShort     = "NAME_TOO_SHORT"
	errorCodeOther            = "OTHER"
)

type handlerError struct {
	ErrorCode    errorCode `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
}

func (e handlerError) Error() string {
	return e.ErrorMessage
}
