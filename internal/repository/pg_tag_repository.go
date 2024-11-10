package repository

import (
	"context"
	"database/sql"
	"scoreplay/internal/domain"
)

type PostgresTagRepository struct {
	db *sql.DB
}

func NewPostgresTagRepository(db *sql.DB) TagRepository {
	return &PostgresTagRepository{
		db: db,
	}
}

func (r *PostgresTagRepository) CreateTag(ctx context.Context, name string) (*domain.Tag, error) {
	return nil, nil
}

func (r *PostgresTagRepository) ListTags(ctx context.Context) ([]*domain.Tag, error) {
	return nil, nil
}
