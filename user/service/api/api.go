package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	tokenTTL       time.Duration
	userStorage    *storage.UserStorage
	sessionStorage *storage.SessionStorage
	revokedStorage *storage.RevokedTokenStorage
	passkeySvc     *passkey.Service

	emailActivationEnabled bool
	emailSender            EmailSender
	activationURLFormat    string
	activationSubject      string
}

// defaultTokenTTL mirrors service.DefaultTokenTTL but is duplicated here so
// the api package does not need to import its parent. Kept at 30 days to
// preserve historical behaviour.
const defaultTokenTTL = 30 * 24 * time.Hour

// defaultActivationSubject mirrors service.DefaultActivationSubject.
const defaultActivationSubject = "Confirm your account"

// NewHandlers returns a new Handlers instance with the given options.
// When opts.TokenTTL is zero the package default (30 days) is used.
func NewHandlers(opts Options) *Handlers {
	ttl := opts.TokenTTL
	if ttl <= 0 {
		ttl = defaultTokenTTL
	}
	subject := opts.ActivationSubject
	if subject == "" {
		subject = defaultActivationSubject
	}
	return &Handlers{
		signingKey:             opts.SigningKey,
		tokenTTL:               ttl,
		userStorage:            opts.UserStorage,
		sessionStorage:         opts.SessionStorage,
		revokedStorage:         opts.RevokedStorage,
		passkeySvc:             opts.PasskeyService,
		emailActivationEnabled: opts.EmailActivationEnabled,
		emailSender:            opts.EmailSender,
		activationURLFormat:    opts.ActivationURLFormat,
		activationSubject:      subject,
	}
}

// Mount registers the user API routes on the given router.
func (s *Handlers) Mount(r platform.Router) {
	r.Group(func(r platform.Router) {
		r.Post("/api/user/register", s.Register)
		r.Post("/api/user/token/create", s.CreateToken)
		r.Post("/api/user/token/refresh", s.RefreshToken)
		r.Post("/api/user/token/revoke", s.RevokeToken)

		r.Post("/api/user/email/activate", s.ActivateEmail)
		r.Post("/api/user/email/resend", s.ResendActivation)

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

	// Activation gate: only enforced when the policy is on, so existing
	// callers without the toggle see no behaviour change.
	if s.emailActivationEnabled {
		activated, aerr := s.userStorage.IsActivated(r.Context(), user.ID)
		if aerr != nil {
			return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to check activation")}
		}
		if !activated {
			return &RequestError{StatusCode: http.StatusForbidden, Err: model.ErrUserNotActivated}
		}
	}

	token, err := auth.NewJWT(s.signingKey).Create(user.ID, s.tokenTTL)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create token")}
	}

	resp := struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
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
	claims, err := jwtAuth.Claims(authHeader)
	if err != nil {
		return &RequestError{StatusCode: http.StatusUnauthorized, Err: errors.New("invalid token")}
	}

	// Reject already-revoked tokens before issuing a new one.
	if s.revokedStorage != nil && claims.JTI != "" {
		revoked, rerr := s.revokedStorage.IsRevoked(r.Context(), claims.JTI)
		if rerr != nil {
			return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to check revocation")}
		}
		if revoked {
			return &RequestError{StatusCode: http.StatusUnauthorized, Err: errors.New("token revoked")}
		}
	}

	token, err := jwtAuth.Create(claims.UserID, s.tokenTTL)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create token")}
	}

	// Best-effort revoke the presented token so it can't be reused
	// alongside the freshly-issued one. Failure here is non-fatal:
	// the caller still gets a valid new token.
	if s.revokedStorage != nil && claims.JTI != "" {
		_ = s.revokedStorage.Revoke(r.Context(), claims.JTI, claims.UserID, time.Unix(claims.ExpiresAt, 0))
	}

	resp := struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
	}

	platform.JSON(w, r, http.StatusOK, resp)
	return nil
}

