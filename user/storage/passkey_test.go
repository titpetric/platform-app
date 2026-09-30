package storage

import (
	"testing"

	"github.com/titpetric/platform-app/internal/assert"
)

func TestNewPasskeyStorage(t *testing.T) {
	s := NewPasskeyStorage(nil)

	assert.NotNil(t, s)
}
