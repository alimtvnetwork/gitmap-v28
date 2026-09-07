# 70-vmware-shared-folders-and-ubuntu-profiles.md: VMware Shared Folder Integration, Ubuntu Build-Essential Profiles & Server-Cmd Cluster Remote Execution

## 1. Executive Summary

This plan introduces three major functional capabilities into the Gitmap ecosystem alongside disciplined download staging rules:
1. **VMware Shared Folders & Tools (`gitmap vmware shared enable`)**: Automates detection and installation of `open-vm-tools`/`open-vm-tools-desktop`, verifies `/usr/bin/vmhgfs-fuse --enabled`, mounts host shares at `/mnt/hgfs` (`vmhgfs-fuse .host:/ /mnt/hgfs -o allow_other`), creates desktop symlink `~/Desktop/SharedDirectories -> /mnt/hgfs`, and ensures `@reboot` crontab persistence.
2. **Ubuntu Common "build-essential" Profile & Multi-Tool Tracking**: Provides a compact installer profile for Ubuntu/Debian developer environments containing `build-essential`, `vim`, `wget`, `nano`, `curl`, `file`, `git`, `zlib1g`, `zlib1g-dev`, `libssl-dev`, `git-core`, `sshpass`, `zsh`, `git-lfs`, and `snapd` (explicitly excluding deprecated `libpcre3-dev` to ensure compatibility with Ubuntu 24.04+). Persists `build-essential` and each constituent package into SQLite `InstalledTool` table with detected versions.
3. **Linux Download Staging & Persistent Keep Architecture**: Implements a strict two-tier download architecture: transient in-flight downloads stage inside `/tmp` (or `mktemp -d`), while verified persistent assets relocate into `~/.gitmap-installation` (or `/root/.gitmap-installation` for root) with `EXDEV` cross-device link resilience (`io.Copy` fallback).
4. **Server-Cmd / SSH / Kubernetes Cluster Remote Execution (`gitmap server-cmd`)**: Based on `aukgit/kubernetes-training-v1` `05-server-cmds/02-run-cmd-v2.sh`, provides `gitmap server-cmd` (aliases: `server-cmds`, `scmd`) bridging cluster daemon and SSH pools across target selectors (`all`, `control`, `workers`, `worker-N`), supporting direct commands, `--sudo` privilege escalation, and on-the-fly script execution (`/tmp/on-the-fly-cmd/script-<ts>.sh`).

## 2. Task-Specific Rules & Constraints

1. **Rule 1 (Strict Relative Paths):** Zero `file:///` URIs or machine-specific absolute paths in markdown files, artifacts, or code. All paths must be relative to repository root (`/`).
2. **Rule 2 (Coding Guidelines Compliance):** Functions $\le 15$ lines, blank line before every return statement, affirmative booleans (`is_*`, `has_*`), no magic numbers, and zero swallowed errors (`err != nil` must be handled or wrapped via `apperror`).
3. **Rule 3 (Strict 200-Line File Limit):** All new Go source files must remain $\le 200$ lines. Decompose logic into focused single-responsibility files.
4. **Rule 4 (AST Parity & Constants):** All top-level CLI verbs must be declared in `gitmap/constants/constants_cli.go` under `// gitmap:cmd top-level`, synchronized in `gitmap/constants/cmd_constants_test.go` `topLevelCmds()`, and verified with `TestTopLevelCmdRegistryMatchesAST`.
5. **Rule 5 (No Premature Releases or Bumps):** Commits must remain standard development commits. Do not edit `version.json` or cut release tags.
6. **Rule 6 (Cross-Device EXDEV Safety):** Moving files between `/tmp` and `~/.gitmap-installation` must account for `/tmp` mounted on `tmpfs`. If `os.Rename` fails with `EXDEV`, fall back to buffered copy and unlink.

## 3. Subsystems & Architecture

### A. VMware Shared Folders (`gitmap vmware shared enable`)
- **CLI Commands**: `gitmap vmware` (alias `vm`), subcommands `shared enable` (aliases `shared-enable`, `mount`), `status`, `shared status`.
- **Hypervisor Probe**: Inspect `/sys/class/dmi/id/sys_vendor`, `/sys/class/dmi/id/product_name`, or `systemd-detect-virt` for `vmware`.
- **Tooling Verification**: Check `open-vm-tools`, `vmhgfs-fuse --enabled`. Auto-install via `apt-get` if missing.
- **Mount & Symlink**: Ensure `/mnt/hgfs` exists, execute `vmhgfs-fuse .host:/ /mnt/hgfs -o allow_other`, and create `~/Desktop/SharedDirectories` pointing to `/mnt/hgfs` (resolving `SUDO_USER` if run under sudo).
- **Crontab Persistence**: Read `crontab -l`, append `@reboot /usr/bin/vmhgfs-fuse .host:/ /mnt/hgfs -o allow_other` idempotently, and reinstall.

### B. Two-Tier Staging & Persistent Keep Architecture
- **Tier 1 (Transient Staging)**: In-flight downloads download into `/tmp/gitmap-stage-*` or `filepath.Join(os.TempDir(), "gitmap-downloads")`.
- **Tier 2 (Persistent Keep)**: `~/.gitmap-installation` (or `/root/.gitmap-installation` for root). Structure:
  - `downloads/` — Cached installers and archives.
  - `scripts/` — Provisioning and helper scripts.
  - `logs/` — Installation execution logs.
