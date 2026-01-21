package storage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/schema"
)

// Storage provides database operations for the blog module
type Storage struct {
	db *sqlx.DB
}

// New creates a new Storage instance, runs migrations.
func New(ctx context.Context) (*Storage, error) {
	db, err := DB(ctx)
	if err != nil {
		return nil, err
	}
	return NewStorage(ctx, db)
}

// NewStorage creates a new Storage instance, rungs migrations.
func NewStorage(ctx context.Context, db *sqlx.DB) (*Storage, error) {
	if err := Migrate(ctx, db, schema.Migrations); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

// GetArticleBySlug retrieves an article by its slug
func (s *Storage) GetArticleBySlug(ctx context.Context, slug string) (*model.Article, error) {
	return GetArticleBySlug(ctx, s.db, slug)
}

// GetArticles retrieves all articles
func (s *Storage) GetArticles(ctx context.Context, start, length int) ([]model.Article, error) {
	return GetArticles(ctx, s.db, start, length)
}

// SearchArticles performs a full-text search on articles
func (s *Storage) SearchArticles(ctx context.Context, query string) ([]model.Article, error) {
	return SearchArticles(ctx, s.db, query)
}

// InsertArticle inserts a new article
func (s *Storage) InsertArticle(ctx context.Context, article *model.Article) error {
	return InsertArticle(ctx, s.db, article)
}

// CountArticles returns the total count of articles
func (s *Storage) CountArticles(ctx context.Context) (int, error) {
	return CountArticles(ctx, s.db)
}
