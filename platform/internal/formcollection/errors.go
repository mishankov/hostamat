package formcollection

import "errors"

var (
	errAccessDenied          = errors.New("access denied")
	errNoFormsSelected       = errors.New("no forms selected")
	errFormNotSupported      = errors.New("form is not supported")
	errFormNotFound          = errors.New("form not found")
	errFormInactive          = errors.New("form is inactive")
	errTooManyFields         = errors.New("too many fields")
	errFieldValueTooLong     = errors.New("field value too long")
	errUnsupportedSubmission = errors.New("unsupported submission encoding")
)
