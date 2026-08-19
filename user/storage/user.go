package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/oida"
	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/ulid"
	"golang.org/x/crypto/bcrypt"

	emailmodel "github.com/titpetric/platform-app/email/model"
	emailstorage "github.com/titpetric/platform-app/email/storage"
	"github.com/titpetric/platform-app/user/model"
)

// UserStorage implements the model.Storage interface using MySQL via sqlx.
type UserStorage struct {
	db *sqlx.DB

	activationSubject   string
	activationURLFormat string
}

// dummyHash is a precomputed bcrypt hash used to spend approximately the
// same CPU on a "user not found" path as on a real password check. This
// is a defense against timing-based user enumeration via the login
// endpoint. The value is generated once at package init.
var dummyHash []byte

func init() {
	dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password-for-timing"), bcrypt.DefaultCost)
}

// NewUserStorage returns a new UserStorage backed by the given sqlx.DB.
func NewUserStorage(handle *sqlx.DB) *UserStorage {
	return &UserStorage{
		activationSubject:   "Activate your user",
		activationURLFormat: "https://{SITE_DOMAIN}/activation/%s",
		db:                  handle,
	}
}

// NewUserStorageErr returns a user storage and any error creating the database.
func NewUserStorageErr(ctx context.Context) (*UserStorage, error) {
	handle, err := DB(ctx)
	if err != nil {
		return nil, err
	}

	return NewUserStorage(handle), nil
}

// Create inserts a new user and their authentication credentials.
// Returns an error if authentication information is missing.
func (s *UserStorage) Create(ctx context.Context, req *model.UserCreateRequest) (*model.User, error) {
	ctx, span := oida.StartAuto(ctx, s.Create)
	defer span.End()

	if err := s.validateCreate(req); err != nil {
		return nil, err
	}

	if _, err := s.GetByUsername(ctx, req.Username); err == nil {
		return nil, model.ErrUsernameTaken
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("check username: %w", err)
	}

	_, span2 := oida.Start(ctx, "bcrypt.GenerateFromPassword")
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	span2.End()
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	userID := ulid.String()

	if err := s.insertUserAndAuth(ctx, req, userID, string(hashed), activatedNow, ""); err != nil {
		return nil, err
	}

	return s.Get(ctx, userID)
}

