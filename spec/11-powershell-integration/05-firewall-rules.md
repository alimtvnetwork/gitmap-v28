# Windows Firewall Configuration

> **Version:** 1.1.0  
> **Updated:** 2026-02-04  
> **Status:** Active

---

## Overview

The PowerShell runner can automatically configure Windows Firewall inbound rules to allow network access to the Go backend server.

---

## Automatic Configuration

### Using the Script

```powershell

# Run as Administrator

.\run.ps1 -OpenFirewall
```

This creates inbound rules for all ports in `powershell.json`:

```json
{
  "ports": [8080, 8081]
}
```

### Created Rules

For each port, the script creates:

| Property | Value |
|----------|-------|
| DisplayName | `LLM Runner (Go Backend) TCP {port}` |
| Direction | Inbound |
| Action | Allow |
| Protocol | TCP |
| LocalPort | {port} |
| Profile | Private, Domain |

---

## Manual Configuration

If automatic configuration fails:

### PowerShell (Admin)

```powershell

# Create rule for port 8080

New-NetFirewallRule `
    -DisplayName "LLM Runner (Go Backend) TCP 8080" `
    -Direction Inbound `
    -Action Allow `
    -Protocol TCP `
    -LocalPort 8080 `
    -Profile Private,Domain

# Create rule for port 8081

New-NetFirewallRule `
    -DisplayName "LLM Runner (Go Backend) TCP 8081" `
    -Direction Inbound `
    -Action Allow `
    -Protocol TCP `
    -LocalPort 8081 `
    -Profile Private,Domain
```

### Windows Firewall GUI

1. Open **Windows Defender Firewall with Advanced Security**
2. Click **Inbound Rules** → **New Rule...**
3. Select **Port** → **Next**
4. Select **TCP**, enter port (e.g., `8080`) → **Next**
5. Select **Allow the connection** → **Next**
6. Check **Domain** and **Private** (uncheck Public) → **Next**
7. Name: `LLM Runner (Go Backend) TCP 8080` → **Finish**
