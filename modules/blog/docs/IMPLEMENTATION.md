# Blog Module Implementation

This document describes the implementation of the blog module for the platform.

## Architecture

The blog module follows the platform architecture with clear separation of concerns:

```
blog/
├── README.md            # Project overview
├── Taskfile.yaml        # Build automation
├── blog.go              # Module implementation (lifecycle)
├── handlers.go          # HTTP handlers
├── docs/                # Documentation
│   ├── ARCHITECTURE.md
│   ├── IMPLEMENTATION.md
│   ├── SETUP.md
│   ├── PORTING.md
│   ├── MARKDOWN.md
│   ├── api.md           # Generated API docs
│   └── testing-coverage.md  # Generated coverage report
├── template/
│   ├── base.go         # Base layout template
│   ├── post.go         # Post layout template
│   ├── component.go    # Component interface
│   └── helpers.go      # Template helper functions
├── storage/
│   ├── db.go           # Database connection via platform.Database
│   ├── storage.go      # Storage interface
│   └── articles.go     # Article operations
├── markdown/
│   └── markdown.go     # Markdown rendering with syntax highlighting
├── model/
│   └── article.go      # Article data models
├── schema/
│   └── 001_articles.sql # SQLite schema
└── theme/              # Vue templates
```

## Database Access Pattern

The blog module uses `platform.Database.Connect()` to get a database connection:

```go
// In storage/db.go
func DB(ctx context.Context) (*sqlx.DB, error) {
	return platform.Database.Connect(ctx, "blog")
}
```

This pattern:
- Uses the "blog" named connection from environment: `PLATFORM_DB_BLOG`
- Reuses connections between module calls
- Is safe for concurrent use
- Eliminates explicit connection management

## Module Lifecycle

### Start()
1. Gets database connection via `storage.DB(ctx)`
2. Creates `Storage` instance
3. Initializes database schema
4. Scans `./data` directory for `.md` files
5. Parses YAML front matter from each file
6. Indexes articles in both memory map and database

### Mount()

Registers HTTP routes after Start() completes:

**JSON API Routes:**
- `GET /api/blog/articles` - List all articles
- `GET /api/blog/articles/{slug}` - Get single article
- `GET /api/blog/search?q=query` - Search articles

**HTML Routes:**
- `GET /blog/` - Article list (HTML)
- `GET /blog/{slug}` - Article detail (HTML)

### Stop()

No explicit cleanup needed - platform manages database lifecycle.

## Handlers

The `Handlers` struct manages HTTP request handling:

```go
type Handlers struct {
	repository *storage.Storage
	helpers    *template.TemplateHelpers
}
```

Handlers are created in `Mount()` with the module's storage instance:

```go
h := NewHandlers(m.repository, nil)
```

Each handler:
- Uses `r.Context()` for context passing
- Calls storage methods for data
- Sets appropriate headers (Content-Type, Cache-Control)
- Renders templates for HTML responses

## Storage Pattern

The `Storage` struct provides an interface to database operations:

```go
type Storage struct {
	db *sqlx.DB
}

// Methods:
// - GetArticleBySlug(ctx, slug) (*Article, error)
// - GetAllArticles(ctx) ([]Article, error)
// - SearchArticles(ctx, query) ([]Article, error)
// - InsertArticle(ctx, article) error
// - CountArticles(ctx) (int, error)
// - InitSchema(ctx) error
```

All database functions accept `context.Context` for cancellation and observability.

## Templates

### Layout Templates

**base.vuego** - Root HTML layout with:
- Meta tags (title, description, OG image)
- CSS and JS bundles
- Theme component (site-header, site-footer)
- Content slot

**post.vuego** - Article detail layout with:
- Article title and metadata
- Published date and reading time
- Article content
- Back-to-articles link

### Template Rendering

Templates are rendered in handlers:

```go
// HTML rendering
postData := template.PostFromArticle(article, htmlContent, h.helpers)
html, err := template.Post(r.Context(), postData)
w.Write([]byte(html))
```

Template data includes:
- Article content (rendered markdown as HTML)
- Helper functions for formatting
- Meta information and navigation

## Markdown File Processing

Articles are sourced from markdown files in the `data/` directory.

### Front Matter Format

```yaml
---
title: Article Title
description: Brief description
date: 2024-01-01
layout: post
ogImage: /path/to/image.png
source: optional-external-source
---

# Article content starts here
```

### Processing Steps

1. **Parse YAML** - Extract front matter metadata
2. **Generate ID** - Create unique ID from slug + timestamp
3. **Parse Date** - Convert YAML date string to time.Time
4. **Set Layout** - Use layout from metadata or default to "post"
5. **Generate URL** - Create URL path `/blog/{slug}/`
6. **Store** - Save to both memory map and database

## Database Schema

SQLite in-memory database with:

- **articles table** - Full article data with indexes
- **articles_fts** - Virtual FTS5 table for full-text search
- **Triggers** - Keep FTS index in sync

Indexes on:
- `date DESC` - For sorting by publication date
- `slug` - For fast lookup by slug
- `layout` - For filtering by layout type
- `created_at DESC` - For recent articles

## Configuration

Environment variables:

```bash
# Database connection (required)
PLATFORM_DB_BLOG="sqlite:///tmp/blog.db"

# Or for in-memory (development)
PLATFORM_DB_BLOG="sqlite://:memory:"
```

## Usage Example

```go
// In main application setup
func init() {
	// Module auto-registers via platform.Register()
}

func main() {
	ctx := context.Background()

	// Platform automatically calls:
	// 1. m.Start(ctx) - Initialize module
	// 2. m.Mount(router) - Register routes

	platform.Start(ctx)
}
```

## API Endpoints

### List Articles

```
GET /api/blog/articles
Content-Type: application/json

Response:
{
  "articles": [
    {
      "id": "my-article-20240101000000",
      "slug": "my-article",
      "title": "My Article",
      "description": "Article description",
      "date": "2024-01-01T00:00:00Z",
      "url": "/blog/my-article/"
    }
  ],
  "total": 1,
  "page": 1,
  "pageSize": 1
}
```

### Get Article

```
GET /api/blog/articles/{slug}
Content-Type: application/json

Response:
{
  "id": "my-article-20240101000000",
  "slug": "my-article",
  "title": "My Article",
  "description": "Article description",
  "content": "Raw markdown content",
  "date": "2024-01-01T00:00:00Z",
  "url": "/blog/my-article/",
  "layout": "post"
}
```

### Search Articles

```
GET /api/blog/search?q=keyword
Content-Type: application/json

Response:
{
  "articles": [...],
  "total": 1,
  "query": "keyword"
}
```

### HTML Routes

```
GET /blog/                  # Article list page
GET /blog/my-article/       # Article detail page
```

## Key Design Decisions

1. **Separation of Concerns**
   - Module manages lifecycle (Start/Stop/Mount)
   - Handlers manage HTTP/template rendering
   - Storage manages database access
   - Templates manage HTML generation

2. **Storage Abstraction**
   - All DB access goes through Storage struct
   - Eliminates tight coupling
   - Enables testing with mock storage

3. **Template Flexibility**
   - Templates don't know about HTTP layer
   - Can be reused in other contexts (CLI, gRPC, etc.)
   - Helper functions provide formatting

4. **In-Memory + Database**
   - Memory map for fast slug lookups
   - Database for complex queries and indexing
   - Both synced during Start()
   - Source of truth: markdown files on disk

5. **Named Database Connection**
   - Enables multiple modules to use different databases
   - Least privilege principle
   - Centralized configuration via environment
