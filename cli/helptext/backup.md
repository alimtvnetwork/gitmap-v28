# `gitmap backup`

Create, list, restore, and prune Git-backed cloud backups and local on-disk snapshots.

## Simulation

```
$ gitmap backup create --note "Pre-migration backup"
  ✓ Cloud backup snapshot created and pushed successfully!
  ● Snapshot ID: snapshot-2026-09-03-184512
  ● Remote Repo: https://github.com/alimtvnetwork/gitmap-cloud-backup
  ● Timestamp:   2026-09-03 18:45:12
```

## Subcommands

```
gitmap backup create [--note <text>]           # Push cloud backup snapshot
gitmap backup ls [--json] [--local]            # List cloud snapshots or local tree
gitmap backup restore [1-N|snapshot-id]        # Restore databases & profiles
gitmap backup rm <1-N|snapshot-id>             # Delete remote cloud snapshot
gitmap backup status                           # Check cloud repo connection & stats
gitmap backup prune --keep=N                   # Prune local fix-repo snapshots
```

## Flags

| Flag | Description |
|------|-------------|
| `--note <text>` | Optional annotation recorded in snapshot manifest |
| `--repo <name>` | Custom remote backup repository slug |
| `--local` | List local fix-repo backup directories instead of cloud |
| `--keep=N` | Retain newest N snapshots per repo during prune |
| `--older-than=DAYS` | Drop snapshots older than specified days |
| `-y`, `--yes` | Confirm snapshot restoration or deletion non-interactively |

## Examples

```bash

# Push full backup to private cloud repository

gitmap backup create --note "Weekly system backup"

# List available cloud backup snapshots

gitmap backup ls

# Restore databases from latest or numbered snapshot

gitmap backup restore 1 -y

# Check remote cloud backup status

gitmap backup status

# Prune local fix-repo snapshots keeping newest 5

gitmap backup prune --keep=5 --local
```
