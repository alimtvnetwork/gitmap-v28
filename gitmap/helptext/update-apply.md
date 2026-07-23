# gitmap update apply

Upgrade one scanned repo to the latest release tag known to gitmap.

## Alias

ua

## Usage

    gitmap update apply <repo> [flags]
    gitmap ua <repo> [flags]

`<repo>` accepts any of:

- Repo slug (`owner/name`)
- Nickname / group alias
- Absolute or relative path on disk

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --yes / -y | false | Skip the confirmation prompt |
| --dry-run | false | Print the upgrade plan without touching the working tree |
| --tag \<vX.Y.Z\> | (latest) | Target a specific tag instead of the latest known |
| --strategy fetch-checkout\|source-release | fetch-checkout | `fetch-checkout` runs `git fetch --tags && git checkout <tag>`. `source-release` delegates to the repo's own `release-self` when the row's source is `release`. |
| --stash | false | `git stash` local changes before checkout, restore on success |
| --json | false | Emit a structured result object |

## Prerequisites

- Repo appears in `gitmap list --update` (or the tag is reachable via `--tag`)
- Working tree is clean, or `--stash` is passed

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Upgrade succeeded |
| 1 | Upgrade failed; a `PendingTask` of type `Upgrade` was enqueued for retry via `do-pending` |
| 2 | Nothing to do (already on latest, or `<repo>` not scanned) |

## Examples

### Basic upgrade (interactive)

    $ gitmap update apply acme/api
      → Upgrading acme/api: v1.11.0 → v1.14.0 (3 releases)
      Proceed? (y/N): y
      ✓ Fetched tags (+4)
      ✓ Checked out v1.14.0
      ✓ Recorded completed task #48

### Non-interactive, dry run

    $ gitmap ua acme/api --dry-run
      (dry-run) Would run:
        cd /home/a/work/api
        git fetch --tags
        git checkout v1.14.0
      (dry-run) No changes written.

### Stash local edits and upgrade

    $ gitmap ua acme/api -y --stash
      ✓ Stashed 2 file(s)
      ✓ Checked out v1.14.0
      ✓ Popped stash

### Failure enqueues a pending task

    $ gitmap ua acme/api -y
      ✗ git fetch failed: could not resolve host
      → Enqueued PendingTask #61 (type=Upgrade). Retry with:
          gitmap do-pending 61

### JSON output

    $ gitmap ua acme/api -y --json
    {"repo":"acme/api","from":"v1.11.0","to":"v1.14.0","status":"ok","pendingTaskId":null}

## See Also

- [list-update](list-update.md) — Discover upgradable repos
- [update-all](update-all.md) — Upgrade every upgradable repo
- [do-pending](do-pending.md) — Retry Upgrade tasks that failed
- [pull](pull.md) — Non-tagged fast-forward alternative

## Scripting (JSON)

    gitmap help --json --filter update-apply

Schema: `spec/08-json-schemas/update-apply.schema.json` (v6.80.0+).
