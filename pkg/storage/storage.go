package storage

// Storage is an interface for saving, retrieving, and managing user pages.
type Storage interface {
	Save(p *Page) error
	GetRandomUnread(userName string) (*Page, error)
	MarkAsRead(p *Page) error
	IsExists(p *Page) (bool, error)
	Remove(p *Page) error
}

// Page represents a user-saved link with its read status.
type Page struct {
	URL      string
	UserName string
	Read     bool
}
