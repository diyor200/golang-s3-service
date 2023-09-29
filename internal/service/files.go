package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang-diplom-work/internal/entity"
	"golang-diplom-work/pkg/storage"
	"os"
	"strings"
)

type FileService struct {
	provider storage.Provider
	file     entity.File
}

func NewFileService(provider storage.Provider) *FileService {
	return &FileService{provider: provider}
}

func (fs *FileService) Upload(ctx context.Context, file entity.File) (string, error) {
	f, err := os.Open(file.Name)
	if err != nil {
		return "", err
	}
	info, _ := f.Stat()
	log.Info("file info: %v", info)
	defer f.Close()

	if fs.checkFile(f.Name(), info.Size()) {
		return fs.provider.Upload(ctx, storage.UploadInput{
			Name: file.Name,
			Size: info.Size(),
			File: f,
		})
	}
	return "", ErrCannotUploadFile
}

func (fs *FileService) checkFile(file string, size int64) bool {
	extensions := [3]string{"jpeg", "jpg", "png"}
	recommendedSize := int64(1024 * 1024 * 8)
	fileParts := strings.Split(file, ".")
	ext := fileParts[len(fileParts)-1]

	if recommendedSize < size {
		return false
	}

	for _, extension := range extensions {
		if ext == extension {
			return true
		}
	}
	return false
}
