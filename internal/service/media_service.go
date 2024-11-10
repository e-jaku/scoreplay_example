package service

import (
	"context"
	"scoreplay/internal/domain"
	"scoreplay/internal/repository"
	"scoreplay/internal/storage"
)

var _ MediaService = (*mediaServiceImpl)(nil)

type MediaService interface {
	CreateMedia(ctx context.Context, name string, tags []string) (*domain.Media, error)
	ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error)
}

type mediaServiceImpl struct {
	repository repository.TagRepository
	storage    storage.Storage
}

func NewMediaService(repository repository.TagRepository, storage storage.Storage) MediaService {
	return &mediaServiceImpl{
		repository: repository,
		storage:    storage,
	}
}

func (s *mediaServiceImpl) CreateMedia(ctx context.Context, name string, tags []string) (*domain.Media, error) {
	// after the upload create file url string
	return &domain.Media{}, nil
}

func (s *mediaServiceImpl) ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error) {
	return nil, nil
}
