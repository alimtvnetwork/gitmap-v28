# Master Plan: 365 — GitMap Installation Suggestion Issue

- **Plan Slug:** `365-gitmap-installation-suggestion-issue`
- **Root Slug:** `gitmap-installation-suggestion-issue`
- **Specification:** [02-spec/21-app/365-gitmap-installation-suggestion-issue/readme.md](../../../02-spec/21-app/365-gitmap-installation-suggestion-issue/readme.md)
- **Status:** `pending`
- **Creation Date:** `31-Oct-26`
- **Budget:** `400`

---

## 1. Executive Summary

This plan remediates the intermittent appearance and disappearance of installation suggestions in the terminal during GitMap installation across 6 failure modes. It enforces deterministic visual reporting in `install.ps1`, brings complete feature parity to `local-install.ps1`, expands Cobra argument completion to the short alias `gm`, supports unversioned repository URLs in release hints, and enriches `gitmap install gitmap` with actionable guidance.

---

## 2. Subtask Breakdown

| Code | Slug | Title | Status |
|---|---|---|---|
| TASK-01 | `gitmap-installation-suggestion-issue - Task 01` | Specification and Root Cause Analysis Documentation | completed |
| TASK-02 | `gitmap-installation-suggestion-issue - Task 02` | Master Execution Plan and Subtask Artifacts Creation | in-progress |
| TASK-03 | `gitmap-installation-suggestion-issue - Task 03` | Installer Scripts Fix — Profile Suggestion Persistence & Always-Report Logic | pending |
| TASK-04 | `gitmap-installation-suggestion-issue - Task 04` | Local Installer Parity — Integrate Suggestions into `local-install.ps1` | pending |
| TASK-05 | `gitmap-installation-suggestion-issue - Task 05` | Cobra Argument Completion for Short Aliases & Release Install Hints | pending |
| TASK-06 | `gitmap-installation-suggestion-issue - Task 06` | Build, Minor Bump `v6.639.0`, Release & Fleet Verification | pending |

---

## 3. Invariants & Rules

1. **Relative Git Paths:** Strict adherence to repo-relative paths (`cli/...`, `02-spec/...`).
2. **Coding Guidelines:** Functions <= 15 lines, positive booleans (`is*`, `has*`), zero nested ifs.
3. **No Unrequested Tests:** Zero test suite runs in external downstream repos.
4. **GitMap Semantic Commits:** Release using `gitmap cpr "v6.639.0 ..."`.
