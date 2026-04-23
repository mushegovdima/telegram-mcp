#!/usr/bin/env bash
# Reads credentials from macOS Keychain and launches the MCP server.
TOKEN=$(security find-generic-password -s "telegram-mcp-token" -a "$USER" -w 2>/dev/null)
CHAT_ID=$(security find-generic-password -s "telegram-mcp-chat-id" -a "$USER" -w 2>/dev/null)
if [[ -z "$TOKEN" || -z "$CHAT_ID" ]]; then
  echo "ERROR: Credentials not found in Keychain. Run install.sh first." >&2
  exit 1
fi
export TELEGRAM_BOT_TOKEN="$TOKEN"
export TELEGRAM_DEFAULT_CHAT_ID="$CHAT_ID"
exec "$(dirname "$0")/telegram-mcp"
