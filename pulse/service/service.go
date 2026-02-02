package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"

	"github.com/titpetric/platform-app/user"
	"github.com/titpetric/platform-app/user/service/auth"
	userstorage "github.com/titpetric/platform-app/user/storage"

	"github.com/titpetric/platform-app/pulse/schema"
	"github.com/titpetric/platform-app/pulse/service/keycounter"
	"github.com/titpetric/platform-app/pulse/storage"
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

		token := os.Getenv("PULSE_AUTH")
		if token == "" {
			return fmt.Errorf("pulse: PULSE_AUTH undefined")
		}

		userID, err := auth.NewJWT(user.SigningKey()).UserID(token)
		if err != nil {
			return err
		}

		userdata, err := p.userStorage.Get(ctx, userID)
		if err != nil {
			return err
		}
		if err := userdata.Validate(); err != nil {
			return err
		}

		opts := &keycounter.Options{
			FlushInterval: p.Options.RecordDuration,
			FlushFn: func(val int32) {
				if val <= 0 {
					return
				}

				ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()

				ctx = user.SetSessionUser(ctx, userdata)

				if err := p.storage.Pulse(ctx, int64(val)); err != nil {
					telemetry.CaptureError(ctx, err)
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
