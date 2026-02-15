package view

import (
	"embed"
)

// FS contains embedded templates.
//
//go:embed *.vuego all:data all:layouts all:assets
var FS embed.FS
