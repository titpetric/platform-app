package pulse

import (
	"github.com/titpetric/platform-app/pulse/service"
)

func NewModule(path string) *service.PulseModule {
	return service.NewPulseModule(service.Options{
		Path: path,
	})
}
