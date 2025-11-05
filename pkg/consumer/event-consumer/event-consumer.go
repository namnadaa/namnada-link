package eventconsumer

import (
	"URLbot/pkg/events"
	"log/slog"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	procceser events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, procceser events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		procceser: procceser,
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
			slog.Error("Start: consumer unexpected error", "err", err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		slog.Info("got new message", "text", event.Text)

		err := c.procceser.Process(event)
		if err != nil {
			slog.Error("can't handle event", "err", err)
			continue
		}
	}

	return nil
}
