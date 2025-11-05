package eventconsumer

import (
	"URLbot/pkg/events"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			slog.Error("Start: consumer unexpected error", "err", err)
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		err = c.handleEvents(gotEvents)
		if err != nil {
			return fmt.Errorf("too many errors %v", err)
		}
	}
}

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

	if failed > 5 {
		return fmt.Errorf("failed to procces %d events", failed)
	}

	return nil
}
