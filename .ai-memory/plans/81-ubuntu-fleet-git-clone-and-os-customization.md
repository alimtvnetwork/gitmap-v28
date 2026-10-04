# Plan 81: Ubuntu Fleet Git Clone, Desktop Ergonomics & Cross-OS Migration

> **Plan Status:** Active  
> **Traceability IDs:** Subtask 211-01 .. Subtask 211-05  
> **Spec Reference:** [02-spec/21-app/211-ubuntu-fleet-git-clone-and-os-customization/01-architecture-spec.md](../../02-spec/21-app/211-ubuntu-fleet-git-clone-and-os-customization/01-architecture-spec.md)  
> **Target Node:** Ubuntu U1 (`192.168.1.22`)  
> **Execution Location:** `d:\work\repo-secrets\04-ubuntu-migration\`  

---

## 1. Architectural Context & Subtask Mapping

| Subtask ID | Focus Area | Target Deliverable | Status |
| :--- | :--- | :--- | :--- |
| **Subtask 211-01** | SSH & Git Fleet Sync | `clone-repos-to-u1.sh`, `clone-repos-to-u1.ps1`, `gitmap-linux.json` | In-Progress |
| **Subtask 211-02** | Desktop Ergonomics & Shortcuts | `configure-ubuntu-desktop.sh` (140% font scaling, Windows keybindings) | In-Progress |
| **Subtask 211-03** | VMware Shared Folders | `setup-vmware-shared-folders.sh` (`mnt-hgfs.mount`, `mnt-hgfs.automount`) | In-Progress |
| **Subtask 211-04** | Antigravity Cross-OS Sync | `sync-antigravity-settings.sh` (`/d/work` symlink, path normalization) | In-Progress |
| **Subtask 211-05** | Master Runner & Documentation | `master-ubuntu-setup.ps1`, `step-by-step-log.md`, full OS setup plan | In-Progress |

---

## 2. 5-Subtask Detailed Breakdown

### Subtask 211-01: SSH Fleet Remote Git Clone Pipeline
- Configure `Host u1` in `C:\Users\Administrator\.ssh\config`.
- Generate sanitized `gitmap-linux.json` mapping all 71 repositories to forward-slash Unix paths.
- Implement `clone-repos-to-u1.sh` and wrapper `clone-repos-to-u1.ps1` to clone all missing 45 repos into `/home/a/git-work/` over SSH.

### Subtask 211-02: Ubuntu 140% Font Scaling and Windows Keybindings
- Implement `configure-ubuntu-desktop.sh`:
  - Set `org.gnome.desktop.interface text-scaling-factor 1.4`.
  - Bind `Win+D` to show desktop (`show-desktop`).
  - Bind `Win+Tab` to toggle overview (`toggle-overview`).
  - Bind `Ctrl+Win+Left` and `Ctrl+Win+Right` to workspace switching.
  - Bind `Alt+Tab` to ungrouped window cycling (`switch-windows`).
  - Bind `Ctrl+Shift+D` / `Ctrl+Win+D` to new workspace jump.

### Subtask 211-03: VMware Shared Folder Systemd Automount at Boot
- Implement `setup-vmware-shared-folders.sh`:
  - Verify `open-vm-tools` and `open-vm-tools-desktop`.
  - Configure `user_allow_other` in `/etc/fuse.conf`.
  - Deploy systemd units `/etc/systemd/system/mnt-hgfs.mount` and `/etc/systemd/system/mnt-hgfs.automount`.
  - Enable and start units.

### Subtask 211-04: Cross-OS Antigravity Settings and Conversation Sync
- Implement `sync-antigravity-settings.sh`:
  - Create `/d/work` symlink on Ubuntu pointing to `/home/a/git-work`.
  - Export and transfer Antigravity settings and conversation transcripts from Windows to Ubuntu.
  - Normalize workspace URIs from `file:///d:/work/` to `file:///home/a/git-work/`.

### Subtask 211-05: Master Orchestrator Script and Detailed Markdown Log
- Implement `master-ubuntu-setup.ps1` to execute all steps sequentially with status reporting.
- Author comprehensive engineering log `d:\work\repo-secrets\04-ubuntu-migration\step-by-step-log.md`.
- Document future full OS setup blueprint (VS Code themes, flameshot, browser, gitmap tools).
