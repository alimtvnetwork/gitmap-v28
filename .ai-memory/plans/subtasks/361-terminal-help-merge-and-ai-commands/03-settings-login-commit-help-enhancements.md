# Subtask 03: Settings, Login, and Commit Terminal Help Enhancements

- **Plan:** `361-terminal-help-merge-and-ai-commands`
- **Slug:** `gitmap-feature-implementation-and-update - Task 03`
- **Status:** `pending`

## Steps
1. In `cli/cmdcommit/commit_help_menu.go`, add `--ai-agent <name>` flag and tips, and include `merge-right`, `merge-left`, `merge-in` (`commit-in`).
2. In `cli/cmdagy/agy_settings_help.go`, add `set ai_auto.default_name`, `set ai_auto.default_email`, and `set ai_auto.default_agent` to usage lines and description entries.
3. In `cli/cmd/rich_help_dispatcher.go`, wire `"settings"`, `"setting"`, `"speed-settings"` to `cmdagy.RenderAgySettingsHelp()`.
4. Create `cli/helpdoc/settings.md`, `cli/helpdoc/merge.md`, `cli/helpdoc/merge-in.md`.
5. Update `cli/helpdoc/login.md` with full details on `login set ai-auto <agent>` and `login set global ai-auto`.
