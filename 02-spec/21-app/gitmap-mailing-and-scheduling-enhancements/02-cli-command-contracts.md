# GitMap Mailing & Scheduling CLI Command Contracts

**Specification Path:** `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/02-cli-command-contracts.md`  
**Parent Index:** `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/readme.md`  
**Status:** DRAFT / RATIFIED  
**Target Subsystems:** `cli/cmdschedule`, `cli/cmdmail`, `cli/cmd`

---

## 1. Command Tree Structure & Routing

GitMap exposes scheduling and mailing capabilities through two primary command hierarchies:

```
gitmap
├── schedule (aliases: sc, crontab, cron)
│   ├── mail-check (alias: mail)     -> Multi-mailbox inspection & scheduled polling
│   ├── email                        -> Automated recurring/delayed email dispatcher
│   ├── watch-folder                 -> Event-driven filesystem trigger to notifier
│   └── ui                           -> Embedded Schedule & Notifier Settings Web UI
└── mail (alias: email)
    ├── send                         -> Direct email transmission via SMTP / lemon
    ├── check                        -> Interactive mailbox status & unread query
    ├── test                         -> Connectivity & authentication diagnostic
    ├── template                     -> Email template manager (list, get, save, del)
    ├── notify                       -> Manual trigger, watch-rule management, logs
    └── ui                           -> Embedded Mailing Composer Web UI
```

---

## 2. Schedule Command Contracts (`cmdschedule`)

### 2.1 `gitmap schedule mail-check` (alias: `gitmap schedule mail`)
Performs scheduled or periodic multi-mailbox polling to track unread message counts, status health, and message digests.

#### Syntax
```bash
gitmap schedule mail-check [flags]
gitmap sc mail [flags]
```

#### Flags
| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--all` | `-a` | boolean | `true` | Poll all active configured mailboxes in `MailboxConfig`. |
| `--mailbox` | `-m` | string | `""` | Target specific mailbox name or account identifier. |
| `--interval` | `-i` | string | `""` | Register recurring polling timer (e.g., `5m`, `1h`, `12h`). |
| `--unread-only` | `-u` | boolean | `false` | Filter summary output to mailboxes with unread messages > 0. |
| `--mask` | | boolean | `true` | Enforce privacy masking on sender and account addresses. |
| `--limit` | `-l` | integer | `10` | Maximum recent headers/messages to inspect per mailbox. |
| `--json` | `-j` | boolean | `false` | Emit output in Universal JSON Response Envelope. |

#### Output Example (Terminal Masked)
```text
[GITMAP MAIL-CHECK] Evaluated 2 active mailboxes at 2026-10-20 14:30:00
+-------------------+-------------+--------------+--------+---------------------+
| MAILBOX ACCOUNT   | UNREAD / TOT| STATUS       | MASKED | LAST CHECKED        |
+-------------------+-------------+--------------+--------+---------------------+
| a***n@c***y.com   | 3 / 142     | OK           | YES    | 2026-10-20 14:30:00 |
| d***s@g***p.io    | 0 / 89      | OK           | YES    | 2026-10-20 14:30:00 |
+-------------------+-------------+--------------+--------+---------------------+
```

---

### 2.2 `gitmap schedule email`
Registers a scheduled, recurring, or delayed email task in GitMap's scheduler database.

#### Syntax
```bash
gitmap schedule email [name] [flags]
```

#### Flags
| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--to` | `-t` | string | required | Comma-separated recipient addresses (masked in logs). |
| `--subject` | `-s` | string | required | Email subject line (supports template tags). |
| `--body` | `-b` | string | `""` | Inline email body content. |
| `--file` | `-f` | string | `""` | Relative path to body file or markdown document. |
| `--template` | | string | `""` | Named email template registered in template library. |
| `--interval` | `-i` | string | `""` | Recurring execution interval (e.g. `24h`, `7d`). |
| `--cron` | `-c` | string | `""` | Standard 5-field cron expression (e.g. `0 9 * * 1-5`). |
| `--delay` | `-d` | string | `""` | One-shot execution delay (e.g. `30m`, `2h`). |
| `--mailbox` | `-m` | string | `"default"` | Account identifier used as dispatch sender. |
| `--json` | `-j` | boolean | `false` | Emit Universal JSON Response Envelope. |

#### Exit Codes
- `0`: Task registered successfully in `ScheduledEmail`.
- `1`: Validation error (missing required recipient or subject).
- `2`: Invalid cron expression or duration parser failure.

---

### 2.3 `gitmap schedule watch-folder`
Configures a filesystem directory watcher that monitors file modifications, additions, or deletions and dispatches email notifications.

#### Syntax
```bash
gitmap schedule watch-folder [path] [flags]
```

