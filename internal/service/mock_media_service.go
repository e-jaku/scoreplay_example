package service

import (
	"context"
	"io"
	"scoreplay/internal/domain"

	"github.com/google/uuid"
)

var _ MediaService = (*MockedMediaService)(nil)

type MockedMediaService struct {
	Err error
}

func (s *MockedMediaService) CreateMedia(ctx context.Context, name string, tags []string, file io.Reader, fileType string) (*domain.Media, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	return &domain.Media{
		ID:      uuid.New(),
		Name:    name,
		FileURL: fileType, // for simplicity
	}, nil
}

func (s *MockedMediaService) ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	return []*domain.Media{}, nil
}
