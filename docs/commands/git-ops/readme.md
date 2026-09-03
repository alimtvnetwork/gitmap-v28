# Git Operations & Pull Engine

Execute Git operations across individual, grouped, or all repositories.

## Commands

### `gitmap pull [name]`
* **Alias:** `p`
* Pulls changes with non-unlink failure diagnostics and multi-repo tabular status display.
* Flags: `--all` (alias: `gitmap pull-all`), `--group <name>`, `--verbose`.

### `gitmap fix [repo] [action]`
* Remediates pull conflicts and dirty working tree collisions (`stash`, `wip`, `discard`).

### `gitmap push`
* **Alias:** `ph`
* Push commits with `--ssh` / `--https` transport flags.

### `gitmap status`
* **Alias:** `st`
* Tabular overview of dirty/clean states, uncommitted files, ahead/behind counts, and stashes.

### `gitmap watch`
* **Alias:** `w`
* Live auto-refreshing terminal dashboard of tracked repository statuses.

### `gitmap exec <args...>`
* **Alias:** `x`
* Executes arbitrary Git commands across all tracked repositories simultaneously.

### `gitmap has-any-updates`
* **Aliases:** `hau`, `hac`
* Checks if remote repositories contain unpulled commits.

### `gitmap latest-branch`
* **Alias:** `lb`
* Identifies the most recently updated remote branch.

### `gitmap lfs-common`
* **Alias:** `lfsc`
* Automatically tracks common binary file types with Git LFS.
