package schema

import (
	"embed"
)

// Migrations contains sql migrations contained in this folder.
//
//go:embed *.up.sql
var Migrations embed.FS
