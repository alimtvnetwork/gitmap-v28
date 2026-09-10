# Subtask 04: Shell Installers Multi-Tier Cleanup Hardening

## Objective
Harden all shell installation and uninstallation scripts:
1. `install-quick.sh`:
   - Replace the fragile single glob `for bad_dir in ...` with the multi-tiered cleanup engine:
     - Tier 1: Python 3 `os.listdir()` byte scanner (highest reliability, immune to glob/newline bugs).
     - Tier 2: Perl fallback (`opendir`/`readdir`).
     - Tier 3: POSIX `find` with ANSI and newline character matching.
     - Tier 4: Bash fallback with `shopt -s nullglob`.
   - Add `trap cleanup_corrupted_install_dirs EXIT INT TERM` so cleanup ALWAYS runs even on failure or abort.
   - Ensure `prompt_dir()` sends all banners, prompts, and hints strictly to `stderr` (`>&2`).
   - Fix `sanitize_install_dir` with portable ANSI escape stripping and reject any string containing control characters (`\x00`–`\x1F`).
2. `install.sh` & `gitmap/scripts/install.sh`:
   - Keep both files strictly identical.
   - Synchronize the hardened multi-tier `cleanup_corrupted_install_dirs` and `sanitize_install_dir`.
3. `uninstall-quick.sh`:
   - Add `cleanup_corrupted_install_dirs` to `uninstall-quick.sh` so running uninstaller sweeps away orphaned corrupted directories.

## Files Affected
- `install-quick.sh`
- `install.sh`
- `gitmap/scripts/install.sh`
- `uninstall-quick.sh`
