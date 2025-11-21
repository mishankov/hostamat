package site

import "errors"

var (
	errLongName  = errors.New("long site name")
	errShortName = errors.New("short site name")
)
