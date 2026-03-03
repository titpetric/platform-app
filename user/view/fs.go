package view

import (
	"embed"
	"io/fs"
)

//go:embed *.vuego
var templateFS embed.FS

// Templates returns the embedded templates for the user module.
func Templates() fs.FS {
	return templateFS
}
