package service

import (
	"context"
	"io"
	"scoreplay/internal/domain"
)

type MediaService interface {
	CreateMedia(ctx context.Context, name string, tags []string, file io.Reader, fileType string) (*domain.Media, error)
	ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error)
}

type TagService interface {
	CreateTag(ctx context.Context, name string) (*domain.Tag, error)
	ListTags(ctx context.Context) ([]*domain.Tag, error)
}
