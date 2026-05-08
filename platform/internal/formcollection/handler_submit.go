package formcollection

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type submitService interface {
	Submit(ctx context.Context, token string, values url.Values) (*Submission, error)
}

type SubmitHandler struct {
	service submitService
}

func NewSubmitHandler(service submitService) *SubmitHandler {
	return &SubmitHandler{service: service}
}

func (h *SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "multipart/form-data") {
		respondSubmissionPage(w, http.StatusBadRequest, "Не удалось отправить форму", "Формы с файлами пока не поддерживаются.")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxSubmissionBodyBytes)
	if err := r.ParseForm(); err != nil {
		respondSubmissionPage(w, http.StatusBadRequest, "Не удалось отправить форму", "Проверьте поля формы и попробуйте еще раз.")
		return
	}

	_, err := h.service.Submit(r.Context(), r.PathValue("token"), r.PostForm)
	if err != nil {
		switch {
		case errors.Is(err, errFormNotFound), errors.Is(err, errFormInactive):
			respondSubmissionPage(w, http.StatusNotFound, "Форма недоступна", "Эта форма больше не принимает ответы.")
		case errors.Is(err, errTooManyFields):
			respondSubmissionPage(w, http.StatusBadRequest, "Не удалось отправить форму", "В форме слишком много полей.")
		case errors.Is(err, errFieldValueTooLong):
			respondSubmissionPage(w, http.StatusBadRequest, "Не удалось отправить форму", "Одно из полей слишком длинное.")
		case errors.Is(err, errUnsupportedSubmission):
			respondSubmissionPage(w, http.StatusBadRequest, "Не удалось отправить форму", "Этот формат отправки пока не поддерживается.")
		default:
			respondSubmissionPage(w, http.StatusInternalServerError, "Не удалось отправить форму", "Попробуйте повторить отправку позже.")
		}
		return
	}

	respondSubmissionPage(w, http.StatusOK, "Форма отправлена", "Спасибо. Ваш ответ сохранен.")
}
