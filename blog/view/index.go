package view

import (
	"github.com/titpetric/platform-app/blog/model"
)

// IndexData holds the data required for rendering the index page.
type IndexData struct {
	Articles []model.Article
	Data     map[string]any
	Total    int
}

// NewIndexData creates IndexData from a list of articles.
func NewIndexData(articles []model.Article) *IndexData {
	return &IndexData{
		Articles: articles,
		Data:     make(map[string]any),
		Total:    len(articles),
	}
}

// Map converts IndexData to a map for template rendering.
func (d *IndexData) Map() map[string]any {
	data := d.Data
	data["articles"] = d.Articles
	data["total"] = d.Total
	data["module"] = "blog"
	return data
}
