# commit-push

Stages all changes, commits with the given message, and pushes to the remote in one command.

## Why use this instead of shell commands?

Instead of running three separate commands:

```bash
git add -A
git commit -m "your message"
git push
```

You should use:

```bash
gitmap commit-push "your message"
```

## Aliases

- `gitmap cp "message"`

## Semantic Commit Shortcuts & Full Forms

GitMap provides dedicated 1-step stage, commit, and push shortcuts for semantic git history:

- `cpf <msg>`: `commit-push-feature` (Feature: <msg>)
- `cpb <msg>`: `commit-push-bug` (Bug: <msg>)
- `cpc <msg>`: `commit-push-chore` (Chore: <msg>)
- `cpr <msg>`: `commit-push-release` (Release: <msg>)
- `pcp <msg>`: `pull-commit-push` (pull rebase, stage, commit, and push)
- `pas`: `pull-all-ssh` (pull all repositories using SSH transport)

## Examples

```bash
gitmap commit-push "fix: resolve null pointer in scanner"
gitmap cp "docs: update readme with new install steps"
gitmap commit-push "refactor: extract shared DB resolver to cmd_db.go"
```
