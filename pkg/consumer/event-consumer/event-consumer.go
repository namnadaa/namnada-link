package eventconsumer

import (
	"URLbot/pkg/events"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// Consumer implements the event-consuming logic using a Fetcher and Processor.
type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

// New creates and returns a new Consumer with the given fetcher, processor, and batch size.
func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

// Start begins the event processing loop. It periodically fetches events, processes
// them concurrently, and terminates on context cancellation.
func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			gotEvents, err := c.fetcher.Fetch(c.batchSize)
			if err != nil {
				slog.Error("Start: consumer unexpected error", "err", err)
				continue
			}

			if len(gotEvents) == 0 {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(1 * time.Second):
				}

				continue
			}

			err = c.handleEvents(gotEvents)
			if err != nil {
				return fmt.Errorf("too many errors %v", err)
			}
		}
	}
}

// handleEvents processes a slice of events concurrently.
// If more than 5 events fail during processing, it returns an error.
func (c *Consumer) handleEvents(ev []events.Event) error {
	var failed int32
	var wg sync.WaitGroup

	for _, event := range ev {
		wg.Add(1)
		go func(e events.Event) {
			defer wg.Done()

			slog.Info("got new message", "text", e.Text)

			err := c.processor.Process(e)
			if err != nil {
				slog.Error("can't handle event", "err", err)
				atomic.AddInt32(&failed, 1)
			}
		}(event)
	}
	wg.Wait()

	if failed >= 5 {
		return fmt.Errorf("failed to process %d events", failed)
	}

	return nil
}
