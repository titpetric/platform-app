package storage

import (
	"testing"

	"github.com/titpetric/platform/pkg/require"
)

func TestNewUserStorage(t *testing.T) {
	s := NewUserStorage(nil)

	require.NotNil(t, s)
}
