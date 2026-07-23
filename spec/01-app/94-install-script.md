# Install Scripts

## Overview

One-liner installer scripts that download, verify, and install the `gitmap-v27`
binary from GitHub Releases. Each script supports version pinning, checksum
verification, and automatic PATH registration.

---

## Repository

| Field       | Value                                              |
|-------------|----------------------------------------------------|
| GitHub Repo | `alimtvnetwork/gitmap-v27`                 |
| Binary Name | `gitmap-v27` (`gitmap.exe` on Windows)                 |
| Asset Format| `gitmap-{os}-{arch}.zip` (Windows), `gitmap-{os}-{arch}.tar.gz` (Unix) |
| Checksums   | `checksums.txt` (SHA-256, one line per asset)      |

---

## Windows — `install.ps1`

### One-Liner (Full Bootstrap)

The recommended one-liner follows the Chocolatey install pattern: it bypasses
the execution policy for the current process, enforces TLS 1.2+, and
downloads-then-executes the installer script. This ensures the command works
on locked-down machines, older Windows versions, and fresh installs where
`irm` may not be available.

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap-v27/scripts/install.ps1'))
```

### Short-Form (PowerShell 5+ / Modern Systems)

If the machine already has TLS 1.2 defaults and unrestricted execution
policy (e.g., developer workstations), the short form also works:

```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap-v27/scripts/install.ps1 | iex
```

### Why the Full Bootstrap?

| Concern                  | `irm \| iex`      | Full bootstrap           |
|--------------------------|--------------------|--------------------------|
| Execution policy blocked | Fails              | Bypasses (process scope) |
| TLS 1.2 not default      | May fail on old OS | Forces TLS 1.2+         |
| PowerShell 3.x compat   | No (`irm` = PS3+) | Yes (`WebClient` = PS2+) |
| Corporate firewalls      | May fail silently  | Explicit protocol set    |

### Parameters

| Parameter    | Type   | Default                        | Description                        |
|--------------|--------|--------------------------------|------------------------------------|
| `Version`    | string | latest (via GitHub API)        | Pin a specific release tag         |
| `InstallDir` | string | `$env:LOCALAPPDATA\gitmap-v27`     | Target directory for the binary    |
| `Arch`       | string | auto-detect                    | Force `amd64` or `arm64`           |
| `NoPath`     | switch | false                          | Skip adding install dir to PATH    |

### Flow

1. Resolve version — fetch latest tag from GitHub API or use pinned value.
2. Resolve architecture — read `PROCESSOR_ARCHITECTURE` or use override.
3. Download `gitmap-windows-{arch}.zip` and `checksums.txt`.
4. Verify SHA-256 checksum against `checksums.txt`.
5. Extract zip to install directory (rename-first if binary is running).
6. Add install directory to user PATH (unless `--NoPath`).
7. Print an install summary with installed version, binary path, install directory,
   and PATH target/status.

### File

`gitmap-v27/scripts/install.ps1`

---

## Linux / macOS — `install.sh`

### One-Liner

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap-v27/scripts/install.sh | bash
```

### Version-Pinned

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap-v27/scripts/install.sh | bash -s -- --version v2.55.0
```

### Custom Directory

```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap-v27/scripts/install.sh | bash -s -- --dir /opt/gitmap-v27 --version v2.55.0
```

### Parameters (CLI Flags)

| Flag        | Default              | Description                        |
|-------------|----------------------|------------------------------------|
| `--version` | latest (GitHub API)  | Pin a specific release tag         |
| `--dir`     | `~/.local/bin`       | Target directory for the binary    |
| `--arch`    | auto-detect          | Force `amd64` or `arm64`           |
| `--no-path` | false                | Skip adding install dir to PATH    |

### Flow

1. Detect OS (`linux` or `darwin`); reject Windows with redirect to `install.ps1`.
2. Detect architecture (`uname -m` → `amd64` or `arm64`).
3. Resolve version — query GitHub API via `curl` or `wget`, or use `--version`.
4. Download `gitmap-{version}-{os}-{arch}.tar.gz` and `checksums.txt`.
5. Verify SHA-256 checksum (`sha256sum` or `shasum -a 256`).
6. If `.tar.gz` not found in checksums, fall back to `.zip` variant.
7. Extract archive to temp directory; search for binary using 4-priority
   matching: exact name → platform-specific → versioned pattern
   (e.g., `gitmap-v27.55.0-linux-amd64`) → any executable.
8. Rename-first strategy for safe upgrades of running binaries.
9. Set executable permission (`chmod +x`).
10. Auto-detect shell (bash/zsh/fish) and append PATH entry to the
    correct profile file (`~/.bashrc`, `~/.zshrc`, or `~/.config/fish/config.fish`).
11. Print an install summary with installed version, binary path, install directory,
    detected shell, PATH target/status, and reload guidance.

### File

`gitmap-v27/scripts/install.sh`

---

## Checksum Verification

Both scripts download `checksums.txt` from the same release. Each line
follows the format:

```
<sha256-hash>  <filename>
```

The script matches the downloaded asset filename, compares hashes, and
aborts with a clear error on mismatch.

---

## Architecture Detection

| Platform | Source                          | Mapping                          |
|----------|---------------------------------|----------------------------------|
| Windows  | `$env:PROCESSOR_ARCHITECTURE`   | `AMD64`/`x86` → `amd64`, `ARM64` → `arm64` |
| Linux    | `uname -m`                      | `x86_64` → `amd64`, `aarch64` → `arm64`    |
| macOS    | `uname -m`                      | `x86_64` → `amd64`, `arm64` → `arm64`      |

---

## PATH Registration

| Platform | Method                                          |
|----------|-------------------------------------------------|
| Windows  | `[Environment]::SetEnvironmentVariable` (User) + `SendMessageTimeout` broadcast |
| Linux    | Auto-appends `export PATH` to `~/.bashrc` or `~/.profile`                        |
| macOS    | Auto-appends `export PATH` to `~/.zshrc` or `~/.bash_profile`                    |
| Fish     | Auto-appends `fish_add_path` to `~/.config/fish/config.fish`                     |

Windows modifies the registry-backed user PATH immediately and broadcasts
the change via `WM_SETTINGCHANGE`. Unix scripts auto-detect the active
shell and append a PATH entry to the appropriate profile file. The
`--no-path` / `-NoPath` flag skips this step on both platforms.

The post-install summary must always show:
- installed version
- binary path
- install directory
- PATH target and whether it was added, already present, or skipped
- shell/profile guidance where applicable

---

## Constraints

- No external dependencies beyond `curl`/`PowerShell` and `tar`/`Expand-Archive`.
- Scripts exit non-zero on any failure (download, checksum, extract).
- No interactive prompts — fully automatable.
- Temp files cleaned up in all exit paths.

---

## Related

- [CLI Interface](02-cli-interface.md)
- [Install Bootstrap](83-install-bootstrap.md)
- [Build & Deploy](09-build-deploy.md)
- [Future Features](82-future-features.md)
- [Release Workflow](../../.github/workflows/release.yml)
- [Release Workflow](../../.github/workflows/release.yml)

## Cross-References (Generic Specifications)

| Topic | Generic Spec | Covers |
|-------|-------------|--------|
| Install scripts | [03-install-scripts.md](../07-generic-release/03-install-scripts.md) | Version-pinned installers, SHA-256 verification, PATH registration |
| Release pipeline | [02-release-pipeline.md](../07-generic-release/02-release-pipeline.md) | Script generation via placeholder substitution |
| Checksums | [04-checksums-verification.md](../07-generic-release/04-checksums-verification.md) | SHA-256 generation and verification |
