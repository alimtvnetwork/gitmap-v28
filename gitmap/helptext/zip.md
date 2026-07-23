# gitmap zip

Create an archive from one or more sources. Sources may be local paths,
HTTP(S) URLs (downloaded to a temp dir first), or git URLs (shallow-
cloned to a temp dir first). The output format is derived from the
`--out` extension.

## Usage

```bash
gitmap zip <src...> --out <file> [flags]
```

## Flags

| Flag                | Description                                           |
|---------------------|-------------------------------------------------------|
| `--out`, `-o`       | Output file path (extension drives format selection). |
| `--best`            | Highest compression (slowest, smallest).              |
| `--standard`, `-s`  | Default compression (the default when no flag set).   |
| `--fast`            | Fastest compression (largest output).                 |
| `--include <globs>` | Comma-separated; only matching paths are included.    |
| `--exclude <globs>` | Comma-separated; matching paths are skipped.          |

`--best`, `--fast`, and `--standard` are mutually exclusive.

## Format support

| Extension     | Writeable | Notes                            |
|---------------|-----------|----------------------------------|
| `.zip`        | yes       | DEFLATE, selective compression.  |
| `.tar`        | yes       | Uncompressed.                    |
| `.tar.gz`     | yes       | Best=9, Standard=−1, Fast=1.     |
| `.tar.bz2`    | yes       | Best=9, Standard=6, Fast=1.      |
| `.tar.xz`     | yes       | Library default level.           |
| `.tar.zst`    | yes       | Library default level.           |
| `.7z`, `.rar` | **no**    | Read-only — use zip or tar.\*.   |

## Examples

```bash
gitmap zip ./src ./docs --out release.zip --best
gitmap zip https://example.com/data.zip --out bundle.tar.gz
gitmap zip git@github.com:foo/bar.git --out bar.zip
```

## History

Every invocation writes one row to the `ArchiveHistory` SQLite table
including the original sources (not the resolved temp paths), output
path, format, compression mode, status, and timestamps.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter zip
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
