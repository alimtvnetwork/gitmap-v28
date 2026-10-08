# Subtask 04: Global Dispatch Interceptor, AppError Integration & CLI Automation

> **Subtask ID:** Subtask-04  
> **Parent Plan:** `.ai-memory/plans/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design.md`  
> **Target Path:** `.ai-memory/plans/subtasks/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design/04-error-interceptor-and-cli-integration.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Wire the automated suggestion engine into global dispatch in `cli/cmd/root.go:runDispatch`.
2. Integrate structured suggestions directly into `appfault.AppError` and the global error handler `handleGlobalError`.
3. Eliminate ad-hoc and manual suggestion printing (`fmt.Printf("   Did you mean: ...")`) across all subcommands.
4. Integrate contextual remediation hints into `pull_remediation_hint.go` and `pipeline-ai`.

---

## 2. Detailed Technical Plan

### 2.1 Global CLI Interceptor (`cli/cmd/root.go`)
- In `runDispatch(args []string)`:
  - If token does not match any registered command:
    * Invoke `suggestion.DefaultEngine().ResolveCommand(token)`.
    * If suggestions are found, format and print via `suggestion.RenderBox(os.Stderr, group)`.
    * Return structured `appfault.NewNotFoundError` with suggestions attached.
  - Intercept unknown subcommands within command handlers using the same pattern.

### 2.2 AppError Extension (`cli/appfault/app_error.go`)
- Add field:
  ```go
  type AppError struct {
      Code        string                  `json:"code"`
      Message     string                  `json:"message"`
      StackTrace  string                  `json:"stackTrace,omitempty"`
      Suggestions []suggestion.Suggestion `json:"suggestions,omitempty"`
  }
  ```
- Add fluent method:
  ```go
  func (e *AppError) WithSuggestions(items ...suggestion.Suggestion) *AppError {
      e.Suggestions = append(e.Suggestions, items...)
      return e
  }
  ```
- In `handleGlobalError(err error)`:
  - If `appErr.Suggestions` is not empty, automatically call `suggestion.RenderBox` or `suggestion.RenderCompact` based on terminal context or JSON flag.

### 2.3 Elimination of Manual Suggestion Prints
- Search across `cli/` for manual suggestion strings:
  * `cli/cmd/rootusagecompact.go`
  * `cli/cmdpull/pull_remediation_hint.go`
  * `cli/cmdscan/`
- Replace manual formatting with structured calls to `engine.ResolveRemediation` and `AppError.WithSuggestions`.

---

## 3. Verification Criteria
- [ ] Mistyped commands (e.g. `gitmap scann`, `gitmap clne`) automatically trigger Catppuccin suggestion box.
- [ ] Zero manual `fmt.Printf("Did you mean")` statements remaining in command routing logic.
- [ ] `gitmap --json <invalid-command>` outputs structured suggestions JSON without crashing.
- [ ] End-to-end integration verified via unit tests.
