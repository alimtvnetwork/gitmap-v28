# Context Menu Icons and Non-Git Error Handling - Component & UI Specification

## 1. Component Contracts

### 1.1 `cli/cmdlogviewer/server.go`
- `handleNonGitShowLog(args []string, isNoBrowser bool) error`:
  - Extracts sanitized target directory from args.
  - Emits friendly notice: `  ✗ '<target>' is not a Git repository.`.
  - Invokes `findChildRepos(target)` to discover nested repositories.
  - Emits detected repositories list or friendly tip.
  - If `!isNoBrowser && FallbackDashboardUIFn != nil`:
    - Emits `  ⚡ Launching GitMap Web UI dashboard at http://127.0.0.1:8080 ...`.
    - Returns `FallbackDashboardUIFn("search", 8080)`.
  - Returns `nil`.
- `findChildRepos(dir string) []string`:
  - Probes direct child repos; if none found, probes nested children (depth 2) skipping noise directories (`node_modules`, `vendor`, `.git`, `AppData`, `Windows`, `Program Files`, `$Recycle.Bin`).

### 1.2 `cli/cmddiff/single_diff.go`
- `runSinglePathDiff(target string) error`:
  - If `!isGitRepo`, emits `  ✗ '<target>' is not a Git repository.` and guidance tip, and returns `nil`.

### 1.3 `cli/constants/constants_installctx.go`
- Defines high-definition icon constants:
  - `CtxIconShowLog` = `"shell32.dll,264"` (History/Clock)
  - `CtxIconDiff` = `"shell32.dll,246"` (Comparison)
  - `CtxIconBranchCompare` = `"shell32.dll,133"` (Tree hierarchy)
  - `CtxIconCommit` = `"shell32.dll,297"` (Green checkmark)
  - `CtxIconPull` = `"shell32.dll,146"` (Arrow down)
  - `CtxIconPush` = `"shell32.dll,147"` (Arrow up)
  - `CtxIconScan` = `"shell32.dll,22"` (Search)
  - `CtxIconClone` = `"shell32.dll,134"` (Copy/Clone)
  - `CtxIconRelease` = `"shell32.dll,265"` (Package/Star)
  - `CtxIconRepos` = `"shell32.dll,4"` (Folder stack)
  - `CtxIconWebUI` = `"shell32.dll,14"` (Globe)
  - `CtxIconSettings` = `"shell32.dll,315"` (Gear)
  - `CtxIconTerminal` = `"shell32.dll,220"` (Console)
  - `CtxIconDocs` = `"shell32.dll,266"` (Book)
  - `CtxIconHelp` = `"shell32.dll,23"` (Question mark)
  - `CtxIconInstall` = `"shell32.dll,260"` (Box)
  - `CtxIconMacro` = `"shell32.dll,250"` (Play script)
  - `CtxIconVisibility` = `"shell32.dll,174"` (Lock)

### 1.4 `cli/cmdinstall/installctxentries.go`
- Assigns distinct icons across direct entries, background entries, categories, and child actions.
- Respects `isGitActionsHidden()` to suppress direct git actions when configured.

### 1.5 `cli/cmdinstall/installctx_git_actions.go`
- Manages `IsGitActionsHidden() bool` and `SetGitActionsHidden(hidden bool) error`.
- Persists state in Windows registry under `constants.CtxRootKeyDirectory` (`HideGitActions` REG_DWORD).

### 1.6 `cli/cmdcontextmenu/contextmenu.go`
- Subcommands: `hide-git-actions`, `show-git-actions`, `toggle-git-actions`.
- Flag: `--hide-git-actions` for `gitmap context-menu install`.
- Status output displays Git actions visibility.

### 1.7 `cli/cmdui/ui_server_contextmenu.go` & `ui_assets_markup_settings_b.go`
- REST API endpoint `/api/context-menu/toggle-git-actions`.
- Settings UI button: "Toggle Git Actions Visibility" with dynamic label and toast feedback.
