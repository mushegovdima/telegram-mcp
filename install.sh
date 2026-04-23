#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY_SERVER="$REPO_DIR/telegram-mcp"
LAUNCHER="$REPO_DIR/telegram-mcp-launcher.sh"

KEYCHAIN_TOKEN_SERVICE="telegram-mcp-token"
KEYCHAIN_CHAT_SERVICE="telegram-mcp-chat-id"

# ── Build ────────────────────────────────────────────────────────────────────
echo "Building MCP server..."
cd "$REPO_DIR"
go build -o "$BINARY_SERVER" ./cmd/server/
echo "Binary ready."

# ── jq check ────────────────────────────────────────────────────────────────
if ! command -v jq &>/dev/null; then
  echo "ERROR: jq is not installed. Run: brew install jq"
  exit 1
fi

# ── Ask for credentials ──────────────────────────────────────────────────────
echo ""
SAVED_TOKEN=$(security find-generic-password -s "$KEYCHAIN_TOKEN_SERVICE" -a "$USER" -w 2>/dev/null || true)
SAVED_CHAT_ID=$(security find-generic-password -s "$KEYCHAIN_CHAT_SERVICE" -a "$USER" -w 2>/dev/null || true)

if [[ -n "$SAVED_TOKEN" && -n "$SAVED_CHAT_ID" ]]; then
  echo "Using saved credentials from Keychain (token: ${SAVED_TOKEN:0:10}..., chat_id: $SAVED_CHAT_ID)"
  echo "Press Enter to keep them, or type new values."
  echo ""
fi

echo "Telegram Bot Token (from @BotFather) [${SAVED_TOKEN:0:10}...]:"
read -rs INPUT_TOKEN
echo ""
BOT_TOKEN="${INPUT_TOKEN:-$SAVED_TOKEN}"

echo "Your Telegram User ID (send any message to @userinfobot to find it) [$SAVED_CHAT_ID]:"
read -r INPUT_CHAT_ID
DEFAULT_CHAT_ID="${INPUT_CHAT_ID:-$SAVED_CHAT_ID}"

if [[ -z "$BOT_TOKEN" || -z "$DEFAULT_CHAT_ID" ]]; then
  echo "ERROR: Both bot token and user ID are required."
  exit 1
fi

# ── Save to Keychain ─────────────────────────────────────────────────────────
security add-generic-password -U -s "$KEYCHAIN_TOKEN_SERVICE" -a "$USER" -w "$BOT_TOKEN"      2>/dev/null
security add-generic-password -U -s "$KEYCHAIN_CHAT_SERVICE"  -a "$USER" -w "$DEFAULT_CHAT_ID" 2>/dev/null
echo "✓ Credentials stored in macOS Keychain."

# ── Create launcher script ───────────────────────────────────────────────────
cat > "$LAUNCHER" <<'LAUNCHER_EOF'
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
LAUNCHER_EOF
chmod 700 "$LAUNCHER"
echo "✓ Launcher script created."

# ── Patch claude_desktop_config.json ────────────────────────────────────────
CONFIG="$HOME/Library/Application Support/Claude/claude_desktop_config.json"
if [ ! -f "$CONFIG" ]; then
  mkdir -p "$(dirname "$CONFIG")"
  echo '{}' > "$CONFIG"
fi

TMP=$(mktemp)
jq --arg launcher "$LAUNCHER" \
  '.mcpServers.telegram = {"command": $launcher}' \
  "$CONFIG" > "$TMP"
mv "$TMP" "$CONFIG"
echo "✓ Claude Desktop MCP connector updated (no secrets in config)."

# ── Restart Claude Desktop ───────────────────────────────────────────────────
if pgrep -x "Claude" &>/dev/null; then
  echo "Restarting Claude Desktop..."
  osascript -e 'quit app "Claude"'
  sleep 2
fi
open -a Claude
echo "✓ Claude Desktop launched."

echo ""
echo "Done! Credentials are stored in macOS Keychain."

