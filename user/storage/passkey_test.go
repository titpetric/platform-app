package storage

import (
	"testing"

	"github.com/titpetric/platform/pkg/require"
)

func TestNewPasskeyStorage(t *testing.T) {
	s := NewPasskeyStorage(nil)

	require.NotNil(t, s)
}
