# Git Workflow

## Overview

Standards for branching strategies, commit message conventions, and
merge/rebase policies to keep history clean and collaboration predictable.

---

## Branching Strategy

### Branch Types

| Prefix        | Purpose                          | Lifetime     |
|---------------|----------------------------------|--------------|
| `main`        | Production-ready code            | Permanent    |
| `release/v*`  | Version stabilization            | Until tagged |
| `feature/*`   | New functionality                | Until merged |
| `fix/*`       | Bug fixes                        | Until merged |
| `refactor/*`  | Code restructuring (no behavior) | Until merged |
| `chore/*`     | Tooling, CI, dependencies        | Until merged |

### Naming Rules

- Use lowercase with hyphens: `feature/add-retry-logic`.
- Include a ticket or spec ID when available: `fix/54-lock-stale-pid`.
- Keep names under 50 characters.
- Never use personal names or dates in branch names.

### Branch Lifetime

- Feature branches live at most **one sprint** (two weeks).
- Long-lived branches require weekly rebases onto `main`.
- Delete remote branches immediately after merge.

---

## Commit Message Conventions

### Format

```
<type>: <subject>

<optional body>
```

- **Subject line**: imperative mood, no period, ≤72 characters.
- **Body**: wrapped at 72 characters, explains *why* not *what*.

### Types

| Type       | Usage                                    |
|------------|------------------------------------------|
| `feat`     | New feature or capability                |
| `fix`      | Bug fix                                  |
| `refactor` | Code change with no behavior difference  |
| `test`     | Adding or updating tests                 |
| `docs`     | Documentation only                       |
| `chore`    | Build, CI, tooling, dependencies         |
| `perf`     | Performance improvement                  |

### Examples

```
feat: add exponential backoff to asset upload

Retries up to 3 times with 1s/2s/4s delays on 5xx or
network errors. HTTP 4xx (except 429) fails immediately.
```

```
fix: resolve stale lock file detection on Windows

os.FindProcess always succeeds on Windows; switch to
tasklist-based PID verification.
```

### Commit Granularity

- Each commit compiles and passes tests independently.
- One logical change per commit — do not mix features and fixes.
- Refactors are separate commits from behavior changes.

---

## Merge and Rebase Policies

### Default Strategy: Rebase and Merge

Preferred for feature and fix branches to maintain a linear history:

1. Rebase the branch onto `main` locally.
2. Force-push the rebased branch.
3. Merge via fast-forward on the remote.

### When to Use Merge Commits

Use a merge commit (no fast-forward) only for:

- Release branches merging into `main`.
- Long-running branches with shared collaboration.
- Preserving a meaningful integration point in history.

### Squash Policy

- Squash when a feature branch has noisy WIP commits.
- Do not squash when individual commits carry distinct meaning.
- Never squash across multiple logical changes.

### Conflict Resolution

- The branch author resolves conflicts, not the reviewer.
- Rebase onto the latest `main` before requesting review.
- Never resolve conflicts in the merge commit itself.

---

## Pull Request Standards

### Before Opening

- Rebase onto current `main`.
- Ensure all tests pass locally.
- Self-review the diff for leftover debug code.

### PR Description

Include:

- **What** changed (one-sentence summary).
- **Why** it changed (link to spec or issue).
- **How to test** (manual steps or test commands).

### Review Expectations

- At least one approving review before merge.
- Reviewer checks logic, edge cases, and naming.
- Address all comments before merging — do not defer.

---

## Tag and Release Conventions

### Tag Format

- Semantic versioning: `vMAJOR.MINOR.PATCH` (e.g., `v2.36.7`).
- Pre-release tags: `v2.37.0-rc.1`.
- Always annotated tags with a message: `git tag -a v2.37.0 -m "v2.37.0"`.

### Release Branch Flow

1. Create `release/vX.Y.Z` from `main`.
2. Apply final fixes on the release branch.
3. Tag the release commit.
4. Merge back into `main` (merge commit).
5. Delete the release branch.

---

## Constraints

- All branches follow the naming convention.
- All commits follow the type-subject format.
- No direct pushes to `main` — all changes go through PRs.
- Branches deleted from remote within 24 hours of merge.
- Tags are annotated, never lightweight.
- Rebase is the default; merge commits require justification.
