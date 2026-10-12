# `gitmap release remove-between` (`rrb`)

Prune intermediate Git tags and GitHub releases between two release bounds while anchoring milestone versions.

<div align="center">

```
  [4.162]  ───────────────── ( intermediate tags purged ) ─────────────────  [4.190]
  Anchor (KEPT)           v4.162.1 .. v4.189.0 deleted            Anchor (KEPT)
```

</div>

---

## Overview

When developing in high-frequency CI/CD environments, repositories accumulate many intermediate pre-release, beta, or test tags. When preparing clean releases, `gitmap release remove-between` (`rrb`) safely prunes the intermediate tag clutter while keeping:
1. The **lower bound** anchor milestone (e.g. `v4.162.0`).
2. The **upper bound** anchor milestone (e.g. `v4.190.0`).
3. Any tag registered in the **pin registry** (`.gitmap/release/pinned-tags.json`) or passed via `--pin`.

---

## Syntax & Aliases

```bash
# Canonical long form
gitmap release remove-between <start> <end> [flags]

# High-speed alias
gitmap rrb <start> <end> [flags]
```

### Argument Formats
- Comma-separated: `gitmap rrb 4.162, 4.190`
- Space-separated: `gitmap rrb 4.162 4.190`
- SemVer prefixed: `gitmap rrb v4.162.0 v4.190.0`

---

## Flags

| Flag | Description |
|---|---|
| `--confirm` | Bypass interactive confirmation prompt and proceed immediately |
| `-y, --yes` | Shorthand alias for `--confirm` |
| `--dry-run` | Explicit preview mode: outputs evaluation and exits without prompting or deleting |
| `--pin <list>` / `--keep <list>` | Comma-separated list of additional tags to protect |
| `--keep-remote` | Only delete local tags; preserve remote origin tags |

---

## Operational Pipeline

1. **AST/SemVer Validation:** Validates `<start>` and `<end>`. Inverted inputs are automatically normalized (`min .. max`).
2. **Automatic Dry-Run Evaluation:** Every execution calculates and displays candidate tags queued for removal, milestone anchors preserved, and pinned tags kept.
3. **Interactive Confirmation Guard:** Prompts the user (`[y/N]`). The **only** way to bypass this prompt is via `--confirm` or `-y`.
4. **Anchor & Pin Protection:** Excludes `<start>`, `<end>`, and all pinned tags from deletion.
5. **GitHub Release Deletion:** Deletes corresponding GitHub Release objects via `gh release delete` if configured.
6. **Remote Git Deletion:** Pushes batch deletion requests to `origin`.
7. **Local Git Deletion:** Deletes local tags via `git tag -d`.
8. **Sidecar Cleanup:** Removes `.gitmap/release/<tag>.json` metadata files.
9. **Repository Optimization:** Executes `git fetch --prune --prune-tags origin` and `git gc --prune=now`.

---

## Examples

```bash
# Evaluates and prints auto dry-run preview, then asks for confirmation:
gitmap rrb 4.162, 4.190

# Bypass confirmation prompt directly:
gitmap rrb 4.162, 4.190 --confirm

# Bypass confirmation prompt with shorthand -y:
gitmap rrb 4.162 4.190 -y

# Pure preview mode (no prompt, zero side-effects):
gitmap rrb 4.162, 4.190 --dry-run

# Preserve an additional build version:
gitmap rrb 4.162, 4.190 --pin 4.175.0 --confirm
```
