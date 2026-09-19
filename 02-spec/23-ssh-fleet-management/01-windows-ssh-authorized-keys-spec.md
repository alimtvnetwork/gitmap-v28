# Windows OpenSSH Authorized Keys Architecture & Injection Specification

**Version:** 1.0.0  
**Updated:** 2026-09-20  
**AI Confidence:** Production-Ready  
**Ambiguity:** None  

---

## Keywords

`windows-ssh` · `openssh` · `authorized-keys` · `fleet-management` · `administrators-authorized-keys` · `icacls` · `service-lifecycle`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| Target File Present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

This specification defines the architectural rules, security models, file locations, privilege detection mechanics, and idempotent key injection algorithms for deploying OpenSSH public keys onto Windows nodes in the GitMap SSH fleet management subsystem.

---

## 1. Dual Authorized Keys Architecture on Windows

Unlike POSIX systems (Linux, macOS) where all public keys reside within `~/.ssh/authorized_keys`, Microsoft Win32-OpenSSH enforces a strict separation between standard users and accounts with administrative privileges.

By default, the Windows OpenSSH daemon configuration (`$env:ProgramData\ssh\sshd_config`) contains:

```sshd
Match Group administrators
       AuthorizedKeysFile __PROGRAMDATA__/ssh/administrators_authorized_keys
```

This directive overrides the standard user profile authorized keys file whenever the connecting user belongs to the local `Administrators` group.

### Key Target Mapping

| Connecting User Role | Target File Path | Strict Security / ACL Requirement |
| :--- | :--- | :--- |
| **Local / Domain Administrator** | `$env:ProgramData\ssh\administrators_authorized_keys` | Inheritance removed (`/inheritance:r`). Full control granted solely to `Administrators` and `SYSTEM` (`/grant "Administrators:F"` `/grant "SYSTEM:F"`). |
| **Standard User** | `$env:USERPROFILE\.ssh\authorized_keys` | Standard profile access permissions. Directory `.ssh` created if missing. |

---

## 2. Administrator & Privilege Detection

When injecting public keys over an existing SSH session, the active session may operate with filtered tokens under User Account Control (UAC). Therefore, elevation detection must verify both token elevation and group membership.

```powershell
# 1. Primary Check: Token Elevation
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

# 2. Secondary Check: Localgroup Administrators Membership
if (-not $isAdmin) {
    $members = net localgroup administrators 2>$null
    if ($members -match $env:USERNAME) {
        $isAdmin = $true
    }
}
```

If `$isAdmin` is true:
- The target file is `$env:ProgramData\ssh\administrators_authorized_keys`.
- Directory `$env:ProgramData\ssh` is created if missing.
- ACL enforcement via `icacls` is strictly executed.

If `$isAdmin` is false:
- The target file is `$env:USERPROFILE\.ssh\authorized_keys`.
- Directory `$env:USERPROFILE\.ssh` is created if missing.

---

## 3. Mandatory Access Control List (ACL) Rules

OpenSSH on Windows strictly validates file ownership and access control lists before reading `administrators_authorized_keys`. If any non-administrative user, authenticated users group, or inherited world-readable rule is attached, `sshd.exe` rejects authentication with error:

```
Authentication refused: bad ownership or modes for file .../administrators_authorized_keys
```

### Remediation via `icacls`

The deployment script executes `icacls` with inheritance broken and explicit grants:

```cmd
icacls $keysFile /inheritance:r /grant "Administrators:F" /grant "SYSTEM:F"
```

1. `/inheritance:r` — strips all inherited permissions from parent directories (`C:\ProgramData\ssh` or `C:\ProgramData`).
2. `/grant "Administrators:F"` — grants full access rights to members of the local Administrators group.
3. `/grant "SYSTEM:F"` — grants full access rights to the NT AUTHORITY\SYSTEM account required for `sshd` worker process impersonation.

---

## 4. Key Format Validation & Token-Based Deduplication

SSH public keys adhere to the standard format:
`<algorithm-type> <base64-key-data> [optional-comment]`

