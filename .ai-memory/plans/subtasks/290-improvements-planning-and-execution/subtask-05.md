# Subtask Plan 05: Check & Propagate Clipboard and Stream Write Errors (`cli/cmdpipeline/pipeline_history.go`)

- **Subtask ID:** `subtask-05`
- **Subtask Name:** `05-pipeline-history-clipboard-error-propagation`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdpipeline/pipeline_history.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdpipeline/pipeline_history.go`, clipboard and terminal output operations for positional pipeline commit inspection violate **Rule 1 (Error Return Sovereignty)**:

1. In `copyPositionalReportToClipboard(terminalOutput string, isFailure bool)`:
   ```go
   cleanOutput := stripANSI(terminalOutput)
   _ = writeClipboard(cleanOutput)
   label := "positional pipeline report"
   if isFailure {
       label = "positional pipeline error report"
   }
   fmt.Printf("\n  %s📋 Copied %s to clipboard%s\n", constants.ColorCyan, label, constants.ColorReset)
   ```
   The return value of `writeClipboard(cleanOutput)` is blanked out (`_ = writeClipboard(...)`). Even if the platform clipboard tool (e.g. `pbcopy`, `xclip`, `wl-copy`, Windows PowerShell clipboard) fails or is unavailable in headless CI/CD, the function unconditionally prints that the report was copied to clipboard!
2. `renderPositionalCommitTerminal` is a void function that calls `copyPositionalReportToClipboard`. In `dispatchInspectGroupOutput`, which declares an `error` return type, it calls `renderPositionalCommitTerminal` and returns `nil`.
3. To adhere to Error Return Sovereignty:
   - `copyPositionalReportToClipboard` must return `error`. When `writeClipboard` fails, it must emit a warning notice or return the error, avoiding false success announcements.
   - `renderPositionalCommitTerminal` must return `error` and propagate it to `dispatchInspectGroupOutput`.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `copyPositionalReportToClipboard`
- **File:** `cli/cmdpipeline/pipeline_history.go`
- **Signature:** `func copyPositionalReportToClipboard(terminalOutput string, isFailure bool) error`
- **Requirements:**
  - Check error from `writeClipboard(cleanOutput)`.
  - If error is non-nil, emit warning or return wrapped error without false success message.
  - Return `nil` on success.
  - Function length <= 15 lines.

```go
func copyPositionalReportToClipboard(terminalOutput string, isFailure bool) error {
	cleanOutput := stripANSI(terminalOutput)
	if err := writeClipboard(cleanOutput); err != nil {
		fmt.Printf("\n  %s⚠ Could not copy to clipboard: %v%s\n", constants.ColorYellow, err, constants.ColorReset)

		return err
	}

	label := "positional pipeline report"
	if isFailure {
		label = "positional pipeline error report"
	}
	fmt.Printf("\n  %s📋 Copied %s to clipboard%s\n", constants.ColorCyan, label, constants.ColorReset)

	return nil
}
```

### 2.2 Refactor `renderPositionalCommitTerminal`
- **File:** `cli/cmdpipeline/pipeline_history.go`
- **Signature:** `func renderPositionalCommitTerminal(repo string, runs []ghRunItem, groups []CommitPipelineGroup, group *CommitPipelineGroup, offset int) error`
- **Requirements:**
  - Construct inspector content and recent summary into `strings.Builder`.
  - Print collapsed output to stdout.
  - Return result of `copyPositionalReportToClipboard(output, isCommitGroupFailure(group))`.
  - Function length <= 12 lines.

```go
func renderPositionalCommitTerminal(repo string, runs []ghRunItem, groups []CommitPipelineGroup, group *CommitPipelineGroup, offset int) error {
	var sb strings.Builder
	renderCommitInspectorContent(&sb, repo, runs, group, offset)
	renderRecentCommitsSummaryToBuilder(&sb, groups, 5)
	output := CollapseConsecutiveEmptyLines(sb.String())
	fmt.Print(output)

	return copyPositionalReportToClipboard(output, isCommitGroupFailure(group))
}
```

### 2.3 Refactor `dispatchInspectGroupOutput`
- **File:** `cli/cmdpipeline/pipeline_history.go`
- **Signature:** `func dispatchInspectGroupOutput(repo string, runs []ghRunItem, groups []CommitPipelineGroup, group *CommitPipelineGroup, offset int, isJSON bool) error`
- **Requirements:**
  - If `isJSON`, return `printPositionalJSON(group)`.
  - Otherwise, return `renderPositionalCommitTerminal(repo, runs, groups, group, offset)`.
  - Function length <= 8 lines.

```go
func dispatchInspectGroupOutput(repo string, runs []ghRunItem, groups []CommitPipelineGroup, group *CommitPipelineGroup, offset int, isJSON bool) error {
	if isJSON {
		return printPositionalJSON(group)
	}

	return renderPositionalCommitTerminal(repo, runs, groups, group, offset)
}
```

---

## 3. Implementation Verification Checklist

- [x] `_ = writeClipboard(cleanOutput)` removed.
- [x] `copyPositionalReportToClipboard` returns `error` and guards against false success reporting.
- [x] `renderPositionalCommitTerminal` returns `error`.
- [x] `dispatchInspectGroupOutput` propagates error from `renderPositionalCommitTerminal`.
- [x] Function lengths <= 15 lines.
- [x] Zero nested `if` statements.
- [x] Blank line before `return` and after `}` preserved.
