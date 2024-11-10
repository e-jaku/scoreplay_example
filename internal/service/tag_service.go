package service

import (
	"context"
	"scoreplay/internal/domain"
	"scoreplay/internal/repository"
)

var _ TagService = (*tagServiceImpl)(nil)

type TagService interface {
	CreateTag(ctx context.Context, name string) (*domain.Tag, error)
	ListTags(ctx context.Context) ([]*domain.Tag, error)
}

type tagServiceImpl struct {
	repository repository.TagRepository
}

func NewTagService(repository repository.TagRepository) TagService {
	return &tagServiceImpl{
		repository: repository,
	}
}

func (s *tagServiceImpl) CreateTag(ctx context.Context, name string) (*domain.Tag, error) {
	return &domain.Tag{}, nil
}

func (s *tagServiceImpl) ListTags(ctx context.Context) ([]*domain.Tag, error) {
	return nil, nil
}
