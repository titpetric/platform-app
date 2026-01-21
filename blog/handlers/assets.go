package handlers

import (
	"io/fs"
	"net/http"

	"github.com/titpetric/platform"
)

// RegisterAssets registers static asset routes with the router
func RegisterAssets(r platform.Router, themeFS fs.FS) {
	assetFS := http.FileServer(http.FS(themeFS))

	r.Group(func(r platform.Router) {
		// CSS assets
		r.Get("/assets/css/*", func(w http.ResponseWriter, r *http.Request) {
			assetFS.ServeHTTP(w, r)
		})

		// Font assets
		r.Get("/assets/fonts/*", func(w http.ResponseWriter, r *http.Request) {
			assetFS.ServeHTTP(w, r)
		})

		// Icon assets
		r.Get("/assets/icons/*", func(w http.ResponseWriter, r *http.Request) {
			assetFS.ServeHTTP(w, r)
		})

		// Favicon assets
		r.Get("/assets/favicon/*", func(w http.ResponseWriter, r *http.Request) {
			assetFS.ServeHTTP(w, r)
		})

		// Robots.txt
		r.Get("/assets/robots.txt", func(w http.ResponseWriter, r *http.Request) {
			assetFS.ServeHTTP(w, r)
		})

		// Web manifest
		r.Get("/assets/site.webmanifest", func(w http.ResponseWriter, r *http.Request) {
			assetFS.ServeHTTP(w, r)
		})
	})
}
