# storage

Display disk drives, storage volumes, mount points, and filesystem space utilization.

## Aliases

`disk`, `df`

## Usage

    gitmap storage [flags]

## Description

The `storage` command inspects local hard drives, SSDs, partitions, and mounted volumes cross-platform on Windows, macOS, and Linux. It reports total capacity, used storage, available free space, and percentage utilization.

## Options

    -a, --all      Include hidden, system, and virtual drives
    -j, --json     Output storage statistics in JSON format

## Examples

```bash
# View storage consumption across all physical volumes
gitmap storage

# View storage metrics including system and virtual partitions
gitmap storage --all

# Output storage statistics as JSON
gitmap storage --json
```
