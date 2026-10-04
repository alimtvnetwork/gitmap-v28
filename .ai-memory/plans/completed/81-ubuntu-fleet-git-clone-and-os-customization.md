# Plan 81: Ubuntu Fleet Git Clone, Desktop Ergonomics & Cross-OS Migration (Completed)

> **Plan Status:** Completed  
> **Completed At:** 2026-10-04  
> **Traceability IDs:** Subtask 211-01 .. Subtask 211-05  
> **Spec Reference:** [02-spec/21-app/211-ubuntu-fleet-git-clone-and-os-customization/01-architecture-spec.md](../../../02-spec/21-app/211-ubuntu-fleet-git-clone-and-os-customization/01-architecture-spec.md)  
> **Target Node:** Ubuntu U1 (`node-u1`)  
> **Deliverable Path:** `./repo-secrets\04-ubuntu-migration\`  

---

## 1. Architectural Summary & Results

All 5 core objectives have been executed, tested live on Ubuntu U1, and verified:
1. **SSH Fleet & Git Clone Pipeline**:
   - Registered `Host u1` in `%USERPROFILE%\.ssh\config` pointing to `id_rsa.backup-devorg` and `node-u1`.
   - Sanitized all 71 repositories from `gitmap.json` into Unix format `gitmap-linux.json` with forward slashes `/`.
   - Executed `clone-repos-to-u1.sh` and `clone-repos-to-u1.ps1` over SSH, safely cloning the 45 missing repositories into `/home/a/git-work/` while skipping all existing repositories.
   - Total verified repositories in `/home/a/git-work`: **74** (all 71 manifest repositories + legacy variations).
2. **Ubuntu Desktop Visual Ergonomics & Keybindings**:
   - Executed `configure-ubuntu-desktop.sh` over SSH with user D-Bus bus session address.
   - Applied 140% font scaling: `org.gnome.desktop.interface text-scaling-factor 1.4` (verified `1.4` live).
   - Configured Windows muscle-memory shortcuts:
     - `Win+D`: Show Desktop (`show-desktop`)
     - `Win+Tab`: Task View Overview (`toggle-overview`)
     - `Ctrl+Win+Left` / `Ctrl+Win+Right`: Switch Workspaces (`switch-to-workspace-left/right`)
     - `Alt+Tab`: Ungrouped Individual Window Switching (`switch-windows`)
     - `Ctrl+Shift+D` / `Ctrl+Win+D`: Jump to New Workspace (`switch-to-workspace-last`)
3. **VMware Shared Folders Automount**:
   - Deployed systemd units `/etc/systemd/system/mnt-hgfs.mount` and `/etc/systemd/system/mnt-hgfs.automount` with `fuse.vmhgfs-fuse` and `allow_other,uid=1000,gid=1000`.
   - Verified `systemctl is-active mnt-hgfs.automount` is `active` and boot-persistent.
4. **Antigravity Cross-OS Settings & Path Normalization**:
   - Created root symlink `/d/work` -> `/home/a/git-work` on Ubuntu for zero-rewrite transcript compatibility.
   - Transferred and normalized Antigravity IDE configuration from `%APPDATA%\Antigravity\User\` to `~/.config/Antigravity/User/`.
5. **Master Orchestrator & Retrospective Log**:
   - Provided `master-ubuntu-setup.ps1` chaining all tasks together cleanly.
   - Authored comprehensive engineering log `step-by-step-log.md` detailing thinking, executed commands, outputs, and errors faced/resolved.
   - Authored `02-full-os-setup-blueprint.md` outlining the roadmap for future complete OS setup (VS Code, themes, flameshot, chrome, gitmap tools).

---

## 2. Completed Subtasks & Deliverables

| Subtask ID | Focus Area | Deliverables | Verification |
| :--- | :--- | :--- | :--- |
| **Subtask 211-01** | SSH & Git Fleet Sync | `clone-repos-to-u1.sh`, `clone-repos-to-u1.ps1`, `gitmap-linux.json` | PASS (74 repos in /home/a/git-work) |
| **Subtask 211-02** | Desktop Ergonomics & Shortcuts | `configure-ubuntu-desktop.sh` | PASS (text-scaling-factor 1.4, keybindings) |
| **Subtask 211-03** | VMware Shared Folders | `setup-vmware-shared-folders.sh` | PASS (mnt-hgfs.automount active) |
| **Subtask 211-04** | Antigravity Cross-OS Sync | `sync-antigravity-settings.sh`, `sync-antigravity-settings.ps1` | PASS (/d/work symlink, normalized URIs) |
| **Subtask 211-05** | Master Runner & Documentation | `master-ubuntu-setup.ps1`, `step-by-step-log.md`, blueprint | PASS (complete log and master runner) |

---

## 3. Acceptance Verification Evidence

```text
CONNECTED
Repo count: 74
Scaling factor: 1.3999999999999999
lrwxrwxrwx 1 root root 16 Oct  4 21:06 /d/work -> /home/a/git-work
active
Show desktop: ['<Super>d']
Switch windows: ['<Alt>Tab']
Switch workspace left: ['<Control><Super>Left']
Overview: ['<Super>Tab']
```
