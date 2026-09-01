# Clone Next

## Status

Implemented (v2.75.0). Flatten-by-default behavior active.

## Command

```text
gitmap-v28 clone-next <v++|v+1|vN> [flags]
```

## Alias

```text
cn
```

## Responsibility

From inside an existing Git repository, derive the source repository from
`remote.origin.url`, resolve the next or explicit versioned target repository,
clone it into the parent directory using the **base name folder** (version
suffix stripped), register it with GitHub Desktop, record the version
transition in the database, and optionally remove the current local folder.
If `--create-remote` is passed, the command will also create the target
GitHub repository before cloning when it does not exist.

## Flatten-by-Default Behavior

Starting from v2.75.0, `clone-next` always flattens:

1. The target clone folder is the **base name** without version suffix
   (e.g., `macro-ahk` instead of `macro-ahk-v16`).
2. If the base name folder already exists, it is **removed entirely** and
   re-cloned fresh (no prompt).
3. The remote URL still points to the versioned repo (e.g., `macro-ahk-v16`).
4. Version columns (`CurrentVersionTag`, `CurrentVersionNum`) are updated
   on the `Repos` row.
5. A `RepoVersionHistory` row is inserted tracking the transition.
6. `GITMAP_SHELL_HANDOFF` is set to the flattened path.

The previous `--flatten` flag is no longer required — this is the default.

## Source of Truth

The command must use the Git remote as the authoritative source for repo name
resolution.

1. Read the current repo URL using:
   - `git config --get remote.origin.url`
   - fallback: `.git/config` origin parsing if needed
2. Parse the host, owner/org, and repo name from that remote URL.
3. Derive the base repo name and current version from the **remote repo name**.
4. Use the current local folder name only for:
   - determining the parent directory for the clone
   - prompting/removing the current folder after success

This avoids incorrect behavior when the local folder name and remote repo name
are not perfectly aligned.

## Terminology

| Term | Meaning |
|------|---------|
| Base name | Repo name without version suffix, e.g. `macro-ahk` |
| Current version | Version implied by the current remote repo name |
| Target version | Version requested by the user |
| Target repo | New repo name in the form `<base>-vN` |
| Flattened folder | Local folder using base name only |

## Version Arguments

| Argument | Meaning | Example |
|----------|---------|---------|
| `v++` | Increment current version by 1 | `macro-ahk-v11` → `macro-ahk-v12` |
| `v+1` | Alias for increment-by-one | `coding-guidelines-v7` → `coding-guidelines-v8` |
| `vN` | Jump directly to an explicit version | `macro-ahk-v12` + `v15` → `macro-ahk-v15` |

## Version Rules

1. `v++` and `v+1` mean the same thing.
2. `vN` must accept only positive integers (`v1`, `v2`, `v15`, ...).
3. `v0`, negative values, and malformed inputs must fail with a clear error.
4. If the current repo has no suffix, the unsuffixed repo is treated as the
   original repo and the first increment target is `-v2`.

### No-Suffix Behavior

| Current repo | Argument | Target repo | Local folder |
|--------------|----------|-------------|--------------|
| `macro-ahk` | `v++` | `macro-ahk-v2` | `macro-ahk/` |
| `macro-ahk` | `v+1` | `macro-ahk-v2` | `macro-ahk/` |
| `macro-ahk` | `v15` | `macro-ahk-v15` | `macro-ahk/` |

## Target Resolution

After parsing the current remote:

1. Compute the target version.
2. Build the target repo name: `<base-name>-v<target-version>`.
3. Build the target local path: `<parent-directory>/<base-name>` (flattened).
4. Build the target remote URL by preserving the same host, owner/org, and URL
   scheme as the current remote.

### URL Examples

