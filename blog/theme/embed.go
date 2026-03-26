package theme

import (
	"embed"
	"io/fs"
)

//go:embed all:assets all:components all:layouts
var templateFS embed.FS

// Templates returns the embedded templates for Basecoat.
func Templates() fs.FS {
	return templateFS
}
