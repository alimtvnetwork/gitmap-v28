# Subtask 01: cmdnodes Core Implementation & Remote File Staging

## Scope
1. Implement `cli/cmdnodes/nodes_clone_types.go`:
   - `NodesCloneKind` enum (`clone`, `cfr`, `cfrp`).
   - `NodesCloneOptions` struct (parsed flags, remaining args, target exclusions, dryRun, json).
   - `RemoteCloneNodeResult` struct.
2. Implement `cli/cmdnodes/nodes_clone_file.go`:
   - Detect if arguments refer to a file or if `gitmap.json` exists locally.
   - Read file bytes and compute target remote path (`D:\work\<filename>` for Windows, `~/work/<filename>` for Unix/Linux).
   - Stage file to remote node over SSH via `cmdssh.StreamFileToRemote`.
3. Implement `cli/cmdnodes/nodes_clone_remote.go`:
   - Filter remote connections (exclude current local machine).
   - Launch concurrent workers per remote node.
   - Execute `gitmap <kind> <args>` on remote node.
   - Format and render results in clean termtable.
4. Implement `cli/cmdnodes/nodes_clone.go`:
   - Master entrypoint `RunNodesClone(args []string) error`.
   - Execute local clone and dispatch remote async workers concurrently.
5. Implement `cli/cmdnodes/nodes_clone_help.go`:
   - Comprehensive help screens for `nodes clone`, `nodes cfr`, `nodes cfrp`.
