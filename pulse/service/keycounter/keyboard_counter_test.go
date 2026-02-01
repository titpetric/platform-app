package keycounter

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyboardCounter_CancelsCleanly(t *testing.T) {
	if !testing.Verbose() {
		return
	}

	ctx, cancel := context.WithCancel(t.Context())

	var flushCalls int32
	flushFn := func(n int32) {
		atomic.AddInt32(&flushCalls, n)
		t.Logf("keys: %d", n)
	}

	done := make(chan error, 1)

	t.Log("Listening to inputs for 8s, 2x 5s intervals (first = ticker, second = defer)")

	go func() {
		err := KeyboardCounter(ctx, NewOptions(flushFn, 5*time.Second))
		done <- err
	}()

	// Ensure the goroutine is running
	time.Sleep(8 * time.Second)

	cancel()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("KeyboardCounter did not exit after context cancellation")
	}
}
