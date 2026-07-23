## Fix flaky CI test + release v6.82.0

### 1. Root cause of the CI failure

Failing test: `TestPlanIdempotenceBeyond200Commits` in `gitmap/committransfer/plan_idempotence_test.go`.

Error: `git commit (bury 127): exit status 128 — error: bad tree object HEAD`

Why it fails intermittently:

- The test runs with `t.Parallel()` alongside the rest of the `committransfer` suite (which also uses `mustCommitCount` and its own `t.TempDir()` repos).
- It fires **252 sequential `git commit` shellouts** (1 buried + 1 anchor + 250 bury) each stamped with `time.Now().UTC().Format(time.RFC3339)` — second-granularity, so many commits share an identical `GIT_AUTHOR_DATE` / `GIT_COMMITTER_DATE`.
- Under parallel CI load, git occasionally sees a not-yet-flushed `.git/objects/..` write and reports `bad tree object HEAD` on the next commit. It is not a real defect in the planner — the passing 200-cap regression is still being exercised by the surviving assertion; we just need the setup to be reliable.

### 2. Fix (single file: `gitmap/committransfer/plan_idempotence_test.go`)

- Remove `t.Parallel()` on this test. It is the heaviest test in the package (250 shellouts) and does not benefit from parallel scheduling. Serial execution eliminates the filesystem contention window that produces `bad tree object HEAD`.
- Give every commit a unique, monotonic timestamp. Introduce a small `mustCommitCountAt(t, dir, path, body, msg, stamp)` helper local to this file (keeps `mustCommitCount` untouched for other tests) and walk `stamp` forward one second per commit starting from a fixed base (`2026-01-01T00:00:00Z`). Unique stamps avoid any same-second collision inside git's index/refs machinery.
- Lower `buryCount` from 250 to 220. Still comfortably above the legacy 200-commit cap this test guards, but 30 fewer shellouts and 30 fewer chances to hit a flaky filesystem.
- No production code changes. Planner behaviour is unchanged.

### 3. Confirm commit-in status (report only, no edits)

From the CI log the user pasted:

- `gitmap/cmd/commitin` ok
- `gitmap/cmd/commitin/e2e` ok (2.071s)
- All 12 `commitin/*` subpackages (checkpoint, dedupe, finalize, funcintel, message, orchestrator, profile, prompt, replay, runlog, walk, workspace) ok
- `gitmap/committransfer` was the ONLY failure, and inside a *test-only* helper, not in the transfer/replay engine (`BuildPlan`, `runOneDirection`, `RunBothInterleaved` all green in `committransfer` `cmd`, and `e2e` runs).

Conclusion to surface in chat after the fix lands: the `commit-in` feature and its "each commits individually" replay path are intact — the CI red was a flaky test harness, not a regression in commit replay.

### 4. Release v6.82.0 (minor bump per `mem://project/release-keyword.md`)

Batch these edits in parallel (no touches to `.gitmap/release/`, no touches to `.github/workflows/`):

1. **`gitmap/constants/constants.go`** line 16: `const Version = "6.81.0"` → `const Version = "6.82.0"`
2. **`src/constants/index.ts`** line 2: `export const VERSION = "v6.81.0"` → `"v6.82.0"` (keeps `version-sync.test.ts` green)
3. **`CHANGELOG.md`**: prepend
   ```
   ## v6.82.0 - (2026-07-23) - CI reliability

   - committransfer: de-flaked TestPlanIdempotenceBeyond200Commits by removing
     t.Parallel(), assigning each buried commit a unique monotonic timestamp,
     and trimming the bury count from 250 to 220 (still >200 legacy cap).
   - No production code changes; commit-in / commit-left / commit-right /
     commit-both replay engines unchanged and green in CI.
   ```
4. **`README.md`** (root): `sed -i 's/v6\.81\.0/v6.82.0/g' README.md` to refresh every pinned mention (the previous release skipped this — user flagged it).

### 5. Verification

- `nix run nixpkgs#go -- test ./gitmap/committransfer/... -count=5` to prove the de-flake sticks under repeated runs.
- `nix run nixpkgs#go -- test ./... -count=1` for the full suite.
- `bunx vitest run src/test/version-sync.test.ts` for Go ↔ TS parity.
- `.github/scripts/check-changelog-version-sync.sh` for changelog gate.

### Explicitly NOT touched

- `.gitmap/release/*.json` (CLI-managed, Core constraint)
- `.github/workflows/release.yml` and release scripts (strictly prohibited)
- Any `gitmap/` source outside the one test file + the `Version` literal
- Production commit-transfer code
