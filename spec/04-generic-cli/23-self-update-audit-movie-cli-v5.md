# Audit: movie-cli-v5 Self-Update vs Gitmap Gold Standard

> AI-shareable instruction document. Hand this to any assistant fixing
> `movie-cli-v5` (or any similar CLI). Each section names the **exact defect**
> in the audited repo and the **exact fix** from the gold standard
> ([22-self-update-gold-standard.md](22-self-update-gold-standard.md)).

**Audited source**:
- `https://github.com/alimtvnetwork/movie-cli-v5/blob/main/install.ps1`
- `https://github.com/alimtvnetwork/movie-cli-v5/blob/main/run.ps1`

**Reference**: gitmap-v27 repo (this codebase), `gitmap-v27/cmd/update.go`,
`gitmap-v27/cmd/updatescript.go`, `run.ps1` (`-Update` branch).

---

## Verdict

`movie-cli-v5` has a partial deploy story but **no real self-update flow**.
`install.ps1` is reused as the update path, which means every update:

- runs from the *user's CWD* not from the active binary location,
- copy-overwrites the deployed binary (will fail the moment any process holds
  it open on Windows),
- skips the handoff entirely (no separate worker process),
- has no cleanup of stale artifacts,
- has no version-equality check across PATH after deploy,
- relies on `Read-Host` in `run.ps1` (line 184) — fatal in non-interactive sessions.

---

## Defect-by-Defect Comparison

### D1 — No two-phase handoff

**movie-cli-v5**: `install.ps1` calls `& $runScript` directly inside the same
process tree. The currently running `movie.exe` (if upgrading from PATH) holds
its install path locked. There is no temp copy, no worker subcommand.

**Gold standard**: `tool update` copies itself to
`<installDir>/<tool>-update-<pid>.exe`, launches that copy with the hidden
`update-runner` subcommand using `cmd.Run()` (foreground, blocking), and lets
*that* copy perform the deploy. The original lock is released because the
worker is a different file.

**Fix**: implement `movie update` and `movie update-runner` in Go (or whatever
language) that mirror `gitmap-v27/cmd/update.go` lines 18–157. Do **not** keep
`install.ps1` as the update path.

---

### D2 — Copy-overwrite of locked binary (`Deploy-Binary` in run.ps1)

**movie-cli-v5** `run.ps1` Deploy-Binary (lines 567–700ish):

```
Rename-Item -Path $destFile -NewName "$binaryName.bak" -Force
Copy-Item -Path $SourceBinary -Destination $destFile -Force
```

Renaming `movie.exe → movie.exe.bak` works on Windows, BUT the script then
deletes the `.bak` immediately (`Remove-Item -Path $backupFile`) instead of
keeping `.old` as a rollback artifact. There is also no `-Update` switch path
that uses **rename-first as the primary strategy** — it does it for backup,
not for lock-avoidance.

**Gold standard**: in update mode the deploy is **rename-first**, the renamed
file becomes `<tool>.exe.old` (kept), the new binary is copied in, and a
20×500 ms copy-retry loop is **fallback only**. `.old` files are removed by
`tool update-cleanup` after success.

**Fix**: split deploy into two modes inside `run.ps1`:
- normal deploy: current logic is fine
- `-Update` mode: rename to `.old` (keep), copy new, never delete `.old` here
  (cleanup command does it).

---

### D3 — `Read-Host` inside `run.ps1` (line 184)

```
$choice = Read-Host "  Enter choice (S/D/C/Q)"
```

