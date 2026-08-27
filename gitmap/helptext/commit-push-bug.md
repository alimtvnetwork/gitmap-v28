# commit-push-bug

Stages all changes, commits with a "Bug: " prefix, and pushes. Use this for bug fix commits so the git history is clearly categorized.

## Why use this instead of shell commands?

Instead of manually typing the prefix:

```bash
git add -A
git commit -m "Bug: fix null pointer in scanner"
git push
```

You should use:

```bash
gitmap commit-push-bug "fix null pointer in scanner"
```

The commit message will be: `Bug: fix null pointer in scanner`

## Aliases

- `gitmap cpb "description"`

## Examples

```bash
gitmap commit-push-bug "fix race condition in version flush"
gitmap cpb "resolve sqlite driver mismatch"
gitmap commit-push-bug "handle nil pointer when DB is empty"
```
