package view

import (
	"encoding/json"
	"fmt"
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

// NewIndexData creates IndexData from a list of articles.
func NewIndexData(articles []model.Article) *IndexData {
	return &IndexData{
		Title:       "Blog",
		Description: "Read my latest articles and posts",
		Articles:    articles,
		Total:       len(articles),
	}
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

func fillTemplateData(w *map[string]any) error {
	if err := loadFile(w, "navigation", "appdata/config/navigation.json"); err != nil {
		log.Println("warn:", err)
	}
	if err := loadFile(w, "themes", "appdata/config/themes.json"); err != nil {
		log.Println("warn:", err)
	}
	if err := loadFileYaml(w, "meta", "appdata/config/meta.yml"); err != nil {
		log.Println("warn:", err)
	}

	// missing files aren't fatal, technically.
	// problem with tests is that the working dir changes
	// so we should swallow errors here.
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
