package admin

import (
	"errors"
	"fmt"
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
	Time        string `json:"time"`
	Layout      string `json:"layout"`
	Draft       bool   `json:"draft"`
	CustomYaml  string `json:"customYaml"`
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
	slugRegex     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	datetimeRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2}(:\d{2})?)?$`)
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

	// Date validation - must be YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS format if provided
	if r.Date != "" {
		if !datetimeRegex.MatchString(r.Date) {
			return errors.New("date must be in YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS format")
		}
		if _, err := parseDateTime(r.Date); err != nil {
			return errors.New("invalid date/time")
		}
	}

	// Sanitize content - remove script tags and event handlers
	r.Content = sanitizeContent(r.Content)

	return nil
}

func parseDateTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", s)
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
		if t := r.parseCombinedDateTime(); t != nil {
			article.Date = t
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
		if t := r.parseCombinedDateTime(); t != nil {
			existing.Date = t
		}
	}

	return existing
}

// parseCombinedDateTime combines date and time fields into a single time.Time.
func (r *ArticleRequest) parseCombinedDateTime() *time.Time {
	if r.Date == "" {
		return nil
	}

	dateTime := r.Date
	if r.Time != "" {
		dateTime = r.Date + " " + r.Time
	} else {
		dateTime = r.Date + " 00:00:00"
	}

	if t, err := parseDateTime(dateTime); err == nil {
		return &t
	}
	return nil
}

// formatDateTimeForFrontmatter formats date+time for YAML frontmatter.
func (r *ArticleRequest) formatDateTimeForFrontmatter() string {
	if r.Date == "" {
		return ""
	}
	if r.Time != "" {
		return r.Date + " " + r.Time
	}
	return r.Date
}

// BuildMarkdownContent builds the markdown content with front matter.
func (r *ArticleRequest) BuildMarkdownContent() string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString("title: \"" + escapeYAML(r.Title) + "\"\n")

	if r.Description != "" {
		sb.WriteString("description: \"" + escapeYAML(r.Description) + "\"\n")
	}

	if dateStr := r.formatDateTimeForFrontmatter(); dateStr != "" {
		sb.WriteString("date: \"" + dateStr + "\"\n")
	}

	if r.Layout != "" && r.Layout != "post" {
		sb.WriteString("layout: \"" + r.Layout + "\"\n")
	}

	if r.Draft {
		sb.WriteString("draft: true\n")
	}

	// Add custom YAML fields if provided
	if r.CustomYaml != "" {
		customYaml := strings.TrimSpace(r.CustomYaml)
		if customYaml != "" {
			sb.WriteString(customYaml)
			if !strings.HasSuffix(customYaml, "\n") {
				sb.WriteString("\n")
			}
		}
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
