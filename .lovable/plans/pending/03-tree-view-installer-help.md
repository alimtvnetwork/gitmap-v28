# Parent Task 03: Tree-View Installer, Profile Summary, & Help Import/Export

## Overview
Add a high-contrast, UTF-8 box-drawing tree-view output to the installer and profile commands,
improve help text to include the import/export section, and ensure scheduler & macro commands
also show a structured tree view of their composition during execution and `ls` listing.

## Assets
- Reference images saved to `.lovable/assets/help-ui/`

## Coding Rules (enforced)
- Functions: max 15 lines.
- Booleans: `is`, `has`, `can`, `should` prefix. No negatives.
- Names: strictly semantic. No `data`, `temp`, `obj`, `result`.
- AppError for all error wrapping.
- Box-drawing chars: `├──`, `└──`, `│   `, `    `.
- Color palette (from `constants_terminal.go`): GREEN=success/checkmarks, CYAN=tree branches/items, YELLOW=timestamps/sub-headings, DIM=descriptions/dividers, RED=errors.

## Subtasks

### Task 1: installer_tree.go — Tree-View After Install
Create `gitmap/cmd/installer_tree.go` that:
1. Defines a recursive `InstallerTreeNode` struct with `Title`, `Description`, `Children []InstallerTreeNode`.
2. Implements `printInstallerTree(root InstallerTreeNode, prefix string, isLast bool)` using UTF-8 glyphs.
3. Color rules: `├──` / `└──` in CYAN, title in WHITE, description in DIM.
4. Calls from `executeSmartInstall` (after success) to render the post-install summary.
5. Also renders via `printInstallHistory` (new function) for `install ls` output.

### Task 2: install_profile_tree.go — Profile Hierarchy Tree View
Create `gitmap/cmd/install_profile_tree.go` that:
1. Defines `ProfileComposition` with `Name`, `Alias`, `Description`, `BaseProfile *ProfileComposition`, `Tools []ToolEntry`.
2. Defines `ToolEntry` with `Slug`, `Description`.
3. Implements `resolveProfileTree(name string) (ProfileComposition, bool)` — looks up known Ubuntu profiles.
4. Implements `printProfileTree(p ProfileComposition)` which renders the full hierarchy using box-drawing with colors.
5. Implements `printProfileInstallSummary(slug string)` that calls profile resolution + tree print.
6. Called from `executeSmartInstall` when slug starts with `ubuntu+` or `ubuntu-`.

### Task 3: installer_history_tree.go — History with Tree Expansion
Create `gitmap/cmd/installer_history_tree.go` that:
1. Implements `printInstallerHistoryTree(db *store.DB)` — reads the installation ledger from SQLite.
2. Groups by latest distinct entry, sorts by `MAX(timestamp) DESC`.
3. For each entry: if it is a profile slug → call `resolveProfileTree` and render children. Otherwise: print single `└──` metadata line.
4. Separates entries with a `------` DIM divider.
5. Called from `runInstallerLs` as an opt-in when `--tree` flag is provided.

### Task 4: macro_tree.go — Macro Composition Tree Output
Create `gitmap/cmd/macro_tree.go` that:
1. Implements `printMacroTree(name string)` which loads a macro via `macro.LoadMacro(name)` and renders each step as a tree node.
2. Uses `├──` for all steps except last (`└──`).
3. Calls from `runExecuteCmd` before executing the macro (as a one-time summary header).

### Task 5: schedule_tree.go — Scheduler Tree Output
Create `gitmap/cmd/schedule_tree.go` that:
1. Implements `printScheduleTree(taskName, interval, shellType string, steps []string)` to show a structured summary of what the schedule will run.
2. Called from `gitmap schedule <name>` interactive confirmation step.

### Task 6: Help Text — Import/Export Section
Modify `gitmap/cmd/rootusage*.go` (or `constants_helpgroups.go`) to add a visible Import/Export section to the `gitmap help` output:
```
Import / Export:
  import-export (ie)  Export or import gitmap tracked repos, aliases, and groups
  export              Export tracked repos and settings to a JSON snapshot
  import              Import a JSON snapshot to restore tracking state
```

### Task 7: Release
- Bump MINOR version in `version.json` (6.113.0 -> 6.114.0).
- Update `changelog.md` and `readme.md` badges.

## Verification Checklist
- [ ] `go build ./...` inside `gitmap/` passes cleanly.
- [ ] All functions <= 15 lines.
- [ ] All booleans use `is`/`has`/`can`/`should`.
- [ ] No generic names.
- [ ] Tree output uses correct UTF-8 glyphs + color constants.
- [ ] Import/Export appears in `gitmap help`.
- [ ] No test files modified.
