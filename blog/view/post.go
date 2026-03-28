package view

import (
	"time"

	"github.com/titpetric/platform-app/blog/model"
)

// PostData holds the data required for rendering the post layout.
type PostData struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	OgImage     string     `json:"ogImage"`
	Content     string     `json:"content"`
	Date        *time.Time `json:"date"`
	Class       string     `json:"class"`
	LoggedIn    bool       `json:"loggedIn"`
}

// NewPostData creates PostData from an Article.
func NewPostData(article *model.Article, content string, loggedIn bool) *PostData {
	return &PostData{
		Slug:        article.Slug,
		Title:       article.Title,
		Description: article.Description,
		OgImage:     article.OgImage,
		Content:     content,
		Date:        article.Date,
		Class:       "prose",
		LoggedIn:    loggedIn,
	}
}

// Map converts PostData to a map[string]any.
func (d *PostData) Map() map[string]any {
	m := make(map[string]any)
	m["slug"] = d.Slug
	m["title"] = d.Title
	m["description"] = d.Description
	m["ogImage"] = d.OgImage
	m["content"] = d.Content
	m["date"] = d.Date
	m["class"] = d.Class
	m["module"] = "blog"
	m["loggedIn"] = d.LoggedIn
	m["page"] = map[string]any{
		"url": "/blog/" + d.Slug + "/",
	}
	return m
}
