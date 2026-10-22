# Subtask 17: Add TokyoTag Formatting to HasChange Multi-Dimension Output

## Status: Completed
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdhaschange/haschange.go`

---

## 1. Objective

Enhance the multi-dimension output (`--all`) in `cli/cmdhaschange/haschange.go` to display Tokyo Night styled badges using `termout.TokyoTag` when running in an interactive terminal, while strictly preserving pure stdout output (`true` / `false`) in single-dimension mode for shell scripting integration.

---

## 2. Scope & Implementation Details

1. **Single-Mode Output Purity Invariant**:
   - `printHasChangeOne(target, mode string)` MUST NOT be modified in a way that introduces ANSI codes, badges, or formatting wrappers.
   - Piped and non-piped execution of single mode must strictly output raw `"true"` or `"false"`.

2. **Modernize `printHasChangeAll(target string)`**:
   - Query repository state: `dirty := isWorkingTreeDirty(target)` and `ahead, behind, hasUpstream := readAheadBehind(target)`.
   - When no upstream is configured (`!hasUpstream`):
     - If `termout.ColorAllowed()`, render:
       `fmt.Println(termout.TokyoTag("dirty", boolStr(dirty)) + " " + termout.TokyoTag("upstream", "missing"))`
     - Otherwise, print standard fallback: `fmt.Printf(constants.MsgHCAllNoUpstream, boolStr(dirty))`.
   - When upstream exists:
     - If `termout.ColorAllowed()`, render:
       `fmt.Println(formatHasChangeStyledTags(dirty, ahead, behind))`
     - Otherwise, print standard plain string: `fmt.Printf(constants.MsgHCAllFmt, boolStr(dirty), boolStr(ahead > 0), boolStr(behind > 0))`.

3. **Helper Decomposition**:
   - `formatHasChangeStyledTags(dirty bool, ahead, behind int) string`:
     - Assembles:
       - `termout.TokyoTag("dirty", boolStr(dirty))`
       - `termout.TokyoTag("ahead", boolStr(ahead > 0))`
       - `termout.TokyoTag("behind", boolStr(behind > 0))`
     - Joins tags with a single space and returns.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Functions must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: All logic must use guard clauses and early returns.
- **Positive Booleans**: Strict use of positive booleans (`hasUpstream`, `isDirty`).
- **Standard Output Purity**: Single mode stdout remains strictly parseable by scripts.
- **Zero Builds & Tests**: Strict prohibition on compiler executions.

---

## 4. Verification Checklist

- [x] `gitmap hc <repo> --all` formats output with `termout.TokyoTag` in TTY / color mode.
- [x] Non-TTY / piped output for `--all` preserves plain `dirty=true ahead=false behind=true`.
- [x] Single mode `gitmap hc <repo>` retains pure `"true"` / `"false"` stdout.
- [x] Decomposed helper functions conform to $\le 15$-line rule.
- [x] Blank line before `return` strictly preserved.
