# Package blog

```go
import (
	"github.com/titpetric/platform-example/blog"
}
```

## Types

```go
// Generator generates static HTML files from blog content
type Generator struct {
	module    *Module
	outputDir string
}
```

```go
// Handlers handles HTTP requests for the blog module
type Handlers struct {
	repository *storage.Storage
	views      *view.Views
}
```

```go
// Module implements the blog module for the platform
type Module struct {
	platform.UnimplementedModule

	// Data directory for markdown files
	dataDir string

	// Storage for database operations
	repository *storage.Storage

	// Articles index for in-memory access
	articles map[string]*model.Article

	// Theme fs that combines embedded theme and live theme/ folder.
	themeFS fs.FS
}
```

```go
type OverlayFS struct {
	Upper fs.FS
	Lower fs.FS
}
```

## Function symbols

- `func NewGenerator (m *Module, outputDir string) *Generator`
- `func NewHandlers (repo *storage.Storage, themeFS fs.FS) (*Handlers, error)`
- `func NewModule (dataDir string) *Module`
- `func NewOverlayFS (upper,lower fs.FS) *OverlayFS`
- `func (*Generator) Generate (ctx context.Context) error`
- `func (*Handlers) GetArticleHTML (w http.ResponseWriter, r *http.Request)`
- `func (*Handlers) GetArticleJSON (w http.ResponseWriter, r *http.Request)`
- `func (*Handlers) GetAtomFeed (w http.ResponseWriter, r *http.Request)`
- `func (*Handlers) IndexHTML (w http.ResponseWriter, r *http.Request)`
- `func (*Handlers) ListArticlesHTML (w http.ResponseWriter, r *http.Request)`
- `func (*Handlers) ListArticlesJSON (w http.ResponseWriter, r *http.Request)`
- `func (*Handlers) SearchArticlesJSON (w http.ResponseWriter, r *http.Request)`
- `func (*Module) Mount (_ context.Context, r platform.Router) error`
- `func (*Module) Name () string`
- `func (*Module) ScanMarkdownFiles (ctx context.Context) (int, error)`
- `func (*Module) SetRepository (repo *storage.Storage)`
- `func (*Module) Start (ctx context.Context) error`
- `func (*Module) Stop (context.Context) error`
- `func (*OverlayFS) Glob (pattern string) ([]string, error)`
- `func (*OverlayFS) Open (name string) (fs.File, error)`
- `func (*OverlayFS) ReadDir (name string) ([]fs.DirEntry, error)`

### NewGenerator

NewGenerator creates a new Generator instance

```go
func NewGenerator(m *Module, outputDir string) *Generator
```

### NewHandlers

NewHandlers creates a new Handlers instance with the given storage

```go
func NewHandlers(repo *storage.Storage, themeFS fs.FS) (*Handlers, error)
```

### NewModule

NewModule creates a new blog module instance

```go
func NewModule(dataDir string) *Module
```

### Generate

Generate generates all static HTML files

```go
func (*Generator) Generate(ctx context.Context) error
```

### GetArticleHTML

GetArticleHTML returns a single article as HTML

```go
func (*Handlers) GetArticleHTML(w http.ResponseWriter, r *http.Request)
```

### GetArticleJSON

GetArticleJSON returns a single article as JSON

```go
func (*Handlers) GetArticleJSON(w http.ResponseWriter, r *http.Request)
```

### GetAtomFeed

GetAtomFeed returns an Atom XML feed of all articles

```go
func (*Handlers) GetAtomFeed(w http.ResponseWriter, r *http.Request)
```

### IndexHTML

IndexHTML returns an HTML index page listing blogs

```go
func (*Handlers) IndexHTML(w http.ResponseWriter, r *http.Request)
```

### ListArticlesHTML

ListArticlesHTML returns an HTML list of articles

```go
func (*Handlers) ListArticlesHTML(w http.ResponseWriter, r *http.Request)
```

### ListArticlesJSON

ListArticlesJSON returns a JSON list of all articles

```go
func (*Handlers) ListArticlesJSON(w http.ResponseWriter, r *http.Request)
```

### SearchArticlesJSON

SearchArticlesJSON performs full-text search on articles

```go
func (*Handlers) SearchArticlesJSON(w http.ResponseWriter, r *http.Request)
```

### Mount

Mount registers the blog routes with the router

```go
func (*Module) Mount(_ context.Context, r platform.Router) error
```

### Name

Name returns the module name

```go
func (*Module) Name() string
```

### ScanMarkdownFiles

ScanMarkdownFiles scans the data directory for markdown files and indexes them Returns the count of scanned files

```go
func (*Module) ScanMarkdownFiles(ctx context.Context) (int, error)
```

### SetRepository

SetRepository sets the repository on the module

```go
func (*Module) SetRepository(repo *storage.Storage)
```

### Start

Start initializes the blog module by scanning markdown files and building the index

```go
func (*Module) Start(ctx context.Context) error
```

### Stop

Stop is called when the module is shutting down

```go
func (*Module) Stop(context.Context) error
```

### NewOverlayFS

```go
func NewOverlayFS(upper, lower fs.FS) *OverlayFS
```

### Glob

```go
func (*OverlayFS) Glob(pattern string) ([]string, error)
```

### Open

```go
func (*OverlayFS) Open(name string) (fs.File, error)
```

### ReadDir

```go
func (*OverlayFS) ReadDir(name string) ([]fs.DirEntry, error)
```
