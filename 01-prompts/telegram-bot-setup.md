# Two-Way Telegram Bot & Speed Notification Setup (`gitmap telegram` / `gitmap email`)

Use this canonical prompt and operational guide to connect a two-way Telegram chatbot and Email speed notification channel to any local or remote GitMap / Antigravity machine.

---

## 1. Quick Setup (`gitmap telegram setup`)

1. Create or retrieve your Telegram Bot Token from `@BotFather` and your target `CHAT_ID`.
2. Register the bot credentials in GitMap and sync with Antigravity Manager:

```bash
gitmap telegram setup --token "<TELEGRAM_BOT_TOKEN>" --chat "<TELEGRAM_CHAT_ID>"
# Or via agy namespace:
gitmap agy telegram setup --token "<TELEGRAM_BOT_TOKEN>" --chat "<TELEGRAM_CHAT_ID>"
```

3. Verify readiness and masked credentials:

```bash
gitmap telegram status
gitmap telegram status --json
```

---

## 2. Outbound & Two-Way Chat Commands (`send` & `poll`)

### Send Instant Message from Machine to Chat
```bash
gitmap telegram send "Build & pipeline green on $(gitmap machine ls)"
```

### Poll & Execute Incoming Two-Way Commands from Telegram
```bash
gitmap telegram poll
# Simulate or test a specific chat command directly:
gitmap telegram poll --command "/lap"
```

### Supported Two-Way Telegram Chat Commands
When chatting with your machine's Telegram bot, send any of the following commands and GitMap will execute and reply back to the chat:

| Chat Command | Action Executed on Machine |
|--------------|----------------------------|
| `/status` | Reports machine status, active Antigravity IDE PID, and workspace health |
| `/rp` | Lists currently running projects (`gitmap agy rp ls`) with 24h sequence IDs |
| `/lap` | Lists last active projects (`gitmap agy lap 24`) with conversation tree & prompts |
| `/rwi <seq> <prompt>` | Injects and reruns `<prompt>` into the project/conversation matching `<seq>` (`#1`, `P1`) |
| `/asw` | Triggers fast-forward Antigravity account switch check (`15%` remaining credit threshold) |
| `/help` | Replies with the full list of supported two-way bot commands |

---

## 3. Email Speed Setup & Unified Settings (`gitmap email` & `gitmap settings`)

```bash
# Configure SMTP Email speed notifications
gitmap email setup --smtp "smtp.gmail.com:587" --from "bot@example.com" --to "dev@example.com" --password "<APP_PASSWORD>"
gitmap email status
gitmap email test

# Configure default lookback window for 'lap' (default 24h) and account-switch threshold (default 15%)
gitmap settings --lap-hours 24 --threshold 15 --alias dev-win-01
gitmap settings --json
```
