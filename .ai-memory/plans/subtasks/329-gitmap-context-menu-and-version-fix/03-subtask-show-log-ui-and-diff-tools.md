# Subtask 329.3 — Show Log Dual-Pane UI & Beyond Compare Launcher Integration

> **Parent Plan:** `.ai-memory/plans/pending/329-gitmap-context-menu-and-version-fix.md`  
> **Target Packages:** `cli/cmdlogviewer/`, `cli/cmdui/`, `cli/cmdgit/`  
> **Status:** Completed

---

## 1. Objectives
1. Implement `gitmap show-log [path]` and `gitmap log ui [path]` command entrypoints.
2. Build an embedded web UI modeling TortoiseGit's Log Messages window:
   - Upper Pane: Commit history table (commit hash, author, date, message).
   - Middle pane: Selected commit full details (hash, author email, commit message body).
   - Lower pane: Modified files list (Action: Added/Modified/Deleted, file path, changes).
   - Toolbar: Refresh button, branch selector / filter, branch comparison trigger.
3. Integrate external diff tool launcher (`cli/cmdgit/compare_launcher.go`):
   - Auto-detect Beyond Compare (`BComp.exe`, `BCompare.exe`, `/usr/bin/bcompare`).
   - On double-clicking a file in the lower pane, launch Beyond Compare (or fallback diff tool) comparing that file version against parent commit or working tree.
4. Add unit tests for log parser and diff tool invocation.

---

## 2. Completed Implementation & File Inventory

- `cli/cmdlogviewer/types.go`: CommitItem, CommitDetail, FileChange, DiffRequest/Response, BranchCompareRequest/Response, BranchesResponse.
- `cli/cmdlogviewer/git_log_parser.go`: `QueryCommits`, `QueryCommitDetail`, git log/show parsing with status and stats.
- `cli/cmdlogviewer/git_diff_extractor.go`: `TriggerExternalDiff`, `TriggerBranchDiff`, temporary file extraction for left/right revisions.
- `cli/cmdlogviewer/branches.go`: `QueryBranches`, branch validation.
- `cli/cmdlogviewer/server_handlers.go`: REST API endpoints for `/api/log/commits`, `/api/log/commit`, `/api/log/diff`, `/api/log/branches`, `/api/log/branch-compare`, `/api/log/ping`.
- `cli/cmdlogviewer/server.go`: Embedded loopback HTTP server, target path resolution, auto-open browser support.
- `cli/cmdlogviewer/ui_assets_html.go` & `ui_assets_script.go` & `ui_assets.go`: Dark theme dual/tri-pane UI with responsive tables, keyboard navigation, real-time search, double-click diff, branch comparison modal.
- `cli/cmdlogviewer/dispatch.go`: Show log command routing and alias matching.
- `cli/cmdgit/compare_launcher.go`: Exported `FindCompareTool`, `LaunchToolWithPaths`, and added async launcher `LaunchToolWithPathsAsync`.
- `cli/cmdgit/cmd_log.go`: Forwarding hook for `gitmap log ui [path]`.
- `cli/cmdui/show_log_shim.go`: Integration bridge into cmdui.
- `cli/cmd/dispatch_git_ops.go`: Root dispatch support for `show-log`, `showlog`, and `log ui`.

---

## 3. Verification & Test Gate Results

- Unit tests in `cli/cmdlogviewer`: 100% PASS (10 unit tests covering parser, handlers, ping, branch validation, diff validation, live git repo queries).
- Unit tests in `cli/cmdgit`: 100% PASS (5 unit tests covering args, Beyond Compare discovery, flag parsing).
- Full CLI build: `go build .` PASS (zero errors).
- FastGate (`03-ai-scripts/50-fastgate.py`): 100% PASS across all 19 files in 0.028s (<1.5s target).
  - Relative Paths: PASS
  - Nested Ifs: PASS (0 violations)
  - Boolean Guidelines: PASS (positive booleans only)
  - Go File Size: PASS (<500 lines)
