package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform-app/blog/model"
)

// GetArticleBySlug retrieves a single article by slug
func GetArticleBySlug(ctx context.Context, db *sqlx.DB, slug string) (*model.Article, error) {
	query := `SELECT * FROM article	WHERE slug=? LIMIT 1`

	var article model.Article
	err := db.GetContext(ctx, &article, query, slug)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

// GetArticles retrieves all articles ordered by date descending
func GetArticles(ctx context.Context, db *sqlx.DB, start, length int) ([]model.Article, error) {
	var article *model.Article
	query := article.Select(model.WithOrderBy("date DESC"), model.WithLimit(start, length))

	var articles []model.Article

	if err := db.SelectContext(ctx, &articles, query); err != nil {
		return nil, err
	}

	return articles, nil
}

// SearchArticles performs a search on articles by title, description, or content
func SearchArticles(ctx context.Context, db *sqlx.DB, find string) ([]model.Article, error) {
	searchTerm := "%" + find + "%"

	var article *model.Article
	query := article.Select(
		model.WithWhere("title LIKE ? or description LIKE ? or slug LIKE ?"),
		model.WithOrderBy("date DESC"),
	)

	var articles []model.Article
	err := db.SelectContext(ctx, &articles, query, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// InsertArticle inserts a new article into the database
func InsertArticle(ctx context.Context, db *sqlx.DB, article *model.Article) error {
	now := time.Now()

	article.SetCreatedAt(now)
	article.SetUpdatedAt(now)
	if article.Date == nil {
		article.SetDate(now)
	}

	query := article.Insert(model.WithStatement("INSERT OR REPLACE INTO"))

	_, err := db.NamedExecContext(ctx, query, article)

	return err
}

// CountArticles returns the total number of articles
func CountArticles(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM article")
	return count, err
}
