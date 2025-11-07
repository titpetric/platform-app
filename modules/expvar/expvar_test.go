package expvar

import (
	"testing"

	chi "github.com/go-chi/chi/v5"
	"github.com/titpetric/platform/pkg/require"
)

func TestHandler(t *testing.T) {
	h := NewHandler()

	require.NotNil(t, h)
	require.NoError(t, h.Mount(chi.NewRouter()))
}
