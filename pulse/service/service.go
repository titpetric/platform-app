package service

import (
	"context"
	"log"
	"time"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"

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
}

func NewPulseModule(opts Options) *PulseModule {
	return &PulseModule{
		Options: opts,
	}
}

func (p *PulseModule) Name() string {
	return Name
}

func (p *PulseModule) Start(ctx context.Context) error {
	db, err := storage.DB(ctx)
	if err != nil {
		return err
	}

	if err := storage.Migrate(ctx, db, schema.Migrations); err != nil {
		return err
	}

	p.storage = storage.NewStorage(db)

	if p.Options.Record {
		log.Printf("[record] enabled recording keypress detail every %s", p.Options.RecordDuration)

		opts := &keycounter.Options{
			FlushInterval: p.Options.RecordDuration,
			FlushFn: func(val int32) {
				if val <= 0 {
					return
				}

				ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()

				if err := p.storage.Pulse(ctx, int64(val)); err != nil {
					log.Println("error storing pulse:", err)
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
