# Plan 361: Terminal Help Parity for Merge, Commit-In, and AI Commands

- **Plan Slug:** `terminal-help-merge-and-ai-commands`
- **Root Slug:** `gitmap-feature-implementation-and-update`
- **Spec Reference:** `02-spec/21-app/361-terminal-help-merge-and-ai-commands/01-architecture-spec.md`
- **Status:** `completed`
- **Created At:** 2026-10-31
- **Completed At:** 2026-10-31

---

## 1. Goal

Implement comprehensive terminal help (`thelp`) parity across GitMap:
1. Alias `thelp` and `th` to the core help dispatcher so `gitmap thelp` works everywhere.
2. Route `merge` subcommands (`right`, `left`, `both`, `in`) and rich help menu via `cmdmerge.RenderMergeHelp()`.
3. Support `merge-in` as a first-class alias for `commit-in`.
4. Update `settings` and `login` terminal help with all AI Auto commands and configuration keys.
5. Provide curated examples for `merge`, `merge-right`, `merge-left`, `merge-in`, `commit-in`, `thelp`, `specific`, `login`, and `settings`.
6. Run tests, bump minor version, and perform release ceremony.

---

## 2. Subtasks Breakdown

### Task 01: Dispatcher Aliases & Help Architecture [completed]
- **Slug:** `gitmap-feature-implementation-and-update - Task 01`
- **Subtask File:** `.ai-memory/plans/subtasks/361-terminal-help-merge-and-ai-commands/01-dispatcher-aliases-and-help-architecture.md`
- **Description:** Add `thelp` and `th` to `utilityCoreBaseEntries()`, add `merge-in` to `commit-in` entries, map in `normalizeHelpTopic`, and handle `thelp specific`.

### Task 02: Rich Merge Help Menu & Subcommand Dispatcher [completed]
- **Slug:** `gitmap-feature-implementation-and-update - Task 02`
- **Subtask File:** `.ai-memory/plans/subtasks/361-terminal-help-merge-and-ai-commands/02-rich-merge-help-menu-and-subcommand-dispatcher.md`
- **Description:** Create `cli/cmdmerge/merge_help_menu.go` (`RenderMergeHelp`), wire into `rich_help_dispatcher.go` and `cmdmerge.DispatchMoveMerge`.

### Task 03: Settings, Login, and Commit Terminal Help Enhancements [completed]
- **Slug:** `gitmap-feature-implementation-and-update - Task 03`
- **Subtask File:** `.ai-memory/plans/subtasks/361-terminal-help-merge-and-ai-commands/03-settings-login-commit-help-enhancements.md`
- **Description:** Update `cli/cmdcommit/commit_help_menu.go`, `cli/cmdagy/agy_settings_help.go`, create `cli/helpdoc/settings.md`, `cli/helpdoc/merge.md`, `cli/helpdoc/merge-in.md`, and update `cli/helpdoc/login.md`.

### Task 04: Curated Help Examples Registry & Tests [completed]
- **Slug:** `gitmap-feature-implementation-and-update - Task 04`
- **Subtask File:** `.ai-memory/plans/subtasks/361-terminal-help-merge-and-ai-commands/04-curated-help-examples-and-tests.md`
- **Description:** Register curated examples in `cli/helpdoc/examples_merge.go` and `cli/helpdoc/examples_ops.go` for all target commands, verify via targeted tests.

### Task 05: Release Ceremony & Verification [completed]
- **Slug:** `gitmap-feature-implementation-and-update - Task 05`
- **Subtask File:** `.ai-memory/plans/subtasks/361-terminal-help-merge-and-ai-commands/05-release-ceremony-and-verification.md`
- **Description:** Bump minor version to `v6.635.0`, update changelog and release notes, commit via `gitmap cpr`, and push to remote.
