package service

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
	"github.com/titpetric/vuego"
	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/service/admin"
	"github.com/titpetric/platform-app/blog/service/api"
	"github.com/titpetric/platform-app/blog/service/web"
	"github.com/titpetric/platform-app/blog/storage"
)

// BlogModule implements the blog module for the platform.
type BlogModule struct {
	platform.UnimplementedModule

	// Data directory for markdown files
	dataDir string

	// Storage for database operations
	repository *storage.Storage

	// GitFS for content management with automatic commits
	contentFS *storage.GitFS

	// Articles index for in-memory access
	articles map[string]*model.Article

	mountFns []func(platform.Router)
}

// NewBlogModule creates a new blog module instance.
func NewBlogModule() *BlogModule {
	return &BlogModule{
		dataDir:  "./src",
		articles: make(map[string]*model.Article),
	}
}

// Name returns the module name.
func (m *BlogModule) Name() string {
	return "blog"
}

// Mount registers the blog routes with the router.
func (m *BlogModule) Mount(_ context.Context, r platform.Router) error {
	for _, mountFn := range m.mountFns {
		mountFn(r)
	}
	return nil
}

// Start initializes the blog module by scanning markdown files and building the index.
func (m *BlogModule) Start(ctx context.Context) error {
	// Create storage instance (includes migrations)
	var err error
	m.repository, err = storage.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	// Initialize GitFS for content directory
	m.contentFS, err = storage.NewGitFS(m.dataDir)
	if err != nil {
		return fmt.Errorf("failed to initialize content filesystem: %w", err)
	}
	fmt.Printf("[blog] initialized git content store at %s\n", m.contentFS.Root())

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

	return m.initHandlers(ctx)
}

func (m *BlogModule) initHandlers(ctx context.Context) error {
	// Check if local theme directory exists (for development)
	var customFS fs.FS
	if _, err := os.Stat("theme"); err == nil {
		customFS = os.DirFS("theme")
	}

	platformOpts := platform.OptionsFromContext(ctx)

	webFS := vuego.NewOverlayFS(customFS, FS(platformOpts.ConfigFS))
	adminFS, err := AdminFS(platformOpts.ConfigFS)
	if err != nil {
		return fmt.Errorf("failed initializing admin views: %w", err)
	}

	m.mountFns = []func(platform.Router){
		web.NewHandlers(m.repository, webFS).Mount,
		api.NewHandlers(m.repository).Mount,
		admin.NewHandlers(m.repository, m.contentFS, adminFS).Mount,
	}

	return nil
}

// Stop is called when the module is shutting down.
func (m *BlogModule) Stop(context.Context) error {
	// Nothing to clean up - database is managed by platform
	return nil
}

// SetRepository sets the repository on the module.
func (m *BlogModule) SetRepository(repo *storage.Storage) {
	m.repository = repo
}

// ScanMarkdownFiles scans the data directory for markdown files and indexes them.
// It returns the count of scanned files.
func (m *BlogModule) ScanMarkdownFiles(ctx context.Context) (int, error) {
	return m.scanMarkdownFiles(ctx)
}

// scanMarkdownFiles scans the data directory for markdown files and indexes them.
// It returns the count of scanned files.
func (m *BlogModule) scanMarkdownFiles(ctx context.Context) (int, error) {
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

// parseMarkdownFile parses a markdown file and extracts metadata.
func (m *BlogModule) parseMarkdownFile(filePath string) (*model.Article, error) {
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

	// Set draft status from metadata
	var draft int64
	if meta.Draft {
		draft = 1
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
		Draft:       draft,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	return article, nil
}

// generateID creates a unique ID from slug.
func generateID(slug string) string {
	return slug + "-" + time.Now().Format("20060102150405")
}
