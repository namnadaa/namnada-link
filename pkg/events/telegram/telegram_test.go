package telegram_test

import (
	"URLbot/pkg/clients/telegram"
	"URLbot/pkg/events"
	tg "URLbot/pkg/events/telegram"
	"errors"
	"reflect"
	"testing"
)

type mockTelegramClient struct {
	updates []telegram.Update
	err     error
}

func (m *mockTelegramClient) GetUpdates(offset, limit int) ([]telegram.Update, error) {
	return m.updates, m.err
}

func (m *mockTelegramClient) SendMessage(chatID int, text string) error {
	return nil
}

func TestProcessor_Fetch(t *testing.T) {
	tests := []struct {
		name    string
		client  *mockTelegramClient
		limit   int
		want    []events.Event
		wantErr bool
	}{
		{
			name: "multiply updates",
			client: &mockTelegramClient{
				updates: []telegram.Update{
					{
						ID: 1,
						Message: &telegram.Message{
							Text: "test 1",
							From: telegram.From{Username: "User 1"},
							Chat: telegram.Chat{ID: 10},
						},
					},
					{
						ID: 2,
						Message: &telegram.Message{
							Text: "test 2",
							From: telegram.From{Username: "User 2"},
							Chat: telegram.Chat{ID: 20},
						},
					},
				},
			},
			limit: 10,
			want: []events.Event{
				{
					Type: events.Message,
					Text: "test 1",
					Meta: tg.Meta{
						ChatID:   10,
						UserName: "User 1",
					},
				},
				{
					Type: events.Message,
					Text: "test 2",
					Meta: tg.Meta{
						ChatID:   20,
						UserName: "User 2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty updates",
			client: &mockTelegramClient{
				updates: []telegram.Update{},
			},
			limit:   10,
			want:    nil,
			wantErr: false,
		},
		{
			name: "client error",
			client: &mockTelegramClient{
				err: errors.New("failed"),
			},
			limit:   10,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tg.New(tt.client, nil)
			got, gotErr := p.Fetch(tt.limit)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Fetch() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Fetch() succeeded unexpectedly")
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessor_Process(t *testing.T) {
	tests := []struct {
		name    string
		client  *mockTelegramClient
		event   events.Event
		wantErr bool
	}{
		{
			name: "valid message type",
			client: &mockTelegramClient{
				updates: []telegram.Update{
					{
						ID: 1,
						Message: &telegram.Message{
							Text: "test 1",
							From: telegram.From{Username: "User 1"},
							Chat: telegram.Chat{ID: 10},
						},
					},
				},
			},
			event: events.Event{
				Type: events.Message,
				Text: "test 1",
				Meta: tg.Meta{
					ChatID:   10,
					UserName: "User 1",
				},
			},
			wantErr: false,
		},
		{
			name: "unknown event type",
			client: &mockTelegramClient{
				updates: []telegram.Update{},
			},
			event:   events.Event{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tg.New(tt.client, nil)
			gotErr := p.Process(tt.event)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Process() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Process() succeeded unexpectedly")
			}
		})
	}
}
