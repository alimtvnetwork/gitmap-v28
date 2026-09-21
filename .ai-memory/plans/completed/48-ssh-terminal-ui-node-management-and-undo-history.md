# Plan 48: SSH Terminal UI, Node Management, RM/Reset Enhancements, and Undo History [COMPLETED]

## Summary of Completed Work

1. **SSH Help & Zero-Argument Parity:**
   - Unified `gitmap ssh` and `gitmap ssh help` terminal UI.
   - Updated `MsgSSHAvailableCommands` in `cli/constants/constants_ssh.go` with:
     - 4-space left padding for commands, 2-space padding for section titles.
     - Newline gap after `Available SSH subcommands:`.
     - Vibrant ANSI theming (`ColorCyan` title, `ColorWhite` commands, `ColorDim` descriptions, `ColorYellow` examples).
     - Clear examples demonstrating `-y` flags for remove, reset, and undo.
     - Trailing bottom padding.
   - Refactored `checkSSHHelp` into `cli/cmdssh/ssh_help_check.go` so both `gitmap ssh` and `gitmap ssh help` render the identical command list.

2. **Public Key Masking on Console:**
   - Implemented `cli/cmdssh/ssh_key_display.go` with `formatDisplayPublicKey(pubKey string, isRaw bool) string` and `hasRawFlag(args []string) bool`.
   - Public keys on console now display as `ssh-rsa AAAAB3...[redacted, pass --raw to view]... user@host` unless `--raw` is passed.
   - The raw key continues to be pushed to the OS clipboard for seamless pasting to GitHub.
   - Blacked out the exposed `ssh-rsa` key in uploaded screenshot artifacts.

3. **Compare UI Overhaul (`gitmap ssh compare` / `matrix`):**
   - Redesigned `cli/cmdssh/ssh_compare_table.go` from an overwide 192-char table into clean, responsive boxed cards fitting comfortably within 80-120 column terminals.
   - Extracted data providers into `cli/cmdssh/ssh_compare_data.go`, maintaining 100% test compatibility with `ssh_compare_table_test.go`.

4. **Verb Guard & False Join Bug Fix:**
   - Added `isReservedSSHVerb` in `cli/cmdssh/ssh_rm_resolver.go` covering `clear`, `reset`, `rm`, `remove`, `delete`, `ls`, `list`, `nodes`, `node`, `help`, `join`, `sj`, `add`, `enroll`, `keys`, `key`.
   - Protected `ParseSSHTarget` and `resolveSSHTargetFromHost` in `cli/cmdssh/ssh_parser.go` to reject reserved verbs as targets.
   - Intercepted `gitmap ssh join clear` and `gitmap sj clear` in `cli/cmdssh/sshjoin_cmd.go` to route directly to node clearing.
   - Enhanced `cli/cmdssh/ssh_ls_cmd.go` to dispatch `nodes clear`, `nodes rm`, `nodes reset`.

5. **Interactive & Confirmation-Based Removal:**
   - Implemented `cli/cmdssh/sshjoin_rm_cmd.go`, `cli/cmdssh/sshjoin_rm_exec.go`, `cli/cmdssh/ssh_rm_interactive.go`, and `cli/cmdssh/ssh_rm_resolver.go`.
   - Supports:
     - `gitmap ssh rm nodes [target]`
     - `gitmap ssh nodes rm [target]`
     - `gitmap ssh nodes clear`
     - `gitmap ssh join clear`
     - `gitmap ssh rm all`
     - `gitmap ssh remove all nodes`
     - `gitmap ssh rm keys [name]`
     - Bare `gitmap ssh rm`: prompts user interactively (`[1] Nodes`, `[2] Keys`).
   - Resolves node candidates by alias (exact/prefix), IP (exact/prefix), ID (exact/prefix), and user@host.
   - Displays candidate preview table with `RenderSSHHostsTable` and requires `[y/N]` confirmation unless `-y` / `--yes` is supplied.

6. **SSH Reset Command (`gitmap ssh reset [-y]`):**
   - Implemented `cli/cmdssh/ssh_reset_cmd.go` to wipe all registered nodes and flush reachability connection caches.
   - Snapshots all nodes to task history prior to wiping.
   - Requires `-y` / `--yes` or user confirmation.

7. **Task Enqueue, History DB & Undo/Restore Architecture:**
   - Created `cli/cmdssh/ssh_history_types.go`, `cli/cmdssh/ssh_history_db.go`, and `cli/cmdssh/ssh_history_query.go`.
   - Stores node snapshots in `data/history/task/sql.db` table `ssh_task_history`.
   - Implemented `cli/cmdssh/ssh_undo_cmd.go` providing `gitmap ssh undo` (undoes most recent deletion/reset) and `gitmap ssh restore <task-id>`.
   - Restores nodes idempotently with `store.UpsertSSHHost`.

## Verification Outcomes

- `python .github/scripts/go-format-check.py --check-only`: PASS (0 unformatted across 3,168 Go files).
- `python linter-scripts/check-boolean-guidelines.py`: PASS (0 violations).
- `python linter-scripts/check-nested-ifs.py`: PASS (0 violations).
- `python linter-scripts/check-newline-styling.py --all`: PASS (All files Unix LF).
- `python linter-scripts/check-relative-paths.py`: PASS (0 violations).
- Functions <= 15 lines across all new and modified Go files.
- Files <= 100 lines across all new and decomposed Go files.
