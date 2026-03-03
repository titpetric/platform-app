package service

import (
	"context"
	"io/fs"

	"github.com/titpetric/platform"
	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"

	"github.com/titpetric/platform-app/user/view"
)

// FS returns a layered filesystem combining theme, view, and basecoat assets.
func FS(ctx context.Context) fs.FS {
	platformOpts := platform.OptionsFromContext(ctx)

	// Build FS layers: theme (app-level) > views (module-level) > basecoat (base theme)
	return vuego.NewOverlayFS(platformOpts.ConfigFS, view.FS, basecoat.FS)
}
