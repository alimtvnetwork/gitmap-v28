# 02-task-vmware-install-command.md: VMware Install Subcommand & Tooling Integration

**Parent Plan:** [.lovable/plans/completed/80-vmware-shared-mount-fix-install-and-root-help.md](../../completed/80-vmware-shared-mount-fix-install-and-root-help.md)  
**Status:** Completed  
**Objective:** Add `gitmap vmware install` subcommand and integrate `vmware` into the `gitmap install` tooling engine.

## Scope & Implementation Details
- In `gitmap/cmd/vmware.go`:
  - Register `install`, `in` subcommands in `dispatchVmwareSubcommand`.
  - Update `printVmwareUsage()` to show `install (in)` subcommand.
- In `gitmap/cmd/vmware_install.go`:
  - Implement `runVmwareInstall(args []string) error` with `--dry-run` and `-y` flags.
  - On Linux/Ubuntu: install `open-vm-tools open-vm-tools-desktop` via `apt-get`.
  - Enable and start `open-vm-tools` systemd service.
- In `gitmap/constants/constants_install.go`:
  - Define `ToolVMware = "vmware"`.
  - Add Package IDs: `PackageAptVMware = "open-vm-tools open-vm-tools-desktop"`, `PackageChocoVMware = "vmware-workstation-player"`, `PackageWingetVMware = "VMware.WorkstationPlayer"`, `PackageBrewVMware = "vmware-fusion"`.
  - Add descriptions, category mapping (`CategoryDevOps`), and tool aliases (`open-vm-tools`, `vmtools`, `vmware-tools`, `vm`).
- In `gitmap/cmd/install_packages.go` & `gitmap/cmd/install_packages_extra.go`:
  - Add `ToolVMware` entries to `chocoPackages`, `wingetPackages`, `aptPackages`, `brewPackages`.
  - Add tool aliases to `resolveToolAliases()`.
- In `gitmap/cmd/installverify.go`:
  - Map `ToolVMware` to binary name `vmware-toolbox-cmd` / `vmhgfs-fuse`.
