# Subtask 04: Expand scripts-fixer Tools and Commands

**Parent Plan:** [77-scripts-fixer-installation-split-db-and-tooling-engine.md](../../pending/77-scripts-fixer-installation-split-db-and-tooling-engine.md)
**Status:** Completed
**Target Files:**
- `gitmap/constants/constants_install.go`
- `gitmap/cmd/install.go`
- `gitmap/cmd/installtools.go`
- `gitmap/cmd/install_packages.go`
- `gitmap/cmd/install_packages_extra.go`
- `gitmap/cmd/install_packages_test.go`

## Objectives
1. Add tool constants and descriptions in `gitmap/constants/constants_install.go` matching `scripts-fixer`:
   - Languages & Runtimes: `rust`, `dotnet`, `java`, `flutter`
   - Local AI: `ollama`, `llama-cpp`, `python-libs`
   - DevOps & Containers: `docker`, `kubernetes`, `jenkins`
   - Terminal & Utilities: `zsh`, `flameshot`, `conemu`, `vlc`
2. Configure package manager resolution maps (`choco`, `winget`, `apt`, `brew`, `snap`) for the newly added tools.
3. Wire tool validation in `gitmap/cmd/install.go`.

## Verification
- Unit test suite `gitmap/cmd/install_packages_test.go` verifies all 14 new tools across package managers (choco, winget, apt, brew), alias resolution (k8s, kubectl, dotnet-sdk, jdk, openjdk, rustup, cargo, llamacpp), descriptions, and categories.
- `go test ./constants/... ./cmd/...` passes with all tests succeeding.
