# Plan 02 — `gitmap cn -f` Force-Flatten

**Date:** let's start now 2026-04-21 19:30 (UTC+8)
**Target version:** v3.50.0
**Trigger:** User report — `cn v+1 -f` from inside already-flattened `macro-ahk/` should remove cwd and re-clone into `macro-ahk/` (not fall back to `macro-ahk-v22/`).

## Problem

`clone-next` flattens by default into `<base>/` (e.g. `macro-ahk/`). Once a repo is already flat, the user's shell cwd IS the target folder. Windows file-locks the cwd, so `os.RemoveAll(targetPath)` fails. The current code degrades silently to a versioned-folder fallback (`MsgFlattenFallback`), which defeats the whole point of staying flat.

## Solution

Add `-f` / `--force` flag with a strict contract:

1. **Pre-clone chdir.** If `cwd == targetPath`, `os.Chdir(parentDir)` BEFORE the existence check. This releases the file handle on Windows.
2. **No fallback.** When `-f` is set, the locked-folder fallback path is replaced with a hard error + `os.Exit(1)`. User gets either a flat layout or a clear failure — never a surprise versioned folder.
3. **Post-clone chdir-back.** After successful clone, `os.Chdir(targetPath)` so the user is back where they expect, and `GITMAP_SHELL_HANDOFF` carries through.

## Files Touched

| File | Change |
|---|---|
| `spec/01-app/87-clone-next-flatten.md` | Document `-f` flag + use case + interaction table |
| `gitmap/cmd/clonenextflags.go` | Add `Force bool` field, `--force` / `-f` parsing |
| `gitmap/cmd/clonenext.go` | New `forceReleaseLockOnCwd` helper; guard fallback branch on `!Force` |
| `gitmap/constants/constants_clonenext.go` | New msgs: `MsgForceReleasing`, `ErrCloneNextForceFailed`, flag-desc, help-line |
| `gitmap/helptext/clone-next.md` | New row in flags table + new example block |
| `gitmap/completion/zsh.go`, `powershell.go` | Add `-f` / `--force` to hint arrays |
| `gitmap/constants/constants.go` | Bump `Version` to `3.50.0` |
| `.lovable/memory/features/clone-next-flatten.md` | Updated with v3.50.0 force section |

## Out of Scope

- Cross-dir `gitmap cn <repo> v++ -f` — already chdirs into the repo via `tryCrossDirCloneNext`, so the same logic kicks in for free; no extra wiring needed but worth a follow-up smoke test.
- Batch `cn --all -f` — each batch entry runs in its own dir, same code path, no extra wiring.
- Linux/macOS — `os.RemoveAll` on cwd works there, so `-f` is effectively a no-op (but harmless). Still useful as an explicit "I really mean flat" signal.

## Test Strategy

Manual on Windows (per project convention — Windows-first):
1. From `D:\repos\macro-ahk\` (flattened from v21), run `gitmap cn v++ -f`.
2. Expect: `→ Force-flatten: leaving ... to release lock...`, then removal, then clone into `macro-ahk/`, then `→ Now in macro-ahk`.
3. Negative: lock the folder with another process (e.g. `Get-ChildItem` in another shell holding it), run `cn v++ -f` → expect `ErrCloneNextForceFailed` and exit 1, NOT a versioned-folder fallback.

## Open Questions

None.
