package domain

import (
	"github.com/google/uuid"
)

type Media struct {
	ID      uuid.UUID `db:"id" json:"id,omitempty"`
	Name    string    `db:"name" json:"name"`
	Tags    []string  `db:"tags" json:"tags"`
	FileURL string    `db:"file_url" json:"fileUrl"`
}
