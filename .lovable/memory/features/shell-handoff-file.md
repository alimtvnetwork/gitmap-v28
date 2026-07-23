---
name: Shell Handoff File Mechanism
description: GITMAP_HANDOFF_FILE sentinel-file pattern lets clone-next, as, and cd hand a target dir back to the parent shell wrapper. Replaces broken GITMAP_SHELL_HANDOFF env var (v3.103.0).
type: feature
---

# Shell Handoff via Sentinel File (v3.103.0)

## Problem

Previously `clone-next` did `os.Setenv("GITMAP_SHELL_HANDOFF", path)` to
ask the parent shell to cd into the new flattened folder. **This was a
no-op**: a child process cannot mutate the parent shell's environment.
The line was documented as the contract but never functional.

## Solution: sentinel-file pattern

The wrapper function (defined in `constants.CDFuncBash` / `CDFuncZsh` /
`CDFuncPowerShell`) now exports `GITMAP_HANDOFF_FILE=<temp>` before
invoking the real binary, then reads the file after the binary exits
and `cd`s the parent shell to the recorded path.

### Flow

1. User runs `gitmap clone-next foo` (or `as`, or `cd`).
2. Wrapper function:
   - `mktemp` → `$handoff` (PowerShell uses `[Path]::Combine` + GUID)
   - Sets `GITMAP_HANDOFF_FILE=$handoff` and `GITMAP_WRAPPER=1`
   - Invokes `command gitmap "$@"`
3. Binary calls `cmd.WriteShellHandoff(targetPath)`:
   - Reads `EnvGitmapHandoffFile` from env
   - If unset → no-op (legacy behaviour preserved)
   - If set → writes `targetPath` verbatim to that file
4. Wrapper after binary exits:
   - If `$handoff` is non-empty → `cd` to its contents
   - `rm -f $handoff` cleanup

### Wired commands

| Command | Path written | File |
|---------|--------------|------|
| `clone-next` | `targetPath` (flattened folder) | `gitmap/cmd/clonenext.go` |
| `as` | repo top-level (`gitTopLevel()`) | `gitmap/cmd/as.go` |
| `cd <name>` | resolved repo path | `gitmap/cmd/cdops.go::runCDLookup` |
| `cd repos` | picked repo path | `gitmap/cmd/cdops.go::runCDRepos` |

`cd` already used a stdout-capture mechanism via the wrapper. It now
**also** writes the handoff file for parity, so any future wrapper
upgrade can drop the stdout dance.

### Backwards compatibility

- Without the wrapper installed, `GITMAP_HANDOFF_FILE` is unset →
  `WriteShellHandoff` is a silent no-op.
- The legacy `GITMAP_WRAPPER=1` detector and `cd` stdout protocol both
  remain intact.

## Files

- `gitmap/cmd/shellhandoff.go` — `WriteShellHandoff(path)` helper
- `gitmap/cmd/shellhandoff_test.go` — 3 unit tests (no-op / writes / empty)
- `gitmap/constants/constants_cd.go` — `EnvGitmapHandoffFile` constant + updated `CDFunc*` wrappers
- `gitmap/cmd/clonenext.go` — replaced broken `os.Setenv` line
- `gitmap/cmd/as.go` — added handoff after `registerAlias`
- `gitmap/cmd/cdops.go` — added handoff in `runCDLookup` + `runCDRepos`

## Why not extend to `update`?

`update` does not navigate to a new directory — it self-replaces the
binary and leaves the user's cwd untouched. A handoff would have no
sensible target.
