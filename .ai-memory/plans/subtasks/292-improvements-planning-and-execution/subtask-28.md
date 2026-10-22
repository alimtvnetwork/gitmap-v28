# Subtask Plan 28: Target Discovery & Dirty Extraction (`cli/cmdcpar/cpar_target.go`)

> **Subtask ID:** `subtask-28`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdcpar/cpar_target.go`, `cli/cmdcpar/cpar.go`  
> **Status:** Pending  

---

## 1. Objectives & Scope
1. Extract repository target resolution and dirty state discovery from `cli/cmdcpar/cpar.go` into a new companion file `cli/cmdcpar/cpar_target.go`.
2. Move `CollectDirtyRepositories`, `resolveRelativePath`, `appendDirty`, `resolveAllRecords`, and `queryStoreRecords` into `cpar_target.go`.
3. Reduce `cli/cmdcpar/cpar.go` from 292 lines to ~230 lines, keeping both files comfortably below sizing limits.
4. Maintain full package-internal visibility so existing callers in `cli/cmdcpar` continue working without modification.
5. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
`cli/cmdcpar/cpar.go` currently mixes multiple responsibilities:
- Command argument parsing and option resolution.
- Target discovery via SQLite queries and fallback resolution.
- Working tree dirty state inspection.
- Interactive user consent prompting and table rendering.
- Batch git execution (`git add`, `git commit`, `git push`).
- Background TaskQueue telemetry persistence.

Extracting target discovery into `cpar_target.go`:
1. Decouples repository discovery from batch command execution.
2. Isolates database interaction and path normalization in a focused module.
3. Decreases the line count of `cpar.go` significantly, advancing the repository-wide decomposition goal.

---

## 3. Concrete Implementation Details

### 3.1 New File `cli/cmdcpar/cpar_target.go`
Create `cli/cmdcpar/cpar_target.go` in package `cmdcpar`:

```go
package cmdcpar

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CollectDirtyRepositories scans all tracked repositories and returns those with dirty states.
func CollectDirtyRepositories() []DirtyRepoSummary {
	records := resolveAllRecords()
	var dirty []DirtyRepoSummary
	for _, r := range records {
		relPath := resolveRelativePath(r)
		diag := gitutil.InspectDirtyState(relPath)
		if diag.IsDirty {
			dirty = appendDirty(dirty, r, relPath, diag)
		}
	}

	return dirty
}

func resolveRelativePath(r model.ScanRecord) string {
	hasRel := len(r.RelativePath) > 0
	if hasRel {
		return r.RelativePath
	}

	hasName := len(r.RepoName) > 0
	if hasName {
		return r.RepoName
	}

	return "."
}

func appendDirty(dirty []DirtyRepoSummary, r model.ScanRecord, relPath string, diag gitutil.DirtyDiagnosis) []DirtyRepoSummary {
	return append(dirty, DirtyRepoSummary{
		RepoName:  r.RepoName,
		RepoPath:  relPath,
		Diagnosis: diag,
	})
}

func resolveAllRecords() []model.ScanRecord {
	records := queryStoreRecords()
	hasRecords := len(records) > 0
	if hasRecords {
		return records
	}

	return []model.ScanRecord{{RelativePath: ".", RepoName: "."}}
}

func queryStoreRecords() []model.ScanRecord {
	s, err := store.OpenDefault()
	hasErr := err != nil
	if hasErr {
		return nil
	}

	defer s.Close()

	records, _ := s.ListRepos()

	return records
}
```

### 3.2 Update `cli/cmdcpar/cpar.go`
Remove the extracted functions (`CollectDirtyRepositories`, `resolveRelativePath`, `appendDirty`, `resolveAllRecords`, `queryStoreRecords`) from `cli/cmdcpar/cpar.go`.

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** No exit handlers. Handled cleanly with empty or fallback slices.
- **Rule 2 (Function Limit):** Every extracted function is between 6 and 14 lines, well within <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Linear guard clauses (`hasRel`, `hasName`, `hasRecords`, `hasErr`).
- **Rule 4 (Positive Booleans):** All guards use positive naming (`hasRel`, `hasName`, `hasRecords`, `diag.IsDirty`).
- **Rule 6 (Vertical Whitespace):** Proper blank lines around control blocks and before returns.

---

## 5. Verification Criteria
- [ ] File `cli/cmdcpar/cpar_target.go` created in package `cmdcpar`.
- [ ] Functions `CollectDirtyRepositories`, `resolveRelativePath`, `appendDirty`, `resolveAllRecords`, and `queryStoreRecords` present in `cpar_target.go`.
- [ ] Functions removed from `cli/cmdcpar/cpar.go`.
- [ ] `cpar.go` line count decreased by ~60 lines.
- [ ] Zero changes to external signatures or public behavior.