| Current remote | Target remote | Local folder |
|----------------|---------------|--------------|
| `https://github.com/alimtvnetwork/macro-ahk-v11.git` | `https://github.com/alimtvnetwork/macro-ahk-v12.git` | `macro-ahk/` |
| `git@github.com:alimtvnetwork/macro-ahk-v11.git` | `git@github.com:alimtvnetwork/macro-ahk-v12.git` | `macro-ahk/` |

## Optional GitHub Creation (`--create-remote`)

By default, `clone-next` assumes the target remote already exists and proceeds
directly to `git clone`. When the `--create-remote` flag is set, the command
checks whether the target GitHub repository exists and creates it if missing
**before** attempting to clone. This requires `GITHUB_TOKEN` to be set.

### Behavior when `--create-remote` is set

1. Check whether the target remote repository exists.
2. If it does not exist and the host is GitHub, create it under the same
   owner/org as the source repo.
3. The created repo should use the target repo name exactly.
4. The command must not attempt `git clone` first when the target repo is known
   to be missing.
5. If repo creation fails, stop with a clear error and do not prompt for local
   deletion.

### Visibility (when creating)

The preferred behavior is to inherit the visibility of the source repository.
If that cannot be determined safely, the command should fail with a clear error
instead of guessing.

## Workflow

1. Confirm the current directory is a Git repo.
2. Resolve `remote.origin.url`.
3. Parse the current remote repo name and current version.
4. Resolve the target version from `v++`, `v+1`, or `vN`.
5. Compute the target repo name and flattened local path (base name only).
6. If the flattened folder already exists, remove it entirely.
7. If `--create-remote` is set, check whether the target remote exists and
   create it if missing.
8. Clone the target repo into the flattened folder.
9. Record version transition in database (`RepoVersionHistory`).
10. Update `Repos` row with `CurrentVersionTag` and `CurrentVersionNum`.
11. Register the cloned repo with GitHub Desktop unless `--no-desktop` is set.
12. If `--delete` is set and the current folder differs from the flattened path,
    remove the current versioned folder.
13. Set `GITMAP_SHELL_HANDOFF` to the flattened path.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--delete` | false | Remove the current versioned folder after clone (when different from flattened path) |
| `--keep` | false | Keep the current folder and skip the removal prompt |
| `--no-desktop` | false | Skip GitHub Desktop registration |
| `--create-remote` | false | Create the target GitHub repo if it does not exist (requires `GITHUB_TOKEN`) |
| `--ssh-key <name>` / `-K <name>` | (none) | Use a named SSH key for Git operations |
| `--verbose` | false | Show detailed clone-next diagnostics |

## Examples

### Example 1: Simple clone with `v++` (flattened)

```text
D:\wp-work\riseup-asia\macro-ahk-v11> gitmap-v28 cn v++

Removing existing macro-ahk for fresh clone...
Cloning macro-ahk-v12 into macro-ahk (flattened)...
✓ Cloned macro-ahk-v12 into macro-ahk
✓ Recorded version transition v11 -> v12
✓ Registered macro-ahk-v12 with GitHub Desktop
```

### Example 2: Jump to a specific version with auto-delete

```text
D:\wp-work\riseup-asia\macro-ahk-v12> gitmap-v28 cn v15 --delete

Cloning macro-ahk-v15 into macro-ahk (flattened)...
✓ Cloned macro-ahk-v15 into macro-ahk
✓ Recorded version transition v12 -> v15
✓ Registered macro-ahk-v15 with GitHub Desktop
✓ Removed macro-ahk-v12
```

### Example 3: Repo without an existing suffix

```text
D:\wp-work\riseup-asia\macro-ahk> gitmap-v28 cn v++

Removing existing macro-ahk for fresh clone...
Cloning macro-ahk-v2 into macro-ahk (flattened)...
✓ Cloned macro-ahk-v2 into macro-ahk
✓ Recorded version transition v1 -> v2
✓ Registered macro-ahk-v2 with GitHub Desktop
```

### Example 4: Create remote repo before clone

```text
D:\wp-work\riseup-asia\macro-ahk-v12> gitmap-v28 cn v15 --create-remote --delete

