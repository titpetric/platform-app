package view

import (
	"embed"
)

// FS contains embedded templates.
//
//go:embed *.vuego all:data all:layouts
var FS embed.FS
