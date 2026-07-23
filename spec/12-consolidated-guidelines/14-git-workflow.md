# 14 — Git Workflow

Branching strategies, commit conventions, and merge policies.

## Branch Types

| Prefix | Purpose | Lifetime |
|--------|---------|----------|
| `main` | Production-ready | Permanent |
| `release/v*` | Version stabilization | Until tagged |
| `feature/*` | New functionality | Until merged |
| `fix/*` | Bug fixes | Until merged |
| `refactor/*` | Code restructuring | Until merged |
| `chore/*` | Tooling, CI, deps | Until merged |

Lowercase with hyphens. Under 50 characters. Delete after merge.

## Commit Format

```
<type>: <subject>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`. Subject: imperative mood, no period, ≤72 chars.

Each commit compiles and passes tests independently. One logical change per commit.

## Merge Policy

Default: rebase and merge for linear history. Merge commits only for release branches. Squash for noisy WIP commits.

## PR Standards

| Metric | Target | Hard Limit |
|--------|--------|------------|
| Changed lines | ≤200 | ≤400 |
| Files changed | ≤5 | ≤10 |

Rebase before review. Self-review the diff. At least one approving review.

## Tags

SemVer: `vMAJOR.MINOR.PATCH`. Always annotated. Pre-release: `v2.37.0-rc.1`.

---

Source: `spec/05-coding-guidelines/14-git-workflow.md`
