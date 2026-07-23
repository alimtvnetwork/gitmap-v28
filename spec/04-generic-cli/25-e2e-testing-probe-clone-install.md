# 25 — End-to-End Testing: Probe, Clone, and Install Scripts

> **Status:** Authoritative draft (2026-04-26).
> **Audience:** Any AI agent or human implementer responsible for adding
> end-to-end (e2e) test coverage to the URL-based discovery, probe, and
> clone surface of any host repo following the gitmap-v27 framework.
> **Related specs:**
> - [12-testing.md](12-testing.md) — base unit/integration test layout
> - [../07-generic-release/09-generic-install-script-behavior.md](../07-generic-release/09-generic-install-script-behavior.md) — install-script contract under test
> - [../01-app/88-clone-direct-url.md](../01-app/88-clone-direct-url.md) — direct-URL clone behavior
> - [../01-app/95-installer-script-find-latest-repo.md](../01-app/95-installer-script-find-latest-repo.md) — sibling-version probe rationale
> - [../01-app/103-probe-depth.md](../01-app/103-probe-depth.md) — probe internals

The keywords **MUST**, **MUST NOT**, **SHOULD**, **MAY** follow RFC 2119.

---

## 0. How to use this document

When asked to implement e2e tests for the URL-driven flows in any host
repository:

1. Read this spec end-to-end.
2. Identify the host's package paths for `probe/`, `cloner/`, and the
   install-script directory. Substitute them in the file paths below.
3. Build the local-bare-repo fixture helper described in §3 once.
   Reuse it across all three suites.
4. Implement every test class marked **MUST**. Tests marked **SHOULD**
   may be deferred only with a written justification in the host's PR.
5. Mirror the §10 acceptance checklist in CI.

---

## 1. Scope

This spec covers e2e tests for three layers that together implement the
"give us a URL → resolve → install" pipeline:

| Layer | Package | Behaviors under test |
|-------|---------|----------------------|
| Probe | `gitmap-v27/probe/` | `ls-remote` happy path, shallow-clone fallback, empty-tag remote, malformed URL, temp-dir cleanup |
| Cloner (direct URL) | `gitmap-v27/cloner/` | URL classification, folder derivation, exists-conflict, successful clone, DB upsert |
| Install scripts | `scripts/install.sh`, `scripts/install.ps1` | Strict-tag mode, 20-parallel sibling probe, latest-release fallback, main-HEAD last-resort |

**Out of scope:** unit tests for pure helpers (covered by `12-testing.md`),
real network calls to `github.com` (forbidden in CI — see §3.4),
package-manager publishing flows.

---

## 2. Test placement and naming

Per [`12-testing.md`](12-testing.md):

```
tests/
├── e2e_probe_test/
│   ├── lsremote_test.go           ← happy + empty + malformed
│   ├── shallow_fallback_test.go   ← ls-remote fail → clone path
│   └── tempdir_cleanup_test.go    ← /tmp/gitmap-probe-* removed
├── e2e_cloner_test/
│   ├── direct_url_test.go         ← derive folder, classify, clone
│   ├── exists_conflict_test.go    ← target dir already present
│   └── db_upsert_test.go          ← record persisted after clone
└── e2e_install_test/
    ├── strict_tag_test.go         ← --version <tag> never falls back
    ├── sibling_probe_test.go      ← 20-parallel v<N+i> HEAD probe
    ├── latest_release_test.go     ← release-page fallback
    └── main_head_fallback_test.go ← last-resort branch HEAD
```

- Each test file **MUST** declare its own package (`e2e_probe_test`, etc.).
- Function names follow `Test<Layer>_<Scenario>` (e.g.
  `TestProbe_LsRemoteHappyPath`).
- Table-driven tests **MUST** be used for any scenario with ≥3 input
  variants.

---

## 3. Local bare-repo fixtures (no network)

All e2e tests **MUST** operate against locally-constructed bare
repositories under `t.TempDir()`. No test is permitted to touch the
public internet.

### 3.1 The `fixture` helper package

Create `tests/internal/fixture/fixture.go` exposing:

