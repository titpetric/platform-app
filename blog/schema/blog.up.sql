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

-- Blog settings table matching meta.yml structure
CREATE TABLE IF NOT EXISTS setting (
    `user_id` TEXT PRIMARY KEY,
    `meta_lang` TEXT DEFAULT 'en',
    `meta_url` TEXT DEFAULT '',
    `meta_author_name` TEXT DEFAULT '',
    `meta_subtitle` TEXT DEFAULT '',
    `meta_headshot` TEXT DEFAULT '/assets/images/headshot.jpg',
    `posts_per_page` INTEGER DEFAULT 10,
    `social_github` TEXT DEFAULT '',
    `social_twitter` TEXT DEFAULT '',
    `social_linkedin` TEXT DEFAULT '',
    `feature_webmention` INTEGER DEFAULT 0,
    `feature_pingback` INTEGER DEFAULT 0,
    `feature_comments` INTEGER DEFAULT 0,
    `feature_rss` INTEGER DEFAULT 1,
    `seo_title_suffix` TEXT DEFAULT '',
    `seo_default_image` TEXT DEFAULT '',
    `analytics_id` TEXT DEFAULT '',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
