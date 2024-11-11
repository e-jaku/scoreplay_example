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
	mediaRepository repository.MediaRepository
	tagRepository   repository.TagRepository
	storage         storage.Storage
}

func NewMediaService(mediaRepository repository.MediaRepository, tagRepository repository.TagRepository, storage storage.Storage) MediaService {
	return &mediaServiceImpl{
		mediaRepository: mediaRepository,
		tagRepository:   tagRepository,
		storage:         storage,
	}
}

func (s *mediaServiceImpl) CreateMedia(ctx context.Context, name string, tags []string, media []byte) (*domain.Media, error) {
	// after the upload create file url string

	// validate tags exist

	// vlaidate file???

	// upload file and get url

	// store in create media

	s.mediaRepository.CreateMedia(ctx, name, tags, "tbd")

	return &domain.Media{}, nil
}

func (s *mediaServiceImpl) ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error) {
	media, err := s.mediaRepository.ListMediaByTagId(ctx, tagId)
	if err != nil {
		return nil, xerrors.Errorf("failed to list media by tag id")
	}
	return media, nil
}
