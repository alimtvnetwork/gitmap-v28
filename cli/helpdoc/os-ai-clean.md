# gitmap os ai-clean

Scan, analyze, and safely purge Google Antigravity AI caches, conversation transcripts, system generated artifacts, and temporary installer files.

## Usage

```bash
gitmap os ai-clean [flags]
gitmap os clean [flags]
gitmap ai-clean [flags]
```

## Description

`gitmap os ai-clean` identifies and safely reclaims disk space consumed by autonomous AI coding workflows and installer caches. Over time, extensive AI agent sessions can accumulate gigabytes of conversation transcripts, reasoning logs, and temporary sandbox files.

### Cleaned Target Categories
1. **Google Antigravity Brain Caches**: Subagent caches and workspace traces under `~/.gemini/antigravity/brain/*/`.
2. **System Generated Tasks**: Agent logs and execution transcripts under `.system_generated/tasks/` and `.system_generated/logs/`.
3. **Operating System Temp AI Dumps**: Transient temporary directories (`gemini-*`, `antigravity-*`) in system temp space.
4. **GitMap Installer Caches**: Cached binary releases and temporary archives under `~/.gitmap/cache/` and `~/.gitmap/temp/`.

### Interactive Safety Guardrails
Before deleting any files, `gitmap os ai-clean` renders a preflight table showing each target path, matched file count, and total reclaimable bytes, followed by an interactive confirmation prompt:
```text
Proceed with AI cache cleanup? [y/N]:
```
In automated CI/CD or scripting environments, the interactive prompt can be bypassed using `-y` or `--yes`.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--yes` | `-y` | false | Bypass interactive `[y/N]` preflight confirmation and proceed with deletion |
| `--dry-run` | `-n` | false | Simulate cache scanning and display reclaimable sizes without deleting files |
| `--json` | | false | Output scan results and freed space metrics as structured JSON |
| `--help` | `-h` | | Show command help |

## Standalone Python Utility
A zero-dependency Python 3 counterpart is also available in `scripts/os-ai-clean.py` across `gitmap` and `scripts-fixer` for systems without the compiled binary:
```bash
python scripts/os-ai-clean.py --dry-run
python scripts/os-ai-clean.py -y
```

## Examples

```bash
# 1. Interactive scan and cleanup with preflight confirmation table
gitmap os ai-clean

# 2. Simulate cleanup without deleting any files
gitmap os ai-clean --dry-run
gitmap os ai-clean -n

# 3. Non-interactive cleanup for automated cron jobs or CI/CD pipelines
gitmap os ai-clean -y
gitmap os ai-clean --yes

# 4. Output scanned target statistics as structured JSON
gitmap os ai-clean --dry-run --json

# 5. Run via top-level or clean aliases
gitmap ai-clean -y
gitmap os clean
```

See also: `gitmap os`, `gitmap os storage`, `gitmap agy`
