# gitmap amend

Rewrite commit author information for one or more commits.

## Alias

am

## Usage

    gitmap amend [commit-hash] --name <name> --email <email> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --name \<name\> | — | New author name |
| --email \<email\> | — | New author email |
| --branch \<name\> | current | Target branch |
| --dry-run | false | Preview without executing |
| --force | false | Skip confirmation prompt |

## Prerequisites

- Must be inside a Git repository

## Examples

### Example 1: Amend the last commit's author

    gitmap amend --name "John Doe" --email "john@example.com"

**Output:**

    Amending commit abc1234...
    Before: Old Author <old@email.com>
    After:  John Doe <john@example.com>
    Continue? [y/N]: y
    ✓ 1 commit amended on branch 'main'

### Example 2: Amend a specific commit

    gitmap amend abc1234 --name "Jane Smith" --email "jane@company.com"

**Output:**

    Amending commit abc1234...
    Before: Wrong Name <wrong@email.com>
    After:  Jane Smith <jane@company.com>
    Continue? [y/N]: y
    ✓ 1 commit amended on branch 'main'

### Example 3: Dry-run preview (no changes)

    gitmap am abc1234 --name "Jane" --email "jane@co.com" --dry-run

**Output:**

    [DRY RUN] Amending commit abc1234
    [DRY RUN] Before: Old Name <old@email.com>
    [DRY RUN] After:  Jane <jane@co.com>
    [DRY RUN] Branch: main
    No changes made.

### Example 4: Force amend without confirmation

    gitmap amend --name "CI Bot" --email "ci@example.com" --force

**Output:**

    Amending commit def5678...
    ✓ 1 commit amended (forced, no confirmation)

## See Also

- [amend-list](amend-list.md) — View previous amendment records
- [history](history.md) — View command history

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter amend
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
