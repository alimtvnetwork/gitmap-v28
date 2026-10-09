# gitmap ignore / gitmap fix-ignore-all

Audit, manage, configure, and enforce `.gitignore` hygiene across repositories with high-speed SQLite Split-DB check caching and customizable check intervals.

## Usage

```bash
# Audit or remediate ignore rules in current repository or across tracked repositories
gitmap ignore [subcommand|flags]
gitmap fix-ignore-all [flags]
gitmap fia [flags]
```

## Aliases

- `ignore`, `fix-ignore-all`, `fia`, `fias`

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `config` | Display active check interval, SQLite Split-DB cache path, and cache record metrics |
| `config set interval <dur>` | Set ignore audit frequency threshold (e.g. `24h`, `1d`, `12h`, `30m`, or `0`/`off`) |
| `cache` | List cached repository check timestamps, statuses (`clean`, `skipped`, `remediated`), and durations |
| `cache clear` | Invalidate all cached repository ignore inspection records |

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--interval <dur>` | `-i` | Override check interval for this execution (e.g. `--interval 12h`, `--interval 0` to bypass) |
| `--force` | `-f` | Bypass cache and force a complete re-scan of all repositories |
| `--yes` | `-y` | Automatically accept ignore remediation prompts without interactive confirmation |
| `--json` | `-j` | Output status, cache records, or configuration in structured JSON format |
| `--help` | `-h` | Display usage instructions |

## SQLite Split-DB Cache

To prevent redundant, heavy git subprocess invocations during rapid pull and status loops, GitMap caches inspection outcomes in an isolated SQLite Split-DB:
- **Location:** `.gitmap/data/gitignore/cache/sql.db`
- **Default TTL:** `24h` (configurable via `gitmap ignore config set interval <duration>`)
- **Status Codes:**
  - `clean`: Repository `.gitignore` is fully compliant with zero unignored tracked artifacts.
  - `skipped`: Interactive remediation was declined or deferred.
  - `remediated`: Untracked/ignored file cleanup was successfully committed.

## Examples

### 1. View Ignore Configuration and Cache Statistics

```bash
gitmap ignore config
```

Output:

```text
  GitIgnore Audit Configuration
  =============================
  • Check Interval:   24h (default: 24 hours)
  • Split-DB Path:    .gitmap/data/gitignore/cache/sql.db
  • Cached Records:   62 repositories (60 clean, 2 remediated)
  • Status:           Cache active (checks within 24h will be skipped)
```

### 2. Update Check Frequency to 12 Hours

```bash
gitmap ignore config set interval 12h
```

### 3. List All Cached Repository Checks

```bash
gitmap ignore cache
```

### 4. Force Bypass Cache and Re-audit Immediately

```bash
gitmap ignore --force
```

### 5. Automated Remediation Across Repositories (Non-Interactive)

```bash
gitmap fix-ignore-all --yes
```
