# Subtask 05: Changelog Synchronization, Help Text Parity & CI/CD Verification

## Status
Pending

## Context & Objectives
Ensure complete documentation, help parity, and full CI/CD quality gate compliance:
1. **Help Text Creation & Golden File Parity**:
   - Create `gitmap/helptext/vmware.md` (<= 120 lines, 3-8 line execution simulation block).
   - Create `gitmap/helptext/server-cmd.md` (<= 120 lines, 3-8 line execution simulation block).
   - Register both in `gitmap/helptext/catalog.go`.
   - Update `gitmap/cmd/rootusage_groups.go` so `vmware` and `server-cmd` appear in group listings.
   - Run `go test ./gitmap/helptext/... -run Golden -count=1` to guarantee golden help text parity.
2. **Changelog & Documentation**:
   - Update `changelog.md` with features: VMware shared folders, Ubuntu build-essential profile, staging directory rules, and server-cmd remote execution.
   - Mirror update to `gitmap/changelog.md` and `src/data/changelog.ts`.
3. **Coding Guidelines Verification**:
   - Run `python linter-scripts/check-relative-paths.py`.
   - Run `python 03-ai-scripts/14-version-sync-checker.py`.
   - Run `python 03-ai-scripts/06-cicd-local-runner.py` ensuring all 33 quality gates pass with exit code 0.

## Files to Create / Modify
- [NEW] `gitmap/helptext/vmware.md`
- [NEW] `gitmap/helptext/server-cmd.md`
- [MODIFY] `gitmap/helptext/catalog.go`
- [MODIFY] `gitmap/cmd/rootusage_groups.go`
- [MODIFY] `changelog.md`
- [MODIFY] `gitmap/changelog.md`
- [MODIFY] `src/data/changelog.ts`

## Verification Steps
- `go test ./gitmap/helptext/... -run Golden -count=1` exits 0.
- `python 03-ai-scripts/06-cicd-local-runner.py` exits 0 across all gates.
