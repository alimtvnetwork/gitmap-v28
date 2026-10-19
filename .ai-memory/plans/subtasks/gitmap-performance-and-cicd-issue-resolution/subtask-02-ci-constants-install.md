# Subtask 02: Decompose cli/constants/constants_install.go Below 300-Line Cap

- **Subtask ID**: `subtask-02-ci-constants-install`
- **Parent Task**: `gitmap-performance-and-cicd-issue-resolution`
- **Owned Files**:
  - `cli/constants/constants_install.go`
  - `cli/constants/constants_install_tools.go`
  - `cli/constants/constants_install_pkgs.go`
- **Assigned Worker**: Worker 02
- **Status**: PENDING

---

## 1. Objectives
1. Analyze `cli/constants/constants_install.go` (currently 755 lines).
2. Decompose into modular, cohesive files:
   - `cli/constants/constants_install.go`: CLI install commands, base tools, and core help text ($\le 250$ lines).
   - `cli/constants/constants_install_tools.go`: Secondary and ecosystem tool constants ($\le 250$ lines).
   - `cli/constants/constants_install_pkgs.go`: Package manager constants, script templates, and command-line flags ($\le 250$ lines).
3. Ensure all identifiers and package declarations remain in package `constants` with zero breaking changes.
4. Verify all files are strictly $\le 300$ lines.

---

## 2. Verification
- `wc -l cli/constants/constants_install*.go`
- `go test ./cli/constants/... -v`
