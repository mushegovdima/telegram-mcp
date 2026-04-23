package mcp

import (
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mushegovdima/telegram-mcp/internal/telegram"
)

// NewServer constructs the MCP server and registers all Telegram tools.
// defaultChatID is used when chat_id is not provided in the tool call.
func NewServer(tg *telegram.Client, defaultChatID string) *server.MCPServer {
	s := server.NewMCPServer(
		"telegram-mcp",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	s.AddTool(mcplib.NewTool("send_message",
		mcplib.WithDescription("Send a text message to a Telegram chat. If chat_id is omitted, sends to the configured default user."),
		mcplib.WithString("chat_id",
			mcplib.Description("Telegram chat ID (numeric) or @username. Leave empty to send to the default user."),
		),
		mcplib.WithString("text",
			mcplib.Required(),
			mcplib.Description("Message text to send (Markdown supported)"),
		),
		mcplib.WithString("parse_mode",
			mcplib.Description("Optional parse mode: Markdown or HTML (default: plain text)"),
		),
	), sendMessageHandler(tg, defaultChatID))

	s.AddTool(mcplib.NewTool("get_updates",
		mcplib.WithDescription("Fetch recent messages sent to the bot (last N updates)"),
		mcplib.WithNumber("limit",
			mcplib.Description("Number of updates to fetch (1–100, default 10)"),
		),
		mcplib.WithNumber("offset",
			mcplib.Description("Identifier of the first update to return (for pagination)"),
		),
	), getUpdatesHandler(tg))

	s.AddTool(mcplib.NewTool("get_chat_info",
		mcplib.WithDescription("Get information about a chat (title, type, id) — useful to look up chat_id for a group or channel"),
		mcplib.WithString("chat_id",
			mcplib.Required(),
			mcplib.Description("Telegram chat ID (numeric) or @username"),
		),
	), getChatInfoHandler(tg))

	return s
}
