package service

import (
	"io/fs"

	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"

	"github.com/titpetric/platform-app/blog/view"
)

// FS returns a layered filesystem combining config and blog view assets.
func FS(configFS fs.FS) fs.FS {
	return vuego.NewOverlayFS(configFS, view.Templates())
}

// AdminFS roots the template FS in admin/ for custom layouts.
// It combines the config, view and basecoat filesystems.
func AdminFS(configFS fs.FS) (fs.FS, error) {
	adminFS, err := fs.Sub(view.Templates(), "admin")
	if err != nil {
		return nil, err
	}

	return vuego.NewOverlayFS(configFS, adminFS, basecoat.Templates()), nil
}
