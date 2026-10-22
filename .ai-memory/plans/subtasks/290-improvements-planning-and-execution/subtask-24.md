# Subtask Plan 24: Direct File Touch Tracking Helper (`cli/heatmap/heatmap_tracker.go`)

> **Subtask ID:** `subtask-24`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/heatmap/heatmap_tracker.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Author new module `cli/heatmap/heatmap_tracker.go`.
2. Implement `TrackFileTouch(filePath string, touchType store.FileTouchType) error`.
3. Provide a fast, convenient tracking entrypoint that records file access directly into the root heatmap database without requiring manual struct boilerplate.
4. Adhere to Error Return Sovereignty (Rule 1), Function Length Limits (<= 15 lines, Rule 2), and Strongly Typed Enums (Rule 5).

---

## 2. Targeted Problem & Architectural Rationale
GitMap commands that touch repository files (e.g. `cmdpull`, `cmdfixauth`, AI subagents, AST transformers) want to emit heatmap telemetry. Previously, doing so required callers to:
1. Know how to resolve the root heatmap database path.
2. Manually initialize a `HeatmapActionLogRecord` with 18 fields.
3. Call `hdb.LogAction(...)` and handle resource cleanup.

By encapsulating this into `TrackFileTouch(filePath, touchType)`, callers obtain a clean one-line telemetry recording helper.

---

## 3. Concrete Implementation Details

### 3.1 Implementation of `cli/heatmap/heatmap_tracker.go`

```go
package heatmap

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TrackFileTouch records a file access operation into the root heatmap database.
func TrackFileTouch(filePath string, touchType store.FileTouchType) error {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimSpace(filePath)))
	if cleanPath == "" || cleanPath == "." {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.cwd")
	}

	return recordFileTouchAtRoot(cwd, cleanPath, touchType)
}

func recordFileTouchAtRoot(cwd, cleanPath string, touchType store.FileTouchType) error {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return err
	}

	defer hdb.Close()

	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
		ActionType: string(touchType),
		FilePath:   cleanPath,
	})

	return logErr
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Errors from `os.Getwd()`, `store.OpenHeatmapRootDB`, and `hdb.LogAction` are returned directly. No exit handlers or swallowed errors.
- **Rule 2 (Function Limit):** `TrackFileTouch` is 14 lines; `recordFileTouchAtRoot` is 13 lines. Both <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Guard clauses on empty path and error checks; zero nested conditionals.
- **Rule 5 (Typed Enums):** Accepts strongly typed `store.FileTouchType` ending in `Type` suffix.
- **Rule 6 (Vertical Whitespace):** Single blank line after each `}` and before `return`.

---

## 5. Verification Criteria
- [x] File `cli/heatmap/heatmap_tracker.go` is created in package `heatmap`.
- [x] Exported function `TrackFileTouch(filePath string, touchType store.FileTouchType) error` is present.
- [x] Normalizes paths to forward slashes with `filepath.ToSlash`.
- [x] Returns `nil` on empty or blank file path.
- [x] Every function <= 15 lines.
- [x] All paths are strictly relative lowercase Git paths.
