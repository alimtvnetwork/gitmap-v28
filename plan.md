# gitmap: Future Work Roadmap

Generated: 2026-07-23 14:06 UTC. Companion to `.lovable/memory/reports/20260723-rejog-reliability.md`.

This roadmap is written for hand-off to another AI. Every task is scoped to a specific project, has explicit dependencies, and lists acceptance criteria. Do not implement anything from here until the user selects a task from the "Next task selection" section.

Projects:
- **gitmap-cli**: Go CLI under `gitmap/`.
- **gitmap-updater**: standalone updater under `gitmap-updater/`.
- **docs-site**: React docs site under `src/`.
- **ci-pipeline**: workflows under `.github/`.
- **spec-hygiene**: everything under `spec/` and `.lovable/`.

Phases:
- **P0 Stabilize**: fix onboarding lies and spec ambiguity before any large work.
- **P1 Feature completion**: finish planned features whose specs already exist.
- **P2 Polish**: hygiene, tests, docs.

---

## P0 Stabilize

### P0-1: Sync current-version references
- **Project**: spec-hygiene
- **Objective**: Make `.lovable/overview.md` and `.lovable/memory/index.md` state the same version as `gitmap/constants/constants.go`.
- **Dependencies**: none.
- **Expected outputs**: edits to `.lovable/overview.md` and `.lovable/memory/index.md`. A one-line note naming `gitmap/constants/constants.go` as the source of truth.
- **Acceptance criteria**:
  - Grep for the old versions (v3.1.0, v5.9.0) returns 0 matches in `.lovable/`.
  - Both files display the version currently in `constants.go`.
  - Files still parse as valid markdown.

### P0-2: Deduplicate numeric prefixes in spec/01-app and spec/02-app-issues
- **Project**: spec-hygiene
- **Objective**: Every spec file has a unique numeric prefix, and a README maps prefixes to topics.
- **Dependencies**: none.
- **Expected outputs**:
  - `spec/01-app/README.md` and `spec/02-app-issues/README.md` mapping prefix to file.
  - Renames for duplicated prefixes: 26, 27, 89, 90, 95, 96, 100, 108, 109, 110, 111 (list in reliability report section 1.4).
  - Grep-based redirect stubs or an updated cross-reference table in any spec that cites a renamed file.
- **Acceptance criteria**:
  - `ls spec/01-app | awk -F- '{print $1}' | sort | uniq -d` returns empty.
  - Same for `spec/02-app-issues`.
  - No broken internal links (`rg '\[.*\]\(.*01-app/.*\)'` all resolve).

### P0-3: Resolve function-length contradiction
- **Project**: spec-hygiene
- **Objective**: One canonical function-length rule.
- **Dependencies**: none.
- **Expected outputs**: edits to `spec/05-coding-guidelines/01-*.md`, `spec/05-coding-guidelines/02-go-code-style.md`, and `spec/12-consolidated-guidelines/02-go-code-style.md`.
- **Acceptance criteria**: All three files agree on the same maximum lines per function, and cross-link each other.

### P0-4: Split Core rules from session flags in mem://index.md
- **Project**: spec-hygiene
- **Objective**: Prevent "NO-QUESTIONS MODE" and its 40-task budget from silently persisting across sessions.
- **Dependencies**: none.
- **Expected outputs**: `mem://index.md` reorganized into `## Core (permanent)` and `## Session flags (expires)` sections; No-Questions Mode moved to Session flags with an explicit expiry.
- **Acceptance criteria**: Permanent Core rules would apply to any project or session; Session flags are clearly marked.

---

## P1 Feature completion

### P1-1: Ship spec 100 `scan-all` and spec 101 `pull-all`
- **Project**: gitmap-cli
- **Objective**: Implement the two planned bulk commands.
- **Dependencies**: P0-2 (prefix map so `100-scan-all` vs `100-clone-pick` is unambiguous).
- **Expected outputs**: new `cmd/scan_all.go`, `cmd/pull_all.go`, constants entries in `constants_cli.go`, tests, feature memory files.
- **Acceptance criteria**: specs 100 (scan-all) and 101 acceptance criteria pass; help text present; parallel worker respects `--max-concurrency`.

