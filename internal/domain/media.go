package domain

import (
	"github.com/google/uuid"
)

type Media struct {
	ID      uuid.UUID
	Name    string
	Tags    []string
	FileURL string
}
