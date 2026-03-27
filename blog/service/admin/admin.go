package admin

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/titpetric/platform"
	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/storage"
	"github.com/titpetric/platform-app/blog/view"
)

// Handlers provides HTTP handlers for blog admin endpoints.
type Handlers struct {
	repository *storage.Storage
	contentFS  *storage.GitFS
	views      *view.AdminViews
	themeFS    fs.FS
}

// NewHandlers returns a new Handlers instance.
func NewHandlers(repo *storage.Storage, contentFS *storage.GitFS, themeFS fs.FS) *Handlers {
	return &Handlers{
		repository: repo,
		contentFS:  contentFS,
		views:      view.NewAdminViews(themeFS),
		themeFS:    vuego.NewOverlayFS(themeFS, basecoat.Templates()),
	}
}

// Mount registers the admin routes on the given router.
func (h *Handlers) Mount(r platform.Router) {
	r.Group(func(r platform.Router) {
		// Admin HTML Routes
		r.Get("/admin/", h.DashboardHTML)
		r.Get("/admin/drafts", h.ListDraftsHTML)
		r.Get("/admin/scheduled", h.ListScheduledHTML)
		r.Get("/admin/published", h.ListPublishedHTML)
		r.Get("/admin/articles/{slug}", h.EditArticleHTML)
		r.Get("/admin/articles/{slug}/edit", h.EditArticleHTML)
		r.Get("/admin/new", h.NewArticleHTML)

		// Admin JSON API Routes
		r.Get("/admin/drafts.json", h.ListDraftsJSON)
		r.Get("/admin/scheduled.json", h.ListScheduledJSON)
		r.Get("/admin/published.json", h.ListPublishedJSON)
		r.Get("/admin/articles/{slug}.json", h.GetArticleJSON)

		r.Get("/admin/articles/{slug}/check", h.CheckSlugJSON)
		r.Post("/admin/articles", h.CreateArticleJSON)
		r.Put("/admin/articles/{slug}", h.UpdateArticleJSON)
		r.Delete("/admin/articles/{slug}", h.DeleteArticleJSON)
	})
}

// DashboardHTML renders the admin dashboard.
func (h *Handlers) DashboardHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.dashboardHTML(w, r))
}

func (h *Handlers) dashboardHTML(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	drafts, err := h.repository.CountDraftArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count drafts", err)
	}

	scheduled, err := h.repository.CountScheduledArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count scheduled", err)
	}

	published, err := h.repository.CountPublishedArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count published", err)
	}

	data := view.NewAdminDashboardData(drafts, scheduled, published)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.views.Dashboard(data).Render(ctx, w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// ListDraftsHTML renders the drafts list page.
func (h *Handlers) ListDraftsHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listDraftsHTML(w, r))
}

func (h *Handlers) listDraftsHTML(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page, pageSize := parsePagination(r)
	offset := (page - 1) * pageSize

	total, err := h.repository.CountDraftArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count drafts", err)
	}

	articles, err := h.repository.GetDraftArticles(ctx, offset, pageSize)
	if err != nil {
		return ErrInternal("failed to fetch drafts", err)
	}

	data := view.NewAdminListData("Drafts", articles, total, page, pageSize)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.views.List(data).Render(ctx, w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// ListScheduledHTML renders the scheduled articles list page.
func (h *Handlers) ListScheduledHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listScheduledHTML(w, r))
}

func (h *Handlers) listScheduledHTML(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page, pageSize := parsePagination(r)
	offset := (page - 1) * pageSize

	total, err := h.repository.CountScheduledArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count scheduled", err)
	}

	articles, err := h.repository.GetScheduledArticles(ctx, offset, pageSize)
	if err != nil {
		return ErrInternal("failed to fetch scheduled", err)
	}

	data := view.NewAdminListData("Scheduled", articles, total, page, pageSize)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.views.List(data).Render(ctx, w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// ListPublishedHTML renders the published articles list page.
func (h *Handlers) ListPublishedHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listPublishedHTML(w, r))
}

func (h *Handlers) listPublishedHTML(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page, pageSize := parsePagination(r)
	offset := (page - 1) * pageSize

	total, err := h.repository.CountPublishedArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count published", err)
	}

	articles, err := h.repository.GetPublishedArticles(ctx, offset, pageSize)
	if err != nil {
		return ErrInternal("failed to fetch published", err)
	}

	data := view.NewAdminListData("Published", articles, total, page, pageSize)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.views.List(data).Render(ctx, w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// EditArticleHTML renders the article edit form.
func (h *Handlers) EditArticleHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.editArticleHTML(w, r))
}

func (h *Handlers) editArticleHTML(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	slug := platform.URLParam(r, "slug")

	article, err := h.repository.GetArticleBySlug(ctx, slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	// Read the markdown content and strip frontmatter (metadata is from DB)
	content, err := h.contentFS.ReadFile(article.Filename)
	if err != nil {
		return ErrInternal("failed to read article content", err)
	}

	// Strip frontmatter and extract custom YAML
	bodyContent := view.StripFrontMatter(content)
	customYaml := view.ExtractCustomYAML(content)
	data := view.NewAdminEditData(article, string(bodyContent), customYaml)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.views.Edit(data).Render(ctx, w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// NewArticleHTML renders the new article form.
func (h *Handlers) NewArticleHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.newArticleHTML(w, r))
}

func (h *Handlers) newArticleHTML(w http.ResponseWriter, r *http.Request) error {
	data := view.NewAdminEditData(nil, "", "")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.views.Edit(data).Render(r.Context(), w); err != nil {
		return fmt.Errorf("render failed: %w", err)
	}
	return nil
}

// ListDraftsJSON returns drafts as JSON.
func (h *Handlers) ListDraftsJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listDraftsJSON(w, r))
}

