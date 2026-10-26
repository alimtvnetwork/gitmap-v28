# Subtask 332.3 — Context Menu Install & Repositories Sections and Terminal Help

> **Parent Plan:** `.ai-memory/plans/pending/332-gitmap-ui-and-functionality-improvements.md`  
> **Target Packages:** `cli/cmdinstall/`, `cli/cmd/`  
> **Status:** Pending  

---

## 1. Objectives
1. In `cli/cmdinstall/installctxentries.go`:
   - Add `Install ▸` submenu in `ctxCategoryEntries()` containing `Install Wizard`, `Dev Workstation Profile`, `Tools List`, `Antigravity`, and `Clean Duplicate PowerShell`.
   - Add `Repositories ▸` submenu containing `Open Repositories Web UI`, `List Tracked Repositories`, `Status All Repositories`, and `Pull All Repositories`.
2. In `cli/cmd/rootusage_groups.go`:
   - Update `printGroupGitOps()` to include `merge-right` (`mr`), `merge-left` (`ml`), `merge-both` (`mb`), `compare` (`cmp`), and `show-log`.
3. In `cli/cmd/roothelp_clusters.go`:
   - Update Core Repo Operations cluster with `merge-right`, `merge-left`, `compare`, and `show-log`.
4. Add unit tests for context menu structure and terminal help rendering.
