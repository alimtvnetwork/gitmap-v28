# gitmap update all

Upgrade every repo that `gitmap list --update` reports as behind. Iterates
`update apply` per row, continues past per-repo failures, and enqueues a
`PendingTask` of type `Upgrade` for each failed row.

## Alias

uall

## Usage

    gitmap update all [flags]
    gitmap uall [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --yes / -y | false | Skip the single top-level confirmation prompt |
| --dry-run | false | Print the full plan and exit; nothing is fetched or checked out |
| --only \<kind\> | (all) | Restrict to one source: `release`, `import`, or `git-tag` |
| --parallel \<n\> | 4 | Number of repos upgraded concurrently. Clamped to `[1, 16]` |
| --strategy fetch-checkout\|source-release | fetch-checkout | Passed through to each `update apply` call |
| --stash | false | Passed through to each `update apply` call |
| --stop-on-error | false | Abort the remaining queue after the first failure |
| --json | false | Emit a JSON summary object |

## Prerequisites

- A prior `gitmap scan`
- Network access (or use `--dry-run` first)

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All upgrades succeeded |
| 1 | Partial success. Failed rows are visible via `gitmap pending` |
| 2 | Nothing upgradable (list was empty) |

## Examples

### Prompted upgrade of everything

    $ gitmap update all
      → 3 repo(s) will be upgraded:
          alimtvnetwork/gitmap-v27  v6.75.0 → v6.80.0
          acme/api                  v1.11.0 → v1.14.0
          acme/web                  (new)   → v2.0.0
      Proceed? (y/N): y
      ✓ alimtvnetwork/gitmap-v27 → v6.80.0
      ✓ acme/api → v1.14.0
      ✗ acme/web → clone failed: permission denied
          Enqueued PendingTask #62 (type=Upgrade)

      ─────────────────────────────────────────────
      2 succeeded · 1 failed · 0 skipped   (see: gitmap pending)

### Dry-run, release-tracked only

    $ gitmap uall --only release --dry-run
      (dry-run) 1 repo would be upgraded:
          alimtvnetwork/gitmap-v27  v6.75.0 → v6.80.0

### Non-interactive, parallel 8, JSON

    $ gitmap uall -y --parallel 8 --json
    {"total":3,"ok":2,"failed":1,"skipped":0,"pending":[62]}

### Stop on first failure

    $ gitmap uall -y --stop-on-error

## See Also

- [list-update](list-update.md) — Preview what would be upgraded
- [update-apply](update-apply.md) — Upgrade one repo at a time
- [pending](pending.md) — Inspect failed upgrades
- [do-pending](do-pending.md) — Retry failed upgrades

## Scripting (JSON)

    gitmap help --json --filter update-all

Schema: `spec/08-json-schemas/update-apply.schema.json` (shared with
`update apply`; `update all` wraps the same result object in a summary
envelope).
