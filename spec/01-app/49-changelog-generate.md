# Changelog Generate

## Overview

The `changelog-generate` command auto-generates changelog entries by reading commit messages between two Git tags. It outputs Markdown-formatted changelog sections that can be previewed or written directly to `CHANGELOG.md`.

## Command

```
gitmap-v27 changelog-generate [--from <tag>] [--to <tag>] [--write]
gitmap-v27 cg [--from <tag>] [--to <tag>] [--write]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--from` | second-latest tag | Start tag (older boundary) |
| `--to` | latest tag | End tag or HEAD |
| `--write` | `false` | Prepend output to CHANGELOG.md instead of printing |

## Tag Resolution

When no flags are provided, the command automatically resolves the range:

1. List all `v*` tags sorted by version descending
2. Use the second tag as `--from` and the first tag as `--to`
3. If only one tag exists, use it as `--from` with HEAD as `--to`

When `--from` is specified without `--to`, the range extends to HEAD and the section is labeled "Unreleased".

## Commit Extraction

Uses `git log --format=%s --no-merges <from>..<to>` to collect commit subjects:

- Merge commits are excluded via `--no-merges`
- Only the subject line (first line) of each commit is used
- Empty lines are filtered out

## Output Format

The generated section follows standard CHANGELOG.md format:

```markdown
## v2.24.0

- Add TUI log viewer with detail panel
- Add release rollback on push failure
- Fix watch interval validation edge case
```

## Write Mode

With `--write`, the generated section is prepended to `CHANGELOG.md`:

1. Read existing file content (or start fresh if missing)
2. Prepend the new section with a blank line separator
3. Write back atomically

## Implementation

| File | Purpose |
|------|---------|
| `release/changeloggen.go` | Core logic: commit extraction, formatting, tag listing |
| `cmd/changeloggen.go` | CLI handler: flag parsing, dispatch, file writing |
| `constants/constants_changelog.go` | Constants for command, flags, messages, errors |
| `helptext/changelog-generate.md` | Embedded help text |

## Constraints

- Functions ≤ 15 lines
- Files ≤ 200 lines
- All strings in constants package
- Tags validated before use
- No-op on empty commit range (prints message, exits 0)
