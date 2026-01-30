package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/pulse/storage"
)

type Handlers struct {
	Options
}

func NewHandlers(opts Options) *Handlers {
	return &Handlers{
		Options: opts,
	}
}

func (h *Handlers) Mount(r platform.Router) {
	r.Get("/pulse", h.IndexPage)
	r.Get("/pulse/{username}", h.UserPage)

	r.Post("/api/pulse/ingest", h.PostIngest)
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
	return nil
}

func (h *Handlers) UserPage(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(r.Context(), w, h.userPage(w, r))
}

func (h *Handlers) userPage(w http.ResponseWriter, r *http.Request) error {
	return nil
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

	repo, err := storage.NewStorage(ctx)
	if err != nil {
		return err
	}

	return repo.Pulse(ctx, body.Count)
}
