package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdates  = "getUpdates"
	sendMessage = "sendMessage"
)

// Client represents a Telegram Bot API client.
type Client struct {
	scheme   string
	host     string
	basePath string
	client   http.Client
}

// NewClient creates a new Telegram Bot API client with the given host and token.
func NewClient(scheme, host, token string) *Client {
	return &Client{
		scheme:   scheme,
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// newBasePath returns the base API path using the provided bot token.
func newBasePath(token string) string {
	return "bot" + token
}

// GetUpdates retrieves new updates (messages, commands, etc.) from Telegram.
func (c *Client) GetUpdates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdates, q)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	var res UpdateResponse

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if !res.Ok {
		return nil, fmt.Errorf("telegram API returned ok=false")
	}

	return res.Result, nil
}

// SendMessage sends a text message to the specified chat ID.
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessage, q)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	return nil
}

// doRequest performs an HTTP GET request to the Telegram API using the given method and query parameters.
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("doRequest [%s]: create request failed: %v", method, err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doRequest [%s]: request execution failed: %v", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doRequest [%s]: unexpected status %d", method, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doRequest [%s]: read response failed: %v", method, err)
	}

	return body, nil
}