```go
type Repo struct {
    Dir     string // bare repo path, suitable as a clone URL (file://...)
    URL     string // "file://" + Dir
    Tags    []string
}

// NewBareRepo initialises an empty bare repo and seeds it with the
// supplied tags pointing at a single seed commit. Tags MUST be created
// in semver order so `git tag --sort=-v:refname` returns them descending.
func NewBareRepo(t *testing.T, tags ...string) *Repo

// NewEmptyBareRepo returns a bare repo with no tags and no commits.
func NewEmptyBareRepo(t *testing.T) *Repo

// NewBareRepoNoTags returns a bare repo with one commit on `main` but
// zero tags — exercises the "remote exists but has no tags" branch.
func NewBareRepoNoTags(t *testing.T) *Repo
```

### 3.2 Construction recipe

```go
// pseudocode — implementer fills in exec.Command details
git init --bare <tmp>/bare.git
git init <tmp>/seed
(cd <tmp>/seed
   echo seed > README.md
   git add . && git -c user.email=t@t -c user.name=t commit -m seed
   git remote add origin file://<tmp>/bare.git
   git push origin HEAD:refs/heads/main
   for tag in $tags; do git tag $tag && git push origin $tag; done)
```

### 3.3 URL form

The fixture **MUST** expose `file://<absolute-path>` as the clone URL.
This is sufficient for `git ls-remote`, `git clone --depth 1`, and
`git ls-remote refs/heads/v<N+i>` probes — exactly the operations the
production code performs.

### 3.4 No-network guard

Each e2e suite **MUST** fail fast if it detects accidental network
egress. A `TestMain` shim that sets `GIT_ALLOW_PROTOCOL=file` and
unsets `HTTP_PROXY`/`HTTPS_PROXY` is sufficient. CI **MUST** run the
e2e jobs with network disabled where the runner allows.

---

## 4. Probe layer e2e tests

Under test: `gitmap-v27/probe/probe.go` (`RunOne`) and
`gitmap-v27/probe/clone.go` (`tryShallowClone`), driven through the public
CLI surface `gitmap-v27 probe <URL>` (and `gitmap-v27 probe <URL> --json`) so
the tests exercise the same code path real users hit.

### 4.0 Per-scenario contract (read first)

Every P-scenario in §4.1–§4.3 **MUST** be expressed with the same five
sections so the test body can be generated mechanically:

1. **Fixture preconditions** — exact fixture-builder calls and any
   filesystem state that must exist *before* the command runs.
2. **CLI invocation** — the literal `argv` passed to the gitmap-v27 binary
   under test (built once per `TestMain` via `go build -o ./gitmap-e2e`).
   `${URL}` is the bare-repo `file://` URL from §3.
3. **Expected stdout / stderr / exit code** — asserted with substring
   matches against `constants.MsgProbe*` / `constants.ErrProbe*` (never
   hard-coded literals, per §7.5).
4. **Expected DB delta** — exact row(s) inserted into `VersionProbe`
   (and any tagging on `Repo`). Compared via `db.LatestVersionProbe`
   (or equivalent helper), never raw SQL.
5. **Cleanup assertions** — invariants that must hold *after* the
   command returns: temp-dir count delta, no orphan `gitmap-probe-*`
   dirs, no stray `git` child processes, no rows in unexpected tables.

A shared helper `runProbeCLI(t, args ...string) cliResult` returns
`{Stdout, Stderr, ExitCode, Duration}` and registers the §4.3 leak
guard via `t.Cleanup`.

### 4.1 Required test classes (MUST)

#### P1 — ls-remote returns highest semver tag

- **Preconditions:** `repo := fixture.NewBareRepo(t, "v1.0.0", "v1.0.5", "v1.0.20")`.
  DB pre-seeded with one `Repo` row whose `HTTPSUrl == repo.URL` (use
  `db.UpsertRepo`). `VersionProbe` table empty for this `RepoId`.
- **CLI:** `gitmap-e2e probe ${repo.URL}`
- **Stdout:** contains `fmt.Sprintf(constants.MsgProbeOkFmt, <slug>, "v1.0.20", constants.ProbeMethodLsRemote)`
  and `fmt.Sprintf(constants.MsgProbeDoneFmt, 1, 0, 0)`.
- **Stderr:** empty.
- **Exit code:** `0`.
- **DB delta:** exactly one new `VersionProbe` row with
  `NextVersionTag == "v1.0.20"`, `NextVersionNum == 20`,
  `Method == constants.ProbeMethodLsRemote`, `IsAvailable == 1`,
  `Error == ""`. `Repo.ScanFolderId` unchanged.
