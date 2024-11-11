package repository

import (
	"context"
	"database/sql"
	"scoreplay/internal/domain"

	"github.com/lib/pq"
	"golang.org/x/xerrors"
)

var _ TagRepository = (*PostgresTagRepository)(nil)

type PostgresTagRepository struct {
	db *sql.DB
}

func NewPostgresTagRepository(db *sql.DB) TagRepository {
	return &PostgresTagRepository{
		db: db,
	}
}

func (r *PostgresTagRepository) CreateTag(ctx context.Context, name string) (*domain.Tag, error) {
	tag := &domain.Tag{
		Name: name,
	}

	query := `INSERT INTO tag (name) VALUES ($1) RETURNING id`
	if err := r.db.QueryRowContext(ctx, query, name).Scan(&tag.ID); err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *PostgresTagRepository) ListTags(ctx context.Context) ([]*domain.Tag, error) {
	query := `SELECT id, name FROM tag LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, 1000)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		tag := &domain.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *PostgresTagRepository) GetTags(ctx context.Context, tagIds []string) ([]*domain.Tag, error) {
	query := `SELECT id, name FROM tag WHERE id = ANY($1::uuid[])`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(tagIds))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var tags []*domain.Tag
	for rows.Next() {
		tag := &domain.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if len(tags) != len(tagIds) {
		return nil, xerrors.New("some of the provided tags do not exist")
	}

	return tags, nil
}
