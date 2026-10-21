# Subtask 05: Comprehensive Verification of New Mail & Scheduling Features

- **Subtask ID:** `05-subtask-new-features-e2e-verification`
- **Parent Task:** `gitmap-testing-and-schedule-management`
- **Specification Reference:** `02-spec/21-app/gitmap-testing-and-schedule-management/05-mail-and-schedule-e2e-matrix.md`
- **Acceptance Reference:** `02-spec/21-app/gitmap-testing-and-schedule-management/06-verification-and-acceptance-criteria.md`
- **Assigned Worker:** Worker 03
- **Status:** DONE

---

## 1. Owned Files (Exclusive)

The following test suite files are strictly owned by Worker 03:
- `cli/cmdmail/cmd_mail_feature_test.go`
- `cli/cmdschedule/schedule_feature_test.go`

Worker 03 may also update existing test fixtures or helpers in `cli/cmdmail/` and `cli/cmdschedule/` if required to assert coverage across newly added features.

---

## 2. Read-Only Context Files

- `02-spec/21-app/gitmap-testing-and-schedule-management/05-mail-and-schedule-e2e-matrix.md`
- `02-spec/21-app/gitmap-testing-and-schedule-management/06-verification-and-acceptance-criteria.md`
- `02-spec/21-app/gitmap-testing-and-schedule-management/readme.md`
- `.ai-memory/plans/gitmap-testing-and-schedule-management.md`
- `cli/cmdmail/folder_watcher.go`
- `cli/cmdmail/storage.go`
- `cli/cmdmail/storage_rules.go`
- `cli/cmdmail/mail_sender.go`
- `cli/cmdmail/cmd_mail.go`
- `cli/cmdmail/cmd_mail_notify.go`
- `cli/cmdmail/mail_ui_server.go`
- `cli/cmdschedule/schedule_email.go`
- `cli/cmdschedule/schedule_mail_check.go`
- `cli/cmdschedule/privacy_masker.go`
- `cli/cmdschedule/schedule_ui_server.go`

---

## 3. Subtask Objectives & Step-by-Step Test Execution Plan

### Step 1: Folder Watcher Mutation Lifecycle (`cli/cmdmail/cmd_mail_feature_test.go`)
1. Implement `TestFolderWatcherMutationLifecycle`:
   - Initialize baseline snapshot with synthetic file map (`notes.txt`, `config.json`).
   - Simulate addition of `new_file.go`: assert `DetectFileMutations` outputs `EventCreate` with matching path and size.
   - Simulate modification of `notes.txt` (changed size and timestamp): assert output contains `EventModify` with details `"Modified (size: ... -> ...)"`.
   - Simulate removal of `config.json`: assert output contains `EventDelete`.
   - Simulate unchanged snapshot: assert output returns empty slice (`len(events) == 0`).
2. Implement `TestIntervalGuardSafety`:
   - Pass requested intervals of `0s`, `1s`, `2s`, `4s`: assert `IntervalGuard.WasClamped == true` and `EffectiveInterval == 5 * time.Second`.
   - Pass valid requested interval `10s`: assert `WasClamped == false` and `EffectiveInterval == 10 * time.Second`.

### Step 2: Split-DB Schema & CRUD Verification (`cli/cmdmail/cmd_mail_feature_test.go`)
1. Initialize an isolated in-memory or temporary SQLite test store via `OpenMailStore(t.TempDir() + "/test_mail.db")`.
2. Test `FolderWatcherRule`:
   - Insert new rule (`ruleName: "test-rule"`, `IsEnabled: true`); assert retrieval returns `IsEnabled == true`.
   - Update rule with conflicting `ruleName`: assert on-conflict update overwrites attributes without duplicate keys.
3. Test `NotificationLog`:
   - Save notification record (`IsDelivered: true`, `triggerSource: "test"`).
   - List logs with limit: assert records sorted `createdAt DESC` and `IsDelivered` decodes as Go boolean `true`.
4. Test `MailTemplate`:
   - Save template with `IsHTML: true`.
   - Query by name via `GetTemplate`: assert subject, body, and boolean match.
   - Delete template: assert removal and subsequent `ErrNoRows`.
5. Test `MailboxProfile`:
   - Save profile with `IsActive: true`, `IsMasked: true`.
   - List profiles: assert positive boolean integer values `(1, 1)` map to Go `true`.

