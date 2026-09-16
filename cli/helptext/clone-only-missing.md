# gitmap clone-only-missing

Clone only missing repositories from JSON manifests or URLs, skipping existing repositories on disk without pulling or modifying them.

## Usage

```bash
gitmap clone-only-missing [source] [target-dir] [flags]
gitmap com [source] [target-dir] [flags]
```

## Description

`gitmap clone-only-missing` (alias: `com`) acts identically to `gitmap clone`, reading JSON manifests or URLs, but strictly clones only repositories that are missing from disk. Any repository whose directory already exists on disk is safely skipped without executing `git pull` or altering existing local files.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--ssh` | | false | Clone using SSH transport (`git@github.com:...`) with automatic host key acceptance |
| `--https` | | false | Force HTTPS transport |
| `--target` | `-t` | current dir | Target directory to clone missing repositories into |
| `--workers` | `-w` | 4 | Max concurrent clone operations |
| `--verbose` | `-v` | false | Show verbose clone execution details |
| `--yes`, `-y` | | false | Assume yes for confirmations |
| `--help` | `-h` | | Show command help |

## Examples

```bash
# 1. Clone missing repositories from default JSON manifest
gitmap clone-only-missing

# 2. Using short alias 'com'
gitmap com

# 3. Clone missing repositories over SSH
gitmap com --ssh

# 4. Clone missing repositories from specific JSON manifest
gitmap com .gitmap/output/gitmap.json

# 5. Clone missing repository from a direct URL (skips if folder exists)
gitmap com https://github.com/alimtvnetwork/gitmap-v28

# 6. Clone missing repositories over SSH with custom concurrency
gitmap com --ssh -w 8
```

See also: `gitmap clone`, `gitmap pull-all`, `gitmap status`
