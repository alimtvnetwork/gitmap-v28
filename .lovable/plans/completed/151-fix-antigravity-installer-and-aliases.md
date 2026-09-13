# Plan 151: Fix Antigravity & Agy Installer and Aliases

- **Slug**: `fix-antigravity-installer-and-aliases`
- **Status**: completed
- **Target**: Fix `gitmap install agy` and `gitmap agy install` failures by correcting alias mapping, default target routing, and desktop fallback.

## Context & Root Cause
When running `gitmap install agy` on Linux, GitMap previously reported:
```
Installing Google Antigravity Desktop IDE...
✗ Post-install verification failed for antigravity.
gitmap: [E9000:EXECUTION] cmd.verifyInstallation: post-install verification failed for "antigravity"
```
The root causes were:
1. `cli/cmdinstall/install_packages.go`: `"agy"` and `"ag"` aliases were mapped to `constants.ToolAntigravity` instead of `constants.ToolAgy`.
2. `cli/cmdinstall/agy_install.go`: `runAgyInstallCmd` defaulted target to `"ide"` and routed `"agy"` to `runInstallAntigravityWithOpts`.
3. `cli/cmdinstall/installantigravity.go`: If desktop archive returns 404, it immediately crashed without falling back to `agy` CLI installer.
4. `cli/cmdinstall/installantigravity_other.go`: Desktop finder on Linux only searched for `"antigravity"`, unlike Windows which checked `"agy"` as well.

## Subtasks
- [x] `01-fix-agy-tool-alias-and-default-target.md`: Correct alias mappings in `install_packages.go` and default target in `agy_install.go`.
- [x] `02-fix-antigravity-desktop-fallback-and-verification.md`: Implement clean `fallbackToAgyCli` in `installantigravity.go` and update `installantigravity_other.go`.
- [x] `03-verification-and-test-suite.md`: Update test suites, run unit and simulated tests, build binary and sync 4 locations.
