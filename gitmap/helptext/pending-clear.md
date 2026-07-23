# gitmap pending clear

Remove orphaned or illegal pending tasks so the next `clone` / `clone-next`
run is not blocked by a leftover entry from a previous crash.

## Usage

    gitmap pending clear [<mode>|<id>] [--dry-run] [--yes|-y]

## Modes

| Mode | What it removes |
|------|-----------------|
| `orphans` (default) | Pending tasks whose `TargetPath` does not exist on disk anymore |
| `illegal` | Pending tasks whose `TargetPath` looks like a URL (e.g. `https:\github.com\...`) or contains illegal Windows path chars (`:` after the drive letter, `?`, `*`, `<`, `>`, `|`, `"`) |
| `all` | Every pending task (use with care — confirm prompt always shown unless `--yes`) |
| `<id>` | A single task by numeric ID (e.g. `gitmap pending clear 17`) |

When no mode is given, `orphans` is assumed.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--dry-run` | false | Print what would be deleted, don't touch the database |
| `--yes` / `-y` | false | Skip the confirmation prompt |

## Why

The clone pipeline records a pending task before it begins, so a crash
mid-clone (file-lock, broken filesystem path, OS reboot) can leave the
DB pointing at a target that no longer makes sense. The next clone then
hits `pending task already exists for Clone at <bad-path>` and refuses
to proceed. `pending clear` is the deterministic escape hatch.

A common trigger on Windows: PowerShell silently splits unquoted
comma-separated URLs, the second URL becomes the "folder name", and
the resulting target path looks like `D:\work\https:\github.com\...` —
illegal because `:` is reserved after the drive letter. v3.85.0 added
defensive guards in the clone command, but pre-existing rows from
older crashes still need this command to clean up.

## Examples

### Clear orphans (default)

    gitmap pending clear

Output:

      ╔══════════════════════════════════════╗
      ║       gitmap pending clear           ║
      ╚══════════════════════════════════════╝

      → Mode: orphans
      → Scanned 4 pending task(s)
      • #1  type=Clone    reason=orphan-target-missing target=D:\work\old-repo
      → Delete the 1 task(s) above? (yes/N): yes
      ✓ Deleted task #1 (Clone)

      ✓ Cleared 1/4 pending task(s).

### Clear illegal targets, no prompt

    gitmap pending clear illegal --yes

### Preview only

    gitmap pending clear all --dry-run

### Clear one specific task

    gitmap pending clear 17

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success (including no matches found) |
| 1 | Database error, invalid mode/ID, or user canceled the prompt |

## See Also

- `gitmap pending` — list pending tasks
- `gitmap do-pending` (`dp`) — retry pending tasks
- `gitmap clone` — what creates pending tasks in the first place

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter pending-clear
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
