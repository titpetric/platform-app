package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/pulse/client"
	"github.com/titpetric/platform-app/pulse/schema"
	"github.com/titpetric/platform-app/pulse/service/keycounter"
	"github.com/titpetric/platform-app/pulse/storage"
	userstorage "github.com/titpetric/platform-app/user/storage"
)

// Name is the service module name.
const Name = "pulse"

type PulseModule struct {
	platform.UnimplementedModule
	Options

	storage *storage.Storage

	userStorage *userstorage.UserStorage
}

func NewPulseModule(opts Options) *PulseModule {
	return &PulseModule{
		Options: opts,
	}
}

func (p *PulseModule) Name() string {
	return Name
}

func (p *PulseModule) setupStorage(ctx context.Context) error {
	if err := p.setupUserStorage(ctx); err != nil {
		return err
	}
	if err := p.setupPulseStorage(ctx); err != nil {
		return err
	}
	return nil
}

func (p *PulseModule) setupUserStorage(ctx context.Context) error {
	db, err := userstorage.DB(ctx)
	if err != nil {
		return err
	}

	p.userStorage = userstorage.NewUserStorage(db)
	return nil
}

func (p *PulseModule) setupPulseStorage(ctx context.Context) error {
	db, err := storage.DB(ctx)
	if err != nil {
		return err
	}

	if err := storage.Migrate(ctx, db, schema.Migrations); err != nil {
		return err
	}

	p.storage = storage.NewStorage(db)
	return nil
}

func (p *PulseModule) Start(ctx context.Context) error {
	if err := p.setupStorage(ctx); err != nil {
		return err
	}

	if p.Options.Record {
		log.Printf("[record] enabled recording keypress detail every %s", p.Options.RecordDuration)

		c := client.New(p.Options.Server)
		if err := c.EnsureToken(); err != nil {
			return fmt.Errorf("authentication required: %w (run 'pulse login' first)", err)
		}

		log.Printf("[record] sending to server: %s", p.Options.Server)

		lastRefresh := time.Now()

		opts := &keycounter.Options{
			FlushInterval: p.Options.RecordDuration,
			FlushFn: func(val int32) {
				if val <= 0 {
					return
				}

				if time.Since(lastRefresh) > 24*time.Hour {
					if err := c.RefreshToken(); err != nil {
						log.Printf("[record] token refresh failed: %v", err)
					} else {
						lastRefresh = time.Now()
						log.Println("[record] token refreshed")
					}
				}

				if err := c.SendPulse(int64(val)); err != nil {
					telemetry.CaptureError(ctx, err)
					log.Printf("[record] send failed: %v", err)
				}
			},
		}

		go func() {
			log.Println("[keycounter] started")
			keycounter.KeyboardCounter(ctx, opts)
			log.Println("[keycounter] exited")
		}()

	} else {
		log.Println("[record] recording keypress detail disabled")
	}

	return nil
}

func (p *PulseModule) Mount(ctx context.Context, r platform.Router) error {
	handlers := NewHandlers(p.storage)
	handlers.Mount(r)
	return nil
}
