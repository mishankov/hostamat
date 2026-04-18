package customdomain

import "errors"

var (
	errInvalidHostname = errors.New("Укажите корректный домен без схемы и пути")
	errMissingHostname = errors.New("Домен обязателен")
	errHostnameTaken   = errors.New("Этот домен уже подключен к другому сайту")
	errAccessDenied    = errors.New("Доступ запрещен")
	errNotFound        = errors.New("Домен не найден")
	errMissingSiteID   = errors.New("Отсутствует ID сайта")
	errMissingDomainID = errors.New("Отсутствует ID домена")
)

type errorCode = string

const (
	errorCodeInvalidBody     errorCode = "INVALID_BODY"
	errorCodeMissingHostname errorCode = "MISSING_HOSTNAME"
	errorCodeInvalidHostname errorCode = "INVALID_HOSTNAME"
	errorCodeHostnameTaken   errorCode = "HOSTNAME_TAKEN"
	errorCodeMissingSiteID   errorCode = "MISSING_SITE_ID"
	errorCodeMissingDomainID errorCode = "MISSING_DOMAIN_ID"
	errorCodeAccessDenied    errorCode = "ACCESS_DENIED"
	errorCodeNotFound        errorCode = "NOT_FOUND"
	errorCodeOther           errorCode = "OTHER"
)

type handlerError struct {
	ErrorCode    errorCode `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
}

func (e handlerError) Error() string {
	return e.ErrorMessage
}
