# 329.3 — Show Log UI & Beyond Compare Diff Integration Specification

> **Status:** Drafted  
> **Subsystem:** `cli/cmdlogviewer/`, `cli/cmdui/`, `cli/cmdgit/`  
> **Parent Ledger:** `02-spec/21-app/329-gitmap-context-menu-and-version-fix/00-master-audit-ledger.md`  
> **Rule Compliance:** Leaf Error Sovereignty, Positive Booleans (`is*`/`has*`), Functions <= 15 lines.

---

## 1. Overview & Operational Requirements

The user requested an interactive, visual "Show Log" interface modeling TortoiseGit's Log Messages window (as captured in `assets/screenshots/gitmap-context-menu-and-version-fix-03.png` and `04.png`).

### Key Ergonomics & User Flow
1. **Invocation:** User right-clicks a file or directory in Windows Explorer / Nautilus / Finder and selects `GitMap ▸ Show log` (or runs `gitmap show-log [path]` or `gitmap log ui [path]`).
2. **Dual-Pane Layout:**
   - **Upper Pane (Commits Table):** Displays commit history affecting the selected path (Graph indicator, Short Hash, Author, Date/Time, Subject). Filterable by branch or keyword.
   - **Middle Pane (Commit Details):** Displays the full commit message, hash, author email, and parent hashes of the currently selected commit.
   - **Lower Pane (Modified Files Table):** Displays the files affected by the selected commit, with Action type (`Added`, `Modified`, `Deleted`), Path, and line stat delta (`+N / -M`).
3. **Interactive Actions:**
   - **Refresh Button:** Reloads git history on demand without restarting the server or browser.
   - **View Options:** Toggle between compact and detailed commit graphs, filter by branch.
   - **Double-Click Diff Launching:** Double-clicking any file in the lower pane immediately opens an external visual compare tool comparing that file version against its parent commit (or working tree).
   - **Branch Comparison:** Compare the selected path across two local or remote branches.

---

## 2. External Diff Tool Discovery & Beyond Compare Priority

GitMap leverages `cli/cmdgit/compare_launcher.go` to locate and execute external diff tools.

### Beyond Compare Discovery Order
1. Check PATH for `bcomp`, `bcompare`, `BComp.exe`, `BCompare.exe`.
2. Windows standard installation locations:
   - `C:\Program Files\Beyond Compare 5\BComp.exe`
   - `C:\Program Files\Beyond Compare 5\BCompare.exe`
   - `C:\Program Files\Beyond Compare 4\BComp.exe`
   - `C:\Program Files\Beyond Compare 4\BCompare.exe`
   - `C:\Program Files (x86)\Beyond Compare 4\BComp.exe`
   - `C:\Program Files (x86)\Beyond Compare 4\BCompare.exe`
3. Linux / macOS locations:
   - `/usr/bin/bcompare`, `/usr/local/bin/bcompare`
   - `/Applications/Beyond Compare.app/Contents/MacOS/bcomp`
4. **Fallback Chain:** If Beyond Compare is not found, dynamically fallback to VS Code (`code --diff`), Meld, WinMerge, or GitMap's built-in web diff viewer.

### Diff Execution Workflow
When diffing a file from a commit:
1. Extract file content from `<commit>^:<file>` and `<commit>:<file>` into temporary files in `os.TempDir()`.
2. Invoke `launchToolWithPaths(tool, []string{tempLeft, tempRight})`.

---

## 3. Web Service & API Contract

The Show Log UI is served via GitMap's embedded HTTP server (`cli/cmdui/ui_server.go` or `cli/cmdlogviewer/`):
- `GET /api/log/commits?path=<path>&limit=100`: Returns JSON list of commits affecting `path`.
- `GET /api/log/commit?hash=<hash>`: Returns commit details and list of modified files with actions.
- `POST /api/log/diff`: Accepts `{ "commit": "<hash>", "path": "<file>", "tool": "auto" }` and triggers the diff launcher.
- `GET /api/log/branches`: Returns list of local and remote branches for branch comparison.
- `POST /api/log/branch-compare`: Accepts `{ "branchA": "<b1>", "branchB": "<b2>", "path": "<file>" }` and triggers visual diff.
