# gitmap amend-list

List previous author amendment records.

## Alias

al

## Usage

    gitmap amend-list [--json]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |

## Prerequisites

- Run `gitmap amend` at least once (see amend.md)

## Examples

### Example 1: List all amendments

    gitmap amend-list

**Output:**

     #  COMMIT    OLD AUTHOR                    NEW AUTHOR                    DATE
     1  abc1234   old@email.com                 john@example.com              2025-03-10
     2  def5678   other@email.com               jane@company.com              2025-03-09
     3  ghi9012   wrong-name@email.com          ci-bot@internal.com           2025-03-08
    3 amendments recorded

### Example 2: JSON output for scripting

    gitmap al --json

**Output:**

    [
      {
        "commit": "abc1234",
        "old_name": "Old Author",
        "old_email": "old@email.com",
        "new_name": "John Doe",
        "new_email": "john@example.com",
        "date": "2025-03-10T14:30:00Z"
      },
      {
        "commit": "def5678",
        "old_name": "Other Person",
        "old_email": "other@email.com",
        "new_name": "Jane Smith",
        "new_email": "jane@company.com",
        "date": "2025-03-09T10:15:00Z"
      }
    ]

### Example 3: No amendments recorded

    gitmap amend-list

**Output:**

    No amendments recorded.
    → Use 'gitmap amend --name "Name" --email "email"' to amend a commit.

## See Also

- [amend](amend.md) — Rewrite commit author information
- [history](history.md) — View command history

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter amend-list
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
