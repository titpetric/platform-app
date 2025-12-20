# Blog Module Architecture

Complete overview of the blog module implementation for the platform.

## Directory Structure

```
blog/
├── README.md                # Project overview
├── docs/
│   ├── ARCHITECTURE.md      # This file
│   ├── IMPLEMENTATION.md    # Implementation details
│   ├── SETUP.md             # Setup and configuration
│   ├── PORTING.md           # WebC to VueGo migration guide
│   └── MARKDOWN.md          # Markdown rendering
│
├── blog.go                  # Module implementation
├── handlers.go              # HTTP request handlers
│
├── model/
│   └── article.go          # Article data types
│
├── storage/
│   ├── db.go               # Database connection helper
│   ├── storage.go          # Storage interface
│   └── articles.go         # SQL operations
│
├── template/
│   ├── base.go             # Base layout rendering
│   ├── post.go             # Post layout rendering
│   ├── component.go        # Component interface
│   └── helpers.go          # Template helper functions
│
├── schema/
│   └── 001_articles.sql    # Database schema definition
│
├── theme/
│   ├── layouts/
│   │   ├── base.vuego      # HTML root layout
│   │   └── post.vuego      # Article detail layout
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
│   └── assets/             # CSS, JS, images
│
├── data/                   # Markdown article files (38 articles)
└── config/                 # Configuration JSON files
```

## Core Components

### 1. Module (blog.go)

The blog module implements the `platform.Module` interface:

```go
type Module struct {
	platform.UnimplementedModule
	dataDir    string                    // Path to markdown files
	repository *storage.Storage          // Database operations
	articles   map[string]*model.Article // In-memory index
}

// Implements:
// - Name() string
// - Start(ctx, platform) error
// - Mount(ctx, router) error
// - Stop(ctx) error
```

**Responsibilities:**
- Lifecycle management (Start/Stop)
- Route registration (Mount)
- Markdown file scanning and parsing
- Article indexing

### 2. Handlers (handlers.go)

HTTP request handlers:

```go
type Handlers struct {
	repository *storage.Storage
	helpers    *template.TemplateHelpers
}

// Methods:
// - ListArticlesJSON(w, r)      GET /api/blog/articles
// - GetArticleJSON(w, r)        GET /api/blog/articles/{slug}
// - SearchArticlesJSON(w, r)    GET /api/blog/search
// - ListArticlesHTML(w, r)      GET /blog/
// - GetArticleHTML(w, r)        GET /blog/{slug}
```

**Responsibilities:**
- Request parsing and validation
- Database queries via storage
- Template rendering
- Response serialization (JSON/HTML)
- Cache control headers

### 3. Storage Package

Three-layer storage abstraction:

#### db.go

```go
func DB(ctx context.Context) (*sqlx.DB, error)

// Gets database connection from platform.Database.Connect()
// Named connection: "blog"
// Reuses connections, no explicit close needed
```

#### storage.go

```go
type Storage struct {
	db *sqlx.DB
}

// Methods provide a clean interface:
// - GetArticleBySlug(ctx, slug) (*Article, error)
// - GetAllArticles(ctx) ([]Article, error)
// - SearchArticles(ctx, query) ([]Article, error)
// - InsertArticle(ctx, article) error
// - CountArticles(ctx) (int, error)
// - InitSchema(ctx) error
```

#### articles.go

```go
// Package-level functions for SQL operations
// Called by Storage methods
// Direct sqlx.DB access
```

**Benefits:**
- Clean API for handlers
- Dependency injection friendly
- Testable (easy to mock)
- Reusable (can be used from other modules)

### 4. Models (model/)

Data types:

```go
// Article represents a blog post
type Article struct {
	ID          string // Unique ID
	Slug        string // URL-friendly identifier
	Title       string
	Description string
	Content     string // Raw markdown
	Date        time.Time
	OGImage     string
	Layout      string // Template name
	Source      string // Optional external source
	URL         string // Generated path
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Metadata from YAML front matter
type Metadata struct {
	Title       string
	Description string
	OGImage     string
	Date        string
	Layout      string
	Source      string
}
```

### 5. Template Layer (template/)

Four files handle template rendering:

#### base.go

```go
type BaseData struct {
	Title, Description, OGImage string
	Content                     string
	Meta, Page                  interface{}
	Classnames                  string
	// Helper functions...
}

func Base(ctx Context, data *BaseData) (string, error)

// Renders base.vuego layout
```

#### post.go

```go
type PostData struct {
	Title, Description string
	Content            string
	Date               time.Time
	// Helper functions for formatting...
}

func Post(ctx Context, data *PostData) (string, error)

// Renders post.vuego layout
```

#### component.go

```go
type Component interface {
	Render(ctx Context, w io.Writer) error
}

// Implementations:
// - ArticleListComponent
// - Individual component types
```

#### helpers.go

```go
type TemplateHelpers struct {
	Meta, Page  interface{}
	MetaTitle   func(string) string
	PostDate    func(time.Time) string
	ReadingTime func(string) string
	// More helpers...
}

func DefaultHelpers(meta, page interface{}) *TemplateHelpers

// Provides sensible defaults
```

## Data Flow

### Startup Flow

```
Platform.Start()
    ↓
Module.Start(ctx, platform)
    ↓
storage.DB(ctx) → platform.Database.Connect(ctx, "blog")
    ↓
Storage.InitSchema(ctx) → Create SQLite tables
    ↓
Module.scanMarkdownFiles(ctx)
    ├→ filepath.WalkDir("./data")
    ├→ parseMarkdownFile() for each .md
    └→ Storage.InsertArticle() for each
    ↓
Platform.Mount()
    ↓
Module.Mount(ctx, router)
    ↓
Register handlers with router
```

