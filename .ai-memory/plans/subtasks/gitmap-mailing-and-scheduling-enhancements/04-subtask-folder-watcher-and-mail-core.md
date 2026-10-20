# Subtask 04: Folder Watcher Event Trigger & Mail CLI Engine

- **Task ID:** `Subtask-04`
- **Subtask Slug:** `04-subtask-folder-watcher-and-mail-core`
- **Assigned Worker:** Worker 02
- **Status:** DONE
- **Parent Plan:** `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`
- **Target Specifications:**
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/02-cli-command-contracts.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/03-folder-watcher-and-triggers.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/06-storage-and-split-db-schema.md`

---

## 1. Scope & Exclusive File Ownership

### 1.1 Owned Files (Strictly Disjoint — Exclusive to Worker 02)
- `cli/cmdmail/cmd_mail.go`
- `cli/cmdmail/mail_sender.go`
- `cli/cmdmail/mail_checker.go`
- `cli/cmdmail/folder_watcher.go`
- `cli/cmdmail/storage.go`

### 1.2 Read-Only Context References
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/03-folder-watcher-and-triggers.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/04-privacy-and-email-masking.md`
- `cli/cmdspace/send_email_smtp.go`
- `cli/store/schedule_split_db.go`
- `cli/apperror/apperror.go`

---

## 2. Core Objectives

1. **Dedicated Mailing CLI Verb (`cmd_mail.go`):** Implement the root Cobra command `gitmap mail` and subcommands `send`, `check`, `test`, `template`, `notify`.
2. **Mail Delivery Engine (`mail_sender.go`):** Author an extensible SMTP and API delivery pipeline supporting plain-text, HTML bodies, attachments, template variables, and retry backoff.
3. **Mailbox Inspection Engine (`mail_checker.go`):** Author multi-protocol mailbox querying (IMAP/POP3/REST) for health checks, unread message counts, and snippet fetching.
4. **Folder Watcher Engine (`folder_watcher.go`):** Implement cross-platform directory polling using `os.Stat`/`ModTime`, `IntervalGuard` (min 5s, default 8s), burst event debouncing (3s quiet window), and event trigger dispatching.
5. **Split-DB Storage Engine (`storage.go`):** Implement persistence methods for `EmailWatcherRule`, `EmailWatcherSnapshot`, and `MailHistoryLog` in SQLite Split-DB.

---

## 3. Step-by-Step Implementation Work Breakdown

### Phase 3.1: Top-Level `gitmap mail` Hierarchy (`cmd_mail.go`)
- **Action 3.1.1:** Author root command `MailCmd`:
  - Usage: `gitmap mail [command]`
  - Short: "Manage email delivery, mailbox inspection, templates, and event notification."
- **Action 3.1.2:** Implement and wire subcommands:
  - `mail send`: Send ad-hoc email via SMTP credentials.
  - `mail check`: Check configured mailboxes for unread messages.
  - `mail test`: Verify SMTP/IMAP server connectivity.
  - `mail template`: List, inspect, and preview email templates.
  - `mail notify`: Trigger notification dispatches (including folder watcher subcommands).
- **Action 3.1.3:** Register non-nil stubs for autocompletion parity in `cli/completion`.

### Phase 3.2: SMTP Delivery Pipeline (`mail_sender.go`)
- **Action 3.2.1:** Define mail payload model:
  ```go
  type MailMessage struct {
      From        string            `json:"from"`
      To          []string          `json:"to"`
      Cc          []string          `json:"cc,omitempty"`
      Bcc         []string          `json:"bcc,omitempty"`
      Subject     string            `json:"subject"`
      BodyText    string            `json:"bodyText"`
      BodyHTML    string            `json:"bodyHtml,omitempty"`
      Attachments []string          `json:"attachments,omitempty"`
      Headers     map[string]string `json:"headers,omitempty"`
  }
  ```
- **Action 3.2.2:** Author `SendEmail(msg MailMessage) (*DeliveryReport, error)` with TLS handshake, authentication fallback, and exponential retry.
- **Action 3.2.3:** Add template variable interpolation helper (`{{.TargetFolder}}`, `{{.TotalEvents}}`, `{{.Timestamp}}`).

### Phase 3.3: Mailbox Status Inspector (`mail_checker.go`)
- **Action 3.3.1:** Implement `CheckMailbox(accountName string) (*MailboxStatus, error)`.
- **Action 3.3.2:** Implement credential resolution from `installation.db` or `repo-secrets/11-go-email-reader/`.
- **Action 3.3.3:** Return unread count, recent subjects, and latency metrics.

### Phase 3.4: Folder Watcher & Event Debouncer (`folder_watcher.go`)
- **Action 3.4.1:** Author `IntervalGuard` enforcing `MinPollingInterval = 5s` and `DefaultPollingInterval = 8s`.
- **Action 3.4.2:** Implement `DirectoryScanner`:
  - Walk directory tree using `filepath.WalkDir`.
  - Capture map of relative paths to `FileItemState` (size, `ModTime`, `IsDir`).
  - Compare with prior snapshot to detect `create`, `modify`, `delete` events.
- **Action 3.4.3:** Implement `EventDebouncer`:
  - Accumulate events in thread-safe buffer.
  - Use `time.NewTimer(DefaultDebounceWindow)` (3s) reset on each mutation.
  - When timer fires, flush batch and invoke `TriggerCallback`.
- **Action 3.4.4:** Connect trigger dispatch to `mail_sender.SendEmail()` with formatted activity digest.

### Phase 3.5: Storage & Schema Persistence (`storage.go`)
- **Action 3.5.1:** Author schema initialization for `EmailWatcherRule`, `EmailWatcherSnapshot`, and `MailHistoryLog`.
- **Action 3.5.2:** Implement CRUD operations:
  - `SaveWatcherRule(rule WatcherConfig) error`
  - `ListWatcherRules() ([]WatcherConfig, error)`
  - `DeleteWatcherRule(ruleID string) error`
  - `SaveSnapshot(ruleID string, items []FileItemState) error`
  - `LoadSnapshot(ruleID string) (map[string]FileItemState, error)`
  - `AppendMailLog(log MailHistoryEntry) error`

---

## 4. Invariants & Guardrails

1. **Interval Clamping:** Polling frequencies requested under 5s must be clamped to 5s to safeguard system resources.
2. **Pure Polling Architecture:** Do not invoke CGo or OS-specific native handles; rely exclusively on standard library `os.Stat` and `filepath.WalkDir`.
3. **Database Conventions:** Tables must be PascalCase (`EmailWatcherRule`), columns camelCase (`ruleId`), and booleans positive (`isActive`, `isRecursive`).
4. **Strict Relative Paths:** All references in comments, code imports, and docs must use relative paths.

---

## 5. Verification Evidence & Quality Gates

- **Gate 1 (CLI Help Hierarchy):** `gitmap mail --help` outputs all 5 subcommands cleanly.
- **Gate 2 (IntervalGuard Enforcement):** Supplying `--interval 2s` logs clamping advisory and runs at 5s interval.
- **Gate 3 (Burst Debouncing):** Rapid file creations within 1s window emit exactly 1 batched notification email.
- **Gate 4 (Split-DB Integrity):** Watcher rules and snapshots round-trip successfully through SQLite Split-DB without corruption.
