-- Enable foreign key support
PRAGMA foreign_keys = ON;

-- Snippet Table
CREATE TABLE snippets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- Unique row ID
    uuid TEXT NOT NULL,                    -- Time-based UUID (same for different versions)
    name TEXT,                              -- Friendly name (optional)
    code TEXT NOT NULL,                     -- Code snippet
    language TEXT NOT NULL,                 -- Programming language
    tags TEXT,                              -- Comma-separated tags for searching/filtering (optional)
    description TEXT,                        -- Description (optional)
    source TEXT,                             -- Source (site, project, etc.) (optional)
    date_added DATETIME DEFAULT CURRENT_TIMESTAMP, -- Date added
    version INTEGER NOT NULL DEFAULT 1,      -- Versioning number
    superseded_by INTEGER,                  -- ID of the next version (optional)
    FOREIGN KEY (superseded_by) REFERENCES snippets(id) ON DELETE SET NULL
);

-- Index for faster search by UUID (since it's not unique)
CREATE INDEX idx_snippets_uuid ON snippets(uuid);

-- Trigger to delete all versions of snippets when one is deleted.
CREATE TRIGGER delete_snippet_history
BEFORE DELETE ON snippets
FOR EACH ROW
BEGIN
    DELETE FROM snippets WHERE uuid = OLD.uuid;
END;

-- Group Table
CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- Unique row ID
    group_name TEXT NOT NULL,              -- Friendly group name
    description TEXT,                       -- Group description (optional)
    uuid_list TEXT NOT NULL,                -- Comma-separated list of snippet UUIDs
    date_added DATETIME DEFAULT CURRENT_TIMESTAMP, -- Date the group was added
    date_updated DATETIME DEFAULT CURRENT_TIMESTAMP -- Date last updated
);

-- Trigger to update date_updated when a group is modified
CREATE TRIGGER update_group_timestamp
AFTER UPDATE ON groups
FOR EACH ROW
BEGIN
    UPDATE groups SET date_updated = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

