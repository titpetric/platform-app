package admin

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/titpetric/platform"
	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/schema"
	"github.com/titpetric/platform-app/blog/storage"
	"github.com/titpetric/platform-app/blog/view"
	"github.com/titpetric/platform-app/user"
)

// adminSlugPattern matches lowercase alphanumeric slugs with hyphens.
var adminSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// isValidSlugAdmin validates that a slug is well-formed.
func isValidSlugAdmin(slug string) bool {
	if slug == "" || len(slug) > maxSlugLength {
		return false
	}
	return adminSlugPattern.MatchString(slug)
}

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

// requireLoginRedirect is middleware that checks for session cookie and redirects to login if missing.
// It must run BEFORE the user middleware since that middleware short-circuits on auth failure.
func requireLoginRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Mount registers the admin routes on the given router.
func (h *Handlers) Mount(r platform.Router) {
	r.Group(func(r platform.Router) {
		// Check for session cookie first and redirect if missing
		r.Use(requireLoginRedirect)
		// Then load session data into context
		r.Use(user.NewMiddleware(user.AuthCookie()))

		// Admin HTML Routes
		r.Get("/admin", h.DashboardHTML)
		r.Get("/admin/blog/drafts", h.ListDraftsHTML)
		r.Get("/admin/blog/scheduled", h.ListScheduledHTML)
		r.Get("/admin/blog/published", h.ListPublishedHTML)
		r.Get("/admin/blog/articles/{slug}", h.EditArticleHTML)
		r.Get("/admin/blog/articles/{slug}/edit", h.EditArticleHTML)
		r.Get("/admin/blog/new", h.NewArticleHTML)
		r.Get("/admin/blog/settings", h.SettingsHTML)

		// Admin JSON API Routes (grouped under /api/admin)
		r.Get("/api/admin/blog/drafts", h.ListDraftsJSON)
		r.Get("/api/admin/blog/scheduled", h.ListScheduledJSON)
		r.Get("/api/admin/blog/published", h.ListPublishedJSON)
		r.Get("/api/admin/blog/articles/{slug}", h.GetArticleJSON)
		r.Get("/api/admin/blog/articles/{slug}/check", h.CheckSlugJSON)
		r.Post("/api/admin/blog/articles", h.CreateArticleJSON)
		r.Put("/api/admin/blog/articles/{slug}", h.UpdateArticleJSON)
		r.Delete("/api/admin/blog/articles/{slug}", h.DeleteArticleJSON)
		r.Post("/api/admin/blog/articles/{slug}/publish", h.PublishArticleJSON)

		// Settings API
		r.Get("/api/admin/blog/settings", h.GetSettingsJSON)
		r.Get("/api/admin/blog/settings/schema", h.GetSettingsSchemaJSON)
		r.Post("/api/admin/blog/settings", h.SaveSettingsJSON)
	})
}

// DashboardHTML renders the admin dashboard.
func (h *Handlers) DashboardHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.dashboardHTML(w, r))
}

