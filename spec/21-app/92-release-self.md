# Spec 60 — Release Self

## Overview

`gitmap-v28 release-self` (alias `rs`, `rself`) provides explicit self-release capability.
It resolves the gitmap-v28 executable's own source repository and performs a full
release workflow from that directory, regardless of the user's current working
directory.

Additionally, `gitmap-v28 release` auto-detects when it is run outside a Git
repository and falls back to self-release mode.

## Commands

### release-self (rs / rself)

    gitmap-v28 release-self [version] [flags]

Explicitly triggers a self-release. Accepts all flags supported by `release`.

### release (auto-fallback)

When `gitmap-v28 release` is invoked outside a Git repository, it automatically
enters self-release mode instead of failing.

## Behavior

### 1. Source Repository Discovery

The command resolves the source repository using a two-tier strategy:

**Strategy 1 — Executable path:**
1. Call `os.Executable()` to get the running binary path.
2. Resolve symlinks via `filepath.EvalSymlinks()`.
3. Walk up the directory tree from the executable's location to find the
   nearest `.git` directory — that directory is the source repo root.
4. On success, persist the resolved path to the `Settings` table
   (`source_repo_path` key) for future fallback.

**Strategy 2 — Database fallback:**
1. If the executable path strategy fails (e.g., binary moved/installed
   outside source tree), read `source_repo_path` from the `Settings` table.
2. Verify the stored path still contains a `.git` root.

If both strategies fail, the command exits with an error:
`could not locate gitmap-v28 source repository`.

### 2. Same-Directory Skip

If the resolved source repo root matches the current working directory,
skip the directory switch entirely. Print:
`→ Self-release: already in source repo /path` and proceed directly.

### 3. Directory Switch

1. Record the caller's working directory via `os.Getwd()`.
2. `os.Chdir()` into the resolved source repo root.
3. Execute the full release workflow (identical to `release`).
4. `os.Chdir()` back to the original working directory.
5. Print `✓ Returned to <original-path>`.

### 4. Flag Passthrough

All flags supported by `release` are accepted:
`--assets`, `--commit`, `--branch`, `--bump`, `--notes`, `--draft`,
`--dry-run`, `--compress`, `--checksums`, `--no-assets`, `--targets`,
`--list-targets`, `--zip-group`, `-Z`, `--bundle`, `--no-commit`,
`--yes` / `-y`, `--verbose`.

### 5. Output

Self-release prints a preamble before the standard release output:

    → Self-release: switching to /path/to/gitmap-source
    <standard release output>
    ✓ Returned to /original/working/directory

## Error Scenarios

| Scenario | Behavior |
|----------|----------|
| Executable path unresolvable + no DB fallback | Exit 1: `could not locate gitmap-v28 source repository` |
| DB path stale (no .git) | Falls through to error |
| Release fails | Standard release error handling (rollback); still returns to original dir |
| Return chdir fails | Warning printed; exit 0 (release succeeded) |

## Implementation

### Package: `release`

New exported function:

```go
func ExecuteSelf(opts Options) error
```

Resolves the source repo, switches directories, calls `Execute(opts)`,
and switches back.

### Package: `cmd`

New file `releaseself.go` with `runReleaseSelf(args)` that reuses
`parseReleaseFlags` and calls `release.ExecuteSelf`.

### Package: `cmd` — release fallback

In `runRelease`, before `requireOnline()`, check `release.IsInsideGitRepo()`.
If false, delegate to `runReleaseSelf(args)` and return.

### Constants

- `CmdReleaseSelf = "release-self"`
- `CmdReleaseSelfAlias = "rself"`
- `CmdReleaseSelfAlias2 = "rs"`
- `SettingSourceRepoPath = "source_repo_path"`
- Messages: `MsgSelfReleaseSwitch`, `MsgSelfReleaseReturn`, `MsgSelfReleaseSameDir`

## Acceptance Criteria

1. `gitmap-v28 rs --bump patch` releases gitmap-v28 itself from any directory.
2. `gitmap-v28 release` outside a Git repo triggers self-release automatically.
3. `gitmap-v28 release` inside a Git repo behaves exactly as before.
4. After self-release, the user is returned to their original directory with confirmation.
5. All release flags work identically in self-release mode.
6. If executable is outside source tree, DB fallback resolves the repo path.
7. If already in the source repo directory, skip the chdir round-trip.
