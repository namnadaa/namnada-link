package memory_test

import (
	"URLbot/pkg/storage"
	"URLbot/pkg/storage/memory"
	"testing"
)

func TestStorage_Save(t *testing.T) {
	tests := []struct {
		name     string
		pages    []*storage.Page
		checkURL string
		user     string
		want     bool
	}{
		{
			name: "save single page",
			pages: []*storage.Page{
				{
					URL:      "https://example.com",
					UserName: "Alex",
				},
			},
			checkURL: "https://example.com",
			user:     "Alex",
			want:     true,
		},
		{
			name: "save duplicate page",
			pages: []*storage.Page{
				{
					URL:      "https://example.com",
					UserName: "Bob",
				},
				{
					URL:      "https://example.com",
					UserName: "Bob",
				},
			},
			checkURL: "https://example.com",
			user:     "Bob",
			want:     true,
		},
		{
			name: "different users",
			pages: []*storage.Page{
				{
					URL:      "https://example.com",
					UserName: "Alex",
				},
				{
					URL:      "https://example.com",
					UserName: "Bob",
				},
			},
			checkURL: "https://example.com",
			user:     "Bob",
			want:     true,
		},
		{
			name: "not saved page",
			pages: []*storage.Page{
				{
					URL:      "https://example.com",
					UserName: "Alex",
				},
			},
			checkURL: "https://falseexample.com",
			user:     "Alex",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := memory.New()

			for _, page := range tt.pages {
				err := s.Save(page)
				if err != nil {
					t.Fatalf("Save() failed: %v", err)
				}
			}

			page := storage.Page{
				URL:      tt.checkURL,
				UserName: tt.user,
			}

			exists, err := s.IsExists(&page)
			if err != nil {
				t.Fatalf("IsExists() failed: %v", err)
			}

			if exists != tt.want {
				t.Errorf("IsExists = %v; want %v", exists, tt.want)
			}
		})
	}
}

func TestStorage_GetRandomUnread(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		page     *storage.Page
		wantErr  bool
	}{
		{
			name:     "no page",
			userName: "Alex",
			page: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
				Read:     false,
			},
			wantErr: true,
		},
		{
			name:     "unread page",
			userName: "Alex",
			page: &storage.Page{
				URL:      "https://example.com",
				UserName: "Alex",
				Read:     false,
			},
			wantErr: false,
		},
		{
			name:     "read page",
			userName: "Alex",
			page: &storage.Page{
				URL:      "https://example.com",
				UserName: "Alex",
				Read:     true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := memory.New()

			err := s.Save(tt.page)
			if err != nil {
				t.Fatalf("failed to save page: %v", err)
			}

			got, gotErr := s.GetRandomUnread(tt.userName)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetRandomUnread() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("GetRandomUnread() succeeded unexpectedly")
			}

			if got.UserName != tt.userName {
				t.Errorf("got %s - want %s user name", got.UserName, tt.userName)
			}

			if got.Read {
				t.Error("expected unread page, but got read=true")
			}
		})
	}
}

func TestStorage_MarkAsRead(t *testing.T) {
	tests := []struct {
		name    string
		toSave  *storage.Page
		toMark  *storage.Page
		wantErr bool
	}{
		{
			name: "not read",
			toSave: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
				Read:     false,
			},
			toMark: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
			},
			wantErr: false,
		},
		{
			name: "already read",
			toSave: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
				Read:     true,
			},
			toMark: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
			},
			wantErr: false,
		},
		{
			name: "no page",
			toSave: &storage.Page{
				URL:      "https://example.com",
				UserName: "Alex",
				Read:     true,
			},
			toMark: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := memory.New()

			err := s.Save(tt.toSave)
			if err != nil {
				t.Fatalf("failed to save page: %v", err)
			}

			err = s.MarkAsRead(tt.toMark)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("MarkAsRead() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MarkAsRead() succeeded unexpectedly")
			}
		})
	}
}

func TestStorage_Remove(t *testing.T) {
	tests := []struct {
		name     string
		toSave   *storage.Page
		toRemove *storage.Page
		wantErr  bool
	}{
		{
			name: "exist",
			toSave: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
			},
			toRemove: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
			},
			wantErr: false,
		},
		{
			name: "does not exist",
			toSave: &storage.Page{
				URL:      "https://example.com",
				UserName: "Alex",
			},
			toRemove: &storage.Page{
				URL:      "https://example.com",
				UserName: "Bob",
			},
			wantErr: true,
		},
		{
			name: "nil page",
			toSave: &storage.Page{
				URL:      "https://example.com",
				UserName: "Alex",
			},
			toRemove: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := memory.New()

			err := s.Save(tt.toSave)
			if err != nil {
				t.Fatalf("failed to save page: %v", err)
			}

			err = s.Remove(tt.toRemove)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Remove() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Remove() succeeded unexpectedly")
			}
		})
	}
}
