# Mail & Scheduling Subsystem E2E Verification Matrix

**Specification ID:** `05-mail-and-schedule-e2e-matrix`  
**Parent Specification:** `02-spec/21-app/gitmap-testing-and-schedule-management/readme.md`  
**Parent Plan:** `.ai-memory/plans/gitmap-testing-and-schedule-management.md`  
**Target Subsystems:** `cli/cmdmail/`, `cli/cmdschedule/`  
**Status:** APPROVED  

---

## 1. Overview & Test Architecture

This document defines the comprehensive end-to-end (E2E) verification test matrix for GitMap's mail management and task scheduling subsystems. All tests execute in a hermetic environment without external network calls, utilizing mock email transmitters, in-memory/isolated SQLite databases, and synthetic filesystem structures.

```
Subsystem Integration Topology:
┌────────────────────────────────┐       ┌────────────────────────────────┐
│       cli/cmdmail Engine       │       │    cli/cmdschedule Engine      │
│  - Folder Watcher Lifecycle    │       │  - Task Schedule Email         │
│  - Split-DB Schema & CRUD      │       │  - Mailbox Status Inspection   │
│  - CLI Routing & Templates     │       │  - Privacy Masking Engine      │
└───────────────┬────────────────┘       └───────────────┬────────────────┘
                │                                        │
                └───────────────┬────────────────────────┘
                                │
                 ┌──────────────▼──────────────┐
                 │  Zero-CDN Web UI Dashboards │
                 │  - /mail/ui  &  /schedule/ui│
                 │  - 405 Method Guards & REST │
                 └─────────────────────────────┘
```

---

## 2. Folder Watcher Mutation Lifecycle Matrix (`cli/cmdmail/`)

The folder watcher detects file modifications via directory snapshots (`ScanFolderTree`) and event diffing (`DetectFileMutations`). Tests verify event identification, interval guarding, and notification suppression.

| Test ID | Scenario | Input / Action | Expected Event Output | Verification Assertions |
|---|---|---|---|---|
| `TC-FW-01` | File Creation | Add `new_record.log` (128 bytes) to clean folder | `EventCreate` (`"create"`) | `len(events) == 1`, `events[0].EventType == EventCreate`, `events[0].Path == "new_record.log"`, `SizeBytes == 128` |
| `TC-FW-02` | File Modification | Change byte length of `notes.txt` from 50 to 120 | `EventModify` (`"modify"`) | `len(events) == 1`, `events[0].EventType == EventModify`, `events[0].Details` matches `"Modified (size: 50 -> 120)"` |
| `TC-FW-03` | File Deletion | Remove `obsolete.tmp` from monitored directory | `EventDelete` (`"delete"`) | `len(events) == 1`, `events[0].EventType == EventDelete`, `events[0].Path == "obsolete.tmp"`, `SizeBytes == oldSize` |
| `TC-FW-04` | Zero Events | Execute cycle with unchanged directory state | Empty slice (`[]`) | `len(events) == 0`, `batch.HasEvents == false`, zero email dispatches via mock transmitter |
| `TC-FW-05` | Interval Clamp | Pass `requested = 2 * time.Second` | Clamped to 5s minimum | `guard.WasClamped == true`, `guard.EffectiveInterval == 5 * time.Second`, warning emitted |
| `TC-FW-06` | Default Interval | Pass `requested = 0` (unspecified) | Fallback default 8s | `guard.WasClamped == false`, `guard.EffectiveInterval == 8 * time.Second` |
| `TC-FW-07` | Excluded Paths | Add files under `.git/` or `node_modules/` | Ignored | `isSkippedPath` returns true, files excluded from `ScanFolderTree` snapshot |

---

## 3. SQLite Split-DB CRUD Operations Matrix

Tests verify transactional integrity across PascalCase tables, strictly typed integers for positive booleans, and conflict handling.

