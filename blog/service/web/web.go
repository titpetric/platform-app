package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/blog/markdown"
	"github.com/titpetric/platform-app/blog/storage"
	"github.com/titpetric/platform-app/blog/view"
)

// Handlers provides HTTP handlers for blog web endpoints.
type Handlers struct {
	repository *storage.Storage
	views      *view.Views
	themeFS    fs.FS
}

// NewHandlers returns a new Handlers instance.
func NewHandlers(repo *storage.Storage, themeFS fs.FS) *Handlers {
	return &Handlers{
		repository: repo,
		views:      view.NewViews(themeFS),
		themeFS:    themeFS,
	}
}

// Repository returns the storage repository
func (h *Handlers) Repository() *storage.Storage {
	return h.repository
}

// Views returns the views instance
func (h *Handlers) Views() *view.Views {
	return h.views
}

// Mount registers the blog web routes on the given router.
func (h *Handlers) Mount(r platform.Router) {
	// Register static assets
	h.registerAssets(r)

	r.Group(func(r platform.Router) {
		// Public HTML Routes
		r.Get("/", h.IndexHTML)
		r.Get("/blog", h.ListArticlesHTML)
		r.Get("/blog/", h.ListArticlesHTML)
		r.Get("/blog/{slug}", h.GetArticleHTML)
		r.Get("/blog/{slug}/", h.GetArticleHTML)

		// Feed Routes
		r.Get("/feed.xml", h.GetAtomFeed)

		// Admin HTML Routes
		r.Get("/admin/articles", h.ListArticlesAdminHTML)
	})
}

// IndexHTML returns an HTML index page listing blogs
func (h *Handlers) IndexHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.indexHTML(w, r))
}

func (h *Handlers) indexHTML(w http.ResponseWriter, r *http.Request) error {
	articles, err := h.repository.GetArticles(r.Context(), 0, 5)
	if err != nil {
		return ErrInternal("failed to fetch articles", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")

	// Create index component to render list
	indexData := view.NewIndexData(articles)

	if err := h.views.Index(indexData).Render(r.Context(), w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// ListArticlesHTML returns an HTML list of articles
func (h *Handlers) ListArticlesHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listArticlesHTML(w, r))
}

func (h *Handlers) listArticlesHTML(w http.ResponseWriter, r *http.Request) error {
	articles, err := h.repository.GetArticles(r.Context(), 0, 9999)
	if err != nil {
		return ErrInternal("failed to fetch articles", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")

	// Create blog list and render
	blogData := view.NewIndexData(articles)

	if err := h.views.Blog(blogData).Render(r.Context(), w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// GetArticleHTML returns a single article as HTML
func (h *Handlers) GetArticleHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.getArticleHTML(w, r))
}

func (h *Handlers) getArticleHTML(w http.ResponseWriter, r *http.Request) error {
	slug := platform.URLParam(r, "slug")

	article, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	content, err := os.ReadFile(article.Filename)
	if err != nil {
		return ErrNotFound("article content not found", err)
	}

	contentWithoutFrontMatter := view.StripFrontMatter(content)
	mdRenderer := markdown.NewRenderer()
	htmlContent := mdRenderer.Render(contentWithoutFrontMatter)

	// Create PostData and render
	postData := view.NewPostData(article, string(htmlContent))

	if err := h.views.Post(postData).Render(r.Context(), w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// GetAtomFeed returns an Atom XML feed of all articles
func (h *Handlers) GetAtomFeed(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.getAtomFeed(w, r))
}

func (h *Handlers) getAtomFeed(w http.ResponseWriter, r *http.Request) error {
	articles, err := h.repository.GetArticles(r.Context(), 0, 20)
	if err != nil {
		return ErrInternal("failed to fetch articles", err)
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	if err := h.views.AtomFeed(r.Context(), w, articles, nil); err != nil {
		return fmt.Errorf("feed generation failed: %w", err)
	}
	return nil
}

// ListArticlesAdminHTML renders the admin articles list as HTML using CMS layout
func (h *Handlers) ListArticlesAdminHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listArticlesAdminHTML(w, r))
}

func (h *Handlers) listArticlesAdminHTML(w http.ResponseWriter, r *http.Request) error {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	// Default pagination values
	pageNum := 1
	pageSz := 10

	// Parse and validate page parameter
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}

	// Parse and validate pageSize parameter
	if pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			pageSz = ps
		}
	}

	// Calculate offset
	offset := (pageNum - 1) * pageSz

	// Get total count
	total, err := h.repository.CountArticles(r.Context())
	if err != nil {
		return ErrInternal("failed to count articles", err)
	}

	// Get paginated articles
	articles, err := h.repository.GetArticles(r.Context(), offset, pageSz)
	if err != nil {
		return ErrInternal("failed to fetch articles", err)
	}

	// Prepare data for template rendering
	data := map[string]interface{}{
		"title":       "Articles - Blog Admin",
		"total":       total,
		"page":        pageNum,
		"pageSize":    pageSz,
		"breadcrumbs": []map[string]string{},
		"articles":    articles,
	}

	// Render using layout renderer with blog_admin_articles.vuego
	err = h.views.Loader.Load("blog_admin_articles.vuego").Fill(data).Render(r.Context(), w)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	return nil
}
