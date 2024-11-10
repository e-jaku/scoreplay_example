package service

import (
	"context"
	"scoreplay/internal/domain"
)

type MediaService interface {
	CreateMedia(ctx context.Context, name string, tags []string, media []byte) (*domain.Media, error)
	ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error)
}

type TagService interface {
	CreateTag(ctx context.Context, name string) (*domain.Tag, error)
	ListTags(ctx context.Context) ([]*domain.Tag, error)
}
