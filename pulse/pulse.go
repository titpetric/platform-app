// Package pulse provides keystroke activity tracking and reporting.
package pulse

import (
	"github.com/titpetric/platform-app/pulse/service"
)

// NewModule creates a new pulse service module.
func NewModule() *service.PulseModule {
	return service.NewPulseModule()
}
