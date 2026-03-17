package file

import "errors"

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileNotFound    = errors.New("file not found")
)

type File struct {
	Name string
	Size int
}
