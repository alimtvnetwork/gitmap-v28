# Subtask 02: VMware Shared Folders & Tools CLI Integration

## Status
Completed

## Context & Objectives
Implement the CLI command `gitmap vmware shared enable` (and companion subcommands):
1. **Hypervisor & Environment Detection**:
   - Check if running on Linux (`runtime.GOOS == "linux"`). On other OSes, return graceful explanation.
   - Verify VMware environment by probing `/sys/class/dmi/id/sys_vendor`, `/sys/class/dmi/id/product_name`, or `systemd-detect-virt`.
2. **Tooling Detection & Installation**:
   - Check if `/usr/bin/vmhgfs-fuse` or `open-vm-tools` exists.
   - If missing, offer/trigger `sudo apt-get install -y open-vm-tools open-vm-tools-desktop`.
   - Verify `/usr/bin/vmhgfs-fuse --enabled`.
3. **Mount Point Provisioning**:
   - Ensure `/mnt/hgfs` directory exists (`mkdir -p /mnt/hgfs`).
   - Execute mount: `sudo vmhgfs-fuse .host:/ /mnt/hgfs -o allow_other` (or `-o allow_other -o auto_unmount`).
4. **Desktop Symlink Creation**:
   - Resolve user desktop directory (checking `$SUDO_USER` if run under sudo so the desktop link lands on the actual user's desktop, e.g. `/home/<user>/Desktop/SharedDirectories`, falling back to `~/Desktop/SharedDirectories`).
   - Create symlink idempotently: `ln -s /mnt/hgfs <Desktop>/SharedDirectories`.
5. **Startup Crontab Persistence**:
   - Inspect `crontab -l`.
   - Check if an entry for `vmhgfs-fuse .host:/ /mnt/hgfs` already exists.
   - If not present, append `@reboot /usr/bin/vmhgfs-fuse .host:/ /mnt/hgfs -o allow_other` idempotently and reload crontab.
6. **CLI Command Routing**:
   - `gitmap vmware` (alias `vm`)
   - `gitmap vmware shared enable` (aliases: `shared-enable`, `mount`)
   - `gitmap vmware status` and `gitmap vmware shared status`

## Files to Create / Modify
- [MODIFY] `gitmap/constants/constants_cli.go`: Add `CmdVmware = "vmware"`, `CmdVmwareAlias = "vm"`, and subcommands under skip.
- [MODIFY] `gitmap/constants/cmd_constants_test.go`: Add `CmdVmware` and `CmdVmwareAlias` to `topLevelCmds()`.
- [MODIFY] `gitmap/cmd/roottooling.go`: Register `CmdVmware` in `toolingInstallEntries()`.
- [MODIFY] `gitmap/cmd/rootutility.go`: Map `CmdVmwareAlias` in `canonicalCommandName`.
- [NEW] `gitmap/cmd/vmware.go` (<= 200 lines): Umbrella command dispatcher and help routing.
- [NEW] `gitmap/cmd/vmware_shared.go` (<= 200 lines): Tool detection, `/mnt/hgfs` mount, desktop symlink, and crontab persistence.
- [NEW] `gitmap/cmd/vmware_test.go` (<= 200 lines): Unit tests for command parsing, flag validation, and mock detection.

## Verification Steps
- AST parity test: `go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1` exits 0.
- Command tests: `go test -v ./gitmap/cmd/... -run "TestVmware"` exits 0.