- **Cleanup:** §4.3 temp-dir delta = 0; no `gitmap-probe-*` directories
  remain; the bare repo at `repo.Dir` is byte-identical to its
  pre-test snapshot (compare via `fixture.HashTree`).

#### P2 — Remote has commits but zero tags

- **Preconditions:** `repo := fixture.NewBareRepoNoTags(t)`. DB has one
  `Repo` row for `repo.URL`.
- **CLI:** `gitmap-e2e probe ${repo.URL}`
- **Stdout:** contains `fmt.Sprintf(constants.MsgProbeNoneFmt, <slug>, constants.ProbeMethodLsRemote)`
  and `fmt.Sprintf(constants.MsgProbeDoneFmt, 0, 1, 0)`.
- **Stderr:** empty.
- **Exit code:** `0`.
- **DB delta:** one new `VersionProbe` row with `NextVersionTag == ""`,
  `NextVersionNum == 0`, `IsAvailable == 0`, `Error == ""`,
  `Method == constants.ProbeMethodLsRemote`.
- **Cleanup:** §4.3 invariant; no shallow-clone temp dir was created
  (assert via temp-dir snapshot — proves the no-tags branch did not
  fall through to clone).

#### P3 — Repo row exists with empty clone URL

- **Preconditions:** insert a `Repo` row whose `HTTPSUrl == ""` and
  `SSHUrl == ""` (slug `orphan`). No fixture bare repo needed.
- **CLI:** `gitmap-e2e probe <slug-path>` where `<slug-path>` resolves
  to that repo via `db.FindByPath`.
- **Stdout:** contains `fmt.Sprintf(constants.MsgProbeDoneFmt, 0, 0, 1)`.
- **Stderr:** contains `fmt.Sprintf(constants.ErrProbeMissingURL, "orphan")`.
- **Exit code:** `0` (the loop tallies a failure but does not abort).
- **DB delta:** one new `VersionProbe` row with
  `Method == constants.ProbeMethodNone`, `IsAvailable == 0`,
  `Error == fmt.Sprintf(constants.ErrProbeMissingURL, "orphan")`.
- **Cleanup:** §4.3 invariant; no `git` subprocess was spawned (assert
  by wrapping `PATH` to a `git` shim that records invocations).

#### P4 — Malformed URL

- **Preconditions:** DB has a `Repo` row whose `HTTPSUrl == "not-a-url"`.
- **CLI:** `gitmap-e2e probe <that-repo-path>`
- **Stdout:** contains `fmt.Sprintf(constants.MsgProbeFailFmt, <slug>, <error-substring>)`
  and `fmt.Sprintf(constants.MsgProbeDoneFmt, 0, 0, 1)`.
- **Stderr:** empty (per-repo errors go to stdout via `MsgProbeFailFmt`).
- **Exit code:** `0`. **MUST NOT** panic (asserted by the absence of
  `panic:` in combined output).
- **DB delta:** one new `VersionProbe` row, `IsAvailable == 0`, `Error`
  non-empty, `Method` ∈ {`ls-remote`, `shallow-clone`}.
- **Cleanup:** §4.3 invariant; any temp clone dir created during the
  shallow-clone attempt is removed.

#### P5 — Annotated-tag dereference (`v1.0.0^{}`)

- **Preconditions:** build the bare repo manually (helper variant
  `fixture.NewBareRepoAnnotated(t, "v1.0.0")`) using
  `git tag -a v1.0.0 -m x`. Confirm via
  `git ls-remote ${repo.URL}` that output contains both `refs/tags/v1.0.0`
  and `refs/tags/v1.0.0^{}`.
- **CLI:** `gitmap-e2e probe ${repo.URL}`
- **Stdout:** contains `"v1.0.0"` (the `^{}` suffix **MUST NOT** appear).
- **Exit code:** `0`.
- **DB delta:** `VersionProbe.NextVersionTag == "v1.0.0"` exactly.
- **Cleanup:** §4.3 invariant.

#### P6 — Pre-release sort order

- **Preconditions:** `fixture.NewBareRepo(t, "v1.0.0", "v1.0.1-rc1", "v1.0.1")`.
- **CLI:** `gitmap-e2e probe ${repo.URL} --json`
- **Stdout:** valid JSON array of length 1; the entry has
  `nextVersionTag == "v1.0.1"` and `nextVersionNum == 1`.
- **Stderr:** empty.
- **Exit code:** `0`.
- **DB delta:** one `VersionProbe` row matching the JSON.
- **Cleanup:** §4.3 invariant.

