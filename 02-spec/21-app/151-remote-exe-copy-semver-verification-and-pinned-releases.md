# Spec 151: Remote Binary Installation & Copy via Low-Level SSH, Remote Fleet SemVer Verification & Comparison, and Specific Version Pinning/Downgrades for GitMap & Antigravity Manager

> **Spec ID:** `SPEC-151`  
> **Version:** `v6.327.0`  
> **Status:** Implemented & Verified  
> **Date:** 2026-09-24  

---

## 0. User Request (Verbatim) & Actionable Requirements

```text
Can you install or copy exe file from one machine to another using SSH and lower commands? Confirm that. And also, how do you verify and confirm that the other machine has the Git map and which is running on the same version or above version? How do you know that? That's one interesting question. So you need to show me how you do that, that's the first understanding, and then we can discuss furthermore. Okay? And also, for all the updates or install, we can give a specific version that I want you to understand how to do that. And also, at the same time, let's say we have the AGM. So I could do AGM version LS. That would show me all the GitHub version tags. I could pick from this and install. I could do AGM install specific version, and that would install that version. You need to have that capability. Same thing should go for Git map update as well. I could go to Git map's older version as well if I choose to. Okay? Can you integrate this type of thing? We can confirm and then we can discuss furthermore.
```

### Actionable Deliverables:
1. **Low-Level SSH Exe Copy & Install Confirmation**:
   - Provide architectural explanation and native command execution (`gitmap ssh cp <src> <node>:<dest>`) showing how binary executables are transferred using pure SSH channel primitives (Base64 encoding, piped chunking, and memory-to-disk write) without SMB, CIFS, NFS, or FTP services.
2. **Remote Fleet Version Verification & SemVer Comparison**:
   - Provide real-time remote inspection (`gitmap ssh nodes --version` / `gitmap ssh node-version`) that probes remote hosts and compares remote installed GitMap versions against local running version (`constants.Version`):
     - `● IDENTICAL (SAME)`: Remote version equals local version.
     - `▲ ABOVE (NEWER)`: Remote version is higher than local version.
     - `▼ BELOW (OLDER)`: Remote version is lower than local version.
     - `○ NOT INSTALLED` / `○ OFFLINE`: Binary not found on PATH or machine unreachable.
3. **Antigravity Manager (AGM) Release Listing & Version Pinning**:
   - `gitmap agm version ls` (aliases: `agm versions`, `agm tags`, `agm ls`): Queries GitHub release tags from `alimtvnetwork/Antigravity-Manager` and displays release dates, tags, latest flag, and status in a formatted table.
   - `gitmap agm install <version>` (e.g. `gitmap agm install 4.7.1` or `--version 4.7.1`): Installs the specific pinned AGM version.
   - `gitmap agm update <version>` (e.g. `gitmap agm update 4.7.1` or `--version 4.7.1`): Updates or downgrades AGM to that specific version.
4. **GitMap Release Listing & Version Pinning / Downgrades**:
   - `gitmap version ls` (aliases: `gitmap versions`, `gitmap update ls`, `gitmap update versions`): Queries GitHub release tags for `gitmap-v28` and displays active, newer, and older versions.
   - `gitmap update <version>` (e.g. `gitmap update v6.324.0` or `--version v6.324.0`): Pins the target version, downloads the canonical remote installer (`install.ps1` or `install.sh`), and passes `-Version <ver>` / `--version <ver>` to install or downgrade to any chosen historical release.
5. **Temporary E2E Validation**:
   - Isolated under `cli/tests/e2e/version_ls_and_node_comparison_tempe2e_test.go` guarded by `//go:build tempe2e` and `RUN_TEMP_E2E=1`.

---

## 1. Technical Architecture & Protocols

### 1.1 Low-Level Exe Copy via Pure SSH (`gitmap ssh cp`)

GitMap does not require SMB, CIFS shares, FTP, or open incoming firewall ports. Binary files (`.exe`, ELF binaries, archives) are transferred directly through the encrypted SSH session:

