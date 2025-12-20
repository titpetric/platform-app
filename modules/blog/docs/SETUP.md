# Blog Module Setup Guide

## Prerequisites

- Go 1.20+
- SQLite 3.x
- Platform package: `github.com/titpetric/platform`

## Installation

### 1. Dependencies

```bash
go get github.com/jmoiron/sqlx
go get github.com/mattn/go-sqlite3
go get github.com/russross/blackfriday/v2
go get github.com/titpetric/vuego
go get gopkg.in/yaml.v3
```

### 2. Module Registration

The blog module auto-registers via `init()` function in blog.go:

```go
func init() {
	platform.Register(NewModule("./data"))
}
```

No additional registration code needed in your main application.

## Configuration

### Environment Variables

Set the database connection for the blog module:

```bash
# SQLite file-based database
export PLATFORM_DB_BLOG="sqlite:///var/lib/blog.db"

# SQLite in-memory (development only)
export PLATFORM_DB_BLOG="sqlite://:memory:"

# PostgreSQL
export PLATFORM_DB_BLOG="postgres://user:pass@localhost/blog_db"

# MySQL
export PLATFORM_DB_BLOG="mysql://user:pass@localhost/blog_db"
```

### Data Directory

By default, the module scans `./data` for markdown files:

```go
NewModule("./data") // Scans ./data directory
```

Create your data directory:

```bash
mkdir -p ./data
```

## Creating Articles

### Markdown File Format

Create files in the `data/` directory with `.md` extension:

**data/my-first-article.md:**

```markdown
---
title: My First Article
description: This is my first article
date: 2024-01-15
layout: post
ogImage: /social/my-first-article.png
---

# My First Article

This is the article content. It supports **markdown** formatting.

- Bullet points
- Code blocks
- And more!
```

### Front Matter Fields

| Field       | Type   | Required | Notes                             |
|-------------|--------|----------|-----------------------------------|
| title       | string | Yes      | Article title                     |
| description | string | No       | Meta description                  |
| date        | string | No       | Format: YYYY-MM-DD                |
| layout      | string | No       | Template layout (default: "post") |
| ogImage     | string | No       | Open Graph image URL              |
| source      | string | No       | External source/attribution       |

### File Naming

- Use lowercase with hyphens: `my-article.md`
- The filename becomes the slug: `my-article`
- URL generated: `/blog/my-article/`

## Module Lifecycle

### Startup

When the application starts:

1. Platform calls `m.Start(ctx)`
2. Module gets database connection from platform
3. Initializes SQLite schema
4. Scans `./data` for `.md` files
5. Parses YAML front matter
6. Indexes articles in database and memory

### Request Handling

When a request comes in:

1. Platform calls `m.Mount(router)` (once)
2. Handlers registered for blog routes
3. Requests routed to handler methods
4. Handlers query storage
5. Templates rendered
6. HTML/JSON response sent

### Shutdown

When application stops:

1. Platform calls `m.Stop(ctx)`
2. Module cleanup (database managed by platform)
3. All handlers stop receiving requests

## Database Initialization

The schema is automatically created during `Start()`:

```go
// In storage/storage.go InitSchema()
// Creates:
// - articles table
// - articles_fts virtual table
// - indexes on date, slug, layout, created_at
```

No manual migration needed for the blog module.

## API Routes

### JSON Endpoints

```bash
# List all articles
curl http://localhost:8080/api/blog/articles

# Get single article
curl http://localhost:8080/api/blog/articles/my-article

# Search articles
curl http://localhost:8080/api/blog/search?q=keyword
```

### HTML Endpoints

```bash
# Article list page
curl http://localhost:8080/blog/

# Article detail page
curl http://localhost:8080/blog/my-article/
```

## Testing

### Unit Tests

Test storage operations:

```go
func TestGetArticleBySlug(t *testing.T) {
	ctx := context.Background()
	db, _ := sqlx.Open("sqlite3", ":memory:")
	defer db.Close()

	storage := storage.NewStorage(db)
	storage.InitSchema(ctx)

	article := &model.Article{
		ID:    "test-1",
		Slug:  "test",
		Title: "Test Article",
	}
	storage.InsertArticle(ctx, article)

	retrieved, _ := storage.GetArticleBySlug(ctx, "test")
	assert.Equal(t, "Test Article", retrieved.Title)
}
```

### Integration Tests

Test with platform:

```go
func TestBlogModule(t *testing.T) {
	opts := platform.NewTestOptions()
	plat := platform.New(opts)

	m := NewModule("./testdata")
	plat.Register(m)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	plat.Start(ctx)
	defer plat.Stop(context.Background())

	// Test handlers here
}
```

## Troubleshooting

### Articles Not Indexed

Check:
1. Data directory exists: `ls -la ./data`
2. Markdown files have correct extension: `.md`
3. YAML front matter is valid YAML
4. Database connection works: Check `PLATFORM_DB_BLOG` env var

```bash
# Check data files
ls -la ./data/*.md

# Check database (if file-based)
sqlite3 /var/lib/blog.db "SELECT COUNT(*) FROM articles;"
```

### Template Rendering Fails

Check:
1. Template files exist in `theme/layouts/` and `theme/components/`
2. vuego is installed: `go list -m github.com/titpetric/vuego`
3. Template syntax is valid (vuego syntax, not webc)

### Database Connection Error

Check:
1. `PLATFORM_DB_BLOG` environment variable is set
2. Database file/server is accessible
3. Driver is imported (see Dependencies section)

```bash
# Test SQLite connection
sqlite3 ":memory:" "SELECT 1;"

# Test env var
echo $PLATFORM_DB_BLOG
```

## Performance Tuning

### Caching

HTML responses include cache headers:

```go
w.Header().Set("Cache-Control", "public, max-age=3600") // Articles: 1 hour
w.Header().Set("Cache-Control", "public, max-age=300")  // Search: 5 minutes
```

### Database Indexes

The schema creates indexes on:
- `date DESC` - For list operations
- `slug` - For detail page lookups
- `layout` - For filtering by template
- `created_at DESC` - For recent articles

### Memory Usage

Articles are kept in memory for fast O(1) slug lookups:

```go
m.articles map[string]*model.Article  // In-memory cache
```

For very large sites (10k+ articles), this can be optimized by removing the memory cache.

## Next Steps

1. Create sample markdown files in `./data/`
2. Start the application
3. Test endpoints: `curl http://localhost:8080/api/blog/articles`
4. Customize templates in `theme/`
5. Extend handlers with additional features
