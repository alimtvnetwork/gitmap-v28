# gitmap make-all-private

Bulk-flip every matching repo on an owner to **private** in one call.
Mirror image of `gitmap make-all-public` — same flags, same exit
codes, same audit trail, same undo/redo/history wiring.

```
gitmap make-all-private <owner-or-url> <patterns> [-Y|--yes] [--verbose]
```

## What it does

1. Resolves the owner from `<owner-or-url>`.
2. Preflights the provider CLI: `gh` / `glab` must be on `PATH` AND
   authenticated (`gh auth status` / `glab auth status`).
3. Lists every repo under the owner (≤1000 per call) and matches it
   against the comma-separated `<patterns>` (glob + `!negation`).
4. Prints the match table and asks for confirmation (skip with `-Y`).
   The prompt accepts `y`, `n`, or exclusion expressions like `1,3-5`.
5. Per repo: reads current visibility → skips if already private →
   applies → verifies. Each step is persisted in the audit DB.

## Flags

| Flag | Behavior |
|------|----------|
| `-Y`, `--yes` | Skip the interactive confirmation prompt. |
| `--verbose` | Echo every shell command to stderr before running it. |

## Examples

```
gitmap make-all-private alice "archive-*"
gitmap make-all-private alice "old-*,!old-keep" -Y
gitmap make-all-private https://gitlab.com/alice "*" --verbose
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | All target repos already private OR all flips succeeded |
| 4 | Owner resolution / provider detection failed |
| 5 | Provider CLI missing OR not authenticated OR every repo failed |
| 6 | Bad flag / missing positional argument |
| 7 | User aborted at the confirmation prompt |
| 9 | Partial — at least one repo flipped AND at least one failed |

## See also

- `gitmap make-all-public` — opposite direction.
- `gitmap visibility-undo` (`vu`) — reverse the most recent run.
- `gitmap visibility-redo` (`vr`) — replay an undone run.
- `gitmap visibility-history` (`vish`) — list past runs and their IDs.

## Scripting (JSON)

```bash
gitmap help --json --filter make-all-private
```
