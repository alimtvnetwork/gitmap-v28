# gitmap visibility-undo (`vu`)

Reverse a previous bulk visibility run by re-applying each repo's
captured `PrevVisibility`. The undo is itself logged as a fresh run
with `CommandKind='VisibilityUndo'`, so `gitmap visibility-redo` can
roll it back and `vu` again undoes the redo.

```
gitmap visibility-undo [--run <id>] [--force] [--dry-run] [--verbose]
gitmap vu             [--run <id>] [--force] [--dry-run] [--verbose]
```

## What it does

1. Picks the source run — either `--run <id>` (exact lookup in
   `MakeAllVisibilityRun`) or the most recent run with `OkCount>0`.
2. Preflights the provider CLI: `gh`/`glab` on PATH AND authenticated.
3. For each persisted `Ok` row from the source run:
   - **Drift guard** (default): reads the *current* visibility and
     skips with `DRIFT SKIP` if it no longer matches the `NewVisibility`
     we set last time — i.e. someone (or something) flipped that repo
     out-of-band after the original run.
   - Re-applies the row's `PrevVisibility` via the provider CLI.
   - Verifies the change took effect.
4. Persists a new `MakeAllVisibilityRun` row (`CommandKind='VisibilityUndo'`,
   `PatternList='visibility-undo:source-run=<id>'`) so the chain stays
   auditable through `gitmap visibility-history`.

## Flags

| Flag | Behavior |
|------|----------|
| `--run <id>` | Reverse the run with this exact `MakeAllVisibilityRunId` (discover IDs via `vh`). |
| `--force` | Bypass the drift guard. Reverses even when current visibility differs from `NewVisibility` — destroys out-of-band changes. |
| `--dry-run` | Print the planned per-repo reversal without touching the provider. |
| `--verbose` | Echo every shell command to stderr before running it. |

## Examples

```
gitmap vu                       # reverse the most recent flip
gitmap vu --dry-run             # preview before committing
gitmap vu --run 42              # reverse run #42 specifically
gitmap vu --run 42 --force      # ignore drift skips
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | All reversals succeeded (or every row was a no-op/drift skip) |
| 5 | Provider CLI missing/unauth OR every reversal failed OR audit DB failed |
| 6 | Bad `--run <id>` value (missing, non-numeric, ≤0) |
| 7 | No undoable run found |
| 9 | Partial — at least one reversal succeeded AND at least one failed |

## See also

- `gitmap visibility-redo` (`vr`) — roll a `vu` back.
- `gitmap visibility-history` (`vish`) — discover `--run <id>` values.
- `gitmap make-all-public` / `make-all-private` — the runs `vu` reverses.

## Scripting (JSON)

```bash
gitmap help --json --filter visibility-undo
```
