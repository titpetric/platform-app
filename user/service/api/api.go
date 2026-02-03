package api

import (
	"net/http"

	"github.com/titpetric/platform"
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
		r.Post("/api/user/token/create", s.CreateToken)
		r.Post("/api/user/token/refresh", s.RefreshToken)
		r.Post("/api/user/token/revoke", s.RevokeToken)
	})
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
	return nil
}

func (s *Service) refreshToken(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Service) revokeToken(w http.ResponseWriter, r *http.Request) error {
	return nil
}
