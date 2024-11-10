package repository

import (
	"context"
	"database/sql"
	"scoreplay/internal/domain"
)

var _ MediaRepository = (*PostgresMediaRepository)(nil)

type PostgresMediaRepository struct {
	db *sql.DB
}

func NewPostgresMediaRepository(db *sql.DB) MediaRepository {
	return &PostgresMediaRepository{
		db: db,
	}
}

func (r *PostgresMediaRepository) CreateMedia(ctx context.Context, name string, tags []string, fileUrl string) (*domain.Media, error) {
	return nil, nil
}

func (r *PostgresMediaRepository) ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error) {
	return nil, nil
}
