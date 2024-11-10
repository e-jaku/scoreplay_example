package service

import (
	"context"
	"scoreplay/internal/domain"
	"scoreplay/internal/repository"
	"scoreplay/internal/storage"

	"golang.org/x/xerrors"
)

var _ MediaService = (*mediaServiceImpl)(nil)

type mediaServiceImpl struct {
	repository repository.MediaRepository
	storage    storage.Storage
}

func NewMediaService(repository repository.MediaRepository, storage storage.Storage) MediaService {
	return &mediaServiceImpl{
		repository: repository,
		storage:    storage,
	}
}

func (s *mediaServiceImpl) CreateMedia(ctx context.Context, name string, tags []string, media []byte) (*domain.Media, error) {
	// after the upload create file url string

	// validate tags exist

	// vlaidate file???

	// upload file and get url

	// store in create media

	s.repository.CreateMedia(ctx, name, tags, "tbd")

	return &domain.Media{}, nil
}

func (s *mediaServiceImpl) ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error) {
	media, err := s.repository.ListMediaByTagId(ctx, tagId)
	if err != nil {
		return nil, xerrors.Errorf("failed to list media by tag id")
	}
	return media, nil
}