### 4.2 Shallow-clone fallback (MUST)

#### P7 — ls-remote fails, shallow-clone is reached and also fails

- **Preconditions:** `dir := filepath.Join(t.TempDir(), "notagit")`,
  `os.MkdirAll(dir, 0o755)`. DB has a `Repo` row with
  `HTTPSUrl == "file://" + dir`.
- **CLI:** `gitmap-e2e probe <that-repo-path>`
- **Stdout:** contains `fmt.Sprintf(constants.MsgProbeFailFmt, <slug>, <err>)`.
- **Stderr:** empty.
- **Exit code:** `0`.
- **DB delta:** one `VersionProbe` row with
  `Method == constants.ProbeMethodShallowClone`, `IsAvailable == 0`,
  and `Error` matching `constants.ErrProbeCloneFail` format
  (assert with `strings.HasPrefix` after stripping the `%v`).
- **Cleanup:** §4.3 invariant — *strictly* zero `gitmap-probe-*`
  entries in `os.TempDir()` after the run, even though shallow-clone
  created one mid-flight. Test fails if any remain.

#### P8 — ls-remote succeeds with zero tags; shallow-clone MUST NOT run

- **Preconditions:** `repo := fixture.NewBareRepoNoTags(t)`. Wrap `git`
  on `PATH` with a shim that records every invocation.
- **CLI:** `gitmap-e2e probe ${repo.URL}`
- **Stdout:** as P2.
- **Exit code:** `0`.
- **DB delta:** one `VersionProbe` row, `Method == ls-remote`,
  `IsAvailable == 0`.
- **Cleanup:** the recorded git shim log **MUST NOT** contain a
  `clone` invocation. §4.3 invariant holds.

### 4.3 Temp-dir cleanup invariant (MUST)

#### P9 — Zero leak across the full P1–P8 matrix

- **Preconditions:** `before := fixture.SnapshotTempProbeDirs()` taken
  in `TestMain` before any P-scenario runs.
- **Mechanism:** `assertNoTempLeak(t)` is registered via `t.Cleanup`
  inside the shared `runProbeCLI` helper, so every P1–P8 test
  re-asserts the invariant on its own.
- **CLI:** N/A (cross-cutting).
- **Stdout/stderr/exit:** N/A.
- **DB delta:** none beyond what each individual scenario records.
- **Cleanup assertion:** `after := fixture.SnapshotTempProbeDirs()`
  in `TestMain`'s teardown. `len(after) == len(before)` and the
  set difference is empty. Any extra entry fails the suite with the
  offending paths printed.

### 4.4 Optional (SHOULD)

- **P10**: concurrent `RunOne` calls on the same fixture do not interfere.
- **P11**: very long tag list (1000 tags) returns the expected top in
  under a hard timebox (e.g. 2s on the CI runner).

---

## 5. Cloner direct-URL e2e tests

Under test: the direct-URL path in `gitmap-v27/cloner/` (see
`spec/01-app/88-clone-direct-url.md`).

### 5.1 Required test classes (MUST)

| ID | Scenario | Fixture | Expected |
|----|----------|---------|----------|
| C1 | HTTPS-style URL → folder name derived | `https://github.com/owner/my-repo.git` | Folder = `my-repo` |
| C2 | URL with `.git` suffix and without | both forms | Same folder name |
| C3 | SCP-style `git@host:owner/repo.git` | literal | Folder = `repo`, classified as URL |
| C4 | URL with `:branch` suffix | `https://.../repo:develop` | URL stripped, branch surfaced |
| C5 | Successful clone into `--target-dir` | `NewBareRepo("v1.0.0")` URL | Working tree exists, `.git/` present |
| C6 | Target folder already exists, not git | pre-create dir | Exits with error, no partial clone |
| C7 | Target folder already exists, IS git, cache hit | run C5 twice | Second call short-circuits, prints `skipped (cached)` |
| C8 | Custom folder-name override | URL + `--folder my-alias` | Clone lands at `<target>/my-alias` |

### 5.2 DB upsert verification (MUST)

ID **C9**: after C5 succeeds, open the SQLite DB created in `<target>/.gitmap/`,
query the `Repository` table, and assert exactly one row whose
`HTTPSUrl` matches the fixture URL. The test **MUST** use the same DB
helper the production code uses — no raw SQL in the test body.

