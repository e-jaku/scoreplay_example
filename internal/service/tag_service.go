package service

import (
	"context"
	"scoreplay/internal/domain"
	"scoreplay/internal/repository"

	"golang.org/x/xerrors"
)

var _ TagService = (*tagServiceImpl)(nil)

type tagServiceImpl struct {
	repository repository.TagRepository
}

func NewTagService(repository repository.TagRepository) TagService {
	return &tagServiceImpl{
		repository: repository,
	}
}

func (s *tagServiceImpl) CreateTag(ctx context.Context, name string) (*domain.Tag, error) {
	tag, err := s.repository.CreateTag(ctx, name)
	if err != nil {
		return nil, xerrors.Errorf("failed to create tag: %w", err)
	}
	return tag, nil
}

func (s *tagServiceImpl) ListTags(ctx context.Context) ([]*domain.Tag, error) {
	tags, err := s.repository.ListTags(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to list tags: %w", err)
	}
	return tags, nil
}
