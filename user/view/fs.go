package view

import (
	"embed"
)

// FS contains embedded templates.
//
//go:embed *.vuego all:layouts
var FS embed.FS
