# File-Sync Watch Tasks

Manage named file-sync watch tasks for one-way folder synchronization.

## Alias

tk

## Usage

    gitmap task <subcommand> [flags]

## Subcommands

| Subcommand | Description                                               |
|------------|-----------------------------------------------------------|
| create     | Create a new sync task                                    |
| list       | List all saved tasks                                      |
| run        | Start a task's sync loop                                  |
| show       | Show details of a task                                    |
| delete     | Remove a saved task                                       |
| history    | View task execution history ([limit] [offset], e.g. -10)  |
| undo       | Revert the most recently completed task or task ID        |
| redo       | Re-execute the most recently undone task                  |
| clear      | Clear orphaned or illegal pending tasks                   |

## Flags (create)

| Flag   | Default    | Description                |
|--------|------------|----------------------------|
| --src  | (required) | Source directory path       |
| --dest | (required) | Destination directory path |

## Flags (run)

| Flag       | Default | Description                              |
|------------|---------|------------------------------------------|
| --interval | 5       | Sync interval in seconds (minimum 2)     |
| --verbose  | false   | Show detailed sync output                |
| --dry-run  | false   | Preview sync actions without copying     |

## Prerequisites

- Source directory must exist before creating or running a task
- `.gitignore` patterns in the source root are respected during sync

## Examples

### Create and verify a sync task

    $ gitmap task create my-sync --src ./src --dest ./backup
      Task 'my-sync' created.

    $ gitmap task show my-sync
      Name:     my-sync
      Source:   ./src
      Dest:     ./backup

### Run a sync task with verbose output

    $ gitmap tk run my-sync --interval 10 --verbose
      Task 'my-sync' running — syncing every 10s (Ctrl+C to stop)
      Synced: main.go
      Synced: config/settings.json
      Synced: utils/helpers.go
      All files up to date.
      ^C
      Task 'my-sync' stopped.

### List and delete tasks

    $ gitmap task list
      Tasks:
        my-sync              ./src → ./backup
        docs-mirror          ./docs → /mnt/share/docs

    $ gitmap task delete docs-mirror
      Task 'docs-mirror' deleted.

### View task execution history

    $ gitmap task history
    $ gitmap task history 50
    $ gitmap task history 50 -10
    $ gitmap task history -20

### Undo or redo a task

    $ gitmap task undo
    $ gitmap task undo 12
    $ gitmap task redo

## See Also

- [watch](watch.md) — Live-refresh dashboard of repo status
- [exec](exec.md) — Run git commands across repos

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter task
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