// CreatePending behaves like Create but leaves the new user un-activated
// and assigns an activation token that the caller is expected to deliver
// to the user (typically by email). The user cannot pass the activation
// gate (see IsActivated, Activate) until the token is exchanged.
func (s *UserStorage) CreatePending(ctx context.Context, req *model.UserCreateRequest) (*model.User, error) {
	ctx, span := oida.StartAuto(ctx, s.CreatePending)
	defer span.End()

	if err := s.validateCreate(req); err != nil {
		return nil, err
	}

	if _, err := s.GetByUsername(ctx, req.Username); err == nil {
		return nil, model.ErrUsernameTaken
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("check username: %w", err)
	}

	_, span2 := oida.Start(ctx, "bcrypt.GenerateFromPassword")
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	span2.End()
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	userID := ulid.String()
	token := newActivationToken()

	if err := s.insertUserAndAuth(ctx, req, userID, string(hashed), activatedNever, token); err != nil {
		return nil, err
	}

	user, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// activationState lets insertUserAndAuth express "activate now" vs
// "leave pending" without overloading bool semantics.
type activationState int

const (
	activatedNow activationState = iota
	activatedNever
)

// validateCreate centralises the cheap request validation shared by
// Create and CreatePending.
func (s *UserStorage) validateCreate(req *model.UserCreateRequest) error {
	if req.Valid() {
		return nil
	}
	if req.Username == "" {
		return errors.New("missing authentication info: username is required")
	}
	if req.Email == "" || req.Password == "" {
		return errors.New("missing authentication info: email and password are required")
	}
	if req.FullName == "" {
		return errors.New("missing authentication info: full name is required")
	}
	return errors.New("missing authentication info")
}

// insertUserAndAuth performs the user + user_auth inserts in one tx and
// then sets the activation columns according to the requested state.
// Keeping these in a single transaction means a partial failure cannot
// leave a user without credentials or an orphan user_auth row.
func (s *UserStorage) insertUserAndAuth(ctx context.Context, req *model.UserCreateRequest, userID, hashedPassword string, state activationState, activationToken string) error {
	return platform.Transaction(ctx, s.db, func(ctx context.Context, tx *sqlx.Tx) error {
		now := time.Now()

		userData := *req.User()
		userData.ID = userID
		userData.SetCreatedAt(now)
		userData.SetUpdatedAt(now)

		if _, err := tx.NamedExecContext(ctx, userData.Insert(), userData); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		userAuth := *req.UserAuth()
		userAuth.UserID = userData.ID
		userAuth.Password = hashedPassword
		userAuth.SetCreatedAt(now)
		userAuth.SetUpdatedAt(now)

		if _, err := tx.NamedExecContext(ctx, userAuth.Insert(), userAuth); err != nil {
			return fmt.Errorf("create user_auth: %w", err)
		}

		// The activation columns are not in UserAuthFields (they were
		// added by user_activation.up.sql after the generator ran), so
		// we patch them in with raw SQL instead of relying on Insert().
		switch state {
		case activatedNow:
			if _, err := tx.ExecContext(ctx, `UPDATE user_auth SET activated_at = ? WHERE user_id = ?`, now, userID); err != nil {
				return fmt.Errorf("activate user_auth: %w", err)
			}
		case activatedNever:
			if _, err := tx.ExecContext(ctx, `UPDATE user_auth SET activation_token = ?, activation_sent_at = ? WHERE user_id = ?`, activationToken, now, userID); err != nil {
				return fmt.Errorf("set activation token: %w", err)
			}
		}
		return nil
	})
}

// Update modifies an existing user and updates the updated_at timestamp.
func (s *UserStorage) Update(ctx context.Context, u *model.User) (*model.User, error) {
	ctx, span := oida.StartAuto(ctx, s.Update)
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
	ctx, span := oida.StartAuto(ctx, s.Get)
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
	ctx, span := oida.StartAuto(ctx, s.GetByUsername)
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
	ctx, span := oida.StartAuto(ctx, s.GetByStub)
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
	ctx, span := oida.StartAuto(ctx, s.GetGroups)
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
	ctx, span := oida.StartAuto(ctx, s.Authenticate)
	defer span.End()

	if !userAuth.Valid() {
		return nil, errors.New("missing authentication info: email and password are required")
	}

	query := `SELECT user_id, password FROM user_auth WHERE email=? LIMIT 1`

	dbAuth := &model.UserAuth{}
	if err := s.db.GetContext(ctx, dbAuth, query, userAuth.Email); err != nil {
		// Spend bcrypt time on the not-found path so the response
		// duration does not reveal whether the email exists. The
		// outcome of this compare is discarded; we still return the
		// original error so callers (including tests) keep getting
		// sql.ErrNoRows when the user does not exist.
		_, span := oida.Start(ctx, "bcrypt.CompareHashAndPassword.dummy")
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(userAuth.Password))
		span.End()
		return nil, fmt.Errorf("authenticate lookup: %w", err)
	}

	// instrument a cpu-heavy operation with an inner span
	err := func() error {
		_, span := oida.Start(ctx, "bcrypt.CompareHashAndPassword")
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

// IsActivated reports whether the given user has cleared the email
// activation gate. Returns true when the user exists and has a non-null
// activated_at. The pre-feature backfill in user_activation.up.sql means
// every user created before activation was introduced reads as activated.
func (s *UserStorage) IsActivated(ctx context.Context, userID string) (bool, error) {
	ctx, span := oida.StartAuto(ctx, s.IsActivated)
	defer span.End()

	var activatedAt *time.Time
	err := s.db.GetContext(ctx, &activatedAt, `SELECT activated_at FROM user_auth WHERE user_id=?`, userID)
	if err != nil {
		return false, fmt.Errorf("is activated: %w", err)
	}
	return activatedAt != nil, nil
}

// Activate exchanges an activation token for an activated user. The
// token is single-use: once consumed it is cleared from the row so
// re-presenting it yields ErrInvalidActivationToken.
func (s *UserStorage) Activate(ctx context.Context, token string) (*model.User, error) {
	ctx, span := oida.StartAuto(ctx, s.Activate)
	defer span.End()

	if token == "" {
		return nil, model.ErrInvalidActivationToken
	}

	var row struct {
		UserID      string     `db:"user_id"`
		ActivatedAt *time.Time `db:"activated_at"`
	}
	err := s.db.GetContext(ctx, &row, `SELECT user_id, activated_at FROM user_auth WHERE activation_token=? AND activation_token<>'' LIMIT 1`, token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrInvalidActivationToken
	}
	if err != nil {
		return nil, fmt.Errorf("activate lookup: %w", err)
	}

	// If the user is somehow already activated (race or stale token),
	// still clear the token so it can't be reused, and return success
	// — the desired end state is "activated, token gone".
	if _, err := s.db.ExecContext(ctx, `UPDATE user_auth SET activated_at = COALESCE(activated_at, ?), activation_token = '' WHERE user_id = ?`, time.Now(), row.UserID); err != nil {
		return nil, fmt.Errorf("activate: %w", err)
	}

	return s.Get(ctx, row.UserID)
}

// ResetActivation issues a fresh activation token for the user
// with the given email. Returns the new token. Errors with sql.ErrNoRows
// if no such user exists, and ErrUserAlreadyActivated if the user has
// already activated (callers should not resend in that case).
func (s *UserStorage) ResetActivation(ctx context.Context, email string) error {
	ctx, span := oida.StartAuto(ctx, s.ResetActivation)
	defer span.End()

	var row struct {
		UserID      string     `db:"user_id"`
		ActivatedAt *time.Time `db:"activated_at"`
	}
	if err := s.db.GetContext(ctx, &row, `SELECT user_id, activated_at FROM user_auth WHERE email=? LIMIT 1`, email); err != nil {
		return err
	}
	if row.ActivatedAt != nil {
		return model.ErrUserAlreadyActivated
	}

	token := newActivationToken()
	if _, err := s.db.ExecContext(ctx, `UPDATE user_auth SET activation_token = ?, activation_sent_at = NULL WHERE user_id = ?`, token, row.UserID); err != nil {
		return err
	}

	if mailErr := s.sendActivationEmail(ctx, email, token); mailErr != nil {
		return mailErr
	}

	return nil
}

// sendActivationEmail composes the activation email body and dispatches
// it through the configured EmailSender. Returns an error if no
// sender is configured or if delivery fails.
func (s *UserStorage) sendActivationEmail(ctx context.Context, to, token string) error {
	emailSender, err := emailstorage.NewEmailStorageErr(ctx)
	if err != nil {
		return err
	}

	body := s.activationBody(token)

	return emailSender.Create(ctx, emailmodel.NewEmail(to, s.activationSubject, body))
}

// activationBody renders the email body. If the ActivationURLFormat
// option is set, the token is interpolated into it and offered as a
// clickable link; otherwise the token is included verbatim with a
// short instruction.
func (s *UserStorage) activationBody(token string) string {
	if s.activationURLFormat != "" {
		return fmt.Sprintf("Please confirm your account by following this link:\n\n%s\n", fmt.Sprintf(s.activationURLFormat, token))
	}
	return fmt.Sprintf("Please confirm your account using the following token:\n\n%s\n", token)
}

// newActivationToken produces an opaque url-safe activation token.
// 32 random bytes (256 bits) base64url-encoded is overkill for guess
// resistance but fits comfortably in URLs and headers.
func newActivationToken() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		// rand.Read fails only on a broken kernel CSPRNG; nothing the
		// caller can do, and a non-panic fallback would silently
		// generate guessable tokens. Better to fail loud.
		panic("user/storage: crypto/rand failed: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(b[:])
}

// List returns all active (non-deleted) users.
func (s *UserStorage) List(ctx context.Context) ([]model.User, error) {
	ctx, span := oida.StartAuto(ctx, s.List)
	defer span.End()

	var users []model.User
	query := `SELECT * FROM user WHERE deleted_at IS NULL ORDER BY username`
	if err := s.db.SelectContext(ctx, &users, query); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

var _ model.UserStorage = (*UserStorage)(nil)
