package storage

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform-app/blog/model"
)

// GetArticleBySlug retrieves a single article by slug. Includes drafts and scheduled articles.
func GetArticleBySlug(ctx context.Context, db *sqlx.DB, slug string) (*model.Article, error) {
	query := `SELECT * FROM article WHERE slug=? LIMIT 1`

	var article model.Article
	err := db.GetContext(ctx, &article, query, slug)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

// GetPublishedArticleBySlug retrieves a single published article by slug.
// Drafts and scheduled (future-dated) articles are not returned.
func GetPublishedArticleBySlug(ctx context.Context, db *sqlx.DB, slug string) (*model.Article, error) {
	query := `SELECT * FROM article WHERE slug=? AND draft=0 AND date<=? LIMIT 1`

	var article model.Article
	err := db.GetContext(ctx, &article, query, slug, time.Now())
	if err != nil {
		return nil, err
	}

	return &article, nil
}

// SearchPublishedArticles performs a case-insensitive substring search on published articles.
// The query is escaped so user-supplied `%` and `_` are treated literally.
func SearchPublishedArticles(ctx context.Context, db *sqlx.DB, find string) ([]model.Article, error) {
	searchTerm := "%" + escapeLike(find) + "%"

	query := `SELECT * FROM article
		WHERE draft = 0 AND date <= ?
		AND (title LIKE ? ESCAPE '\' OR description LIKE ? ESCAPE '\' OR slug LIKE ? ESCAPE '\')
		ORDER BY date DESC`

	var articles []model.Article
	err := db.SelectContext(ctx, &articles, query, time.Now(), searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// escapeLike escapes LIKE wildcard characters in a search string.
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// GetArticles retrieves all articles ordered by date descending.
func GetArticles(ctx context.Context, db *sqlx.DB, start, length int) ([]model.Article, error) {
	var article *model.Article
	query := article.Select(model.WithOrderBy("date DESC"), model.WithLimit(start, length))

	var articles []model.Article

	if err := db.SelectContext(ctx, &articles, query); err != nil {
		return nil, err
	}

	return articles, nil
}

// SearchArticles performs a search on articles by title, description, or slug.
// User-supplied wildcards (`%`, `_`) are escaped.
func SearchArticles(ctx context.Context, db *sqlx.DB, find string) ([]model.Article, error) {
	searchTerm := "%" + escapeLike(find) + "%"

	var article *model.Article
	query := article.Select(
		model.WithWhere(`title LIKE ? ESCAPE '\' OR description LIKE ? ESCAPE '\' OR slug LIKE ? ESCAPE '\'`),
		model.WithOrderBy("date DESC"),
	)

	var articles []model.Article
	err := db.SelectContext(ctx, &articles, query, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// InsertArticle inserts a new article into the database.
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

// CountArticles returns the total number of articles.
func CountArticles(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM article")
	return count, err
}

// GetPublishedArticles retrieves published articles (not draft, date <= now).
func GetPublishedArticles(ctx context.Context, db *sqlx.DB, start, length int) ([]model.Article, error) {
	var article *model.Article
	query := article.Select(
		model.WithWhere("draft = 0 AND date <= ?"),
		model.WithOrderBy("date DESC"),
		model.WithLimit(start, length),
	)

	var articles []model.Article
	if err := db.SelectContext(ctx, &articles, query, time.Now()); err != nil {
		return nil, err
	}
	return articles, nil
}

// CountPublishedArticles returns the count of published articles.
func CountPublishedArticles(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM article WHERE draft = 0 AND date <= ?", time.Now())
	return count, err
}

// GetDraftArticles retrieves draft articles.
func GetDraftArticles(ctx context.Context, db *sqlx.DB, start, length int) ([]model.Article, error) {
	var article *model.Article
	query := article.Select(
		model.WithWhere("draft = 1"),
		model.WithOrderBy("updated_at DESC"),
		model.WithLimit(start, length),
	)

	var articles []model.Article
	if err := db.SelectContext(ctx, &articles, query); err != nil {
		return nil, err
	}
	return articles, nil
}

// CountDraftArticles returns the count of draft articles.
func CountDraftArticles(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM article WHERE draft = 1")
	return count, err
}

// GetScheduledArticles retrieves scheduled articles (not draft, date > now).
func GetScheduledArticles(ctx context.Context, db *sqlx.DB, start, length int) ([]model.Article, error) {
	var article *model.Article
	query := article.Select(
		model.WithWhere("draft = 0 AND date > ?"),
		model.WithOrderBy("date ASC"),
		model.WithLimit(start, length),
	)

	var articles []model.Article
	if err := db.SelectContext(ctx, &articles, query, time.Now()); err != nil {
		return nil, err
	}
	return articles, nil
}

// CountScheduledArticles returns the count of scheduled articles.
func CountScheduledArticles(ctx context.Context, db *sqlx.DB) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM article WHERE draft = 0 AND date > ?", time.Now())
	return count, err
}

// UpdateArticle updates an existing article.
func UpdateArticle(ctx context.Context, db *sqlx.DB, article *model.Article) error {
	article.SetUpdatedAt(time.Now())

	query := article.Update(model.WithWhere("id = :id"))
	_, err := db.NamedExecContext(ctx, query, article)
	return err
}

// DeleteArticle deletes an article by slug.
func DeleteArticle(ctx context.Context, db *sqlx.DB, slug string) error {
	var article *model.Article
	query := article.Delete(model.WithWhere("slug = ?"))
	_, err := db.ExecContext(ctx, query, slug)
	return err
}

// GetArticleByID retrieves a single article by ID.
func GetArticleByID(ctx context.Context, db *sqlx.DB, id string) (*model.Article, error) {
	query := `SELECT * FROM article WHERE id=? LIMIT 1`

	var article model.Article
	err := db.GetContext(ctx, &article, query, id)
	if err != nil {
		return nil, err
	}
	return &article, nil
}
