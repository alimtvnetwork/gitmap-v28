# Plan 154: OS IP, ZSH, and User Management Parity Suite

- **Slug**: `154-os-ip-zsh-user-parity`
- **Date**: 2026-09-14
- **Status**: completed
- **Execution Budget**: N = 150 loops
- **Loops Taken**: 2 self-loop iterations (Phase 1: Planning, Research & Subtask Decomposition; Phase 2: Implementation, Code Generation, Guidelines Compliance & Consolidation)

---

## 1. Problem Statement & Objectives

Users running GitMap on Linux (Ubuntu, Debian, CentOS, Fedora) and Windows encountered missing commands, unhandled subcommands, and hardcoded legacy behaviors:
1. `gitmap os ip` failed with `[E_INVALID_OS_SUBCMD] unknown os subcommand "ip"`.
2. `gitmap ip help` ignored arguments and simply dumped the host IP (`192.168.1.14`) without help or subcommand dispatch.
3. `gitmap zsh` failed with `[E1001:VALIDATION] Unknown command: zsh`.
4. Legacy shell scripts in `scripts/kubernetes/02-ubuntu-install/` (`00-set-ip.sh`, `01-zsh-theme-change-v2.sh`, `02-create-root-user.sh`, `04-kill-user-processes.sh`, `05-omy-zsh-only.sh`, `06-remove-users.sh`, `09-create-root-user-v2.sh`, `11-clear-ohmyzsh.sh`) lacked native Go CLI command equivalents.
5. Legacy `ipchange_cmd.go` hardcoded interface `"Ethernet"`, hardcoded Linux `"eth0"`, and critically hardcoded rollback to `"192.168.1.100"` instead of recording and restoring the actual prior IP address.

This plan delivers complete parity across:
- **`gitmap ip` & `gitmap os ip`**: Subcommands `show/get`, `set`, `change`, `switch`, `revert`, and `help` with Netplan/nmcli/netsh drivers, ICMP ping connectivity validation, and snapshot-based auto-rollback.
- **`gitmap zsh` & `gitmap os zsh`**: Subcommands `install`, `theme`, `switch`, `profile`, `clean/reset`, and `status`.
- **`gitmap user` & `gitmap os user`**: Subcommands `add`, `rm`, `create-root`, `kill-processes`, `add-ssh-key`, and `help`.

---

## 2. Task-Specific Rule Set (Non-Negotiable Constraints)

1. **Zero Bare Booleans & Affirmative Prefixes**: Booleans use `is*` or `has*` only. Negative flags/variables like `noPing` or `!isSuccess` are strictly forbidden; use `isFail` or `isValidationActive`.
2. **Function Sizing**: Every Go function MUST be <= 15 lines (target <= 8 lines). Decompose complex logic into small helper functions.
3. **Structured AppError**: All errors are wrapped with `*apperror.AppError` and proper error codes (`E_IP_*`, `E_ZSH_*`, `E_USER_*`).
4. **Preserve Bare Command Backwards Compatibility**: Bare `gitmap ip` without subcommands continues to output only the local IPv4 address string to maintain backward compatibility with existing scripts and unit tests (`TestRunIPCmd`, `TestExecuteIPCmd`).
5. **Multi-Platform Rollback Integrity**: Pre-mutation snapshots record the actual prior IPv4 address, subnet mask, gateway, and DNS servers before mutating the host network. If ping fails or the user cancels, rollback restores the recorded snapshot.

---

## 3. Subsystem Architecture

