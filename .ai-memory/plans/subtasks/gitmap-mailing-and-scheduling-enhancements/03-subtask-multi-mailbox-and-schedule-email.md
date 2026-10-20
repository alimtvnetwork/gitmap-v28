# Subtask 03: Multi-Mailbox Checker, Schedule Email CLI & Privacy Masker

- **Task ID:** `Subtask-03`
- **Subtask Slug:** `03-subtask-multi-mailbox-and-schedule-email`
- **Assigned Worker:** Worker 01
- **Status:** PENDING
- **Parent Plan:** `.ai-memory/plans/gitmap-mailing-and-scheduling-enhancements.md`
- **Target Specifications:**
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/02-cli-command-contracts.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/04-privacy-and-email-masking.md`

---

## 1. Scope & Exclusive File Ownership

### 1.1 Owned Files (Strictly Disjoint — Exclusive to Worker 01)
- `cli/cmdschedule/privacy_masker.go`
- `cli/cmdschedule/schedule_mail_check.go`
- `cli/cmdschedule/schedule_email.go`
- `cli/cmdschedule/schedule_cmd.go`

### 1.2 Read-Only Context References
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/04-privacy-and-email-masking.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/02-cli-command-contracts.md`
- `cli/cmdspace/schedule_send_email.go`
- `cli/apperror/apperror.go`
- `cli/jsonenv/envelope.go`

---

## 2. Core Objectives

1. **Privacy Masking Engine (`privacy_masker.go`):** Implement the deterministic email address masking algorithm (`u***r@d***n.com`), list masking, and unmask authorization flag handling (`--show-full-email`).
2. **Multi-Mailbox Checking (`schedule_mail_check.go`):** Implement `gitmap schedule mail-check` to poll single or multiple configured mailboxes, retrieve unread counts, snippets, and connection health.
3. **Automated Scheduled Emails (`schedule_email.go`):** Implement `gitmap schedule email` to register recurring interval or cron-based email dispatches with subject, body templates, and recipient routing.
4. **Cobra CLI Integration (`schedule_cmd.go`):** Wire `mail-check` and `email` subcommands into GitMap's existing `schedule` command tree, ensuring clean flag inheritance and typed JSON envelopes.

---

## 3. Step-by-Step Implementation Work Breakdown

### Phase 3.1: Privacy Masker Implementation (`privacy_masker.go`)
- **Action 3.1.1:** Author `MaskEmailAddress(raw string) string` implementing the standard:
  - Local part: Keep first char, last char, mask middle with `***`.
  - Domain part: Mask domain name labels, preserve final TLD.
  - Fallback for short addresses (`<= 2` chars) and invalid tokens.
- **Action 3.1.2:** Author `PrivacyMasker` struct with `AllowUnmasked bool` field.
- **Action 3.1.3:** Author `MaskRecipientList(emails []string) []string` and `MaskCommaSeparated(raw string) string`.
- **Action 3.1.4:** Add helper `ResolveUnmaskFlag(cmd *cobra.Command) bool` checking `--show-full-email` flag and `GITMAP_ALLOW_UNMASKED_PII` environment variable.

### Phase 3.2: Multi-Mailbox Checker Subcommand (`schedule_mail_check.go`)
- **Action 3.2.1:** Define mailbox status models:
  ```go
  type MailboxCheckResult struct {
      AccountName  string    `json:"accountName"`
      Address      string    `json:"address"`
      UnreadCount  int       `json:"unreadCount"`
      LastChecked  time.Time `json:"lastChecked"`
      Status       string    `json:"status"`
      ErrorMessage string    `json:"errorMessage,omitempty"`
  }
  ```
- **Action 3.2.2:** Implement mailbox retrieval from Split-DB settings or secret vaults (`repo-secrets/11-go-email-reader/`).
- **Action 3.2.3:** Implement `runScheduleMailCheck(cmd *cobra.Command, args []string) error`:
  - Check `--account` filter or `--all` flag.
  - Apply `PrivacyMasker.MaskAddress()` to all rendered addresses by default.
  - Render terminal results using boxed table formatting (`termtable`).
  - Emit typed JSON envelope when `--json` flag is provided.

### Phase 3.3: Automated Scheduled Email Registration (`schedule_email.go`)
- **Action 3.3.1:** Define scheduled mail task structure:
  ```go
  type ScheduledEmailJob struct {
      JobID          string        `json:"jobId"`
      Subject        string        `json:"subject"`
      Recipients     []string      `json:"recipients"`
      CronExpression string        `json:"cronExpression,omitempty"`
      Interval       time.Duration `json:"interval,omitempty"`
      BodyTemplate   string        `json:"bodyTemplate"`
      IsActive       bool          `json:"isActive"`
      CreatedAt      time.Time     `json:"createdAt"`
  }
  ```
- **Action 3.3.2:** Implement `runScheduleEmail(cmd *cobra.Command, args []string) error`:
  - Flags: `--to`, `--subject`, `--body`, `--template-file`, `--cron`, `--interval`, `--show-full-email`.
  - Validate that either `--cron` or `--interval` is specified.
  - Validate recipient syntax before persistence.
  - Persist task definition to SQLite Split-DB (`installation.db` / `gitmap.db`).
  - Mask recipient output in confirmation banner unless `--show-full-email` is set.

### Phase 3.4: Subcommand Registration in `schedule_cmd.go`
- **Action 3.4.1:** Register `mailCheckCmd` under `scheduleCmd`.
  - Bind flags: `--account`, `--all`, `--show-full-email`, `--json`.
- **Action 3.4.2:** Register `emailCmd` under `scheduleCmd`.
  - Bind flags: `--to`, `--subject`, `--body`, `--template-file`, `--cron`, `--interval`, `--show-full-email`.
- **Action 3.4.3:** Ensure non-nil execution stubs for autocompletion parity (`completion.AllCommands()`).

---

## 4. Invariants & Guardrails

1. **Zero Plaintext PII:** Plaintext emails must never appear in standard stdout or logs unless explicitly authorized.
2. **Strict Relative Paths:** All internal references must use relative Git paths (e.g. `cli/cmdschedule/...`).
3. **PascalCase / camelCase Parity:** SQLite table definitions and JSON model tags must strictly respect naming guidelines.
4. **AppError Returns:** CLI handlers must return structured `apperror.AppError` instances rather than raw `fmt.Errorf`.

---

## 5. Verification Evidence & Quality Gates

- **Gate 1 (Masking Accuracy):** `MaskEmailAddress("john.doe@example.com")` returns `"j***e@e***e.com"`.
- **Gate 2 (Mail Check Execution):** `gitmap schedule mail-check --all` displays formatted table with masked addresses.
- **Gate 3 (Unmask Override):** `gitmap schedule mail-check --all --show-full-email` renders authentic addresses.
- **Gate 4 (Email Schedule Persistence):** `gitmap schedule email --to dev@example.com --interval 1h --subject "Health"` registers the job successfully.
