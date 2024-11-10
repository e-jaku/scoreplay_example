package domain

import "github.com/google/uuid"

type Tag struct {
	ID   uuid.UUID `db:"id" json:"id,omitempty"`
	Name string    `db:"name" json:"name"`
}
