# gitmap db-listing

List, discover, and inspect registered and split SQLite databases across GitMap.

## Aliases

`dbl`, `dbls`, `db-list`

## Usage

```bash
gitmap db-listing [flags]
gitmap dbl [flags]
gitmap db listing [flags]
```

## Description

`gitmap db-listing` provides a unified inventory of all SQLite databases managed or tracked by GitMap. It displays aliases, database categories, disk usage sizes, table counts, and filesystem paths.

The command tracks:
- **Master DB**: Central SQLite database (`gitmap.db`) storing repository definitions and configurations.
- **Installation DB**: Global package and toolchain registry (`installation.db`).
- **Split-DB Tiers**: Subsystem databases including `tasks`, `pipeline`, `heatmap`, `history`, and `search`.
- **Custom Aliases**: User-defined aliases registered via `gitmap db-read <path> -as <alias>`.

## Flags

| Flag | Default | Description |
|---|---|---|
| `-a, --all` | false | Include discovered split and profile databases beyond registered aliases |
| `-r, --refresh` | false | Refresh file size stats and table counts directly from disk |
| `-j, --json` | false | Output results in structured JSON envelope format |
| `-h, --help` | false | Show command help text |

## Flags Detail

### `--all` (`-a`)
By default, `db-listing` lists the primary registered aliases. Adding `--all` recursively discovers split databases located under `.gitmap/` or standard data storage paths.

### `--refresh` (`-r`)
Queries the filesystem and database headers to recalculate real-time file sizes and table counts, updating registry cache entries.

### `--json` (`-j`)
Formats database metrics into a structured JSON envelope containing arrays of databases with their names, categories, paths, sizes in bytes, and table counts.

## Examples

### List registered databases in tabular view

```bash
gitmap db-listing
```

### Discover all split and profile databases

```bash
gitmap db-listing --all
```

### Refresh disk size stats and table counts

```bash
gitmap dbl --refresh
```

### Output database listing in JSON format

```bash
gitmap db-listing --json
```

### Combined flags using short alias

```bash
gitmap dbl -a -r
```

## See Also

- `gitmap db` — Inspect master database architecture and disk usage
- `gitmap db-read` — Read and inspect SQLite databases
- `gitmap db-view` — Query multi-table joins and filtered views

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter db-listing
```
