# Subtask 01: Spinner Indentation, Padding & `txt` Ghost Removal

## Scope & Intent
- Shift the blue/cyan spinner square rightwards by 2 spaces so column 2 aligns with the white square `■` in `  ■ Probing sibling repos...`.
- Ensure step numbers `[1/5]`, `[4/5]` begin at column 4, aligning with the word `Probing`.
- Add an empty line (`\n`) before the first `SUCCESS` badge so there is a clean top padding gap.
- Eliminate `txt` ghost text and double-flushed stage output by replacing `pterm.MultiPrinter` background loop with sequential stage lifecycle and clean line erasure (`\r\033[2K`).

## Owned Files
- `cli/cmdupdate/update_ui.go`
- `cli/cmdupdate/update_ui_types.go`

## Invariants
- Zero builds or full test suites (`go test ./...` banned).
- Targeted test only: `go test ./cmdupdate -run Test...`.
