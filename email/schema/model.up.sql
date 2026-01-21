--
-- This is an append only file. Added statements will run during migration.
--

CREATE TABLE IF NOT EXISTS email (
    id TEXT PRIMARY KEY,
    recipient TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    error TEXT,
    created_at DATETIME NOT NULL,
    sent_at DATETIME,
    retry_count INTEGER DEFAULT 0,
    retry_error TEXT,
    retry_at DATETIME
);

CREATE TABLE IF NOT EXISTS email_sent (
    id TEXT PRIMARY KEY,
    recipient TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    status TEXT DEFAULT 'sent',
    error TEXT,
    created_at DATETIME NOT NULL,
    sent_at DATETIME NOT NULL,
    retry_count INTEGER DEFAULT 0,
    retry_error TEXT,
    retry_at DATETIME
);

CREATE TABLE IF NOT EXISTS email_failed (
    id TEXT PRIMARY KEY,
    recipient TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    status TEXT DEFAULT 'failed',
    error TEXT,
    created_at DATETIME NOT NULL,
    sent_at DATETIME,
    retry_count INTEGER DEFAULT 0,
    retry_error TEXT,
    retry_at DATETIME
);
