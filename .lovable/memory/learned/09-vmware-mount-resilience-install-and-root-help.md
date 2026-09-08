# VMware Shared Folders Mount Resilience, Install Command & Root Help Integration

## 1. Context & Objectives
Resolved VMware guest shared folders mounting failures on Linux, added the missing `gitmap vmware install` command, and registered VMware in the root CLI help text:
- Remediated `Error -107 cannot open connection! (exit status 149)` during `gitmap vmware shared enable`.
- Implemented `gitmap vmware install` (alias `in`) to install `open-vm-tools` and `open-vm-tools-desktop` via `apt-get` and enable `open-vm-tools` systemd service.
- Integrated `ToolVMware = "vmware"` across `gitmap install` package managers (`choco`, `winget`, `apt`, `brew`).
- Registered `vmware (vm)` in root help groups, compact listing, search filter rows, and updated `helptext/vmware.md`.

## 2. Key Technical Findings & Learned Patterns

### 1. VMware Linux Guest Error -107 (`ENOTCONN`) Architecture
- **Problem**: `vmhgfs-fuse` exits with `Error -107 cannot open connection! (exit status 149)` when the guest kernel transport endpoint cannot reach the host hypervisor.
- **Root Cause**:
  1. The host hypervisor (VMware Workstation, Player, Fusion, ESXi) does not have Shared Folders enabled in Virtual Machine Settings (Options -> Shared Folders).
  2. No folders are configured or enabled on the host side.
  3. Stale or interrupted FUSE mount left a dead endpoint at `/mnt/hgfs`.
  4. The `open-vm-tools` systemd service is inactive or not loaded.
- **Resilience Engine**:
  - **Pre-mount Unmount**: Call `unmountIfMounted(mountPoint)` using `fusermount -u` and `sudo umount -l` before attempting mount.
  - **Service Start**: Call `ensureVMwareServiceRunning()` to ensure `systemctl start open-vm-tools` is active.
  - **Fallback Mounting**: Try primary `vmhgfs-fuse` first; on failure, attempt fallback `mount -t fuse.vmhgfs-fuse .host:/ <mountPoint> -o allow_other`.
  - **Actionable Host Diagnostics**: Query `vmware-hgfsclient`. If Error -107 or no shares detected, emit clear step-by-step remediation directing the user to VMware VM Settings -> Options -> Shared Folders -> "Always enabled" -> Add Folder.

### 2. Desktop Package Requirement (`open-vm-tools-desktop`)
- **Problem**: Installing only `open-vm-tools` misses the X11/Wayland desktop integration agent (`vmtoolsd` user daemon), resulting in failure to resolve Desktop symlinks and auto-mounting permissions.
- **Solution**: Always install both `open-vm-tools` and `open-vm-tools-desktop` together via `apt-get install -y open-vm-tools open-vm-tools-desktop`.

### 3. Root Help Registration Surface
- **Problem**: Adding a CLI command requires updating multiple help surfaces for parity:
  1. `constants/constants_cli.go`: `HelpVmware` summary text.
  2. `constants/constants_helpgroups.go`: `CompactIntegrations` string and `HelpGroupKeys`.
  3. `cmd/rootusage_groups.go`: `printGroupIntegrations()` line rendering.
  4. `cmd/rootusagefilter_rows.go`: `allHelpRows()` group row registration.
  5. `cmd/rootusagecompact.go`: `compactGroups()` mapping entry.
  6. `helptext/vmware.md`: Markdown documentation with subcommand tables and troubleshooting guides.
