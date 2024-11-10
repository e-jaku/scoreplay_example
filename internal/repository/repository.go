package repository

import (
	"context"
	"scoreplay/internal/domain"
)

type TagRepository interface {
	CreateTag(ctx context.Context, name string) (*domain.Tag, error)
	ListTags(ctx context.Context) ([]*domain.Tag, error)
}

type MediaRepository interface {
	CreateMedia(ctx context.Context, name string, tags []string, fileUrl string) (*domain.Media, error)
	ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error)
}

type Repository interface {
	TagRepository
	MediaRepository
}
