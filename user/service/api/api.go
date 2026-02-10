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
	"github.com/titpetric/platform-app/user/storage"
)

type Service struct {
	opts        Options
	userStorage *storage.UserStorage
}

func NewService(userStorage *storage.UserStorage, opts Options) *Service {
	return &Service{
		userStorage: userStorage,
		opts:        opts,
	}
}

func (s *Service) Mount(r platform.Router) {
	r.Group(func(r platform.Router) {
		r.Post("/api/user/register", s.Register)
		r.Post("/api/user/token/create", s.CreateToken)
		r.Post("/api/user/token/refresh", s.RefreshToken)
		r.Post("/api/user/token/revoke", s.RevokeToken)
	})
}

func (s *Service) Register(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.register(w, r))
}

func (s *Service) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
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

func (s *Service) CreateToken(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.createToken(w, r))
}

func (s *Service) RefreshToken(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.refreshToken(w, r))
}

func (s *Service) RevokeToken(w http.ResponseWriter, r *http.Request) {
	s.errorHandler(w, r, s.revokeToken(w, r))
}

func (s *Service) createToken(w http.ResponseWriter, r *http.Request) error {
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
	token, err := auth.NewJWT(s.opts.SigningKey).Create(user.ID, ttl)
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

func (s *Service) refreshToken(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return &RequestError{StatusCode: http.StatusUnauthorized, Err: errors.New("missing authorization header")}
	}

	jwtAuth := auth.NewJWT(s.opts.SigningKey)
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

func (s *Service) revokeToken(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Service) register(w http.ResponseWriter, r *http.Request) error {
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
		return &RequestError{StatusCode: http.StatusInternalServerError, Err: errors.New("failed to create user")}
	}

	ttl := 30 * 24 * time.Hour
	token, err := auth.NewJWT(s.opts.SigningKey).Create(user.ID, ttl)
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
