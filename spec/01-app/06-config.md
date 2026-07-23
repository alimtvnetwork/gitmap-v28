# Config

## File Location

Default: `./data/config.json`
Override: `--config <path>` flag.

## Schema

```json
{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": "./.gitmap/output",
  "excludeDirs": [".cache", "node_modules", "vendor", ".venv"],
  "notes": "",
  "release": {
    "targets": [
      {"goos": "windows", "goarch": "amd64"},
      {"goos": "linux", "goarch": "amd64"},
      {"goos": "darwin", "goarch": "arm64"}
    ],
    "checksums": true,
    "compress": false
  }
}
```

## Fields

| Field         | Type          | Default            | Description                              |
|---------------|---------------|--------------------|------------------------------------------|
| defaultMode   | string        | "https"            | "https" or "ssh"                         |
| defaultOutput | string        | "terminal"         | "terminal", "csv", or "json"             |
| outputDir     | string        | "./.gitmap/output"  | Where output files are written           |
| excludeDirs   | []string      | []                 | Directory names to skip                  |
| notes         | string        | ""                 | Default note for all records             |
| release       | ReleaseConfig | {}                 | Release-specific settings (see below)    |

### Release Config

| Field     | Type           | Default          | Description                                    |
|-----------|----------------|------------------|------------------------------------------------|
| targets   | []ReleaseTarget| [] (all 6)       | Override cross-compile OS/arch matrix           |
| checksums | bool           | false            | Generate SHA256 checksums.txt for assets        |
| compress  | bool           | false            | Wrap assets in .zip (Windows) or .tar.gz        |

Each `ReleaseTarget` has `goos` (string) and `goarch` (string) fields.

## Merge Rules

1. Load config file (if it exists).
2. Apply CLI flags on top — flags always win.
3. If config file is missing, use built-in defaults silently.
