package blog

import (
	"github.com/titpetric/platform-app/blog/service"
)

// NewModule creates a new blog service module.
func NewModule() *service.BlogModule {
	return service.NewBlogModule()
}
