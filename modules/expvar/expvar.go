package expvar

import (
	"context"
	"expvar"
	"os"
	"sync"
	"time"

	"github.com/titpetric/platform"
)

var publishMu sync.Mutex

type Handler struct {
	platform.UnimplementedModule

	Enabled bool
}

func NewHandler() *Handler {
	return &Handler{
		Enabled: os.Getenv("PLATFORM_ENABLE_EXPVAR") == "true",
	}
}

func (m *Handler) Name() string {
	return "expvar"
}

func (m *Handler) Start(context.Context) error {
	if !m.Enabled {
		return nil
	}

	publishMu.Lock()
	defer publishMu.Unlock()

	start := time.Now()
	if expvar.Get("uptime") == nil {
		expvar.Publish("uptime", expvar.Func(func() interface{} {
			return time.Since(start).Seconds()
		}))
	}
	return nil
}

func (m *Handler) Mount(_ context.Context, r platform.Router) error {
	if !m.Enabled {
		return nil
	}

	r.Get("/debug/vars", expvar.Handler().ServeHTTP)
	return nil
}
