# gitmap clone-next

Clone the next or a specific versioned iteration of the current repository into the parent directory, using the base name (no version suffix) as the local folder.

## Alias

cn

## Usage

    gitmap clone-next <v++|vN> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --force, -f | false | Force flatten when cwd IS the target folder (chdir to parent first; refuses versioned-folder fallback) |
| --delete | false | Auto-remove current versioned folder after clone |
| --keep | false | Keep current folder without prompting |
| --no-desktop | false | Skip GitHub Desktop registration |
| --ssh-key \<name\> | (none) | Use a named SSH key for the clone |
| --verbose | false | Write detailed debug log |
| --create-remote | false | Create target GitHub repo if missing (requires GITHUB_TOKEN) |
| --csv \<path\> | (none) | Batch mode: read repo paths from CSV (one per row) |
| --all | false | Batch mode: cn every git repo one level under cwd |
| --max-concurrency N | 0 (auto = NumCPU) | Batch mode: run up to N repos in parallel. `1` = sequential, `0` = auto. |
| --output \<mode\> | (off) | `terminal` prints the standardized branch/from/to/command block on **stdout** before the clone runs (one block per repo in batch mode). Git progress and version-transition lines stay on **stderr**. |
| --no-vscode-sync | false | Skip syncing the freshly-cloned (flattened) folder into VS Code Project Manager `projects.json`. Default is to sync after every successful `cn`. In batch mode (`--csv` / `--all`), the flag suppresses sync for every repo in the batch. |

`projects.json` `rootPath` values are canonicalized via `filepath.Clean` +
`filepath.EvalSymlinks` (resolves Windows 8.3 short names + symlinks) so
the same physical folder cannot produce duplicate sidebar entries. When
`EvalSymlinks` errors (path not yet on disk, broken symlink), the cleaned
absolute path is used and the clone still succeeds — soft-fail policy.
See `gitmap clone --help` ("Windows path canonicalization & EvalSymlinks
soft-fail") for the full rule set.

## Prerequisites

- Must be run inside a Git repository with a remote origin configured

## Flatten Behavior

By default, clone-next clones into the base name folder (without version suffix).
For example, running `gitmap cn v++` inside `macro-ahk-v11` will:
1. Clone `macro-ahk-v12` into `macro-ahk/` (not `macro-ahk-v12/`)
2. If `macro-ahk/` already exists, remove it first
3. The remote URL still points to `macro-ahk-v12` on GitHub
4. Record the version transition (v11 -> v12) in the database

## Examples

### Example 1: Increment version by one

    gitmap cn v++

**Output:**

    Removing existing macro-ahk for fresh clone...
    Cloning macro-ahk-v12 into macro-ahk (flattened)...
    ✓ Cloned macro-ahk-v12 into macro-ahk
    ✓ Recorded version transition v11 -> v12
    ✓ Registered macro-ahk-v12 with GitHub Desktop

### Example 2: Jump to a specific version with auto-delete

    gitmap cn v15 --delete

**Output:**

    Cloning macro-ahk-v15 into macro-ahk (flattened)...
    ✓ Cloned macro-ahk-v15 into macro-ahk
    ✓ Recorded version transition v12 -> v15
    ✓ Registered macro-ahk-v15 with GitHub Desktop
    ✓ Removed macro-ahk-v12

### Example 3: Lock detection when folder is in use

    gitmap cn v++ --delete

**Output:**

    Removing existing macro-ahk for fresh clone...
    Cloning macro-ahk-v12 into macro-ahk (flattened)...
    ✓ Cloned macro-ahk-v12 into macro-ahk
    ✓ Recorded version transition v11 -> v12
    ✓ Registered macro-ahk-v12 with GitHub Desktop
    Warning: could not remove macro-ahk-v11: unlinkat: access denied
    Checking for processes locking macro-ahk-v11...
    The following processes are using this folder:
      • Code.exe (PID 14320)
      • explorer.exe (PID 5928)
    Terminate these processes to allow deletion? [y/N] y
    Terminating Code.exe (PID 14320)...
    ✓ Terminated Code.exe
    Terminating explorer.exe (PID 5928)...
    ✓ Terminated explorer.exe
    Retrying folder removal...
    ✓ Removed macro-ahk-v11

### Example 4: Force-flatten from inside the already-flat folder

You're working in `D:\repos\macro-ahk\` (flattened from v21) and want
to bump to v22 without ending up in `macro-ahk-v22/`.

    gitmap cn v++ -f

**Output:**

    → Force-flatten: leaving D:\repos\macro-ahk to release lock...
    Removing existing macro-ahk for fresh clone...
    Cloning macro-ahk-v22 into macro-ahk (flattened)...
    ✓ Cloned macro-ahk-v22 into macro-ahk
    ✓ Recorded version transition v21 -> v22
    ✓ Registered macro-ahk-v22 with GitHub Desktop
    → Now in macro-ahk

If `-f` is omitted in this scenario, gitmap falls back to creating
`macro-ahk-v22/` (because the shell holds a file lock on the cwd) and
prints a hint to use `-f` next time.

If even `-f` cannot remove the folder (some other process holds a
handle), gitmap aborts with a clear error rather than silently
falling back to a versioned folder name.

## See Also


- [clone](clone.md) — Clone repos from output files
- [desktop-sync](desktop-sync.md) — Sync repos to GitHub Desktop
- [ssh](ssh.md) — Manage named SSH keys

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter clone-next
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
