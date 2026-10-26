# Master Audit Ledger: GitMap UI and Functionality Improvements

> **Module:** `02-spec/21-app/332-gitmap-ui-and-functionality-improvements`  
> **Sequence ID:** 332  
> **Date:** 2026-10-26  
> **Status:** Draft / Active Specification  
> **Author:** Antigravity Architect  

---

## 1. Verbatim User Request & Analysis

> "Hi there. In the `gitmap` menu, the right-click one, let's say D for anything else, when we try to showcase folder from Show Log, it just does not do things properly. I think we need to work on that. The UI that actually created there does not have a scroll or things like that. That is a pretty bad UI. And it is not modern UI, that's for sure. You need to work on it. Also, if I click on Show Diff, I don't see anything. If I do show diff, it just gives the error, and error does not have the stack trace. That's pretty bad. The context menu does not have an install section, and also from the browser, we should have an exit button where we could exit this out and close the browser as well. By clicking on the exit button, it would automatically close the browser, close the endpoint, things like that. It is not happening. It's just stuck together. There should be copy option. We should have a right-click option in the web UI to revert something or open the text or the items. Currently, it's not really viewable. Double-click should open to a diff merge or merge comparison view tool, like beyond compare or anything. That is not also happening, so it does not make much sense what we have right now. The UI is not responsive that you have created. You need to fix that. Also in the `gitmap`, there is no install section. I mentioned there should be an install section, which is not shown. The repo list is not there. There should be a repo list section in the context menu. Also try to add the help text for the merge right, merge left in the help text of the terminal. A lot of the new commands added which is not there in the help text of the terminal. Make sure find those and fix those."

---

## 2. Forensic Analysis & Root Cause Findings

### 2.1 Problem 1: Show Log on Folder Does Not Work Properly
- **Mechanism:** When launched on a folder (e.g. `cli/cmd`), `git log` executed without path normalization or trailing separator hygiene, and if sub-commits only touched subfolders without clean relative normalization from the repository root, `git log -- <targetPath>` returned 0 or truncated commits.
- **Root Cause:** Path resolution in `ResolveTargetRepoAndPath()` and `QueryCommits()` passed unnormalized paths to `git log`. If a folder has subpaths or backslashes on Windows, `git log` fails to match tree entries.

### 2.2 Problem 2: UI Lacks Scrolling, Not Modern, and Not Responsive
- **Mechanism:** In `cli/cmdlogviewer/ui_assets_html.go`, `.pane-upper`, `.pane-middle`, and `.pane-lower` had rigid flex ratios with fixed minimum heights (`min-height: 140px`) and missing responsive media queries (`@media (max-width: 900px)`). The commit body and files list did not have independent, styled scrollbars or responsive collapse on narrow displays.
- **Root Cause:** Absence of responsive layout grid/flex rules, custom modern web scrollbars (`::-webkit-scrollbar`), and fluid height allocation.

### 2.3 Problem 3: Show Diff Fails with Error & Missing Stack Trace
- **Mechanism:** `gitmap diff` currently requires two endpoints (`gitmap diff LEFT RIGHT`). In the context menu, right-clicking a file or folder invoked `gitmap diff "%1"`, which failed with `Usage: gitmap diff LEFT RIGHT [flags]`.
- **Root Cause:**
  1. `gitmap diff` dispatch in `cmddiff` did not recognize single-argument invocations (`gitmap diff <path>`). When a single path is given, it must compare working tree changes for that path against `HEAD` via external compare tool (Beyond Compare / VS Code) or show formatted working tree diff with complete AppError stack traces.
  2. The web API `/api/log/diff` returned generic JSON errors without detailed `apperror` stack traces when compare tools failed.

### 2.4 Problem 4: Context Menu Missing Install Section & Repo List Section
- **Mechanism:** The context menu in `installctxentries.go` had direct actions and submenus for `Scan`, `Clone`, `Release`, and `Visibility`, but had no dedicated `Install ▸` submenu or `Repositories ▸` list submenu.
- **Root Cause:** `ctxCategoryEntries()` omitted an `Install ▸` submenu (e.g. `Install Wizard`, `Tools`, `Profiles`) and had only a single static Web UI entry for repos rather than a structured `Repositories ▸` listing.

### 2.5 Problem 5: Missing Browser Exit & Shutdown Endpoint
- **Mechanism:** The `show-log` server ran indefinitely until the terminal process was terminated. There was no "Exit" button in the Web UI to terminate the local HTTP server and close the browser tab.
- **Root Cause:** Missing `/api/log/exit` (or `/api/shutdown`) HTTP handler in `cmdlogviewer/server.go` and missing shutdown button in `ui_assets_html.go`.

### 2.6 Problem 6: Missing Context Menu / Copy / Revert in Web UI
- **Mechanism:** Users had no right-click context menu in the Show Log web table to copy commit hash, copy subject, revert commit, or open file.
- **Root Cause:** Absence of custom context menu event listeners (`contextmenu`) in `ui_assets_script.go`.

### 2.7 Problem 7: Missing Double-Click to External Diff Tool (Beyond Compare)
- **Mechanism:** Double-clicking files in the modified files table failed if the diff request body was rejected or if Beyond Compare path was not resolved properly on Windows.
- **Root Cause:** Double-click handler triggered `/api/log/diff` which failed on certain path formats and did not fallback gracefully.

### 2.8 Problem 8: Missing Help Text in Terminal for `merge-right`, `merge-left`, and New Commands
- **Mechanism:** `merge-right` (`mr`), `merge-left` (`ml`), `merge-both` (`mb`), `compare` (`cmp`), and `show-log` were implemented in code but completely omitted from the terminal help screens (`printGroupGitOps()`, `printUsage()`, and semantic clusters).
- **Root Cause:** Help registration was not updated when new commands were introduced.
