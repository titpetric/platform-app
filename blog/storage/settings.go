package storage

import (
	"context"
	"time"

	"github.com/titpetric/pdo/client"
	"github.com/titpetric/platform-app/blog/model"

	"database/sql"
)

// GetSettingByUserID retrieves settings for a specific user.
func GetSettingByUserID(ctx context.Context, db *sql.DB, userID string) (*model.Setting, error) {
	var setting model.Setting
	query := setting.Select(model.WithWhere("user_id = ?"), model.WithLimit(0, 1))

	err := client.NewClient(db).Get(ctx, &setting, query, userID)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}

// GetGlobalSettings retrieves the global settings (user_id = 'global').
func GetGlobalSettings(ctx context.Context, db *sql.DB) (*model.Setting, error) {
	return GetSettingByUserID(ctx, db, "global")
}

// SaveSetting inserts or updates settings for a user.
func SaveSetting(ctx context.Context, db *sql.DB, setting *model.Setting) error {
	now := time.Now()
	setting.SetUpdatedAt(now)
	setting.SetCreatedAt(now)

	query := setting.Insert(model.WithStatement("INSERT OR REPLACE INTO"))
	err := client.NewClient(db).Exec(ctx, query, setting)
	return err
}

// Storage methods for settings

// GetSettingByUserID retrieves settings for a specific user.
func (s *Storage) GetSettingByUserID(ctx context.Context, userID string) (*model.Setting, error) {
	return GetSettingByUserID(ctx, s.db, userID)
}

// GetGlobalSettings retrieves the global settings.
func (s *Storage) GetGlobalSettings(ctx context.Context) (*model.Setting, error) {
	return GetGlobalSettings(ctx, s.db)
}

// SaveSetting inserts or updates settings.
func (s *Storage) SaveSetting(ctx context.Context, setting *model.Setting) error {
	return SaveSetting(ctx, s.db, setting)
}
