# Plan 329 — GitMap Context Menu and Version Fix

> **Status:** Completed  
> **Sequence ID:** 329  
> **Slug:** `gitmap-context-menu-and-version-fix`  
> **Specification:** `02-spec/21-app/329-gitmap-context-menu-and-version-fix/`  
> **Subtasks Directory:** `.ai-memory/plans/subtasks/329-gitmap-context-menu-and-version-fix/`

---

## 1. Executive Summary & Problem Context
1. **Version Downgrade Bug:** `gitmap update` downgraded `v6.613.0` to `v6.610.0` due to equality-only version checks and unvalidated sibling probe resolution.
2. **Context Menu Feature Expansion:** Provide right-click context menu installation across Windows, Linux, and macOS (`gitmap install context-menu`), offering TortoiseGit-like features (`Commit to main`, `Diff`, `Diff later`, `Diff with previous version`, `Show log`, `Branch comparison`) and an Install Wizard sub-item on empty space.
3. **Show Log Dual-Pane UI:** Interactive dual-pane commit & file history UI with automatic Beyond Compare launcher integration on double-click.

---

## 2. Subtasks Roadmap

| Subtask ID | File | Description | Owner | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | `01-subtask-version-update-downgrade-fix.md` | Monotonic SemVer update guard, `--allow-downgrade` flag, and sibling probe monotonicity filter. | Worker 01 | Completed |
| **Subtask-02** | `02-subtask-context-menu-cross-platform.md` | Cross-platform context menu engine (`gitmap install context-menu`), registry/script actions, and install wizard launcher. | Worker 02 | Completed |
| **Subtask-03** | `03-subtask-show-log-ui-and-diff-tools.md` | Dual-pane Show Log UI (`gitmap show-log`), commit/file APIs, and Beyond Compare launcher integration. | Worker 03 | Completed |

---

## 3. Verification & Acceptance Criteria
- [x] `gitmap update` skips updating when target version <= current version unless `--allow-downgrade` is provided.
- [x] `gitmap install context-menu` registers the updated Git submenu and Install Wizard sub-item.
- [x] `gitmap show-log` displays commits and changed files with double-click opening Beyond Compare if installed.
- [x] Zero absolute paths in repository.
- [x] FastGate static linter passes 100%.
- [x] Unit test suites pass 100% (`./cmdupdate`, `./cmdcontextmenu`, `./cmdinstall`, `./cmdlogviewer`).
