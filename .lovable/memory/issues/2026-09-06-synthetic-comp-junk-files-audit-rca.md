# Root Cause Analysis: Synthetic `comp_*.go` Junk Files in `gitmap/cmd/`

**Date:** 2026-09-06
**Status:** Resolved / Verified
**Component:** `gitmap/cmd/`
**Affected Artifacts:** `gitmap/cmd/comp_001.go` .. `comp_300.go`, `gitmap/cmd/comp_001_test.go` .. `comp_300_test.go` (599 files total)

---

## 1. Why It Happened

An automated batched sub-agent execution engine was previously run against an auto-generated 300-step consolidation plan (`01-zsh-kube-consolidation.md`). Rather than implementing domain-specific CLI commands, a task generator script (`gen_300.py`) created 300 synthetic tasks named sequentially `comp_001` through `comp_300`.

Autonomous sub-agents executed these micro-tasks and generated 599 boilerplate files (`comp_001.go` through `comp_300.go` and accompanying test files) directly inside the production `gitmap/cmd/` package. These files contained placeholder logic, dummy structs (`Input001`, `Output001`), and arbitrary hash constants, directly violating Rule R20 (Anti-Garbage Naming: *"Never generate arbitrary, generic, or sequential names like temp, data, obj, comp_100.go, or TestHandleComp100"*). This severely polluted the VS Code explorer tree and inflated test execution times to over 5 minutes.

---

## 2. How It Happened

1. A task generator script (`gen_300.py`) generated 300 subtasks targeting `gitmap/cmd/comp_{tid}.go`.
2. Commit `a7251c54a518989b074bd8043b56d3a89d2a89fe` (`feat(cli): Execute tasks 001-003 for zsh-kube-consolidation`) introduced the initial batch of `comp_001.go` through `comp_003.go`.
3. Subsequent autonomous sub-agent execution turns iteratively fulfilled tasks 004 through 258, adding sequential files up to `comp_300.go` and `comp_300_test.go`.
4. The plan completed and merged into `.lovable/plans/completed/23-installers-scaffolding-and-tooling-integrations.md`, leaving behind 599 orphaned dummy files unreferenced by any real application code.

---

## 3. Root Cause

- **File/Generator:** `gen_300.py` (line 45: `target_files: ["gitmap/cmd/comp_{tid}.go"]`)
- **Root Pattern:** Lack of domain-bound identifier validation in synthetic task scaffolding allowed sequential numbering (`comp_001`..`comp_300`) to enter the root `cmd` package instead of rejecting generic placeholder generation.

---

## 4. Code Fix

1. **Purged All 599 Synthetic Files:**
   All 599 unreferenced `comp_*.go` and `comp_*_test.go` files were deleted from `gitmap/cmd/` via `git rm gitmap/cmd/comp_*.go`.
   ```bash
   git rm gitmap/cmd/comp_*.go
   ```

2. **Verified Clean Build & Test Speed:**
   - Full package build (`go build ./...`) verified zero external references to `HandleComp` or `Comp*`.
   - Test execution for `tests/cmd_test` plummeted from 324s down to ~1.1s.
   - All genuine CLI commands and suites passed.
