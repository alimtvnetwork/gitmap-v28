# Clear Release JSON

## Purpose

Remove a specific release metadata JSON file from the `.gitmap/release/` directory. Supports a `--dry-run` flag to preview the target file without deleting it.

## Command

    gitmap-v27 clear-release-json <version> [--dry-run]

## Alias

    crj

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dry-run` | bool | `false` | Preview which file would be removed without deleting it |

## Version Resolution

The `<version>` argument is parsed through `release.Parse`, which applies standard semver normalisation:

1. A leading `v` prefix is optional — `2.20.0` and `v2.20.0` are equivalent.
2. Partial versions are zero-padded — `v2` becomes `v2.0.0`, `v2.1` becomes `v2.1.0`.
3. Pre-release suffixes (e.g. `v3.0.0-rc.1`) are preserved as-is.

The resolved version determines the file path: `.gitmap/release/vX.Y.Z.json`.

## Behaviour

### Normal Mode

1. Parse and validate the version argument.
2. Construct the path `.gitmap/release/vX.Y.Z.json`.
3. Check the file exists — if not, print an error and exit 1.
4. Remove the file.
5. Print a success message.

### Dry-Run Mode

1. Parse and validate the version argument.
2. Construct the path `.gitmap/release/vX.Y.Z.json`.
3. Check the file exists — if not, print an error and exit 1.
4. Print `[dry-run] Would remove <path>` and exit 0 without deleting.

## Edge Cases

| Scenario | Behaviour |
|----------|-----------|
| No version argument | Print usage message and exit 1 |
| Invalid version string (e.g. `abc`) | Print `Error: '<input>' is not a valid version.` and exit 1 |
| File does not exist | Print `Error: no release file found for vX.Y.Z` and exit 1 |
| File exists but is read-only | `os.Remove` fails; print `Error: could not remove release file: <err>` and exit 1 |
| `--dry-run` with missing file | Same as normal missing-file error — dry-run still validates existence |
| `--dry-run` with valid file | Print preview message and exit 0; file is untouched |
| Partial version `v2` | Normalised to `v2.0.0`; targets `.gitmap/release/v2.0.0.json` |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | File removed successfully, or dry-run preview printed |
| 1 | Missing argument, invalid version, file not found, or removal failed |

## Output Formats

This command produces only terminal output. It does not support `--json` or `--csv`.

| Constant | Format String |
|----------|---------------|
| `MsgClearReleaseDone` | `✓ Removed .gitmap/release/%s.json` |
| `MsgClearReleaseDryRun` | `[dry-run] Would remove %s` |
| `ErrClearReleaseUsage` | `Usage: gitmap-v27 clear-release-json <version> [--dry-run]` |
| `ErrClearReleaseNotFound` | `Error: no release file found for %s` |
| `ErrClearReleaseFailed` | `Error: could not remove release file: %v` |

## Implementation

| File | Responsibility |
|------|----------------|
| `cmd/clearreleasejson.go` | Flag parsing (`parseClearReleaseJSONFlags`) and handler (`runClearReleaseJSON`) |
| `constants/constants_messages.go` | All message and error format strings |
| `release/metadata.go` | `ReleaseExists`, `metaFilePath` — shared path construction |
| `release/semver.go` | `Parse` — version normalisation and validation |
| `helptext/clear-release-json.md` | Embedded help text displayed with `--help` |

## See Also

- [release](12-release-command.md) — Create a release
- [release data model](13-release-data-model.md) — `.gitmap/release/` file layout and schemas
- [list-releases](21-list-releases.md) — Show stored releases
- [CLI interface](02-cli-interface.md) — Full command reference
