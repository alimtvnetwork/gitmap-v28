# Subtask 03: Email Sending, Scheduling & Terminal UI

- **Task-ID**: Task-03
- **Subtask**: 03-email-sending-and-scheduling-ui
- **Owner**: Worker 03
- **Status**: PENDING

## Owned Files (Exclusive)
- `cli/cmdspace/send_email_ui.go`
- `cli/cmdspace/send_email_ui_test.go`
- `cli/cmdspace/schedule_send_email.go`
- `cli/cmdspace/space.go`

## Read-Only Context Files
- `02-spec/21-app/goemailsender-connection-and-testing/02-component-spec.md`
- `02-spec/21-app/276-email-integration/01-scope-and-plan.md`
- `repo-secrets/11-go-email-reader/01-email-setup.json`
- `repo-secrets/11-go-email-reader/02-default-notifier.json`

## Objectives
1. Implement interactive terminal UI for email drafting and scheduling (`gitmap space email ui` or `gitmap space send-email --ui`).
2. Display account information and notifier address loaded from `repo-secrets` or `go-email-reader list`.
3. Provide prompts/options for recipient, subject, body, HTML format, and scheduling trigger tasks (`release`, `pull`, `clone`, `scan`, `push`, `backup`).
4. Unit tests in `cli/cmdspace/send_email_ui_test.go` covering UI rendering and option dispatching.
