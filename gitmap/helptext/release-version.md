# release-version

Install a specific version of gitmap using the pinned release-version scripts.

## Synopsis

**PowerShell (Windows):**
```powershell
release-version.ps1 -Version <tag> [-AllowFallback] [-Quiet] [-NoPath] [-NoSelfInstall] [-InstallDir <path>]
```

**Bash (Linux/macOS):**
```bash
release-version.sh --version <tag> [--allow-fallback] [--quiet] [--no-path] [--no-self-install] [--dir <path>]
```

## What it does

1. **Version validation** — Verifies the requested version exists on GitHub releases.
2. **OS/Arch detection** — Automatically detects your platform and architecture.
3. **Asset download** — Fetches the correct binary archive from the release.
4. **SHA256 verification** — Validates the download against published checksums.
5. **Extraction** — Unpacks to the install directory.
6. **Self-install chain** — Runs `gitmap self-install` to finalize PATH and cleanup.

## Missing version behavior

If the requested version is **not a published release**:

| Context | Behavior |
|---------|----------|
| **Interactive terminal** | Prompts with the 5 most recent releases to choose from (or N to quit) |
| **Non-interactive** (no TTY, piped input, CI) | **Exits with code 1** immediately |

**Important:** Piped installs like `irm ... | iex` or `curl ... | bash` are
non-interactive. If the version is missing, the script will exit 1 without
prompting. Use `--allow-fallback` for automated environments.

## Flags

| Flag | Description |
|------|-------------|
| `-Version` / `--version` | **(Required)** The version tag to install (e.g., `v3.38.0`) |
| `-AllowFallback` / `--allow-fallback` | If version missing, use newest patch in same `vMAJOR.MINOR` series without prompting |
| `-Quiet` / `--quiet` | Suppress all prompts and progress output |
| `-NoPath` / `--no-path` | Skip PATH modification |
| `-NoSelfInstall` / `--no-self-install` | Download and extract only; skip `gitmap self-install` chain |
| `-InstallDir` / `--dir` | Custom install directory override |

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Installed and verified successfully |
| 1 | Requested version missing (and no fallback selected/allowed) |
| 2 | Network or download error |
| 3 | Checksum verification failed |
| 4 | OS/architecture not supported |
| 5 | PATH update failed (warning only) |
| 6 | Self-install chain failed |
| 7 | Verified version mismatch |

## Examples

**Install specific version (interactive prompt on missing):**
```powershell
release-version.ps1 -Version v3.38.0
```

**Install with automatic fallback (CI-friendly):**
```powershell
release-version.ps1 -Version v3.38.0 -AllowFallback
```

**Silent install with custom directory:**
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v27/releases/download/v3.38.0/release-version-v3.38.0.sh | bash -s -- --version v3.38.0 --quiet --dir /opt/gitmap
```

**Pinned snapshot install (recommended for reproducible builds):**
```powershell
irm https://github.com/alimtvnetwork/gitmap-v27/releases/download/v3.38.0/release-version-v3.38.0.ps1 | iex
```

**Generic install with version parameter:**
```powershell
irm https://gitmap.dev/scripts/release-version.ps1 | iex; Install-Gitmap -Version v3.38.0
```

## See also

- `self-install` — Re-install the current binary
- `update` — Pull the latest build from the source repo
- [Spec 105: Release-Version Script](../spec/01-app/105-release-version-script.md) — Full specification

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter release-version
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
