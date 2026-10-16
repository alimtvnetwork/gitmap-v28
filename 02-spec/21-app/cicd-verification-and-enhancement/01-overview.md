# 01 — Overview: cicd-verification-and-enhancement

**Version:** 1.0.0
**Updated:** 2026-10-17

---

## 1. Owner order (verbatim, program 274 charter)

> "CI/CD verification + enhancement program: confirm CI green, fix failures,
> minor bump + release until green; 10 new command suggestions; review last
> 100 commits; 20 codebase improvements excluding rejected ideas; E2E testing
> of recent features; critical feedback."

(Quoted from the program brief handed to Wave 1 — this is the charter text of
record for program 274.)

## 2. The seven deliverables, enumerated

| # | Deliverable | Wave / WS | Output home |
|---|-------------|-----------|-------------|
| 1 | Confirm CI green; fix every failure | WS2 | Fixes committed; `gitmap pe` 100% green |
| 2 | Minor bump + release, repeat until green | WS5 | Release published; CI green on the release commit |
| 3 | 10 new command suggestions | WS3 | `.ai-memory/` proposals (proposal-only) |
| 4 | Review last 100 commits for feature ideas | WS3 | Method in `03-commit-review-method.md` |
| 5 | 20 codebase improvements, excluding rejected ideas | WS3 | `.ai-memory/` proposals (proposal-only) |
| 6 | E2E testing of recent features | WS4 | Plan in `05-e2e-test-plan.md` |
| 7 | Critical feedback on the E2E results | later wave | `06-e2e-critical-feedback.md` (reserved — NOT created in Wave 1) |

## 3. Program workstreams

| WS | Owner | Work | Blocking |
|----|-------|------|----------|
| WS1 | This spec | Spec + plan + subtasks (BEFORE any implementation) | Nothing — Wave 1 |
| WS2 | CI worker | CI green loop: `gitmap pe` → diagnose via gh logs → fix → commit+push → re-check, until 100% green | WS5 blocked on this |
| WS3 | Review worker | Last-100-commit review; 10 command suggestions; 20 improvement proposals (proposals ONLY, no code) | Hard-exclusion list in `04-improvement-criteria.md` |
| WS4 | E2E worker | CLI-level E2E of recent features; dry-run/`--help` only, no live side effects | Writes plan evidence; `06-e2e-critical-feedback.md` left for a later wave |
| WS5 | Lead | Minor bump + release; repeat until CI is green on the release commit | BLOCKED until WS2 reports 100% green |

## 4. Hard rules (non-negotiable, all waves)

- `gitmap aum search` / `gitmap find` / `gitmap cat` / `gitmap ps` / `gitmap py` /
  `gitmap llm train` ONLY. TOTAL BAN on rg, ripgrep, grep, git grep,
  Select-String.
- NEVER run raw `git` commands (no git status/add/commit/push/log/diff).
  Commits go through `gitmap cpf "<module> - <summary>"` (hyphen, NO colons);
  the lead coordinates commit sequencing across workers.
- RELATIVE paths only in every artifact. NEVER absolute paths.
- NEVER `go build`, `go test`, or any test suite. No builds.
- NEVER commit binaries. NEVER delete files. ALIASES ARE KEPT (no pruning).
- All filenames lowercase, kebab-case, two-digit numeric prefixes in specs.
- Files ≤300 lines per the repo mindset.
- Mindset (`.ai-memory/what-to-read.md`): staging-repo velocity — fix forward,
  versions are checkpoints, DRY single sources of truth, enforcement beats
  documentation.
- NEVER cut a new bump/release while a Release run is in flight.
- `version.json` is the single source of truth for versions; the owner's
  "minor bump" means a small increment via the standard bump flow, not a major.

## 5. Definition of done

The program is DONE only when ALL of these hold:

1. `gitmap pe` (and the CI dashboard) show **100% green** — every workflow,
   every platform, including the coverage-floor gate added in program 272
   (fail on >2% coverage drop vs baseline).
2. A release is **published** from a **minor version bump** (from 6.557.0 per
   `version.json` at program start), and **CI is green on the release commit**
   itself. If the release commit is red, bump and re-release until green.
3. 10 command suggestions and 20 improvement proposals are recorded in
   `.ai-memory/` (proposal-only), each passing the criteria in
   `04-improvement-criteria.md` and none touching the hard-exclusion list.
4. E2E of the recent features ran per `05-e2e-test-plan.md` with evidence
   recorded; critical feedback handed to the later wave that owns
   `06-e2e-critical-feedback.md`.
5. Verified personally by the lead — never trust a worker's self-reported
   "green" (standing rule: always verify build/release status independently).

## 6. Ground truth at program start (2026-10-17)

- `version.json` `Version`: **6.557.0**; `LastCommitSha`:
  `fd487dd0b5eab289a3847e9f58428add7907bdb5`.
- Program 272 (v6.553.0) added: coverage floor gate in CI, build-time
  `go:generate` `nc` catalog (105 entries), deprecation policy removed.
- Program 273 (no release): 4 new themes (navy-purple default, tokyo-night,
  vscode-dark-plus, monokai-pro), `gitmap help example <cmd>` 3-part format
  for 527 verbs/aliases, new `mr` / `merge write` command.
- Standing open items (do NOT fold into this program without owner approval):
  Send/Queue buttons broken on owner's machine; PR #6 'Sync v5' do-not-merge;
  release-pipeline durable fix proposed 2026-10-10/11, unanswered.

## 7. File map of this spec

| File | Covers |
|------|--------|
| `01-overview.md` | This file — owner order, workstreams, hard rules, definition of done |
| `02-ci-green-loop-procedure.md` | The exact CI-green loop for WS2 + log-reading procedure |
| `03-commit-review-method.md` | How WS3 reviews the last 100 commits for feature ideas |
| `04-improvement-criteria.md` | What makes a good proposal + the owner-rejected hard-exclusion list |
| `05-e2e-test-plan.md` | CLI-level E2E plan for recent features (WS4) |
