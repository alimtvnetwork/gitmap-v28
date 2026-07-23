# gitmap make-all-public

Bulk-flip every matching repo on an owner to **public** in one call.

```
gitmap make-all-public <owner-or-url> <patterns> [-Y|--yes] [--verbose]
```

The companion of `gitmap make-all-private`. Both commands share one
audit table (`MakeAllVisibilityRun` / `MakeAllVisibilityResult`) so
every flip is reversible via `gitmap visibility-undo`.

## What it does

1. Resolves the owner from `<owner-or-url>` (bare owner like `alice`,
   a profile URL, or any repo URL whose owner segment we can parse).
2. Verifies that `gh` / `glab` is on `PATH` **and** authenticated
   (`gh auth status` / `glab auth status` — fails fast on missing auth).
3. Lists every repo under the owner (capped at 1000 per call).
4. Matches the comma-separated `<patterns>` (glob + `!negation`
   supported) against repo names and prints the match table.
5. Prompts for confirmation unless `-Y` is passed. The prompt accepts
   `y`, `n`, or exclusion expressions like `1,3-5` to drop entries.
6. Per repo: reads current visibility → skips if already public →
   applies via the provider CLI → verifies the change took effect.
7. Persists every step in the audit DB so `vu`, `vr`, and `vish` can
   replay or inspect the run later.

## Flags

| Flag | Behavior |
|------|----------|
| `-Y`, `--yes` | Skip the interactive confirmation prompt. |
| `--verbose` | Echo every shell command to stderr before running it. |

## Examples

```
gitmap make-all-public alice "demo-*"
gitmap make-all-public alice "demo-*,proto-*,!proto-secret"
gitmap make-all-public https://github.com/alice "*" -Y --verbose
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | All target repos already public OR all flips succeeded |
| 4 | Owner resolution / provider detection failed |
| 5 | Provider CLI missing OR not authenticated OR every repo failed |
| 6 | Bad flag / missing positional argument |
| 7 | User aborted at the confirmation prompt |
| 9 | Partial — at least one repo flipped AND at least one failed |

## See also

- `gitmap make-all-private` — opposite direction, same machinery.
- `gitmap visibility-undo` (`vu`) — reverse the most recent run.
- `gitmap visibility-redo` (`vr`) — replay an undone run.
- `gitmap visibility-history` (`vish`) — list past runs and their IDs.

## Scripting (JSON)

```bash
gitmap help --json --filter make-all-public
```
