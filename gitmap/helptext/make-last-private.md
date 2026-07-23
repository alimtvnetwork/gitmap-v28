# `gitmap make-last-public` / `make-last-private`

Flip visibility on exactly one repo: the highest `-vN` sibling under
`<base>` for the given owner. Short aliases: `MLPUB`, `MLPRI`.

## Usage

```
gitmap make-last-public  <owner-or-url> <base>  [-Y|--yes]
gitmap make-last-private <owner-or-url> <base>  [-Y|--yes]
gitmap MLPUB             <owner-or-url> <base>  [-Y|--yes]
gitmap MLPRI             <owner-or-url> <base>  [-Y|--yes]
```

## How `<base>` is resolved

1. If `<base>` already matches `<name>-v<N>` it's used verbatim.
2. Otherwise gitmap consults the local `OwnerRepoNameIndex` (populated
   by every `make-all-*` run and any prior `make-last-*` run) for the
   highest `-vN` row whose `BaseName == <base>`.
3. On a cache miss the owner repo list is refreshed (warming the
   index) and step 2 is retried.

## Examples

```
# Flip the newest macro-ahk-vN repo public:
gitmap make-last-public https://github.com/alimtvnetwork macro-ahk

# Same, scripted (no prompt):
gitmap MLPUB alimtvnetwork macro-ahk -Y

# Exact form — no lookup needed:
gitmap make-last-private alimtvnetwork macro-ahk-v51
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0    | Changed, skipped (already at target), or no match |
| 2    | Missing arg / bad flag |
| 5    | User declined confirmation |
| 7    | Apply failed |
