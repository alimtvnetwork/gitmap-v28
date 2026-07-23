# Contributing to GitMap

Thank you for your interest in contributing to GitMap. This guide covers the
development workflow, coding standards, and pull-request requirements.

---

## Getting Started

```bash
git clone https://github.com/alimtvnetwork/gitmap-v27.git gitmap
```

```bash
cd gitmap
./setup.sh          # installs hooks, linter, and downloads Go deps
```

`setup.sh` performs:

- Installs the **pre-commit hook** (runs `golangci-lint` before each commit).
- Verifies the Go toolchain and installs `golangci-lint` if missing.
- Downloads Go module dependencies.

---

## Development Workflow

### 1. Create a Branch

Branch from the latest `main`. Use the correct prefix:

| Prefix | Purpose |
|---|---|
| `feature/<desc>` | New functionality |
| `bugfix/<desc>` | Non-urgent fix |
| `hotfix/<desc>` | Urgent production fix |
| `refactor/<desc>` | Code restructuring |
| `chore/<desc>` | Build, CI, or tooling |

Names are lowercase, hyphen-separated, 2–4 words (e.g. `feature/add-export-command`).

### 2. Write Code

Follow the project coding standards:

- **Functions** ≤ 15 lines (excluding blanks/comments). Split larger functions.
- **Files** ≤ 200 lines. Split by responsibility.
- **No magic strings** — use constants from the `constants` package.
- **Positive conditionals** — `if ready` not `if !notReady`.
- **Blank line before `return`** (except single-line bodies).
- **Boolean names** start with `is` or `has`.
- See [`spec/05-coding-guidelines/`](spec/05-coding-guidelines/) for the full ruleset.

### 3. Make Targets

Run these locally before pushing:

```bash
make lint       # golangci-lint (5 min timeout)
make test       # go test -v -count=1 ./...
make build      # compile binary with version info
make vulncheck  # govulncheck vulnerability scan
make all        # lint → test → build
```

### 4. Commit Messages

Format: `<type>: <subject>`

| Type | Usage |
|---|---|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Restructuring (no behavior change) |
| `docs` | Documentation only |
| `test` | Adding or updating tests |
| `chore` | Build, CI, tooling, dependency updates |
| `perf` | Performance improvement |
| `style` | Formatting (no logic change) |

Rules:

- Subject ≤ 72 characters, imperative mood, no trailing period.
- One logical change per commit.
- No `WIP` commits — squash before opening a PR.

```
feat: add CSV export command
fix: resolve null pointer in scanner for empty directory
docs: add shell completion spec
```

---

## Pull Request Requirements

### Before Opening

- [ ] Code compiles and all tests pass locally (`make all`).
- [ ] Self-reviewed the diff — no debug code, commented-out blocks, or orphan TODOs.
- [ ] Commit messages follow the `<type>: <subject>` convention.
- [ ] New or changed behavior has corresponding tests.
- [ ] Documentation updated where applicable.
- [ ] No unrelated changes bundled into the PR.

### PR Size Limits

| Metric | Target | Hard Limit |
|---|---|---|
| Changed lines | ≤ 200 | ≤ 400 |
| Files changed | ≤ 5 | ≤ 10 |
| Commits | ≤ 3 | ≤ 5 |

PRs exceeding hard limits must be split before review. Exceptions: generated
code, migrations, or vendor updates (with justification).

### Description Template

```markdown
## What

One-sentence summary of the change.

## Why

Link to spec, issue, or business rationale.

## How to Test

1. Step-by-step manual verification.
2. Or: `make test`.

## Screenshots (if UI)

Before/after screenshots or screen recordings.
```

### Branch Rules

- Rebased onto current `main` before requesting review.
- No merge commits in the PR branch.

---

## Review Process

### Standard Flow

1. Author opens PR with the completed checklist.
2. Assign at least one reviewer with domain knowledge.
3. Reviewer approves or requests changes within **one business day**.
4. Author resolves all comments — nothing deferred.
5. Final approval required before merge.
6. Merge via **squash merge** (default).

