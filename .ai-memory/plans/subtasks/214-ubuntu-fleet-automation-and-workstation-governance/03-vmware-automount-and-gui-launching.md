# Subtask 03: VMware Automount & Remote GUI Launching

> **Task Reference:** `214-ubuntu-fleet-automation-and-workstation-governance`  
> **Parent Plan:** [.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md](./.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md)  
> **Target Node:** Ubuntu 24.04 LTS (`u1` / `ubuntu-fleet-01`)  

---

## 1. Objective & Requirements
- Fix the issue where VMware shared folders (`/mnt/hgfs`) do not automatically mount on system boot.
- Enable remote launching of graphical applications (like Antigravity IDE) from headless SSH sessions into the active user session.

---

## 2. Implementation & Commands
- Implemented systemd units `/etc/systemd/system/mnt-hgfs.mount` and `/etc/systemd/system/mnt-hgfs.automount` with `TimeoutIdleSec=0`.
- Configured `/etc/fuse.conf` with `user_allow_other` and mounted with `allow_other,uid=1000,gid=1000`.
- Remote GUI invocation implemented via `systemd-run --user antigravity`.
- Execution command:
  ```powershell
  pwsh -File $SECRETS_DIR/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action vmware
  pwsh -File $SECRETS_DIR/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action gui -GuiCommand "antigravity"
  ```

---

## 3. Verification & Live Evidence
- Live `systemctl is-active mnt-hgfs.automount` returns `active`.
- GUI application launches successfully via `systemd-run --user`.
- Status: **DONE**
