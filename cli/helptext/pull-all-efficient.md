# gitmap pull-all-efficient

Efficient multi-repo pull with automated inactivity heuristic. Analyzes pull history in `gitmap-pull.db`; if a repository has had zero commits across 20+ pulls within a 24-hour window, it is skipped.

## Alias

pae, pull-ae

## Usage

    gitmap pull-all-efficient [flags]
    gitmap pae [flags]
    gitmap pull-ae [flags]
    gitmap pull all-efficient [flags]
    gitmap pull ae [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --status | false | Render full tabular summary of pulled repositories |
| --status-table | false | Alias for --status |
| --ssh | false | Enforce SSH transport for pulled repositories |
| --https | false | Enforce HTTPS transport for pulled repositories |
| -t, --target | "" | Query remote SSH fleet machines (e.g. `-t ssh` or `-t <node>`) |
| --json | false | Emit structured JSON summary object (`PullEfficientSummary`) |

## Key Features

- **Inactivity Detection**: Evaluates the last 20 to 30 pull sessions in `gitmap-pull.db`.
- **24-Hour Active Window**: Only skips repositories that have 20 consecutive zero-change pulls within the last 24 hours.
- **Short-Form Expansion Banner**: Running `pae` or `pull-ae` announces the full command name, working directory, and GitMap version.
- **Concise Default Output**: Suppresses verbose tabular overhead by default, presenting a clean list of updated repositories and skipped inactive items.
- **Full Table on Demand**: Use `--status` or `gitmap pull-all-efficient-table` (`paet`) to render the full batch table.
- **Remote SSH JSON Query**: Use `-t ssh` to query remote SSH nodes concurrently; remote nodes return JSON summaries which are formatted and displayed locally using native `termtable`.

## Examples

### Example 1: Efficient pull via short form

```bash
gitmap pae
```

### Example 2: Efficient pull with full status table

```bash
gitmap pae --status
```

### Example 3: Remote SSH fleet query formatted as native table

```bash
gitmap pull all-efficient -t ssh
gitmap pae -t ssh
```

### Example 4: Output structured JSON for automation

```bash
gitmap pae --json
```

## See also

- `pull` — single-repo or traditional pull (defaults to pull all outside git repos)
- `pull-all` — unconditional batch pull of all tracked repositories
- `pull-all-efficient-table` — efficient pull displaying full table by default
