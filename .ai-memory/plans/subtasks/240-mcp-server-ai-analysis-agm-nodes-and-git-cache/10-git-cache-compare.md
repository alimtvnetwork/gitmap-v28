# Subtask 10: Git Commit Graph SQLite Cache & Instant Branch Compare Subsystem

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `10-git-cache-compare`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-component-and-cli-spec.md` (Section 6)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Implement the **SQLite Split-DB Git Commit Cache Engine**, cached commit logs (`gitmap log`), and sub-millisecond branch ahead/behind diff comparison (`gitmap branch compare <branchA> <branchB>`, aliases `gitmap branch cmp`, `gitmap cmp`).

### The Problem
In multi-branch and multi-agent development workflows, calculating branch divergences and displaying git history via subshell commands (`git rev-list`, `git log`) incurs 50–150ms of process startup overhead per invocation. Repeated calls during pipeline monitoring or interactive rebasing waste compute cycles and degrade agent speed.

### The Solution
Implement an isolated Split SQLite commit graph database in `.gitmap/data/git-cache/<slug>/sql.db`:
1. **0ms Fast Path:** Reads the current head commit hashes directly from `.git/refs/heads/<branch>` in <1ms. Checks composite cache key `<branchA>...<branchB>:<hashA>:<hashB>`. If cached, returns ahead/behind counts and commit lists in **<1ms** with zero git process spawns.
2. **On-Demand Graph Synchronization:** On cache miss, executes a single rev-list query, writes missing commit objects into `GitCommit`, and stores divergence results in `GitBranchCompareCache`.
3. **Cached Log Viewer (`gitmap log`):** Renders recent commit history from local SQLite tables with pagination, author filtering, and `--json` support.
4. **Terminal UI Rendering:** Formats branch ahead/behind divergence in an aligned summary card using `cli/termtable` and `cli/termpad`.

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/cmdgitcache/gitcache_types.go` | New | Data models: `GitCommitRecord`, `BranchCompareResult`, and query option structs. |
| `cli/cmdgitcache/gitcache_db.go` | New | SQLite connection manager, DDL migration, and CRUD queries for `sql.db`. |
| `cli/cmdgitcache/gitcache_log.go` | New | Cached git log queries, formatting, and `RunCachedLogCommand()`. |
| `cli/cmdgitcache/gitcache_compare.go` | New | Ahead/behind calculation engine, 0ms fast path lookup, and `termtable` card rendering. |
| `cli/cmdgitcache/gitcache_test.go` | New | Unit tests for cache hits/misses, divergence math, and SQLite storage integrity. |
| `cli/cmd/rootdispatch.go` | Modify | Register `log`, `branch compare`, `branch cmp`, and `cmp` routing. |

---

## 3. Detailed Implementation Requirements

### 3.1 SQLite Schema (`.gitmap/data/git-cache/<slug>/sql.db`)
```sql
CREATE TABLE IF NOT EXISTS GitCommit (
    CommitHash TEXT PRIMARY KEY,
    AuthorName TEXT NOT NULL,
    AuthorEmail TEXT NOT NULL,
    AuthorDate DATETIME NOT NULL,
    CommitterName TEXT NOT NULL,
    CommitterDate DATETIME NOT NULL,
    CommitMessage TEXT NOT NULL,
    ParentHashes TEXT,
    TreeHash TEXT NOT NULL,
    IsMergeCommit INTEGER NOT NULL DEFAULT 0,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS GitBranchCompareCache (
    CompareKey TEXT PRIMARY KEY,
    BaseBranch TEXT NOT NULL,
    TargetBranch TEXT NOT NULL,
    AheadCount INTEGER NOT NULL,
    BehindCount INTEGER NOT NULL,
    AheadHashes TEXT NOT NULL,
    BehindHashes TEXT NOT NULL,
    CommonMergeBase TEXT NOT NULL,
    CachedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 3.2 0ms Fast Path Engine (`cli/cmdgitcache/gitcache_compare.go`)
```go
package cmdgitcache

import (
	"fmt"
	"time"
)

func CompareBranchesFast(slug, baseBranch, targetBranch string, isForceRefresh bool) (*BranchCompareResult, error) {
	start := time.Now()
	baseHash, err := readRefHash(baseBranch)
	if err != nil {
		return nil, err
	}
	targetHash, err := readRefHash(targetBranch)
	if err != nil {
		return nil, err
	}

	compareKey := fmt.Sprintf("%s...%s:%s:%s", baseBranch, targetBranch, baseHash, targetHash)
	db, err := OpenGitCacheDB(slug)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if !isForceRefresh {
		if cached, isHit := lookupCompareCache(db, compareKey); isHit {
			cached.IsCacheHit = true
			cached.CalculationMs = time.Since(start).Milliseconds()
			return cached, nil
		}
	}

	// Compute via git engine and populate SQLite cache
	return computeAndCacheBranchDivergence(db, slug, baseBranch, targetBranch, baseHash, targetHash, compareKey, start)
}
```

---

## 4. Acceptance Criteria

- [ ] SQLite database `.gitmap/data/git-cache/<slug>/sql.db` initializes with correct schema.
- [ ] Direct ref hash reading (`.git/refs/heads/...`) detects head state in <1ms.
- [ ] Cache hit returns branch comparison in <1ms (0ms fast path).
- [ ] Diverged branch commits are categorized into ahead vs behind accurately.
- [ ] Terminal UI displays styled box using `termtable` and `termpad`.
- [ ] `gitmap log` supports `--cached`, `--max`, `--author`, and `--json`.
- [ ] Unit tests pass with >85% coverage.

---

## 5. Verification Commands

```powershell
# 1. Run unit tests
go test -v ./cli/cmdgitcache/...

# 2. Test branch comparison between main and HEAD
gitmap cmp main HEAD

# 3. Test cached log output
gitmap log -n 5

# 4. Verify coding guidelines
python 03-ai-scripts/05-guideline-autofixer.py cli/cmdgitcache --check-only
```
