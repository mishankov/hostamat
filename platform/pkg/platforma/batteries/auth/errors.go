package auth

import "errors"

var ErrInvalidPassword = errors.New("invalid password")
var ErrWrongUserOrPassword = errors.New("wrong user or password")
var ErrCurrentPasswordIncorrect = errors.New("current password is incorrect")
