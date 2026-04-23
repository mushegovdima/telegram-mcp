package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	mcpserver "github.com/mushegovdima/telegram-mcp/internal/mcp"
	"github.com/mushegovdima/telegram-mcp/internal/telegram"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "TELEGRAM_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	defaultChatID := os.Getenv("TELEGRAM_DEFAULT_CHAT_ID")

	tg := telegram.New(token)
	s := mcpserver.NewServer(tg, defaultChatID)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
