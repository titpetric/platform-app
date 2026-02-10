package web

import "github.com/titpetric/platform-app/user/model"

type (
	Links struct {
		Login    string `json:"login"`
		Logout   string `json:"logout"`
		Register string `json:"register"`
		Recover  string `json:"recover"`
	}

	Data struct {
		SessionUser  *model.User `json:"sessionUser"`
		User         string      `json:"user"`
		Email        string      `json:"email"`
		Username     string      `json:"username"`
		ErrorMessage string      `json:"errorMessage"`
		FullName     string      `json:"fullName"`
		Links        Links       `json:"links"`
	}

	LoginData    = Data
	LogoutData   = Data
	RegisterData = Data
)
