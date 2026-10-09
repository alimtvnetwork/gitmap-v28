# commit-push-chore

Stages all changes, commits with a "Chore: " prefix, and pushes. Use this for maintenance, dependency updates, and housekeeping commits so the git history is clearly categorized.

## Why use this instead of shell commands?

Instead of manually typing the prefix:

```bash
git add -A
git commit -m "Chore: update dependencies and clean caches"
git push
```

You should use:

```bash
gitmap commit-push-chore "update dependencies and clean caches"
```

The commit message will be: `Chore: update dependencies and clean caches`

## Aliases

- `gitmap cpc "description"`

## Semantic Commit Full Forms

GitMap provides standardized 1-step stage, commit, and push shortcuts:

- `cpf <msg>`: `commit-push-feature` (Feature: <msg>)
- `cpb <msg>`: `commit-push-bug` (Bug: <msg>)
- `cpc <msg>`: `commit-push-chore` (Chore: <msg>)
- `cpr <msg>`: `commit-push-release` (Release: <msg>)
- `pcp <msg>`: `pull-commit-push` (pull rebase, stage, commit, and push)
- `pas`: `pull-all-ssh` (pull all repositories using SSH transport)

## Examples

```bash
gitmap commit-push-chore "cleanup unused styles and temporary assets"
gitmap cpc "bump dependencies to latest stable release"
gitmap cpc "cli - update help text and documentation"
```