func (h *Handlers) listDraftsJSON(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page, pageSize := parsePagination(r)
	offset := (page - 1) * pageSize

	total, err := h.repository.CountDraftArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count drafts", err)
	}

	articles, err := h.repository.GetDraftArticles(ctx, offset, pageSize)
	if err != nil {
		return ErrInternal("failed to fetch drafts", err)
	}

	return writeJSON(w, &model.ArticleList{
		Articles: articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// ListScheduledJSON returns scheduled articles as JSON.
func (h *Handlers) ListScheduledJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listScheduledJSON(w, r))
}

func (h *Handlers) listScheduledJSON(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page, pageSize := parsePagination(r)
	offset := (page - 1) * pageSize

	total, err := h.repository.CountScheduledArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count scheduled", err)
	}

	articles, err := h.repository.GetScheduledArticles(ctx, offset, pageSize)
	if err != nil {
		return ErrInternal("failed to fetch scheduled", err)
	}

	return writeJSON(w, &model.ArticleList{
		Articles: articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// ListPublishedJSON returns published articles as JSON.
func (h *Handlers) ListPublishedJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listPublishedJSON(w, r))
}

func (h *Handlers) listPublishedJSON(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page, pageSize := parsePagination(r)
	offset := (page - 1) * pageSize

	total, err := h.repository.CountPublishedArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count published", err)
	}

	articles, err := h.repository.GetPublishedArticles(ctx, offset, pageSize)
	if err != nil {
		return ErrInternal("failed to fetch published", err)
	}

	return writeJSON(w, &model.ArticleList{
		Articles: articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// CheckSlugJSON checks if a slug is available for use.
func (h *Handlers) CheckSlugJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.checkSlugJSON(w, r))
}

func (h *Handlers) checkSlugJSON(w http.ResponseWriter, r *http.Request) error {
	slug := platform.URLParam(r, "slug")

	_, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return writeJSON(w, map[string]bool{"available": true})
	}

	return writeJSON(w, map[string]bool{"available": false})
}

// GetArticleJSON returns a single article as JSON.
func (h *Handlers) GetArticleJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.getArticleJSON(w, r))
}

func (h *Handlers) getArticleJSON(w http.ResponseWriter, r *http.Request) error {
	slug := platform.URLParam(r, "slug")

	article, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	return writeJSON(w, article)
}

// CreateArticleJSON creates a new article.
func (h *Handlers) CreateArticleJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.createArticleJSON(w, r))
}

func (h *Handlers) createArticleJSON(w http.ResponseWriter, r *http.Request) error {
	var req ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return ErrBadRequest("invalid request body", err)
	}

	if err := req.Validate(); err != nil {
		return ErrBadRequest(err.Error(), nil)
	}

	article := req.ToArticle()

	// Write markdown file
	content := req.BuildMarkdownContent()
	filename := article.Slug + ".md"
	if err := h.contentFS.WriteFile(filename, []byte(content), 0o644, fmt.Sprintf("Create article: %s", article.Title)); err != nil {
		return ErrInternal("failed to write article file", err)
	}

	article.Filename = filename

	// Insert into database
	if err := h.repository.InsertArticle(r.Context(), article); err != nil {
		return ErrInternal("failed to create article", err)
	}

	w.WriteHeader(http.StatusCreated)
	return writeJSON(w, article)
}

// UpdateArticleJSON updates an existing article.
func (h *Handlers) UpdateArticleJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.updateArticleJSON(w, r))
}

func (h *Handlers) updateArticleJSON(w http.ResponseWriter, r *http.Request) error {
	slug := platform.URLParam(r, "slug")

	existing, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	var req ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return ErrBadRequest("invalid request body", err)
	}

	article := req.UpdateArticle(existing)

	// Write markdown file
	content := req.BuildMarkdownContent()
	if err := h.contentFS.WriteFile(existing.Filename, []byte(content), 0o644, fmt.Sprintf("Update article: %s", article.Title)); err != nil {
		return ErrInternal("failed to write article file", err)
	}

	// Update in database
	if err := h.repository.UpdateArticle(r.Context(), article); err != nil {
		return ErrInternal("failed to update article", err)
	}

	return writeJSON(w, article)
}

// DeleteArticleJSON deletes an article.
func (h *Handlers) DeleteArticleJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.deleteArticleJSON(w, r))
}

func (h *Handlers) deleteArticleJSON(w http.ResponseWriter, r *http.Request) error {
	slug := platform.URLParam(r, "slug")

	existing, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	// Remove markdown file
	if err := h.contentFS.Remove(existing.Filename, fmt.Sprintf("Delete article: %s", existing.Title)); err != nil {
		return ErrInternal("failed to remove article file", err)
	}

	// Delete from database
	if err := h.repository.DeleteArticle(r.Context(), slug); err != nil {
		return ErrInternal("failed to delete article", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func parsePagination(r *http.Request) (page, pageSize int) {
	page = 1
	pageSize = 10

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	return page, pageSize
}

func writeJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}
