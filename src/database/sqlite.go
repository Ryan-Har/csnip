package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/Ryan-Har/csnip/common"
	"github.com/Ryan-Har/csnip/common/models"
	"github.com/Ryan-Har/csnip/database/sqlite"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteHandler struct {
	database   *sql.DB
	queries    *sqlite.Queries
	version    int32 //version of database schema
	writeMutex *sync.Mutex
}

func NewSQLiteHandler() (DatabaseInteractions, error) {
	var dbHandler DatabaseInteractions
	dbLoc := "./my.db"
	db, err := openSQLiteDB(dbLoc)
	if err != nil {
		return dbHandler, err
	}

	dbHandler = &SQLiteHandler{
		database:   db,
		queries:    sqlite.New(db),
		version:    1,
		writeMutex: &sync.Mutex{},
	}
	return dbHandler, nil
}

func openSQLiteDB(dbLoc string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbLoc)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		// Handle error setting foreign_keys
		return nil, fmt.Errorf("unable to enforce foreign keys in db: %v", err.Error())
	}
	return db, nil
}

func (s SQLiteHandler) PopulateHelloWorldSnippets() error {
	examples := common.GetHelloWorldExamples()
	for lang, code := range examples {
		err := s.AddNewSnippet(models.CodeSnippet{
			Name:        "Hello World Example in " + lang,
			Code:        code,
			Language:    lang,
			Tags:        "example,generated",
			Description: "Hello World Example in " + lang,
			Source:      "generated",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s SQLiteHandler) AddNewSnippet(m models.CodeSnippet) error {
	m.Uuid = uuid.New()
	m.Version = 1

	createParams := codeSnippetModelToDbCreateSnippetParams(m)
	_, err := s.queries.CreateSnippet(context.Background(), createParams)
	return err
}

// updates the uuid with the changedSnippet
func (s SQLiteHandler) UpdateSnippet(u uuid.UUID, changedSnippet models.CodeSnippet) (models.CodeSnippet, error) {
	//initialise return snippet
	var returnSnippet models.CodeSnippet

	//get last change as oldSnippet
	oldSnippet, err := s.queries.GetSnippetByUUID(context.Background(), u.String())
	if err != nil {
		return returnSnippet, err
	}

	oldCodeSnippet := convertSqliteSnippetToCodeSnippet(oldSnippet)
	snippetToUpdate := normaliseCodeSnippetStruct(changedSnippet, oldCodeSnippet)
	createParams := codeSnippetModelToDbCreateSnippetParams(snippetToUpdate)

	//begin transaction
	tx, err := s.database.BeginTx(context.Background(), nil)
	if err != nil {
		return returnSnippet, fmt.Errorf("failed to start transaction: %w", err)
	}

	q := s.queries.WithTx(tx)

	createdSnippet, err := q.CreateSnippet(context.Background(), createParams)
	if err != nil {
		tx.Rollback()
		return returnSnippet, fmt.Errorf("failed to insert snippet: %w", err)
	}

	supersededParams := sqlite.MarkSnippetSupersededParams{
		SupersededBy: sql.NullInt64{
			Int64: createdSnippet.ID,
			Valid: true,
		},
		ID: oldSnippet.ID,
	}

	err = q.MarkSnippetSuperseded(context.Background(), supersededParams)
	if err != nil {
		tx.Rollback()
		return returnSnippet, fmt.Errorf("failed to mark snippet superceded: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return returnSnippet, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return convertSqliteSnippetToCodeSnippet(createdSnippet), nil

}

// GetSnippets returns a list of snippets.
// item is paginated for efficiency, inputs are the page number needed and the limit for response.
func (s SQLiteHandler) GetSnippets(page int64, limit int64) ([]models.CodeSnippet, error) {
	var responseSnippets []models.CodeSnippet
	offset := (page - 1) * limit

	pageParams := sqlite.ListSnippetsByPageParams{
		Offset: offset,
		Limit:  limit,
	}

	dbSnippets, err := s.queries.ListSnippetsByPage(context.Background(), pageParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responseSnippets, ErrNoSnippetsFound
		}
		return responseSnippets, fmt.Errorf("failed to retrieve snippets: %w", err)
	}

	for _, snippet := range dbSnippets {
		responseSnippets = append(responseSnippets, convertSqliteSnippetToCodeSnippet(snippet))
	}

	return responseSnippets, nil
}

// GetSnippetsByLanguage returns a list of snippets
func (s SQLiteHandler) GetSnippetsByLanguage(lang string) ([]models.CodeSnippet, error) {
	var responseSnippets []models.CodeSnippet

	dbSnippets, err := s.queries.GetSnippetByLanguage(context.Background(), lang)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responseSnippets, ErrNoSnippetsFound
		}
		return responseSnippets, fmt.Errorf("failed to retrieve snippets: %w", err)
	}

	for _, snippet := range dbSnippets {
		responseSnippets = append(responseSnippets, convertSqliteSnippetToCodeSnippet(snippet))
	}

	return responseSnippets, nil
}

// GetSnippetsByTag returns a list of snippets where the tag string provided patially matches the list of tags in the database
func (s SQLiteHandler) GetSnippetsByTag(tag string) ([]models.CodeSnippet, error) {
	var responseSnippets []models.CodeSnippet

	dbSnippets, err := s.queries.GetSnippetByTag(context.Background(), tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responseSnippets, ErrNoSnippetsFound
		}
		return responseSnippets, fmt.Errorf("failed to retrieve snippets: %w", err)
	}

	for _, snippet := range dbSnippets {
		responseSnippets = append(responseSnippets, convertSqliteSnippetToCodeSnippet(snippet))
	}

	return responseSnippets, nil
}

// GetSnippetsByLanguageAndTag returns a list of snippets where language matches and the tag string provided patially matches the list of tags in the database
func (s SQLiteHandler) GetSnippetsByLanguageAndTag(lang string, tag string) ([]models.CodeSnippet, error) {
	var responseSnippets []models.CodeSnippet

	params := sqlite.GetSnippetByLanguageAndTagParams{
		Language: lang,
		INSTR:    tag,
	}

	dbSnippets, err := s.queries.GetSnippetByLanguageAndTag(context.Background(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responseSnippets, ErrNoSnippetsFound
		}
		return responseSnippets, fmt.Errorf("failed to retrieve snippets: %w", err)
	}

	for _, snippet := range dbSnippets {
		responseSnippets = append(responseSnippets, convertSqliteSnippetToCodeSnippet(snippet))
	}

	return responseSnippets, nil
}

// GetSnippetsByUUID returns a single snippet matching the UUID
func (s SQLiteHandler) GetSnippetByUUID(u uuid.UUID) (models.CodeSnippet, error) {
	dbSnippet, err := s.queries.GetSnippetByUUID(context.Background(), u.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CodeSnippet{}, ErrNoSnippetsFound
		}
		return models.CodeSnippet{}, fmt.Errorf("failed to retrieve snippets: %w", err)
	}

	return convertSqliteSnippetToCodeSnippet(dbSnippet), nil
}

