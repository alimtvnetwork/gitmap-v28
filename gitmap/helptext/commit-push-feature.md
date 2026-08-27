# commit-push-feature

Stages all changes, commits with a "Feature: " prefix, and pushes. Use this for feature commits so the git history is clearly categorized.

## Why use this instead of shell commands?

Instead of manually typing the prefix:

```bash
git add -A
git commit -m "Feature: add new search command with regex support"
git push
```

You should use:

```bash
gitmap commit-push-feature "add new search command with regex support"
```

The commit message will be: `Feature: add new search command with regex support`

## Aliases

- `gitmap cpf "description"`

## Examples

```bash
gitmap commit-push-feature "add commit-push commands for quick git workflows"
gitmap cpf "implement SplitDB file indexer with parallel workers"
gitmap commit-push-feature "wire repo-regex search to CLI endpoints"
```
