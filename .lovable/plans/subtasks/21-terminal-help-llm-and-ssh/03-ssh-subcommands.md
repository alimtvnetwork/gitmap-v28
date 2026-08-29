# Subtask 21.03: SSH Subcommands Expansion in Help & CLI

## Target Files
- `gitmap/constants/constants_ssh.go`
- `gitmap/cmd/ssh.go`
- `gitmap/cmd/rootusage_groups.go`

## Actions
- [ ] Update `MsgSSHAvailableCommands` in `gitmap/constants/constants_ssh.go` to list:
  - `ssh create [name]`
  - `ssh list (ls)`
  - `ssh status (st)`
  - `ssh copy (cp)`
  - `ssh cat (view)`
  - `ssh delete (rm)`
  - `ssh config`
  - `ssh join` (join/broadcast/distribute cluster nodes)
  - `ssh login <user@ip>` (connect & install tools)
  - `ssh alias <alias> <user@ip>` (manage host aliases)
  - `ssh exec <cmd>` (execute command on remote nodes)
- [ ] Ensure `gitmap/cmd/ssh.go` routes `join`, `alias`, `login`, `exec`.

## Acceptance Criteria
- [ ] `gitmap ssh` prints full subcommands list.