| Table Name | Operation | Input Signature | Expected DB State | Verification Query |
|---|---|---|---|---|
| `FolderWatcherRule` | Save New | `RuleName: "ci-logs"`, `FolderPath: "./logs"`, `IsEnabled: true` | Record inserted with positive boolean | `SELECT IsEnabled FROM FolderWatcherRule WHERE ruleName = 'ci-logs'` equals `1` |
| `FolderWatcherRule` | Update Conflict | Same `ruleName` with `debounceSeconds: 15` | Existing record updated, zero duplicates | `SELECT debounceSeconds FROM FolderWatcherRule WHERE ruleName = 'ci-logs'` equals `15` |
| `FolderWatcherRule` | List | Query all configured rules | Slice of `FolderWatcherRuleRecord` sorted DESC | `len(rules) >= 1`, `IsEnabled` decoded as Go `true` |
| `NotificationLog` | Record | `TriggerSource: "watcher"`, `IsDelivered: true` | Immutable log row saved with timestamp | `SELECT IsDelivered FROM NotificationLog WHERE logId = ?` equals `1` |
| `NotificationLog` | List Limit | Query with `limit: 5` | Maximum 5 rows returned in reverse chronological order | `len(logs) <= 5`, records ordered by `createdAt DESC` |
| `MailTemplate` | Save / Upsert | `TemplateName: "alert"`, `IsHTML: true` | Template stored with positive boolean | `SELECT IsHTML FROM MailTemplate WHERE templateName = 'alert'` equals `1` |
| `MailTemplate` | Get By Name | `GetTemplate("alert")` | Matching template record returned | `t.TemplateName == "alert"`, `t.IsHTML == true` |
| `MailTemplate` | Delete | `DeleteTemplate("alert")` | Record successfully removed | Subsequent `GetTemplate("alert")` returns `sql.ErrNoRows` |
| `MailboxProfile` | Save & List | `Alias: "dev"`, `IsActive: true`, `IsMasked: true` | Profile stored with active and masked flags | `SELECT IsActive, IsMasked FROM MailboxProfile` returns `(1, 1)` |

---

## 4. Mail Dispatch & Format Variations Matrix

Tests verify command-line flag handling in `gitmap mail send`, payload building, and subcommand routing.

| Subcommand / Flag | Ingestion / Context | Expected Payload / MIME Structure | Assertions |
|---|---|---|---|
| `mail send --body` | `--to "u@x.com" --subject "Test" --body "Hello"` | Standard RFC-822 plain text envelope | `Content-Type: text/plain; charset=UTF-8`, body contains `"Hello"` |
| `mail send --file` | `--file ./fixtures/msg.txt` | Reads external file content into message body | File contents parsed, errors properly handled if path does not exist |
| `mail send --html` | `--html --body "<h1>Report</h1>"` | HTML MIME envelope | `Content-Type: text/html; charset=UTF-8`, header contains HTML declaration |
| `mail send --dry-run` | `--dry-run --to "admin@org.com"` | Simulates validation, suppresses transport | Mock transmitter receives zero dispatches, output contains `DRY RUN:` |
| `mail template ls` | Subcommand dispatch | Reads template records from `MailStore` | Formatted console list, exits with code 0 |
| `mail template view` | `--name "deploy"` | Fetches and displays template details | Outputs Subject, HTML boolean, and Body |
| `mail notify watch` | `--path ./dist --once` | Executes single pass snapshot comparison | Runs `ExecuteWatchCycle`, logs completed cycle count |
| `mail notify ls` | Subcommand dispatch | Reads rules and notification history | Formats table with Rule ID, Subject, and Delivery status |

---

## 5. Scheduler Email & Mail-Check Flags Matrix

Tests verify interval validation, flag combinations, and empty dataset handling in `cli/cmdschedule/`.

