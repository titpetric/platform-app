--
-- This is an append only file. Added statements will run during migration.
--

CREATE TABLE IF NOT EXISTS email (
    id TEXT PRIMARY KEY,
    recipient TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at DATETIME NOT NULL,
    sent_at DATETIME,
    error TEXT,
    retry_count INTEGER DEFAULT 0,
    last_error TEXT,
    last_retry DATETIME
);

CREATE TABLE IF NOT EXISTS email_sent (
    id TEXT PRIMARY KEY,
    recipient TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    status TEXT DEFAULT 'sent',
    created_at DATETIME NOT NULL,
    sent_at DATETIME NOT NULL,
    error TEXT,
    retry_count INTEGER DEFAULT 0,
    last_error TEXT,
    last_retry DATETIME
);

CREATE TABLE IF NOT EXISTS email_failed (
    id TEXT PRIMARY KEY,
    recipient TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    status TEXT DEFAULT 'failed',
    created_at DATETIME NOT NULL,
    sent_at DATETIME,
    error TEXT,
    retry_count INTEGER DEFAULT 0,
    last_error TEXT,
    last_retry DATETIME
);
