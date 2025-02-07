// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: snippets.sql

package sqlite

import (
	"context"
	"database/sql"
)

const createSnippet = `-- name: CreateSnippet :one
INSERT INTO snippets (
    uuid, name, code, language, tags, description, source, date_added, version, superseded_by
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?, NULL
) RETURNING id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by
`

type CreateSnippetParams struct {
	Uuid        string
	Name        sql.NullString
	Code        string
	Language    string
	Tags        sql.NullString
	Description sql.NullString
	Source      sql.NullString
	Version     int64
}

// Creates the first version of a snippet
func (q *Queries) CreateSnippet(ctx context.Context, arg CreateSnippetParams) (Snippet, error) {
	row := q.db.QueryRowContext(ctx, createSnippet,
		arg.Uuid,
		arg.Name,
		arg.Code,
		arg.Language,
		arg.Tags,
		arg.Description,
		arg.Source,
		arg.Version,
	)
	var i Snippet
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.Code,
		&i.Language,
		&i.Tags,
		&i.Description,
		&i.Source,
		&i.DateAdded,
		&i.Version,
		&i.SupersededBy,
	)
	return i, err
}

const deleteSnippet = `-- name: DeleteSnippet :exec
DELETE FROM snippets WHERE id = ?
`

// Delete a snippet by its ID
func (q *Queries) DeleteSnippet(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSnippet, id)
	return err
}

const deleteSnippetByUUID = `-- name: DeleteSnippetByUUID :exec
DELETE FROM snippets WHERE uuid = ?
`

// Delete all versions of a snippet by UUID
func (q *Queries) DeleteSnippetByUUID(ctx context.Context, uuid string) error {
	_, err := q.db.ExecContext(ctx, deleteSnippetByUUID, uuid)
	return err
}

const getSnippetByID = `-- name: GetSnippetByID :one
SELECT id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by FROM snippets WHERE id = ?
`

// Get a snippet by its ID
func (q *Queries) GetSnippetByID(ctx context.Context, id int64) (Snippet, error) {
	row := q.db.QueryRowContext(ctx, getSnippetByID, id)
	var i Snippet
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.Code,
		&i.Language,
		&i.Tags,
		&i.Description,
		&i.Source,
		&i.DateAdded,
		&i.Version,
		&i.SupersededBy,
	)
	return i, err
}

const getSnippetByLanguage = `-- name: GetSnippetByLanguage :many
SELECT id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by FROM snippets WHERE language = ? 
AND superseded_by IS NULL
ORDER BY id DESC
`

