package web

import (
	"fmt"
	"net/http"
	"strconv"
)

// ListArticlesAdminHTML renders the admin articles list as HTML using CMS layout
func (h *Handlers) ListArticlesAdminHTML(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, fmt.Sprintf("failed to count articles: %v", err), http.StatusInternalServerError)
		return
	}

	// Get paginated articles
	articles, err := h.repository.GetArticles(r.Context(), offset, pageSz)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch articles: %v", err), http.StatusInternalServerError)
		return
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
		http.Error(w, fmt.Sprintf("failed to render template: %v", err), http.StatusInternalServerError)
		return
	}
}
