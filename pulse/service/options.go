package service

import (
	"os"
	"time"

	"github.com/titpetric/cli"
)

// Options contains options for running pulse.
// Best to keep this a flat structure w/ simple types.
type Options struct {
	Record         bool
	RecordDuration time.Duration
	Server         string
}

func (c *Options) Bind(p *cli.FlagSet) {
	defaultServer := os.Getenv("PULSE_SERVER")
	if defaultServer == "" {
		defaultServer = "http://localhost:8080"
	}
	p.BoolVar(&c.Record, "record", false, "Record device input activity")
	p.DurationVarP(&c.RecordDuration, "duration", "d", 5*time.Minute, "Recording interval duration")
	p.StringVar(&c.Server, "server", defaultServer, "Pulse server URL for recording")
}
