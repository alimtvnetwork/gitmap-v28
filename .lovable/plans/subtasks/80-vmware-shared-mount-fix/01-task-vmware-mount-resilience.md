# 01-task-vmware-mount-resilience.md: VMware Mount Engine Resilience & Error -107 Diagnostics

**Parent Plan:** [.lovable/plans/completed/80-vmware-shared-mount-fix-install-and-root-help.md](../../completed/80-vmware-shared-mount-fix-install-and-root-help.md)  
**Status:** Completed  
**Objective:** Remediate `Error -107 cannot open connection!` mount failure with pre-mount unmount, service restart, alternative mount fallback, and diagnostic guidance.

## Scope & Implementation Details
- In `gitmap/cmd/vmware_shared.go`:
  - `unmountIfMounted(mountPoint string)`: If mountPoint is currently mounted or has a broken fuse endpoint, run `fusermount -u` or `umount -l`.
  - `ensureVMwareToolsRunning()`: Check and trigger `systemctl start open-vm-tools` or `systemctl restart open-vm-tools`.
  - `getHostShares()`: Run `vmware-hgfsclient` to check for active host shares.
  - `mountHostShareWithFallback(mountPoint string)`:
    - Primary attempt: `sudo vmhgfs-fuse -o allow_other -o auto_unmount .host:/ <mountPoint>`
    - If failed, attempt fallback: `sudo mount -t fuse.vmhgfs-fuse .host:/ <mountPoint> -o allow_other`
    - If still failing with Error -107: check `getHostShares()`. If empty, format rich, user-friendly diagnostic explaining how to enable shared folders in VMware Workstation / Player settings.
- Ensure `checkVMwarePrerequisites()` installs both `open-vm-tools` AND `open-vm-tools-desktop`.
- Coding guidelines: functions <= 15 lines, blank lines before return, positive booleans, 0 nested ifs.
