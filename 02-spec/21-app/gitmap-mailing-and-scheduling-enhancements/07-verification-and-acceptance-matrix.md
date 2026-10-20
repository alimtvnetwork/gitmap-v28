# 07-verification-and-acceptance-matrix: Quality Gates & Verification Specification

- **Spec ID:** `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/07-verification-and-acceptance-matrix.md`
- **Subsystem:** Test Suites, Verification Gates & Acceptance Criteria
- **Status:** APPROVED
- **Related Specs:**
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/01-architecture-overview.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/05-embedded-web-uis.md`
  - `02-spec/21-app/gitmap-mailing-and-scheduling-enhancements/06-storage-and-split-db-schema.md`
- **Implementation Subtask:** `.ai-memory/plans/subtasks/gitmap-mailing-and-scheduling-enhancements/06-subtask-verification-and-acceptance.md`

---

## 1. Executive Summary & Verification Protocol

This document defines the test suites, verification gates, and acceptance criteria governing the implementation of GitMap's mailing, multi-mailbox checking, event-driven folder watcher notifications, privacy address masking, and embedded web dashboards.

### 1.1 Non-Negotiable Test Invariants
1. **Zero External Network Dependencies:** All unit and integration tests MUST execute against local mock servers (`httptest.Server`, in-memory SMTP/IMAP fakes). Tests must never attempt real outbound internet connections.
2. **Deterministic Cleanup:** File watcher and storage tests MUST use isolated temporary directories and ensure all timers and goroutines terminate gracefully.
3. **Strict Zero-CDN Guarantee:** Embedded UI tests assert that generated HTML/JS contains zero remote CDN URLs (`http://`, `https://`).
4. **Execution Boundaries:** Tests must run cleanly within bounded execution limits without lingering port bindings.

---

## 2. Quality Gates Matrix

| Gate ID | Area | Success Criteria | Verification Mechanism |
|---|---|---|---|
| **G1** | Multi-Mailbox Checking | Evaluates multiple configured mailboxes and accurately returns unread counts and message headers. | `TestMailboxChecking_MultiAccount` |
| **G2** | Scheduled Email Task | Correctly persists interval/cron schedule records in `MailSchedule` and triggers execution upon interval elapsed. | `TestScheduleEmail_IntervalDispatch` |
| **G3** | Folder Watcher Triggers | File mutations (create, modify, delete) trigger debounced email notifications to configured recipients. | `TestFolderWatcher_DebouncedTrigger` |
| **G4** | Privacy Masking Engine | Masking transforms emails into `u***r@d***n.com` across terminal output, logs, and public REST responses. | `TestPrivacyMasker_OutputSanitization` |
| **G5** | Mail CLI Hierarchy | Verbs (`send`, `check`, `test`, `template`, `notify`, `ui`) parse flags and dispatch correctly. | `TestMailCLI_CommandRouting` |
| **G6** | Embedded Web UIs | Both `/mail/ui` and `/schedule/ui` serve responsive zero-CDN dashboards with dark theme tokens. | `TestEmbeddedUI_ZeroCDNAndTokens` |

---

## 3. Detailed Test Suites

### 3.1 Suite 1: Multi-Mailbox Checking (`cli/cmdschedule/`, `cli/cmdmail/`)
- `TestMailboxChecking_MultiAccount`: Registers 3 mock IMAP accounts; verifies concurrent polling retrieves accurate unread counts without cross-account state pollution.
- `TestMailboxChecking_UnreadFilter`: Verifies filtering by unread status, date range, and sender query.
- `TestMailboxChecking_TimeoutResilience`: Verifies graceful failure and retry backoff when an account connection times out.

### 3.2 Suite 2: Folder Watcher & Event Debouncing (`cli/cmdmail/`, `cli/cmdschedule/`)
- `TestFolderWatcher_FileCreationTrigger`: Creates a new file in a watched directory; verifies notification event fires after debounce window.
- `TestFolderWatcher_BurstDebouncing`: Generates 20 rapid file modifications within a 2-second burst; asserts exactly ONE aggregated notification email is dispatched with the full modified files list.
- `TestFolderWatcher_GlobPatternFilter`: Ensures only files matching the specified pattern (e.g. `*.go`, `*.json`) trigger notifications.
- `TestFolderWatcher_IntervalGuard`: Verifies minimum polling interval of 5 seconds is strictly enforced.

