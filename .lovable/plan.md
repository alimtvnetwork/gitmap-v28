## Release v6.81.0 (minor bump from 6.80.1)

Follows the canonical "release" procedure in `mem://project/release-keyword.md` and `version-bump-procedure.md`. Batches 4 file edits, no touches to `.gitmap/release/` (managed by the CLI itself), no touches to the CI release pipeline.

### Files to update

1. **`gitmap/constants/constants.go`** (line 16)
   - `const Version = "6.80.1"` → `const Version = "6.81.0"`

2. **`src/constants/index.ts`** (line 2)
   - `export const VERSION = "v6.80.1"` → `"v6.81.0"`
   - Kept in lockstep so `src/test/version-sync.test.ts` stays green.

3. **`CHANGELOG.md`** (top)
   - Insert new `## v6.81.0 — (2026-07-23) — Minor release` section summarizing the shipped work since 6.80.0:
     - `gitmap doctor fix-repo` probe suite (gofmt-present, gofmt-runs, argv-budget, chunker-selftest)
     - `fix-repo --dry-run` batch preview with NEAR-LIMIT / OVER-LIMIT tags
     - `fix-repo --verbose` progress with per-batch ETA
     - `--gofmt-max-cmd-len N` tuning flag (floor 512)
     - Golden-test fixes for `hd`, `list-update`, `update-all`, `update-apply` helptext

4. **`README.md`** (root)
   - `sed -i 's/v6\.80\.1/v6.81.0/g' README.md` to refresh every pinned mention.

### Verification

- Run `bunx vitest run src/test/version-sync.test.ts` to confirm Go ↔ TS parity.
- Run `.github/scripts/check-changelog-version-sync.sh` locally-equivalent check by grepping `## v6.81.0` in CHANGELOG.

### Explicitly NOT touched

- `.gitmap/release/*.json` (CLI-managed, per Core constraint)
- `.github/workflows/release.yml` and release scripts (strictly prohibited)
- Any other `gitmap/` source file
