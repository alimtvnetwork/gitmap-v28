# Plan 07 — Generic Install-Script Behavior Rollout

**Spec of record:** `spec/07-generic-release/09-generic-install-script-behavior.md`
**Memory:** `mem://features/generic-install-behavior`
**Created:** 2026-04-22
**Status:** Spec written; code changes NOT YET applied (awaiting user approval).

---

## Goal

Bring all four install scripts in this repo into compliance with the new
canonical contract, and provide a reusable plan that any sibling repo
(`<stem>-v<M>`) can follow.

## Affected files (this repo)

| File                              | Current behavior                     | Required change                                                                 |
|-----------------------------------|--------------------------------------|---------------------------------------------------------------------------------|
| `install-quick.sh`                | Sequential `-v<N+i>` probe, fail-fast, ceiling 30 | Replace probe with 20-parallel max-hit-wins. Add `--source` flag.       |
| `install-quick.ps1`               | Same as above (PowerShell mirror)    | Mirror change. Use `Start-Job` or `ForEach-Object -Parallel` (PS 7+) with PS 5.1 fallback to `Invoke-WebRequest -Method Head` in a runspace pool. |
| `gitmap/scripts/install.sh`       | Honors `--version` (skips probe + latest) — already strict-mode-compliant. Discovery still sequential. | Strengthen strict-mode error message to canonical wording (§3.7). Replace probe with 20-parallel. Add Phase C `main` HEAD fallback. |
| `gitmap/scripts/install.ps1`      | Same as above                        | Same as above; PowerShell mirror.                                               |

## Auxiliary files

| File                                                  | Change                                                                                |
|-------------------------------------------------------|---------------------------------------------------------------------------------------|
| `gitmap/release/installsnippet.go` + `constants_release.go` | No code change required — snippet already passes `--version`, which triggers §3 correctly. Add a comment pointing at spec/07-generic-release/09. |
| `spec/01-app/95-installer-script-find-latest-repo.md` | Add a banner at the top: "§4 fail-fast clause superseded by spec/07-generic-release/09 §4.1 (20-parallel max-hit-wins). This document remains valid for context." |
| `spec/07-generic-release/08-pinned-version-install-snippet.md` | Add a "See also" pointer to spec 09. No semantic change. |
| `gitmap/cmd/` install-related help text                | Audit for stale flag names; add `--discovery-window`, `--source` to help if user-facing. |

## Phased rollout

1. **Phase 0 — DONE (this turn).** Spec + memory + plan written. Zero code touched.
2. **Phase 1 — Audit (no writes).** On approval, run a read-only audit
   across the four scripts producing a diff-preview report (no edits).
3. **Phase 2 — Strict-mode hardening.** Tighten error messages and add
   the canonical exit-1 wording. Smallest blast radius. Add tests:
   `--version v0.0.0-nope` MUST exit 1 with no network calls beyond the
   single 404'd download attempt.
4. **Phase 3 — Parallel discovery.** Replace sequential probe with
   parallel HEAD fan-out (window=20). Add `--discovery-window` flag.
   Add tests: probe-count assertion, gap tolerance, loop-guard.
5. **Phase 4 — Phase C main fallback.** Add main-HEAD installation path
   for repos with zero releases. Add `[warn]` banner and SHA recording.
6. **Phase 5 — `--source` flag + docs.** Wire `--source main|latest`,
   update `helptext/`, update `gitmap` CLI changelog (`src/data/changelog.ts`).

Each phase is a separate user-approved turn. No phase begins until the
prior one is verified.

## Verification per phase

* `bash install-quick.sh --help` and `pwsh install-quick.ps1 -?` print
  the new flags.
* `bash install-quick.sh --version v3.12.0 --dry-run` (if dry-run added)
  shows zero `[discovery]` lines.
* `bash install-quick.sh --version v0.0.0-nope` exits 1 with the
  canonical message and zero probe traffic.
* `bash install-quick.sh --discovery-window 5` issues exactly 5 HEADs
  concurrently (verify with `tcpdump`/`mitmproxy` in a sandbox, or by
  asserting `[discovery]` line count).
* New unit tests in `gitmap/cmd/` (or a fresh `gitmap/install/` package
  if logic is extracted) cover the §3 fail-closed contract.

## Generic application to other repos

Any sibling repo wanting to adopt this contract:

1. Copy `spec/07-generic-release/09-generic-install-script-behavior.md`
   into its own `spec/` tree verbatim. The text is repo-agnostic; only
   the `<owner>`, `<stem>`, `<binary>`, `<installerPath>` placeholders
   need substitution in code, not in the spec.
2. Implement §3 first (strict-tag fail-closed). It is the smallest and
   highest-leverage change.
3. Implement §4 Phase A only if the repo follows the `-v<N>` naming
   pattern. Otherwise skip Phase A entirely (spec permits this).
4. Mirror the §10 acceptance checklist into CI.

## Open questions (none blocking)

* Whether `--source` should also accept `commit:<sha>` for power users.
  Defer until a real request lands.
* Whether to add a `--prerelease` flag to opt into draft/prerelease
  picks in Phase B. Defer.
