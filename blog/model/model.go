package model

import "time"

// Metadata represents the YAML front matter of a markdown file.
type Metadata struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	OgImage     string `yaml:"ogImage"`
	Date        string `yaml:"date"`
	Layout      string `yaml:"layout"`
	Source      string `yaml:"source"`
	Draft       bool   `yaml:"draft"`
}

// ArticleList represents a paginated list of articles.
type ArticleList struct {
	Articles []Article `json:"articles"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

// ArticleStatus represents the publication status of an article.
type ArticleStatus string

// Article status constants.
const (
	// StatusDraft indicates an unpublished draft.
	StatusDraft ArticleStatus = "draft"
	// StatusScheduled indicates a scheduled article (future date).
	StatusScheduled ArticleStatus = "scheduled"
	// StatusPublished indicates a published article.
	StatusPublished ArticleStatus = "published"
)

// IsDraft returns true if the article is marked as a draft.
func (a *Article) IsDraft() bool {
	return a.Draft != 0
}

// IsScheduled returns true if the article has a future publication date.
func (a *Article) IsScheduled() bool {
	if a.IsDraft() {
		return false
	}
	if a.Date == nil {
		return false
	}
	return a.Date.After(time.Now())
}

// IsPublished returns true if the article is published (not draft, date <= now).
func (a *Article) IsPublished() bool {
	if a.IsDraft() {
		return false
	}
	if a.Date == nil {
		return true
	}
	return !a.Date.After(time.Now())
}

// Status returns the current publication status of the article.
func (a *Article) Status() ArticleStatus {
	if a.IsDraft() {
		return StatusDraft
	}
	if a.IsScheduled() {
		return StatusScheduled
	}
	return StatusPublished
}
