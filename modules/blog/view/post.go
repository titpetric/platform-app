package view

import (
	"time"

	"github.com/titpetric/platform-app/modules/blog/model"
)

// PostData holds the data required for rendering the post layout
type PostData struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	OgImage     string     `json:"ogImage"`
	Content     string     `json:"content"`
	Date        *time.Time `json:"date"`
	Classnames  string     `json:"classnames"`
}

// NewPostData creates PostData from an Article
func NewPostData(article *model.Article, content string) *PostData {
	return &PostData{
		Slug:        article.Slug,
		Title:       article.Title,
		Description: article.Description,
		OgImage:     article.OgImage,
		Content:     content,
		Date:        article.Date,
		Classnames:  "prose",
	}
}

// Map converts PostData to a map[string]any
func (d *PostData) Map() map[string]any {
	return map[string]any{
		"slug":        d.Slug,
		"title":       d.Title,
		"description": d.Description,
		"ogImage":     d.OgImage,
		"content":     d.Content,
		"date":        d.Date,
		"classnames":  d.Classnames,
	}
}
