# Subtask 01: GoEmailSender Discovery, Windows Path Resolution & Installation Validation

- **Task-ID**: Task-01
- **Subtask**: 01-goemailsender-discovery-and-install
- **Owner**: Worker 01
- **Status**: PENDING

## Owned Files (Exclusive)
- `cli/cmdspace/send_email.go`
- `cli/cmdspace/send_email_windows.go`
- `cli/cmdspace/send_email_test.go`

## Read-Only Context Files
- `02-spec/21-app/goemailsender-connection-and-testing/01-architecture-spec.md`
- `02-spec/21-app/goemailsender-connection-and-testing/02-component-spec.md`
- `02-spec/21-app/276-email-integration/01-scope-and-plan.md`

## Objectives
1. Fix `resolveGoEmailReaderBin` and `isExecutableFile` in `cli/cmdspace/send_email.go` and `cli/cmdspace/send_email_windows.go` to cleanly detect `.exe` on Windows and inspect `os.Getenv("GOPATH")` / `filepath.Join(gopath, "bin")` as well as `os.UserHomeDir()` / `go/bin`.
2. Expand probed candidate binary names to include `go-email-reader`, `GoEmailSender`, and `goemailsender`.
3. Add targeted unit tests in `cli/cmdspace/send_email_test.go` proving resolution on both Windows and Unix.
