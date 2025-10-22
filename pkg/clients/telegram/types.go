package telegram

// UpdateResponse represents the top-level response from the Telegram API for getUpdates.
// The Result field contains a list of updates.
type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update represents a single update from Telegram, typically a new incoming message.
type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message"`
}

// Message represents a Telegram message sent by a user, including the text, from and chat info.
type Message struct {
	Text string `json:"text"`
	From From   `json:"From"`
	Chat Chat   `json:"chat"`
}

// From represents the sender of a Telegram message.
// It contains the Telegram username of the user who sent the message.
type From struct {
	Username string `json:"username"`
}

// Chat represents information about the chat where the message was sent.
type Chat struct {
	ID int `json:"id"`
}
