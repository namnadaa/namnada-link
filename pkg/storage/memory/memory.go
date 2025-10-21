package memory

import (
	"URLbot/pkg/storage"
	"errors"
	"math/rand"
	"sync"
)

// ErrNoUnreadPages is returned when there are no unread pages for the user.
var ErrNoUnreadPages = errors.New("no unread pages")

// ErrNoPagesFound is returned when a page is not found for the user.
var ErrNoPagesFound = errors.New("page not found")

// Storage is an in-memory implementation of Storage interface.
type Storage struct {
	mu    sync.RWMutex
	pages map[string][]*storage.Page
}

// New creates a new in-memory storage.
func New() *Storage {
	return &Storage{
		pages: make(map[string][]*storage.Page),
	}
}

// Save stores a page for a given user.
func (s *Storage) Save(p *storage.Page) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, page := range s.pages[p.UserName] {
		if page.URL == p.URL {
			return nil
		}
	}

	s.pages[p.UserName] = append(s.pages[p.UserName], p)
	return nil
}

// GetRandomUnread returns a random unread page for a user.
func (s *Storage) GetRandomUnread(userName string) (*storage.Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pages := s.pages[userName]
	unread := make([]*storage.Page, 0, len(pages))
	for _, p := range pages {
		if !p.Read {
			unread = append(unread, p)
		}
	}

	if len(unread) == 0 {
		return nil, ErrNoUnreadPages
	}

	return unread[rand.Intn(len(unread))], nil
}

// MarkAsRead marks a page as read.
func (s *Storage) MarkAsRead(p *storage.Page) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, page := range s.pages[p.UserName] {
		if page.URL == p.URL {
			page.Read = true
			return nil
		}
	}
	return ErrNoPagesFound
}

// IsExists checks whether a page is already stored.
func (s *Storage) IsExists(p *storage.Page) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, page := range s.pages[p.UserName] {
		if page.URL == p.URL {
			return true, nil
		}
	}
	return false, nil
}

// Remove deletes a page.
func (s *Storage) Remove(p *storage.Page) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pages := s.pages[p.UserName]
	for i, page := range pages {
		if page.URL == p.URL {
			s.pages[p.UserName] = append(pages[:i], pages[i+1:]...)
			return nil
		}
	}
	return ErrNoPagesFound
}
