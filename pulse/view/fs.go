package view

import (
	"embed"
	"io/fs"
)

//go:embed *.vuego
var templateFS embed.FS

// Templates returns the embedded pulse templates.
func Templates() fs.FS {
	return templateFS
}
