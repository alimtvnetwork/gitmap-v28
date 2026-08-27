# commit-push-pull

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
gitmap commit-push-pull "your message"
```

If the pull fails due to conflicts, GitMap will stop and notify you to resolve them manually before retrying.

## Aliases

- `gitmap cpp "message"`

## Examples

```bash
gitmap commit-push-pull "fix: resolve merge conflict in constants.go"
gitmap cpp "feat: add new search command"
gitmap commit-push-pull "chore: sync with upstream before pushing release"
```
