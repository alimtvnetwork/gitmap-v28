# Subtask 03: Git Actions Visibility Toggle & Web UI Integration

- **Status**: completed

## Objectives
- In `cli/cmdinstall/installctx_git_actions.go`:
  - Implement `IsGitActionsHidden() bool` and `SetGitActionsHidden(hidden bool) error`.
  - Persist setting in registry under `constants.CtxRootKeyDirectory` (`HideGitActions`).
- In `cli/cmdinstall/installctxentries.go`:
  - When `isGitActionsHidden()` is true, omit direct git entries (`01_show_log` through `06_push`).
- In `cli/cmdinstall/exports.go`:
  - Export `IsGitActionsHidden` and `SetGitActionsHidden`.
- In `cli/cmdcontextmenu/contextmenu.go`:
  - Register `hide-git-actions`, `show-git-actions`, `toggle-git-actions` subcommands and `--hide-git-actions` install flag.
  - Surface git action visibility in status report.
- In `cli/cmdui/ui_server_contextmenu.go`:
  - Add `IsGitActionsHidden` field in status response and register `/api/context-menu/toggle-git-actions`.
- In `cli/cmdui/ui_assets_markup_settings_b.go`:
  - Add "Toggle Git Actions" button in Context Menu settings card with dynamic status and toast notifications.

## Verification
- Test `gitmap context-menu hide-git-actions` and `gitmap context-menu show-git-actions`.
- Verified: Both commands and unit tests passed. Status reports visibility accurately.