### 5.3 Audit-mode parity (SHOULD)

ID **C10**: invoke the cloner with `--audit` against a manifest that
references the bare-repo URL. Assert the printed report classifies the
record as `clone (+)` before C5 and `cached (=)` after C5, and that
`--audit` writes nothing to disk and makes no `git` invocation
(check via a fake `PATH` that errors if `git` is called).

---

## 6. Install-script e2e tests

Under test: `scripts/install.sh` and `scripts/install.ps1` against the
contract in [`spec/07-generic-release/09-generic-install-script-behavior.md`](../07-generic-release/09-generic-install-script-behavior.md).

Tests are written in Go (so they run alongside the rest of the suite)
but invoke the scripts via `exec.Command("bash", scriptPath, ...)` /
`exec.Command("pwsh", scriptPath, ...)`. The scripts **MUST** be made
testable by allowing `GITMAP_RELEASE_BASE_URL` (or equivalent) to be
overridden to point at a local HTTP server fixture.

### 6.1 Local release-server fixture

Create `tests/internal/fixture/relsrv.go`:

```go
type ReleaseServer struct {
    URL      string                  // base URL of the test server
    Releases map[string][]byte       // tag -> tarball bytes
    HEAD     map[string]int          // path -> status (for sibling probe)
}

func NewReleaseServer(t *testing.T) *ReleaseServer
func (s *ReleaseServer) AddRelease(tag string, payload []byte)
func (s *ReleaseServer) SetSiblingProbeStatus(version string, status int)
```

The server **MUST** respond to:

- `HEAD /releases/tag/v<N+i>` → status from `HEAD` map (default 404)
- `GET  /releases/download/<tag>/<asset>` → bytes from `Releases`
- `GET  /releases/latest` → 302 redirect to highest registered tag
- `GET  /raw/<branch>/...` → main-HEAD fallback assets

### 6.2 Required test classes (MUST)

| ID | Mode | Setup | Expected |
|----|------|-------|----------|
| I1 | Strict tag | `--version v3.0.0`, server has v3.0.0 | Installs v3.0.0, exit 0, no probe traffic |
| I2 | Strict tag, missing | `--version v9.9.9`, server returns 404 | Exit 1, **MUST NOT** probe siblings, **MUST NOT** fall back |
| I3 | Discovery, sibling hit | no `--version`, current = v3.0.0, server returns 200 for `v3.0.4`, 404 for v3.0.1..3 and v3.0.5..20 | Installs v3.0.4 (max sibling hit) |
| I4 | Discovery, no siblings | all 20 HEADs 404, latest-release endpoint returns v3.0.0 | Installs v3.0.0 |
| I5 | Discovery, no release at all | siblings 404, `/releases/latest` 404 | Falls back to main HEAD raw assets, exit 0 |
| I6 | Discovery, partial sibling failures | some HEADs 500, others 404, one 200 at v3.0.7 | Installs v3.0.7; 500 responses **MUST NOT** be treated as success |

### 6.3 Parallelism invariants (MUST)

ID **I7**: instrument the test server to record arrival timestamps for
the 20 sibling HEAD requests. Assert that the spread between first and
last arrival is below a wall-clock threshold (e.g. 500ms) — proving the
20 probes ran in parallel, not serially.

ID **I8**: assert the script issues **exactly 20** sibling probes when
no early hit shortcuts the loop, and **at most 20** when an early hit
occurs (the spec allows but does not require cancellation of in-flight
probes).

### 6.4 Cross-shell parity (SHOULD)

ID **I9**: every test in §6.2 **SHOULD** run twice — once against
`install.sh` under `bash`, once against `install.ps1` under `pwsh` —
with identical assertions. CI may skip the pwsh leg on platforms where
PowerShell is unavailable, but **MUST** record the skip explicitly.

---

## 7. Shared invariants across all three suites

The following invariants **MUST** hold for every test in §4–§6:

1. **No network.** A test that issues a DNS lookup for a public host
   fails. Enforce via `GIT_ALLOW_PROTOCOL=file` and a custom HTTP
   transport that rejects non-loopback addresses.
2. **No global state.** Every test uses `t.TempDir()` and its own
   fixture instance. No test reads or writes `$HOME`, `$PWD`, or any
   shared cache directory.
3. **Deterministic timing.** No `time.Sleep` over 50ms. Use channels
   or `t.Eventually`-style polling.
