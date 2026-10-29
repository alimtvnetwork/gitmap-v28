# GitMap UI and Functionality Improvements

**Version:** 6.630.0
**Updated:** 2026-10-29

- **Sequence:** 354
- **Slug:** `gitmap-ui-and-functionality-improvements`
- **Status:** pending
- **Specification:** `02-spec/21-app/354-gitmap-ui-and-functionality-improvements/03-context-menu-and-retirement-spec.md`

File ownership is disjoint. Each path below belongs to one subtask. An owner edits only the files on that subtask.

## Subtasks

### ST-01 theme hover

- **Owner:** Worker 01
- **Status:** pending
- **Files:** `cli/cmdui/ui_assets_theme.go`
- **Work:** Hover rules in the theme asset move controls (`translateX` on navigation buttons, `translateY` on cards and buttons). Remove hover movement from this file. Color and border changes on hover may remain. Leave `.ai-memory/strictly-avoid.md` to ST-05.

### ST-02 ssh ping mount and online probe

- **Owner:** Worker 02
- **Status:** pending
- **Files:** `cli/cmdui/ui_server.go`, `cli/cmdui/ui_server_ssh.go`
- **Work:** The server mux mounts SSH node, export, import, and deploy-key routes and does not mount `/api/ssh/ping` or `/api/ssh/test`. Mount those two routes on the server mux. Node summaries currently set online to true for every connection. Replace that constant with a real online probe and return the probed state. Leave the settings script to ST-03.

### ST-03 stop stubbing testNode; settings debug toggle

- **Owner:** Worker 03
- **Status:** pending
- **Files:** `cli/cmdui/ui_assets_js_settings.go`
- **Work:** The dead-handler list assigns `window.testNode` to a "Not implemented yet" toast and covers the real ping function. Remove `testNode` from that stub list. Add a settings debug toggle in this same file. Do not edit the SSH server files (ST-02) or the nodes script.

### ST-04 help contrast

- **Owner:** Worker 04
- **Status:** pending
- **Files:** `cli/cmdui/ui_server_help.go`, `cli/cmdui/ui_assets_js_prompts_help.go`
- **Work:** Help topic names, summaries, and prompt badges use low-contrast foregrounds against the theme background. Raise that contrast in these two files so help text stays readable on the dark theme and the light theme.

### ST-05 strictly avoid hover movement

- **Owner:** Worker 05
- **Status:** pending
- **Files:** `.ai-memory/strictly-avoid.md`
- **Work:** Record the rule that hover must not move the hovered element. This file is the rule. The theme asset change belongs to ST-01.

### ST-06 context menu

- **Owner:** Held until the context-menu Go files are released
- **Status:** blocked
- **Files:** `cli/cmdinstall/installctxentries.go`, `cli/cmdinstall/installctx.go`, `cli/cmddiff/single_diff.go`, `cli/cmdgit/compare_launcher.go`
- **Work:** Specified in `02-spec/21-app/354-gitmap-ui-and-functionality-improvements/03-context-menu-and-retirement-spec.md`. Another agent owns these files. Do not edit them from this plan while they are dirty.
- **Required when unblocked:** Write a registry separator (`CommandFlags` `0x00000020`) on Commit to main so the diff group (Diff, Diff later, Diff with previous version, Show log, Blame) is separated from the commit group. Call `Set-Location` only when `%V` is a directory. When a file has no diff against `HEAD`, print a clear no-changes line and do not launch Beyond Compare.

### ST-07 catalog deletion

- **Owner:** Held until heatmap clearance
- **Status:** blocked pending heatmap clearance
- **Files:** `02-spec/21-app/340-next-100-ui-improvements`, `02-spec/21-app/341-next-100-ui-improvements`, `02-spec/21-app/342-next-100-ui-improvements`, `02-spec/21-app/343-next-100-ui-improvements`, `02-spec/21-app/344-next-100-ui-improvements`, `02-spec/21-app/345-next-100-ui-improvements`, `02-spec/21-app/346-next-100-ui-improvements`
- **Work:** After the heatmap no longer lists these generic catalogs as active, delete those seven directories from the working tree in a normal commit. Do not rewrite git history.

## Blocked record: time restriction

- **Id:** 349
- **Path:** `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling`
- **Status:** blocked
- **Reason:** This directory is the time-restriction content. The heatmap lists that task as active.
- **Rule:** Keep the directory. Do not delete it. Do not rewrite git history. ST-07 must not include this path. The 340–346 range ends at 346.

## Disjoint ownership

| Subtask | Owner | Paths |
| --- | --- | --- |
| ST-01 | Worker 01 | `cli/cmdui/ui_assets_theme.go` |
| ST-02 | Worker 02 | `cli/cmdui/ui_server.go`, `cli/cmdui/ui_server_ssh.go` |
| ST-03 | Worker 03 | `cli/cmdui/ui_assets_js_settings.go` |
| ST-04 | Worker 04 | `cli/cmdui/ui_server_help.go`, `cli/cmdui/ui_assets_js_prompts_help.go` |
| ST-05 | Worker 05 | `.ai-memory/strictly-avoid.md` |
| ST-06 | Held | `cli/cmdinstall/installctxentries.go`, `cli/cmdinstall/installctx.go`, `cli/cmddiff/single_diff.go`, `cli/cmdgit/compare_launcher.go` |
| ST-07 | Held | `02-spec/21-app/340-next-100-ui-improvements` through `02-spec/21-app/346-next-100-ui-improvements` |

No path appears on more than one subtask. Sequence 349 is a blocked record and is outside ST-07.
