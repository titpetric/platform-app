package service

import (
	"github.com/titpetric/platform-app/modules/user/storage"
)

// Service encapsulates what we need to get from the handler.
type Service struct {
	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage
}

// NewService takes in required dependencies to support the MVC framework.
func NewService(u *storage.UserStorage, s *storage.SessionStorage) *Service {
	svc := &Service{
		UserStorage:    u,
		SessionStorage: s,
	}
	return svc
}