4. **Zero leaks.** `t.Cleanup` removes every temp dir, kills every
   spawned process, and closes every server. The §4.3 temp-dir guard
   applies to all three suites.
5. **Error messages are asserted, not just types.** When the spec
   prescribes a user-facing message (e.g. `constants.ErrProbeCloneFail`
   formatting), the test **MUST** assert against the constant — not a
   hard-coded literal — so message changes update the constant in one
   place.

---

## 8. Constants and fixtures registry

To keep test code free of magic strings (per the project-wide
constants policy), introduce:

- `tests/internal/constants/constants_test.go`
  - `TagV1_0_0 = "v1.0.0"` etc. for any tag mentioned in ≥2 tests
  - `FolderMyRepo = "my-repo"`
  - `URLOwnerRepo = "https://github.com/owner/my-repo.git"`
- `tests/internal/fixture/probe_payloads.go`
  - canned `ls-remote` outputs for parser tests

No test file may inline a tag string or URL that appears in another
test file — promote it to the registry instead.

---

## 9. Running the suites

```bash
# Unit + existing tests stay where they are.
go test ./...

# E2E suites — slower, opt-in flag for local dev iteration.
go test -tags=e2e ./tests/e2e_probe_test/...
go test -tags=e2e ./tests/e2e_cloner_test/...
go test -tags=e2e ./tests/e2e_install_test/...

# CI runs everything.
go test -tags=e2e -race ./...
```

The `e2e` build tag **MUST** gate every file in `tests/e2e_*_test/`.
This keeps `go test ./...` fast for everyday work while ensuring CI
runs the whole matrix.

---

## 10. Acceptance checklist

A PR adding or modifying URL-handling code is acceptable only if:

- [ ] `tests/internal/fixture/` exposes `NewBareRepo`,
      `NewEmptyBareRepo`, `NewBareRepoNoTags`, and `NewReleaseServer`.
- [ ] All probe scenarios P1–P9 are implemented and pass.
- [ ] All cloner scenarios C1–C9 are implemented and pass.
- [ ] All install scenarios I1–I8 are implemented and pass for `bash`.
- [ ] §7 invariants are enforced via shared helpers, not copy-pasted
      per test.
- [ ] No test issues a real-network request (verified by CI network
      isolation or transport guard).
- [ ] No test file contains a tag or URL literal that appears in
      another test file.
- [ ] CI workflow runs `go test -tags=e2e -race ./...` on at least one
      Linux runner and one Windows runner.

---

## 11. Open extension points

These items are intentionally deferred but documented so a future AI
agent can pick them up without re-deriving context:

- **Mock SSH server** for SCP-style URL coverage beyond C3 (currently
  classification-only). Requires an embedded SSH server fixture.
- **Flaky-network simulator** that injects 1% packet loss into the
  release server to validate retry/backoff in install scripts once
  retries are added.
- **chocolatey/winget package install tests** — out of scope here per
  §1, but should follow the same fixture pattern when added.

---

## Appendix A — SHOULD-coverage backlog (deferred scenarios)

This appendix promotes previously hand-waved "nice to have" items into
specified, paste-ready scenarios. Every entry follows the same five-
section contract used in §4 (preconditions, CLI, stdout/stderr/exit,
DB delta, cleanup). Items here are **SHOULD**, not **MUST** — a host
PR may defer any of them with written justification, but once
implemented they belong under `tests/e2e_*_test/` next to the MUST
suites and **MUST NOT** weaken the §7 invariants.

IDs use the next free number per layer (P12+, C11+, I10+) so existing
references stay stable.

### A.1 Concurrent probe races

#### P12 — Parallel `RunOne` against the same fixture

- **Preconditions:** `repo := fixture.NewBareRepo(t, "v1.0.0", "v1.0.5", "v1.0.20")`.
  DB row for `repo.URL` exists. Spawn `N := 16` goroutines.
- **CLI:** each goroutine runs `gitmap-e2e probe ${repo.URL} --json`
  via `runProbeCLI`; launched simultaneously through a
  `sync.WaitGroup` + start barrier (`close(start)`).
- **Stdout/stderr/exit:** every invocation exits `0`; every JSON
  payload reports `nextVersionTag == "v1.0.20"`,
  `method == "ls-remote"`, `isAvailable == true`.
- **DB delta:** exactly `N` new `VersionProbe` rows for that `RepoId`,
  all with the same `NextVersionTag`. Assert via
  `db.CountVersionProbes(repoID) == before + N`. No row has a NULL
  `Method` or partially-written `Error`.
