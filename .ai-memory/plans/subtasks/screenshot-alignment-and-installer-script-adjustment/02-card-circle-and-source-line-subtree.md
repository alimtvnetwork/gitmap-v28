# Subtask 02: Card Circle Alignment & "Added source line" Subtree

## Scope & Intent
- Fix the broken circle / top-right rounded corner disconnection on the summary card.
  - In `cli/cmdupdate/update_ui_notes.go`, change `✔` (`\u2714`, which expands to 2 columns in Windows Terminal DirectWrite emoji font) to standard 1-column `✓` (`\u2713`), perfectly aligning line 1 right border `│` with `╮` at column 77.
- Transform "Added source line" output in `cli/completion/install.go`:
  - Print `  ✓ Added source line:` once.
  - Render updated profile paths as an indented subtree (`    ├── ` / `    └── `).
  - Print a trailing newline (`\n`) after completing the tree before handing over to the user prompt.
- Update `cli/constants/constants_completion.go` with subtree prefix constants if necessary.

## Owned Files
- `cli/cmdupdate/update_ui_notes.go`
- `cli/completion/install.go`
- `cli/constants/constants_completion.go`

## Invariants
- Zero builds or full test suites (`go test ./...` banned).
- Targeted test only: `go test ./completion -run Test...` and `go test ./cmdupdate -run Test...`.
