package models

import (
	"github.com/google/uuid"
	"time"
)

type CodeSnippet struct {
	ID           int64
	Uuid         uuid.UUID
	Name         string
	Code         string
	Language     string
	Tags         string
	Description  string
	Source       string
	DateAdded    time.Time
	Version      int64
	SupersededBy int64
}
