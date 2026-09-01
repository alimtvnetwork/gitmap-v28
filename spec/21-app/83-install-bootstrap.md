# Install Bootstrap One-Liner

## Overview

The project provides a Chocolatey-style one-liner that safely bootstraps the
`gitmap-v28` installer on any Windows machine — including locked-down
environments, older OS versions, and fresh installs.

---

## The Pattern

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('<script-url>'))
```

This is the same pattern used by:
- **Chocolatey** (`community.chocolatey.org/install.ps1`)
- **Scoop** (`get.scoop.sh`)
- **winget bootstrap** scripts

---

## Breakdown

| Segment | Purpose |
|---------|---------|
| `Set-ExecutionPolicy Bypass -Scope Process -Force` | Allows script execution for the current process only — no permanent system change |
| `[SecurityProtocol] -bor 3072` | Bitwise-OR adds TLS 1.2 (0xC00) to the allowed protocols, required for GitHub HTTPS on older .NET |
| `New-Object System.Net.WebClient` | Uses the .NET WebClient class available since PowerShell 2.0 (Windows 7+) |
| `.DownloadString(url)` | Fetches the script as a string in memory |
| `iex (...)` | Invoke-Expression executes the downloaded script in the current session |

---

## gitmap-v28 One-Liner

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap-v28/scripts/install.ps1'))
```

### With Version Pinning

```powershell
$env:GITMAP_VERSION = "v2.25.0"; Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap-v28/scripts/install.ps1'))
```

---

## Compatibility Matrix

| Environment | `irm \| iex` | Full Bootstrap |
|-------------|:------------:|:--------------:|
| PowerShell 2.x (Win 7) | No | Yes |
| PowerShell 3.x+ | Yes | Yes |
| PowerShell 5.x (Win 10) | Yes | Yes |
| PowerShell 7.x (Core) | Yes | Yes |
| Restricted execution policy | No | Yes |
| AllSigned execution policy | No | Yes |
| TLS 1.0-only default (.NET < 4.6) | No | Yes |
| Corporate proxy (basic) | Partial | Yes |

---

## Security Considerations

- `Set-ExecutionPolicy Bypass -Scope Process` affects **only** the current
  PowerShell process — it does not modify machine or user policy.
- The `-bor 3072` is additive — it does not remove existing protocols.
- The downloaded script is executed in memory and is subject to the same
  checksum verification documented in `94-install-script.md`.

---

## Related

- [Install Scripts](94-install-script.md) — full installer spec
- [Build & Deploy](09-build-deploy.md)
