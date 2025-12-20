package service

import (
	"net/http"

	"github.com/titpetric/platform-app/modules/user"
	usermodel "github.com/titpetric/platform-app/modules/user/model"
)

type Permissions struct {
	User *usermodel.User

	Create bool
}

func NewPermissions(r *http.Request) Permissions {
	var isLoggedIn bool

	user, ok := user.GetSessionUser(r.Context())
	if ok {
		isLoggedIn = true
	}

	return Permissions{
		User:   user,
		Create: isLoggedIn,
	}
}
