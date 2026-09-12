# RCA: Unused `isHelpArg` Function in `cmdvmware/vmware.go`

## 1. Why it happened
During package modularization and extraction of VMware CLI handlers into `cmdvmware`, a local helper function `isHelpArg` was copied into `vmware.go` while `checkHelp(constants.CmdVmware, args)` was used for top-level help inspection. `golangci-lint`'s `unused` analyzer flagged the dead private function as an error under strict CI settings.

## 2. How it happened
The `Run(args []string)` entrypoint in `gitmap/cmdvmware/vmware.go` delegates help detection to `checkHelp`, which immediately prints help text and exits if `-h`, `--help`, or `help` is passed. Consequently, the helper function `isHelpArg` was never invoked from any codepath within the package.

## 3. Root Cause
- File: `gitmap/cmdvmware/vmware.go`
- Lines: 25-27
- Unused private symbol `func isHelpArg(arg string) bool`

## 4. Code Fix
Removed `isHelpArg` from `gitmap/cmdvmware/vmware.go`.

### Before:
```go
func isHelpArg(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}
```

### After:
```go
// (removed)
```
