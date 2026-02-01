package pulse

import (
	"github.com/titpetric/platform-app/pulse/service"
)

func NewModule(opts service.Options) *service.PulseModule {
	return service.NewPulseModule(opts)
}
