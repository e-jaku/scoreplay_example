package repository

import (
	"context"
	"database/sql"
	"scoreplay/internal/domain"

	"github.com/lib/pq"
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
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	media := &domain.Media{
		Name:    name,
		FileURL: fileUrl,
	}

	insertMediaQuery := `INSERT INTO media (name, file_url) VALUES ($1, $2) RETURNING id`
	if err := r.db.QueryRowContext(ctx, insertMediaQuery, name, fileUrl).Scan(&media.ID); err != nil {
		return nil, err
	}

	insertMediaTagsQuery := `
        INSERT INTO media_tag (media_id, tag_id)
        SELECT $1, unnest($2::uuid[])
        ON CONFLICT DO NOTHING
    `
	if _, err := tx.ExecContext(ctx, insertMediaTagsQuery, media.ID, pq.Array(tags)); err != nil {
		return nil, err
	}

	return media, nil
}

func (r *PostgresMediaRepository) ListMediaByTagId(ctx context.Context, tagId string) ([]*domain.Media, error) {
	query := ` SELECT m.id AS media_id, m.name AS media_name, m.file_url, ARRAY_AGG(t.name) AS tags
	FROM media m JOIN media_tags mt ON m.id = mt.media_id
	JOIN tags t ON mt.tag_id = t.id
	WHERE m.id IN (	SELECT media_id FROM media_tags WHERE tag_id = $1)
	GROUP BY m.id Limit $2;`

	rows, err := r.db.QueryContext(ctx, query, tagId, 1000)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medias []*domain.Media
	for rows.Next() {
		media := &domain.Media{}
		err := rows.Scan(&media.ID, &media.Name, media.FileURL, media.Tags)
		if err != nil {
			return nil, err
		}
		medias = append(medias, media)
	}

	return medias, nil
}
