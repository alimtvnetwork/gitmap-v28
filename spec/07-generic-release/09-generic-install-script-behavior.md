# 09 — Generic Install-Script Behavior (AI-Shareable Contract)

> **Status:** Authoritative draft (2026-04-22).
> **Audience:** Any AI agent or human maintainer implementing or auditing
> an installation script in any repository. This document is intentionally
> **repository-agnostic**. Substitute `<owner>`, `<stem>`, `<binary>`, and
> `<installerPath>` with values from the host repo.
> **Supersedes (in spirit, not by deletion):** the strict-tag clause in
> [`08-pinned-version-install-snippet.md`](08-pinned-version-install-snippet.md)
> and the discovery section of
> [`../01-app/95-installer-script-find-latest-repo.md`](../01-app/95-installer-script-find-latest-repo.md).
> Those documents remain valid; this spec is the single normative summary
> a foreign AI should be handed.
> **Tested by:** [`../04-generic-cli/25-e2e-testing-probe-clone-install.md`](../04-generic-cli/25-e2e-testing-probe-clone-install.md) — §6 covers strict-tag, sibling probe, latest-release, and main-HEAD fallback.

---

## 0. How to use this document

When asked to implement or audit an installer in **any** repository:

1. Read this file end-to-end.
2. Identify the host repo's `<owner>`, `<stem>`, `<binary>`, branch name,
   and installer path.
3. Apply every clause marked **MUST**. Clauses marked **SHOULD** may be
   skipped only with a written justification in the host repo's spec.
4. Mirror the §11 acceptance checklist in the host repo's CI.

The keywords **MUST**, **MUST NOT**, **SHOULD**, **MAY** follow RFC 2119.

---

## 1. Scope

