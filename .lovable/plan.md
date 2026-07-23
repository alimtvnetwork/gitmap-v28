## Context (verified)

- `gitmap/cmd/fixrepo_gofmt.go` (137 lines) already has `chunkPathsForGofmt`, `shortenGofmtPaths`, `invokeGofmt`. Chunking uses `constants.FixRepoGofmtMaxCmdLen` (30000).
- `fixRepoOptions` in `gitmap/cmd/fixrepo.go` L18-31 owns flags; `isDryRun` and `isVerbose` already exist. Flags parsed in `gitmap/cmd/fixrepo_flags.go` (see `consumeOneFixRepoArg` L74, config-arg pattern L185, restrict-arg L217).
- Gofmt is invoked from `runFixRepo` L58 with `runFixRepoGofmt(result.goFiles, opts)`. Dry-run currently short-circuits at L34-38 with `FixRepoMsgGofmtSkip`.
- Doctor lives in `gitmap/cmd/doctor_run.go` + `defaultDoctorChecks()`; each check is a `DoctorCheck{Name, Run, FixHint}`. Sub-args aren't wired: `runDoctor` only reads `--json` / `--fix`. There is no existing sub-command dispatch under `doctor`.
- Constants for fix-repo gofmt live in `gitmap/constants/constants_fixrepo.go`; help text in `gitmap/constants/constants_fixrepohelp.go`.
- History: `.lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md` + `.lovable/plan.md` capture the Windows argv-overflow root cause and the shipped chunker.

## Scope

Four features + docs + tests + release. All changes are additive; no behavior change when new flags are absent.

---

### 1) `gitmap doctor fix-repo` subcommand

**Where:** new `gitmap/cmd/doctor_fixrepo.go`, dispatched from `runDoctor` in `gitmap/cmd/doctor_run.go`.

**How:**

