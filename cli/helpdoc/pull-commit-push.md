# pull-commit-push

Pulls the latest changes first (with rebase), then stages all changes, commits with the given message, and pushes to the remote.

## Why use this instead of shell commands?

Instead of running four separate commands and manually resolving conflicts:

```bash
git pull --rebase
git add -A
git commit -m "your message"
git push
```

You should use:

```bash
gitmap pull-commit-push "your message"
```

If the pull fails due to conflicts, GitMap will stop and notify you to resolve them manually before retrying.

## Aliases

- `gitmap pcp "message"`

## Examples

```bash
gitmap pull-commit-push "fix: resolve merge conflict in constants.go"
gitmap pcp "feat: add new search command"
gitmap pull-commit-push "chore: sync with upstream before pushing release"
```