### P1-2: Resolve open SSH bug (`.lovable/issues/01-ssh-repo-cloned-as-https.md`)
- **Project**: gitmap-cli
- **Objective**: Fix reclone losing SSH transport (spec `110-clone-ssh-flag` + plan `03-reclone-transport-and-vscode-open`).
- **Dependencies**: P0-1.
- **Expected outputs**: patch to clone/reclone code, regression test that reclones an ssh:// repo and asserts the resulting `.git/config` origin is ssh://.
- **Acceptance criteria**: automated test passes; issue file moved to `.lovable/solved-issues/`.

### P1-3: `commit-in` command
- **Project**: gitmap-cli
- **Objective**: Implement per `spec/03-commit-in/`.
- **Dependencies**: P0-2, user typing `next` (the spec is gated).
- **Expected outputs**: full command family, migrations for `.gitmap/commit-in/profiles/`, ShaMap dedupe.
- **Acceptance criteria**: spec `01-overview-and-glossary` invariants hold; profile schema matches `spec/03-commit-in/05-profiles-and-json-shape.md`.

### P1-4: `update-remote-probe` (spec 111)
- **Project**: gitmap-cli
- **Objective**: Ship remote-probe update flow.
- **Dependencies**: P1-2 not required; P0-2 required.
- **Expected outputs**: cmd + tests, reconciled with spec `110-update-remote-installer`.
- **Acceptance criteria**: spec `111-update-remote-probe` acceptance criteria pass; no overlap with `110`.

### P1-5: Bulk visibility (spec 113 + 116)
- **Project**: gitmap-cli
- **Objective**: `mapub` / `mapri` bulk visibility flip.
- **Dependencies**: P1-2 (uses the transport work).
- **Expected outputs**: cmd, GitHub API integration, dry-run mode.
- **Acceptance criteria**: spec `113-clone-parent-escape-and-bulk-visibility.md` + `116-bulk-visibility-mapub-mapri.md` criteria pass.

---

## P2 Polish

### P2-1: Backfill unit tests for `task`, `env`, `install`
- **Project**: gitmap-cli
- **Objective**: Close gap flagged in `.lovable/pending-issues/`.
- **Dependencies**: none.
- **Expected outputs**: new `_test.go` files covering happy path + error path per command.
- **Acceptance criteria**: coverage floor (`.github/coverage.floor`) raised; CI green.

### P2-2: Compress `.lovable/prompts/03-*` through `20-next-task.md`
- **Project**: spec-hygiene
- **Objective**: One canonical `next-task.md`; older per-invocation copies moved to `.lovable/prompts/archive/`.
- **Dependencies**: none.
- **Expected outputs**: single canonical file + archive folder.
- **Acceptance criteria**: repo policy from user memory ("no per-invocation archive mirrors") satisfied.

### P2-3: Deduplicate `.lovable/spec/commands/` vs `spec/01-app/`
- **Project**: spec-hygiene
- **Objective**: Single source of truth.
- **Dependencies**: P0-2.
- **Expected outputs**: `.lovable/spec/commands/README.md` redirecting to root `spec/01-app/`.

### P2-4: Contract-test index
- **Project**: spec-hygiene
- **Objective**: One page listing every contract test that guards an invariant (hosted-docs fallback, wrapper marker, stablejson, cn folder-arg dispatcher, clone-next flatten).
- **Dependencies**: P0-2.
- **Expected outputs**: `spec/01-app/CONTRACT-TESTS.md`.

### P2-5: Docs-site content parity with shipped features
- **Project**: docs-site
- **Objective**: Ensure pages under `src/pages/` cover every command listed in `spec/01-app/` up to the current version.
- **Dependencies**: P0-1, P0-2.
- **Expected outputs**: audit table + missing pages.

---

## Next task selection

Pick one of these to start implementing. Each is ready to begin with no upstream blockers.

1. **P0-1 Sync current-version references.** Smallest, unblocks every other onboarding. Two file edits.
2. **P0-3 Resolve function-length contradiction.** Also small, removes an active review-time contradiction.
3. **P0-2 Deduplicate numeric prefixes.** Larger but high leverage; enables P1-1, P1-4, P2-3, P2-4.
4. **P1-2 Fix SSH reclone transport bug.** Only tier-2 item ready today; a concrete user-visible fix.
5. **P2-1 Backfill unit tests for `task`, `env`, `install`.** Independent, mechanical, restores CI coverage floor.

Say which number (or name) you want and I will start implementation.
