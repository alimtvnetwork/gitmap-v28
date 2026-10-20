# schedule email

Register recurring or one-shot scheduled email delivery tasks directly within the GitMap scheduling subsystem.

## Aliases

`schedule-mail`, `sched-email`

---

## Usage

```bash
gitmap schedule email --to <emails> --subject <text> [--body <text>] [--file <path>] [--every <interval>]
```

---

## Flags

| Flag | Description |
|---|---|
| `--to <emails>`, `-t <emails>` | Recipient email addresses (comma-separated, masked in CLI outputs) |
| `--subject <text>`, `-s <text>` | Scheduled email subject line |
| `--body <text>`, `-b <text>` | Message body content |
| `--file <path>`, `-f <path>` | Read email body from a file |
| `--every <interval>`, `-i <interval>` | Recurring interval (e.g. `1d`, `2h`, `30m`, default: `1d`) |
| `--account <name>`, `-a <name>` | Outbound sender profile alias |
| `--json`, `-j` | Format registration response as Universal JSON Envelope |

---

## Examples

```bash
# Schedule daily fleet summary email
gitmap schedule email --to "admin@example.com" --subject "Fleet Status" --body "Daily status digest" --every 1d

# Schedule hourly monitoring digest from file
gitmap schedule email --to "ops@example.com" --subject "Hourly Alerts" --file ./alerts.txt --every 1h
```
