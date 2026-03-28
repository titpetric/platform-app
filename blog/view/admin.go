package view

import (
	"strings"
	"time"

	"github.com/titpetric/platform-app/blog/model"
)

// AdminDashboardData holds data for the admin dashboard page.
type AdminDashboardData struct {
	Title          string `json:"title"`
	DraftsCount    int    `json:"draftsCount"`
	ScheduledCount int    `json:"scheduledCount"`
	PublishedCount int    `json:"publishedCount"`

	Drafts    []model.Article `json:"drafts"`
	Scheduled []model.Article `json:"scheduled"`
	Published []model.Article `json:"publishes"`

	LoggedIn bool `json:"loggedIn"`
}

// NewAdminDashboardData creates AdminDashboardData with counts and article lists.
func NewAdminDashboardData(draftsCount, scheduledCount, publishedCount int, drafts, scheduled, published []model.Article) *AdminDashboardData {
	return &AdminDashboardData{
		Title:          "Blog Dashboard",
		DraftsCount:    draftsCount,
		ScheduledCount: scheduledCount,
		PublishedCount: publishedCount,
		Drafts:         drafts,
		Scheduled:      scheduled,
		Published:      published,
		LoggedIn:       true,
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
		"loggedIn": true, // Admin area requires login
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
		Content:    strings.TrimSpace(bodyContent),
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
		"title":       d.Title,
		"bodyContent": d.Content,
		"isNew":       d.IsNew,
	}

	utcdate := time.Now().UTC()

	if d.Article != nil {
		data["article"] = d.Article
		data["slug"] = d.Article.Slug
		data["articleTitle"] = d.Article.Title
		data["description"] = d.Article.Description
		data["layout"] = d.Article.Layout
		data["draft"] = d.Article.Draft != 0

		if d.Article.Date != nil {
			utcdate = d.Article.Date.UTC()
		}

		data["customYaml"] = d.CustomYaml
	} else {
		// Provide default values for new article form
		data["slug"] = ""
		data["articleTitle"] = ""
		data["description"] = ""
		data["layout"] = "post"
		data["draft"] = false
		data["customYaml"] = ""
	}

	data["date"] = utcdate.Format(time.DateTime)
	data["time"] = utcdate.Format("15:04:05")
	data["loggedIn"] = true // Admin area requires login

	return data
}
