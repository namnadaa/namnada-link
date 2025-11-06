package consumer

import "context"

// Consumer defines the interface for starting the event processing loop.
type Consumer interface {
	Start(ctx context.Context) error
}