### Step 3: Mail Dispatch & CLI Subcommand Routing (`cli/cmdmail/cmd_mail_feature_test.go`)
1. Test `--file` flag handling:
   - Create temporary file containing multi-line message body.
   - Invoke `parseSendArgs` with `--file <path>`: assert `opts.Body` matches file content.
2. Test `--html` flag handling:
   - Invoke `BuildMIMEMessage` with `opts.IsHTML = true`: assert header contains `Content-Type: text/html; charset=UTF-8`.
   - Invoke with `opts.IsHTML = false`: assert `Content-Type: text/plain; charset=UTF-8`.
3. Test CLI routing:
   - Route `gitmap mail template ls`, `view`, `add`, `rm` via `RunMail`: assert zero exit errors.
   - Route `gitmap mail notify watch --once`, `ls`: assert proper execution flow.
4. Test Mock Transmitter:
   - Swap `CurrentTransmitter` with mock implementation in tests.
   - Invoke `runMailSend` with `--dry-run`: assert mock receives zero transmissions.

### Step 4: Schedule Email & Mail-Check Flags (`cli/cmdschedule/schedule_feature_test.go`)
1. Test duration & interval validation in `RunScheduleEmail`:
   - Test missing `--to`: assert validation error `"missing required flag: --to"`.
   - Test missing `--subject`: assert validation error `"missing required flag: --subject"`.
   - Test missing `--every`: assert validation error `"missing required flag: --every"`.
   - Test valid parameters: assert task registered in SQLite with specified interval.
2. Test flag variations for `RunScheduleMailCheck`:
   - Flag `-a` / `--all`: assert `opts.IsAll == true`.
   - Flag `-u` / `--unread-only`: assert `opts.IsUnreadOnly == true`.
   - Flag `-j` / `--json`: assert output emits valid Universal JSON envelope (`Status`, `Attributes`, `Results`).
3. Test empty table rendering:
   - Invoke `renderMailCheckTable` with `len(results) == 0`.
   - Capture stdout: assert presence of `(no active mailboxes matched query)`.

### Step 5: Privacy Masking & PII Redaction (`cli/cmdschedule/schedule_feature_test.go`)
1. Single-character & edge case addresses:
   - `MaskEmail("a@b.com")` -> assert equals `"a*@b*.com"`.
   - `MaskEmail("ab@cd.org")` -> assert equals `"a*@c*.org"`.
   - `MaskEmail("abc@def.net")` -> assert equals `"a*c@d*f.net"`.
   - `MaskEmail("john.doe@company.com")` -> assert equals `"j***e@c***y.com"`.
2. Multi-segment domain masking:
   - `MaskEmail("user@sub.domain.corp.com")`: assert each subdomain label is masked while TLD is intact.
3. Unmasking controls:
   - Set `GITMAP_ALLOW_UNMASKED_PII=1`: assert `ResolveIsMasked(false)` returns `false`.
   - Pass `--show-full-email`: assert `ResolveIsMasked(true)` returns `false`.

### Step 6: Zero-CDN Web Dashboards & HTTP Method Security (`cli/cmdmail/cmd_mail_feature_test.go`, `cli/cmdschedule/schedule_feature_test.go`)
1. Zero external CDN scan:
   - Scan `mailAppHTMLHeader`, `mailAppHTMLBody`, and `BuildScheduleDashboardHTML()`.
   - Assert zero occurrences of `http://`, `https://`, or scheme-relative `//` URLs for external assets.
2. HTTP 405 Method Not Allowed validation:
   - Test `GET /api/mail/send`: assert response code `405 Method Not Allowed` and error message `"POST method required"`.
   - Test `GET /api/schedule/toggle`: assert response code `405 Method Not Allowed`.
   - Test `PUT /api/schedule/watchers`: assert response code `405 Method Not Allowed`.

---

## 4. Acceptance Criteria & Quality Gates

- [ ] All feature unit tests pass cleanly: `go test -v -race ./cli/cmdmail/... ./cli/cmdschedule/...`.
- [ ] Zero real OS shutdown or reboot commands executed.
- [ ] Zero external network calls made during test execution.
- [ ] Every individual file strictly $\le 300$ lines with 100% relative Git paths.
