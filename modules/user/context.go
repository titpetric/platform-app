package user

import (
	"github.com/titpetric/platform/pkg/httpcontext"

	"github.com/titpetric/platform-app/modules/user/model"
)

type (
	sessionIDKey struct{}
	sessionKey   struct{}
	userKey      struct{}
)

var (
	userContext      = httpcontext.NewValue[*model.User](userKey{})
	sessionContext   = httpcontext.NewValue[*model.UserSession](sessionKey{})
	sessionIDContext = httpcontext.NewValue[string](sessionIDKey{})
)
