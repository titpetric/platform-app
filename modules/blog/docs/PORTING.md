# Blog Module Porting Guide

This document describes the complete porting of the blog system from Eleventy/WebC to a Go platform module using vuego templates.

## Overview

The blog has been completely refactored to use Go with the `github.com/titpetric/platform` module system. Key components include:

1. **Templates**: All `.webc` files ported to `.vuego` format
2. **Models**: Go types for article metadata and data structures
3. **Database**: SQLite schema for in-memory indexing
4. **Module**: Platform module implementing the Module interface

## Directory Structure

```
blog/
├── blog.go                  # Main module implementation
├── model/
│   └── article.go          # Go types for articles and metadata
├── schema/
│   └── 001_articles.sql    # SQLite schema definition
├── theme/
│   ├── layouts/
│   │   ├── base.vuego      # Base HTML layout
│   │   └── post.vuego      # Post layout
│   ├── components/
│   │   ├── article-list.vuego
│   │   ├── inline-svg.vuego
│   │   ├── lite-youtube.vuego
│   │   ├── page-timer.vuego
│   │   ├── site-footer.vuego
│   │   ├── site-header.vuego
│   │   ├── target-toggler.vuego
│   │   ├── theme-machine.vuego
│   │   └── info-cta.vuego
│   ├── pages/
│   │   ├── blog.vuego
│   │   ├── index.vuego
│   │   └── 404.vuego
│   └── assets/            # CSS, JS, images
├── data/                  # Markdown source files
└── config/                # Configuration JSON files
```

## Template Migration (webc → vuego)

### Syntax Changes

| webc                    | vuego                   |
|-------------------------|-------------------------|
| `{{ value }}`           | `{{ value }}`           |
| `:attr="value"`         | `:attr="value"`         |
| `@text="value"`         | `{{ value }}`           |
| `@raw="value"`          | `v-html="value"`        |
| `v-if`                  | `v-if`                  |
| `v-for="item of items"` | `v-for="item of items"` |
| `webc:if`               | `v-if`                  |
| `webc:for`              | `v-for`                 |
| `webc:nokeep`           | (removed)               |
| `webc:root="override"`  | (removed)               |
| `webc:setup`            | (removed - use props)   |

### Example Conversions

**base.webc → base.vuego:**

```webc
<title @text="metaTitle(title)"></title>
```

becomes:

```vuego
<title>{{ metaTitle(title) }}</title>
```

**site-header.webc → site-header.vuego:**

```webc
<li webc:if="index < 3" webc:for="(item, index) of menu">
  <a :href="item.url" @text="item.label"></a>
</li>
```

becomes:

```vuego
<li v-for="(item, index) of navigation.menu" v-if="index < 3">
  <a :href="item.url">{{ item.label }}</a>
</li>
```

## Go Module Implementation

### Module Structure

The blog module implements the `platform.Module` interface:

```go
type Module struct {
	db       *sql.DB
	dataDir  string
	articles map[string]*model.Article
}
```

### Key Methods

1. **Name()** - Returns "blog"
2. **Mount(ctx, router)** - Registers HTTP routes for:
   - GET `/api/blog/articles` - List all articles (JSON)
   - GET `/api/blog/articles/{slug}` - Get single article (JSON)
   - GET `/api/blog/search` - Full-text search (JSON)
   - GET `/blog/` - Blog list page (HTML)
   - GET `/blog/{slug}` - Article page (HTML)

3. **Start(ctx, platform)** - Lifecycle method that:
   - Initializes in-memory SQLite database
   - Creates database schema
   - Scans `data/` directory for `.md` files
   - Parses YAML front matter from each file
   - Indexes articles in database and memory map

4. **Stop(ctx)** - Gracefully closes database connection

### Markdown File Processing

Each markdown file in the `data/` directory is processed as follows:

1. **YAML Front Matter Extraction**

   ```yaml
   ---
   title: Article Title
   description: Brief description
   date: 2024-01-01
   layout: post
   ogImage: /path/to/image.png
   ---
   ```

2. **Metadata Parsing** to `model.Metadata` struct
3. **Article Creation** with computed fields:
   - `ID`: Generated from slug and timestamp
   - `URL`: `/blog/{slug}/`
   - `Layout`: From YAML or defaults to "post"

4. **Database Insertion** using prepared statements

### Database Schema

SQLite `:memory:` database with:
- **articles table**: Full article data with indexes on date, slug, layout
- **articles_fts**: Virtual FTS5 table for full-text search
- **Triggers**: Keep FTS index in sync with articles table

## Model Types

### Article

```go
type Article struct {
	ID          string // Unique identifier
	Slug        string // URL-friendly identifier
	Title       string
	Description string
	Content     string // Raw markdown content
	Date        time.Time
	OGImage     string
	Layout      string // Template layout name
	Source      string // External source if applicable
	ReadingTime string // Calculated reading time
	URL         string // Full article URL
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
```

### Metadata

```go
type Metadata struct {
	Title       string
	Description string
	OGImage     string
	Date        string
	Layout      string
	Source      string
}
```

## Registration

The module is registered in the `init()` function:

```go
func init() {
	platform.Register(NewModule("./data"))
}
```

This allows automatic discovery and registration with the platform.

## API Endpoints

### List Articles

```
GET /api/blog/articles
```

Returns JSON array of articles sorted by date (newest first).

### Get Article

```
GET /api/blog/articles/{slug}
```

Returns full article JSON including content.

### Search Articles

```
GET /api/blog/search?q=query
```

Returns matching articles based on title, description, or slug.

## Template Variables

When rendering templates, the following variables are available:

- `articles` - Array of Article objects
- `article` - Current Article object (detail pages)
- `navigation` - Navigation configuration
- `meta` - Site metadata
- `page` - Current page information
- `themes` - Available themes

## Migration Checklist

- [X] Port all .webc layout files to .vuego
- [X] Port all .webc component files to .vuego
- [X] Port all .webc page files to .vuego
- [X] Create Go model package with Article types
- [X] Create SQLite schema
- [X] Implement blog.go module
- [X] Implement Mount() for route registration
- [X] Implement Start() for markdown scanning and indexing
- [X] Implement Stop() for graceful shutdown
- [ ] Implement HTML template rendering with vuego
- [ ] Add content markdown rendering (HTML conversion)
- [ ] Add reading time calculation
- [ ] Add search functionality with FTS
- [ ] Add pagination support
- [ ] Add cache headers and optimization

## Notes

- The module maintains articles in both an in-memory map (for fast lookups) and SQLite database (for complex queries)
- The database is `:memory:` - articles are re-indexed on every application start
- Markdown content is stored raw and would need a markdown renderer (e.g., `github.com/russross/blackfriday`) for HTML output
- Template rendering with vuego is abstracted away from the HTTP handlers to keep the module clean
- All HTTP handlers are independent of template implementation for maximum flexibility
