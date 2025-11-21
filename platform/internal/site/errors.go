package site

import "errors"

var (
	errLongName  = errors.New("Название сайта должно быть не более 25 символов")
	errShortName = errors.New("Название сайта должно быть не менее 1 символа")
)