**Why it breaks update**: when invoked from a hidden worker process during
update, stdin is not interactive. The script blocks forever or fails. This is
explicitly forbidden in the gold standard ([§2 prohibitions](22-self-update-gold-standard.md#hard-prohibitions)).

**Fix**: gate any `Read-Host` behind `if (-not $Update -and -not $env:CI)`.
Provide non-interactive defaults for every prompt path.

---

### D4 — No `update-cleanup` subcommand

`movie-cli-v5` deletes `.bak` immediately and never has any concept of stale
handoff copies (`movie-update-*.exe`) — because it has no handoff. After
adopting D1, you will accumulate temp copies. You need:

```
movie update-cleanup
```

Idempotent. Removes `<installDir>/movie-update-*.exe`, `movie.exe.old`,
`%TEMP%\movie-update-*.exe`, and `movie.exe.bak`. Runs automatically at end of
successful update; safe to invoke manually.

Reference: `spec/08-generic-update/06-cleanup.md`.

---

### D5 — Repo path resolution is not persisted

`install.ps1` re-derives repo path from CWD on every call (`Get-Location`,
fall back to clone into `./movie-cli-v3`). On the second run the user must
be in the right directory or it clones again next to wherever they invoked
the script.

**Gold standard** resolution chain (stop at first hit):

1. `--repo-path` flag
2. Embedded ldflags constant
3. SQLite / state DB
4. TTY prompt
5. `<tool>-updater` (release-based) fallback
6. Fail with 4-option actionable message

**Fix**: persist resolved path in a state DB or `~/.movie/state.json` and
check it first. See `gitmap-v27/cmd/update.go` `resolveRepoPath()`.

---

### D6 — No post-update version equality check

`run.ps1` runs `& $destFile version` and prints whatever comes out, but it
**does not assert** that the binary on PATH (`Get-Command movie`) matches the
deployed binary version. The script literally warns "PATH resolves 'movie' to
a different binary" and continues with exit 0. That hides broken updates.

**Fix**: after deploy, fail the script (`exit 1`) if:

```
(& (Get-Command movie).Source version) -ne (& $destFile version)
```

Print before → after explicitly:

```
v0.4.2 → v0.5.0
✓ Updated to v0.5.0
```

---

### D7 — Repo name mismatch (`movie-cli-v3` vs `movie-cli-v5`)

`install.ps1` clones from `https://github.com/alimtvnetwork/movie-cli-v3.git`
into folder `movie-cli-v3`, but the actual repo is `movie-cli-v5`. Self-update
will clone the **wrong project** on a fresh machine. Pure bug, fix the URL and
folder name and pin them in a single `$Repo*` constant block.

---

### D8 — No script encoding handling

`install.ps1` and `run.ps1` are not generated by the binary, but if `movie`
ever generates an update script (it should, per gold standard), it MUST be
written with **UTF-8 BOM** so PowerShell handles Unicode glyphs in output.
See `gitmap-v27/cmd/updatescript.go::writeScriptToTemp` (BOM = `0xEF 0xBB 0xBF`).

---

### D9 — No async/foreground discipline documented

Nothing in the codebase prevents a future contributor from "fixing" the
locked-binary problem by detaching with `Start-Process -NoNewWindow` or
`cmd.Start()` + `os.Exit(0)`. That kills the terminal session on Windows.

**Fix**: add a lint/grep guard in CI:

```
! grep -RIn 'cmd.Start(' --include='*.go' ./cmd/update*.go
! grep -RIn 'Start-Process' run.ps1
```

And document the prohibition in `CONTRIBUTING.md`.

---

## Migration Plan for movie-cli-v5

Order matters. Each step is independently testable.

1. **Fix D7** (repo URL/name) — 5 minute change, unblocks everything else.
2. **Fix D3** (gate `Read-Host` behind interactive check) — required before
   any non-interactive invocation works.
3. **Add `-Update` switch to `run.ps1`** with rename-first deploy (D2). Keep
   normal deploy untouched.
4. **Implement Go `movie update` + hidden `movie update-runner`** mirroring
   `gitmap-v27/cmd/update.go`. Use foreground `cmd.Run()` (D1, D9).
5. **Implement `movie update-cleanup`** (D4) and call it at end of update.
6. **Implement repo-path resolution chain + persistence** (D5).
7. **Add post-update version equality assertion** to `run.ps1 -Update`
   tail (D6).
8. **Document prohibitions** in CONTRIBUTING.md and add CI grep guards (D9).
9. **Add `helptext/update.md`** with the 6 canonical example scenarios from
   `gitmap-v27/helptext/update.md`.

---

## Quick Diff Summary

| # | Defect | File | Severity |
|---|---|---|---|
| D1 | No handoff, deploys from locked binary location | `install.ps1`, `run.ps1` | 🔴 Critical |
| D2 | Copy-overwrite, no rename-first in update mode | `run.ps1::Deploy-Binary` | 🔴 Critical |
| D3 | `Read-Host` blocks non-interactive update | `run.ps1:184` | 🔴 Critical |
| D4 | No `update-cleanup` subcommand | (missing) | 🟠 High |
| D5 | Repo path not persisted, re-derived from CWD | `install.ps1` | 🟠 High |
| D6 | No PATH==deploy version equality check | `run.ps1` | 🟠 High |
| D7 | Wrong repo URL (`v3` vs `v5`) | `install.ps1` | 🟡 Medium |
| D8 | No UTF-8 BOM policy for generated scripts | (future) | 🟡 Medium |
| D9 | No prohibition guard against async detach | (CI) | 🟡 Medium |

---

## Hand-off Note for AI Assistants

When implementing the fixes above:

- Read [22-self-update-gold-standard.md](22-self-update-gold-standard.md)
  in full before writing any code.
- Read `gitmap-v27/cmd/update.go`, `gitmap-v27/cmd/updatescript.go`, and the
  `-Update` branch of `run.ps1` as the literal template.
- Do **not** invent alternative strategies (named pipes, scheduled tasks,
  WMI, etc.). They have been considered and rejected — see
  `spec/02-app-issues/02-update-flow-spec-alignment.md`.
- Every PR touching update flow MUST update all of:
  `spec/03-general/02f-self-update-orchestration.md`,
  `spec/04-generic-cli/22-self-update-gold-standard.md`,
  `spec/08-generic-update/*`, and the audited project's own checklist.
