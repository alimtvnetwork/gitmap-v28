# Subtask 03: Ubuntu Build-Essential Profile & Multi-Tool SQLite Tracking

## Status
Completed

## Context & Objectives
Implement the compact Ubuntu developer profile `build-essential`:
1. **Package Profile Invariants**:
   - Packages to install: `build-essential`, `vim`, `wget`, `nano`, `curl`, `file`, `git`, `zlib1g`, `zlib1g-dev`, `libssl-dev`, `git-core`, `sshpass`, `zsh`, `git-lfs`, `snapd`.
   - **CRITICAL**: Explicitly exclude `libpcre3-dev` (incompatible with modern Ubuntu 24.04+ releases; replacing with native PCRE2 where needed or omitting deprecated package).
2. **Compact CLI Invocations**:
   - `gitmap install build-essential` (aliases: `gitmap in build-essential`, `gitmap install ubuntu-common`, `gitmap in ub-common`).
3. **Constituent Tool Probing & SQLite Persistence**:
   - After invoking `apt-get install -y <packages>`, probe each constituent tool (`build-essential`, `gcc`, `g++`, `make`, `git`, `git-lfs`, `curl`, `wget`, `vim`, `nano`, `zsh`, `file`, `sshpass`, `snapd`).
   - Call `db.SaveInstalledTool(pkg, version, "apt")` for:
     - The umbrella profile `build-essential`.
     - Every detected constituent tool with its specific detected version.
   - Print a tree hierarchy summary (`printProfileTree`) highlighting installed constituent tools.

## Files to Create / Modify
- [MODIFY] `gitmap/constants/constants_install.go`: Define `ToolBuildEssential = "build-essential"`, `ToolUbuntuCommon = "ubuntu-common"`, update tool categories and descriptions.
- [MODIFY] `gitmap/cmd/install_profile_tree.go`: Add `buildUbuntuBuildEssentialProfile() ProfileComposition` and register in `resolveProfileTree`.
- [NEW] `gitmap/cmd/install_buildessential.go` (<= 200 lines): Profile installer handler, constituent probe logic, and SQLite `InstalledTool` multi-entry recording.
- [NEW] `gitmap/cmd/install_buildessential_test.go` (<= 200 lines): Unit tests for profile resolution, constituent probing, and package exclusion checks.

## Verification Steps
- `go test -v ./gitmap/cmd/... -run "TestBuildEssential"` passes cleanly.
- Verify that `libpcre3-dev` is NOT present in any package list.
- Verify `InstalledTool` table receives entries for both the profile and constituent tools.
