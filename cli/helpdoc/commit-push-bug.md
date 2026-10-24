# commit-push-bug

Stages all changes, commits with an uppercase `FIX(MODULE): ` prefix, and pushes. Use this for bug fix commits so the git history is clearly categorized.

## Synopsis

```bash
gitmap commit-push-bug "(module): description" [flags]
gitmap cbf "(module): description" [flags]
gitmap cpb "(module): description" [flags]
```

## Aliases

- `gitmap cbf "(module): description"` (first-class bug fix alias)
- `gitmap cpb "(module): description"`

## Examples

```bash
gitmap cbf "(db): fix race condition in version flush"
gitmap cpb "(sqlite): resolve driver mismatch"
gitmap commit-push-bug "(scanner): handle nil pointer when DB is empty"
```

## See Also

- `gitmap hd cbf` — High-speed bug fix commit and push alias.
- `gitmap hd commit-format` — Uppercase Conventional Commit formatting guide.