- **Cleanup:** §4.3 invariant; the bare repo's tree hash is unchanged
  (`fixture.HashTree(repo.Dir)` matches its pre-test snapshot —
  proves no goroutine corrupted the shared remote).

#### P13 — Concurrent probes across distinct fixtures (no cross-talk)

- **Preconditions:** create 4 fixtures with disjoint tag sets
  (`v1.*`, `v2.*`, `v3.*`, `v4.*`). DB has one `Repo` row per fixture.
- **CLI:** `gitmap-e2e probe --all --json` (single invocation; the
  internal probe loop is currently sequential per §99 of the
  v3.8.0 plan, so this test is a regression guard for the planned
  Phase 2.5 worker pool).
- **Expected:** JSON array length 4; each entry's `nextVersionTag`
  matches its fixture's max tag.
- **DB delta:** 4 new `VersionProbe` rows, one per `RepoId`. No row
  carries a tag from another fixture's namespace.
- **Cleanup:** §4.3 invariant; no `gitmap-probe-*` directory survived;
  no `git` shim invocation referenced a URL outside its fixture.

### A.2 Partial-clone resume / interruption

#### P14 — Shallow-clone interrupted mid-flight, retried successfully

- **Preconditions:** `repo := fixture.NewBareRepo(t, "v1.0.0")`.
  Wrap `git` on `PATH` with a shim that, on the **first** `clone`
  invocation, writes a partial `.git/objects/pack/tmp_pack_*` file
  into the destination and exits with `signal: killed` (exit 137).
  Subsequent invocations pass through to real `git`.
- **CLI:** run `gitmap-e2e probe ${repo.URL}` twice in sequence.
- **Expected (first run):** stdout contains
  `MsgProbeFailFmt`; exit code `0`; `VersionProbe.Error` matches
  `constants.ErrProbeCloneFail`. The orphan partial dir under
  `os.TempDir()` **MUST** still be cleaned up by the production code's
  `defer os.RemoveAll`.
- **Expected (second run):** stdout contains `MsgProbeOkFmt` for
  `v1.0.0`; new `VersionProbe` row with `IsAvailable == 1`.
- **DB delta:** exactly two new `VersionProbe` rows total — one
  failure, one success. No row references the partial-pack path.
- **Cleanup:** §4.3 snapshot delta = 0 across both runs; the git shim
  log shows exactly two `clone` invocations.

#### C11 — Cloner resumes after a failed first attempt leaves a stub dir

- **Preconditions:** pre-create `<target>/my-repo/` with a single empty
  `.git/` directory (no `HEAD`, no objects) — simulates a prior
  killed clone.
- **CLI:** `gitmap-e2e clone ${repo.URL} --target-dir <target>`
- **Expected stdout:** contains the cloner's "stale partial directory
  detected, removing" message (must be promoted to a constant if not
  already — `constants.MsgClonerStalePartial`).
- **Exit code:** `0`.
- **DB delta:** one `Repository` row inserted (or upserted) with the
  fixture URL.
- **Cleanup:** `<target>/my-repo/` contains a valid working tree
  (`git -C <target>/my-repo rev-parse HEAD` exits `0`); no
  `gitmap-probe-*` or `gitmap-clone-*` temp dir remains.

### A.3 Malformed tag handling

#### P15 — Non-semver tag in remote MUST be ignored, not crash

- **Preconditions:** custom fixture
  `fixture.NewBareRepoRawTags(t, "v1.0.0", "release-candidate", "1.0", "v1.0.5", "vNEXT")`.
  Only `v1.0.0` and `v1.0.5` are valid `vMAJOR.MINOR.PATCH` tags.
- **CLI:** `gitmap-e2e probe ${repo.URL} --json`
- **Expected JSON:** `nextVersionTag == "v1.0.5"`, `nextVersionNum == 5`.
  The malformed tags **MUST NOT** appear in stdout/stderr.
- **Stderr:** empty.
- **Exit code:** `0`. **MUST NOT** panic.
- **DB delta:** one `VersionProbe` row, `NextVersionTag == "v1.0.5"`.
- **Cleanup:** §4.3 invariant.

#### P16 — Tag containing shell metacharacters MUST be inert

