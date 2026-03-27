package view

import (
	"github.com/titpetric/platform-app/blog/model"
)

// AdminDashboardData holds data for the admin dashboard page.
type AdminDashboardData struct {
	Title     string
	Drafts    int
	Scheduled int
	Published int
}

// NewAdminDashboardData creates AdminDashboardData with counts.
func NewAdminDashboardData(drafts, scheduled, published int) *AdminDashboardData {
	return &AdminDashboardData{
		Title:     "Blog Dashboard",
		Drafts:    drafts,
		Scheduled: scheduled,
		Published: published,
	}
}

// Map converts AdminDashboardData to a map[string]any.
func (d *AdminDashboardData) Map() map[string]any {
	return map[string]any{
		"title":     d.Title,
		"drafts":    d.Drafts,
		"scheduled": d.Scheduled,
		"published": d.Published,
	}
}

// AdminListData holds data for the admin article list pages.
type AdminListData struct {
	Title    string
	Articles []model.Article
	Total    int
	Page     int
	PageSize int
}

// NewAdminListData creates AdminListData for a paginated list.
func NewAdminListData(title string, articles []model.Article, total, page, pageSize int) *AdminListData {
	return &AdminListData{
		Title:    title,
		Articles: articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

// Map converts AdminListData to a map[string]any.
func (d *AdminListData) Map() map[string]any {
	return map[string]any{
		"title":    d.Title,
		"articles": d.Articles,
		"total":    d.Total,
		"page":     d.Page,
		"pageSize": d.PageSize,
	}
}

// TotalPages returns the total number of pages.
func (d *AdminListData) TotalPages() int {
	if d.PageSize == 0 {
		return 0
	}
	return (d.Total + d.PageSize - 1) / d.PageSize
}

// AdminEditData holds data for the article edit page.
type AdminEditData struct {
	Title      string
	Article    *model.Article
	Content    string
	CustomYaml string
	IsNew      bool
}

// NewAdminEditData creates AdminEditData for editing or creating an article.
// rawContent is the full file content including frontmatter (for custom YAML extraction).
// bodyContent is the content without frontmatter.
func NewAdminEditData(article *model.Article, bodyContent string, customYaml string) *AdminEditData {
	data := &AdminEditData{
		Article:    article,
		Content:    bodyContent,
		CustomYaml: customYaml,
		IsNew:      article == nil,
	}

	if article != nil {
		data.Title = "Edit: " + article.Title
	} else {
		data.Title = "New Article"
	}

	return data
}

// Map converts AdminEditData to a map[string]any.
func (d *AdminEditData) Map() map[string]any {
	data := map[string]any{
		"title":   d.Title,
		"content": d.Content,
		"isNew":   d.IsNew,
	}

	if d.Article != nil {
		data["article"] = d.Article
		data["slug"] = d.Article.Slug
		data["articleTitle"] = d.Article.Title
		data["description"] = d.Article.Description
		data["layout"] = d.Article.Layout
		data["draft"] = d.Article.Draft != 0
		if d.Article.Date != nil {
			data["date"] = d.Article.Date.UTC().Format("2006-01-02")
			data["time"] = d.Article.Date.UTC().Format("15:04:05")
		} else {
			data["date"] = ""
			data["time"] = ""
		}
		data["customYaml"] = d.CustomYaml
	} else {
		// Provide default values for new article form
		data["slug"] = ""
		data["articleTitle"] = ""
		data["description"] = ""
		data["layout"] = "post"
		data["draft"] = false
		data["date"] = ""
		data["time"] = ""
		data["customYaml"] = ""
	}

	return data
}