- **Atomic Promotion**: `PromoteStagedFile(src, dst)` handles atomic rename with `EXDEV` fallback copy.

### C. Ubuntu "build-essential" Profile & Multi-Tool Tracking
- **Profile Slug**: `build-essential` (aliases: `be`, `ubuntu-common`).
- **Packages**: `build-essential`, `vim`, `wget`, `nano`, `curl`, `file`, `git`, `zlib1g`, `zlib1g-dev`, `libssl-dev`, `git-core`, `sshpass`, `zsh`, `git-lfs`, `snapd`. Explicitly exclude `libpcre3-dev`.
- **Constituent Probing**: After `apt-get install`, run probes for `build-essential`, `gcc`, `g++`, `make`, `git`, `git-lfs`, `curl`, `wget`, `zsh`, `vim`, `nano`.
- **SQLite Persistence**: Call `db.SaveInstalledTool(pkg, version, "apt")` for the umbrella profile AND each constituent tool.

### D. Server-Cmd Cluster & SSH Delegation (`gitmap server-cmd`)
- **CLI Commands**: `gitmap server-cmd` (aliases `server-cmds`, `scmd`).
- **Target Selectors**:
  - `all`: Targets both control plane servers and client worker nodes (`cluster.ServersClients`).
  - `control`: Targets control plane servers only (`cluster.ServersOnly`).
  - `workers`: Targets worker client nodes only (`cluster.ClientsOnly`).
  - `<node-id-or-alias>`: Targets specific node by ID, display ID, or hostname.
- **Execution Modes**:
  - Direct command execution over SSH session.
  - `--script`: Packages remote script to `/tmp/on-the-fly-cmd/script-<ts>.sh`, executes remotely, and cleans up.
  - `--sudo`: Runs with non-interactive `sudo -S` escalation.
  - Dual pool resolution: Database `ClusterNode` / `SSHConnection` with fallback to `--config <path>` JSON file.

## 4. Subtasks Breakdown

1. [01-install-staging-and-keep-dir.md](subtasks/70-vmware-shared-folders-and-ubuntu-profiles/01-install-staging-and-keep-dir.md): Implement two-tier staging helper (`gitmap/downloaderconfig/staging.go`), temporary directory management, `~/.gitmap-installation` path resolution, and `EXDEV` atomic relocation.
2. [02-vmware-shared-command.md](subtasks/70-vmware-shared-folders-and-ubuntu-profiles/02-vmware-shared-command.md): Implement `gitmap vmware` and `gitmap vmware shared enable` in `gitmap/cmd/vmware.go` and `gitmap/cmd/vmware_shared.go`, tool probing, mount execution, desktop symlinking, and crontab persistence.
3. [03-build-essential-profile.md](subtasks/70-vmware-shared-folders-and-ubuntu-profiles/03-build-essential-profile.md): Implement `build-essential` profile in `gitmap/cmd/install_profile_tree.go` and `gitmap/cmd/install_buildessential.go`, multi-tool probing, and SQLite `InstalledTool` multi-entry persistence.
4. [04-server-cmd-integration.md](subtasks/70-vmware-shared-folders-and-ubuntu-profiles/04-server-cmd-integration.md): Implement `gitmap server-cmd` in `gitmap/cmd/servercmd.go` and `gitmap/cmd/servercmd_run.go`, target routing (`all`, `control`, `workers`), `--script` on-the-fly runner, and `--sudo` escalation.
5. [05-changelog-docs-and-ci.md](subtasks/70-vmware-shared-folders-and-ubuntu-profiles/05-changelog-docs-and-ci.md): Help text authoring (`gitmap/helptext/vmware.md`, `gitmap/helptext/server-cmd.md`), changelog synchronization (`changelog.md`, `gitmap/changelog.md`, `src/data/changelog.ts`), AST parity verification, and full CI quality gate verification.

## 5. Acceptance Criteria

- [x] AST Parity: `go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1` exits 0.
- [x] Unique Top-Level Constants: `go test ./gitmap/constants/... -run TestTopLevelCmdConstantsAreUnique -count=1` exits 0.
- [x] Help Text Parity: `go test ./gitmap/helptext/... -run Golden -count=1` exits 0.
- [x] All new and modified Go files are $\le 200$ lines, functions $\le 15$ lines, blank lines before returns.
- [x] Staging helper cleanly handles `/tmp` isolation, `~/.gitmap-installation` directory tree, and `EXDEV` cross-device links.
- [x] `gitmap vmware shared enable` checks hypervisor, verifies tools, mounts `/mnt/hgfs`, creates desktop symlink, and sets `@reboot` crontab idempotently.
- [x] `gitmap install build-essential` installs complete Ubuntu common package list (omitting `libpcre3-dev`) and registers all constituent tools in SQLite `InstalledTool`.
- [x] `gitmap server-cmd` executes commands across `all`, `control`, and `workers` targets with `--script` and `--sudo` options.
- [x] `python 03-ai-scripts/06-cicd-local-runner.py` passes all gates with exit code 0.
