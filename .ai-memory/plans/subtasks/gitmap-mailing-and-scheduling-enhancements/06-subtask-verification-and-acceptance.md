# Subtask 06: Targeted Unit Tests & Verification Suite

- **Subtask ID:** `subtask-06-verification-and-acceptance`
- **Parent Task:** `gitmap-mailing-and-scheduling-enhancements`
- **Assigned Worker:** Worker Swarm
- **Status:** DONE

---

## 1. Owned Files (Exclusive)

The following test suite files are strictly owned by Worker Swarm:
- `cli/cmdmail/mail_sender_test.go`
- `cli/cmdmail/folder_watcher_test.go`
- `cli/cmdmail/ui_test.go`
- `cli/cmdmail/storage_test.go`
- `cli/cmdschedule/schedule_mail_check_test.go`
- `cli/cmdschedule/schedule_email_test.go`
- `cli/cmdschedule/privacy_masker_test.go`
- `cli/cmdschedule/ui_test.go`

---

## 2. Read-Only Context Files

- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/07-verification-and-acceptance-matrix.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/05-embedded-web-uis.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/06-storage-and-split-db-schema.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/02-cli-command-contracts.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/03-folder-watcher-and-triggers.md`
- `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/04-privacy-and-email-masking.md`

---

## 3. Subtask Objectives & Step-by-Step Test Execution Plan

### Step 1: Multi-Mailbox & Scheduler Tests (`cli/cmdschedule/`)
1. Implement `cli/cmdschedule/schedule_mail_check_test.go`:
   - Mock multi-mailbox account registry with 3 distinct accounts.
   - Test `schedule mail-check --all` parses flags, runs concurrent polling routines, and formats summary table.
   - Assert error handling for offline mailboxes does not crash other mailbox checks.
2. Implement `cli/cmdschedule/schedule_email_test.go`:
   - Test interval parsing (e.g. `1h`, `30m`, `8s`) and cron parsing (`0 9 * * 1`).
   - Test scheduling database insertion in `MailSchedule` table.
   - Test dispatch triggering with mock mail dispatcher.

### Step 2: Privacy Masker Verification (`cli/cmdschedule/privacy_masker_test.go`)
1. Test standard address masking:
   - `john.doe@company.com` -> `j***e@c***y.com`
   - `alice@domain.org` -> `a***e@d***n.org`
2. Test short local parts and edge cases:
   - Single-character local: `a@b.com` -> `a*@b***m`
   - Two-character local: `ab@domain.com` -> `a*@d***n.com`
3. Test batch string & JSON array masking:
   - `MaskEmailList([]string{"a@x.com", "b@y.com"})`
4. Test unmasking security:
   - Ensure unmasking is disallowed unless explicitly enabled via flag.

### Step 3: Folder Watcher & Debounce Suite (`cli/cmdmail/folder_watcher_test.go`)
1. File Creation & Modification:
   - Spin up isolated temporary test directory via `t.TempDir()`.
   - Start watcher routine with 100ms test debounce window.
   - Create 10 files rapidly; assert exactly 1 batch notification event is emitted.
2. Glob Pattern Matching:
   - Test pattern `*.md`: mutations to `notes.md` trigger events; mutations to `app.bin` are ignored.
3. Polling Interval Guard:
   - Assert configured intervals below 5 seconds automatically enforce the minimum 5s guard in production mode.

### Step 4: Storage & Schema Verification (`cli/cmdmail/storage_test.go`)
1. Schema DDL Assertion:
   - Verify PascalCase tables exist: `MailSchedule`, `MailboxProfile`, `FolderWatcherRule`, `NotificationLog`, `MailTemplate`.
2. Positive Booleans Validation:
   - Query table pragma; assert boolean columns (`IsEnabled`, `IsActive`, `IsMasked`, `HasAttachment`, `IsHTML`, `IsDelivered`) are integers.
3. WAL Mode Verification:
   - Execute `PRAGMA journal_mode;` asserting result is `wal`.
4. Concurrency Test:
   - Run parallel writes from 8 goroutines to assert zero lock contention errors.

### Step 5: Embedded Web UI Test Suite (`cli/cmdmail/ui_test.go`, `cli/cmdschedule/ui_test.go`)
1. Zero-CDN Scan:
   - Assert HTML/JS source code contains zero `http://` or `https://` URLs.
2. Dark Theme Verification:
   - Assert CSS includes `#090d16` canvas background and `#38bdf8` accent.
3. REST Endpoint Tests (`httptest.NewServer`):
   - Test `GET /api/mail/accounts` returns masked profiles.
   - Test `POST /api/mail/send` validates payload.
   - Test `GET /api/schedule/watchers` returns active rules.
   - Test `POST /api/schedule/watchers` creates new rule.

---

## 4. Acceptance Sign-Off Criteria

- [ ] All 8 test files compile and pass with 0 failures (`go test ./cli/cmdmail/... ./cli/cmdschedule/...`).
- [ ] No race conditions reported (`go test -race`).
- [ ] No network access required during test execution.
- [ ] Every individual source and test file is strictly <= 300 lines.