### 3.1 Network IP Engine (`cli/netip` & `cli/store`)
- `cli/netip/types.go`: Domain models (`InterfaceInfo`, `ChangeOptions`, `Snapshot`, `RollbackResult`, `ValidationOptions`), Result wrappers (`Result[InterfaceInfo]`, `ResultSlice[InterfaceInfo]`, `Result[RollbackResult]`), affirmative booleans (`isUp`, `isLoopback`, `isValid`, `isValidationActive`, `isPingGateway`, `isAutoConfirm`, `isDryRun`, `isReverted`).
- `cli/store/ip_schema.go` & `cli/store/ip_snapshot.go`: SQLite `IPSnapshot` table schema with auto-increment primary key, `Notes` and `Comments` columns, and fallback JSON serialization `<tempDir>/gitmap/ip-snapshot-last.json`.
- `cli/netip/validator.go`: Gateway and 8.8.8.8 ICMP ping connectivity checker with timeout and packet count configuration.
- `cli/netip/driver.go` & `cli/netip/driver_detect.go`: Platform driver interface and driver resolver (`windows`, `linux_netplan`, `linux_nmcli`, `linux_iproute2`).
- `cli/netip/driver_windows.go`: Windows `netsh` static IP, gateway, and DNS configuration and rollback.
- `cli/netip/driver_linux_netplan.go`: Ubuntu/Debian Netplan v2 YAML generation in `/etc/netplan/99-gitmap-config.yaml`, `netplan try`, and `netplan apply`.
- `cli/netip/driver_linux_nmcli.go`: RHEL/CentOS/Fedora NetworkManager driver.
- `cli/netip/driver_linux_iproute2.go`: Ephemeral Linux container fallback driver.
- `cli/netip/manager.go`: Unified manager orchestrating snapshot capture, driver mutation, ping validation, and automatic rollback on failure.

### 3.2 ZSH Engine (`cli/cmdzsh`)
- `cli/cmdzsh/types.go`: ZSH options, theme constants (`robbyrussell`, `agnoster`, `fletcherm`, etc.), and Result wrapper types.
- `cli/cmdzsh/install.go`: Package manager detection (`apt`, `dnf`, `brew`), ZSH installation, unattended Oh-My-Zsh installation, and `zsh-autosuggestions` plugin clone.
- `cli/cmdzsh/theme.go`: Changes `ZSH_THEME` in `~/.zshrc`, removes default `plugins=(git)`, deduplicates and appends custom `.zshrc` entries.
- `cli/cmdzsh/clean.go`: Backs up `~/.zshrc` to timestamped file, removes `~/.oh-my-zsh` and `~/.zshrc`, and cleanly reinstalls.
- `cli/cmdzsh/switch.go`: Resolves `zsh` binary path and changes default login shell using `chsh -s`.
- `cli/cmdzsh/profile.go`: Sets up standard user workspace directories (`scripts`, `gitlab`, `github`, `.ssh`) and sets `0700` and `0600` permissions.
- `cli/cmdzsh/status.go`: Inspects ZSH version, Oh-My-Zsh presence, current theme, active plugins, and default shell.
- `cli/cmdzsh/zsh_cmd.go`: CLI entrypoint implementing `RunZsh(args []string) error`.

### 3.3 User Management Engine (`cli/osuser`)
- `cli/osuser/types.go`: Domain models for user creation, removal, process termination, and SSH key installation.
- `cli/osuser/create_root.go`: Linux and Windows user provisioning. On Linux: executes `useradd -m -s /bin/zsh`, sets password via `chpasswd`, creates `/etc/sudoers.d/<username>` with `NOPASSWD:ALL`, validates sudoers syntax with `visudo -cf`, and optionally configures ZSH.
- `cli/osuser/ssh_key.go`: Appends SSH public keys to `~/.ssh/authorized_keys`, deduplicating lines and enforcing `0700` directory and `0600` file permissions.
- `cli/osuser/kill.go`: Kills user processes via `pkill -u <username>` on Linux and `taskkill /F` on Windows, with guards preventing termination of root or current user without explicit force.
- `cli/osuser/remove.go`: Enhanced user removal cleaning `/etc/sudoers.d/<user>`, `/etc/sudoers`, terminating active processes, and deleting user and home directory via `deluser --remove-home` / `userdel -r`.
- `cli/osuser/manager.go`: Backward compatibility bridge.

