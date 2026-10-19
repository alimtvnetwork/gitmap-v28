# Subtask 03: Decompose cli/cmdinstall/install_packages.go and Tests Below 300-Line Cap

- **Subtask ID**: `subtask-03-ci-cmdinstall-packages`
- **Parent Task**: `gitmap-performance-and-cicd-issue-resolution`
- **Owned Files**:
  - `cli/cmdinstall/install_packages.go`
  - `cli/cmdinstall/install_packages_resolve.go`
  - `cli/cmdinstall/install_packages_test.go`
  - `cli/cmdinstall/install_packages_resolve_test.go`
- **Assigned Worker**: Worker 03
- **Status**: PENDING

---

## 1. Objectives
1. Split `cli/cmdinstall/install_packages.go` (403 lines):
   - Move `managerResolvers`, `resolvePackageName`, `resolveWingetPackage`, `resolveAptPackage`, `resolveBrewPackage`, `resolveSnapPackage`, `resolveChocoPackage`, and package alias tables into `cli/cmdinstall/install_packages_resolve.go` ($\le 250$ lines).
   - Keep execution orchestrator and runners in `cli/cmdinstall/install_packages.go` ($\le 250$ lines).
2. Split `cli/cmdinstall/install_packages_test.go` (431 lines):
   - Move resolver and package mapping unit tests into `cli/cmdinstall/install_packages_resolve_test.go` ($\le 250$ lines).
   - Keep runner and execution tests in `cli/cmdinstall/install_packages_test.go` ($\le 250$ lines).
3. Verify that `python 03-ai-scripts/52-file-size-check.py --diff HEAD~1 --limit 300` exits 0.

---

## 2. Verification
- `python 03-ai-scripts/52-file-size-check.py --diff HEAD~1 --limit 300`
- `go test ./cli/cmdinstall/... -v`
