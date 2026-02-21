// Package record implements the pulse keystroke recording command.
package record

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/titpetric/cli"

	"github.com/titpetric/platform-app/pulse/client"
	"github.com/titpetric/platform-app/pulse/service/keycounter"
)

// Name is the command title.
const Name = "Record keystroke activity"

// Options holds record command configuration.
type Options struct {
	Name     string
	Server   string
	Duration string
}

// Bind registers record flags with the flag set.
func (o *Options) Bind(flag *cli.FlagSet) {
	server := "http://localhost:8080"
	if e := os.Getenv("PULSE_SERVER"); e != "" {
		server = e
	}

	name, _ := os.Hostname()

	flag.StringVar(&o.Name, "name", name, "Client name (hostname default)")
	flag.StringVar(&o.Server, "server", server, "Pulse server URL")
	flag.StringVar(&o.Duration, "duration", "5m", "Duration between pulse sends")
}

// NewCommand creates a new record command.
func NewCommand() *cli.Command {
	var opts Options

	return &cli.Command{
		Name:  "record",
		Title: Name,
		Bind:  opts.Bind,
		Run: func(ctx context.Context, args []string) error {
			return Run(ctx, opts)
		},
	}
}

// Run starts recording keystrokes and sending pulses.
func Run(ctx context.Context, opts Options) error {
	// Parse duration
	duration, err := time.ParseDuration(opts.Duration)
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	// Initialize client and load token
	c := client.New(opts.Server)
	if err := c.EnsureToken(); err != nil {
		return fmt.Errorf("authentication required: %w (run 'pulse login' first)", err)
	}

	log.Printf("Recording keypresses, sending every %s to %s", opts.Duration, opts.Server)

	lastRefresh := time.Now()

	var queued int64

	// Setup keyboard counter with flush function
	keyCounterOpts := &keycounter.Options{
		FlushInterval: duration,
		FlushFn: func(val int32) {
			if val <= 0 {
				return
			}

			// Refresh token if needed
			if time.Since(lastRefresh) > 24*time.Hour {
				if err := c.RefreshToken(); err != nil {
					log.Printf("token refresh failed: %v", err)
				} else {
					lastRefresh = time.Now()
					log.Println("token refreshed")
				}
			}

			// get any queued keystrokes
			qVal := atomic.SwapInt64(&queued, 0) + int64(val)

			tries := 3
			for {
				if tries == 0 {
					break
				}
				tries--

				if err := c.SendPulse(qVal, opts.Name); err != nil {
					log.Printf("send failed: %v", err)
				} else {
					log.Printf("sent %d keystroke(s)", val)
					return
				}

				time.Sleep(time.Minute)
			}

			atomic.AddInt64(&queued, qVal)
		},
	}

	// Run keyboard counter
	log.Println("starting keystroke counter")
	if err := keycounter.KeyboardCounter(ctx, keyCounterOpts); err != nil {
		return fmt.Errorf("keystroke counter: %w", err)
	}
	log.Println("keystroke counter exited")

	return nil
}
