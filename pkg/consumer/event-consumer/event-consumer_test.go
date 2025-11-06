package eventconsumer_test

import (
	eventconsumer "URLbot/pkg/consumer/event-consumer"
	"URLbot/pkg/events"
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type mockFetcher struct {
	calls  int32
	events [][]events.Event
	idx    int
}

func (m *mockFetcher) Fetch(batchSize int) ([]events.Event, error) {
	atomic.AddInt32(&m.calls, 1)

	if m.idx >= len(m.events) {
		return []events.Event{}, nil
	}

	result := m.events[m.idx]
	m.idx++
	return result, nil
}

type mockProcessor struct {
	called []events.Event
	err    error
}

func (m *mockProcessor) Process(event events.Event) error {
	m.called = append(m.called, event)
	return m.err
}

func TestConsumer_Start(t *testing.T) {
	tests := []struct {
		name          string
		fetcherEvents [][]events.Event
		processorErr  error
		batchSize     int
		wantCount     int
		wantErr       bool
	}{
		{
			name: "valid case",
			fetcherEvents: [][]events.Event{
				{
					{
						Type: events.Message,
						Text: "event1",
					},
					{
						Type: events.Message,
						Text: "event2",
					},
				},
				{},
			},
			processorErr: nil,
			batchSize:    10,
			wantCount:    2,
			wantErr:      false,
		},
		{
			name: "all processors fail",
			fetcherEvents: [][]events.Event{
				{
					{
						Type: events.Message,
						Text: "fail-event1",
					},
					{
						Type: events.Message,
						Text: "fail-event2",
					},
					{
						Type: events.Message,
						Text: "fail-event3",
					},
					{
						Type: events.Message,
						Text: "fail-event4",
					},
					{
						Type: events.Message,
						Text: "fail-event5",
					},
					{
						Type: events.Message,
						Text: "fail-event6",
					},
				},
				{},
			},
			processorErr: errors.New("mock error"),
			batchSize:    10,
			wantCount:    6,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()

			fetcher := &mockFetcher{
				events: tt.fetcherEvents,
			}

			processor := &mockProcessor{
				err: tt.processorErr,
			}

			consumer := eventconsumer.New(fetcher, processor, tt.batchSize)

			errCh := make(chan error, 1)
			go func() {
				errCh <- consumer.Start(ctx)
			}()

			select {
			case err := <-errCh:
				if tt.wantErr && err == nil {
					t.Error("expected error but got none")
				}
				if !tt.wantErr && err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			case <-time.After(1 * time.Second):
				t.Fatal("Start did not return in time")
			}

			if got := len(processor.called); got != tt.wantCount {
				t.Errorf("processed event count = %d, want %d", got, tt.wantCount)
			}
		})
	}
}
