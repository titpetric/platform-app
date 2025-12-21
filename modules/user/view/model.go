package view

import "github.com/titpetric/platform-app/modules/user/model"

type (
	Data struct {
		SessionUser  *model.User `json:"sessionUser"`
		User         string      `json:"user"`
		Email        string      `json:"email"`
		ErrorMessage string      `json:"errorMessage"`
		FirstName    string      `json:"firstName"`
		LastName     string      `json:"lastName"`
	}

	LoginData    = Data
	LogoutData   = Data
	RegisterData = Data
)
