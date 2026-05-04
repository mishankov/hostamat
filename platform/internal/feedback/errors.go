package feedback

import "errors"

var (
	errMissingMessage = errors.New("Напишите сообщение")
	errLongMessage    = errors.New("Сообщение должно быть не длиннее 5000 символов")
	errLongURL        = errors.New("Адрес страницы слишком длинный")
	errMissingEmail   = errors.New("В аккаунте не найден email для ответа")
	errRateLimited    = errors.New("Можно отправить не более 5 сообщений за час")
)

type errorCode = string

const (
	errorCodeInvalidBody    errorCode = "INVALID_BODY"
	errorCodeMissingMessage errorCode = "MISSING_MESSAGE"
	errorCodeMessageTooLong errorCode = "MESSAGE_TOO_LONG"
	errorCodeURLTooLong     errorCode = "URL_TOO_LONG"
	errorCodeMissingEmail   errorCode = "MISSING_EMAIL"
	errorCodeRateLimited    errorCode = "RATE_LIMITED"
	errorCodeUnauthorized   errorCode = "UNAUTHORIZED"
	errorCodeOther          errorCode = "OTHER"
)

type handlerError struct {
	ErrorCode    errorCode `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
}
