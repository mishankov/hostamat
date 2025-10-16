package user

import "errors"

var ErrInvalidPassword = errors.New("invalid password")
var ErrWrongUserOrPassword = errors.New("wrong user or password")
