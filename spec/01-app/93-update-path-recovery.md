# Update Path Recovery — Source Repo Resolution

## Overview

The `update` command requires access to the gitmap-v28 source repository
to pull changes and rebuild the binary. When the previously saved
path no longer exists on disk (e.g. the user moved or renamed the
folder), the command must recover gracefully rather than failing
with a cryptic `Set-Location` error.

---

## Resolution Strategy (4-Tier Priority)

`resolveRepoPath()` checks each tier in order. The first valid,
existing directory wins.

| Tier | Source                        | Validation             | Persisted |
|------|------------------------------|------------------------|-----------|
| 1    | `--repo-path` CLI flag       | `pathExists()` (none*) | Yes       |
| 2    | Embedded `constants.RepoPath`| `pathExists()`         | No        |
| 3    | SQLite Settings DB           | `pathExists()`         | —         |
| 4    | Interactive user prompt      | `pathExists()`         | Yes       |

\* Tier 1 trusts the user-supplied flag without existence check
because the intent is explicit. The path is saved to the DB for
future runs.

If all four tiers fail, the command delegates to `gitmap-updater`
(if on PATH) or prints the `ErrNoRepoPath` recovery guide.

---

## Interactive Prompt (Tier 4)

When tiers 1–3 produce no valid path:

1. Print a warning to stderr: `MsgUpdatePathMissing`.
2. Print the prompt: `MsgUpdatePathPrompt`.
3. Read a single line from stdin via `bufio.NewReader`.
4. Trim whitespace.
5. Validate the path exists on disk with `pathExists()`.
6. On success: save to DB via `saveRepoPathToDB()` and return.
7. On failure (path does not exist): clone the gitmap-v28 source repo
   into that directory via `cloneRepoInto()`, then re-validate.
8. On successful clone: save to DB and return.
9. On clone failure: print `ErrUpdateCloneFailed` and re-prompt.

---

## Database Storage

Paths are stored in the SQLite Settings table (key-value store)
located in the `data/` directory anchored to the binary's physical
location (resolved via `os.Executable()` + `filepath.EvalSymlinks()`).

| Key                  | Value                              |
|----------------------|------------------------------------|
| `source_repo_path`   | Absolute path to the gitmap-v28 repo  |

This is the same key used by the `release-self` command's
`saveSourceRepoDB()` / `loadSourceRepoDB()` functions, ensuring
both commands share a single source of truth.

---

## File Layout

| File                          | Purpose                                    |
|-------------------------------|--------------------------------------------|
| `cmd/update.go`               | `resolveRepoPath()` with 4-tier dispatch   |
| `cmd/updaterepo.go`           | Path helpers: prompt, DB read/write        |
| `cmd/updatescript.go`         | PowerShell script generation               |
| `constants/constants_update.go`| Path recovery messages and constants      |
| `release/selfrelease_resolve.go`| Shared DB functions for `release-self`   |

---

## Error Handling

| Condition                        | Behavior                                |
|----------------------------------|-----------------------------------------|
| Embedded path missing on disk    | Skip to tier 3 (DB lookup)              |
| DB path missing on disk          | Skip to tier 4 (prompt)                 |
| User enters empty string         | Fall through to updater/error           |
| User enters non-existent path    | Print error, fall through               |
| DB open fails                    | Silently skip (non-fatal)               |
| DB write fails                   | Silently skip (non-fatal)               |

---

## Acceptance Criteria

1. `gitmap-v28 update --repo-path <valid-path>` saves the path to the
   DB and uses it immediately.
2. When the embedded `RepoPath` exists on disk, it is used without
   prompting.
3. When the embedded path is stale but the DB has a valid path,
   the DB path is used without prompting.
4. When both embedded and DB paths are stale, the user is prompted
   interactively.
5. A successfully prompted path is persisted to the DB so
   subsequent runs do not prompt again.
6. Entering a non-existent path at the prompt prints an error and
   falls through to the updater/error guide.
7. The `release-self` and `update` commands share the same
   `source_repo_path` Settings key.
8. All terminal messages use constants from
   `constants_update.go` — no magic strings.