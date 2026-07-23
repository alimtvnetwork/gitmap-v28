# gitmap visibility-history (`vish`)

List the most recent `MakeAllVisibilityRun` rows newest-first — the
discovery layer behind `--run <id>` on `vu` and `vr`.

```
gitmap visibility-history [--limit <N>]
gitmap vish                 [--limit <N>]
```

## What it does

1. Queries `MakeAllVisibilityRun` ordered by `MakeAllVisibilityRunId`
   descending.
2. Prints one row per invocation with: id, kind
   (`MakeAllPublic`/`MakeAllPrivate`/`VisibilityUndo`/`VisibilityRedo`),
   owner, matched/ok/skip/fail/excluded tallies, exit code,
   `StartedAt`.
3. Empty database prints a friendly message and exits `0`.

## Flags

| Flag | Behavior |
|------|----------|
| `--limit <N>` | Cap the output at the N most recent rows (default 20). |

## Examples

```
gitmap vish                       # last 20 runs
gitmap vish --limit 5             # last 5 runs only
```

Sample output:

```
ID    Kind             Owner                 Matched  Ok  Skip Fail Excl Exit  Started
42    VisibilityRedo   alice                       3   3    0    0    0   0   2026-06-06T12:11:09Z
41    VisibilityUndo   alice                       3   3    0    0    0   0   2026-06-06T12:10:44Z
40    MakeAllPublic    alice                       3   3    0    0    0   0   2026-06-06T12:08:12Z
```

Once you have the ID, target it explicitly:

```
gitmap vu --run 40              # reverse run #40
gitmap vr --run 41 --dry-run    # preview replaying undo #41
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success (rows printed OR empty DB) |
| 5 | Audit DB open / select failed |

## See also

- `gitmap visibility-undo` (`vu`) — consumes `--run <id>` from this list.
- `gitmap visibility-redo` (`vr`) — same.
- `gitmap make-all-public` / `make-all-private` — the runs being recorded.

## Scripting (JSON)

```bash
gitmap help --json --filter visibility-history
```
