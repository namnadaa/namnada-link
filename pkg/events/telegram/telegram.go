package telegram

import (
	"URLbot/pkg/clients/telegram"
	"URLbot/pkg/events"
	"URLbot/pkg/storage"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrCantGetMeta      = errors.New("can not get meta")
)

// Processor implements Fetcher interface for receiving Telegram updates
// and converting them into internal Event representations.
type Processor struct {
	client  Client
	offset  int
	storage storage.Storage
}

// Meta contains metadata extracted from an event, such as chat ID and username.
type Meta struct {
	ChatID   int
	UserName string
}

// Client abstracts Telegram API operations used by the bot.
type Client interface {
	GetUpdates(offset, limit int) ([]telegram.Update, error)
	SendMessage(chatID int, text string) error
}

// New creates a new Processor with the given Telegram client and storage.
func New(client Client, storage storage.Storage) *Processor {
	return &Processor{
		client:  client,
		storage: storage,
	}
}

// Fetch retrieves a batch of updates from Telegram, converts them to Event format,
// and updates the offset for the next fetch.
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.client.GetUpdates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}

	if len(updates) == 0 {
		slog.Debug("Fetch: there are no new updates")
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

// Process handles a single event by delegating to the appropriate handler
// based on the event type. Currently supports only Message events.
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

// processMessage extracts metadata from the event and processes the message command.
func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("failed to procces message: %v", err)
	}

	err = p.doCmd(event.Text, meta.UserName, meta.ChatID)
	if err != nil {
		return fmt.Errorf("failed to procces message: %v", err)
	}

	return nil
}

// meta extracts Meta information from the event and validates its type.
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrCantGetMeta
	}

	return res, nil
}

// event converts a Telegram update to an internal Event type.
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.Username,
		}
	}

	return res
}

// fetchType determines the type of the event based on the update content.
func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		slog.Warn("fetchType: incoming message are nil")
		return events.Unknown
	}

	return events.Message
}

// fetchText extracts the message text from a Telegram update.
func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		slog.Warn("fetchText: incoming message are nil")
		return ""
	}

	return upd.Message.Text
}
