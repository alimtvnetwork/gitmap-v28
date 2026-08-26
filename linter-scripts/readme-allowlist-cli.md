# `allowlist-forbidden-string.py` — waive legitimate matches

Generic CLI for adding repo-relative paths to the `allowlist` of a
`[[rule]]` block in [`forbidden-strings.toml`](./forbidden-strings.toml).
Use it whenever a `check-forbidden-strings.py` violation is genuinely
intentional (historical session logs, audit trails, fixtures that *must*
cite the deprecated string verbatim, etc.).

## Why a dedicated CLI

Editing `forbidden-strings.toml` by hand is fine for a single entry,
but it's easy to:

1. Forget the justification comment, leaving a future reader staring at
   a bare path with no idea why it was waived.
2. Add a typo path that doesn't actually contain the pattern — silently
   creating a no-op waiver while the real violation continues to fail
   CI.
3. Drift the file's hand-formatted comment style by round-tripping
   through a TOML writer.

The CLI hard-blocks all three.

## Usage

### Explicit paths

```bash
python3 linter-scripts/allowlist-forbidden-string.py \
    --rule STALE-REPO-SLUG \
    --reason "Documents the v14 → v17 rebrand for the audit trail." \
    .lovable/memory/sessions/2026-04-24-batch-cleanup-and-rebrand.md
```

### Auto mode (waive every current finding under one reason)

```bash
python3 linter-scripts/allowlist-forbidden-string.py \
    --rule STALE-REPO-SLUG \
    --reason "Historical migration documentation — frozen prose." \
    --auto
```

### Preview before writing

```bash
python3 linter-scripts/allowlist-forbidden-string.py \
    --rule LEGACY-CDN-DOMAIN \
    --reason "..." --auto --dry-run
```

### Apply and re-verify in one step

```bash
python3 linter-scripts/allowlist-forbidden-string.py \
    --rule LEGACY-CDN-DOMAIN \
    --reason "..." --auto --check
```

## Guard rails

| Check | Behaviour |
|-------|-----------|
| `--reason` empty or whitespace | exit `2` |
| Unknown `--rule` id | exit `1` (lists available ids) |
| Path doesn't exist on disk | reported, exit `3` |
| Path exists but doesn't match the rule's regex | reported, exit `3` (a waiver here would be a no-op) |
| Path already in the allowlist | silently skipped (idempotent) |
| `--check` and the linter still reports findings | exit `4` |

## Output format

The reason is written verbatim as a one-line `# comment` directly above
the new path entries, so future readers always have the *why* without
needing to dig through git history:

```toml
allowlist = [
  ...
  # Documents the v14 → v17 rebrand for the audit trail.
  ".lovable/memory/sessions/2026-04-24-batch-cleanup-and-rebrand.md",
]
```

## What this CLI deliberately does **not** do

- It does not rewrite or "normalise" the existing TOML formatting.
- It does not remove paths from the allowlist (do that by hand — removal
  should always be a deliberate, reviewed action).
- It does not modify `exclude_dirs` or `exclude_files` (those are
  scan-scope decisions, not waivers).
- It does not create new `[[rule]]` blocks (rules are spec-level
  decisions; add them by hand with a code review).