1. **Local In-Memory Ingestion**: The local executable binary is read into an in-memory byte buffer and Base64 encoded.
2. **Remote In-Memory Materialization**:
   - **Windows Node**:
     ```powershell
     powershell -NoProfile -Command "$d=[Convert]::FromBase64String('<base64_data>'); $p='<remote_path>'; $dir=[IO.Path]::GetDirectoryName($p); if ($dir -and -not (Test-Path $dir)) { [IO.Directory]::CreateDirectory($dir) | Out-Null }; [IO.File]::WriteAllBytes($p, $d)"
     ```
   - **Linux / macOS Node**:
     ```bash
     mkdir -p "$(dirname '<remote_path>')" && printf '%s' '<base64_data>' | base64 -d > '<remote_path>' && chmod +x '<remote_path>'
     ```
3. **Self-Contained Verification**: The remote command executes atomically over SSH `crypto.RunCommand()`. Upon completion, standard exit codes confirm successful materialization and permission assignment (`chmod +x`).

### 1.2 Remote SemVer Detection & Comparison Engine

When executing `gitmap ssh nodes --version`:
1. GitMap dispatches probe commands concurrently across all registered SSH fleet machines:
   - Windows: `where.exe gitmap >nul 2>nul && gitmap.exe --version || ...`
   - Unix/macOS: `which gitmap >/dev/null 2>&1 && gitmap --version || ...`
2. GitMap parses the remote stdout string, extracting clean SemVer tokens (e.g. `v6.325.0` -> `6.325.0`).
3. `CompareSemverStrings(remoteVersion, localVersion)` performs component-wise integer comparison on `Major`, `Minor`, and `Patch`:
   - Returns `0`: Remote is identical to local. Output: `● IDENTICAL (SAME)`.
   - Returns `1`: Remote is newer than local. Output: `▲ ABOVE (NEWER)`.
   - Returns `-1`: Remote is older than local. Output: `▼ BELOW (OLDER)`.
   - If not installed: Output: `○ NOT INSTALLED`.

### 1.3 Specific Release Pinning & Downgrade Engine

Both GitMap and AGM support explicit target version arguments and `--version` flags:

- **GitMap**:
  - `gitmap update <version>` extracts the target semver and assigns `cmdupdate.SetTargetVersion(ver)`.
  - The update workflow fetches the installer script from GitHub:
    - On Windows: Invokes `install.ps1 -Version <version> -InstallDir <dir>`
    - On Unix/macOS: Invokes `install.sh --version <version> --dir <dir>`
  - Both installers download the exact tagged release asset from GitHub releases, extract the binary, and overwrite the local install directory.
- **Antigravity Manager (AGM)**:
  - `gitmap agm install <version>` and `gitmap agm update <version>` pass `-Version '<clean_version>'` to `install.ps1` (Windows) or `--version '<clean_version>'` to `install.sh` (Unix/macOS).

---

## 2. Command Reference

| Command | Description |
| :--- | :--- |
| `gitmap ssh cp <local.exe> <node>:<path>` | Copies local `.exe` or binary to remote machine via pure SSH stream. |
| `gitmap ssh nodes --version` | Scans fleet and prints version table with `COMPARISON (LOCAL: vX.Y.Z)`. |
| `gitmap version ls` / `gitmap versions` | Lists available GitMap GitHub releases with date, tag, and status. |
| `gitmap update <version>` | Installs or downgrades GitMap to the specified version. |
| `gitmap agm version ls` / `gitmap agm versions` | Lists available Antigravity Manager GitHub releases. |
| `gitmap agm install <version>` | Installs specific pinned AGM release version. |
| `gitmap agm update <version>` | Updates or downgrades AGM to specific pinned release version. |

---

## 3. Verification & Quality Gates

- **Unit Tests**:
  - `cli/cmdupdate`: `TestBuildRemoteInstallerCmd` (Windows & Unix arguments).
  - `cli/cmdssh`: `TestBuildRemoteWriteCmd` (Base64 decode and permissions).
- **Temporary E2E Suite**:
  - `cli/tests/e2e/version_ls_and_node_comparison_tempe2e_test.go`:
    - `TestTempE2E_SemverComparisonAndNodeStatus`: Verifies comparison states (`IDENTICAL`, `ABOVE`, `BELOW`, `NOT INSTALLED`).
    - `TestTempE2E_ReleaseTagsTableRendering`: Verifies table output rendering.
    - `TestTempE2E_TargetPinnedVersionAndInstallerArgs`: Verifies pinned version propagation.
    - `TestTempE2E_LowLevelSSHBypassExeCopyMechanism`: Verifies Base64 stream generation and `chmod +x`.
    - `TestTempE2E_AGMVersionListingAndInstallOptions`: Verifies AGM version dispatch.
