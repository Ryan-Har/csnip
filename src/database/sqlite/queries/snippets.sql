-- name: CreateSnippet :one
-- Creates the first version of a snippet
INSERT INTO snippets (
    uuid, name, code, language, tags, description, source, date_added, version, superseded_by
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?, NULL
) RETURNING *;

-- name: GetSnippetByID :one
-- Get a snippet by its ID
SELECT * FROM snippets WHERE id = ?;

-- name: GetSnippetByUUID :one
-- Get last version of a snippet by UUID
SELECT * FROM snippets WHERE uuid = ? ORDER BY version DESC LIMIT 1;

-- name: GetSnippetVersions :many
-- Get all versions of a snippet by UUID
SELECT * FROM snippets WHERE uuid = ? ORDER BY version DESC;

-- name: ListSnippetsByPage :many
-- Get all latest snippets, paginated
SELECT * FROM snippets
WHERE superseded_by IS NULL
ORDER BY id DESC
LIMIT :limit OFFSET :offset;

-- name: GetSnippetByLanguage :many
-- Get last version of a snippets by language
SELECT * FROM snippets WHERE language = ? 
AND superseded_by IS NULL
ORDER BY id DESC;

-- name: GetSnippetByTag :many
-- Get last version of a snippets by tag
SELECT * FROM snippets WHERE instr(tags, ?) > 0 
AND superseded_by IS NULL
ORDER BY id DESC;

-- name: GetSnippetBySource :many
-- Get last version of a snippets by source
SELECT * FROM snippets WHERE instr(source, ?) > 0 
AND superseded_by IS NULL
ORDER BY id DESC;


-- name: GetSnippetByLanguageAndTag :many
-- Get last version of a snippets by language and tag
SELECT * FROM snippets WHERE language = ? 
AND instr(tags, ?) > 0 
AND superseded_by IS NULL
ORDER BY id DESC;

-- name: MarkSnippetSuperseded :exec
-- Marks an old snippet as superseded by a new version
UPDATE snippets
SET superseded_by = ?
WHERE id = ?;

-- name: DeleteSnippet :exec
-- Delete a snippet by its ID
DELETE FROM snippets WHERE id = ?;

-- name: DeleteSnippetByUUID :exec
-- Delete all versions of a snippet by UUID
DELETE FROM snippets WHERE uuid = ?;