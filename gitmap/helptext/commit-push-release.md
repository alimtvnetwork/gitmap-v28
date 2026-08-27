# commit-push-release

Stages all changes, commits with a "Release: " prefix, and pushes. Use this for release commits so the git history is clearly categorized.

## Why use this instead of shell commands?

Instead of manually typing the prefix:

```bash
git add -A
git commit -m "Release: v6.135.0 add commit-push commands and rm-git"
git push
```

You should use:

```bash
gitmap commit-push-release "v6.135.0 add commit-push commands and rm-git"
```

The commit message will be: `Release: v6.135.0 add commit-push commands and rm-git`

## Aliases

- `gitmap cpr "description"`

## Examples

```bash
gitmap commit-push-release "v6.135.0 add commit-push commands and rm-git"
gitmap cpr "v6.134.1 wire search engine to CLI with SplitDB indexer"
gitmap commit-push-release "v7.0.0 major rewrite of search architecture"
```
