package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/user/model"
	"github.com/titpetric/platform-app/user/service/auth"
	"github.com/titpetric/platform-app/user/service/passkey"
	"github.com/titpetric/platform-app/user/storage"
)

// Handlers provides HTTP handlers for user authentication endpoints.
type Handlers struct {
	signingKey     string
	userStorage    *storage.UserStorage
	sessionStorage *storage.SessionStorage
	passkeySvc     *passkey.Service
}

// NewHandlers returns a new Handlers instance with the given options.
func NewHandlers(opts Options) *Handlers {
	return &Handlers{
		signingKey:     opts.SigningKey,
		userStorage:    opts.UserStorage,
		sessionStorage: opts.SessionStorage,
		passkeySvc:     opts.PasskeyService,
	}
}

// Mount registers the user API routes on the given router.
func (s *Handlers) Mount(r platform.Router) {
	r.Group(func(r platform.Router) {
		r.Post("/api/user/register", s.Register)
		r.Post("/api/user/token/create", s.CreateToken)
		r.Post("/api/user/token/refresh", s.RefreshToken)
		r.Post("/api/user/token/revoke", s.RevokeToken)

		r.Post("/api/passkey/register/begin", s.PasskeyRegisterBegin)
		r.Post("/api/passkey/register/finish", s.PasskeyRegisterFinish)
		r.Post("/api/passkey/login/begin", s.PasskeyLoginBegin)
		r.Post("/api/passkey/login/finish", s.PasskeyLoginFinish)
	})
}

// Register handles user registration requests.
func (s *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.register(w, r))
}

func (s *Handlers) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	switch val := err.(type) {
	case *RequestError:
		platform.Error(w, r, val.StatusCode, val.Err)
	default:
		platform.Error(w, r, 503, err)
	}
}

// CreateToken handles token creation requests.
func (s *Handlers) CreateToken(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.createToken(w, r))
}

// RefreshToken handles token refresh requests.
func (s *Handlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.refreshToken(w, r))
}

// RevokeToken handles token revocation requests.
func (s *Handlers) RevokeToken(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.revokeToken(w, r))
}

func (s *Handlers) createToken(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}

	userAuth := model.UserAuth{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := s.userStorage.Authenticate(r.Context(), userAuth)
	if err != nil {
		return &RequestError{StatusCode: http.StatusUnauthorized, Err: errors.New("invalid credentials")}
	}

	ttl := 30 * 24 * time.Hour
	token, err := auth.NewJWT(s.signingKey).Create(user.ID, ttl)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create token")}
	}

	resp := struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{
		Token:     token,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}

	platform.JSON(w, r, http.StatusOK, resp)
	return nil
}

func (s *Handlers) refreshToken(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return &RequestError{StatusCode: http.StatusUnauthorized, Err: errors.New("missing authorization header")}
	}

	jwtAuth := auth.NewJWT(s.signingKey)
	userID, err := jwtAuth.UserID(authHeader)
	if err != nil {
		return &RequestError{StatusCode: http.StatusUnauthorized, Err: errors.New("invalid token")}
	}

	ttl := 30 * 24 * time.Hour
	token, err := jwtAuth.Create(userID, ttl)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create token")}
	}

	resp := struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{
		Token:     token,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}

	platform.JSON(w, r, http.StatusOK, resp)
	return nil
}

func (s *Handlers) revokeToken(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PasskeyRegisterBegin starts passkey registration.
func (s *Handlers) PasskeyRegisterBegin(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.passkeyRegisterBegin(w, r))
}

func (s *Handlers) passkeyRegisterBegin(w http.ResponseWriter, r *http.Request) error {
	var req model.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}

	if req.Username == "" {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: model.ErrUsernameMissing}
	}
	if err := req.ValidateUsername(); err != nil {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: err}
	}
	if req.FullName == "" || req.Email == "" {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("full name and email are required")}
	}

	token, options, err := s.passkeySvc.BeginRegistration(&req)
	if err != nil {
		return err
	}

	platform.JSON(w, r, http.StatusOK, struct {
		Token   string `json:"token"`
		Options any    `json:"options"`
	}{Token: token, Options: options})
	return nil
}

// PasskeyRegisterFinish completes passkey registration.
func (s *Handlers) PasskeyRegisterFinish(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.passkeyRegisterFinish(w, r))
}

func (s *Handlers) passkeyRegisterFinish(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("X-Passkey-Token")
	if token == "" {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("missing passkey token")}
	}

	result, err := s.passkeySvc.FinishRegistration(token, r)
	if err != nil {
		return err
	}

	session, err := s.sessionStorage.Create(r.Context(), result.UserID)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create session")}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  *session.ExpiresAt,
	})

	platform.JSON(w, r, http.StatusCreated, struct {
		UserID string `json:"user_id"`
	}{UserID: result.UserID})
	return nil
}

// PasskeyLoginBegin starts passkey login.
func (s *Handlers) PasskeyLoginBegin(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.passkeyLoginBegin(w, r))
}

func (s *Handlers) passkeyLoginBegin(w http.ResponseWriter, r *http.Request) error {
	token, options, err := s.passkeySvc.BeginLogin()
	if err != nil {
		return err
	}

	platform.JSON(w, r, http.StatusOK, struct {
		Token   string `json:"token"`
		Options any    `json:"options"`
	}{Token: token, Options: options})
	return nil
}

// PasskeyLoginFinish completes passkey login.
func (s *Handlers) PasskeyLoginFinish(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.passkeyLoginFinish(w, r))
}

func (s *Handlers) passkeyLoginFinish(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("X-Passkey-Token")
	if token == "" {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("missing passkey token")}
	}

	result, err := s.passkeySvc.FinishLogin(token, r)
	if err != nil {
		return err
	}

	session, err := s.sessionStorage.Create(r.Context(), result.UserID)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create session")}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  *session.ExpiresAt,
	})

	platform.JSON(w, r, http.StatusOK, struct {
		UserID string `json:"user_id"`
	}{UserID: result.UserID})
	return nil
}

func (s *Handlers) register(w http.ResponseWriter, r *http.Request) error {
	var req model.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}

	if !req.Valid() {
		// Check for specific username validation errors
		if req.Username == "" {
			return &RequestError{StatusCode: http.StatusBadRequest, Err: model.ErrUsernameMissing}
		}
		if err := req.ValidateUsername(); err != nil {
			return &RequestError{StatusCode: http.StatusBadRequest, Err: err}
		}
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("invalid request: all fields are required")}
	}

	user, err := s.userStorage.Create(r.Context(), &req)
	if err != nil {
		if errors.Is(err, model.ErrUsernameTaken) {
			return &RequestError{StatusCode: http.StatusConflict, Err: model.ErrUsernameTaken}
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate") {
			return &RequestError{StatusCode: http.StatusConflict, Err: errors.New("email already exists")}
		}
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	ttl := 30 * 24 * time.Hour
	token, err := auth.NewJWT(s.signingKey).Create(user.ID, ttl)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create token")}
	}

	resp := struct {
		UserID    string `json:"user_id"`
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}

	platform.JSON(w, r, http.StatusCreated, resp)
	return nil
}
