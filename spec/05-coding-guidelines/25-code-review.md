# Code Review

## Overview

Standards for pull request reviews, approval workflows, and automated
quality gates to maintain code quality and knowledge sharing across teams.

---

## Review Checklists

### Author Checklist (Before Requesting Review)

- [ ] Code compiles and all tests pass locally.
- [ ] Self-reviewed the diff — no debug code, commented-out blocks, or TODOs without tickets.
- [ ] Commit messages follow the type-subject convention (`feat:`, `fix:`, etc.).
- [ ] New or changed behavior has corresponding tests.
- [ ] Documentation updated (inline comments, README, spec files) where applicable.
- [ ] No unrelated changes bundled into the PR.

### Reviewer Checklist

- [ ] **Correctness**: Does the logic handle edge cases and error paths?
- [ ] **Naming**: Are variables, functions, and types named clearly per naming conventions?
- [ ] **Size**: Are functions ≤15 lines and files ≤200 lines?
- [ ] **Security**: No secrets in code, no SQL injection vectors, proper input validation.
- [ ] **Performance**: No unnecessary allocations, N+1 queries, or unbounded loops.
- [ ] **Tests**: Are tests meaningful — not just coverage padding?
- [ ] **Style**: Follows project code style (positive conditionals, blank line before return, no magic strings).

---

## PR Standards

### Size Limits

| Metric           | Target    | Hard Limit |
|------------------|-----------|------------|
| Changed lines    | ≤200      | ≤400       |
| Files changed    | ≤5        | ≤10        |
| Commits          | ≤3        | ≤5         |

- PRs exceeding hard limits must be split before review.
- Exception: generated code, migrations, or vendor updates with justification.

### Description Template

```markdown
## What

One-sentence summary of the change.

## Why

Link to spec, issue, or business rationale.

## How to Test

1. Step-by-step manual verification.
2. Or: `go test ./...` / `npm test`.

## Screenshots (if UI)

Before/after screenshots or screen recordings.
```

### Branch Rules

- Branch name follows convention: `feature/*`, `fix/*`, `refactor/*`, `chore/*`.
- Rebased onto current `main` before requesting review.
- No merge commits in the PR branch.

---

## Approval Workflows

### Standard Flow

1. Author opens PR with completed checklist.
2. Assign at least one reviewer with domain knowledge.
3. Reviewer approves or requests changes within **one business day**.
4. Author addresses all comments — no deferred items.
5. Final approval required before merge.
6. Author merges via rebase-and-merge (default).

### Critical Path Flow

For changes affecting authentication, payments, data models, or infrastructure:

1. Minimum **two approving reviews** required.
2. One reviewer must be a domain owner or tech lead.
3. Security-sensitive changes require security review tag.
4. Database migrations require DBA or data team review.

### Review Etiquette

- **Be specific**: "Rename `d` to `duration` for clarity" over "naming is unclear."
- **Distinguish severity**: prefix with `nit:`, `suggestion:`, or `blocker:`.
- **Explain why**: link to a guideline or explain the risk.
- **Acknowledge good work**: positive feedback improves team morale.
- **No bike-shedding**: if it passes the style guide, move on.

---

## Automated Checks

### Required CI Gates

All PRs must pass these checks before merge is enabled:

| Check              | Tool / Method                  | Blocks Merge |
|--------------------|--------------------------------|--------------|
| Lint               | `golangci-lint` / `eslint`     | Yes          |
| Unit tests         | `go test` / `vitest`           | Yes          |
| Build              | `go build` / `vite build`      | Yes          |
| Type check         | `tsc --noEmit`                 | Yes          |
| Coverage threshold | ≥80% lines on changed files    | Yes          |
| Commit format      | commitlint or equivalent       | Yes          |

### Recommended Checks

| Check              | Tool / Method                  | Blocks Merge |
|--------------------|--------------------------------|--------------|
| Security scan      | `npm audit` / `govulncheck`    | Advisory     |
| License compliance | license-checker                | Advisory     |
| PR size            | custom bot / Danger            | Advisory     |
| Stale dependencies | Dependabot / Renovate          | No           |

### Branch Protection Rules

```yaml
# GitHub branch protection
main:
  required_reviews: 1
  dismiss_stale_reviews: true
  require_up_to_date: true
  required_checks:
    - lint
    - test
    - build
    - typecheck
  enforce_admins: true
  restrict_pushes: true
```

### Go Implementation

```go
// reviewcheck validates PR metadata before merge.
func validatePR(pr PullRequest) error {
    if pr.ChangedLines > MaxChangedLines {
        return fmt.Errorf("PR exceeds %d changed lines (%d)", MaxChangedLines, pr.ChangedLines)
    }
    if len(pr.Description) == 0 {
        return fmt.Errorf("PR description is empty")
    }
    if pr.ApprovalCount < MinApprovals {
        return fmt.Errorf("insufficient approvals: %d < %d", pr.ApprovalCount, MinApprovals)
    }
    return nil
}
```

### TypeScript Implementation

```typescript
interface PullRequest {
  changedLines: number;
  filesChanged: number;
  description: string;
  approvalCount: number;
  checksPass: boolean;
}

const MaxChangedLines = 400;
const MaxFilesChanged = 10;
const MinApprovals = 1;

function validatePullRequest(pr: PullRequest): string[] {
  const errors: string[] = [];

  if (pr.changedLines > MaxChangedLines) {
    errors.push(`PR exceeds ${MaxChangedLines} changed lines (${pr.changedLines})`);
  }
  if (pr.filesChanged > MaxFilesChanged) {
    errors.push(`PR exceeds ${MaxFilesChanged} files changed (${pr.filesChanged})`);
  }
  if (pr.description.trim().length === 0) {
    errors.push("PR description is empty");
  }
  if (!pr.checksPass) {
    errors.push("CI checks have not passed");
  }

  return errors;
}
```

---

## Constraints

- All PRs require at least one approving review.
- No self-approvals on production-bound code.
- All CI checks must pass before merge is enabled.
- PRs open longer than 5 business days must be re-rebased or closed.
- Review comments must be resolved — not deferred — before merge.
- Critical path changes require two approvals.
- Automated checks are never bypassed, even by admins.

---

## References

- [Git Workflow](./14-git-workflow.md)
- [CI/CD Patterns](./17-cicd-patterns.md)
- [Documentation Standards](./12-documentation-standards.md)
- [Code Quality Improvement](./01-code-quality-improvement.md)