func (h *Handlers) dashboardHTML(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	draftCount, err := h.repository.CountDraftArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count drafts", err)
	}

	scheduledCount, err := h.repository.CountScheduledArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count scheduled", err)
	}

	publishedCount, err := h.repository.CountPublishedArticles(ctx)
	if err != nil {
		return ErrInternal("failed to count published", err)
	}

	draftArticles, err := h.repository.GetDraftArticles(ctx, 0, 9999)
	if err != nil {
		return ErrInternal("failed to fetch drafts", err)
	}

	scheduledArticles, err := h.repository.GetScheduledArticles(ctx, 0, 9999)
	if err != nil {
		return ErrInternal("failed to fetch scheduled", err)
	}

	publishedArticles, err := h.repository.GetPublishedArticles(ctx, 0, 10)
	if err != nil {
		return ErrInternal("failed to fetch published", err)
	}

	data := view.NewAdminDashboardData(draftCount, scheduledCount, publishedCount, draftArticles, scheduledArticles, publishedArticles)

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
	if !isValidSlugAdmin(slug) {
		return ErrNotFound("article not found", nil)
	}

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
	if !isValidSlugAdmin(slug) {
		return ErrBadRequest("invalid slug format", nil)
	}

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
	if !isValidSlugAdmin(slug) {
		return ErrBadRequest("invalid slug", nil)
	}

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

	// Prevent silent overwrite: storage uses INSERT OR REPLACE so a duplicate
	// slug would otherwise clobber the existing row.
	if _, err := h.repository.GetArticleBySlug(r.Context(), req.Slug); err == nil {
		return ErrConflict("article with this slug already exists", nil)
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
	if !isValidSlugAdmin(slug) {
		return ErrBadRequest("invalid slug", nil)
	}

	existing, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	var req ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return ErrBadRequest("invalid request body", err)
	}

	// Path slug overrides body slug for updates; this ensures the URL slug is authoritative
	req.Slug = existing.Slug

	if err := req.Validate(); err != nil {
		return ErrBadRequest(err.Error(), nil)
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
	if !isValidSlugAdmin(slug) {
		return ErrBadRequest("invalid slug", nil)
	}

	existing, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	// Remove markdown file (best-effort: if file is already gone, continue)
	if existing.Filename != "" {
		if removeErr := h.contentFS.Remove(existing.Filename, fmt.Sprintf("Delete article: %s", existing.Title)); removeErr != nil {
			// Only fail if file existed but couldn't be removed
			if _, statErr := h.contentFS.Stat(existing.Filename); statErr == nil {
				return ErrInternal("failed to remove article file", removeErr)
			}
		}
	}

	// Delete from database
	if err := h.repository.DeleteArticle(r.Context(), slug); err != nil {
		return ErrInternal("failed to delete article", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// PublishArticleJSON sets draft=0 on an article, publishing it immediately or scheduling it.
func (h *Handlers) PublishArticleJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.publishArticleJSON(w, r))
}

func (h *Handlers) publishArticleJSON(w http.ResponseWriter, r *http.Request) error {
	slug := platform.URLParam(r, "slug")
	if !isValidSlugAdmin(slug) {
		return ErrBadRequest("invalid slug", nil)
	}

	article, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	article.Draft = 0

	// Update markdown frontmatter to reflect publish status (drop `draft: true`)
	if article.Filename != "" {
		content, readErr := h.contentFS.ReadFile(article.Filename)
		if readErr == nil {
			updated := removeDraftFromFrontmatter(content)
			if writeErr := h.contentFS.WriteFile(article.Filename, updated, 0o644, fmt.Sprintf("Publish article: %s", article.Title)); writeErr != nil {
				return ErrInternal("failed to update article file", writeErr)
			}
		}
	}

	if err := h.repository.UpdateArticle(r.Context(), article); err != nil {
		return ErrInternal("failed to publish article", err)
	}

	return writeJSON(w, article)
}

// removeDraftFromFrontmatter strips a `draft: true` (or similar) line from the
// YAML frontmatter while preserving the rest of the document.
func removeDraftFromFrontmatter(content []byte) []byte {
	const marker = "---"
	str := string(content)
	if !strings.HasPrefix(str, marker) {
		return content
	}
	rest := str[len(marker):]
	end := strings.Index(rest, marker)
	if end < 0 {
		return content
	}
	fm := rest[:end]
	body := rest[end:]

	var keep []string
	for _, line := range strings.Split(fm, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "draft:") {
			continue
		}
		keep = append(keep, line)
	}
	return []byte(marker + strings.Join(keep, "\n") + body)
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

// SettingsHTML renders the settings admin page.
func (h *Handlers) SettingsHTML(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.settingsHTML(w, r))
}

