# Git Workflow

Standards for branching strategies, commit conventions, merge policies, and release tagging across all projects.

---

## 1. Branching Strategies

### Branch Types

| Branch | Purpose | Lifetime | Naming |
|---|---|---|---|
| `main` | Production-ready code | Permanent | `main` |
| Feature | New functionality | Short (days) | `feature/<short-description>` |
| Bugfix | Non-urgent fixes | Short (days) | `bugfix/<short-description>` |
| Hotfix | Urgent production fix | Very short (hours) | `hotfix/<short-description>` |
| Release | Release preparation | Short (days) | `release/vX.Y.Z` |

### Flow

```
main ← feature/add-scanner
main ← bugfix/fix-csv-output
main ← hotfix/patch-auth-crash
main ← release/v1.2.0
```

### Rules

| Rule | Detail |
|---|---|
| `main` is always deployable | Never push broken code to `main` |
| Branch from `main` | All branches originate from the latest `main` |
| One concern per branch | A branch addresses exactly one feature, fix, or task |
| Short-lived branches | Merge within days, not weeks |
| Delete after merge | Remove branches once merged to keep the repo clean |
| No direct commits to `main` | All changes go through pull requests |

### Naming Conventions

```
feature/add-export-command
bugfix/fix-null-pointer-scan
hotfix/patch-release-lock
release/v2.1.0
```

- Lowercase, hyphen-separated.
- Prefix indicates intent.
- Description is concise (2–4 words).

---

## 2. Commit Conventions

### Format

```
<type>: <subject>

<optional body>
```

### Types

| Type | Usage |
|---|---|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Code restructuring without behavior change |
| `docs` | Documentation only |
| `test` | Adding or updating tests |
| `chore` | Build, CI, tooling, or dependency updates |
| `perf` | Performance improvement |
| `style` | Formatting, whitespace (no logic change) |

### Examples

```
feat: add CSV export command

fix: resolve null pointer in scanner when directory is empty

refactor: extract flag parsing into dedicated functions

docs: add shell completion spec

chore: upgrade cobra to v1.9.0

test: add table-driven tests for slug parser
```

### Rules

| Rule | Detail |
|---|---|
| Subject line ≤ 72 characters | Short and scannable |
| Imperative mood | "add feature" not "added feature" or "adds feature" |
| No trailing period | Subject line has no period at the end |
| Lowercase subject | Start with lowercase after the type prefix |
| Body wraps at 80 characters | Optional body explains why, not what |
| One logical change per commit | Atomic commits that can be reverted independently |
| No `WIP` commits on `main` | Squash or rewrite before merging |

### Go Implementation

```go
// commitlint.go
func ValidateCommitMessage(msg string) error {
    subject := strings.SplitN(msg, "\n", 2)[0]
    if len(subject) > 72 {
        return fmt.Errorf("subject exceeds 72 characters: %d", len(subject))
    }

    validTypes := []string{"feat", "fix", "refactor", "docs", "test", "chore", "perf", "style"}
    parts := strings.SplitN(subject, ": ", 2)
    if len(parts) != 2 {
        return fmt.Errorf("missing type prefix: expected '<type>: <subject>'")
    }

    isValidType := false
    for _, t := range validTypes {
        if parts[0] == t {
            isValidType = true
            break
        }
    }
    if !isValidType {
        return fmt.Errorf("invalid commit type: %q", parts[0])
    }

    return nil
}
```

### TypeScript Implementation

```typescript
const validTypes = ["feat", "fix", "refactor", "docs", "test", "chore", "perf", "style"];

function validateCommitMessage(msg: string): string[] {
  const errors: string[] = [];
  const subject = msg.split("\n")[0];

  if (subject.length > 72) {
    errors.push(`Subject exceeds 72 characters: ${subject.length}`);
  }

  const match = subject.match(/^(\w+): .+$/);
  if (!match) {
    errors.push("Missing type prefix: expected '<type>: <subject>'");
    return errors;
  }

  if (!validTypes.includes(match[1])) {
    errors.push(`Invalid commit type: '${match[1]}'`);
  }

  return errors;
}
```

