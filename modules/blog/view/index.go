package view

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/platform-app/modules/blog/model"
)

// IndexData holds the data required for rendering the index page
type IndexData struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	OGImage     string          `json:"ogImage"`
	Articles    []model.Article `json:"articles"`
	Total       int             `json:"total"`
}

// Map converts IndexData to a map[string]any
func (d *IndexData) Map() map[string]any {
	return map[string]any{
		"title":       d.Title,
		"description": d.Description,
		"ogImage":     d.OGImage,
		"articles":    d.Articles,
		"total":       d.Total,
	}
}

// Index renders the blog index/list page
func (v *Views) Index(ctx context.Context, w io.Writer, data *IndexData) error {
	// Build the context data
	templateData := data.Map()
	for k, v := range v.data {
		if _, ok := templateData[k]; !ok {
			templateData[k] = v
		}
	}

	// Render the index page
	return v.Render(ctx, w, "pages/index.vuego", templateData)
}

func fillTemplateData(w *map[string]any) error {
	if err := loadFile(w, "navigation", "appdata/config/navigation.json"); err != nil {
		return err
	}
	if err := loadFile(w, "themes", "appdata/config/themes.json"); err != nil {
		return err
	}
	if err := loadFileYaml(w, "meta", "appdata/config/meta.yml"); err != nil {
		return err
	}
	return nil
}

func loadFileYaml(w *map[string]any, key string, filename string) error {
	var result any
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", filename, err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&result); err != nil {
		return fmt.Errorf("error loading %s: %w", filename, err)
	}
	(*w)[key] = result
	log.Println("loaded ok:", key, filename)
	return nil
}

func loadFile(w *map[string]any, key string, filename string) error {
	var result any
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", filename, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&result); err != nil {
		return fmt.Errorf("error loading %s: %w", filename, err)
	}
	(*w)[key] = result
	log.Println("loaded ok:", key, filename)
	return nil
}

// Blog renders the blog list page
func (v *Views) Blog(ctx context.Context, w io.Writer, data *IndexData) error {
	// Build the context data
	templateData := data.Map()
	for k, v := range v.data {
		if _, ok := templateData[k]; !ok {
			templateData[k] = v
		}
	}

	var buf bytes.Buffer
	if err := v.Render(ctx, &buf, "pages/blog.vuego", templateData); err != nil {
		return err
	}
	_, err := io.Copy(w, &buf)
	return err
}

// IndexFromArticles creates IndexData from a list of articles
func (v *Views) IndexFromArticles(articles []model.Article) *IndexData {
	return &IndexData{
		Title:       "Blog",
		Description: "Read my latest articles and posts",
		Articles:    articles,
		Total:       len(articles),
	}
}