| Subcommand | Flag Variations | Test Condition | Expected Behavior & Output |
|---|---|---|---|
| `schedule email` | `--to "a@b.com" --subject "S" --every "1h"` | Valid flags provided | Inserts `SchedulerTask` into SQLite DB with `IntervalVal: "1h"`, status success |
| `schedule email` | Missing `--every` flag | Invalid parameters | Returns `ValidationError` with message `"missing required flag: --every"` |
| `schedule email` | Missing `--to` flag | Invalid parameters | Returns `ValidationError` with message `"missing required flag: --to"` |
| `schedule email` | Missing `--subject` flag | Invalid parameters | Returns `ValidationError` with message `"missing required flag: --subject"` |
| `schedule mail-check` | `-a` / `--all` | Checks all configured mailboxes | Evaluates all active profiles, renders ASCII summary table |
| `schedule mail-check` | `-u` / `--unread-only` | Filters mailboxes with unread messages | Omits mailboxes having `UnreadCount == 0` |
| `schedule mail-check` | `-j` / `--json` | Universal JSON envelope requested | Emits valid JSON with top-level `Status: "success"`, `Attributes`, and `Results` |
| `schedule mail-check` | Empty Table | Zero mailboxes configured in DB | Formatted ASCII table containing `(no active mailboxes matched query)` |

---

## 6. Privacy Masking & Redaction Edge Cases

Privacy masking protects PII across terminal streams, log records, and REST API envelopes.

| Email Input | Masking Condition | Expected Output | Rule Mechanics |
|---|---|---|---|
| `a@b.com` | Standard default masking | `a*@b*.com` | Single-character local and domain: `s + "*"` preserved before TLD |
| `ab@cd.org` | Standard default masking | `a*@c*.org` | Two-character segment: first character preserved, appended with `*` |
| `abc@def.net` | Standard default masking | `a*c@d*f.net` | Three-character segment: first and third characters preserved |
| `john.doe@company.com` | Standard default masking | `j***e@c***y.com` | Multi-character segment: boundary characters preserved with `***` |
| `ops@internal.corp.com` | Multi-level subdomains | `o*s@i******l.c***p.com` | Each subdomain label masked independently; TLD preserved |
| `admin@company.com` | `--show-full-email` flag set | `admin@company.com` | Explicit CLI override displays unmasked address |
| `admin@company.com` | `GITMAP_ALLOW_UNMASKED_PII=1` | `admin@company.com` | Environment variable disables masking across all subcommands |
| `invalid-email-format` | Malformed address (no `@`) | `***` | Invalid formats fallback to complete sanitization |

---

## 7. Zero-CDN Embedded Web UIs Verification Matrix

GitMap serves embedded administration dashboards for mail composition and scheduling configuration without external dependencies.

| Dashboard Route | Verification Target | Test Method | Expected Result |
|---|---|---|---|
| `/mail/ui` | Zero External CDN Assets | Regex scan for external URI schemes | Zero occurrences of `http://`, `https://`, or scheme-relative `//` |
| `/schedule/ui` | Zero External CDN Assets | Regex scan for external URI schemes | Zero occurrences of `http://`, `https://`, or scheme-relative `//` |
| `/mail/ui` & `/schedule/ui` | Dark Theme Color Tokens | Body CSS Inspection | Canvas background `#090d16`, accent `#38bdf8`, text `#e2e8f0` |
| `/api/mail/send` | HTTP Method Guard | `GET`, `PUT`, `DELETE` requests | `405 Method Not Allowed`, JSON error: `"POST method required"` |
| `/api/mail/send` | Payload Validation | `POST` with empty `to` array | `400 Bad Request`, JSON error: `"At least one recipient is required"` |
| `/api/schedule/toggle` | HTTP Method Guard | `GET`, `PUT`, `DELETE` requests | `405 Method Not Allowed`, JSON error: `"POST method required"` |
| `/api/schedule/watchers`| HTTP Method Guard | `PUT`, `DELETE` requests | `405 Method Not Allowed`, JSON error: `"GET or POST required"` |
| `/api/schedule/tasks` | REST Task Listing | `GET /api/schedule/tasks` | `200 OK`, JSON array of `ScheduleTaskItem` objects |

---

## 8. Automated Test Execution Commands

To execute the test suite covering this matrix:

```bash
# Execute cmdmail test suite with race detector
go test -v -race ./cli/cmdmail/...

# Execute cmdschedule test suite with race detector
go test -v -race ./cli/cmdschedule/...

# Full repository validation without network dependency
go test -short ./...
```
