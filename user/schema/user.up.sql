-- user: Stores user profile information using ULID as primary key
CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY NOT NULL,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    deleted_at DATETIME,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_user_deleted_at ON user(deleted_at);

-- user_auth: Stores user authentication credentials
CREATE TABLE IF NOT EXISTS user_auth (
    user_id TEXT PRIMARY KEY NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auth_email ON user_auth(email);

-- user_group: Stores user group information using ULID as primary key
CREATE TABLE IF NOT EXISTS user_group (
    id TEXT PRIMARY KEY NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    created_at DATETIME,
    updated_at DATETIME
);

-- user_group_member: Stores user memberships in groups
CREATE TABLE IF NOT EXISTS user_group_member (
    user_group_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    joined_at DATETIME,
    PRIMARY KEY (user_group_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_group_member_user_id ON user_group_member(user_id);

-- user_session: Stores immutable user sessions with ULID IDs and expiration
CREATE TABLE IF NOT EXISTS user_session (
    id TEXT PRIMARY KEY NOT NULL,
    user_id TEXT NOT NULL,
    expires_at DATETIME,
    created_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_user_session_user_id ON user_session(user_id);
CREATE INDEX IF NOT EXISTS idx_user_session_expires_at ON user_session(expires_at);
