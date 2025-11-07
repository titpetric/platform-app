package daily

import (
	"context"
	"encoding/json"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/modules/daily/model"
	"github.com/titpetric/platform-app/modules/daily/storage"
	"github.com/titpetric/platform-app/modules/daily/view"
)

type Module struct {
	platform.UnimplementedModule

	repository *storage.Storage
}

func NewModule() *Module {
	return &Module{}
}

func (*Module) Name() string {
	return "daily"
}

func (m *Module) Start(ctx context.Context) error {
	db, err := storage.DB(ctx)
	if err != nil {
		return err
	}

	if err := Migrate(ctx, db); err != nil {
		return err
	}

	m.repository = storage.NewStorage(db)
	return nil
}

func (m *Module) Mount(r platform.Router) error {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tasks, err := m.repository.List(ctx)
		if err != nil {
			platform.Error(w, r, 503, err)
		}
		view.Daily("Hello...", tasks).Render(ctx, w)
	})

	r.Post("/daily/save", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var todo model.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		created, err := m.repository.Add(ctx, todo)
		if err != nil {
			http.Error(w, "failed to add todo", http.StatusInternalServerError)
			return
		}

		platform.JSON(w, r, http.StatusOK, created)
	})

	r.Post("/daily/complete/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := m.repository.Delete(ctx, chi.URLParam(r, "id"))
		telemetry.CaptureError(ctx, err)

		w.WriteHeader(http.StatusOK)
	})

	return nil
}
