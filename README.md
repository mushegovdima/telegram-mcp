# telegram-mcp

**Send Telegram messages from Claude with one command.**

MCP server that connects Claude Desktop to Telegram Bot API — so Claude can message you, fetch updates, and look up chats without leaving the conversation.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white) ![License](https://img.shields.io/github/license/mushegovdima/telegram-mcp) ![macOS](https://img.shields.io/badge/macOS-installer-black?logo=apple)

**Supported OS:** macOS (installer) · Linux / Windows (manual setup)

## Quick start

1. Get a bot token from [@BotFather](https://t.me/BotFather).
2. Find your Telegram user ID — send any message to [@userinfobot](https://t.me/userinfobot).
3. Run:

```bash
./install.sh
```

That's it. The script builds the binary, stores credentials in macOS Keychain, and connects Claude Desktop automatically.

---

## Project structure

```
cmd/server/main.go          — entrypoint
internal/telegram/client.go — Telegram HTTP client
internal/mcp/server.go      — MCP server + tool registration
internal/mcp/tools.go       — tool handlers
```

## Tools exposed

| Tool | Description |
|---|---|
| `send_message` | Send a text message to a chat |
| `get_updates` | Fetch recent messages sent to the bot |
| `get_chat_info` | Look up chat metadata |

## Build

```bash
go build -o telegram-mcp ./cmd/server/
```

## Connect to Claude Desktop

1. Create a Telegram bot via [@BotFather](https://t.me/BotFather) and copy the token.
2. Find your `chat_id`:
   - Start a conversation with your bot
   - Call `get_updates` — the `message.chat.id` field is your chat ID.
3. Run the installer — it builds the binary, stores credentials in macOS Keychain, patches `claude_desktop_config.json`, and restarts Claude:

```bash
./install.sh
```

4. The Telegram tools will appear in the tool list.

> **Linux / Windows:** the binary works on any platform, but `install.sh` is macOS-only. Add the server to `claude_desktop_config.json` manually and set `TELEGRAM_BOT_TOKEN` / `TELEGRAM_DEFAULT_CHAT_ID` env vars.

### Sending a message

Ask Claude:
> Send "Hello!" to my Telegram

No need to specify a chat ID — `send_message` uses your user ID saved during install by default.

