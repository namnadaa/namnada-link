package telegram

import (
	"URLbot/pkg/storage"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

// Supported Telegram bot commands.
const (
	StartCmd = "/start"  // Shows a welcome message.
	RndCmd   = "/random" // Sends a random unread page.
	ReadCmd  = "/read"   // Marks a page as read.
	RmvCmd   = "/remove" // Removes a saved page.
	ListCmd  = "/list"   // Show all saved pages.
	HelpCmd  = "/help"   // Displays help information.
)

// doCmd handles an incoming command or message text from the user.
// If the text is a valid URL, it saves the page. Otherwise, it executes
// one of the supported bot commands such as /start, /rnd, /read, etc.
func (p *Processor) doCmd(text, username string, chatID int) error {
	text = strings.TrimSpace(text)

	slog.Info("got new command", "text", text, "username", username)

	cmd, arg := parseCmd(text)

	if isAddCmd(cmd) {
		return p.savePage(cmd, username, chatID)
	}

	switch cmd {
	case StartCmd:
		return p.sendHello(chatID)
	case RndCmd:
		return p.sendRandom(username, chatID)
	case ReadCmd:
		if arg == "" {
			return p.client.SendMessage(chatID, msgURLRequired)
		}
		return p.markAsRead(arg, username, chatID)
	case RmvCmd:
		if arg == "" {
			return p.client.SendMessage(chatID, msgURLRequired)
		}
		return p.removePage(arg, username, chatID)
	case ListCmd:
		return p.sendList(username, chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	default:
		return p.client.SendMessage(chatID, msgUnknownCommand)
	}
}

// savePage saves a new page for the given user if it does not already exist.
// After successful saving, it sends a confirmation message back to the user.
func (p *Processor) savePage(pageURL, username string, chatID int) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return fmt.Errorf("failed to check if the page exists: %v", err)
	}

	if isExists {
		return p.client.SendMessage(chatID, msgAlreadyExists)
	}

	err = p.storage.Save(page)
	if err != nil {
		return fmt.Errorf("failed to save page: %v", err)
	}

	err = p.client.SendMessage(chatID, msgSaved)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// sendHello sends a greeting message to the user.
func (p *Processor) sendHello(chatID int) error {
	return p.client.SendMessage(chatID, msgHello)
}

// sendRandom retrieves a random unread page for the user
// and sends its URL as a message. If there are no unread pages, it notifies the user.
func (p *Processor) sendRandom(username string, chatID int) error {
	page, err := p.storage.GetRandomUnread(username)
	if err != nil {
		if errors.Is(err, storage.ErrNoPagesFound) {
			return p.client.SendMessage(chatID, msgNoSavedPages)
		}

		return fmt.Errorf("failed to get random unread page: %v", err)
	}

	err = p.client.SendMessage(chatID, page.URL)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// markAsRead marks a specific page as read for the given user.
// It then sends a confirmation message back to the user.
func (p *Processor) markAsRead(pageURL, username string, chatID int) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	err := p.storage.MarkAsRead(page)
	if err != nil {
		return fmt.Errorf("failed to mark page as read: %v", err)
	}

	err = p.client.SendMessage(chatID, msgMarkedAsRead)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// removePage deletes a saved page for the given user
// and sends a confirmation message to the user.
func (p *Processor) removePage(pageURL, username string, chatID int) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	err := p.storage.Remove(page)
	if err != nil {
		return fmt.Errorf("failed to remove page: %v", err)
	}

	err = p.client.SendMessage(chatID, msgRemoved)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// sendList retrieves and sends the full list of saved pages for the user.
// Each page is shown with a [ ] or [x] prefix indicating unread or read status.
func (p *Processor) sendList(username string, chatID int) error {
	pages, err := p.storage.List(username)
	if err != nil {
		if errors.Is(err, storage.ErrNoPagesFound) {
			return p.client.SendMessage(chatID, msgNoSavedPages)
		}

		return fmt.Errorf("failed to fetch pages list: %v", err)
	}

	var builder strings.Builder
	builder.WriteString("Your saved pages:\n\n")

	for i, page := range pages {
		status := "[ ]"
		if page.Read {
			status = "[x]"
		}
		fmt.Fprintf(&builder, "%d. %s %s\n", i+1, status, page.URL)
	}

	err = p.client.SendMessage(chatID, builder.String())
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// sendHelp sends a help message describing all supported commands and usage instructions.
func (p *Processor) sendHelp(chatID int) error {
	return p.client.SendMessage(chatID, msgHelp)
}

// isAddCmd checks whether the given text should be treated as a "save page" command,
// i.e. whether it is a valid URL.
func isAddCmd(text string) bool {
	return isURL(text)
}

// isURL checks whether the given string is a valid URL that contains a host.
func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

// parseCmd splits the incoming message into a command and its argument.
func parseCmd(text string) (cmd string, arg string) {
	fields := strings.Fields(strings.TrimSpace(text))

	if len(fields) == 0 {
		return "", ""
	}

	cmd = fields[0]

	if len(fields) > 1 {
		arg = fields[1]
	}

	return
}
