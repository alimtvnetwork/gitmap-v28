# Honest Feedback — GitMap Codebase Review

Date: 2026-10-06. Structural survey, not a line-by-line audit.

## Verdict

The codebase's biggest problem is not code quality — it is unchecked growth.
Roughly 3,900 Go files, 100+ CLI packages, 365 help files, ~70 docs-site pages,
300 app specs, and a version at 6.496.0 point to a project that adds far more
than it removes. Most issues below are downstream of that.

## Strengths (with evidence)

- **Discoverable by design.** `.ai-memory/` (overview, what-to-read, plans,
  learned lessons, RCAs) plus `02-spec/` and `docs/commands/` (15 categories)
  make the repo unusually navigable for its size, especially for AI agents.
- **Real engineering infrastructure.** `Makefile` gates (lint, vet, test,
  changelog-check, golden-verify), golden-fixture discipline (`gitmap regoldens`
  as the only sanctioned regeneration path), changelog-from-Conventional-Commits,
  CGo-free SQLite, and cross-platform installers (`install.ps1`/`.sh`,
  `run.ps1`/`.sh`) show sustained investment in reproducibility.
- **Clear core value.** Scanner → manifest (`gitmap.json`) → parallel re-clone
  preserving hierarchy is a coherent, genuinely useful loop for managing a
  large multi-repo tree like `D:\work`.
- **Structured error philosophy.** `apperror` + `cliexit` packages exist instead
  of ad-hoc exits — the right instinct; worth auditing for full adoption.

## Concerns (with evidence)

1. **Duplication is the highest-cost problem.**
   - `gitmap/readme.md` and `gitmap/what-to-read.md` are byte-identical
     (verified: same SHA-256, 2,522 lines each). Every edit must happen twice
     or they drift.
   - Five overlapping clone packages: `cli/cloner`, `cli/clonefrom`,
     `cli/clonenext`, `cli/clonenow`, `cli/clonepick`.
   - Three overlapping state systems: `.gitmap/` files, `gitmap.json`
     manifests, SQLite databases (plus split-DBs) — with no documented
     source of truth per datum.
   - 100+ `cli/constants_*.go` files; likely overlap, never audited as one set.

2. **Docs are impressive but no longer trustworthy.**
   - `.ai-memory/overview.md` still says v3.1.0 / 60+ commands; actual is v6.496.0.
   - `version.json` Title/RepoSlug say "Coding Guidelines" /
     `coding-guidelines-v24` — the wrong repo identity in the file that claims
     to be the single source of truth.
   - Multi-MB `test-inventory.json` / `test-heatmap.json` and 700KB+ audit
     CSVs are committed at root. One accurate page beats ten stale ones.

3. **Style rules likely cost more than they save.** Banning negation and
   `switch`, plus 8–15-line functions across ~3,900 files, almost certainly
   fragments simple logic and fights Go idioms. Consistency is good, but these
   particular rules optimize for machine-checkability over human readability.
   Recommendation: sample real code and relax the rules that produce worse code.

4. **Version numbers have lost meaning.** Bump-on-every-change reached 6.496.0,
   where no user can tell a breaking change from a typo fix. Either make SemVer
   meaningful again or switch to date-based versions.

5. **Committed artifacts are a hygiene smell.** `cli/cli.exe`, `.syso` files,
   `fix_*.py` scripts, `update_runner*.py`, and a `repo-secrets/` directory at
   root suggest the repo doubles as a working folder. Build artifacts belong in
   CI releases, scratch scripts in `/tmp` or a marked archive, secrets nowhere
   near git.

6. **Long install/run scripts are a risk surface.** `install.ps1` (~75KB) and
   `run.ps1` (~69KB) are too large to review as units; they fetch from the
   network and touch system state. Split into sourced modules and audit for
   unnecessary elevation and unpinned URLs.

## What to do first (in order)

1. Consolidate the readme pair — one canonical file.
2. Fix `version.json` repo identity.
3. Remove committed build artifacts and scratch scripts; audit `repo-secrets/`.
4. Merge the clone packages behind one engine.
5. Shrink the readme to an index; delete stale docs instead of syncing them.
6. Revisit the style rules with evidence from real code samples.

Detailed items with priorities live in [action-checklist.md](./action-checklist.md).