// GetSnippetHistoryByUUID returns a the snippet history
func (s SQLiteHandler) GetSnippetHistoryByUUID(u uuid.UUID) ([]models.CodeSnippet, error) {
	var responseSnippets []models.CodeSnippet

	dbSnippets, err := s.queries.GetSnippetVersions(context.Background(), u.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responseSnippets, ErrNoSnippetsFound
		}
		return responseSnippets, fmt.Errorf("failed to retrieve snippets: %w", err)
	}

	for _, snippet := range dbSnippets {
		responseSnippets = append(responseSnippets, convertSqliteSnippetToCodeSnippet(snippet))
	}

	return responseSnippets, nil
}

// DeleteSnippetByUUID returns a single snippet matching the UUID
func (s SQLiteHandler) DeleteSnippetByUUID(u uuid.UUID) error {
	err := s.queries.DeleteSnippetByUUID(context.Background(), u.String())
	if err != nil {
		return fmt.Errorf("failed to delete snippet by uuid: %w", err)
	}
	return nil
}

// Convert Snippet model -> Snippet Create Params (DB)
func codeSnippetModelToDbCreateSnippetParams(m models.CodeSnippet) sqlite.CreateSnippetParams {
	return sqlite.CreateSnippetParams{
		Uuid:        m.Uuid.String(),
		Name:        toNullString(m.Name),
		Code:        m.Code,
		Language:    m.Language,
		Tags:        toNullString(m.Tags),
		Description: toNullString(m.Description),
		Source:      toNullString(m.Source),
		Version:     m.Version,
	}
}

// Convert Snippet (DB) -> Snippet Create Params (DB)
func sqliteSnippetToDbCreateSnippetParams(s sqlite.Snippet) sqlite.CreateSnippetParams {
	return sqlite.CreateSnippetParams{
		Uuid:        s.Uuid,
		Name:        s.Name,
		Code:        s.Code,
		Language:    s.Language,
		Tags:        s.Tags,
		Description: s.Description,
		Source:      s.Source,
		Version:     s.Version,
	}
}

// Convert Snippet (DB) -> CodeSnippet model
func convertSqliteSnippetToCodeSnippet(s sqlite.Snippet) models.CodeSnippet {
	// Parse UUID
	parsedUUID, err := uuid.Parse(s.Uuid)
	if err != nil {
		parsedUUID = uuid.Nil // Defaults to empty UUID if parsing fails
	}

	return models.CodeSnippet{
		ID:           s.ID,
		Uuid:         parsedUUID,
		Name:         getString(s.Name),
		Code:         s.Code,
		Language:     s.Language,
		Tags:         getString(s.Tags),
		Description:  getString(s.Description),
		Source:       getString(s.Source),
		DateAdded:    getTime(s.DateAdded),
		Version:      s.Version,
		SupersededBy: getInt64(s.SupersededBy),
	}
}

// compares code snippets to ensure that any missing fields are retained from the old snippet
func normaliseCodeSnippetStruct(toUpdate models.CodeSnippet, old models.CodeSnippet) models.CodeSnippet {
	if toUpdate.Uuid != old.Uuid {
		toUpdate.Uuid = old.Uuid
	}
	if toUpdate.Name == "" {
		toUpdate.Name = old.Name
	}
	if toUpdate.Code == "" {
		toUpdate.Code = old.Code
	}
	if toUpdate.Language != old.Language {
		toUpdate.Language = old.Language
	}
	if toUpdate.Tags == "" {
		toUpdate.Tags = old.Tags
	}
	if toUpdate.Description == "" {
		toUpdate.Description = old.Description
	}
	if toUpdate.Source == "" {
		toUpdate.Source = old.Source
	}
	toUpdate.Version = old.Version + 1

	return toUpdate
}
