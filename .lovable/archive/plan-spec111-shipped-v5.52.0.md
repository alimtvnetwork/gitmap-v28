
# Remote Installer + Sibling Probing in Go — 5-Step Rollout

Move the `-v<N+i>` sibling-repo probe out of the downloaded `install.{ps1,sh}` and into Go (`gitmap update`), so the update flow no longer depends on shelling out to a shell script for the discovery phase. The remote installer is still used, but only as the final "install the resolved version" step — not as the probe engine.

## Step 1 — Spec

Create `spec/01-app/111-update-remote-probe.md` describing:
- Goal: native Go probe + remote install, single source of truth in Go.
- Probe contract: parse current repo slug (`gitmap-v27`) → derive base (`gitmap`) + `N=23` → fire 20 parallel `HEAD` requests against `https://github.com/alimtvnetwork/gitmap-v<N+1..N+20>` → max-hit index wins → fall back to latest release of current repo → fall back to `main` HEAD of current repo.
- Install contract: once winning repo slug is resolved, download that repo's `install.ps1` / `install.sh` from `raw.githubusercontent.com/.../main/install.{ps1,sh}` and exec with inherited stdio.
- Flags: `--probe-only` (print resolution + exit), `--no-probe` (skip probe, install from current repo only), `--source-rebuild` (unchanged, legacy path).
- Exit codes, timeouts (5s per HEAD, 30s download), and logging format.

## Step 2 — Constants + Probe Package

- Extend `gitmap/constants/constants_update.go` with: `UpdateProbeMaxSiblings=20`, `UpdateProbeTimeoutSec=5`, `UpdateProbeRepoBase="gitmap"`, `UpdateRepoOwner="alimtvnetwork"`, `UpdateRepoHEADTmpl`, `UpdateRawInstallerTmpl`, plus `Msg*` / `Err*` strings for probe lifecycle.
- New file `gitmap/cmd/updateprobe.go` with `resolveLatestRepoSlug() (slug string, source string, err error)`:
  - `parseCurrentRepoSlug()` from the embedded `RepoSlug` constant.
  - `probeSiblings(base, n, max)` — `sync.WaitGroup` + buffered channel of results, `http.Client{Timeout: 5s}`, `HEAD` requests, return highest-N 2xx.
  - `fallbackLatestRelease(slug)` — hit `api.github.com/.../releases/latest`.
  - Final `fallbackMain(slug)`.

## Step 3 — Wire Probe Into Update Flow

- Refactor `gitmap/cmd/updateremoteinstall.go`:
  - Call `resolveLatestRepoSlug()` first, log winning slug + source.
  - Build installer URL from the winning slug (not hardcoded to `gitmap-v27`).
  - Honor `--no-probe` (skip step 1, use current slug) and `--probe-only` (print + return).
- Update `gitmap/cmd/update.go` dispatcher to parse the two new flags via existing flag-reorder helper.

## Step 4 — Tests

- `gitmap/cmd/updateprobe_test.go`:
  - Table-driven `parseCurrentRepoSlug` for `gitmap-v27`, `gitmap-v100`, malformed.
  - `probeSiblings` against a local `httptest.Server` returning 200 for a configured subset, asserting max-hit wins and timeout doesn't deadlock.
  - Fallback chain: all-404 siblings → release fallback called; release 404 → main fallback called.
- Skip live network in CI; everything httptest-based.

## Step 5 — Version Bump + Docs

- Bump to **v5.52.0** in `gitmap/constants/constants.go`, `README.md` pinned block + matrix, `src/constants/index.ts`.
- Append changelog entry to `CHANGELOG.md` and `src/data/changelog.ts`.
- Create `.gitmap/release/v5.52.0.json` and update `.gitmap/release/latest.json`.
- Link new spec from README "Update" section.

## Technical notes

- Probe is read-only `HEAD`; no auth needed for public repos.
- Use `context.WithTimeout` per request; collect results in a fixed-size slice indexed by offset to avoid mutex.
- Embedded `RepoSlug` already exists in `gitmap-updater/cmd/constants.go`; mirror that constant into `gitmap/constants` rather than reaching across module boundaries.
- No behavior change for `--source-rebuild`; legacy handoff path stays intact for power users.

Say **next** to execute Step 1.
