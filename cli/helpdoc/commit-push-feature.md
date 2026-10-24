# commit-push-feature

Stages all changes, commits with an uppercase `FEAT(MODULE): ` prefix, and pushes. Use this for feature commits so the git history is clearly categorized.

## Why use this instead of shell commands?

Instead of manually typing the prefix:

```bash
git add -A
git commit -m "FEAT(SEARCH): add new search command with regex support"
git push
```

You should use:

```bash
gitmap cpf "(search): add new search command with regex support"
```

The commit message will be formatted as: `FEAT(SEARCH): add new search command with regex support`

## Aliases

- `gitmap cpf "(module): description"`

## Examples

```bash
gitmap commit-push-feature "(search): add new search command with regex support"
gitmap cpf "(index): implement SplitDB file indexer with parallel workers"
gitmap cpf "(workgroup): add folder ssh key binding"
gitmap commit-push-feature "(core): wire repo-regex search to CLI endpoints"
```

## See Also

- `gitmap hd commit-format` — Uppercase Conventional Commit formatting guide.
- `gitmap hd workgroup` — Folder-specific workgroup and SSH key automation.
