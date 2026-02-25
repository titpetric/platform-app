package keycounter

import "time"

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
