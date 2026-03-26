package storage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/schema"
)

// Storage provides database operations for the blog module.
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

// NewStorage creates a new Storage instance, runs migrations.
func NewStorage(ctx context.Context, db *sqlx.DB) (*Storage, error) {
	if err := Migrate(ctx, db, schema.Migrations); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

// GetArticleBySlug retrieves an article by its slug.
func (s *Storage) GetArticleBySlug(ctx context.Context, slug string) (*model.Article, error) {
	return GetArticleBySlug(ctx, s.db, slug)
}

// GetArticles retrieves all articles.
func (s *Storage) GetArticles(ctx context.Context, start, length int) ([]model.Article, error) {
	return GetArticles(ctx, s.db, start, length)
}

// SearchArticles performs a full-text search on articles.
func (s *Storage) SearchArticles(ctx context.Context, query string) ([]model.Article, error) {
	return SearchArticles(ctx, s.db, query)
}

// InsertArticle inserts a new article.
func (s *Storage) InsertArticle(ctx context.Context, article *model.Article) error {
	return InsertArticle(ctx, s.db, article)
}

// CountArticles returns the total count of articles.
func (s *Storage) CountArticles(ctx context.Context) (int, error) {
	return CountArticles(ctx, s.db)
}

// GetPublishedArticles retrieves published articles (not draft, date <= now).
func (s *Storage) GetPublishedArticles(ctx context.Context, start, length int) ([]model.Article, error) {
	return GetPublishedArticles(ctx, s.db, start, length)
}

// CountPublishedArticles returns the count of published articles.
func (s *Storage) CountPublishedArticles(ctx context.Context) (int, error) {
	return CountPublishedArticles(ctx, s.db)
}

// GetDraftArticles retrieves draft articles.
func (s *Storage) GetDraftArticles(ctx context.Context, start, length int) ([]model.Article, error) {
	return GetDraftArticles(ctx, s.db, start, length)
}

// CountDraftArticles returns the count of draft articles.
func (s *Storage) CountDraftArticles(ctx context.Context) (int, error) {
	return CountDraftArticles(ctx, s.db)
}

// GetScheduledArticles retrieves scheduled articles (not draft, date > now).
func (s *Storage) GetScheduledArticles(ctx context.Context, start, length int) ([]model.Article, error) {
	return GetScheduledArticles(ctx, s.db, start, length)
}

// CountScheduledArticles returns the count of scheduled articles.
func (s *Storage) CountScheduledArticles(ctx context.Context) (int, error) {
	return CountScheduledArticles(ctx, s.db)
}

// UpdateArticle updates an existing article.
func (s *Storage) UpdateArticle(ctx context.Context, article *model.Article) error {
	return UpdateArticle(ctx, s.db, article)
}

// DeleteArticle deletes an article by slug.
func (s *Storage) DeleteArticle(ctx context.Context, slug string) error {
	return DeleteArticle(ctx, s.db, slug)
}

// GetArticleByID retrieves a single article by ID.
func (s *Storage) GetArticleByID(ctx context.Context, id string) (*model.Article, error) {
	return GetArticleByID(ctx, s.db, id)
}
