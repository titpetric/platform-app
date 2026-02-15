package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"
	"github.com/titpetric/platform/pkg/ulid"
	"golang.org/x/crypto/bcrypt"

	"github.com/titpetric/platform-app/user/model"
)

// UserStorage implements the model.Storage interface using MySQL via sqlx.
type UserStorage struct {
	db *sqlx.DB
}

// NewUserStorage returns a new UserStorage backed by the given sqlx.DB.
func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

// Create inserts a new user and their authentication credentials.
// Returns an error if authentication information is missing.
func (s *UserStorage) Create(ctx context.Context, req *model.UserCreateRequest) (*model.User, error) {
	ctx, span := telemetry.StartAuto(ctx, s.Create)
	defer span.End()

	if !req.Valid() {
		if req.Username == "" {
			return nil, errors.New("missing authentication info: username is required")
		}
		if req.Email == "" || req.Password == "" {
			return nil, errors.New("missing authentication info: email and password are required")
		}
		if req.FullName == "" {
			return nil, errors.New("missing authentication info: full name is required")
		}
		return nil, errors.New("missing authentication info")
	}

	// Check if username already exists
	_, err := s.GetByUsername(ctx, req.Username)
	if err == nil {
		// User with this username exists
		return nil, model.ErrUsernameTaken
	}
	if err != sql.ErrNoRows {
		// Unexpected database error
		return nil, fmt.Errorf("check username: %w", err)
	}

	userAuth := req.UserAuth()

	_, span2 := telemetry.Start(ctx, "bcrypt.GenerateFromPassword")
	hashed, err := bcrypt.GenerateFromPassword([]byte(userAuth.Password), bcrypt.DefaultCost)
	span2.End()
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	userID := ulid.String()

	err = platform.Transaction(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		now := time.Now()

		userData := *req.User()
		userData.ID = userID
		userData.SetCreatedAt(now)
		userData.SetUpdatedAt(now)

		if _, err = tx.NamedExecContext(ctx, userData.Insert(), userData); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		userAuth := *userAuth
		userAuth.UserID = userData.ID
		userAuth.Password = string(hashed)
		userAuth.SetCreatedAt(now)
		userAuth.SetUpdatedAt(now)

		if _, err = tx.NamedExecContext(ctx, userAuth.Insert(), userAuth); err != nil {
			return fmt.Errorf("create user_auth: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return s.Get(ctx, userID)
}

// Update modifies an existing user and updates the updated_at timestamp.
func (s *UserStorage) Update(ctx context.Context, u *model.User) (*model.User, error) {
	ctx, span := telemetry.StartAuto(ctx, s.Update)
	defer span.End()

	u.SetUpdatedAt(time.Now())

	query := `UPDATE user SET full_name=?, deleted_at=?, updated_at=? WHERE id=?`

	_, err := s.db.ExecContext(ctx, query,
		u.FullName, u.DeletedAt, u.UpdatedAt, u.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return u, nil
}

// Get retrieves a user by ULID.
func (s *UserStorage) Get(ctx context.Context, id string) (*model.User, error) {
	ctx, span := telemetry.StartAuto(ctx, s.Get)
	defer span.End()

	u := &model.User{}
	query := `SELECT * FROM user WHERE id=?`
	if err := s.db.GetContext(ctx, u, query, id); err != nil {
		return nil, fmt.Errorf("get user id=%s: %w", id, err)
	}
	return u, nil
}

// GetByUsername retrieves a user by their username.
func (s *UserStorage) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	ctx, span := telemetry.StartAuto(ctx, s.GetByUsername)
	defer span.End()

	u := &model.User{}
	query := `SELECT * FROM user WHERE username=?`
	if err := s.db.GetContext(ctx, u, query, username); err != nil {
		return nil, err
	}
	return u, nil
}

// GetByStub retrieves a user by their slug.
func (s *UserStorage) GetByStub(ctx context.Context, slug string) (*model.User, error) {
	ctx, span := telemetry.StartAuto(ctx, s.GetByStub)
	defer span.End()

	u := &model.User{}
	query := `SELECT * FROM user WHERE slug=?`
	if err := s.db.GetContext(ctx, u, query, slug); err != nil {
		return nil, fmt.Errorf("get user slug=%s: %w", slug, err)
	}
	return u, nil
}

// GetGroups returns all groups the user belongs to.
func (s *UserStorage) GetGroups(ctx context.Context, userID string) ([]model.UserGroup, error) {
	ctx, span := telemetry.StartAuto(ctx, s.GetGroups)
	defer span.End()

	query := `
		SELECT g.id, g.title, g.created_at, g.updated_at
		FROM user_group g
		JOIN user_group_member m ON m.group_id = g.id
		WHERE m.user_id = ?
	`
	var groups []model.UserGroup
	if err := s.db.SelectContext(ctx, &groups, query, userID); err != nil {
		return nil, fmt.Errorf("get user groups: %w", err)
	}
	return groups, nil
}

// Authenticate verifies a user's credentials using bcrypt and returns the user.
func (s *UserStorage) Authenticate(ctx context.Context, userAuth model.UserAuth) (*model.User, error) {
	ctx, span := telemetry.StartAuto(ctx, s.Authenticate)
	defer span.End()

	if !userAuth.Valid() {
		return nil, errors.New("missing authentication info: email and password are required")
	}

	query := `SELECT user_id, password FROM user_auth WHERE email=? LIMIT 1`

	dbAuth := &model.UserAuth{}
	if err := s.db.GetContext(ctx, dbAuth, query, userAuth.Email); err != nil {
		return nil, fmt.Errorf("authenticate lookup: %w", err)
	}

	// instrument a cpu-heavy operation with an inner span
	err := func() error {
		_, span := telemetry.Start(ctx, "bcrypt.CompareHashAndPassword")
		err := bcrypt.CompareHashAndPassword([]byte(dbAuth.Password), []byte(userAuth.Password))
		span.End()

		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = sql.ErrNoRows
		}
		if err != nil {
			return fmt.Errorf("bcrypt compare: %w", err)
		}
		return nil
	}()
	if err != nil {
		return nil, err
	}

	user, err := s.Get(ctx, dbAuth.UserID)
	if err != nil {
		return nil, fmt.Errorf("authenticate get user: %w", err)
	}

	return user, nil
}

// List returns all active (non-deleted) users.
func (s *UserStorage) List(ctx context.Context) ([]model.User, error) {
	ctx, span := telemetry.StartAuto(ctx, s.List)
	defer span.End()

	var users []model.User
	query := `SELECT * FROM user WHERE deleted_at IS NULL ORDER BY username`
	if err := s.db.SelectContext(ctx, &users, query); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

var _ model.UserStorage = (*UserStorage)(nil)
