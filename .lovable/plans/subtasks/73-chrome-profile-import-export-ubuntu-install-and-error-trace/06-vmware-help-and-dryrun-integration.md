# Subtask 06: VMware Help & Dry-Run Integration

## Scope
- Update `gitmap/cmd/vmware.go` and `gitmap/cmd/vmware_shared.go`:
  - Wire `checkHelp("vmware", args)` so `gitmap vmware help` and `gitmap vmware shared help` show structured help.
  - Support `--dry-run` flag in `runVmwareSharedEnable`.
  - Auto-probe `open-vm-tools` and `vmhgfs-fuse` with prompt/auto-install if missing.
- Update `gitmap/constants/constants_cli.go`, `gitmap/cmd/rootusage_groups.go`, and `gitmap/helptext/catalog.go`:
  - Register `HelpVmware` in constants and catalog.
- Verify commands (`gitmap vmware shared enable`, `status`, `gitmap os fix-link`) and render full help text.

## Files Touched
- `gitmap/cmd/vmware.go`
- `gitmap/cmd/vmware_shared.go`
- `gitmap/constants/constants_cli.go`
- `gitmap/cmd/rootusage_groups.go`
- `gitmap/helptext/catalog.go`
