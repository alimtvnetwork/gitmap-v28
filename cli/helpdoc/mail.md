# mail

Send emails, inspect multi-mailbox status, test SMTP TLS handshakes, manage reusable templates, configure event-driven folder watchers, and launch the embedded mailing composer Web UI.

## Aliases

`email`, `emails`

---

## Subcommands

| Subcommand | Description |
|---|---|
| `send` | Send an email directly via SMTP with standard MIME encoding and automatic address masking in outputs |
| `check` | Check status, unread counts, and recent message headers across configured mailboxes |
| `test` | Test SMTP connectivity, port reachability, and implicit TLS / STARTTLS handshakes |
| `template <ls \| view \| add \| rm>` | Manage persistent email templates in Split SQLite storage |
| `notify <watch \| ls \| trigger>` | Configure or run event-driven folder watchers with burst debouncing and email dispatch |
| `ui` | Boot the embedded zero-CDN Mailing Composer Web UI (probes fallback ports starting at 8090) |

---

## Flags

| Flag | Description |
|---|---|
| `--to <emails>`, `-t <emails>` | Comma-separated recipient email addresses |
| `--subject <text>`, `-s <text>` | Subject line for the message |
| `--body <text>`, `-b <text>` | Plaintext or HTML body content |
| `--file <path>`, `-f <path>` | Read email body text from a local file |
| `--html` | Flag message content as HTML (`text/html`) format |
| `--account <name>`, `-a <name>` | Mailbox account or sender profile alias to use |
| `--unread-only`, `-u` | Filter mailbox check results to show accounts with unread messages only |
| `--show-full-email` | Unmask and reveal full email addresses instead of default `u***r@d***n.com` |
| `--json`, `-j` | Format CLI output into the Universal JSON Response Envelope |
| `--port <N>`, `-p <N>` | Port number for SMTP test or UI server listener |
| `--dry-run` | Validate parameters and message construction without dispatching network packets |

---

## Examples

```bash
# Send an email with masked address output
gitmap mail send --to "dev@example.com,ops@example.com" --subject "Build Green" --body "All checks passed"

# Check all mailboxes with unread counts
gitmap mail check --all

# Reveal unmasked email addresses
gitmap mail check --all --show-full-email

# Test SMTP connection with STARTTLS
gitmap mail test --host smtp.example.com --port 587 --tls

# Manage templates
gitmap mail template add deploy-notice --subject "[DEPLOY] Finished" --body "Release completed successfully"
gitmap mail template ls

# Launch the embedded Composer Web UI
gitmap mail ui --port 8090
```
