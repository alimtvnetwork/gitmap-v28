# Consolidated Commits, Immediate Push, and Build Execution Standards

**Added:** 2026-09-12
**Category:** Workflow & Git Hygiene
**Status:** Active Rule (Mandatory)

---

## 1. Context & Motivation

During autonomous guideline runs, committing isolated 1-2 markdown plan files without their underlying code changes pollutes the Git commit history with fragmented entries. Furthermore, leaving commits unpushed across conversational turns creates drift between the local repository and the remote origin, causing confusion in GitHub Desktop and CI dashboards. Finally, running heavy full builds (`npm run build`, `go build ./...`) or the entire 38-gate CI runner during routine styling or naming audits causes extreme latency and resource waste.

---

## 2. Mandatory Rules

### Rule 1: Consolidated Atomic Commits
- **NEVER** commit isolated 1-2 plan/doc markdown files alone.
- Group all modified source files, test fixtures, and associated plan documentation into a single, cohesive, atomic commit.
- Every commit must represent a complete unit of functional work or guideline milestone.

### Rule 2: Immediate Git Push to Remote
- **ALWAYS** execute `git push origin <branch>` immediately after creating any commit.
- Never leave local commits unpushed across turns. The local branch and `origin/main` must stay 100% synchronized.

### Rule 3: No Builds or Full CI Runner During Routine Turns
- **NEVER** execute `npm run build` or `go build ./...` during routine coding guideline audits (styling, naming, booleans, enums).
- **NEVER** execute `python 03-ai-scripts/06-cicd-local-runner.py` during routine turns.
- Use targeted, single-file linters (e.g. `check-newline-styling.py`, `check-boolean-guidelines.py`) for lightweight, instantaneous verification.
- Full builds and pipeline runs are reserved strictly for explicit owner commands or the final pre-release gate.

### Rule 4: Atomic File Change Cache Recording
- Whenever files are modified, record their relative paths into `.lovable/temp/recent-file-changes.json` under atomic file lock:
  ```bash
  python 03-ai-scripts/33-test-inventory-generator.py --record <files...>
  ```
- This maintains the cache needed by the test inventory generator and pre-release validation gates.
