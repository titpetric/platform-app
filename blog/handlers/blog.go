package handlers

import (
	"fmt"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform-app/blog/markdown"
	"github.com/titpetric/platform-app/blog/view"
)

// IndexHTML returns an HTML index page listing blogs
func (h *Handlers) IndexHTML(w http.ResponseWriter, r *http.Request) {
	articles, err := h.repository.GetArticles(r.Context(), 0, 5)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch articles: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")

	// Create index component to render list
	indexData := view.NewIndexData(articles)

	if err := h.views.Index(indexData).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render failed: %v", err), http.StatusInternalServerError)
	}
}

// ListArticlesHTML returns an HTML list of articles
func (h *Handlers) ListArticlesHTML(w http.ResponseWriter, r *http.Request) {
	articles, err := h.repository.GetArticles(r.Context(), 0, 9999)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch articles: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")

	// Create blog list and render
	blogData := view.NewIndexData(articles)

	if err := h.views.Blog(blogData).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render failed: %v", err), http.StatusInternalServerError)
	}
}

// GetArticleHTML returns a single article as HTML
func (h *Handlers) GetArticleHTML(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	content, err := os.ReadFile(article.Filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	contentWithoutFrontMatter := view.StripFrontMatter(content)
	mdRenderer := markdown.NewRenderer()
	htmlContent := mdRenderer.Render(contentWithoutFrontMatter)

	// Create PostData and render
	postData := view.NewPostData(article, string(htmlContent))

	if err := h.views.Post(postData).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render failed: %v", err), http.StatusInternalServerError)
	}
}

// GetAtomFeed returns an Atom XML feed of all articles
func (h *Handlers) GetAtomFeed(w http.ResponseWriter, r *http.Request) {
	articles, err := h.repository.GetArticles(r.Context(), 0, 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch articles: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	if err := h.views.AtomFeed(r.Context(), w, articles, nil); err != nil {
		http.Error(w, fmt.Sprintf("feed generation failed: %v", err), http.StatusInternalServerError)
	}
}
