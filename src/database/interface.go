package database

import (
	"github.com/Ryan-Har/csnip/common/models"
	"github.com/google/uuid"
)

// Generic interface for interacting with databases.
type DatabaseInteractions interface {
	PopulateHelloWorldSnippets() error
	AddNewSnippet(m models.CodeSnippet) error
	UpdateSnippet(u uuid.UUID, changedSnippet models.CodeSnippet) (models.CodeSnippet, error)
	GetSnippets(page int64, limit int64) ([]models.CodeSnippet, error)
	GetSnippetsByLanguage(lang models.Language) ([]models.CodeSnippet, error)
	GetSnippetsByTag(tag string) ([]models.CodeSnippet, error)
	GetSnippetsByLanguageAndTag(lang models.Language, tag string) ([]models.CodeSnippet, error)
	GetSnippetByUUID(u uuid.UUID) (models.CodeSnippet, error)
	GetSnippetHistoryByUUID(u uuid.UUID) ([]models.CodeSnippet, error)
}
