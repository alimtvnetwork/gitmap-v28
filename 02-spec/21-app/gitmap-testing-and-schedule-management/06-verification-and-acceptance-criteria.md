# Verification & Acceptance Criteria Specification

**Specification ID:** `06-verification-and-acceptance-criteria`  
**Parent Specification:** `02-spec/21-app/gitmap-testing-and-schedule-management/readme.md`  
**Parent Plan:** `.ai-memory/plans/gitmap-testing-and-schedule-management.md`  
**Target Subsystems:** Cross-OS Power Scheduling, Docker Harness, Mail Subsystem, Privacy Masker, Embedded UIs  
**Status:** APPROVED  

---

## 1. Quality Gate Architecture

To ensure the reliability, cross-platform portability, and safety of GitMap's scheduling and mailing features, nine mandatory acceptance gates (G1 through G9) are established. Every code modification must satisfy all gates prior to merging and release.

```
Quality Gates Hierarchy:
┌─────────────────────────────────────────────────────────────────┐
│                    Master Quality Gates (G1-G9)                 │
├─────────────────────────┬───────────────────────────────────────┤
│ Cross-OS & Container    │ G1: Cross-OS Command Construction     │
│                         │ G2: Safety Isolation (Zero Host Boot) │
│                         │ G3: Pure Go Linux Compilation         │
│                         │ G4: Docker Container Test Harness     │
├─────────────────────────┼───────────────────────────────────────┤
│ Subsystem Functionality │ G5: Folder Watcher Mutation Lifecycle │
│                         │ G6: Split-DB Schema & Positive Bools  │
│                         │ G7: Privacy Redaction & PII Governance│
│                         │ G8: Zero-CDN Offline Web Dashboards   │
├─────────────────────────┼───────────────────────────────────────┤
│ Codebase Hygiene        │ G9: 300-Line Limit & Relative Paths   │
└─────────────────────────┴───────────────────────────────────────┘
```

---

## 2. Numbered Acceptance Criteria & Verification Gates

### AC-GATE-01 (G1: Cross-OS Power Command Construction)
- **Requirement:** The power scheduling engine must construct valid OS-specific CLI commands for Windows, Linux, and macOS (Darwin) without syntax errors.
- **Verification Rule:**
  * Windows: `shutdown /s /t <seconds> /c <message>` (shutdown), `shutdown /r /t <seconds> /c <message>` (restart), `shutdown /a` (cancel).
  * Linux: `sudo poweroff` / `sudo shutdown -h +<minutes>` (shutdown), `sudo reboot` / `sudo shutdown -r +<minutes>` (restart), `sudo shutdown -c` (cancel).
  * Darwin: `sudo shutdown -h +<minutes>` (shutdown), `sudo shutdown -r +<minutes>` (restart), `sudo killall shutdown` or `sudo shutdown -c` (cancel).
- **Pass Invariant:** Unit tests in `cli/cmdschedule/schedule_os_test.go` must assert exact command argument construction for all three operating systems.

### AC-GATE-02 (G2: Safety Isolation Invariant `cg-isolate-os-tests`)
- **Requirement:** Test execution must NEVER trigger a real host operating system shutdown, reboot, or state corruption.
- **Verification Rule:**
  * All OS commands must be executed via an injectable `OSActionExecutor` interface.
  * Test suites must inject mock executors that record command strings, arguments, and timeouts without invoking `os/exec.Command`.
- **Pass Invariant:** 100% of unit and integration tests run without privileged escalation or system disruptions.

### AC-GATE-03 (G3: Pure Go Linux Cross-Compilation)
- **Requirement:** The entire codebase must compile cleanly for Linux AMD64 without Windows-specific symbols, imports, or build tag leaks.
- **Verification Command:**
  ```bash
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -run=^$ ./cli/cmdschedule/... ./cli/cmdmail/...
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./...
  ```
- **Pass Invariant:** Both commands exit with status code `0` and zero compilation warnings.

### AC-GATE-04 (G4: Containerized Docker Test Harness)
- **Requirement:** A hermetic, self-contained `Dockerfile.test` based on `golang:1.24-alpine` must be provided for containerized Linux validation.
- **Verification Rule:**
  * The Docker container must install required build tools, copy repository files, execute `go test -v ./...`, and report exit status.
  * CI/CD environments with Docker provisioned can execute the test suite in full isolation.
- **Pass Invariant:** Container builds and executes tests cleanly without external dependencies.

### AC-GATE-05 (G5: Folder Watcher Mutation Integrity)
- **Requirement:** `cli/cmdmail/` folder watcher must detect file mutations across all three event categories (`EventCreate`, `EventModify`, `EventDelete`), suppress notifications on zero events, and enforce interval safety.
- **Verification Rule:**
  * File additions trigger `EventCreate`.
  * Byte size or modification time changes trigger `EventModify`.
  * Deletions trigger `EventDelete`.
  * Zero mutations produce empty event slices and suppress mail dispatch.
  * `NewIntervalGuard` clamps requested intervals below 5 seconds to `MinPollingInterval = 5 * time.Second`.
- **Pass Invariant:** Tests in `cli/cmdmail/cmd_mail_test.go` assert all five states deterministically.