- **Preconditions:** fixture seeded with tags `v1.0.0`,
  `v1.0.1; rm -rf /`, `v1.0.2$(whoami)`. (`git` accepts these as long
  as the names follow `git check-ref-format`; for any rejected by git
  the fixture helper falls back to writing `refs/tags/<name>` directly
  into the bare repo's `packed-refs` file.)
- **CLI:** `gitmap-e2e probe ${repo.URL}`
- **Expected:** highest **valid-semver** tag wins (`v1.0.2$(whoami)`
  is non-semver per §A.3/P15 and ignored). Top result is `v1.0.0`
  (the only clean tag) **or** `v1.0.1; rm -rf /` if the parser
  considers a leading `v1.0.1` valid — in either case **no shell
  expansion occurs**: assert that `/` is unmodified
  (`os.Stat("/tmp/pwned")` returns `ErrNotExist`) and that no
  `whoami` subprocess fired (git shim log).
- **Exit code:** `0`. **MUST NOT** panic.
- **DB delta:** one `VersionProbe` row; `NextVersionTag` value is
  written verbatim (no truncation at `;` or `$`).
- **Cleanup:** §4.3 invariant; PATH-shim log contains zero
  `sh`/`bash`/`whoami` invocations.

### A.4 Windows long-path probe directories

These scenarios **MUST** be gated by a `runtime.GOOS == "windows"`
guard and skipped with `t.Skip("windows-only")` elsewhere. They
exercise the §3 fixtures under paths that exceed `MAX_PATH` (260)
once the probe's temp-dir suffix is appended.

#### P17 — Probe under a >260-char temp root succeeds

- **Preconditions:** force the probe's temp root via
  `t.Setenv("TMP", longRoot)` and `t.Setenv("TEMP", longRoot)`,
  where `longRoot` is constructed as
  `filepath.Join(t.TempDir(), strings.Repeat("d", 240))` and created
  with the `\\?\` long-path prefix
  (`os.MkdirAll(`\\?\` + longRoot, 0o755)`). Confirm
  `len(longRoot) + len("gitmap-probe-XXXXXXXX") > 260`.
  Fixture: `repo := fixture.NewBareRepo(t, "v1.0.0")`.
  **Skip condition:** if `core.longpaths` is not enabled in the test's
  scoped git config, `t.Skip` with a clear message — the failure
  belongs to the environment, not the code under test.
- **CLI:** `gitmap-e2e probe ${repo.URL}`
- **Stdout:** contains `MsgProbeOkFmt` with `v1.0.0`.
- **Exit code:** `0`.
- **DB delta:** one `VersionProbe` row, `IsAvailable == 1`.
- **Cleanup:** §4.3 invariant — `os.RemoveAll` succeeded against the
  long path. Test fails if any `gitmap-probe-*` entry remains under
  `longRoot`.

#### P18 — Shallow-clone fallback under a >260-char temp root

- **Preconditions:** as P17 but the URL points at a non-git directory
  (the P7 setup) so the shallow-clone branch fires.
- **CLI:** `gitmap-e2e probe ${repo.URL}`
- **Expected:** stdout contains `MsgProbeFailFmt`; exit `0`;
  `VersionProbe.Method == constants.ProbeMethodShallowClone`,
  `IsAvailable == 0`, `Error` matches `constants.ErrProbeCloneFail`.
- **Cleanup:** the long-path temp dir created by `tryShallowClone` is
  removed; §4.3 invariant holds. Test fails with the offending path
  if cleanup failed (this is the regression we care about — Windows
  long-path `os.RemoveAll` historically silently leaks).

#### C12 — Cloner target directory exceeds MAX_PATH

- **Preconditions:** windows-only; pre-build a target dir whose path
  is 250 chars and let the clone produce nested files that push the
  full path past 260. Requires `git config --global core.longpaths true`
  in the test's scoped HOME.
- **CLI:** `gitmap-e2e clone ${repo.URL} --target-dir <longTarget>`
- **Expected:** exit `0`; working tree present; no
  `filename too long` in stderr.
- **DB delta:** one `Repository` row.
- **Cleanup:** removable via `os.RemoveAll(\\?\<longTarget>)`; no
  orphan temp directories.

### A.5 Coverage matrix update

When any A-scenario lands, add a row to the §10 acceptance checklist
under a new "SHOULD coverage" sub-section. The MUST list **MUST NOT**
be modified — SHOULD items remain optional gates so a host PR can
ship MUST coverage independently.