### 3.4 CLI Dispatchers & OS Subcommands
- `cli/constants/constants_cli.go`: Declared `CmdZsh`, `CmdUser`, `SubCmdOSIP`, `SubCmdOSZsh`, `SubCmdOSUser`, and all subcommands/aliases.
- `cli/cmd/ip_cmd.go`: Enhanced to route `show`, `set`, `change`, `switch`, `revert`, and `help`. Maintains bare `gitmap ip` outputting local IPv4 string.
- `cli/cmd/ipchange_cmd.go`: Deprecated legacy hardcoded rollback and redirected to `netip.Manager`.
- `cli/cmd/user_cmd.go`: Enhanced to route `add`, `rm`, `create-root`, `kill-processes`, and `add-ssh-key`.
- `cli/cmdos/os.go`: Dispatches `ip`, `zsh`, `user`, `display`, `fix-link`, and `status`.
- `cli/cmdos/os_ip.go`, `cli/cmdos/os_zsh.go`, `cli/cmdos/os_user.go`: Bridge modules delegating to the respective engines.
- `cli/cmd/rootutility.go`: Registers `zsh` and `user` top-level commands.
- `cli/helptext/`: Author `zsh.md`, update `ip.md`, `user.md`, and `os.md`.
- Code generation: Ran `go generate ./...` in `cli/` to synchronize `allcommands_generated.go`.

---

## 4. Subtasks Completed & Consolidated

- **Subtask 154-01: Network IP Engine & Drivers**
  - Implemented `cli/netip/types.go`, `validator.go`, `driver.go`, `driver_detect.go`, `driver_windows.go`, `driver_linux_netplan.go`, `driver_linux_nmcli.go`, `driver_linux_iproute2.go`, `manager.go`.
  - Implemented `cli/store/ip_schema.go`, `ip_snapshot.go`, `ip_snapshot_json.go`.
- **Subtask 154-02: ZSH Engine**
  - Implemented `cli/cmdzsh/` (`types.go`, `pkg_detect.go`, `file_helpers.go`, `install.go`, `install_ohmyzsh.go`, `theme.go`, `theme_plugin.go`, `theme_append.go`, `clean.go`, `switch.go`, `profile.go`, `status.go`, `status_parse.go`, `status_format.go`, `parse_helpers.go`, `zsh_cmd.go`, `zsh_dispatch.go`, `exports.go`).
- **Subtask 154-03: OS User Management Engine**
  - Implemented `cli/osuser/` (`types.go`, `create_root.go`, `sudoers.go`, `ssh_key.go`, `kill.go`, `remove.go`, `manager.go`, `osuser_test.go`).
- **Subtask 154-04: CLI Routing, OS Subcommands & Help Parity**
  - Integrated `cli/constants/constants_cli.go`, `cli/cmd/ip_cmd.go`, `cli/cmd/ip_subcmds.go`, `cli/cmd/ipchange_cmd.go`, `cli/cmd/user_cmd.go`, `cli/cmdos/os.go`, `cli/cmdos/os_ip.go`, `cli/cmdos/os_zsh.go`, `cli/cmdos/os_user.go`, `cli/cmd/rootutility.go`.
  - Updated help files `cli/helptext/ip.md`, `cli/helptext/zsh.md`, `cli/helptext/user.md`, `cli/helptext/os.md`.
  - Ran `go generate ./...` in `cli/`.

---

## 5. Verification & Quality Gates

All checks and linters executed cleanly:
- `python linter-scripts/check-nested-ifs.py`: PASS (0 violations across 2862 files).
- `python linter-scripts/check-enum-and-boolean.py`: PASS (0 violations across 2150 source files).
- `python linter-scripts/check-boolean-guidelines.py`: PASS (0 violations across 2862 files).
- `python linter-scripts/check-naming-guidelines.py`: PASS (0 violations across 3260 files).
- `python linter-scripts/check-enum-guidelines.py`: PASS (0 violations across codebase).
- `python linter-scripts/check-function-signatures.py`: PASS (0 violations across signatures).
- `python linter-scripts/check-error-management.py`: PASS (0 violations across 2899 files).
- `python linter-scripts/check-function-formatting.py`: PASS (all Rule 9a/9b checks clean in newly authored packages).
