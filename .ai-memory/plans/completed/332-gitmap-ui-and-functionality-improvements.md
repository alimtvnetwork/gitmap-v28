# Plan 332 — GitMap UI and Functionality Improvements

> **Status:** Pending  
> **Sequence ID:** 332  
> **Slug:** `gitmap-ui-and-functionality-improvements`  
> **Specification:** `02-spec/21-app/332-gitmap-ui-and-functionality-improvements/`  
> **Subtasks Directory:** `.ai-memory/plans/subtasks/332-gitmap-ui-and-functionality-improvements/`  
> **Date:** 2026-10-26  

---

## 1. Executive Summary & Problem Context
The user identified crucial usability gaps across the Show Log Web UI, single-path diff execution, context menu sections, and terminal help:
1. Show Log on folders fails or truncates commits due to unnormalized path matching.
2. The Show Log UI lacks independent scrollbars, modern aesthetic styling, and responsive layout.
3. Show Diff fails when invoked with a single path or file, emitting bare errors without stack traces.
4. The context menu lacks an `Install ▸` section and a `Repositories ▸` list section.
5. The Web UI lacks an `Exit` button to stop the local HTTP server and close the browser.
6. The Web UI lacks right-click actions (copy commit/subject, copy path, revert).
7. Double-click in modified files table does not launch Beyond Compare reliably.
8. Terminal help text is missing `merge-right` (`mr`), `merge-left` (`ml`), `merge-both` (`mb`), `compare` (`cmp`), and `show-log`.

---

## 2. Subtasks Roadmap

| Subtask ID | File | Description | Target | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Subtask-01** | `01-subtask-show-log-folder-and-diff.md` | Fix folder log path normalization, implement single-path working tree diff with stack traces in `cmddiff`, and test diff error reporting. | `cli/cmdlogviewer/`, `cli/cmddiff/` | Pending |
| **Subtask-02** | `02-subtask-modern-scroll-responsive-ui.md` | Overhaul Show Log UI with modern Catppuccin theme, styled scrollbars, responsive CSS, exit endpoint `/api/log/exit`, and web right-click menu with copy & revert. | `cli/cmdlogviewer/` | Pending |
| **Subtask-03** | `03-subtask-context-menu-and-terminal-help.md` | Add `Install ▸` and `Repositories ▸` submenus in `installctxentries.go`, update terminal help in `rootusage_groups.go` and `roothelp_clusters.go` with merge right/left and compare. | `cli/cmdinstall/`, `cli/cmd/` | Pending |

---

## 3. Verification & Acceptance Criteria
- [ ] Show Log on folders (`gitmap show-log <folder>`) correctly displays commits affecting that folder.
- [ ] Web UI has responsive 3-pane layout, styled scrollbars, and works cleanly on mobile/narrow windows.
- [ ] Clicking "Exit" in Web UI cleanly closes browser and terminates local server.
- [ ] Right-clicking in Web UI allows copying commit hash, subject, file path, and triggering revert.
- [ ] Double-clicking file row launches Beyond Compare or configured external diff tool.
- [ ] `gitmap diff <file>` or `gitmap diff` works cleanly on working tree changes with stack traces on errors.
- [ ] Context menu includes `Install ▸` (Wizard, Profiles, Dev Tools) and `Repositories ▸` (Web UI, List, Status, Pull All).
- [ ] Terminal help (`gitmap --help` and `gitmap help gitops`) clearly lists `merge-right`, `merge-left`, `merge-both`, and `compare`.
- [ ] All unit tests pass 100% across touched packages.
