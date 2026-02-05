package client

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenFileOperations(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Cleanup(func() { os.Setenv("HOME", origHome) })
	os.Setenv("HOME", tmpDir)

	c := New("http://localhost:8080")

	err := c.LoadToken()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not logged in")

	token := "test-jwt-token"
	expiresAt := time.Now().Add(24 * time.Hour)
	err = c.SaveToken(token, expiresAt)
	require.NoError(t, err)

	tokenPath := filepath.Join(tmpDir, ".config", "pulse", "token.json")
	_, err = os.Stat(tokenPath)
	require.NoError(t, err)

	c2 := New("http://localhost:8080")
	err = c2.LoadToken()
	require.NoError(t, err)
	assert.Equal(t, token, c2.Token())
}

func TestShouldRefresh(t *testing.T) {
	tests := []struct {
		name     string
		token    *TokenData
		expected bool
	}{
		{
			name:     "nil token",
			token:    nil,
			expected: true,
		},
		{
			name: "expired token",
			token: &TokenData{
				Token:     "expired",
				ExpiresAt: time.Now().Add(-1 * time.Hour),
				SavedAt:   time.Now(),
			},
			expected: true,
		},
		{
			name: "old token needs refresh",
			token: &TokenData{
				Token:     "old",
				ExpiresAt: time.Now().Add(48 * time.Hour),
				SavedAt:   time.Now().Add(-25 * time.Hour),
			},
			expected: true,
		},
		{
			name: "fresh token",
			token: &TokenData{
				Token:     "fresh",
				ExpiresAt: time.Now().Add(48 * time.Hour),
				SavedAt:   time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New("http://localhost:8080")
			c.token = tt.token
			assert.Equal(t, tt.expected, c.ShouldRefresh())
		})
	}
}