#### Flags
| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--to` | `-t` | string | required | Comma-separated notifier recipient emails. |
| `--events` | `-e` | string | `"all"` | Monitored events: `create`, `modify`, `delete`, or `all`. |
| `--debounce` | | string | `"3s"` | Debounce burst grouping window (min `1s`, max `30s`). |
| `--poll-interval` | | string | `"8s"` | Stat poller interval (enforced `IntervalGuard >= 5s`). |
| `--filter` | | string | `"*"` | Glob pattern filter for target files (e.g. `*.go`, `*.json`). |
| `--recursive` | `-r` | boolean | `true` | Recursively monitor nested directory structures. |
| `--rule-name` | `-n` | string | `""` | Descriptive label for registered watch rule. |
| `--json` | `-j` | boolean | `false` | Emit Universal JSON Response Envelope. |

---

## 3. Dedicated Mail Command Contracts (`cmdmail`)

### 3.1 `gitmap mail send`
Directly dispatches an email message via configured SMTP transport or Lemon background daemon.

#### Syntax
```bash
gitmap mail send [flags]
```

#### Flags
| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--to` | `-t` | string | required | Primary recipient email address(es). |
| `--subject` | `-s` | string | required | Email subject line. |
| `--body` | `-b` | string | `""` | Plaintext or HTML email body. |
| `--file` | `-f` | string | `""` | Path to body text or markdown template file. |
| `--cc` | | string | `""` | Carbon copy recipient address(es). |
| `--bcc` | | string | `""` | Blind carbon copy recipient address(es). |
| `--template` | | string | `""` | Template identifier to compile and render. |
| `--account` | `-a` | string | `"default"` | Sender account profile from `MailboxConfig`. |
| `--dry-run` | | boolean | `false` | Validate compilation and transport without sending. |
| `--json` | `-j` | boolean | `false` | Emit Universal JSON Response Envelope. |

---

### 3.2 `gitmap mail check`
Interactively checks connection, inbox status, and recent message headers for specified mailbox accounts.

#### Syntax
```bash
gitmap mail check [account-name] [flags]
```

#### Flags
| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--all` | `-a` | boolean | `false` | Check all configured accounts sequentially. |
| `--unread` | `-u` | boolean | `false` | Display unread message headers only. |
| `--limit` | `-l` | integer | `5` | Maximum number of message snippets to display. |
| `--mask` | | boolean | `true` | Mask recipient/sender email addresses in output. |
| `--json` | `-j` | boolean | `false` | Emit Universal JSON Response Envelope. |

---

### 3.3 `gitmap mail test`
Diagnoses transport connectivity, TLS handshakes, and credential authorization for mail endpoints.

#### Syntax
```bash
gitmap mail test [flags]
```

#### Flags
| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--account` | `-a` | string | `"default"` | Account profile to validate. |
| `--smtp-host`| | string | `""` | Override target SMTP hostname. |
| `--smtp-port`| | integer | `0` | Override target SMTP port (e.g. 587, 465). |
| `--to` | `-t` | string | `""` | Optional test destination to send verification probe. |
| `--json` | `-j` | boolean | `false` | Emit Universal JSON Response Envelope. |

---

### 3.4 `gitmap mail template`
Manages reusable email templates in GitMap's local store.

#### Syntax
```bash
gitmap mail template list [flags]
gitmap mail template get <template-name> [flags]
gitmap mail template save <template-name> --file <path> [flags]
gitmap mail template del <template-name> [flags]
```

---

### 3.5 `gitmap mail notify`
Inspects, triggers, or reviews folder watcher notification logs and active watch rules.

#### Syntax
```bash
gitmap mail notify list [flags]
gitmap mail notify trigger <rule-name> [flags]
gitmap mail notify history [flags]
```

---

### 3.6 Embedded Web UI Commands
Starts embedded zero-CDN web servers for graphical management.

#### Syntax
```bash
gitmap mail ui [--port 9095] [--host 127.0.0.1] [--open-browser]
gitmap schedule ui [--port 9096] [--host 127.0.0.1] [--open-browser]
```

---

## 4. Universal JSON Response Envelope Contract

When `--json` or `-j` is supplied, commands emit structured JSON adhering to GitMap's Universal Response Envelope:

```json
{
  "Status": "success",
  "Attributes": {
    "Command": "schedule-mail-check",
    "ExecutedAt": "2026-10-20T14:30:00Z",
    "IsMasked": true,
    "TotalCount": 2
  },
  "Results": [
    {
      "AccountName": "corporate-primary",
      "SenderAddress": "a***n@c***y.com",
      "UnreadCount": 3,
      "TotalCount": 142,
      "IsActive": true,
      "LastChecked": "2026-10-20T14:30:00Z"
    }
  ]
}
```

---

## 5. Exit Codes & AppError Taxonomy

| Exit Code | AppError Code | Description |
|---|---|---|
| `0` | `nil` | Command completed successfully. |
| `1` | `ErrInvalidArguments` | Missing required parameters, flags, or invalid formats. |
| `2` | `ErrMailboxNotFound` | Specified mailbox profile does not exist in `MailboxConfig`. |
| `3` | `ErrTransportAuthFailed` | SMTP/IMAP TLS or credential authentication failure. |
| `4` | `ErrWatchPathNotFound` | Watcher target directory does not exist or lacks read permission. |
| `5` | `ErrTemplateNotFound` | Requested template identifier does not exist. |
| `6` | `ErrDatabaseOperation` | Split-DB lock timeout or SQLite query error. |
