package storage

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitFS wraps os.DirFS with automatic git commits for write operations.
type GitFS struct {
	root string
	repo *git.Repository
}

// NewGitFS creates a GitFS instance for the given directory.
// If the directory does not exist, it will be created.
// If the directory is not a git repository, it will be initialized.
func NewGitFS(root string) (*GitFS, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git is not installed or not found in PATH: the blog module requires git for content management")
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(absRoot, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	gfs := &GitFS{root: absRoot}
	if err := gfs.initRepo(); err != nil {
		return nil, err
	}
	return gfs, nil
}

// initRepo opens or initializes the git repository.
func (g *GitFS) initRepo() error {
	repo, err := git.PlainOpen(g.root)
	if err == git.ErrRepositoryNotExists {
		return g.initNewRepo()
	}
	if err != nil {
		return fmt.Errorf("failed to open git repository: %w", err)
	}
	g.repo = repo
	return nil
}

// initNewRepo initializes a new git repository using git init command.
func (g *GitFS) initNewRepo() error {
	cmd := exec.Command("git", "init")
	cmd.Dir = g.root
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	repo, err := git.PlainOpen(g.root)
	if err != nil {
		return fmt.Errorf("failed to open newly initialized repository: %w", err)
	}
	g.repo = repo
	return nil
}

// Open implements fs.FS.
func (g *GitFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(g.root, name))
}

// ReadDir implements fs.ReadDirFS.
func (g *GitFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(g.root, name))
}

// ReadFile implements fs.ReadFileFS.
func (g *GitFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(g.root, name))
}

// Stat implements fs.StatFS.
func (g *GitFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(g.root, name))
}

// WriteFile writes content to a file and commits the change.
func (g *GitFS) WriteFile(name string, data []byte, perm fs.FileMode, auditMsg string) error {
	fullPath := filepath.Join(g.root, name)

	// Ensure parent directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the file
	if err := os.WriteFile(fullPath, data, perm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return g.commitFile(name, auditMsg)
}

// Remove removes a file and commits the change.
func (g *GitFS) Remove(name string, auditMsg string) error {
	fullPath := filepath.Join(g.root, name)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return g.commitFile(name, auditMsg)
}

// Rename renames a file and commits the change.
func (g *GitFS) Rename(oldName, newName string, auditMsg string) error {
	oldPath := filepath.Join(g.root, oldName)
	newPath := filepath.Join(g.root, newName)

	// Ensure parent directory exists for new path
	dir := filepath.Dir(newPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	// Stage both old (removal) and new (addition)
	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	if _, err := wt.Add(oldName); err != nil {
		return fmt.Errorf("failed to stage old file removal: %w", err)
	}
	if _, err := wt.Add(newName); err != nil {
		return fmt.Errorf("failed to stage new file: %w", err)
	}

	return g.commit(auditMsg)
}

// commitFile stages a single file and commits.
func (g *GitFS) commitFile(name string, auditMsg string) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	if _, err := wt.Add(name); err != nil {
		return fmt.Errorf("failed to stage file: %w", err)
	}

	return g.commit(auditMsg)
}

// commit creates a commit with the given message.
func (g *GitFS) commit(message string) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	_, err = wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Blog Admin",
			Email: "admin@blog.local",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// Root returns the root directory of the GitFS.
func (g *GitFS) Root() string {
	return g.root
}

// CopyFile copies a file from a reader and commits the change.
func (g *GitFS) CopyFile(name string, src io.Reader, perm fs.FileMode, auditMsg string) error {
	fullPath := filepath.Join(g.root, name)

	// Ensure parent directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, src); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return g.commitFile(name, auditMsg)
}
