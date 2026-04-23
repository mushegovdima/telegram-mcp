package mcp

import (
	"context"
	"fmt"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mushegovdima/telegram-mcp/internal/telegram"
)

func sendMessageHandler(tg *telegram.Client, defaultChatID string) func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	return func(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		chatID := req.GetString("chat_id", defaultChatID)
		text := req.GetString("text", "")
		parseMode := req.GetString("parse_mode", "")

		if chatID == "" {
			return mcplib.NewToolResultError("chat_id is required (or configure a default user ID at server startup)"), nil
		}
		if text == "" {
			return mcplib.NewToolResultError("text is required"), nil
		}

		if err := tg.SendMessage(chatID, text, parseMode); err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		return mcplib.NewToolResultText(fmt.Sprintf("Message sent to %s", chatID)), nil
	}
}

func getUpdatesHandler(tg *telegram.Client) func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	return func(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		limit := req.GetInt("limit", 10)
		offset := req.GetInt("offset", 0)

		updates, err := tg.GetUpdates(limit, offset)
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		return mcplib.NewToolResultText(updates), nil
	}
}

func getChatInfoHandler(tg *telegram.Client) func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	return func(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		chatID := req.GetString("chat_id", "")
		if chatID == "" {
			return mcplib.NewToolResultError("chat_id is required"), nil
		}

		info, err := tg.GetChatInfo(chatID)
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		return mcplib.NewToolResultText(info), nil
	}
}