This spec governs every script whose purpose is to download and install
the host project's binary or runtime onto an end-user machine. Examples
of in-scope scripts (names are illustrative — substitute the host repo's):

| Script class           | Typical filename                      | Invocation style                          |
|------------------------|---------------------------------------|-------------------------------------------|
| Quick / one-liner      | `install-quick.sh`, `install-quick.ps1` | `curl … \| bash`, `irm … \| iex`          |
| Canonical / release    | `scripts/install.sh`, `scripts/install.ps1` | Pinned URL with full flag surface     |
| Feature-specific       | `scripts/error-manage.sh`, etc.       | Same contract as canonical                |
| Bootstrap / re-install | `scripts/install-bootstrap.*`         | Wrapper that re-execs canonical           |

**Out of scope:** in-place self-update inside an already-installed
binary (see [`spec/08-generic-update/`](../08-generic-update/)), package-
manager publishing (Homebrew, winget, apt), and per-language module
installers (`go install`, `npm i -g`).

---

## 2. Core decision tree

Every install script **MUST** make exactly one decision at startup, in
this priority order:

```
            ┌────────────────────────────────────────────┐
            │  Was an explicit version supplied?         │
            │  (--version <tag>, -Version <tag>,         │
            │   $env:VERSION, or release-page snippet)   │
            └────────────────────────────────────────────┘
                  │ yes                    │ no
                  ▼                        ▼
         §3 STRICT TAG MODE         §4 DISCOVERY MODE
       (no fallbacks, no probe)   (20-parallel sibling probe,
                                   then latest release,
                                   then main HEAD as last resort)
```

Implementations **MUST NOT** blend the two paths. Once strict-tag mode
is selected, no discovery, no main-branch fallback, no "nearest tag"
heuristic may execute.

---

## 3. Strict tag mode (MANDATORY semantics)

A tag is **explicit** if any of the following is true at script entry:

* CLI flag `--version <tag>` (bash) or `-Version <tag>` (PowerShell) is
  set to a non-empty value.
* Environment variable `VERSION` (or repo-prefixed equivalent, e.g.
  `<STEM>_VERSION`) is set to a non-empty value.
* The script was rendered from a release-page snippet (see
  [`08-pinned-version-install-snippet.md`](08-pinned-version-install-snippet.md))
  and its template substituted a literal tag at publish time.

When strict-tag mode is active:

1. **MUST** download from `…/releases/download/<tag>/…` directly.
2. **MUST NOT** call `…/releases/latest`.
3. **MUST NOT** invoke the §4 versioned-repo discovery probe — no
   `-v<N+1>` HEAD requests, no sibling lookup, nothing.
4. **MUST NOT** fall back to the main branch HEAD under any error.
5. **MUST NOT** silently substitute a different tag (no "nearest semver",
   no "did you mean v3.11.2?" auto-pick).
6. **MUST** print the resolved URL exactly once before downloading:
   ```
     [strict] requested tag: <tag>
     [strict] download: https://github.com/<owner>/<stem>/releases/download/<tag>/<asset>
   ```
7. On any failure (HTTP 404, checksum mismatch, missing asset), **MUST**
   exit with code `1` and the canonical message:
   ```
     Error: requested release <tag> not found in <owner>/<stem>;
     refusing to fall back per strict-tag contract.
     See spec/07-generic-release/09-generic-install-script-behavior.md §3.
   ```
8. **MUST NOT** delete or modify any pre-existing installation when
   strict-tag mode aborts — leave the user's prior binary untouched.

### 3.1 Rationale

A pinned tag is a **contract**. Users copying an install snippet from
`…/releases/tag/v3.11.1` are asking for v3.11.1 specifically — typically
to reproduce a bug, lock a CI job, or roll back. Any silent fallback
violates that contract and converts a deterministic install into a
guessing game.

---

## 4. Discovery mode (no explicit tag)

When **no** tag was supplied, the installer resolves the install target
in three ordered phases. Each phase is attempted only if the previous
phase yielded no result.

### 4.1 Phase A — Versioned-repo discovery (20 parallel HEADs)

Applies only if the host repo name ends in `-v<N>` where `N` is a
positive integer (per §3 of the legacy
[`95-installer-script-find-latest-repo.md`](../01-app/95-installer-script-find-latest-repo.md)).

```
INPUT:  baselineUrl = https://github.com/<owner>/<stem>-v<N>
        windowSize  = 20      (configurable: --discovery-window <K>)

candidates = [ <owner>/<stem>-v<M>  for M in (N+1) .. (N+windowSize) ]

# Issue all HEAD requests CONCURRENTLY.
results = parallelMap(candidates, head, timeout=5s, maxConcurrency=20)

# Pick the highest M whose HEAD returned 200.
hits        = [ M  for (M, status) in results if status == 200 ]
effectiveN  = max(hits) if hits else N
```

Notes:

* **MUST** issue all probes concurrently (not sequentially) so total
  wall time is bounded by ~one HEAD round-trip + 5 s timeout, not by
  `windowSize × RTT`.
* **MUST NOT** assume version contiguity. The previous fail-fast rule
  in `95-installer-script-find-latest-repo.md` §4 is **superseded** by
  this max-hit-wins rule. Gaps (`v4`, `v6` exists but `v5` does not)
  are tolerated and `v6` wins.
* **MUST** cap concurrency at `min(windowSize, 20)` to stay polite to
  GitHub's anonymous rate limit (60 req/h/IP). With `GITHUB_TOKEN`
  present, **MAY** raise to 50.
* **MUST** log every probe result with the `[discovery]` prefix
  (see §7 of `95-installer-script-find-latest-repo.md`). Concurrent
  logs may interleave; that is acceptable.
* If `effectiveN > N`, **MUST** delegate to the effective repo's
  installer (see §5) and exit with the delegated process's status.
  Re-entry **MUST** be guarded by `INSTALLER_DELEGATED=1`.
* If the host repo name has no `-v<N>` suffix, **MUST** skip Phase A
  entirely and proceed to Phase B.

### 4.2 Phase B — Latest published release

On the effective repo (post-discovery), the installer **MUST**:

1. Call `GET https://api.github.com/repos/<owner>/<stem>/releases/latest`.
2. If a non-prerelease release exists, install its tag using the same
   asset-resolution logic as §3 (but **without** the strict-tag
   fail-fast — proceed to Phase C on missing assets).
3. If the API returns 404 (no published releases) or every release is
   marked prerelease/draft, proceed to Phase C.

### 4.3 Phase C — Main branch HEAD (last resort)

Applies **only** when Phases A and B yielded no installable artifact
(typically: brand-new repo with no tagged release).

* **MUST** install from the `main` branch (or the repo's default branch,
  detected via `GET /repos/<owner>/<stem>` → `default_branch`) by
  building from source or downloading the canonical archive of the
  branch tip.
* **MUST** print a prominent warning:
  ```
    [warn] no published releases found; installing from <branch> HEAD.
    [warn] this is unstable; pin a version with --version <tag> for production.
  ```
* **MUST** record the resolved commit SHA in the post-install summary
  so the install is reproducible.

---

## 5. Delegation contract

When Phase A picks `effectiveN > N`, the **original** installer:

1. Constructs the delegated URL by replacing the stem suffix:
   `https://raw.githubusercontent.com/<owner>/<stem>-v<effectiveN>/<branch>/<installerPath>`
2. Re-invokes that script, **passing through every original flag verbatim**
   (`--dir`, `--no-path`, `--probe-ceiling`, `--discovery-window`, etc.).
   `--version` is never present here (otherwise §3 would have short-
   circuited before §4).
3. Sets `INSTALLER_DELEGATED=1` in the child env. The child **MUST**
   read this and skip §4 Phase A on re-entry. Phases B and C still run.
4. Exits with the child's exit code. **MUST NOT** continue executing
   any post-delegation code in the parent.

---

## 6. Configuration surface (canonical flag names)

All installers in the host repo **SHOULD** expose this exact flag
surface. Synonyms are permitted but the canonical names below **MUST**
be accepted.

| Bash flag                | PowerShell param        | Default | Effect |
|--------------------------|-------------------------|---------|--------|
| `--version <tag>`        | `-Version <tag>`        | (none)  | Activates §3 strict-tag mode. |
| `--dir <path>`           | `-InstallDir <path>`    | OS-default | Override install location. |
| `--arch <arch>`          | `-Arch <arch>`          | autodetect | Override CPU arch. |
| `--no-path`              | `-NoPath`               | off     | Skip PATH registration. |
| `--no-discovery`         | `-NoDiscovery`          | off     | Skip §4 Phase A only. Phases B + C still run. |
| `--discovery-window <K>` | `-DiscoveryWindow <K>`  | `20`    | Number of `-v<N+i>` siblings to probe in parallel. Cap = 20 unless `GITHUB_TOKEN` set (then 50). |
| `--source main\|latest`  | `-Source main\|latest`  | `latest` | Force Phase B or skip directly to Phase C. Strict-tag mode ignores this. |
| `INSTALLER_DELEGATED=1`  | (env)                   | unset   | Loop guard; set automatically by §5 delegation. |
| `GITHUB_TOKEN`           | (env)                   | unset   | Optional auth for higher rate limits. |

---

## 7. Logging contract

* All log lines for discovery **MUST** start with `[discovery]`.
* All log lines for strict-tag mode **MUST** start with `[strict]`.
* All warnings **MUST** start with `[warn]`.
* Errors **MUST** be written to stderr, prefixed `Error:`, and
  contain a pointer to this spec's section number when terminating.
* Post-install summary **MUST** include: binary path, resolved version
  (or commit SHA for Phase C), source repo, and PATH change status.

---

## 8. Failure handling matrix

| Scenario                                      | Mode              | Action                                                                 |
|-----------------------------------------------|-------------------|------------------------------------------------------------------------|
| Strict tag asset 404                          | §3                | exit 1, canonical message, no fallback                                 |
| Strict tag checksum mismatch                  | §3                | exit 1, canonical message, no retry                                    |
| Discovery HEAD timeout                        | §4 A              | treat as MISS, do not abort                                            |
| Discovery rate-limited (HTTP 403)             | §4 A              | treat as MISS, log hint about `GITHUB_TOKEN`                           |
| All Phase-A probes MISS                       | §4 A              | use baseline repo, proceed to Phase B                                  |
| Phase B `releases/latest` returns 404         | §4 B              | proceed to Phase C with `[warn]`                                       |
| Phase C source archive 404                    | §4 C              | exit 1; no further fallback exists                                     |
| Delegated child exits non-zero                | §5                | parent exits with same code, prints `[discovery] child exited <code>`  |
| `INSTALLER_DELEGATED=1` set but no `-v<N>`    | §4 A              | log `[discovery] loop guard active; skipping`, continue to Phase B     |

---

## 9. Security considerations

* Checksum verification (SHA-256) is **MANDATORY** for every downloaded
  archive in Phases §3 and §4.B. **MUST NOT** be skipped via flag.
* Phase C (main HEAD) **SHOULD** record commit SHA but **MAY** skip
  checksum (no asset to compare against). Print `[warn] unverified source`.
* Scripts **MUST NOT** auto-elevate (no implicit `sudo`, no UAC prompt).
  If write to default dir fails, fall back to `~/.local/bin` (Unix) or
  `%LOCALAPPDATA%\<binary>` (Windows) and report the new path.
* Scripts **MUST NOT** persist `GITHUB_TOKEN` to disk or config files.

---

## 10. Acceptance checklist (mirror in host-repo CI)

An installer conforms when:

- [ ] Detects `--version` / `-Version` and short-circuits to §3.
- [ ] In §3, never issues a `releases/latest` API call.
- [ ] In §3, never issues any `-v<N+i>` HEAD probe.
- [ ] In §3, exits 1 with the canonical message on any failure.
- [ ] In §4 A, issues up to 20 HEAD requests concurrently.
- [ ] In §4 A, picks the **maximum** HIT (gaps tolerated).
- [ ] In §4 A, sets `INSTALLER_DELEGATED=1` before re-execing.
- [ ] In §4 B, uses `releases/latest`; in §4 C, uses default branch.
- [ ] Logs use `[strict]` / `[discovery]` / `[warn]` prefixes.
- [ ] Honors `--no-discovery`, `--discovery-window`, `--source`,
  `--dir`, `--arch`, `--no-path`.
- [ ] Negative test: `--version v0.0.0-nope` → exit 1, no probe traffic.
- [ ] Negative test: missing `main` branch and no releases → exit 1
  in Phase C, not silent success.

---

## 11. Cross-references

* [`03-install-scripts.md`](03-install-scripts.md) — installer mechanics
  (download, checksum, PATH).
* [`08-pinned-version-install-snippet.md`](08-pinned-version-install-snippet.md)
  — release-page snippet that triggers §3 strict-tag mode.
* [`../01-app/95-installer-script-find-latest-repo.md`](../01-app/95-installer-script-find-latest-repo.md)
  — original sequential-probe design; §4 of the present spec
  **supersedes** its fail-fast clause with the 20-parallel rule.
* [`../08-generic-update/`](../08-generic-update/) — in-place
  self-update (different concern: same binary, new release).

---

## 12. History

| Date       | Change                                                                                  |
|------------|-----------------------------------------------------------------------------------------|
| 2026-04-22 | Initial draft. Unifies strict-tag contract, 20-parallel discovery, and main-HEAD fallback into one AI-shareable document. |
