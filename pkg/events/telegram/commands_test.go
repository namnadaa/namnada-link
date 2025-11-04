package telegram

import (
	"URLbot/pkg/clients/telegram"
	"URLbot/pkg/storage"
	"strings"
	"testing"
)

type mockClient struct {
	sent []string
	err  error
}

func (m *mockClient) GetUpdates(offset, limit int) ([]telegram.Update, error) {
	return nil, m.err
}

func (m *mockClient) SendMessage(chatID int, text string) error {
	m.sent = append(m.sent, text)
	return nil
}

type mockStorage struct {
	pages []*storage.Page
	err   error
}

func (m *mockStorage) Save(p *storage.Page) error {
	return m.err
}

func (m *mockStorage) GetRandomUnread(userName string) (*storage.Page, error) {
	if len(m.pages) == 0 {
		return nil, storage.ErrNoPagesFound
	}
	return m.pages[0], nil
}

func (m *mockStorage) MarkAsRead(p *storage.Page) error {
	return nil
}

func (m *mockStorage) IsExists(p *storage.Page) (bool, error) {
	return false, nil
}

func (m *mockStorage) Remove(p *storage.Page) error {
	return nil
}

func (m *mockStorage) List(userName string) ([]*storage.Page, error) {
	return m.pages, nil
}

func TestProcessor_doCmd(t *testing.T) {
	tests := []struct {
		name     string
		client   *mockClient
		storage  *mockStorage
		text     string
		username string
		chatID   int
		wantSend string
	}{
		{
			name:     "start command",
			client:   &mockClient{},
			storage:  &mockStorage{},
			text:     "/start",
			username: "alex",
			wantSend: msgHello,
		},
		{
			name:   "random command with no pages",
			client: &mockClient{},
			storage: &mockStorage{
				pages: []*storage.Page{},
			},
			text:     "/random",
			username: "alex",
			wantSend: msgNoSavedPages,
		},
		{
			name:     "read command without arg",
			client:   &mockClient{},
			storage:  &mockStorage{},
			text:     "/read",
			username: "alex",
			wantSend: msgURLRequired,
		},
		{
			name:     "read command with arg",
			client:   &mockClient{},
			storage:  &mockStorage{},
			text:     "/read https://example.com",
			username: "alex",
			wantSend: msgMarkedAsRead,
		},
		{
			name:     "remove command without arg",
			client:   &mockClient{},
			storage:  &mockStorage{},
			text:     "/remove",
			username: "alex",
			wantSend: msgURLRequired,
		},
		{
			name:     "remove command with arg",
			client:   &mockClient{},
			storage:  &mockStorage{},
			text:     "/remove https://example.com",
			username: "alex",
			wantSend: msgRemoved,
		},
		{
			name:   "list command with saved pages",
			client: &mockClient{},
			storage: &mockStorage{
				pages: []*storage.Page{
					{
						URL:  "https://example.com",
						Read: false,
					},
					{
						URL:  "https://readpage.com",
						Read: true,
					},
				},
			},
			text:     "/list",
			username: "alex",
			wantSend: "Your saved pages:",
		},
		{
			name:     "help command",
			client:   &mockClient{},
			storage:  &mockStorage{},
			text:     "/help",
			username: "alex",
			wantSend: msgHelp,
		},
		{
			name:     "save link",
			client:   &mockClient{},
			storage:  &mockStorage{},
			text:     "http://example.com",
			username: "alex",
			wantSend: msgSaved,
		},
		{
			name:     "unknown command",
			client:   &mockClient{},
			storage:  &mockStorage{},
			text:     "/unknown",
			username: "alex",
			wantSend: msgUnknownCommand,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.client
			p := New(client, tt.storage)

			gotErr := p.doCmd(tt.text, tt.username, tt.chatID)
			if gotErr != nil {
				t.Fatalf("doCmd() failed: %v", gotErr)
			}

			if len(client.sent) == 0 {
				t.Fatal("the message was not sent, but it should have been sent")
			}

			got := client.sent[len(client.sent)-1]
			if !strings.Contains(got, tt.wantSend[:min(len(tt.wantSend), 10)]) {
				t.Fatalf("unexpected message: got %v, want %v", got, tt.wantSend)
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
