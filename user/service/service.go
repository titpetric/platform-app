package service

import (
	"io/fs"

	"github.com/titpetric/platform-app/user/storage"
	"github.com/titpetric/platform-app/user/view"
)

// Service encapsulates what we need to get from the handler.
type Service struct {
	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage

	view *view.Renderer
}

// NewService takes in required dependencies to support the MVC framework.
func NewService(templateFS fs.FS, u *storage.UserStorage, s *storage.SessionStorage) *Service {
	svc := &Service{
		UserStorage:    u,
		SessionStorage: s,
		view:           view.NewRenderer(templateFS, nil),
	}
	return svc
}
