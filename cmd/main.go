package cmd

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/titpetric/platform"
	platformcmd "github.com/titpetric/platform/cmd"
	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/blog"
	"github.com/titpetric/platform-app/daily"
	"github.com/titpetric/platform-app/email"
	"github.com/titpetric/platform-app/expvar"
	"github.com/titpetric/platform-app/user"
)

// Main is the entrypoint for the app.
//
// It's expected to have control of the app lifecycle. An application
// exit is not expected to be graceful in case of errors. Main starts
// the platform server with modules loaded beforehand. It is blocking
// until server shutdown from cancellation of the context, or a caught
// SIGTERM, an OS control signal hinting the app should exit.
//
// The variadic parameter allows to inject options from test.
func Main(ctx context.Context, options ...*platform.Options) {
	Register()

	platformcmd.Main(ctx, options...)
}

// Register common middleware.
func Register() {
	platform.Use(middleware.Logger)
	platform.Use(telemetry.Middleware("user"))
	platform.Register(email.NewModule())
	platform.Register(blog.NewModule("data"))
	platform.Register(user.NewModule())
	platform.Register(expvar.NewModule())
	platform.Register(daily.NewModule())
}
