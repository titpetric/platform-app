--
-- This is an append only file. Added statements will run during migration.
--

CREATE TABLE IF NOT EXISTS todo (
    id CHAR(26) PRIMARY KEY,      -- ULID stored as 26-character text
    title TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME
);
