package view

import "github.com/titpetric/platform-app/user/model"

type (
	Links struct {
		Login    string `json:"login"`
		Logout   string `json:"logout"`
		Register string `json:"register"`
	}

	Data struct {
		SessionUser  *model.User `json:"sessionUser"`
		User         string      `json:"user"`
		Email        string      `json:"email"`
		ErrorMessage string      `json:"errorMessage"`
		FirstName    string      `json:"firstName"`
		LastName     string      `json:"lastName"`
		Links        Links       `json:"links"`
	}

	LoginData    = Data
	LogoutData   = Data
	RegisterData = Data
)
