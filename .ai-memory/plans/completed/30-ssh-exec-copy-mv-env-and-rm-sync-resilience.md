# 30 — SSH Exec Polish, Remote Copy/Move, Environment Management & CLI Resilience Suite

- **Slug:** ssh-exec-copy-mv-env-and-rm-sync-resilience
- **Date Completed:** 2026-09-19
- **Status:** completed
- **Budget:** N=250 steps

---

## User Request (Verbatim)

```text
gitmap ssh exec ip
gitmap ssh exec cmd1,cmd2, cmd3 --except machine-name, alias, id # run commands on all nodes if the machine is open if not then it will say first run these mahciens are off and output by running code into thos emachiens

gitmap pipeline errors clear # should clear as well
gitmap ssh copy from-path to-path
gitmap ssh copy help
gitmap ssh copy host-machine-path-or-filein-currentlocation .\test # it would bascially copy file in to the default workdir for all the machines

gitmap ssh copy host-machine-path-or-filein-currentlocation .\test --except alias,.. # it would bascially copy file in to the default workdir for all the machines
gitmap ssh mv host-machine-path-or-filein-currentlocation-or-machine-we-are-runnig-from .\test --except alias,.. # it would bascially move file in to the default workdir for all the machines or specific if full path is given also all paths should respect the `~` meaning from user home directory, %win%, %win-drive%, %temp%, %appdata%

gitmap env add/rm/ls/help

All right. When we are running the SSH execute or Trust Store execute or direct execute as well. So in all cases, it would run the commands on all those machines parallelly. So first it would nicely, uh, show so the UI is not there. So that's the first thing we have to fix. Uh, it should mention like these are the commands we are running on these, these places, things like that. Okay. Um, once the UI commands and things are running, okay, um, it'll say like, "These commands are injected," and once the, uh, running is finished, it'll come back and showcase the result for each one of the machines, and whichever is, uh, shut off, it would let us know at the end, like these machines off. Also starting time, it would let us know that these are the machines that these are off actually. So this needs to fix. Also, the Git map sync command is, uh, not there, so I think there is an issue with Git map sync. Also, there is an issue if the, if the... It's a package or file or folder is not there in the Git map, but it should also be able to remove. Remember that, should fix that code. Uh, packagers should also showcase a command to clear the database if we wanted to. Okay, so the command running, um, with the Git map SSH execute, it does not show the commands very nicely. The padding in the left side is not there. The-- there is no padding, uh, top and bottom, so make sure that these are there. And first it would say like, "These commands are injected for, for these, uh, IPs and machines." Okay, nicely. No need to mention like no password, no key. No need to mention. So nicely put the padding and mention the commands are running or mention like machines are off. Uh, and then it should give a progress, something like this or processing. And once the processing gets back, it would just return back the results. Okay, so this is how, uh, it should work. And also we need to have Git map copy command. So the... Actually it should be SSH copy. So basically we can just put a, put a path from, uh, let's say from path, uh, to toPath. This would be one example. Also, we need the help for this as well. So that means it should have a help with example. Um, also there could be example like, uh, if I'm doing from the host machine. Host machine for a file location. I select that. Then if I just let's say slash, uh, let's say test, it would basically, uh, basically copy the file into the remote work directory, dir for all the machines. Okay. Um, we, we could also do a accept command where we could just do accept, let's say flag, where we could just give the aliasing and other stuff to reduce that host issue. So we could do copy, we could do, uh, let's say move as well. So move out as well. It could be from host machine or the machine that we are running from. Um And we can, we can basically move that file to that location or directory or a specific, uh, if path is given. Also all paths should respect, respect the tilde. That means from the home directory. Uh, there could be also some common directories like, uh, win for the Windows folder or, uh, let's say we have drive for Windows drive. It should have like temp, it should have the app data, it should have the expansion behavior as well. It should deal with the environment variable to expand the variable. Remember that. This is very, very important. All the paths should actually go through this, uh, things. Uh, also we could actually do Git map, uh, env, env add remove. That would really export the environment variables. So that is must for regardless of the OS. Ubuntu and Windows should be must here now
```

---

## Completed Subtasks