- In `runDoctor`, before the current flag loop, peek `args[0]`. If it equals `constants.CmdFixRepo` (or alias `fr`), delegate to `runDoctorFixRepo(args[1:])` and return. Otherwise keep today's behavior.
- `runDoctorFixRepo` runs four probes and prints them with the same `[ok]/[fail]` renderer used by `emitDoctorText` (extract the row printer into a helper `printDoctorRow` so both paths share it):
  1. **gofmt-present**: `exec.LookPath("gofmt")`. On fail: hint to install Go / set PATH.
  2. **gofmt-runs**: `gofmt -l` against a tempdir with one throwaway `.go` file. Confirms the binary is executable (catches the exact `fork/exec ... The filename or extension is too long` class when it's actually a PATH-injected shim).
  3. **argv-budget**: probe the effective command-line cap. On Windows use `constants.FixRepoGofmtMaxCmdLen` as the documented budget and additionally attempt a synthetic exec (`gofmt -l` with N filler `-`-prefixed no-op args generated from a tempdir) doubling until failure to *measure* the real cap; report both "configured budget" and "measured cap". On non-Windows just print the configured budget and skip measurement (ARG_MAX is ~2 MB, not the bottleneck).
  4. **chunker-selftest**: call `chunkPathsForGofmt` with three synthetic inputs (empty, all-fits, forced-overflow using 500 × 200-char fake paths) and assert invariants inline (batch count ≥ 2 on overflow, no batch over budget except single-path overflow). Failures print which invariant broke.
- Supports `--json` (reuses `emitDoctorJSON` shape) and `--budget <N>` to override the budget used during measurement.
- Exit 0 if all probes pass, 1 otherwise.

**Help:** new `gitmap/helptext/doctor-fix-repo.md`.

---

### 2) `--dry-run` batch preview

**Where:** `gitmap/cmd/fixrepo_gofmt.go` + `gitmap/constants/constants_fixrepo.go`.

**How:**

- Change `runFixRepoGofmt`'s dry-run branch (L34-38): instead of printing only `FixRepoMsgGofmtSkip` and returning, when `len(goFiles) > 0` also compute `paths := shortenGofmtPaths(goFiles)` and `batches := chunkPathsForGofmt(paths, effectiveBudget(opts))` and print:
  - `gofmt (dry-run): would run N batch(es) across F file(s), budget=B`
  - Per batch:   `batch i/N: files=K, cmdLen=L bytes (P% of budget)` where cmdLen = sum(len(p)+1) + len("gofmt -w ").
  - Highlight batches ≥ 90% of budget with `NEAR-LIMIT` suffix (uses `constants.ColorYellow`); batches ≥ 100% with `OVER-LIMIT` and `constants.ColorRed`.
- New constants: `FixRepoMsgGofmtDryFmt`, `FixRepoMsgGofmtDryBatchFmt`, `FixRepoGofmtNearLimitPct = 90`.
- Non-dry-run path unchanged.

---

### 3) Verbose progress output

**Where:** `gitmap/cmd/fixrepo_gofmt.go`, plumbed through by passing `opts` (or just `opts.isVerbose`) into `invokeGofmt`.

**How:**

- Change `invokeGofmt(paths []string)` → `invokeGofmt(paths []string, opts fixRepoOptions)` and update its single call site.
- When `opts.isVerbose`:
  - Print `gofmt: N batch(es), F file(s), budget=B` before the loop.
  - Track `start := time.Now()`. Before each batch:   `[i/N] formatting K files (cmdLen=L)...`. After each batch:     `done in Dms; ETA ~Es` where `E = elapsed/i * (N-i)`.
  - Round times to milliseconds (batches finish fast on small repos); use `time.Since(start).Truncate(time.Millisecond)`.
- Non-verbose path prints only the existing final summary line.
- New constants: `FixRepoMsgGofmtVerbHeaderFmt`, `FixRepoMsgGofmtVerbBatchStartFmt`, `FixRepoMsgGofmtVerbBatchDoneFmt`.

---

### 4) `--gofmt-max-cmd-len <N>` flag

**Where:** `gitmap/cmd/fixrepo.go` (opts struct), `gitmap/cmd/fixrepo_flags.go` (parser), `gitmap/cmd/fixrepo_gofmt.go` (consumer), `gitmap/constants/constants_fixrepohelp.go` (help row).

**How:**

- Add `gofmtMaxCmdLen int` to `fixRepoOptions` (0 = use default).
- Parser: extend `consumeOneFixRepoArg` following the `--config <path>` shape (see L185 `consumeFixRepoConfigArg`). Accept `--gofmt-max-cmd-len N` and `--gofmt-max-cmd-len=N`. Validate `N >= 512`; on `< 512` or non-integer, return a parse error using `FixRepoErrBadFlagFmt`. Lower bound 512 keeps at least one real path per batch.
- Helper `effectiveBudget(opts fixRepoOptions) int`: returns `opts.gofmtMaxCmdLen` when > 0, else `constants.FixRepoGofmtMaxCmdLen`. Used by both dry-run preview (feature 2) and the real `invokeGofmt` chunking call.
- Help: add `HelpFRGofmtMaxCmdLen` in `constants_fixrepohelp.go` and include it in the printed help block.

---

### Tests

- Extend `gitmap/cmd/fixrepo_gofmt_test.go`:
  - `TestEffectiveBudget` — default when unset, override when set.
- New `gitmap/cmd/fixrepo_flags_gofmt_test.go`:
  - Parses `--gofmt-max-cmd-len 5000` → `opts.gofmtMaxCmdLen == 5000`.
  - `--gofmt-max-cmd-len=5000` equivalent.
  - `--gofmt-max-cmd-len 100` → error (below floor).
  - Missing value → error.
- New `gitmap/cmd/doctor_fixrepo_test.go`:
  - Chunker-selftest probe returns ok on valid budgets.
  - Dispatch: `runDoctor` with args `["fix-repo","--json"]` routes to fix-repo doctor (assert via captured stdout containing `"probe":"gofmt-present"`).
- Dry-run preview: add a table test in `fixrepo_gofmt_test.go` that captures stdout while calling the dry-run branch with a synthetic goFiles list, asserts batch count line + NEAR-LIMIT tagging when budget forces a 92% batch.

---

### Documentation

- `gitmap/helptext/fix-repo.md` — document `--gofmt-max-cmd-len` and the new verbose/dry-run output shape.
- `gitmap/helptext/doctor-fix-repo.md` — new file for the subcommand.
- `spec/01-app/` — new spec `118-fix-repo-gofmt-tuning.md` covering the four features (mirrors the numbering pattern used by `117-update-awareness.md`).
- `src/data/commands.ts` — add `doctor fix-repo` entry and the new flag row for `fix-repo`.
- `.lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md` — append a "Follow-up (v6.80.1)" section noting the diagnostics and tuning knob shipped.

---

### Release

Bump `6.80.0 → 6.81`in:

- `gitmap/constants/constants.go`
- `.gitmap/release/latest.json`
- `src/constants/index.ts`
- Changelog "Fixed / Added" entry.
- Re-run `bunx vitest src/test/version-sync.test.ts`.

---

## Files touched

```text
gitmap/cmd/doctor_run.go                          (edit: dispatch fix-repo sub)
gitmap/cmd/doctor_fixrepo.go                      (new: 4 probes)
gitmap/cmd/doctor_fixrepo_test.go                 (new)
gitmap/cmd/fixrepo.go                             (edit: opts field)
gitmap/cmd/fixrepo_flags.go                       (edit: parse --gofmt-max-cmd-len)
gitmap/cmd/fixrepo_flags_gofmt_test.go            (new)
gitmap/cmd/fixrepo_gofmt.go                       (edit: dry-run preview, verbose, budget)
gitmap/cmd/fixrepo_gofmt_test.go                  (edit: budget + dry-run tests)
gitmap/constants/constants_fixrepo.go             (edit: 5 new message + pct constants)
gitmap/constants/constants_fixrepohelp.go         (edit: HelpFRGofmtMaxCmdLen)
gitmap/helptext/fix-repo.md                       (edit)
gitmap/helptext/doctor-fix-repo.md                (new)
spec/01-app/118-fix-repo-gofmt-tuning.md          (new)
src/data/commands.ts                              (edit: docs site)
.lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md  (edit: follow-up note)
CHANGELOG + version files                         (edit: 6.80.1)
```

---

## Pending tasks in this project (from `plan.md` and open memory)

- **P0-1** Version sync guardrail (resolved for 6.80.x — keep the vitest guard green each bump).
- **P0-2** Spec dedup / renumber for the numeric-prefix collisions in `spec/01-app/` (116+ files).
- **P0-3** Function-length audit (Go files > ~200 lines flagged by house style).
- **P1-2** SSH clone/reclone transport preservation (`.lovable/issues/01`, `02`).
- **P2-1** Test backfill for `gitmap update all` worker pool.
- **Suggestions workflow** — `.lovable/memory/suggestions/` scaffolded, no entries yet.

---

## 10 additional improvement ideas (not in this plan, listed for triage)

1. **Parallel gofmt batches** — bounded worker pool (GOMAXPROCS) when N batches ≥ 4; keep sequential fallback for clean error attribution.
2. `**gitmap fix-repo --json**` — machine-readable summary (files scanned/changed, batches, timings) for CI consumers.
3. **Persistent gofmt log** — write per-batch stdout/stderr to `.gitmap/logs/fixrepo-gofmt-<ts>.log` when `--verbose`.
4. **Config-file budget override** — allow `gofmtMaxCmdLen` in `fix-repo.config.json` so Windows users don't need to pass the flag every run.
5. `**gitmap doctor --watch**` — re-run probes on file change; useful during environment troubleshooting.
6. **Argv-budget autotune** — on first Windows failure, halve budget and retry once; record the working budget under `.gitmap/state/gofmt-budget`.
7. **Path-shortening telemetry** — verbose mode reports "shortened X/Y paths, avg saving Zc" so tuning is data-driven.
8. `**gitmap fix-repo --only-gofmt**` — skip token rewrite, just run the chunked gofmt pass (useful after a manual bump).
9. **Cross-platform argv probe** — the `argv-budget` doctor probe on Linux/macOS measures via `getconf ARG_MAX`, exposes it in `--json` for parity.
10. **PowerShell parity** — port batching + `-GofmtMaxCmdLen` to `scripts/fix-repo/Rewrite-Engine.ps1` for legacy users.

---

## Validation

1. `go test ./gitmap/cmd/... ./gitmap/constants/...` passes.
2. `bunx vitest run` green (version-sync guard).
3. Manual: on the user's Windows repo, `gitmap doctor fix-repo` all-ok; `gitmap fix-repo --all --dry-run --verbose` shows batch table with per-batch cmdLen and no NEAR-LIMIT flags at default budget; `gitmap fix-repo --all --gofmt-max-cmd-len 8000 --verbose` shows more, smaller batches all succeed.