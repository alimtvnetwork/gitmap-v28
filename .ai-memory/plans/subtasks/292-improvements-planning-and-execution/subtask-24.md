# Subtask Plan 24: Dedicated Agent Lifecycle Action Tracking (`cli/heatmap/heatmap_tracker_agent.go`)

> **Subtask ID:** `subtask-24`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/heatmap/heatmap_tracker_agent.go`  
> **Status:** Pending  

---

## 1. Objectives & Scope
1. Implement `TrackAgentAction(taskSlug, actionType, summary string) error` to record discrete agent lifecycle events into the root heatmap database.
2. Isolate the implementation into a new companion file `cli/heatmap/heatmap_tracker_agent.go` to prevent `cli/heatmap/heatmap_tracker.go` (currently 93 lines) from exceeding the 100-line cap.
3. Handle working directory resolution and database opening/closing defensively.
4. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
Autonomous agent subtasks execute exploration, reasoning, code editing, and verification. Currently, only file touch events (`TrackFileTouch` / `TrackBatchFileTouches`) are tracked in the heatmap tracker. High-level agent lifecycle events (such as plan initialization, subtask start, verification complete) lack a dedicated tracking helper.

Furthermore, `cli/heatmap/heatmap_tracker.go` is currently 93 lines. Adding agent action tracking functions directly to that file would immediately violate the 100-line limit.

Creating `cli/heatmap/heatmap_tracker_agent.go`:
1. Provides a high-level API for agents to record actions with reasoning notes.
2. Ensures `heatmap_tracker.go` remains strictly under 100 lines.
3. Automatically opens the root heatmap database, logs the action record, and closes the connection cleanly.

---

## 3. Concrete Implementation Details

### 3.1 New File `cli/heatmap/heatmap_tracker_agent.go`
```go
package heatmap

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TrackAgentAction records a discrete agent lifecycle event into the root heatmap database.
func TrackAgentAction(taskSlug, actionType, summary string) error {
	cleanSlug := strings.TrimSpace(taskSlug)
	cleanAction := strings.TrimSpace(actionType)
	if cleanSlug == "" || cleanAction == "" {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.cwd")
	}

	return recordAgentActionAtRoot(cwd, cleanSlug, cleanAction, summary)
}

func recordAgentActionAtRoot(cwd, taskSlug, actionType, summary string) error {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return err
	}

	defer hdb.Close()

	_, logErr := hdb.LogAction(store.HeatmapActionLogRecord{
		TaskSlug:       taskSlug,
		ActionType:     actionType,
		ReasoningNotes: summary,
		ReasoningState: "EXPLORATION",
	})

	return logErr
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** All errors (`os.Getwd()`, `OpenHeatmapRootDB`, `LogAction`) return directly to caller without exit handlers.
- **Rule 2 (Function Limit):** `TrackAgentAction` is 12 lines; `recordAgentActionAtRoot` is 14 lines. Both strictly satisfy <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Linear guard clauses with early returns on empty slug/action or directory error.
- **Rule 4 (Positive Booleans):** Positive compound guard `cleanSlug == "" || cleanAction == ""`.
- **Rule 6 (Vertical Whitespace):** Blank lines after `}` and before `return`.

---

## 5. Verification Criteria
- [ ] File `cli/heatmap/heatmap_tracker_agent.go` created in package `heatmap`.
- [ ] Declares `TrackAgentAction(taskSlug, actionType, summary string) error`.
- [ ] Gracefully ignores empty task slug or action type without error.
- [ ] Closes database connection via `defer hdb.Close()`.
- [ ] Returns structured error wrapped with `apperror.WrapSimple`.
- [ ] File is ~40 lines, well below the 100-line cap.
