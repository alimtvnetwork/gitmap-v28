# Subtask 23-03: Subcommand Routing & Full Help Discoverability

## Goal

Ensure 100% of nested subcommands across modular commands (`ssh`, `vscode`, `agy`, `group`, `cluster`) are listed in root help text, routed through dispatchers, and include usage examples:
- `gitmap ssh`: `create`, `list`, `status`, `copy`, `cat`, `delete`, `config`, `join`, `login`, `alias`, `exec`
- `gitmap vscode`: `add`, `rm`, `ls`, `pap`, `plugins`
- `gitmap agy`: `add`, `rm`, `ls`, `stats`, `update`

## Status: DONE

- Verified subcommand dispatchers in `gitmap/cmd/ssh.go`, `gitmap/cmd/vscode_cmd.go`, `gitmap/cmd/agy_cmd.go`.