func (h *Handlers) settingsHTML(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// Get global settings (or create default)
	settings, err := h.repository.GetGlobalSettings(ctx)
	if err != nil {
		// Return empty settings if none exist
		settings = &model.Setting{
			UserID:       "global",
			MetaLang:     "en",
			PostsPerPage: 10,
			FeatureRss:   1,
		}
	}

	// Get schema for settings table
	table, err := schema.GetTable("setting")
	if err != nil {
		return ErrInternal("failed to load schema", err)
	}

	data := map[string]any{
		"title":    "Settings",
		"settings": settings,
		"schema":   table,
		"loggedIn": true,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return h.views.Loader.Load("settings.vuego").Fill(data).Render(ctx, w)
}

// GetSettingsJSON returns global settings as JSON.
func (h *Handlers) GetSettingsJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.getSettingsJSON(w, r))
}

func (h *Handlers) getSettingsJSON(w http.ResponseWriter, r *http.Request) error {
	settings, err := h.repository.GetGlobalSettings(r.Context())
	if err != nil {
		// Return empty settings if none exist
		settings = &model.Setting{UserID: "global"}
	}

	return writeJSON(w, settings)
}

// GetSettingsSchemaJSON returns the schema for the settings table.
func (h *Handlers) GetSettingsSchemaJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.getSettingsSchemaJSON(w, r))
}

func (h *Handlers) getSettingsSchemaJSON(w http.ResponseWriter, r *http.Request) error {
	table, err := schema.GetTable("setting")
	if err != nil {
		return ErrInternal("failed to load schema", err)
	}

	return writeJSON(w, table)
}

// SaveSettingsJSON saves settings from JSON request.
func (h *Handlers) SaveSettingsJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.saveSettingsJSON(w, r))
}

func (h *Handlers) saveSettingsJSON(w http.ResponseWriter, r *http.Request) error {
	var settings model.Setting

	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		return ErrBadRequest("invalid JSON", err)
	}

	if err := validateSettings(&settings); err != nil {
		return ErrBadRequest(err.Error(), nil)
	}

	// Always save as global settings for now
	settings.UserID = "global"

	if err := h.repository.SaveSetting(r.Context(), &settings); err != nil {
		return ErrInternal("failed to save settings", err)
	}

	return writeJSON(w, map[string]any{
		"success": true,
	})
}

// validateSettings ensures settings have reasonable values.
func validateSettings(s *model.Setting) error {
	s.MetaLang = strings.TrimSpace(s.MetaLang)
	s.MetaURL = strings.TrimSpace(s.MetaURL)
	s.SocialGithub = strings.TrimSpace(s.SocialGithub)
	s.SocialTwitter = strings.TrimSpace(s.SocialTwitter)
	s.SocialLinkedin = strings.TrimSpace(s.SocialLinkedin)
	s.AnalyticsID = strings.TrimSpace(s.AnalyticsID)

	if s.PostsPerPage < 0 {
		return fmt.Errorf("posts_per_page must be non-negative")
	}
	if s.PostsPerPage > 1000 {
		return fmt.Errorf("posts_per_page exceeds maximum of 1000")
	}

	const maxFieldLen = 500
	if len(s.MetaLang) > 10 {
		return fmt.Errorf("meta_lang exceeds maximum length")
	}
	if len(s.MetaURL) > maxFieldLen {
		return fmt.Errorf("meta_url exceeds maximum length")
	}
	if len(s.MetaAuthorName) > maxFieldLen {
		return fmt.Errorf("meta_author_name exceeds maximum length")
	}
	if len(s.MetaSubtitle) > maxFieldLen {
		return fmt.Errorf("meta_subtitle exceeds maximum length")
	}
	if s.MetaURL != "" && !strings.HasPrefix(s.MetaURL, "http://") && !strings.HasPrefix(s.MetaURL, "https://") {
		return fmt.Errorf("meta_url must start with http:// or https://")
	}
	return nil
}
