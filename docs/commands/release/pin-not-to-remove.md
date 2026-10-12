# `gitmap release pin-not-to-remove` (`pnr`)

Manage permanent milestone release tag pins to prevent accidental deletion during range cleanup and release rollbacks.

---

## Overview

The `gitmap release pin-not-to-remove` (`pnr`) command manages a persistent JSON registry at `.gitmap/release/pinned-tags.json`.

Tags stored in this registry are unconditionally protected from:
- `gitmap release remove-between` (`rrb`)
- Batch pruning automation
- Automated repository hygiene routines

---

## Syntax & Aliases

```bash
# Canonical long form
gitmap release pin-not-to-remove [tags...] [flags]

# High-speed alias
gitmap pnr [tags...] [flags]
```

---

## Flags

| Flag | Description |
|---|---|
| `--list`, `-l` | List all currently pinned release tags |
| `--remove <tag>`, `--unpin <tag>` | Remove a tag from the pin registry |
| `--clear` | Remove all pins from the registry |

---

## Examples

```bash
# Pin milestone releases so they are never removed
gitmap pnr 4.162, 4.190

# Pin using SemVer notation
gitmap pnr v4.162.0 v4.190.0

# View all active pins
gitmap pnr --list

# Unpin a tag
gitmap pnr --remove 4.162

# Clear all pins
gitmap pnr --clear
```
