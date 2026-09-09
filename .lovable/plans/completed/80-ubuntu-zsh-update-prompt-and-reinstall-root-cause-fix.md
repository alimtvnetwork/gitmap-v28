# Completed Specification: Ubuntu ZSH Update Prompt & Reinstall Root Cause Fix

## 1. Executive Problem Statement
During `gitmap update` on Ubuntu (as well as `gitmap self-install` and `gitmap cd` uninitialized wrapper recovery), `gitmap setup` is automatically invoked as a post-install hook. Because `ensureZshUbuntuStep` in `gitmap/cmd/setup_ubuntu.go` lacked pre-flight detection, TTY guards, and bypass flags:
1. It unconditionally prompted developers with `Install ZSH and Oh-My-Zsh? (y/N): `.
2. If approved, it ran `sudo apt install -y zsh` (requesting root password, and replacing custom or newer ZSH builds with the distribution package).
3. It ran the unattended Oh-My-Zsh installer (which backed up `.zshrc` to `.zshrc.pre-oh-my-zsh`, blowing away custom configurations, PATH, and gitmap's shell wrapper), forcefully set `ZSH_THEME="agnoster"`, and created a fresh Git checkout that subsequently prompted developers on new shells to update Oh-My-Zsh.

---

## 2. Acceptance Criteria
- [x] **Pre-Flight ZSH Detection**: If `zsh` is already installed, `ensureZshUbuntuStep` prints `✓ ZSH is already installed (<version>)` and NEVER prompts to install ZSH or runs `sudo apt install`.
- [x] **Pre-Flight Oh-My-Zsh Detection**: If `~/.oh-my-zsh` exists, `ensureZshUbuntuStep` prints `✓ Oh-My-Zsh is already installed` and NEVER reinstalls OMZ or mutates `.zshrc`.
- [x] **Early Short-Circuit**: When both ZSH and Oh-My-Zsh are present, `ensureZshUbuntuStep` returns immediately without blocking or prompting.
- [x] **Flag & Environment Guards**: `gitmap setup` supports `--skip-zsh` and reads `GITMAP_SKIP_ZSH=1`. When either is set, `ensureZshUbuntuStep` exits immediately.
- [x] **Non-Interactive TTY Guard**: If `os.Stdin` is not a character device (`!isStdinTerminal()`), `ensureZshUbuntuStep` does not block on user input.
- [x] **Update Hook Defense**: `install.sh` and `gitmap/scripts/install.sh` pass `GITMAP_SKIP_ZSH=1 "${bin_path}" setup --skip-zsh`.
- [x] **Self-Install & Wrapper Defense**: `selfinstall.go` and `setupverify.go` invoke `runSetup([]string{"--skip-zsh"})`.
- [x] **Tool Probe Registration**: `constants.ToolZsh` is registered in `toolProbeMap` in `gitmap/cmd/installprobe.go`.
- [x] **100-Line File Cap**: `setup_ubuntu.go` decomposed with `setup_ubuntu_detect.go` and `setup_ubuntu_install.go`, all files <= 100 lines.
- [x] **Unit Tests & Quality Gates**: Comprehensive unit tests added in `setup_ubuntu_test.go` and `installprobe_zsh_test.go`, all linters exit code 0.

---

## 3. Custom Rule Set (Task Constraints)
1. **Zero Password / Zero Sudo During Updates**: Automated update workflows (`gitmap update`, `install.sh`) must NEVER invoke `sudo` or request user passwords.
2. **Preserve User Dotfiles**: Existing `~/.zshrc` and Oh-My-Zsh installations must never be overwritten, renamed, or theme-altered when already configured.
3. **Implicit Booleans & Functions <= 15 Lines**: All functions adhere strictly to repository coding guidelines.