func (s *Handlers) revokeToken(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return &RequestError{StatusCode: http.StatusUnauthorized, Err: errors.New("missing authorization header")}
	}

	claims, err := auth.NewJWT(s.signingKey).Claims(authHeader)
	if err != nil {
		return &RequestError{StatusCode: http.StatusUnauthorized, Err: errors.New("invalid token")}
	}

	if claims.JTI == "" {
		// Token predates JTI rollout; we cannot individually identify
		// it. Surface a 422 so callers know revocation is not possible
		// for this token (rather than silently lying with 204).
		return &RequestError{StatusCode: http.StatusUnprocessableEntity, Err: errors.New("token has no jti claim and cannot be revoked")}
	}

	if s.revokedStorage == nil {
		return &RequestError{StatusCode: http.StatusServiceUnavailable, Err: errors.New("revocation storage not configured")}
	}

	if err := s.revokedStorage.Revoke(r.Context(), claims.JTI, claims.UserID, time.Unix(claims.ExpiresAt, 0)); err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to revoke token")}
	}

	w.WriteHeader(http.StatusNoContent)
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

	// Branch on policy: with activation enabled, the user is created
	// pending and an activation email is sent; with activation
	// disabled, the user is created activated and immediately
	// receives a JWT just like before.
	if s.emailActivationEnabled {
		return s.registerPending(w, r, &req)
	}
	return s.registerActivated(w, r, &req)
}

// registerActivated is the pre-activation-toggle behaviour: create a
// user, mint a JWT, return both.
func (s *Handlers) registerActivated(w http.ResponseWriter, r *http.Request, req *model.UserCreateRequest) error {
	user, err := s.userStorage.Create(r.Context(), req)
	if err != nil {
		return mapRegisterError(err)
	}

	token, err := auth.NewJWT(s.signingKey).Create(user.ID, s.tokenTTL)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create token")}
	}

	platform.JSON(w, r, http.StatusCreated, struct {
		UserID    string `json:"user_id"`
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
	})
	return nil
}

// registerPending creates the user in pending state, fires the
// activation email, and returns a payload that explicitly signals
// "no token yet — go check your inbox". The 202 status reflects
// that the request has been accepted but is not yet complete.
func (s *Handlers) registerPending(w http.ResponseWriter, r *http.Request, req *model.UserCreateRequest) error {
	user, err := s.userStorage.CreatePending(r.Context(), req)
	if err != nil {
		return mapRegisterError(err)
	}

	// Best-effort send. If sending fails the user is left pending and
	// can re-trigger via /api/user/email/resend; surface the failure
	// so the client knows the registration succeeded but mail did not.
	if mailErr := s.userStorage.ResetActivation(r.Context(), req.Email); mailErr != nil {
		return &RequestError{StatusCode: http.StatusBadGateway, Err: fmt.Errorf("user created but activation email failed: %w", mailErr)}
	}

	platform.JSON(w, r, http.StatusAccepted, struct {
		UserID             string `json:"user_id"`
		RequiresActivation bool   `json:"requires_activation"`
	}{
		UserID:             user.ID,
		RequiresActivation: true,
	})
	return nil
}

// mapRegisterError converts a storage error to a RequestError with the
// most useful status code. Extracted so the two register branches
// don't duplicate the conflict-detection logic.
func mapRegisterError(err error) error {
	if errors.Is(err, model.ErrUsernameTaken) {
		return &RequestError{StatusCode: http.StatusConflict, Err: model.ErrUsernameTaken}
	}
	if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate") {
		return &RequestError{StatusCode: http.StatusConflict, Err: errors.New("email already exists")}
	}
	return &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
}

// ActivateEmail completes the activation flow. After a successful
// activation the user is logged in via the response JWT so the
// client can transition straight into the authenticated experience.
func (s *Handlers) ActivateEmail(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.activateEmail(w, r))
}

func (s *Handlers) activateEmail(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}
	if req.Token == "" {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("token is required")}
	}

	user, err := s.userStorage.Activate(r.Context(), req.Token)
	if err != nil {
		if errors.Is(err, model.ErrInvalidActivationToken) {
			return &RequestError{StatusCode: http.StatusNotFound, Err: model.ErrInvalidActivationToken}
		}
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to activate")}
	}

	token, err := auth.NewJWT(s.signingKey).Create(user.ID, s.tokenTTL)
	if err != nil {
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create token")}
	}

	platform.JSON(w, r, http.StatusOK, struct {
		UserID    string `json:"user_id"`
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
	})
	return nil
}

// ResendActivation re-issues an activation token and re-sends the
// email. Returns 204 on success regardless of whether the email
// actually exists, to avoid leaking account presence; logs and
// internal errors are still surfaced via 5xx.
func (s *Handlers) ResendActivation(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.resendActivation(w, r))
}

func (s *Handlers) resendActivation(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}
	if req.Email == "" {
		return &RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("email is required")}
	}

	err := s.userStorage.ResetActivation(r.Context(), req.Email)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows), errors.Is(err, model.ErrUserAlreadyActivated):
		// Swallow these to avoid disclosing account state. The
		// caller observes a 204 in both the "user doesn't exist"
		// and "user already activated" cases.
	default:
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to reset activation")}
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