Example:
`ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGitMapTestKey user@workstation`

### Deduplication Algorithm

Simple string substring matching (`$lines -notcontains $key`) fails when:
1. Trailing or leading whitespace varies.
2. Comments differ between machines while public key payloads are identical.

To guarantee zero duplicate entries:
1. Parse the incoming key by splitting on whitespace tokens (`\s+`).
2. Require at least two tokens (`$tokens.Length -ge 2`); reject malformed keys.
3. Extract Token 1 (`$tokens[1]`), which constitutes the unique base64 payload.
4. Scan each existing line in the target keys file:
   - Split existing line into whitespace tokens.
   - If Token 1 matches the incoming key payload, mark key as already present.
5. If not present, append the incoming key via `Add-Content`.

---

## 5. SSH Daemon (`sshd`) Service Health Check & Activation

Deploying authorized keys is useless if the `sshd` Windows service is stopped, disabled, or set to manual startup.

The injection script checks the Windows service controller:

```powershell
$svc = Get-Service sshd -ErrorAction SilentlyContinue
if ($svc) {
    if ($svc.StartType -ne 'Automatic') {
        Set-Service sshd -StartupType Automatic -ErrorAction SilentlyContinue
    }
    if ($svc.Status -ne 'Running') {
        Start-Service sshd -ErrorAction SilentlyContinue
    }
}
```

This ensures:
1. Service starts automatically on subsequent system reboots (`-StartupType Automatic`).
2. Service is immediately running for inbound SSH connections (`Start-Service sshd`).

---

## 6. Canonical Self-Contained Deployment One-Liner

The complete command executed by `buildWindowsAuthKeyScript` in `cli/cmdssh/ssh_auth_key_deploy.go` runs with zero external script dependencies:

```powershell
powershell -NoProfile -Command "$k = '%s'.Trim(); $tokens = $k -split '\s+'; if ($tokens.Length -lt 2) { exit 1 }; $keyBody = $tokens[1]; $isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator); if (-not $isAdmin) { $members = net localgroup administrators 2>$null; if ($members -match $env:USERNAME) { $isAdmin = $true } }; if ($isAdmin) { $keysFile = Join-Path $env:ProgramData 'ssh\administrators_authorized_keys'; $sshDir = Join-Path $env:ProgramData 'ssh'; if (!(Test-Path $sshDir)) { New-Item -ItemType Directory -Path $sshDir -Force | Out-Null }; if (!(Test-Path $keysFile)) { New-Item -ItemType File -Path $keysFile -Force | Out-Null }; icacls $keysFile /inheritance:r /grant 'Administrators:F' /grant 'SYSTEM:F' | Out-Null } else { $sshDir = Join-Path $env:USERPROFILE '.ssh'; if (!(Test-Path $sshDir)) { New-Item -ItemType Directory -Path $sshDir -Force | Out-Null }; $keysFile = Join-Path $sshDir 'authorized_keys'; if (!(Test-Path $keysFile)) { New-Item -ItemType File -Path $keysFile -Force | Out-Null } }; $lines = Get-Content $keysFile -ErrorAction SilentlyContinue; $hasKey = $false; foreach ($line in $lines) { $t = $line.Trim() -split '\s+'; if ($t.Length -ge 2 -and $t[1] -eq $keyBody) { $hasKey = $true; break } }; if (-not $hasKey) { Add-Content -Path $keysFile -Value $k }; $svc = Get-Service sshd -ErrorAction SilentlyContinue; if ($svc) { if ($svc.StartType -ne 'Automatic') { Set-Service sshd -StartupType Automatic -ErrorAction SilentlyContinue }; if ($svc.Status -ne 'Running') { Start-Service sshd -ErrorAction SilentlyContinue } }"
```

---

## Cross-References

- [SSH Keys Specification](../21-app/50-ssh-keys.md)
- [Coding Guidelines Golang](../02-coding-guidelines/03-golang/00-overview.md)
- [Error Management Architecture](../03-error-manage/02-error-architecture/00-overview.md)
