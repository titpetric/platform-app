package blog

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/titpetric/platform"
	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/platform-app/internal/overlayfs"
	"github.com/titpetric/platform-app/blog/handlers"
	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/storage"
)

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

// NewModule creates a new blog module instance
func NewModule(dataDir string) *Module {
	// Sub into the theme directory since embed.FS includes the directory name
	themeSub, _ := fs.Sub(themeFS, "theme")

	// Check if local theme directory exists
	var overlay fs.FS = themeSub
	if _, err := os.Stat("theme"); err == nil {
		overlay = overlayfs.NewOverlayFS(os.DirFS("theme"), themeSub)
	}

	return &Module{
		dataDir:  dataDir,
		themeFS:  overlay,
		articles: make(map[string]*model.Article),
	}
}

// Name returns the module name
func (m *Module) Name() string {
	return "blog"
}

// Mount registers the blog routes with the router
func (m *Module) Mount(_ context.Context, r platform.Router) error {
	// Register all routes (public, admin, API, HTML, assets)
	return handlers.RegisterRoutes(r, m.repository, m.themeFS)
}

// Start initializes the blog module by scanning markdown files and building the index
func (m *Module) Start(ctx context.Context) error {
	// Create storage instance (includes migrations)
	var err error
	m.repository, err = storage.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	// Scan and index markdown files
	count, err := m.scanMarkdownFiles(ctx)
	if err != nil {
		fmt.Printf("warning: no markdown files scanned: %v", err)
	}

	fmt.Printf("[blog] scanned %d markdown files from %s\n", count, m.dataDir)

	// Verify articles were inserted
	total, err := m.repository.CountArticles(ctx)
	if err != nil {
		return fmt.Errorf("failed to count articles: %w", err)
	}
	fmt.Printf("[blog] verified %d articles in database\n", total)

	return nil
}

// Stop is called when the module is shutting down
func (m *Module) Stop(context.Context) error {
	// Nothing to clean up - database is managed by platform
	return nil
}

// SetRepository sets the repository on the module
func (m *Module) SetRepository(repo *storage.Storage) {
	m.repository = repo
}

// ScanMarkdownFiles scans the data directory for markdown files and indexes them
// Returns the count of scanned files
func (m *Module) ScanMarkdownFiles(ctx context.Context) (int, error) {
	return m.scanMarkdownFiles(ctx)
}

// scanMarkdownFiles scans the data directory for markdown files and indexes them
// Returns the count of scanned files
func (m *Module) scanMarkdownFiles(ctx context.Context) (int, error) {
	count := 0
	err := filepath.WalkDir(m.dataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only process markdown files
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		count++

		article, err := m.parseMarkdownFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Store in memory map
		m.articles[article.Slug] = article

		// Insert into database
		err = m.repository.InsertArticle(ctx, article)
		if err != nil {
			return fmt.Errorf("failed to insert article %s: %w", article.Slug, err)
		}

		return nil
	})
	return count, err
}

// parseMarkdownFile parses a markdown file and extracts metadata
func (m *Module) parseMarkdownFile(filePath string) (*model.Article, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)

	// Extract YAML front matter
	var meta model.Metadata

	// Check if file starts with ---
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
				return nil, fmt.Errorf("failed to parse YAML front matter: %w", err)
			}
		}
	}

	// Generate article ID and slug
	fileName := filepath.Base(filePath)
	slug := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	id := generateID(slug)
	now := time.Now()

	// Parse date
	var stamp *time.Time
	if metaDate, err := time.Parse("2006-01-02", meta.Date); err == nil {
		stamp = &metaDate
	}

	// Set default layout if not provided
	layout := meta.Layout
	if layout == "" {
		layout = "post"
	}

	article := &model.Article{
		ID:          id,
		Slug:        slug,
		Title:       meta.Title,
		Description: meta.Description,
		Filename:    filePath,
		Date:        stamp,
		OgImage:     meta.OgImage,
		Layout:      layout,
		Source:      meta.Source,
		URL:         "/blog/" + slug + "/",
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	return article, nil
}

// generateID creates a unique ID from slug
func generateID(slug string) string {
	return slug + "-" + time.Now().Format("20060102150405")
}
