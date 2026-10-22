# Subtask 20: Isolate Interactive Candidate Selection Prompts to Stderr

## Status: Pending
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdpending/pending_commits_filter.go`

---

## 1. Objective

Enforce the standard error isolation invariant in `cli/cmdpending/pending_commits_filter.go` by redirecting interactive user selection and pagination prompts (`Press [Enter] for next, or 'q' to quit: `) strictly to `os.Stderr`. This guarantees standard output (`os.Stdout`) remains clean, deterministic, and uncontaminated when users or automated agents redirect output to files, pipes, or JSON parsers.

---

## 2. Scope & Implementation Details

1. **Standard Error Isolation Invariant**:
   - Any terminal prompt, question, pagination pause, or interactive candidate menu must write to `os.Stderr`.
   - Data payload emissions (JSON, structured lines, repo lists) must write to `os.Stdout`.

2. **Refactor `isDetailPromptTerminated(scanner *bufio.Scanner) bool`**:
   - Replace standard output write:
     ```go
     // Before
     fmt.Printf("  Press [Enter] for next, or 'q' to quit: ")
     
     // After
     fmt.Fprintf(os.Stderr, "  Press [Enter] for next, or 'q' to quit: ")
     ```

3. **Refactor `renderPendingCommitsInteractive(records []RepoPendingCommitRecord)`**:
   - Route pagination status message `All repositories are clean. Nothing to step through.` and interactive stepping progress `[idx/total]` to `os.Stderr` when interactive stepping is active.
   - Decompose loop body into a sub-15-line helper: `renderInteractiveRecordStep(rec RepoPendingCommitRecord, idx, total int)`.

4. **Return Blank Line Styling**:
   - Fix missing blank lines before return statements across `pending_commits_filter.go`.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Functions must be $\le 15$ lines of executable code.
- **Zero Nested Ifs**: Invert conditions with guard clauses and early returns.
- **Positive Booleans**: Use `isDirtyOnly`, `isClean`.
- **Stderr Isolation**: Interactive prompts must never write to `os.Stdout`.
- **Zero Builds & Tests**: Strict prohibition on compiler builds or tests.

---

## 4. Verification Checklist

- [ ] Interactive prompts in `isDetailPromptTerminated` write strictly to `os.Stderr`.
- [ ] Interactive step status messages write to `os.Stderr`.
- [ ] Standard output remains pure and unpolluted when piping or redirecting.
- [ ] Blank lines precede every return statement.
- [ ] All functions conform to $\le 15$-line limit.
