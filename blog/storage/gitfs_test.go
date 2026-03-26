package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGitFS(t *testing.T) {
	tmpDir := t.TempDir()

	gfs, err := NewGitFS(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, gfs)
	assert.DirExists(t, filepath.Join(tmpDir, ".git"))
}

func TestGitFS_WriteFile(t *testing.T) {
	tmpDir := t.TempDir()

	gfs, err := NewGitFS(tmpDir)
	require.NoError(t, err)

	err = gfs.WriteFile("test.md", []byte("# Test"), 0o644, "Add test file")
	require.NoError(t, err)

	content, err := gfs.ReadFile("test.md")
	require.NoError(t, err)
	assert.Equal(t, "# Test", string(content))
}

func TestGitFS_Remove(t *testing.T) {
	tmpDir := t.TempDir()

	gfs, err := NewGitFS(tmpDir)
	require.NoError(t, err)

	err = gfs.WriteFile("test.md", []byte("# Test"), 0o644, "Add test file")
	require.NoError(t, err)

	err = gfs.Remove("test.md", "Remove test file")
	require.NoError(t, err)

	_, err = gfs.Stat("test.md")
	assert.True(t, os.IsNotExist(err))
}

func TestGitFS_Rename(t *testing.T) {
	tmpDir := t.TempDir()

	gfs, err := NewGitFS(tmpDir)
	require.NoError(t, err)

	err = gfs.WriteFile("old.md", []byte("# Old"), 0o644, "Add old file")
	require.NoError(t, err)

	err = gfs.Rename("old.md", "new.md", "Rename old to new")
	require.NoError(t, err)

	_, err = gfs.Stat("old.md")
	assert.True(t, os.IsNotExist(err))

	content, err := gfs.ReadFile("new.md")
	require.NoError(t, err)
	assert.Equal(t, "# Old", string(content))
}

func TestGitFS_Open(t *testing.T) {
	tmpDir := t.TempDir()

	gfs, err := NewGitFS(tmpDir)
	require.NoError(t, err)

	err = gfs.WriteFile("test.md", []byte("# Test"), 0o644, "Add test file")
	require.NoError(t, err)

	f, err := gfs.Open("test.md")
	require.NoError(t, err)
	defer f.Close()

	info, err := f.Stat()
	require.NoError(t, err)
	assert.Equal(t, "test.md", info.Name())
}

func TestGitFS_ReadDir(t *testing.T) {
	tmpDir := t.TempDir()

	gfs, err := NewGitFS(tmpDir)
	require.NoError(t, err)

	err = gfs.WriteFile("a.md", []byte("# A"), 0o644, "Add a.md")
	require.NoError(t, err)

	err = gfs.WriteFile("b.md", []byte("# B"), 0o644, "Add b.md")
	require.NoError(t, err)

	entries, err := gfs.ReadDir(".")
	require.NoError(t, err)

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.Name() != ".git" {
			names = append(names, e.Name())
		}
	}
	assert.Contains(t, names, "a.md")
	assert.Contains(t, names, "b.md")
}

func TestGitFS_Root(t *testing.T) {
	tmpDir := t.TempDir()

	gfs, err := NewGitFS(tmpDir)
	require.NoError(t, err)

	absRoot, _ := filepath.Abs(tmpDir)
	assert.Equal(t, absRoot, gfs.Root())
}
