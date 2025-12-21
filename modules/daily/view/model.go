package view

import (
	"github.com/titpetric/platform-app/modules/daily/model"
)

type (
	Data struct {
		Tasks       []model.Todo `json:"tasks"`
		SessionUser *model.User  `json:"sessionUser"`
	}
)
