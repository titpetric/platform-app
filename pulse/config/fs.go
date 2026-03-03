package config

import (
	"embed"
	"io/fs"
)

// fs contains embedded config files.
//
//go:embed all:data
var configFS embed.FS

// ConfigFS returns the embedded config filesystem.
func ConfigFS() fs.FS {
	return configFS
}
