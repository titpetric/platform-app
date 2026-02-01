package keycounter

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

// Linux input event constants
const (
	evKey     = 0x01
	keyPress  = 1
	eventSize = 24 // sizeof(struct input_event)
)

// inputEvent mirrors struct input_event from <linux/input.h>
type inputEvent struct {
	Sec   int64
	Usec  int64
	Type  uint16
	Code  uint16
	Value int32
}

// Options holds configuration options for the keyboard counter.
type Options struct {
	FlushFn       func(int32)
	FlushInterval time.Duration
}

// NewOptions will create a new *Options.
func NewOptions(flushFn func(int32), flushInterval time.Duration) *Options {
	return &Options{
		FlushFn:       flushFn,
		FlushInterval: flushInterval,
	}
}

// Flush to options function.
func (o *Options) Flush(c int32) {
	if o.FlushFn != nil {
		o.FlushFn(c)
	}
}

// KeyboardCounter counts keypresses and flushes the total at a fixed interval.
// It blocks until ctx is cancelled.
func KeyboardCounter(ctx context.Context, opts *Options) error {
	if opts == nil {
		return fmt.Errorf("no options configured for counter")
	}
	eventFiles, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		return fmt.Errorf("failed to list input devices: %w", err)
	}
	if len(eventFiles) == 0 {
		return fmt.Errorf("no input devices found")
	}

	var count int32
	done := make(chan struct{})

	// Start readers for each event device
	for _, path := range eventFiles {
		f, err := os.Open(path)
		if err != nil {
			// Skip unreadable devices
			continue
		}

		go func(file *os.File) {
			defer file.Close()

			buf := make([]byte, eventSize)
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				_, err := file.Read(buf)
				if err != nil {
					return
				}

				var ev inputEvent
				binary.Read(
					bytesReader(buf),
					binary.LittleEndian,
					&ev,
				)

				if ev.Type == evKey && ev.Value == keyPress {
					atomic.AddInt32(&count, 1)
				}
			}
		}(f)
	}

	ticker := time.NewTicker(opts.FlushInterval)
	defer ticker.Stop()

	flush := func() {
		n := atomic.SwapInt32(&count, 0)
		if n > 0 {
			opts.Flush(n)
		}
	}

	defer flush()

	for {
		select {
		case <-ctx.Done():
			close(done)
			return ctx.Err()

		case <-ticker.C:
			flush()
		}
	}
}

// bytesReader avoids allocations when decoding input events
func bytesReader(b []byte) *reader {
	return &reader{b: b}
}

type reader struct {
	b []byte
	i int
}

func (r *reader) Read(p []byte) (int, error) {
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
