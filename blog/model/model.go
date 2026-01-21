package model

// Metadata represents the YAML front matter of a markdown file
type Metadata struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	OgImage     string `yaml:"ogImage"`
	Date        string `yaml:"date"`
	Layout      string `yaml:"layout"`
	Source      string `yaml:"source"`
}

// ArticleList represents a paginated list of articles
type ArticleList struct {
	Articles []Article `json:"articles"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}