### Critical Path Flow

Changes to authentication, database schema, CI/CD pipelines, security, or
public APIs require:

- **Two approving reviews**, including a domain owner or tech lead.
- Security-sensitive changes tagged for security review.
- Database migrations require DBA or data-team review.

### Review Etiquette

- Be specific: _"Rename `d` to `duration` for clarity"_ not _"naming is unclear."_
- Prefix with `nit:`, `suggestion:`, or `blocker:` to distinguish severity.
- Explain _why_ — link to a guideline or explain the risk.
- Acknowledge good work.

---

## CI Checks

All PRs must pass these automated gates before merge:

| Check | Tool | Blocks Merge |
|---|---|---|
| Lint | `golangci-lint` / `eslint` | Yes |
| Unit tests | `go test` / `vitest` | Yes |
| Build | `go build` / `vite build` | Yes |
| Vulnerability scan | `govulncheck` | Advisory |

### When CI runs (and when it doesn't)

CI is wired in [`.github/workflows/ci.yml`](.github/workflows/ci.yml). The triggers are:

| Event | Runs CI? | Notes |
|---|:---:|---|
| Push to **any branch** (incl. `main`, feature branches, Lovable branches) | ✅ | Configured via `branches: ['**']` |
| Push to `release/**` | ❌ | Owned by [`release.yml`](.github/workflows/release.yml) — avoids duplicate builds |
| Push of a `v*` tag | ❌ | Same — release pipeline handles it |
| Pull request → `main` | ✅ | Runs as a `pull_request` event |
| Manual "Run workflow" (workflow_dispatch) | ✅ | Useful for re-running with `lint_baseline_disable=true` |

**You do not need to open a PR to get CI feedback** — every push to a feature branch fires the full pipeline.

### When `sha-check` may skip a run

The first real job, **`sha-check`**, looks up a cache entry keyed by `ci-passed-${{ github.sha }}`. On a hit, every downstream job short-circuits and prints `✅ SHA … already passed`. This looks like "CI didn't run" but is intentional dedup.

A run will be deduped (skipped) when:

1. **The same commit SHA already passed CI** on another branch or in another run (e.g. a feature branch was fast-forward-merged into `main` — the merge commit is identical).
2. **A force-push retargets an existing commit** (the SHA was already validated, so no re-run).
3. **Re-running a previously green workflow** for a SHA that's still in the cache.

A run will **not** be deduped (full pipeline executes) when:

- Any file changed → new SHA → cache miss.
- The cache entry was evicted by GitHub (caches expire after ~7 days of no access, or when the repo exceeds 10 GB of cache).
- You manually invalidate by bumping `lint_baseline_cache_version` via "Run workflow".

### How to confirm what happened

Open the Actions run → check the **Diagnostics** job (always runs, never skipped). It prints the exact branch, SHA, event, and actor. The **Summary** tab also shows a "Fresh build" or "Deduped" verdict at the top, so you can tell in one glance whether downstream jobs did real work or short-circuited.

---

## Release Process

Releases follow semantic versioning (`vMAJOR.MINOR.PATCH`):

```bash
make release BUMP=patch   # default: patch
make release BUMP=minor
make release BUMP=major
make release-dry          # preview without executing
```

See the [release spec](spec/01-app/12-release-command.md) for details.

---

## Specs and Architecture

For significant features or architectural changes, create or update a
specification in [`spec/`](spec/) for review **before** implementation. See
the [spec README](spec/README.md) for structure and naming conventions.

---

## References

- [Code Quality Guidelines](spec/05-coding-guidelines/01-code-quality-improvement.md)
- [Go Code Style](spec/05-coding-guidelines/02-go-code-style.md)
- [Naming Conventions](spec/05-coding-guidelines/03-naming-conventions.md)
- [Git Workflow](spec/05-coding-guidelines/28-git-workflow.md)
- [Code Review Standards](spec/05-coding-guidelines/25-code-review.md)
- [CI/CD Patterns](spec/05-coding-guidelines/17-cicd-patterns.md)
