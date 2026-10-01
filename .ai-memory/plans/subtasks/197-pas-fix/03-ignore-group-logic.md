# Subtask 03: Ignore Group Logic

**Objective**: Implement the new `gitmap ignore` group management commands specified in Spec 197.

## Requirements
- Create commands in `cli/cmdignore/`:
  - `add <string>`: Adds ignore rule to default group.
  - `ls`: List groups, rule count, application rules.
  - `add-group` / `remove-group (rm-grp)`: Manage ignore groups.
  - `set-default-group` / `add-grp-to-default (agtd)`: Set default execution groups.
  - `connect-group-with-repo (cgwp)`: Maps a group to a repo alias/path.
  - `apply`: Applies ignore rules to the current folder/repo.
  - `export` / `import`: JSON-based backup/restore of settings.
- Implement storage in `cli/store/split_db_ignore.go` or similar.

## Constraints
- Bounding Box: `cli/cmdignore/*.go`, `cli/store/split_db_ignore*.go`
- Coding Rules: Positive booleans only, `*appfault.AppError`, functions <= 15 lines. No builds/tests.
