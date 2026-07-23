# release-undo (ru)

Undo a previously created release tag.

## Synopsis

```
gitmap release-undo [vX.Y.Z] [--keep-remote] [--dry-run] [-y]
gitmap ru           [vX.Y.Z] [--keep-remote] [--dry-run] [-y]
```

If no version is supplied, the newest `.gitmap/release/v*.json` is used.

## What it does

1. Deletes the local annotated tag (`git tag -d vX.Y.Z`)
2. Deletes the remote tag (`git push origin :refs/tags/vX.Y.Z`) — skip with `--keep-remote`
3. Removes the local sidecar `.gitmap/release/vX.Y.Z.json`

Prints a copy-friendly summary line that can be pasted into a task-completion
report (e.g. "✅ release-undo complete — v6.65.0 removed (local tag, remote
tag, release json)").

## Flags

- `--keep-remote` — only delete locally; leave `origin` tag intact
- `--dry-run`     — preview the steps without applying
- `-y, --yes`     — skip the confirmation prompt

## Examples

```
gitmap ru                     # undo the latest release on this repo
gitmap ru v6.65.0 -y          # undo a specific tag without prompting
gitmap ru v6.65.0 --keep-remote
gitmap ru --dry-run
```
