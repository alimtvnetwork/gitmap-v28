# Spec 158: SSH Installer Payload NSIS Auto-Detection, Silent Execution & Fleet Enrollment Resolution

- **Status:** Approved
- **Owner:** Core CLI & Fleet Engineering
- **Date:** 2026-09-25
- **Version:** v6.338.0
- **Tracking Issue:** Issue #44 (02-spec/22-app-issues/44-ssh-install-exec-nsis-silent-hanging-and-empty-registry-resolution-rca.md)

---

## 1. Context & Problem Statement

When deploying software installers across distributed SSH fleet nodes using `gitmap ssh install-exec <installer-file>`:
1. **Interactive Session 0 Hang on Windows Nodes:**
   Windows OpenSSH sessions run services and spawned commands non-interactively in Session 0 without an attached desktop or GUI interaction. When executable installers (such as `Antigravity.Manager.Tools_4.70.0_x64-setup.exe`, an NSIS installer) were deployed without explicit silent flags, GitMap previously defaulted `.exe` flags to Inno Setup parameters (`/VERYSILENT /SUPPRESSMSGBOXES /NORESTART /SP-`). Because NSIS does not recognize `/VERYSILENT`, it discarded the parameter and attempted to display an interactive GUI wizard in Session 0, causing the process to block indefinitely waiting for user clicks.
2. **Missing Host Registry Guidance in Ad-Hoc Directories:**
   When running `gitmap ssh install-exec` in directories without a populated local or global database (e.g. `Downloads`), the command returned `No matching SSH machines to deploy setup to (filtered by ...)` without explaining that 0 nodes were enrolled in the central database.
3. **Unmocked Network Dials in Unit Test Suites:**
   Tests executing SSH login flows previously made unmocked TCP connection attempts to non-existent IPs (`192.168.1.88`, `10.0.0.5`, `192.0.2.1`) during host key trust scans and OS profiling, which timed out after 10 minutes on CI/CD runner environments.

---

## 2. Architectural Design & Implementation

### 2.1 Payload-Aware Silent Installer Detection

`BuildRemoteInstallerExecCmdWithPayload` inspects the initial byte headers of the uploaded installer binary before execution:
- **NSIS (Nullsoft Scriptable Install System):** Detects magic signature `NullsoftInst`. Automatically supplies `/S` for unattended silent installation.
- **Inno Setup:** Detects `Inno Setup` signature. Automatically supplies `/VERYSILENT /SUPPRESSMSGBOXES /NORESTART /SP-`.
- **WiX / Burn Bootstrappers:** Detects `WixBundle` or `Burn` signatures. Supplies `/quiet /norestart`.
- **Microsoft Installer (.msi):** Automatically generates `$p = '<path>'; Start-Process msiexec.exe -ArgumentList @('/i', $p, '/qn') -Wait -PassThru`.
- **CLI Flag Precedence:** User-supplied explicit CLI arguments (e.g. `/S`, `/qn`, `/DIR=...`) take precedence over auto-detected defaults.

### 2.2 Default Silent Unattended Flag in Remote SSH Sessions

Remote SSH deployments default `opts.IsSilent = true` in `ParseInstallExecArgs`. Operators can opt out using `--no-silent`, `--gui`, or `--interactive`.

### 2.3 Registry Emptiness Diagnostic Guidance

When filtering returns 0 eligible hosts, `renderNoMatchingMachinesMessage` checks if total enrolled hosts is 0:
- If 0 registered hosts exist, outputs clear enrollment instructions:
  `Enroll machines first using: gitmap sjc <user> <ip(alias)...> --pass <password> or gitmap sj add <user@ip> [alias]`

### 2.4 Test Isolation & Network Dial Mocking

`cli/cmdssh` provides test hooks:
- `SetConnectProbeClientForTesting`: Mocks remote SSH client discovery during login commands.
- `SetAutoTrustTargetHostForTesting`: Mocks host key scanning and known_hosts resolution during tests.
- `SetCryptoConnectWithKeyForTesting`: Mocks default SSH key connections during key discovery tests.
- `ssh.ClientConfig.Timeout`: Pinned to `4 * time.Second` in `crypto.ConnectWithPassword` and `crypto.ConnectWithKey` to guarantee fast failover on offline nodes.

---

## 3. Verification & Compliance Matrix

| Target Component | Test Suite / Command | Result | Duration |
|------------------|----------------------|--------|----------|
| NSIS Auto-Detection & Cmd Gen | `TestBuildRemoteInstallerExecCmd_PayloadDetection` | PASS | 0.00s |
| Default Silent Flag Parsing | `TestParseInstallExecArgs_IsSilentDefault` | PASS | 0.00s |
| SSH Login Command Mocking | `TestRunSSHLogin_*` | PASS | 0.16s |
| Known Hosts & Key Discovery | `TestConnectWithDefaultKey_Unreachable` | PASS | 0.00s |
| Live Temporary E2E Fleet Deploy | `TestTempE2E_SSHInstallExecStreamingAndOSProbing` | PASS | 1.86s |
| Live NSIS Setup Deployment | `gitmap ssh install-exec Antigravity.Manager.Tools...` | PASS (Nodes w1, w2, w3) | 2.7s - 3.8s |
| Boolean & Enum Guidelines | `python linter-scripts/check-enum-and-boolean.py` | PASS (0 violations) | 1.75s |
| Nested If Linter | `python linter-scripts/check-nested-ifs.py` | PASS (0 violations) | 0.01s |
