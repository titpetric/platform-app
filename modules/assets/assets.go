package assets

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/titpetric/platform"
)

//go:embed css/* js/*
var embeddedFiles embed.FS

// Module is the assets module for serving embedded CSS and JS.
type Module struct {
	platform.UnimplementedModule
}

// Mount registers the assets routes.
func (Module) Mount(r platform.Router) error {
	cssFS, err := fs.Sub(embeddedFiles, "css")
	if err != nil {
		return err
	}
	jsFS, err := fs.Sub(embeddedFiles, "js")
	if err != nil {
		return err
	}

	r.Get("/assets/css/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/css/", http.FileServer(http.FS(cssFS))).ServeHTTP(w, r)
	})

	r.Get("/assets/js/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/js/", http.FileServer(http.FS(jsFS))).ServeHTTP(w, r)
	})

	return nil
}
