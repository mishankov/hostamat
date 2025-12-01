package file

import "errors"

var (
	ErrInvalidFileType = errors.New("invalid file type")
)

type File struct {
	Name string
	Size int
}