Creating GitHub repo macro-ahk-v15...
✓ Created GitHub repo macro-ahk-v15
Cloning macro-ahk-v15 into macro-ahk (flattened)...
✓ Cloned macro-ahk-v15 into macro-ahk
✓ Recorded version transition v12 -> v15
✓ Registered macro-ahk-v15 with GitHub Desktop
✓ Removed macro-ahk-v12
```

### Example 5: Lock detection during folder removal

```text
D:\wp-work\riseup-asia\macro-ahk-v11> gitmap-v28 cn v++ --delete

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
```

## Error Handling

| Condition | Required behavior |
|-----------|-------------------|
| Not inside a Git repo | Print a clear error and exit 1 |
| `remote.origin.url` missing | Print a clear error and exit 1 |
| Remote URL cannot be parsed | Print a clear error and exit 1 |
| Invalid version argument | Print a clear error and exit 1 |
| Flattened folder removal fails | Print error and exit 1 |
| Target GitHub repo creation fails (`--create-remote`) | Print a clear error and stop before clone |
| Clone fails | Print a clear error and do not update DB |
| DB version tracking fails | Warn to stderr, do not exit (clone succeeded) |
| GitHub Desktop registration fails | Warn, but keep clone success |
| Old folder deletion fails (no locks found) | Warn, but keep clone success |
| Old folder deletion fails (locks found) | List locking processes, prompt to kill, retry removal |

## Implementation Scope

| Component | File |
|-----------|------|
| Command handler | `cmd/clonenext.go` |
| Version history recording | `cmd/clonenexthistory.go` |
| Flag parser | `cmd/clonenextflags.go` |
| Version parser | `clonenext/version.go` |
| Lock detection (shared) | `lockcheck/lockcheck.go` |
| Lock detection (Windows) | `lockcheck/lockcheck_windows.go` |
| Lock detection (Unix) | `lockcheck/lockcheck_unix.go` |
| Version history store | `store/version_history.go` |
| Version history model | `model/version_history.go` |
| Constants | `constants/constants_clonenext.go`, `constants/constants_version_history.go` |
| Help text | `helptext/clone-next.md` |
| Spec | `spec/01-app/59-clone-next.md` |

## Acceptance Criteria

1. `gitmap-v28 cn v++` clones into the base name folder (flattened) by default.
2. `gitmap-v28 cn v+1` behaves exactly like `v++`.
3. `gitmap-v28 cn v15` clones the exact target version into the flattened folder.
4. The source repo name is derived from the Git remote, not the local folder.
5. If the flattened folder exists, it is removed and re-cloned fresh.
6. `Repos.CurrentVersionTag` and `CurrentVersionNum` are updated after clone.
7. A `RepoVersionHistory` row is inserted for each transition.
8. `--create-remote` creates missing target GitHub repos before clone.
9. GitHub Desktop registration happens by default after a successful clone.
10. `--delete` removes the old versioned folder when it differs from the flattened path.
11. `GITMAP_SHELL_HANDOFF` is set to the flattened path.
12. Lock detection works cross-platform for folder removal.

## Deferred Implementation Phases

1. ~~Version parsing and resolution fixes~~ — done
2. ~~Target GitHub repo existence check and creation~~ — done
3. ~~Clone workflow hardening~~ — done
4. ~~Lock detection and process termination~~ — done (v2.52.0)
5. ~~Flatten-by-default with version tracking~~ — done (v2.75.0)
6. Help, completion, and automated test updates

## See Also

- [Clone-Next Flatten](87-clone-next-flatten.md) — Flatten mode DB schema and version tracking
- [Cloner](05-cloner.md) — File-based and direct URL clone behavior
- [Clone Direct URL](88-clone-direct-url.md) — Single repo clone from Git URL
- [Clone Progress](34-clone-progress.md) — Progress tracking for batch clone operations
