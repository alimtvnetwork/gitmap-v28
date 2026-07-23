# gitmap unzip-compact

**Aliases:** `uzc`

Extract a local archive, an HTTP(S) archive URL, or auto-pick the single
archive in the current folder, into a single normalized output folder.

## Usage

```bash
gitmap unzip-compact [src] [dest] [flags]
gitmap uzc          [src] [dest] [flags]
```

- `src`  — local path, https URL, or omitted to auto-detect a single
  archive in the current folder.
- `dest` — destination *parent* folder. Output lands in
  `<dest>/<archive-base-name>/`. Defaults to the current folder.

## Flags

| Flag           | Description                                       |
|----------------|---------------------------------------------------|
| `--list`, `-l` | Print archive entries and exit (no extraction).   |

## Compact-extract algorithm

The extractor wraps everything in `<dest>/<archive-base-name>/` and
flattens up to **4 layers** of duplicate single-child folders, so:

| Archive shape              | Output                       |
|---------------------------|------------------------------|
| `xap.zip → xap/xap/<files>` | `dest/xap/<files>` (1 flatten) |
| `xlt.zip → <files>`         | `dest/xlt/<files>` (wrapped)   |
| `mixed.zip → README + src/` | `dest/mixed/{README,src}`      |

## Format support

- **Read + write:** zip, tar, tar.gz, tar.bz2, tar.xz, tar.zst, gz, bz2,
  xz, zst.
- **Read-only:** 7z, rar.

## History

Every invocation writes one row to the `ArchiveHistory` SQLite table
including command name, source(s), output path, format, status, and
timestamps.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter unzip-compact
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap unzip-compact
```
