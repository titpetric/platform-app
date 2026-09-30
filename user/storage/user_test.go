package storage

import (
	"testing"

	"github.com/titpetric/platform-app/internal/assert"
)

func TestNewUserStorage(t *testing.T) {
	s := NewUserStorage(nil)

	assert.NotNil(t, s)
}
