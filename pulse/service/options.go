package service

import (
	"time"

	"github.com/titpetric/cli"
)

// Options contains options for running pulse.
// Best to keep this a flat structure w/ simple types.
type Options struct {
	Record         bool
	RecordDuration time.Duration
}

func (c *Options) Bind(p *cli.FlagSet) {
	p.BoolVar(&c.Record, "record", false, "Record device input activity")
	p.DurationVarP(&c.RecordDuration, "duration", "d", 5*time.Minute, "Recording interval duration")
}
