# 03-task-root-help-integration.md: VMware Root Help Text & Documentation Integration

**Parent Plan:** [.lovable/plans/completed/80-vmware-shared-mount-fix-install-and-root-help.md](../../completed/80-vmware-shared-mount-fix-install-and-root-help.md)  
**Status:** Completed  
**Objective:** Register `vmware` in root help screens, compact layout, filter rows, and update markdown help documentation.

## Scope & Implementation Details
- In `gitmap/constants/constants_helpgroups.go`:
  - Define `HelpVmware = "  vmware (vm) <sub>           Manage VMware guest shared folders, status, and tools"`.
  - Update `CompactIntegrations` or `CompactEnvTools` to include `vmware (vm)`.
- In `gitmap/cmd/rootusage_groups.go`:
  - In `printGroupIntegrations()`, add `renderLine(constants.HelpVmware)`.
- In `gitmap/cmd/rootusagefilter_rows.go`:
  - In `allHelpRows()`, add `constants.HelpVmware` to `HelpGroupIntegrations`.
- In `gitmap/helptext/vmware.md`:
  - Document the `install (in)` subcommand.
  - Document the `shared enable (mount)` and `shared status` subcommands.
  - Include troubleshooting section detailing how to resolve `Error -107 cannot open connection!` in VMware host settings.
