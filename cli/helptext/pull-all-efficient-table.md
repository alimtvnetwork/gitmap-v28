# gitmap pull-all-efficient-table

Efficient multi-repo pull with automated inactivity heuristic and full tabular display.

## Alias

paet

## Usage

    gitmap pull-all-efficient-table [flags]
    gitmap paet [flags]
    gitmap pull all-efficient-table [flags]
    gitmap pull paet [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --ssh | false | Enforce SSH transport for pulled repositories |
| --https | false | Enforce HTTPS transport for pulled repositories |

## Key Features

- **Automated Table View**: Always renders the full multi-column status table for active repositories.
- **Short-Form Expansion Banner**: Running `paet` displays full command, working directory, and GitMap version.
- **Inactivity Detection**: Skips repositories with 20 consecutive zero-change pulls within the last 24 hours.

## Examples

### Example 1: Efficient pull with full table via short form

    gitmap paet

## See also

- `pull-all-efficient` — efficient pull with concise default output
- `pull-all` — unconditional batch pull of all tracked repositories
