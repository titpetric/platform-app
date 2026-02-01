package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"
	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"

	"github.com/titpetric/platform-app/pulse/storage"
	"github.com/titpetric/platform-app/pulse/view"
	"github.com/titpetric/platform-app/user"
)

type Handlers struct {
	storage *storage.Storage
	vuego   vuego.Template
	fs      fs.FS
}

func NewHandlers(storage *storage.Storage) *Handlers {
	ofs := vuego.NewOverlayFS(view.FS, basecoat.FS)

	return &Handlers{
		fs:      ofs,
		storage: storage,
		vuego:   vuego.NewFS(ofs),
	}
}

func (h *Handlers) Mount(r platform.Router) {
	r.Get("/assets/*", http.FileServer(http.FS(h.fs)).ServeHTTP)

	r.Get("/pulse", h.IndexPage)
	r.Get("/pulse/{username}", h.UserPage)

	r.Group(func(r platform.Router) {
		r.Use(user.NewMiddleware(user.AuthHeader()))
		r.Post("/api/pulse/ingest", h.PostIngest)
	})
}

func (h *Handlers) errorHandler(ctx context.Context, w http.ResponseWriter, err error) {
	if err != nil {
		telemetry.CaptureError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) IndexPage(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(r.Context(), w, h.indexPage(w, r))
}

func (h *Handlers) indexPage(w http.ResponseWriter, r *http.Request) error {
	type viewData struct {
		Menu []any `json:"menu"`
	}

	ctx := r.Context()
	data := viewData{
		Menu: []any{},
	}

	indexPage := vuego.View[viewData](h.vuego, "index.vuego", data)

	return indexPage.Render(ctx, w)
}

func (h *Handlers) UserPage(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(r.Context(), w, h.userPage(w, r))
}

func (h *Handlers) userPage(w http.ResponseWriter, r *http.Request) error {
	type viewData struct{}

	ctx := r.Context()
	data := viewData{}

	indexPage := vuego.View[viewData](h.vuego, "user.vuego", data)

	return indexPage.Render(ctx, w)
}

func (h *Handlers) PostIngest(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(r.Context(), w, h.postIngest(w, r))
}

func (h *Handlers) postIngest(w http.ResponseWriter, r *http.Request) error {
	type ingestBody struct {
		Count int64
	}

	body := ingestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return err
	}

	if body.Count <= 0 {
		return fmt.Errorf("count must be positive: %d", body.Count)
	}

	ctx := r.Context()
	return h.storage.Pulse(ctx, body.Count)
}
