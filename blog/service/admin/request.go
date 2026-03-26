package admin

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"

	"github.com/titpetric/platform-app/blog/model"
)

// ArticleRequest represents a request to create or update an article.
type ArticleRequest struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Date        string `json:"date"`
	Layout      string `json:"layout"`
	Draft       bool   `json:"draft"`
}

// Validation constants.
const (
	maxSlugLength        = 100
	maxTitleLength       = 200
	maxDescriptionLength = 500
	maxContentLength     = 100000
)

// Allowed layouts.
var allowedLayouts = map[string]bool{
	"post": true,
	"page": true,
}

var (
	slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

// Validate validates the request with strict input checking.
func (r *ArticleRequest) Validate() error {
	// Trim whitespace
	r.Slug = strings.TrimSpace(r.Slug)
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)
	r.Layout = strings.TrimSpace(r.Layout)
	r.Date = strings.TrimSpace(r.Date)

	// Slug validation
	if r.Slug == "" {
		return errors.New("slug is required")
	}
	if len(r.Slug) > maxSlugLength {
		return errors.New("slug exceeds maximum length")
	}
	if !slugRegex.MatchString(r.Slug) {
		return errors.New("slug must be lowercase alphanumeric with hyphens")
	}

	// Title validation
	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Title) > maxTitleLength {
		return errors.New("title exceeds maximum length")
	}

	// Description validation
	if len(r.Description) > maxDescriptionLength {
		return errors.New("description exceeds maximum length")
	}

	// Content validation
	if r.Content == "" {
		return errors.New("content is required")
	}
	if len(r.Content) > maxContentLength {
		return errors.New("content exceeds maximum length")
	}

	// Layout validation - must be from allowed set
	if r.Layout != "" && !allowedLayouts[r.Layout] {
		return errors.New("invalid layout")
	}

	// Date validation - must be YYYY-MM-DD format if provided
	if r.Date != "" {
		if !dateRegex.MatchString(r.Date) {
			return errors.New("date must be in YYYY-MM-DD format")
		}
		if _, err := time.Parse("2006-01-02", r.Date); err != nil {
			return errors.New("invalid date")
		}
	}

	// Sanitize content - remove script tags and event handlers
	r.Content = sanitizeContent(r.Content)

	return nil
}

// sanitizeContent removes potentially dangerous HTML content.
// For markdown content, we use a policy that allows safe HTML while
// stripping scripts, event handlers, and other XSS vectors.
func sanitizeContent(content string) string {
	// UGCPolicy allows common formatting tags but strips dangerous content
	p := bluemonday.UGCPolicy()
	return p.Sanitize(content)
}

// ToArticle converts the request to an Article model.
func (r *ArticleRequest) ToArticle() *model.Article {
	layout := r.Layout
	if layout == "" {
		layout = "post"
	}

	article := &model.Article{
		ID:          r.Slug + "-" + time.Now().Format("20060102150405"),
		Slug:        r.Slug,
		Title:       r.Title,
		Description: r.Description,
		Layout:      layout,
		URL:         "/blog/" + r.Slug + "/",
	}

	if r.Draft {
		article.Draft = 1
	}

	if r.Date != "" {
		if t, err := time.Parse("2006-01-02", r.Date); err == nil {
			article.Date = &t
		}
	}

	return article
}

// UpdateArticle updates an existing article with request data.
func (r *ArticleRequest) UpdateArticle(existing *model.Article) *model.Article {
	existing.Title = r.Title
	existing.Description = r.Description

	if r.Layout != "" {
		existing.Layout = r.Layout
	}

	if r.Draft {
		existing.Draft = 1
	} else {
		existing.Draft = 0
	}

	if r.Date != "" {
		if t, err := time.Parse("2006-01-02", r.Date); err == nil {
			existing.Date = &t
		}
	}

	return existing
}

// BuildMarkdownContent builds the markdown content with front matter.
func (r *ArticleRequest) BuildMarkdownContent() string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString("title: \"" + escapeYAML(r.Title) + "\"\n")

	if r.Description != "" {
		sb.WriteString("description: \"" + escapeYAML(r.Description) + "\"\n")
	}

	if r.Date != "" {
		sb.WriteString("date: \"" + r.Date + "\"\n")
	}

	if r.Layout != "" && r.Layout != "post" {
		sb.WriteString("layout: \"" + r.Layout + "\"\n")
	}

	if r.Draft {
		sb.WriteString("draft: true\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString(r.Content)

	return sb.String()
}

func escapeYAML(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
