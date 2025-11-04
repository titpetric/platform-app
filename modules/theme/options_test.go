package theme

import (
	"testing"

	"github.com/titpetric/platform/pkg/require"
)

func TestTheme(t *testing.T) {
	require.NotNil(t, NewOptions())
}
