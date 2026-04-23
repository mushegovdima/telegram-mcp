# telegram-mcp

MCP server that bridges Telegram Bot API to AI clients (Claude Desktop, etc.).

**Supported OS:** macOS

> Linux / Windows: the MCP server binary works on any platform, but `install.sh` is macOS-only. Configure Claude Desktop manually (see below).

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

### Sending a message

Ask Claude:
> Send "Hello!" to my Telegram

No need to specify a chat ID — `send_message` uses your user ID saved during install by default.

