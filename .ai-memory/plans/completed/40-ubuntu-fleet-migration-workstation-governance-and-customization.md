# Milestone 40: Ubuntu Fleet Migration, Workstation Governance, and OS Customization

- **Slug:** `ubuntu-fleet-migration-workstation-governance-and-customization`
- **Milestone Index:** `40`
- **Status:** `COMPLETED`
- **Source Plans Merged:** Plans 52, 76, 81, 206, 211, 212, 214, 215, 222
- **Folded Subtask Folders:**
  - `212-ubuntu-fleet-full-customization-and-embedded-runner` (5 files)
  - `214-ubuntu-fleet-automation-and-workstation-governance` (5 files)
- **Target Subsystems:** `cli/cmdos/`, `cli/cmdfleet/`, `cli/cmdrunner/`, `cli/cmdtweaks/`

---

## 1. Domain Context & Architectural Problem

As development infrastructure transitioned across the fleet from Windows workstations to Ubuntu Linux nodes (e.g. `u1`, developer virtual machines, and cloud instances), significant cross-platform discrepancies surfaced:
1. **Desktop Ergonomics & High-DPI Scaling:** Headless and graphical Ubuntu installations suffered from misaligned desktop scaling on 4K displays, missing GNOME custom shortcuts, and inconsistent window manager settings.
2. **Virtualization Shared Folders:** VMware shared folder mounts (`/mnt/hgfs`) frequently failed to automount on boot or lost permissions across user sessions, preventing access to host workspace repositories.
3. **Cross-OS Shell & Tooling Parity:** PowerShell scripts (`run.ps1`) lacked equivalent native bash runner behavior (`run.sh`), and Oh-My-Zsh plugins generated warning output when executed under non-interactive SSH sessions.
4. **Filesystem Case Sensitivity & Renaming:** Case-insensitive Windows filesystems caused Git to miss case-only renames (e.g. `readme.md` to `README.md`), requiring a dedicated multi-OS renamer engine.

---

## 2. Synthesized Architectural Outcomes

### 2.1 Ubuntu Fleet Migration & Automated Provisioning
- **Zero-to-End Node Setup:** Authored automated provisioning pipelines in `cli/cmdfleet/ubuntu_provision.go` enabling clean onboarding of Ubuntu machines (`gitmap fleet setup ubuntu`):
  - Installs Git, Go toolchain, Node.js/pnpm, Python3, and build essentials.
  - Configures SSH server with key-based authentication and hardened `sshd_config`.
  - Sets up GitMap binary in `/usr/local/bin/gitmap` with bash and zsh tab autocompletion.
- **Cross-OS Shell Normalization:** Added embedded runner engine in `cli/cmdrunner/` providing feature-identical execution across `run.ps1` (PowerShell) and `run.sh` (POSIX bash), driven by a single `run.config.json` specification.

### 2.2 GNOME Desktop Governance & Display Scaling
- **High-DPI Display Scaling:** Added `gitmap os display scale 140` (and `gitmap os dm`) to configure GNOME desktop fractional scaling via `gsettings` (`org.gnome.mutter experimental-features ['scale-monitor-framebuffer']`), setting 140% scaling for optimal typography readability on high-resolution displays.
- **Desktop Ergonomics:** Configures default dark theme, custom terminal keyboard shortcuts, dock positioning (left-aligned, autohide enabled), and favorite application pins via declarative CLI commands.

### 2.3 VMware Shared Folder Automount & Persistence
- **Automount Daemon Unit:** Generated persistent systemd mount units (`open-vm-tools` integration) to automatically mount VMware shared directories to `/mnt/hgfs` during early boot targets:
  - Enforces `uid=1000,gid=1000,allow_other` mount options so regular non-root users retain full read/write permissions.
  - Adds liveness check in `gitmap os vmware status` verifying active shared mount paths before running workspace tasks.

### 2.4 Multi-OS Case Renamer Engine (`gitmap lower-case-fix`)
- **Git Case Preservation:** Implemented `cli/cmdtweaks/lowercase_fix.go` to safely resolve case conflicts between Git index and filesystem:
  - Uses two-phase atomic Git renames (`git mv foo foo_tmp && git mv foo_tmp bar`) to bypass Windows NTFS case-insensitivity.
  - Audits all files across repositories against lowercase naming guidelines, reporting and auto-fixing violations.

### 2.5 Oh-My-Zsh Integration & Non-Interactive SSH Cleanliness
- **Shell Profiling Guards:** Cleaned remote execution wrappers in `cli/cmdssh/` to bypass interactive Oh-My-Zsh themes during non-interactive SSH commands, preventing ANSI escape sequences or theme update notifications from corrupting structured JSON stdout.

---

## 3. Go Type Contracts & Architecture

```go
// UbuntuNodeConfig models target machine provisioning options
type UbuntuNodeConfig struct {
    HostName       string   `json:"hostName"`
    IPAddress      string   `json:"ipAddress"`
    DisplayScale   int      `json:"displayScale"` // e.g. 100, 125, 140, 150
    IsVMwareGuest  bool     `json:"isVmwareGuest"`
    IsDesktopUI    bool     `json:"isDesktopUi"`
    InstalledTools []string `json:"installedTools"`
}

// DisplayManagerSettings captures desktop environment state
type DisplayManagerSettings struct {
    SessionType    string  `json:"sessionType"` // "wayland", "x11"
    ScaleFactor    float64 `json:"scaleFactor"`
    ThemeName      string  `json:"themeName"`
    IsFractionalOK bool    `json:"isFractionalOk"`
}

// RenamerAuditResult tracks file casing normalization
type RenamerAuditResult struct {
    ScannedFiles    int      `json:"scannedFiles"`
    RenamedFiles    int      `json:"renamedFiles"`
    UnchangedFiles  int      `json:"unchangedFiles"`
    ModifiedPaths   []string `json:"modifiedPaths"`
}
```

OS-specific system calls are isolated behind OS build tags (`//go:build linux`, `//go:build windows`) to maintain clean cross-compilation.

---

## 4. Subtask Verification Ledger

| Folded Subtask Directory | Source Subtask Files | Verified Criteria & Delivered Artifacts |
| :--- | :--- | :--- |
| `212-ubuntu-fleet-full-customization-and-embedded-runner` | 5 subtask files | Ubuntu desktop provisioning, embedded `run.sh` cross-platform runner, GNOME dark theme. |
| `214-ubuntu-fleet-automation-and-workstation-governance` | 5 subtask files | Workstation governance commands (`gitmap os`), VMware `/mnt/hgfs` automount persistence, systemd service units. |

---

## 5. Quality & Coding Guideline Compliance

- **Positive Booleans Only:** Enforced `isVmwareGuest`, `isDesktopUi`, `isFractionalOk`, `isToolsInstalled`.
- **System Isolation:** Root privileges (`sudo`) are requested only for system package installations; all user-level configurations remain scoped to `$HOME`.
- **Relative Path Adherence:** All documentation, logs, and audit ledgers use strictly relative Git repository paths.
- **Linter Gates:** Verified with Go static analysis, zero compiler diagnostics, and compliant LF line endings.
