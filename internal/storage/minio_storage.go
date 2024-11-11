package storage

import (
	"context"
	"fmt"
	"io"
	"scoreplay/internal/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/xerrors"
)

var _ Storage = (*MinioStorage)(nil)

type MinioStorage struct {
	Client *minio.Client
	Bucket string
}

func NewMinioStorage(cfg *config.StorageConfig) (Storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretKeyID, ""),
		Secure: false, // just for sample minio running without HTTPS
	})
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize MinIO client: %v", err)
	}

	return &MinioStorage{
		Client: client,
		Bucket: cfg.Bucket,
	}, nil
}

func (s *MinioStorage) UploadMedia(ctx context.Context, file io.Reader, fileType string) (string, error) {
	fileName := uuid.New()
	objectName := fmt.Sprintf("%s%s", fileName, fileType)

	_, err := s.Client.PutObject(ctx, s.Bucket, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", xerrors.Errorf("failed to upload file to MinIO: %v", err)
	}

	return fmt.Sprintf("http://localhost:9000/%s/%s", s.Bucket, objectName), nil
}
