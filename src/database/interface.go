package database

import (
	"errors"
	"github.com/Ryan-Har/csnip/common/models"
	"github.com/google/uuid"
)

// Generic interface for interacting with databases.
type DatabaseInteractions interface {
	PopulateHelloWorldSnippets() error
	AddNewSnippet(m models.CodeSnippet) error
	UpdateSnippet(u uuid.UUID, changedSnippet models.CodeSnippet) (models.CodeSnippet, error)
	GetSnippets(page int64, limit int64) ([]models.CodeSnippet, error)
	GetSnippetsByLanguage(lang string) ([]models.CodeSnippet, error)
	GetSnippetsByTag(tag string) ([]models.CodeSnippet, error)
	GetSnippetsByLanguageAndTag(lang string, tag string) ([]models.CodeSnippet, error)
	GetSnippetByUUID(u uuid.UUID) (models.CodeSnippet, error)
	GetSnippetHistoryByUUID(u uuid.UUID) ([]models.CodeSnippet, error)
	DeleteSnippetByUUID(u uuid.UUID) error
}

// custom errors used by the above interface, used when no results are found in the sql results set
var ErrNoSnippetsFound = errors.New("no snippets found for the given parameters")
