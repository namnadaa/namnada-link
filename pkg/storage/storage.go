package storage

type Storage interface {
	Save(p *Page) error
	GetRandomUnread(userName string) (*Page, error)
	MarkAsRead(p *Page) error
	IsExists(p *Page) (bool, error)
	Remove(p *Page) error
}

type Page struct {
	URL      string
	UserName string
	Read     bool
}
