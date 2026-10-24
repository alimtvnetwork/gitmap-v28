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
- `gitmap cbf "description"` (first-class bug fix alias)

## Examples

```bash
gitmap cbf "[db]: fix race condition in version flush"
gitmap cpb "[sqlite]: resolve driver mismatch"
gitmap commit-push-bug "[scanner]: handle nil pointer when DB is empty"
```

## See Also

- `gitmap hd cbf` — High-speed bug fix commit and push alias.
- `gitmap hd commit-format` — Modular bracket syntax and conventional commit formatting guide.
