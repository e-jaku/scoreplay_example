package storage

import (
	"context"
	"io"
)

type Storage interface {
	UploadMedia(ctx context.Context, file io.Reader, contentType string) (string, error)
}
