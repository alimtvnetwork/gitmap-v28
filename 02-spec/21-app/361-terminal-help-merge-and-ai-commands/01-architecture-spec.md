# Architecture Specification: Terminal Help Parity for Merge, Commit-In, and AI Commands

- **Feature Title:** Terminal Help Parity for Merge, Commit-In, and AI Commands
- **Sequence:** 361
- **Status:** In Progress
- **Created:** 2026-10-31

## 1. Overview & Problem Statement

GitMap provides high-velocity command-line capabilities for multi-repo commit replay, directional merging (`merge-right`, `merge-left`, `merge-both`, `commit-in`), and AI agent identity management (`login set ai-auto`, `settings set ai_auto.*`, `--ai-agent <name>`). However:
1. `thelp` was not aliased to `help`, causing `gitmap thelp` to exit with an unrecognized command error.
2. `gitmap help merge` fell back to a generic filter match list rather than displaying a dedicated rich merge suite help menu.
3. `gitmap merge` command crashed with unrecognized command error instead of displaying merge help or routing subcommands (`right`, `left`, `both`, `in`).
4. `merge-in` had no alias mapping to `commit-in` (`cin`).
5. Terminal help for `login` and `settings` lacked documentation and executable examples for `ai-auto` commands and `--ai-agent` flags.
6. Topic-specific queries like `thelp merge`, `thelp merge-right`, `thelp merge-left`, `thelp merge-in`, `thelp commit-in`, `thelp login`, and `thelp settings` lacked curated example entries.

## 2. Architectural Design & Changes

### A. Terminal Help (`thelp`) Alias
- In `cli/cmd/rootutility.go`, add `"thelp"` and `"th"` to `utilityCoreBaseEntries()` under `constants.CmdHelp`.
- Support `gitmap thelp` as a direct alias for `gitmap help`.
- Support `gitmap thelp specific` (and `gitmap help specific`) with a dedicated guide on querying specific command help.

### B. Merge Command Routing & Rich Help
- In `cli/cmdmerge/dispatchmovemerge.go`:
  - Handle `command == "merge"`:
    - If args start with `right`: dispatch `merge-right`
    - If args start with `left`: dispatch `merge-left`
    - If args start with `both`: dispatch `merge-both`
    - If args start with `in`: dispatch `commit-in`
    - If bare or with `-h`/`--help`/`help`: call `RenderMergeHelp()`
- In `cli/cmdmerge/merge_help_menu.go`:
  - Create `RenderMergeHelp()` utilizing `termout.RenderMenu()` to display a structured 2-column menu for `merge-right`, `merge-left`, `merge-both`, `merge-in`, `commit-right`, `commit-left`, `commit-both`, and `merge write`.
- In `cli/cmd/rich_help_dispatcher.go`:
  - Route `"merge"`, `"merge-right"`, `"merge-left"`, `"merge-both"`, `"merge-in"` to `cmdmerge.RenderMergeHelp()`.

### C. `merge-in` Alias for `commit-in`
- In `cli/cmd/rootcore.go`, add `"merge-in"` and `"mergein"` to `constants.CmdCommitIn` entry.
- In `cli/cmd/rootutility_help.go`, map `"merge-in"` and `"mergein"` to `"commit-in"` in `normalizeHelpTopic`.
- Create `cli/helpdoc/merge.md` and `cli/helpdoc/merge-in.md`.

### D. Settings & Login Terminal Help Parity
- In `cli/cmdagy/agy_settings_help.go`:
  - Add `set ai_auto.default_name`, `set ai_auto.default_email`, and `set ai_auto.default_agent` to usage lines and description entries.
- Create `cli/helpdoc/settings.md` documenting speed settings, AI Auto keys, and CLI usage.
- Update `cli/helpdoc/login.md` documenting `login set ai-auto <agent>`, `login set global ai-auto`, and examples.
- In `cli/cmd/rich_help_dispatcher.go`, route `"settings"`, `"setting"`, `"speed-settings"` to `cmdagy.RenderAgySettingsHelp()`.
- In `cli/cmdcommit/commit_help_menu.go`:
  - Add `--ai-agent <name>` flag and tips.
  - Add `merge-right`, `merge-left`, `merge-in` (`commit-in`) entries and tips.

### E. Curated Examples Registry (`cli/helpdoc/examples_merge.go`, `examples_ops.go`)
- Register curated examples for:
  - `merge`
  - `merge-in`
  - `merge-right`
  - `merge-left`
  - `merge-both`
  - `commit-in`
  - `commit-right`
  - `commit-left`
  - `commit-both`
  - `thelp`
  - `specific`
  - `login`
  - `settings`

## 3. Invariants & Coding Guidelines
- Leaf functions return errors, callers handle exits.
- Function length <= 15 lines (<= 8 lines preferred).
- Positive booleans (`is*`, `has*`).
- Relative Git paths only.
- Semantic commit verbs (`gitmap cpr` for release).
