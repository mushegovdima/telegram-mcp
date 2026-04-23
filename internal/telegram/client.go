package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const baseURL = "https://api.telegram.org/bot"

// Client is a minimal Telegram Bot API HTTP client.
type Client struct {
	token string
	http  *http.Client
}

// New creates a Client with the given bot token.
func New(token string) *Client {
	return &Client{
		token: token,
		http:  &http.Client{},
	}
}

// SendMessage sends a text message to the given chat.
// parseMode is optional ("Markdown", "HTML", or empty for plain text).
func (c *Client) SendMessage(chatID, text, parseMode string) error {
	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}
	if parseMode != "" {
		payload["parse_mode"] = parseMode
	}
	return c.post("sendMessage", payload)
}

// GetUpdates fetches the last `limit` updates starting from `offset`.
// Returns a pretty-printed JSON string.
func (c *Client) GetUpdates(limit, offset int) (string, error) {
	return c.postPretty("getUpdates", map[string]any{
		"limit":   limit,
		"offset":  offset,
		"timeout": 0,
	})
}

// GetChatInfo returns information about a chat as a pretty-printed JSON string.
func (c *Client) GetChatInfo(chatID string) (string, error) {
	return c.postPretty("getChat", map[string]any{"chat_id": chatID})
}

// post calls a Telegram API method and checks for `"ok": false` in the response.
func (c *Client) post(method string, payload map[string]any) error {
	body, err := c.postRaw(method, payload)
	if err != nil {
		return err
	}

	var resp struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}
	if !resp.OK {
		return fmt.Errorf("telegram API error: %s", resp.Description)
	}
	return nil
}

// postPretty calls a Telegram API method and returns the raw result as indented JSON.
func (c *Client) postPretty(method string, payload map[string]any) (string, error) {
	body, err := c.postRaw(method, payload)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, body, "", "  "); err != nil {
		return string(body), nil
	}
	return buf.String(), nil
}

// postRaw marshals the payload, POSTs to the Telegram API, and returns the raw body.
func (c *Client) postRaw(method string, payload map[string]any) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := baseURL + c.token + "/" + method
	resp, err := c.http.Post(url, "application/json", bytes.NewReader(data)) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s: %s", strconv.Itoa(resp.StatusCode), body)
	}

	return body, nil
}
