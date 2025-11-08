package assets

import (
	"embed"
	"net/http"

	"github.com/titpetric/platform"
)

//go:embed css/* js/* fonts/*
var assets embed.FS

// Module is the assets module for serving embedded CSS and JS.
type Module struct {
	platform.UnimplementedModule
}

// Mount registers the assets routes.
func (Module) Mount(r platform.Router) error {
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/", http.FileServer(http.FS(assets))).ServeHTTP(w, r)
	})

	return nil
}