### 1. Subtask 01: SSH Exec UI Polish, Visual Padding, Injection Banner & Offline Detection
- **Files:** `cli/cmdssh/sshexec.go`, `cli/cmdssh/ssh_exec_ui.go`, `cli/cmdssh/ssh_exec_command.go`, `cli/cmdssh/ssh_exec_ui_test.go`
- **Accomplishments:**
  - Implemented pre-flight liveness check (`partitionOnlineOffline`) using TCP dial timeouts to detect offline machines upfront without blocking or throwing noisy credential errors.
  - Rendered a top-padded command injection banner announcing commands to be executed and target machine labels/IPs.
  - Suppressed noisy `"No password or key configured"` messages.
  - Left-padded output execution lines with clear node headers.
  - Summarized all offline nodes at execution completion (`"These machines are off: ..."`).
  - Supported comma-separated commands (`cmd1,cmd2,cmd3` normalized into `cmd1 && cmd2 && cmd3` or `; ` on Windows).
  - Supported `--except` and `--exclude` filters matching ID, alias, host, and user@host.

### 2. Subtask 02: SSH Remote Copy, Move Engine & Universal Path Macro Expansion
- **Files:** `cli/cmdssh/ssh_transfer.go`, `cli/cmdssh/ssh_path_expand.go`, `cli/cmdssh/ssh.go`, `cli/cmdssh/exports.go`, `cli/cmdssh/ssh_path_expand_test.go`, `cli/helptext/ssh.md`
- **Accomplishments:**
  - Implemented `gitmap ssh copy` (`cp`) and `gitmap ssh mv` (`move`) using pure SSH base64 chunked transfer engine.
  - Added relative path resolution defaulting to remote workdirs (e.g. `.\test`).
  - Supported `--except` and `--exclude` machine filtering.
  - Implemented universal cross-platform macro expansion: `~` (home directory), `%win%` (`C:\Windows`), `%win-drive%` (`C:`), `%temp%`, `%appdata%`, and environment variables (`$VAR`, `%VAR%`).
  - Added rich interactive help documentation and practical examples for `gitmap ssh copy help`.

### 3. Subtask 03: Cross-Platform Environment Management & Pipeline Error Clear Command
- **Files:** `cli/cmd/env.go`, `cli/cmdpipeline/pipeline.go`
- **Accomplishments:**
  - Enhanced `gitmap env` to support `add` (alias for set), `rm` / `remove` (alias for delete), `ls` / `list` (alias for list), and `help`.
  - Maintained cross-platform environment persistence and exporting across Ubuntu Linux and Windows.
  - Routed `gitmap pipeline errors clear` and `clear-errors` directly to purge database error tables, logs, and cache without attempting remote log downloads.

### 4. Subtask 04: Resilient Removal (`gitmap rm`), Windows File Lock Protection & Sync Guidance
- **Files:** `cli/cmd/rm.go`, `cli/cmd/rm_resilience.go`, `cli/cmd/rm_resilience_test.go`, `cli/cmd/sync.go`, `cli/cmd/sync_help_test.go`, `cli/helptext/rm.md`, `cli/helptext/sync.md`
- **Accomplishments:**
  - Implemented `safeRemoveWithRetry` and `isProcessLockError` in `cli/cmd/rm_resilience.go`.
  - When Windows file locks (`The process cannot access the file because it is being used by another process` / `unlinkat`) occur, retried with backoff and reported a clean warning, untracking the repo from DB, JSON, and desktop configurations without throwing raw fatal AppError stack traces.
  - Cleaned up missing package handling: if a folder is already missing on disk, cleanly untracks it. If no repo matches, outputs clear guidance and suggestions without crashing.
  - Fixed `gitmap sync` without arguments or with help flags to print usage and exit cleanly with code 0 instead of failing with `E1144` validation error.

---

## Verification & Test Inventory

- `cli/cmdssh/ssh_exec_ui_test.go`: Verified comma-separated multi-command normalization and `--except` exclusion logic.
- `cli/cmdssh/ssh_path_expand_test.go`: Verified path macro expansion for `~`, `%win%`, `%win-drive%`, `%temp%`, `%appdata%`, and env vars.
- `cli/cmd/rm_resilience_test.go`: Verified detection of Windows process lock errors and file access errors.
- `cli/cmd/sync_help_test.go`: Verified `isSyncHelp` tokens.

---

## Coding Guidelines Compliance

- Functions <= 8–15 lines.
- Affirmative booleans (`is*`, `has*`).
- Single return types (zero error tuples).
- Universal `*apperror.AppError` envelopes.
- Unix LF line endings.
