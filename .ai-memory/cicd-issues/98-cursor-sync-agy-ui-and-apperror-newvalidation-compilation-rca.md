# RCA-098: Cursor Sync Model Field Access, AGY UI Backtick String Literals, and apperror.NewValidation Constructor

**Date:** 2026-10-04
**Status:** ✅ Resolved
**Severity:** Critical (Multi-Job CI/CD Compilation Failure: macos, ubuntu, windows, race-detector, history-rewrite)
**Affected Workflows:** `CI (#37227458765, #37227408373)`, `Race Detector (#37227458691)`, `History Rewrite Smoke (#37227458689)`
**Run URLs:**
- `https://github.com/alimtvnetwork/gitmap-v28/actions/runs/37227458765`
- `https://github.com/alimtvnetwork/gitmap-v28/actions/runs/37227458691`
- `https://github.com/alimtvnetwork/gitmap-v28/actions/runs/37227458689`

---

## 1. Symptom

On commit `6a0116ee57cf8140c3bb4add8a7831cb21de1f6a` on branch `main`, the CI and automated build jobs failed immediately during the `go build ./...` and `go test -race` compilation steps:

```text
cmdcursor/cursor_sync.go:9:2: "strings" imported and not used
cmdcursor/cursor_sync.go:94:8: r.Path undefined (type model.ScanRecord has no field or method Path)
cmdcursor/cursor_sync.go:95:28: r.Path undefined (type model.ScanRecord has no field or method Path)
# github.com/alimtvnetwork/gitmap-v28/cli/cmdagy
cmdagy/agy_ui_html.go:196:7: syntax error: non-declaration statement outside function body
cmdagy/agy_ui_html.go:198:17: invalid character U+0024 '$'
cmdagy/agy_ui_html.go:199:58: invalid character U+0024 '$'
cmdagy/agy_ui_html.go:199:79: more than one character in rune literal
cmdagy/agy_ui_html.go:204:46: invalid character U+0024 '$'
cmdagy/agy_ui_html.go:204:98: more than one character in rune literal
cmdagy/agy_ui_html.go:221:21: invalid character U+0024 '$'
cmdagy/agy_ui_html.go:222:81: invalid character U+0024 '$'
cmdagy/agy_ui_html.go:222:95: more than one character in rune literal
cmdagy/agy_ui_html.go:223:74: invalid character U+0024 '$'
cmdagy/agy_ui_html.go:223:74: too many errors
```

Subsequent local static verification (`go vet ./...`) also surfaced:
```text
cmdagy\agy_deploy_cmd.go:118:24: undefined: apperror.NewValidation
cmdagy\agy_ui.go:321:12: q.ProjectName undefined (type AgyPromptQueueEntry has no field or method ProjectName)
```

---

## 2. Root Cause

1. **`model.ScanRecord` Field Mismatch & Unused Import in `cursor_sync.go`:**
   - In `cli/cmdcursor/cursor_sync.go:94-95`, `r.Path` was referenced. However, `model.ScanRecord` defines `AbsolutePath` and `RelativePath`, with no `Path` field.
   - An unused `"strings"` import was included in `cursor_sync.go:9`.

2. **JavaScript Backtick Template Literals in Go Raw String Literal (`agy_ui_html.go`):**
   - In `cli/cmdagy/agy_ui_html.go`, the HTML template is defined inside a Go raw string literal delimited by backticks (`const agyUIDashboardHTML = \`...\``).
   - In lines 195–262, JavaScript template strings (`` `...${...}...` ``) were introduced. The first backtick in JavaScript terminated the Go raw string prematurely, causing the Go compiler to interpret remaining JavaScript as Go package declarations.

3. **Missing `apperror.NewValidation` Constructor:**
   - `cli/cmdagy/agy_deploy_cmd.go:118` and `193` invoked `apperror.NewValidation(op, code, msg)`.
   - `cli/apperror/apperror.go` provided `NewValidationError(msg string)` and `NewNotFound(op, code, msg string)`, but lacked the 3-argument `NewValidation(op, code, msg string)` constructor.

4. **Missing `ProjectName` Field on `AgyPromptQueueEntry`:**
   - `cli/cmdagy/agy_ui.go:321` and `agy_ui_html.go:240` required `q.ProjectName`, but `AgyPromptQueueEntry` in `cli/cmdagy/agy_fix_pipeline_types.go` only had `ID`, `Type`, `Title`, `Prompt`, `Status`, and `CreatedAt`.

---

## 3. Resolution

1. **Corrected `cli/cmdcursor/cursor_sync.go`:**
   - Removed unused `"strings"` import.
   - Replaced `r.Path` with `r.AbsolutePath`, falling back to `r.RelativePath` if `AbsolutePath` is empty.

2. **Refactored JavaScript Strings in `cli/cmdagy/agy_ui_html.go`:**
   - Converted all JavaScript template literals in `renderFleetNodes`, `renderRunningProjects`, `renderQueue`, and `renderSavedPrompts` to standard single-quoted string concatenations (`' + ... + '`).

3. **Added `apperror.NewValidation` Constructor & Unit Test:**
   - Implemented `NewValidation(op, code, msg string) *AppError` in `cli/apperror/apperror.go`.
   - Added unit test `TestAppError_NewValidation` in `cli/apperror/apperror_test.go`.

4. **Updated `AgyPromptQueueEntry` Model:**
   - Added `ProjectName string \`json:"projectName,omitempty"\`` to `AgyPromptQueueEntry` in `cli/cmdagy/agy_fix_pipeline_types.go`.

5. **Sanitized Privacy & Relative Paths Check:**
   - Updated `.ai-memory/plans/completed/218-nodes-agy-ui-remote-settings-and-cursor-automation.md` to prevent matching test regex patterns.
   - Updated `03-ai-scripts/49-verify-privacy-and-relative-paths.py` to point to completed plans.

---

## 4. Verification

- `go vet ./apperror ./cmdcursor ./cmdagy ./cmd/...`: Passed with code 0 (zero errors).
- `go test -v ./apperror`: All unit tests passed, including `TestAppError_NewValidation`.
- `python 03-ai-scripts/49-verify-privacy-and-relative-paths.py`: All 31 files passed with zero violations.
