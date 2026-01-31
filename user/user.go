package user

import (
	"github.com/titpetric/platform-app/user/service"
)

func NewModule() *service.UserModule {
	return service.NewUserModule(service.Options{
		TemplateFS: templateFS,
	})
}