### Request Flow

```
HTTP Request
    ↓
Router dispatches to Handler
    ↓
Handler.GetArticleHTML(w, r)
    ├→ slug := chi.URLParam(r, "slug")
    ├→ article := Storage.GetArticleBySlug(ctx, slug)
    ├→ htmlContent := blackfriday.Run(article.Content)
    ├→ postData := PostFromArticle(article, htmlContent, helpers)
    ├→ html := template.Post(ctx, postData)
    └→ w.Write(html)
    ↓
Response sent to client
```

### Database Schema

```
SQLite (in-memory or file-based)

articles
├── id (PK): text
├── slug (UNIQUE): text
├── title: text
├── description: text
├── content: text (markdown)
├── date: datetime
├── og_image: text
├── layout: text
├── source: text
├── url: text
├── created_at: datetime
└── updated_at: datetime

Indexes:
├── idx_articles_date (DESC)
├── idx_articles_slug
├── idx_articles_layout
└── idx_articles_created_at (DESC)

articles_fts (Virtual FTS5 table)
├── title
├── description
├── content
└── Triggers sync with articles table
```

## Key Design Patterns

### 1. Platform Integration

Uses platform's built-in features:
- **Module interface** - Lifecycle management
- **Router injection** - Route registration
- **Database connection pool** - Named connections via environment

```go
// Platform automatically:
// - Creates module from registry
// - Calls Start() for initialization
// - Calls Mount() for routing
// - Calls Stop() for cleanup
// - Manages database lifecycle
```

### 2. Dependency Injection

Dependencies passed explicitly:

```go
// Module doesn't create its own handlers
h := NewHandlers(m.repository, nil)

// Handlers don't create their own storage
func NewHandlers(repo *storage.Storage, helpers *TemplateHelpers)

// Storage doesn't manage its own DB
func NewStorage(db *sqlx.DB)
```

### 3. Separation of Concerns

Clear boundaries:

| Layer     | Responsibility | Can Access                |
|-----------|----------------|---------------------------|
| Module    | Lifecycle      | Router, Platform          |
| Handlers  | HTTP           | Storage, Templates        |
| Storage   | SQL            | Database                  |
| Templates | Rendering      | Helper functions          |
| Models    | Data types     | (Nothing - plain structs) |

### 4. Context Flow

All async operations accept `context.Context`:

```go
// Enables:
// - Request cancellation
// - Distributed tracing
// - Deadline enforcement
// - Value passing (platform.FromContext)
```

### 5. Error Handling

Errors propagate up:

```go
// Start() errors prevent module initialization
// Handler errors return HTTP errors
// Storage errors bubble up to handlers
// No silent failures
```

## API Contract

### Routes

| Method | Path                      | Response | Cache |
|--------|---------------------------|----------|-------|
| GET    | /api/blog/articles        | JSON     | 5min  |
| GET    | /api/blog/articles/{slug} | JSON     | 1hr   |
| GET    | /api/blog/search?q=X      | JSON     | 5min  |
| GET    | /blog/                    | HTML     | 5min  |
| GET    | /blog/{slug}              | HTML     | 1hr   |

### Content Negotiation

Handlers automatically select format based on Accept header (basic):
- Default: HTML
- Explicit: JSON via Content-Type

### Data Models

All models use:
- JSON tags for API serialization
- `time.Time` for dates
- String IDs for flexibility
- Optional fields are pointers or zero-values

## Performance Characteristics

### Memory
- Articles kept in memory: O(n) where n = article count
- Typical: ~50KB per article in memory
- For 1000 articles: ~50MB

### Database Queries
- List articles: O(n) but indexed by date
- Get article: O(1) with slug index
- Search: O(n log n) with FTS index
- Insert: O(log n) with indexes

### Caching
- HTML responses: 5min-1hr max-age
- JSON responses: 5min max-age
- Browser/CDN will cache appropriately

## Security Considerations

### Input Validation
- Slugs validated as URL parameters
- Search queries sanitized (no SQL injection)
- Markdown content rendered safely (blackfriday escapes HTML)

### Output Encoding
- HTML properly escaped by template engine
- JSON properly encoded by json.Encoder
- SQL parameters bound (sqlx prevents injection)

### Access Control
- No authentication required (public blog)
- Can be extended with middleware via platform

## Testing Strategy

### Unit Tests
- Storage operations with mock database
- Template rendering with fixed data
- Handler logic with mock storage

### Integration Tests
- Full module lifecycle with test platform
- Database schema creation and indexing
- Markdown parsing and front matter extraction

### Examples

See docs/SETUP.md for testing code samples.

## Future Extensions

### Possible Enhancements

1. **Comments system**
   - New comments table
   - Comment handlers
   - Comment templates

2. **Categories/Tags**
   - New tags table
   - Article-tag relationships
   - Tag filtering in queries

3. **Draft mode**
   - Published boolean field
   - Filter unpublished in queries
   - Authentication for draft preview

4. **Search improvements**
   - FTS5 ranking
   - Snippet generation
   - Search analytics

5. **Performance**
   - Cache invalidation on updates
   - CDN integration
   - Database connection pooling tuning

6. **Analytics**
   - View counting
   - Popular articles
   - Search analytics

7. **Admin interface**
   - Create/edit articles via UI
   - Publish scheduling
   - Statistics dashboard

All can be implemented by extending the module while maintaining the current architecture.
