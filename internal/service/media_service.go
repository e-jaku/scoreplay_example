package service

import (
	"context"
	"io"
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

func (s *mediaServiceImpl) CreateMedia(ctx context.Context, name string, tags []string, file io.Reader, fileType string) (*domain.Media, error) {
	// we check that tags provided exist in db before starting uploading media
	dbTags, err := s.tagRepository.GetTags(ctx, tags)
	if err != nil {
		return nil, xerrors.Errorf("failed validating tags of media: %w", err)
	}

	fileURL, err := s.storage.UploadMedia(ctx, file, fileType)
	if err != nil {
		return nil, xerrors.Errorf("failed to upload media to storage")
	}

	createdMedia, err := s.mediaRepository.CreateMedia(ctx, name, tags, fileURL)
	if err != nil {
		return nil, xerrors.Errorf("failed to store media metadata")
	}

	var tagNames []string
	for _, tag := range dbTags {
		tagNames = append(tagNames, tag.Name)
	}

	createdMedia.Tags = tagNames

	return createdMedia, nil
}

func (s *mediaServiceImpl) ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error) {
	media, err := s.mediaRepository.ListMediaByTagId(ctx, tagId)
	if err != nil {
		return nil, xerrors.Errorf("failed to list media by tag id: %w", err)
	}
	return media, nil
}
