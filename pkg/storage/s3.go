package storage

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FileStorage struct {
	uploader *s3manager.Uploader
	bucket   string
}

func NewFileStorage(uploader *s3manager.Uploader, bucket string) *FileStorage {
	return &FileStorage{
		uploader: uploader,
		bucket:   bucket,
	}
}

func (fs *FileStorage) Upload(ctx context.Context, input UploadInput) (string, error) {
	res, err := fs.uploader.Upload(&s3manager.UploadInput{
		Bucket: &fs.bucket,
		Key:    &input.Name,
		Body:   input.File,
	})

	if err != nil {
		return "", err
	}

	return res.Location, nil
}
