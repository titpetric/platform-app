-- Blog article table
CREATE TABLE IF NOT EXISTS article (
    `id` TEXT PRIMARY KEY,
    `slug` TEXT NOT NULL UNIQUE,
    `title` TEXT NOT NULL,
    `filename` TEXT NOT NULL,
    `description` TEXT,
    `date` DATETIME NOT NULL,
    `og_image` TEXT,
    `layout` TEXT DEFAULT 'post',
    `source` TEXT,
    `url` TEXT NOT NULL,
    `draft` INTEGER NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for date-based queries (list, newest first)
CREATE INDEX IF NOT EXISTS idx_article_date ON article(date DESC);

-- Index for slug lookups (detail page)
CREATE INDEX IF NOT EXISTS idx_article_slug ON article(slug);

-- Index for filtering by layout type
CREATE INDEX IF NOT EXISTS idx_article_layout ON article(layout);

-- Index for recent articles
CREATE INDEX IF NOT EXISTS idx_article_created_at ON article(created_at DESC);

-- Index for draft filtering
CREATE INDEX IF NOT EXISTS idx_article_draft ON article(draft);

-- Index for published articles (not draft, date <= now)
CREATE INDEX IF NOT EXISTS idx_article_published ON article(draft, date DESC);
