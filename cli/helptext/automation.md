# gitmap automation

High-performance native Go automation, multi-core search, polyglot newline normalization, and side-by-side execution benchmarks.

## Usage

```bash
gitmap automation [subcommand] [flags]
gitmap auto [subcommand] [flags]
gitmap scripts [subcommand] [flags]
```

## Description

`gitmap automation` provides high-speed compiled Go tooling replacing slow external scripts. It implements thread-safe lazy regex compilation, multi-core directory search with literal fast paths, universal CRLF-to-LF newline and trailing whitespace normalization, sub-millisecond in-memory caching (<0.05ms) with zero disk bloat, and side-by-side performance benchmarking against legacy Python scripts.

## Subcommands

- `search <pattern> [dir]` (alias: `grep`): Multi-core streaming search across repository files with lazy regex or literal match.
- `newlines [paths...]` (alias: `fix-newlines`, `lf`): Polyglot CRLF to LF and trailing whitespace normalizer.
- `cache [status|read|warm|clear]`: Manage sub-millisecond in-memory file cache with zero disk pollution.
- `benchmark [target]` (alias: `bench`, `compare`): Run side-by-side Go vs Python execution benchmarks.

## Flags

### Search Flags (`gitmap automation search`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--regex` | `-r` | `false` | Enable regular expression matching |
| `--ignore-case` | `-i` | `false` | Perform case-insensitive matching |
| `--ext` | `-e` | `nil` | Filter search by file extensions (e.g. `.go, .ts`) |
| `--workers` | `-w` | `CPU count` | Number of concurrent worker threads |

### Newlines Flags (`gitmap automation newlines`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--fix` | `-f` | `false` | Apply normalized line endings and whitespace to disk |
| `--dry-run` | | `false` | Preview modified file counts without writing to disk |

## Examples

```bash
# 1. Multi-core regex search across the codebase with lazy regex compilation
gitmap automation search "func Run" cli --ext .go

# 2. Case-insensitive literal search bypassing the regex engine
gitmap automation search "connection reset" --ignore-case

# 3. Polyglot newline normalization in dry-run mode
gitmap automation newlines --dry-run

# 4. Atomically normalize CRLF to LF and trim whitespace across repository
gitmap automation newlines --fix

# 5. Pre-warm and inspect sub-millisecond in-memory cache
gitmap automation cache warm cli
gitmap automation cache status

# 6. Read file directly from memory cache (<0.05ms)
gitmap automation cache read llm.md

# 7. Run side-by-side Go vs Python performance benchmark
gitmap automation benchmark all
gitmap automation benchmark search
gitmap automation benchmark newlines
```

See also: `gitmap find-files`, `gitmap find-regex-read`, `gitmap os`
