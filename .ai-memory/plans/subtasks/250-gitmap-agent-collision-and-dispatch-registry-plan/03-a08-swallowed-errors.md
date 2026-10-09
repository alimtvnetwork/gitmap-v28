# Subtask 03 — A08 swallowed-error fix (task 250)

## File box
`cli/cmd/root.go` (`writeLastErrorFile` + `persistLastError` audit ONLY). Nothing else.

## Build (per spec 01 §4)
1. `writeLastErrorFile` (`:404`): replace `_ = os.MkdirAll` / `_ = os.WriteFile`
   discards with a concise stderr warning on failure
   (`warning: could not persist last-error log: <err>`) via the existing pipeline.
   The original error being handled is never obscured.
2. Audit `persistToErrorsDB` on the same path for `_ =` discards; fix any found the
   same way. Report the audit result (files/lines checked, discards found/fixed).
3. No signature changes; keep the diff small.

## Rules
GitMap tools only; no rg/grep. `go build ./...` only, never `go test`.

## Deliverable
Diff + audit report + build exit code. Then stop.