// Get last version of a snippets by language
func (q *Queries) GetSnippetByLanguage(ctx context.Context, language string) ([]Snippet, error) {
	rows, err := q.db.QueryContext(ctx, getSnippetByLanguage, language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Snippet
	for rows.Next() {
		var i Snippet
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.Code,
			&i.Language,
			&i.Tags,
			&i.Description,
			&i.Source,
			&i.DateAdded,
			&i.Version,
			&i.SupersededBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSnippetByLanguageAndTag = `-- name: GetSnippetByLanguageAndTag :many
SELECT id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by FROM snippets WHERE language = ? 
AND instr(tags, ?) > 0 
AND superseded_by IS NULL
ORDER BY id DESC
`

type GetSnippetByLanguageAndTagParams struct {
	Language string
	INSTR    string
}

// Get last version of a snippets by language and tag
func (q *Queries) GetSnippetByLanguageAndTag(ctx context.Context, arg GetSnippetByLanguageAndTagParams) ([]Snippet, error) {
	rows, err := q.db.QueryContext(ctx, getSnippetByLanguageAndTag, arg.Language, arg.INSTR)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Snippet
	for rows.Next() {
		var i Snippet
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.Code,
			&i.Language,
			&i.Tags,
			&i.Description,
			&i.Source,
			&i.DateAdded,
			&i.Version,
			&i.SupersededBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSnippetBySource = `-- name: GetSnippetBySource :many
SELECT id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by FROM snippets WHERE instr(source, ?) > 0 
AND superseded_by IS NULL
ORDER BY id DESC
`

// Get last version of a snippets by source
func (q *Queries) GetSnippetBySource(ctx context.Context, instr string) ([]Snippet, error) {
	rows, err := q.db.QueryContext(ctx, getSnippetBySource, instr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Snippet
	for rows.Next() {
		var i Snippet
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.Code,
			&i.Language,
			&i.Tags,
			&i.Description,
			&i.Source,
			&i.DateAdded,
			&i.Version,
			&i.SupersededBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSnippetByTag = `-- name: GetSnippetByTag :many
SELECT id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by FROM snippets WHERE instr(tags, ?) > 0 
AND superseded_by IS NULL
ORDER BY id DESC
`

// Get last version of a snippets by tag
func (q *Queries) GetSnippetByTag(ctx context.Context, instr string) ([]Snippet, error) {
	rows, err := q.db.QueryContext(ctx, getSnippetByTag, instr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Snippet
	for rows.Next() {
		var i Snippet
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.Code,
			&i.Language,
			&i.Tags,
			&i.Description,
			&i.Source,
			&i.DateAdded,
			&i.Version,
			&i.SupersededBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSnippetByUUID = `-- name: GetSnippetByUUID :one
SELECT id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by FROM snippets WHERE uuid = ? ORDER BY version DESC LIMIT 1
`

// Get last version of a snippet by UUID
func (q *Queries) GetSnippetByUUID(ctx context.Context, uuid string) (Snippet, error) {
	row := q.db.QueryRowContext(ctx, getSnippetByUUID, uuid)
	var i Snippet
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.Code,
		&i.Language,
		&i.Tags,
		&i.Description,
		&i.Source,
		&i.DateAdded,
		&i.Version,
		&i.SupersededBy,
	)
	return i, err
}

const getSnippetVersions = `-- name: GetSnippetVersions :many
SELECT id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by FROM snippets WHERE uuid = ? ORDER BY version DESC
`

// Get all versions of a snippet by UUID
func (q *Queries) GetSnippetVersions(ctx context.Context, uuid string) ([]Snippet, error) {
	rows, err := q.db.QueryContext(ctx, getSnippetVersions, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Snippet
	for rows.Next() {
		var i Snippet
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.Code,
			&i.Language,
			&i.Tags,
			&i.Description,
			&i.Source,
			&i.DateAdded,
			&i.Version,
			&i.SupersededBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSnippetsByPage = `-- name: ListSnippetsByPage :many
SELECT id, uuid, name, code, language, tags, description, source, date_added, version, superseded_by FROM snippets
WHERE superseded_by IS NULL
ORDER BY id DESC
LIMIT ?2 OFFSET ?1
`

type ListSnippetsByPageParams struct {
	Offset int64
	Limit  int64
}

// Get all latest snippets, paginated
func (q *Queries) ListSnippetsByPage(ctx context.Context, arg ListSnippetsByPageParams) ([]Snippet, error) {
	rows, err := q.db.QueryContext(ctx, listSnippetsByPage, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Snippet
	for rows.Next() {
		var i Snippet
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.Code,
			&i.Language,
			&i.Tags,
			&i.Description,
			&i.Source,
			&i.DateAdded,
			&i.Version,
			&i.SupersededBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markSnippetSuperseded = `-- name: MarkSnippetSuperseded :exec
UPDATE snippets
SET superseded_by = ?
WHERE id = ?
`

type MarkSnippetSupersededParams struct {
	SupersededBy sql.NullInt64
	ID           int64
}

// Marks an old snippet as superseded by a new version
func (q *Queries) MarkSnippetSuperseded(ctx context.Context, arg MarkSnippetSupersededParams) error {
	_, err := q.db.ExecContext(ctx, markSnippetSuperseded, arg.SupersededBy, arg.ID)
	return err
}
