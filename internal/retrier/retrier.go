package retrier

import (
	"context"
	"time"
)

type RetrierConfig struct {
	InitialWaitTime      uint
	RetriesCount         uint
	WaitTimeIncreaseFunc *func(uint) uint
}

type Retrier struct {
	c          RetrierConfig
	currentTry uint
	stopped    bool
}

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

func (r *Retrier) Try(ctx context.Context) bool {
	defer func() {
		r.currentTry++
	}()

	if r.currentTry == 0 {
		return true
	}

	if r.currentTry > 0 && r.currentTry != (r.c.RetriesCount+1) && !r.stopped {
		additionalWaitTime := (*r.c.WaitTimeIncreaseFunc)(r.currentTry - 1)
		waitTime := time.Second * time.Duration(r.c.InitialWaitTime+additionalWaitTime)

		select {
		case <-time.After(waitTime):
		case <-ctx.Done():
			return false
		}

		return true
	}

	return false
}

func (r *Retrier) Stop() {
	r.stopped = true
}
