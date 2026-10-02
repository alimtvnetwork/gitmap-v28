# Subtask 04: GitMap Ignore Engine & Hierarchical Groups

## Objectives
- Implement and verify `gitmap ignore` (`ig`) subcommands:
  - `add <pattern>`: add pattern to default group.
  - `ls`: list ignore groups and active bindings.
  - `add-group <name>`: create named ignore group.
  - `remove-group` / `rm-grp <name>`: remove named group.
  - `add-grp-to-default` (`agtd`) `<name>`: chain named group after default group.
  - `connect-group-with-repo` (`cgwp`) `<name> <repo-path>` with `--add-with-default` (`--awd`).
  - `apply <path>`: apply ignore rules to repository.
  - `export` / `import`: JSON configuration backup and restore.
- Enforce immutable default patterns (`.gitmap/`, `.gitmap/backup/`).

## Target Files
- `cli/cmdignore/ignore_cli.go`
- `cli/cmdignore/ignore_groups.go`
- `cli/store/split_db_ignore.go`
