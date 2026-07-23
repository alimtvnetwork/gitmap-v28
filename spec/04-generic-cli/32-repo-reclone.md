# Repo-Reclone (single-repo wipe + re-clone)

> Overlays `gitmap reclone` / `rc` / `rec` / `relclone` / `clone-now`.
> Manifest-mode behavior (spec §11-batch-execution) is unchanged.

## Trigger shape

The overlay handler `tryRunRepoReclone` runs **before** manifest parsing
in `runCloneNow`. It owns the command when, AFTER stripping the
`-y` / `--y` / `-yes` / `--yes` flag, the remaining positionals match
one of:

| Positionals | Match? | Target |
|-------------|--------|--------|
| `[]` and cwd has `.git/` | yes | cwd |
| `[<path>]` and `<path>/.git/` exists | yes | abs(`<path>`) |
| anything else | no — fall through to manifest | — |

Manifest paths (`.gitmap/output/gitmap.json`, `.csv`) are NEVER git
repos, so the manifest pipeline is never starved.

## Flow

1. Read `remote.origin.url` from the target via `currentOriginURL`.
   Empty / error → exit 1 (`ErrRepoRecloneNoOrigin`).
2. Print plan line (`MsgRepoReclonePlan`).
3. Confirm — unless `-y`:
   - TTY → prompt `y/N` (`MsgRepoRecloneConfirm`). Anything but `y` →
     exit 1 (`MsgRepoRecloneAborted`).
   - Non-TTY → refuse (`ErrRepoRecloneNonTTY`, exit 1).
4. Release Windows cwd handle (`escapeCwdIfInside`) before delete.
5. `os.RemoveAll(target)` — fatal on failure (`ErrRepoRecloneRemove`).
6. `runCloneCommand(origin, dest)` into the same parent — fatal on
   failure (`ErrRepoRecloneClone`).
7. `WriteShellHandoff(dest)` so the wrapper cd's back in.

## Why overlay instead of new verb

`rc`, `rec`, `reclone`, `relclone`, `clone-now` already route to
`runCloneNow` (see `gitmap/cmd/rootcore.go`). Adding a parallel verb
would either collide on `rc` (breaking existing manifest scripts) or
require yet another name. The overlay is shape-detected and falls
through cleanly, so manifest users see no behavior change.

## Exit codes

| Code | Reason |
|------|--------|
| 0 | Re-clone succeeded |
| 1 | Aborted, missing origin, remove failure, clone failure, or non-TTY without `-y` |

## Files

- `gitmap/cmd/reporeclone.go`
- `gitmap/constants/constants_reporeclone.go`
- `gitmap/cmd/clonenow.go` (interception)
