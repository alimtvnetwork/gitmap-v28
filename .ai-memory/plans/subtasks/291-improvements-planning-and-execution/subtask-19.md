# Subtask 19: Wire CD Selection Table in cdops.go

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdcd/cdops.go`

---

## 1. Objective

Integrate the dedicated selection table from `cd_table.go` into `cli/cmdcd/cdops.go`, replacing unformatted loops in `promptCDPick` and `promptCDReposPick` with formatted rounded tables on `os.Stderr`.

---

## 2. Scope & Implementation Details

1. **Update `promptCDPick`**:
   - Guard check: If `len(records) == 0`, return `"", fmt.Errorf(constants.ErrCDNotFound, name)`.
   - Call `RenderCDSelectionTable(fmt.Sprintf("Multiple repositories matching '%s':", name), records)`.
   - Prompt user on `os.Stderr`: `fmt.Fprintf(os.Stderr, "\nPick [1-%d] (default 1): ", len(records))`.
   - Call `readCDSelection(records, name)`.

2. **Update `promptCDReposPick`**:
   - Guard check: If `len(records) == 0`, return `"", fmt.Errorf("no repositories found")`.
   - Call `RenderCDSelectionTable("Available Repositories:", records)`.
   - Prompt user on `os.Stderr`: `fmt.Fprintf(os.Stderr, constants.MsgCDPickPrompt, len(records))`.
   - Call `readCDSelection(records)`.

3. **Stream Isolation Invariant**:
   - Verify that neither `promptCDPick` nor `promptCDReposPick` writes to `os.Stdout`.
   - Only the final resolved path in `runCD` is printed to `os.Stdout`.

4. **Decomposition**:
   - Keep `promptCDPick` and `promptCDReposPick` $\le 15$ lines of executable code.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Guard clauses and early returns.
- **Error Sovereignty**: Bubble errors up to root command handler.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `promptCDPick` calls `RenderCDSelectionTable`.
- [x] `promptCDReposPick` calls `RenderCDSelectionTable`.
- [x] Prompts and tables output strictly to `os.Stderr`.
- [x] `os.Stdout` remains completely clean for path capture.
- [x] All modified functions $\le 15$ lines.
