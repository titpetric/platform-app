package blog

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/titpetric/platform-app/modules/blog/handlers"
	"github.com/titpetric/platform-app/modules/blog/markdown"
	"github.com/titpetric/platform-app/modules/blog/view"
)

// Generator generates static HTML files from blog content
type Generator struct {
	module    *Module
	outputDir string
}

// NewGenerator creates a new Generator instance
func NewGenerator(m *Module, outputDir string) *Generator {
	return &Generator{
		module:    m,
		outputDir: outputDir,
	}
}

// Generate generates all static HTML files
func (g *Generator) Generate(ctx context.Context) error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy static assets
	fmt.Println("Copying assets...")
	if err := g.copyAssets(); err != nil {
		return fmt.Errorf("failed to copy assets: %w", err)
	}

	// Create handlers for rendering
	h, err := handlers.NewHandlers(g.module.repository, g.module.themeFS)
	if err != nil {
		return fmt.Errorf("failed to create handlers: %w", err)
	}

	// Generate index page
	fmt.Println("Generating index.html...")
	if err := g.generateIndexPage(ctx, h); err != nil {
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	// Generate other pages (404, resume, etc.)
	fmt.Println("Generating static pages...")
	if err := g.generateStaticPages(ctx, h); err != nil {
		return fmt.Errorf("failed to generate static pages: %w", err)
	}

	// Generate individual article pages
	articles, err := g.module.repository.GetArticles(ctx, 0, 9999)
	if err != nil {
		return fmt.Errorf("failed to fetch articles: %w", err)
	}

	mdRenderer := markdown.NewRenderer()
	for _, modelArticle := range articles {
		fmt.Printf("Generating blog/%s/index.html...\n", modelArticle.Slug)

		content, err := os.ReadFile(modelArticle.Filename)
		if err != nil {
			return err
		}

		// Convert markdown to HTML
		contentWithoutFrontMatter := view.StripFrontMatter(content)

		htmlContent := mdRenderer.Render(contentWithoutFrontMatter)

		// Create PostData
		postData := view.NewPostData(&modelArticle, string(htmlContent))

		if err := g.generateArticlePage(ctx, h, postData); err != nil {
			return fmt.Errorf("failed to generate article page for %s: %w", modelArticle.Slug, err)
		}
	}

	// Generate feed.xml
	fmt.Println("Generating feed.xml...")
	if err := g.generateFeed(ctx, h); err != nil {
		return fmt.Errorf("failed to generate feed: %w", err)
	}

	fmt.Printf("✓ Generated %d articles\n", len(articles))
	return nil
}

// generateIndexPage generates the index.html file
func (g *Generator) generateIndexPage(ctx context.Context, h *handlers.Handlers) error {
	articles, err := h.Repository().GetArticles(ctx, 0, 5)
	if err != nil {
		return err
	}

	indexData := view.NewIndexData(articles)

	var buf bytes.Buffer
	if err := h.Views().Index(indexData).Layout(ctx, &buf); err != nil {
		return err
	}

	indexPath := filepath.Join(g.outputDir, "index.html")
	return os.WriteFile(indexPath, buf.Bytes(), 0o644)
}

// generateStaticPages generates all .vuego pages from theme/pages directory recursively
func (g *Generator) generateStaticPages(ctx context.Context, h *handlers.Handlers) error {
	return g.walkPages(ctx, h, "pages", "")
}