### AC-GATE-06 (G6: SQLite Split-DB & Schema Compliance)
- **Requirement:** All SQLite tables must use PascalCase naming, integer positive booleans (`isEnabled`, `isActive`, `isMasked`, `isDelivered`, `isHTML`), and WAL journaling mode.
- **Verification Rule:**
  * Tables: `FolderWatcherRule`, `NotificationLog`, `MailTemplate`, `MailboxProfile`.
  * Banned boolean names: `can*`, `should*`, `was*`, `will*`, or negative prefixes (`isNot*`, `isDisabled`).
  * On-conflict upsert clauses must update existing rows without primary key duplication.
- **Pass Invariant:** Table schema pragmas confirm integer columns and 0 constraint violations.

### AC-GATE-07 (G7: Privacy Redaction & PII Governance)
- **Requirement:** Email addresses displayed in CLI output, logs, or JSON responses must be masked by default (`u***r@d***n.com`).
- **Verification Rule:**
  * Boundary masking preserves first and last characters for standard identifiers.
  * Short handles: single-character local (`a@b.com` -> `a*@b*.com`), two-character (`ab@c.com` -> `a*@c*.com`).
  * Bypass controls: `--show-full-email` flag or `GITMAP_ALLOW_UNMASKED_PII=1` environment variable cleanly unmasks output.
- **Pass Invariant:** Unit tests assert masking compliance across all standard and edge-case email inputs.

### AC-GATE-08 (G8: Zero-CDN Offline Web Dashboards & HTTP Method Guards)
- **Requirement:** Web dashboards (`/mail/ui`, `/schedule/ui`) must operate strictly offline without external network CDN dependencies and enforce strict HTTP method security.
- **Verification Rule:**
  * Zero CDN links: HTML, CSS, and JS sources contain zero references to `http://`, `https://`, or scheme-relative `//` CDNs.
  * Method guards: Mutation endpoints (`/api/mail/send`, `/api/schedule/toggle`) return `405 Method Not Allowed` when invoked via unauthorized HTTP methods (e.g. `GET`).
  * Dark mode theme conforms to standard color tokens (`#090d16` canvas, `#38bdf8` accent).
- **Pass Invariant:** HTTP test suite verifies zero CDN strings and validates 405 error responses.

### AC-GATE-09 (G9: Line Cap Budget & Strict Relative Git Paths)
- **Requirement:** All newly authored or modified specification files, plans, and source files must adhere to strict structural constraints.
- **Verification Rule:**
  * File size limit: Every file must be strictly $\le 300$ lines.
  * Path hygiene: 100% relative Git paths (`02-spec/...`, `.ai-memory/...`, `cli/...`). Absolute filesystem paths and URI scheme prefixes are strictly banned.
  * Naming: All filenames must use strict lowercase and hyphens.
- **Pass Invariant:** Automated line count scan and relative path regex verification pass with 0 errors.

---

## 3. Automated Verification Command Protocol

The following deterministic commands execute all verification gates:

```bash
# 1. Run unit tests across mail and scheduling packages
go test -v -race ./cli/cmdmail/...
go test -v -race ./cli/cmdschedule/...

# 2. Verify pure Go Linux cross-compilation
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -run=^$ ./cli/cmdschedule/...
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./...

# 3. Verify zero external CDN links in embedded UI source code
# Ensure no external network schemes exist in UI assets
go test -v -run=TestZeroCDN ./cli/cmdmail/... ./cli/cmdschedule/...

# 4. Verify line limits and relative paths across specification docs
pwsh -Command "Get-ChildItem -Path 02-spec/21-app/gitmap-testing-and-schedule-management -Filter *.md | ForEach-Object { if ((Get-Content \$_.FullName).Count -gt 300) { Write-Error \"Over 300 lines: \$($_.Name)\" } }"
```

---

## 4. Sign-Off & Verification Matrix Checklist

| Gate ID | Quality Gate Description | Responsible Role | Verification Method | Status |
|---|---|---|---|---|
| `G1` | Cross-OS Power Command Construction | Worker 01 | Unit Test Suite (`schedule_os_test.go`) | PENDING |
| `G2` | Safety Invariant Isolation (`cg-isolate`) | Worker 01 | Mock Executor Verification | PENDING |
| `G3` | Linux AMD64 Cross-Compilation | Worker 02 | `GOOS=linux go build ./...` | PENDING |
| `G4` | Docker Container Test Harness | Worker 02 | `Dockerfile.test` Build & Run | PENDING |
| `G5` | Folder Watcher Mutation Lifecycle | Worker 03 | Unit Test Suite (`folder_watcher_test.go`) | PENDING |
| `G6` | SQLite Split-DB PascalCase & Bools | Worker 03 | DB Schema Pragma Verification | PENDING |
| `G7` | Privacy Redaction & PII Governance | Worker 03 | Unit Test Suite (`privacy_masker_test.go`) | PENDING |
| `G8` | Zero-CDN UIs & 405 Method Guards | Worker 03 | HTTP Integration Tests (`ui_test.go`) | PENDING |
| `G9` | Line Cap Budget (<= 300) & Relative Paths | Lead / Spec Authors | Structural Audit Scripts | APPROVED |
