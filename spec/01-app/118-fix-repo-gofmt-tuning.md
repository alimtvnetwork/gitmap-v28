# 118 — fix-repo gofmt tuning (v6.80.1)

Follow-up to the chunker shipped in v6.80.0
(`.lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md`). Adds
observability and a tuning knob so Windows users can diagnose and
work around setups where the effective `CreateProcess` argv cap is
below the documented 32,767 characters.

## Deliverables

1. **`gitmap doctor fix-repo`** subcommand with four probes:
   `gofmt-present`, `gofmt-runs`, `argv-budget`, `chunker-selftest`.
   Supports `--json` and `--budget N`.
2. **`gitmap fix-repo --dry-run` batch preview**: per-batch cmdLen
   in bytes and percent of budget; `NEAR-LIMIT` at ≥90%,
   `OVER-LIMIT` at ≥100%.
3. **`gitmap fix-repo --verbose` progress**: batch header, per-batch
   start / done lines, rolling ETA computed from average per-batch
   wall time.
4. **`--gofmt-max-cmd-len N`** CLI flag. Floor 512; overrides
   `constants.FixRepoGofmtMaxCmdLen` for the current run.

## Non-goals

- No parallel batch execution (keeps error attribution simple).
- No config-file surface for the budget (may follow if the flag is
  used enough to justify the round-trip).
- No PowerShell parity yet; the legacy `.ps1` still lacks argv
  chunking.

## Files touched

- `gitmap/cmd/fixrepo.go`, `gitmap/cmd/fixrepo_flags.go`,
  `gitmap/cmd/fixrepo_gofmt.go`
- `gitmap/cmd/doctor_run.go`, `gitmap/cmd/doctor_fixrepo.go`
- `gitmap/cmd/rootusageflags.go`
- `gitmap/constants/constants_fixrepo.go`,
  `gitmap/constants/constants_fixrepohelp.go`
- `gitmap/helptext/doctor-fix-repo.md`
- Tests: `gitmap/cmd/fixrepo_gofmt_test.go`

## Validation

- `go test ./gitmap/cmd/... ./gitmap/constants/...` green.
- `bunx vitest run src/test/version-sync.test.ts` green.
- Manual on Windows: `gitmap doctor fix-repo` reports measured cap;
  `gitmap fix-repo --all --dry-run --verbose --gofmt-max-cmd-len 8000`
  prints multi-batch preview + progress.
