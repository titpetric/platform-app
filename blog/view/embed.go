package view

import (
	"embed"
	"io/fs"
)

//go:embed all:admin all:assets all:components all:layouts all:pages blog_admin_articles.vuego
var templateFS embed.FS

// Templates returns the embedded templates for Basecoat.
func Templates() fs.FS {
	return templateFS
}
