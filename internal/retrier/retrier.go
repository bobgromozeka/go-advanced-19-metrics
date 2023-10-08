package retrier

import (
	"context"
	"time"
)

// RetrierConfig Retrier configuration struct.
type RetrierConfig struct {
	WaitTimeIncreaseFunc *func(uint) uint
	InitialWaitTime      time.Duration
	RetriesCount         uint
}

// Retrier Has functionality to repeat actions.
type Retrier struct {
	c          RetrierConfig
	currentTry uint
	stopped    bool
}

// NewRetrier Create new Retrier with specified functionality.
func NewRetrier(config RetrierConfig) Retrier {
	if config.WaitTimeIncreaseFunc == nil {
		f := func(currentRetry uint) uint {
			return currentRetry * 2
		}

		config.WaitTimeIncreaseFunc = &f
	}

	return Retrier{
		config,
		0,
		false,
	}
}

// Try Returns true if next try should be performed.
func (r *Retrier) Try(ctx context.Context) bool {
	defer func() {
		r.currentTry++
	}()

	if r.currentTry == 0 {
		return true
	}

	if r.currentTry > 0 && r.currentTry != (r.c.RetriesCount+1) && !r.stopped {
		additionalWaitTime := (*r.c.WaitTimeIncreaseFunc)(r.currentTry - 1)
		waitTime := time.Second*r.c.InitialWaitTime + time.Duration(additionalWaitTime)

		select {
		case <-time.After(waitTime):
		case <-ctx.Done():
			return false
		}

		return true
	}

	return false
}

// Stop stops Retrier. After this call next Retrier.Try will return false.
func (r *Retrier) Stop() {
	r.stopped = true
}