---

## 3. Merge Policies

### Pull Request Requirements

| Requirement | Detail |
|---|---|
| At least one approval | Reviewer must approve before merge |
| CI passes | All pipeline stages (lint, test, build) must succeed |
| No unresolved comments | All review threads must be resolved |
| Up-to-date with `main` | Branch must be rebased or merged with latest `main` |
| PR description filled | What, why, and how to test |

### Merge Strategy

| Strategy | When to Use |
|---|---|
| Squash merge | Default for feature and bugfix branches |
| Merge commit | Release branches merging to `main` |
| Rebase | Small single-commit PRs with clean history |

### Squash Merge Rules

```
# Final squashed commit message follows commit conventions
feat: add CSV export command

# PR title becomes the squash commit subject
# PR body becomes the squash commit body
```

| Rule | Detail |
|---|---|
| PR title follows commit conventions | The squashed commit inherits the PR title |
| Delete source branch after merge | Automatic cleanup |
| No force-push to `main` | History is append-only |
| Rebase before merge if conflicted | Resolve conflicts in the feature branch |

### Critical Path Reviews

Changes to these areas require two approvals including a tech lead:

- Authentication and authorization
- Database schema changes
- CI/CD pipeline configuration
- Security-sensitive code
- Public API contracts

---

## 4. Release Tagging

### Tag Format

```
vMAJOR.MINOR.PATCH
```

Examples: `v1.0.0`, `v2.3.1`, `v0.1.0-beta.1`

### Tagging Process

```
1. Create release branch: release/vX.Y.Z
2. Update version constant and changelog
3. Merge release branch to main
4. Tag the merge commit: vX.Y.Z
5. Push tag to remote
6. CI builds and publishes release artifacts
7. Delete release branch
```

### Go

```go
// version.go
const Version = "1.2.0"

// Tag creation
// git tag -a v1.2.0 -m "Release v1.2.0"
// git push origin v1.2.0
```

### Automation

```bash
#!/bin/bash
# release-tag.sh
set -euo pipefail

VERSION="$1"
TAG="v${VERSION}"

if git rev-parse "$TAG" >/dev/null 2>&1; then
    echo "Tag $TAG already exists"
    exit 1
fi

git tag -a "$TAG" -m "Release $TAG"
git push origin "$TAG"
echo "Tagged and pushed $TAG"
```

### Rules

| Rule | Detail |
|---|---|
| Semantic versioning only | `vMAJOR.MINOR.PATCH` with optional pre-release suffix |
| Tags are immutable | Never delete or move a tag after push |
| Tag on `main` only | Release tags point to `main` branch commits |
| Annotated tags required | Use `git tag -a`, never lightweight tags |
| Version constant matches tag | Compiled version must equal the Git tag |
| Changelog updated before tagging | Every release has a corresponding changelog entry |
| One tag per release | No duplicate or ambiguous version tags |

### Pre-Release Tags

| Format | Usage |
|---|---|
| `v1.0.0-alpha.1` | Early testing, unstable |
| `v1.0.0-beta.1` | Feature-complete, testing phase |
| `v1.0.0-rc.1` | Release candidate, final verification |

---

## Constraints

| Constraint | Detail |
|---|---|
| No direct pushes to `main` | All changes via pull requests |
| Commit messages follow conventions | `<type>: <subject>` format enforced |
| Squash merge by default | Clean, linear history on `main` |
| Branches deleted after merge | No stale branches accumulate |
| Tags are immutable and annotated | Never rewrite release history |
| CI must pass before merge | No exceptions |
| One concern per branch and commit | Atomic, reversible changes |
| Release tags match version constants | Single source of truth |

---

## References

- [CI/CD Patterns](./17-cicd-patterns.md)
- [Code Review](./25-code-review.md)
- [Dependency Management](./27-dependency-management.md)
- [Documentation Standards](./12-documentation-standards.md)

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
