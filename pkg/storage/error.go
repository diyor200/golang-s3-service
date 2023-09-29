package storage

import (
	"fmt"
)

var (
	ErrFileNotValid     = fmt.Errorf("file should be photo")
	ErrFileBigSize      = fmt.Errorf("file size should not exceed 8MB")
	ErrCannotUploadFile = fmt.Errorf("an error occured while uploading file")
)
