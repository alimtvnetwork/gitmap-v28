# gitmap install muse

Automated cross-platform installation and verification for the Meta Muse AI development companion agent.

## Usage

```bash
gitmap install muse [flags]
gitmap install meta-muse [flags]
gitmap muse install [flags]
```

## Description

The `gitmap install muse` command provisions the official Meta Muse AI development companion binary across Windows, Linux, and macOS workstations. It detects host architecture, runs official authenticated installation scripts, verifies execution integrity, registers the binary in system PATH, and persists installation telemetry to the local GitMap Split-DB store.

## Platform Installation Channels

### Windows

Runs through PowerShell (`pwsh` or `powershell.exe`) with execution policy bypass:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://dev.meta.ai/install.ps1 | iex"
```

Binary destination: `$LOCALAPPDATA\Programs\Muse\muse.exe` or `%USERPROFILE%\.muse\bin\muse.exe`.

### Linux (Ubuntu / Debian / Fedora)

Executes official shell bootstrap pipeline with automatic architecture detection (`amd64` / `arm64`):

```bash
curl -fsSL https://dev.meta.ai/install.sh | bash
```

Binary destination: `~/.local/bin/muse` or `/usr/local/bin/muse`.

### macOS (Darwin)

Executes Unix bootstrap pipeline with Homebrew fallback support:

```bash
curl -fsSL https://dev.meta.ai/install.sh | bash
# Homebrew fallback
brew install meta/muse/muse
```

Binary destination: `/opt/homebrew/bin/muse` or `~/.local/bin/muse`.

## Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--version <tag>` | `-v <tag>` | string | (latest) | Pin and install a specific Meta Muse release version |
| `--dry-run` | `-n` | bool | `false` | Preview download URL and execution command without modifying system |
| `--force` | `-f` | bool | `false` | Force reinstallation even if Meta Muse is already detected |
| `--json` | `-j` | bool | `false` | Output machine-readable JSON telemetry and installation status |
| `--verify-only` | | bool | `false` | Check existing installation path and verify `muse --version` |
| `--install-dir <path>` | | string | (auto) | Custom binary installation directory target |

## Examples

### Default Automated Installation

```bash
gitmap install muse
```

### Preview Installation Command (Dry-Run)

```bash
gitmap install muse --dry-run
```

### Install Specific Pinned Version

```bash
gitmap install muse --version 1.4.3
```

### Force Reinstallation

```bash
gitmap install muse --force
```

### Machine-Readable JSON Output

```bash
gitmap install muse --dry-run --json
```

### Verify Existing Binary

```bash
gitmap muse status
```