### 3.3 Suite 3: Address Privacy Masking (`cli/cmdschedule/`, `cli/cmdmail/`)
- `TestPrivacyMasker_StandardAddress`: Asserts `developer@example.com` transforms to `d***r@e***e.com`.
- `TestPrivacyMasker_ShortLocalPart`: Asserts single and two-character local parts (e.g. `ab@test.org`) mask properly (`a*@t**t.org`).
- `TestPrivacyMasker_SubdomainSupport`: Asserts `user@mail.corp.com` masks correctly to `u***r@m***p.com`.
- `TestPrivacyMasker_RecipientList`: Masks JSON arrays and comma-separated recipient strings uniformly.
- `TestPrivacyMasker_ExplicitUnmaskFlag`: Asserts full email is revealed ONLY when explicit `--unmask` flag is passed.

### 3.4 Suite 4: Embedded Web UIs (`cli/cmdmail/`, `cli/cmdschedule/`)
- `TestEmbeddedUI_ZeroCDNInvariant`: Scans generated HTML, CSS, and JS payloads for `http://` and `https://` URLs, verifying zero external asset calls.
- `TestEmbeddedUI_DarkThemeTokens`: Asserts presence of required CSS tokens (`#090d16`, `#38bdf8`, `#1e293b`).
- `TestEmbeddedUI_MailComposerREST`: Verifies `POST /api/mail/send` validates recipient formats, handles HTML toggle, and returns structured JSON responses.
- `TestEmbeddedUI_ScheduleSettingsREST`: Tests CRUD lifecycle on `/api/schedule/watchers` and `/api/schedule/jobs`.
- `TestEmbeddedUI_GracefulShutdown`: Asserts server terminates cleanly within 3 seconds when shutdown signal is dispatched.

### 3.5 Suite 5: Split-DB Storage & Schema Concurrency (`cli/cmdmail/`, `cli/cmdschedule/`)
- `TestSplitDB_PascalCaseTables`: Asserts tables `MailSchedule`, `MailboxProfile`, `FolderWatcherRule`, `NotificationLog`, `MailTemplate` exist and match casing.
- `TestSplitDB_PositiveBooleans`: Verifies fields `IsEnabled`, `IsActive`, `IsMasked`, `HasAttachment`, `IsHTML`, `IsDelivered` default to valid positive integers (0 or 1).
- `TestSplitDB_WALModeAndPragmas`: Confirms `journal_mode=WAL` and `foreign_keys=ON` are active on connection.
- `TestSplitDB_ConcurrentWriteSafety`: Spawns 10 concurrent goroutines writing notification logs; asserts zero database lock errors.

---

## 4. Numbered Acceptance Criteria (AC-01 through AC-10)

- **AC-01:** `gitmap schedule mail-check --all` inspects all configured mailboxes and outputs masked summary tables.
- **AC-02:** `gitmap schedule email --cron "0 9 * * 1" --to <recipient>` successfully schedules a recurring automated email.
- **AC-03:** `gitmap schedule watch-folder --path <dir> --events create,modify --notify <emails>` creates an active watcher rule.
- **AC-04:** Watcher triggers consolidate rapid mutations occurring within the `--debounce` window (default: 8s) into a single dispatch.
- **AC-05:** All email addresses displayed in terminal stdout, log files, and UI responses are masked by default.
- **AC-06:** `gitmap mail send --to <addr> --subject <subj> --body <text>` sends an email using the active profile.
- **AC-07:** `gitmap mail ui` launches the embedded Composer SPA on `127.0.0.1:8745/mail/ui` without external CDN references.
- **AC-08:** `gitmap schedule ui` launches the Settings SPA on `127.0.0.1:8746/schedule/ui` without external CDN references.
- **AC-09:** Database migrations idempotently bootstrap Split-DB SQLite tables following PascalCase and positive boolean rules.
- **AC-10:** All unit and integration test suites pass with 100% success rate without internet access.