func (g *Generator) walkPages(ctx context.Context, h *handlers.Handlers, dirPath string, relPath string) error {
	entries, err := fs.ReadDir(g.module.themeFS, dirPath)
	if err != nil {
		return fmt.Errorf("failed to read pages directory %s: %w", dirPath, err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		entryRelPath := filepath.Join(relPath, entry.Name())

		if entry.IsDir() {
			// Recursively walk subdirectories
			if err := g.walkPages(ctx, h, entryPath, entryRelPath); err != nil {
				return err
			}
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".vuego") {
			continue
		}

		// Skip index.vuego in root pages (handled separately)
		if entry.Name() == "index.vuego" && relPath == "" {
			continue
		}

		pageName := strings.TrimSuffix(entry.Name(), ".vuego")
		templatePath := filepath.Join("pages", strings.ReplaceAll(entryRelPath, string(filepath.Separator), "/"))

		// Determine output path and prepare data
		var outputPath string
		var templateData map[string]interface{}

		if entry.Name() == "index.vuego" {
			// index.vuego in subdirectories becomes subdir/index.html
			parentDir := strings.TrimSuffix(entryRelPath, string(filepath.Separator)+"index.vuego")
			outputDir := filepath.Join(g.outputDir, parentDir)
			if err := os.MkdirAll(outputDir, 0o755); err != nil {
				return err
			}
			outputPath = filepath.Join(outputDir, "index.html")

			// Special handling for blog/index.vuego
			if parentDir == "blog" || parentDir == "blog"+string(filepath.Separator) {
				articles, err := h.Repository().GetArticles(ctx, 0, 9999)
				if err != nil {
					return fmt.Errorf("failed to fetch articles for blog page: %w", err)
				}
				templateData = map[string]interface{}{
					"articles": articles,
					"total":    len(articles),
				}
			} else {
				templateData = map[string]interface{}{}
			}
		} else {
			// Regular pages become page-name.html
			outputPath = filepath.Join(g.outputDir, pageName+".html")

			// Special handling for blog.vuego
			if pageName == "blog" {
				articles, err := h.Repository().GetArticles(ctx, 0, 9999)
				if err != nil {
					return fmt.Errorf("failed to fetch articles for blog page: %w", err)
				}
				templateData = map[string]interface{}{
					"articles": articles,
					"total":    len(articles),
				}
			} else {
				templateData = map[string]interface{}{}
			}
		}

		// Render the page
		var buf bytes.Buffer
		if err := h.Views().Loader.Load(templatePath).Fill(templateData).Layout(ctx, &buf); err != nil {
			return fmt.Errorf("failed to render page %s: %w", templatePath, err)
		}

		// Write the file
		if err := os.WriteFile(outputPath, buf.Bytes(), 0o644); err != nil {
			return fmt.Errorf("failed to write page %s: %w", outputPath, err)
		}

		fmt.Printf("Generated %s\n", outputPath)
	}

	return nil
}

// generateArticlePage generates an individual article page
func (g *Generator) generateArticlePage(ctx context.Context, h *handlers.Handlers, postData *view.PostData) error {
	var buf bytes.Buffer
	if err := h.Views().Post(postData).Layout(ctx, &buf); err != nil {
		return err
	}

	articleDir := filepath.Join(g.outputDir, "blog", postData.Slug)
	if err := os.MkdirAll(articleDir, 0o755); err != nil {
		return err
	}

	articlePath := filepath.Join(articleDir, "index.html")
	return os.WriteFile(articlePath, buf.Bytes(), 0o644)
}

// generateFeed generates the feed.xml file
func (g *Generator) generateFeed(ctx context.Context, h *handlers.Handlers) error {
	articles, err := h.Repository().GetArticles(ctx, 0, 20)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := h.Views().AtomFeed(ctx, &buf, articles, nil); err != nil {
		return err
	}

	feedPath := filepath.Join(g.outputDir, "feed.xml")
	return os.WriteFile(feedPath, buf.Bytes(), 0o644)
}

// copyAssets copies static assets from theme/assets (both embedded and local) to output directory
func (g *Generator) copyAssets() error {
	assetsDestDir := filepath.Join(g.outputDir, "assets")

	// First, copy embedded theme assets using the overlay FS
	if err := g.copyEmbeddedAssets(assetsDestDir); err != nil {
		return err
	}

	// Then, copy local theme/assets if it exists (local overrides embedded)
	if err := g.copyLocalAssets(assetsDestDir); err != nil {
		return err
	}

	return nil
}

// copyEmbeddedAssets copies assets from the embedded theme filesystem
func (g *Generator) copyEmbeddedAssets(assetsDestDir string) error {
	return fs.WalkDir(g.module.themeFS, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path within assets
		relPath, err := filepath.Rel("assets", path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(assetsDestDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		// Read from embedded FS and write to dest
		data, err := fs.ReadFile(g.module.themeFS, path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0o644)
	})
}

// copyLocalAssets copies assets from local theme/assets directory, overriding embedded assets
func (g *Generator) copyLocalAssets(assetsDestDir string) error {
	assetsSrcDir := filepath.Join("theme", "assets")

	// Check if source exists
	if _, err := os.Stat(assetsSrcDir); os.IsNotExist(err) {
		// Assets directory doesn't exist, skip
		return nil
	}

	// Copy the entire assets directory
	return filepath.Walk(assetsSrcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(assetsSrcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(assetsDestDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0o644)
	})
}
