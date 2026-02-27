package api

import (
	"github.com/titpetric/platform-app/user/service/passkey"
	"github.com/titpetric/platform-app/user/storage"
)

// Options is passed from user service scope.
type Options struct {
	SigningKey     string
	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage
	PasskeyService *passkey.Service
}
