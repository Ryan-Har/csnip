-- name: CreateGroup :one
-- Creates a new group
INSERT INTO groups (
    group_name, description, uuid_list, date_added, date_updated
) VALUES (
    ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetGroupByID :one
-- Get a group by its ID
SELECT * FROM groups WHERE id = ?;

-- name: GetAllGroups :many
-- Get all groups
SELECT * FROM groups ORDER BY date_added DESC;

-- name: UpdateGroup :exec
-- Updates a group's details
UPDATE groups
SET group_name = ?, description = ?, uuid_list = ?, date_updated = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteGroup :exec
-- Delete a group by its ID
DELETE FROM groups WHERE id = ?;
