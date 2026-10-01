# Command Specification: Semantic Flat Commit Suite & Auto-Stage Command

Spec Reference: [02-spec/21-app/200-semantic-flat-commit-and-auto-stage-command.md](../../../02-spec/21-app/200-semantic-flat-commit-and-auto-stage-command.md)
Plan Reference: [.ai-memory/plans/completed/63-semantic-flat-commit-suite.md](../../plans/completed/63-semantic-flat-commit-suite.md)

## 1. Overview

This specification establishes the commands, flags, and operational contracts for the `gitmap commit` semantic flat commit suite:
1. One-step automated staging (`git add -A`) combined with flat git commit creation.
2. Direct support for un-prefixed commit messages without mandatory `feat:` or `fix:` semantic wrappers.
3. Multi-word positional argument parsing allowing natural terminal inputs without mandatory quote wrapping.
4. Optional remote push integration via `--push` / `-p`.
5. Safe preview mode via `--dry-run` / `-n` inspecting changes without disk mutations.
6. Consolidated CLI root dispatching in `cli/cmd/rootcore.go` eliminating previous routing collisions.

---

## 2. Command Index & Usage

### 2.1 Commands & Aliases
- `gitmap commit "<message>"`: Stages all files and commits with the specified message.
- `gitmap cm "<message>"`: Ultra-fast short alias for daily developer commits.
- `gitmap commit-all "<message>"`: Explicit semantic alias for full workspace auto-staging.
- `gitmap ca "<message>"`: Short alias for commit-all.

### 2.2 Flags & Options
- `-m "<message>"`: Explicit commit message flag (optional; unflagged positional arguments are automatically joined with spaces).
- `-p`, `--push`: Pushes to remote repository branch immediately after a successful commit.
- `-n`, `--dry-run`: Previews changes to be committed via `git status --short` without modifying the repository index.
- `-h`, `--help`: Renders the formatted terminal help menu and usage tips.

---

## 3. CLI Examples

```bash
# Standard 1-step commit with quoted message:
gitmap commit "refactor: simplify split-db path resolution"
gitmap cm "fix: handle empty slice in scanner"

# Positional multi-word commit without quotes:
gitmap cm update readme with new install steps

# Auto-stage, commit, and immediately push:
gitmap commit "feat: add user telemetry" --push
gitmap cm "feat: add user telemetry" -p

# Preview staged status without committing:
gitmap cm "wip" --dry-run
gitmap commit -n
```
