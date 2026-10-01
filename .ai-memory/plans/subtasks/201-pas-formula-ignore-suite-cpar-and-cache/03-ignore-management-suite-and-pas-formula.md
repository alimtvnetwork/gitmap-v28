# Subtask 03: Ignore Management Suite & PAS Fleet Extensions

## Scope
- Files:
  - `cli/cmdignore/` (new/updated command package for `gitmap ignore` / `gitmap ig`)
  - `cli/cmd/rootgit.go`, `cli/cmd/rootcore.go`
  - `cli/cmdnodes/`, `cli/cmdssh/`
- Objective:
  - Implement `gitmap fix ignore all [-y]` (`fia`), `gitmap fix-ignore-all [-y]`.
  - Implement `gitmap fix ignores all ssh [-y]`, `gitmap fix-ignores-all-ssh (fias) [-y]` following GitMap PAS Formula.
  - Implement `gitmap ignore` (`ig`) suite:
    - `add`, `scan`, `scan-ssh` (`ss`), `remove`, `edit`, `action`, `ls`, `help`, `ui`, `app`.
    - `add-group`, `remove-group` (`rm-grp`), `set-default-group`, `add-grp-to-default` (`agtd`), `apply`.
    - `connect-group-with-repo` (`cgwp`) <group-name> <path1,alias> `--add-with-default` (`awd`).
    - `export`, `import`.
