# Audit — Reclone-class URL pickers (2026-06-07)

**Scope:** Every command that re-runs `git clone` against an already-known repo. Answers the user's question: "which of CFR / CFRP has SSH-aware reclone, and which doesn't?"

## Verdict table

- `repo-reclone` / `rc` / `rec` / `relclone` / `clone-now` (folder-or-cwd shape — `gitmap/cmd/reporeclone.go:106`) — **HONORS transport.** Reads `currentOriginURL(target)` via `git config --get remote.origin.url` and passes that exact URL to `runCloneCommand` (`reporeclone.go:107,134`). If the origin is `git@…`, the reclone runs over SSH — no HTTPS fallback ever fires. This is the only reclone path that already does the right thing end-to-end for an existing folder.
- `cfr` — `clone-fix-repo` (`gitmap/cmd/clonefixrepo.go:33`) — **PARTIAL.** Honors the user-supplied URL transport verbatim, plus respects `--ssh` / `--https` rewrites via `applyCloneFixRepoScheme` (`clonefixrepo.go:79`). Does NOT auto-detect transport from an existing folder because the command is URL-driven (not folder-driven). When the user pastes an HTTPS URL it clones over HTTPS even if the destination already exists with an SSH origin. Fix: when the destination folder already exists with `.git/`, read `remote.origin.url` and prefer its scheme over the positional URL (or warn on mismatch).
- `cfrp` — `clone-fix-repo-pub` (`gitmap/cmd/clonefixrepo.go:39`) — **PARTIAL, same as cfr.** Shares `runCloneFixRepoPipeline` (`clonefixrepo.go:46`), so identical transport behavior — flag-driven, not folder-driven.
- `clone-now` / `reclone` (manifest shape — `gitmap/cmd/clonenow.go:82`) — **HONORS transport** for HTTPS-origin records and SSH-origin records via the shared picker. The dispatcher routes through `cloner.pickURL` (`gitmap/cloner/summary.go:41`) which checks `rec.Transport == "ssh"` before falling through to HTTPS (fixed in v6.20.0). The scan-wide `--mode=https` default no longer overrides per-record SSH transport. Per-record `Transport` is populated by scan; manifest rows missing the field fall back to URL-prefix sniffing.
- `clone` (direct-URL form — `gitmap/cmd/clone.go:337`) — **HONORS transport** trivially: clones the literal URL the user typed; `upsertDirectClone` (`clone.go:337`) records SSH vs HTTPS by inspecting `PrefixSSH`/`https://` and writes the right column on the `Repo` row. No second-guessing.
- `cloner.pickURL` (shared picker — `gitmap/cloner/summary.go:41`) — **HONORS transport.** Reference implementation for the v6.21+ rule: `Transport == "ssh"` + non-empty `SSHUrl` → SSH; else HTTPS; else SSH fallback.

## Remaining gaps (drive step 2 + step 3 of plan 03)

1. **`cfr` / `cfrp` do not consult the destination folder's existing origin** before issuing the clone. If the user is "re-fixing" a folder that already has `.git/` and an SSH origin, the user's HTTPS positional URL silently downgrades the transport. Fix in step 3 (plan 03): before `executeDirectClone`, when `absPath` already contains `.git/`, classify `remote.origin.url` and prefer the existing transport (with a one-line stderr notice when it diverges from the positional URL).
2. **`Repo.IdentifiedTransport` is not persisted** (Step 2 of plan 03). Today `Repo` carries `HTTPSUrl`/`SSHUrl` only — the classified transport lives in `model.ScanRecord` (`Transport`) and `Repo.IdentifiedTransport` does not exist yet. Until column 007 lands, every reclone classifies from scratch each run; there is no memory across invocations and `gitmap history` cannot surface the transport choice.
3. **No history row is written on reclone.** None of the reclone paths above append to a history table; `gitmap history` currently surfaces `history-purge`/`history-pin` operations only (per `mem://features/history-rewrite`). Step 3 of plan 03 will add a `RecloneHistory` insert at the end of each reclone path (success or failure).

## Picker swap recipe (for the partial commands)

Mirror the v6.21.0 `pickURLForTransport` shape:

```go
// pickRecloneURL prefers the existing folder's origin transport when
// the destination already has a .git/ directory, otherwise uses the
// user-supplied URL verbatim.
func pickRecloneURL(positional string, existingOrigin string) string {
    if existingOrigin == "" {
        return positional
    }
    return existingOrigin // already correct transport for the folder
}
```

## Files referenced

- `gitmap/cmd/reporeclone.go` (lines 27, 33, 62, 81, 97, 106, 134)
- `gitmap/cmd/clonefixrepo.go` (lines 33, 39, 46, 79)
- `gitmap/cmd/clonenow.go` (lines 82, 87, 91, 119)
- `gitmap/cmd/clone.go` (lines 337, 344-348)
- `gitmap/cloner/summary.go` (lines 36-50)
