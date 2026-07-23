# gitmap commit-left

> **Status (v3.102.0):** implemented. Reuses the same Plan/Replay engine
> as `commit-right`, with source/target swapped.

Replay RIGHT's commits onto LEFT as a fresh, cleaned commit sequence.
The "-left" suffix names the **destination**, exactly like
`merge-left` writes files to LEFT.

## Alias

cml

> Spec §13 reserved `cl`, but `cl` is already taken by `changelog`.
> Use `cml` instead. The long-form `commit-left` always works.

## Usage

    gitmap commit-left LEFT RIGHT [flags]

Same endpoint syntax and flag set as `commit-right` (see
[commit-right.md](commit-right.md)). The only difference is direction:
RIGHT → LEFT instead of LEFT → RIGHT.

## Examples

    gitmap commit-left ./repo-A ./repo-B

Replays the commits **from `./repo-B`** onto `./repo-A`. The provenance
footer's `gitmap-replay-source:` line records `./repo-B` (the source).

## See Also

- [commit-right](commit-right.md) — opposite direction (full flag table)
- [commit-both](commit-both.md) — bidirectional (sequential)
- [merge-left](merge-left.md) — file-state mirror (no commit replay)
- spec/01-app/106-commit-left-right-both.md — full design

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter commit-left
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
