package service

import (
	"context"
	"io/fs"

	"github.com/titpetric/platform"
	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"
)

// FS returns a layered filesystem combining theme, view, and basecoat assets.
func FS(ctx context.Context) fs.FS {
	platformOpts := platform.OptionsFromContext(ctx)

	return vuego.NewOverlayFS(platformOpts.ConfigFS, basecoat.Templates())
}
