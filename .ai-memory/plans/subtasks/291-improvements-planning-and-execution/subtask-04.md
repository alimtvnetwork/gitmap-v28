# Subtask 04: Help Topic Leaf Error Sovereignty & Topic Validator Extraction

> **Subtask ID:** Subtask-04  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-04.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Eliminate the hard process termination call `cliexit.Exit(0)` in `cli/cmdhelp/help_topic.go`.
2. Ensure `RunHelpTopic(args []string)` returns `*appfault.AppError` cleanly to the outer command dispatcher without taking exit decisions directly.
3. Architect `cli/cmdhelp/helpcheck.go` to house reusable topic lookup, existence verification, and arg validation helpers, preventing void error discards across help subcommands.
4. Enforce strict coding standards:
   - Functions $\le 15$ lines ($\le 8$ lines preferred).
   - Zero nested conditionals.
   - Positive boolean naming (`isHelpRequested`, `hasTopicMatch`).
   - Blank line before every `return` and after every closing brace `}`.

---

## 2. Detailed Technical Plan

### 2.1 Eliminate `cliexit.Exit(0)` in `cli/cmdhelp/help_topic.go`
- Locate `RunHelpTopic(args []string) *appfault.AppError`:
  - When `isHelpRequested(args)` is true:
    - Print `helpTopicHelpText`.
    - Replace `cliexit.Exit(0)` with `return nil`.
    - This allows root dispatcher in `cli/cmd/root.go` / `rootutility_help.go` to conclude cleanly and log audit metadata.

### 2.2 Architect `cli/cmdhelp/helpcheck.go`
- Create `cli/cmdhelp/helpcheck.go` containing:
  - `findHelpTopic(name string) (helpTopic, bool)`: searches `helpTopics` slice and returns match status using positive boolean `hasFound` ($\le 12$ lines).
  - `validateTopicArgs(args []string) (string, *appfault.AppError)`: checks argument length, normalizes topic string, and returns an `appfault.NewValidationError` if empty ($\le 10$ lines).
  - Refactor `RunHelpTopic` to leverage `validateTopicArgs` and `findHelpTopic`, reducing its cyclomatic complexity.

### 2.3 Style & Guideline Compliance
- Ensure every `return` is preceded by an empty line.
- Verify all boolean parameters and return flags are positive (`isHelpRequested`, `hasFound`).
- Flatten all conditionals using guard clauses and early returns.

---

## 3. Verification Criteria

- [ ] Zero calls to `cliexit.Exit` or `os.Exit` exist in `cli/cmdhelp/help_topic.go`.
- [ ] Help flag requests (`-h`, `--help`) render usage and return `nil` cleanly.
- [ ] `cli/cmdhelp/helpcheck.go` is authored with discrete topic lookup routines $\le 15$ lines.
- [ ] Unknown topic requests return a strongly-typed `*appfault.AppError`.
- [ ] Blank line formatting is strictly preserved across all functions in `cli/cmdhelp/`.
