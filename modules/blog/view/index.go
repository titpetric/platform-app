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
	Articles []model.Article
	Data     map[string]any
	Total    int
}

// NewIndexData creates IndexData from a list of articles.
func NewIndexData(articles []model.Article) *IndexData {
	result := &IndexData{
		Articles: articles,
		Data:     make(map[string]any),
		Total:    len(articles),
	}
	result.Fill()
	return result
}

func (d *IndexData) Map() map[string]any {
	data := d.Data
	data["articles"] = d.Articles
	data["total"] = d.Total
	data["module"] = "user"
	return data
}

func (d *IndexData) Fill() {
	if err := loadFileYaml(&d.Data, "config/meta.yml"); err != nil {
		log.Println("warn:", err)
	}
	if err := loadFile(&d.Data, "themes", "config/themes.json"); err != nil {
		log.Println("warn:", err)
	}
}

func loadFileYaml(w *map[string]any, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", filename, err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(w); err != nil {
		return fmt.Errorf("error loading %s: %w", filename, err)
	}
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
	return nil
}
