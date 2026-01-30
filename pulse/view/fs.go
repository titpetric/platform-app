package view

import (
	"embed"
)

// FS contains embedded templates.
//
//go:embed *.vuego all:data
var FS embed.FS
