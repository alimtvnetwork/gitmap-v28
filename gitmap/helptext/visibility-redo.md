# gitmap visibility-redo (`vr`)

Replay the visibility state that a `gitmap visibility-undo` reverted.
`vr` is the mirror image of `vu` — same flag surface, same drift
guard, same audit chain. It targets the most recent run with
`CommandKind='VisibilityUndo'` (or an exact `--run <id>`).

```
gitmap visibility-redo [--run <id>] [--force] [--dry-run] [--verbose]
gitmap vr             [--run <id>] [--force] [--dry-run] [--verbose]
```

## What it does

1. Picks the source `VisibilityUndo` run (latest with `OkCount>0`, or
   the exact `--run <id>`).
2. Preflights the provider CLI: `gh`/`glab` on PATH AND authenticated.
3. For each persisted `Ok` row from that undo run:
   - **Drift guard** (default): if current visibility no longer matches
     the value the undo set, skip with `DRIFT SKIP` — unless `--force`.
   - Re-applies the row's `PrevVisibility` (which, for an undo source
     row, is the visibility that existed *before* the undo, i.e. the
     post-original-flip state — so the redo restores it).
4. Logs the operation as a fresh run with `CommandKind='VisibilityRedo'`
   and `PatternList='visibility-redo:source-run=<id>'`.

## Flags

| Flag | Behavior |
|------|----------|
| `--run <id>` | Replay this exact `MakeAllVisibilityRunId` (must be a `VisibilityUndo` row). |
| `--force` | Bypass the drift guard. Re-applies even when current state differs from what the undo recorded. |
| `--dry-run` | Print the planned per-repo replay without touching the provider. |
| `--verbose` | Echo every shell command to stderr before running it. |

## Examples

```
gitmap vr                       # replay the most recent undo
gitmap vr --dry-run             # preview before committing
gitmap vr --run 57              # replay undo run #57 specifically
gitmap vr --force               # ignore drift skips
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | All replays succeeded (or every row was a no-op/drift skip) |
| 5 | Provider CLI missing/unauth OR every replay failed OR audit DB failed |
| 6 | Bad `--run <id>` value |
| 7 | No undoable `VisibilityUndo` run found |
| 9 | Partial — at least one replay succeeded AND at least one failed |

## See also

- `gitmap visibility-undo` (`vu`) — the run kind `vr` consumes.
- `gitmap visibility-history` (`vish`) — list runs to find `--run <id>`.

## Scripting (JSON)

```bash
gitmap help --json --filter visibility-redo
```
