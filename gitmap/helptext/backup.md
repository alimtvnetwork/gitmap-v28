# `gitmap backup`

List or prune the on-disk backup tree maintained by `fix-repo`
(`undo` reads from the same tree).

Tree layout:

```
<cwd>/.gitmap/backup/<repo>/v<N>/fix-repo/<UTC-timestamp>/
```

## Subcommands

```
gitmap backup ls                              # group by repo, count + size
gitmap backup prune --keep=N                  # keep newest N snapshots per repo
gitmap backup prune --older-than=DAYS         # drop snapshots older than DAYS
gitmap backup prune --keep=5 --older-than=30  # combine — both rules apply
gitmap backup prune --dry-run ...             # print only, no disk changes
```

`prune` refuses to run without at least one of `--keep` or
`--older-than` so an accidental `gitmap backup prune` is a no-op.

## Examples

```
gitmap backup ls
gitmap backup prune --keep=10 --dry-run
gitmap backup prune --older-than=60
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0    | OK (including "nothing to prune") |
| 1    | Walk failed (permission / IO) |
| 2    | Bad flags or unknown subcommand |
