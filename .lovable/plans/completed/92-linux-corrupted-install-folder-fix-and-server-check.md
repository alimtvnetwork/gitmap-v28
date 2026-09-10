# 92-linux-corrupted-install-folder-fix-and-server-check

## Goal Description
Resolve the critical Ubuntu/Linux and macOS installation issue where interactive banners or stdout pollution captured by command substitutions (`INSTALL_DIR="$(prompt_dir)"`) or unexpanded tildes create corrupted filesystem folders in `$HOME`, `CWD`, or `/tmp` containing raw ANSI escape sequences (`\033[36m`), newlines (`\n`), prompt text (`gitmap quick installer`, `Default:`), or literal `~`.

Provide an end-to-end multi-layer solution:
1. Identify and document the exact root cause of the corrupted folder creation and explain why previous bash glob cleanups failed.
2. Build a native Go detection, asset recovery, and cleanup engine in `gitmap/cmd/` that runs automatically on Linux/macOS startup, inside `gitmap self-install`, in `gitmap doctor`, and as a standalone subcommand (`gitmap clean-corrupted`).
3. Safely handle edge cases: if the user's process is currently inside the corrupted directory, chdir out to a safe directory (`~/.local/bin` or `$HOME`) before removing or renaming; if any Gitmap binary or state file is trapped inside the corrupted folder, recover it to `~/.local/bin` before deletion.
4. Harden all shell scripts (`install-quick.sh`, `install.sh`, `gitmap/scripts/install.sh`, `uninstall-quick.sh`) with robust multi-tiered cleanup (Python 3 inode listing, POSIX find byte matching, and trap handling).
5. Enable remote Linux server verification across clusters via `gitmap server-cmd` and SSH delegation.

---

## Custom Constraints & Rules (Task-Specific)
1. **Protected Paths Invariant (Rule 1):** The cleaner MUST NEVER delete or chdir away from protected root directories: `/`, `$HOME`, `/usr`, `/usr/local`, `/usr/local/bin`, `/bin`, `/tmp`, or `.`. For literal `~` directories, strictly verify `filepath.Base(path) == "~"` and `path != homeDir`.
2. **CWD Escape Before Deletion (Rule 2):** If `os.Getwd()` is inside or equal to the corrupted directory, Gitmap MUST escape to a safe directory (`~/.local/bin` or `$HOME`) before deleting or renaming the directory.
3. **Asset Recovery Guarantee (Rule 3):** If a corrupted directory contains `gitmap`, `gitmap-cli`, or `.gitmap` state, copy it to the canonical target (`~/.local/bin/gitmap` with `0755` permissions) before unlinking.
4. **Function & File Size Caps (Rule 4):** All new Go source files must remain strictly <= 100 lines (hard limit <= 200 lines). Every function must remain strictly <= 15 lines.
5. **Coding Guideline Conformance (Rule 5):** All booleans must use `is` or `has` prefixes only (`isCorrupted`, `hasFiles`, `isProtected`). Zero nested `if` statements. Single return error envelopes with `apperror.Wrap`.

---

## Subtasks Decomposition
- `01-go-corrupted-dirs-detector-and-recovery.md`: Implement detection, pattern matching, protected path checks, and asset recovery in `gitmap/cmd/corrupted_dirs_detector.go` and `gitmap/cmd/corrupted_dirs_recovery.go`.
- `02-go-corrupted-dirs-cleaner-and-subcommand.md`: Implement safe CWD escaping and deletion in `gitmap/cmd/corrupted_dirs_cleaner.go` and register `gitmap clean-corrupted` in `gitmap/cmd/corrupted_dirs_cmd.go`.
- `03-go-cli-lifecycle-hooks-and-doctor.md`: Wire automated cleanup into `Run()` in `gitmap/cmd/root.go`, `selfinstall.go`, and add doctor probe in `gitmap/cmd/doctor.go` / `doctor_run.go`.
- `04-shell-installers-multi-tier-cleanup.md`: Harden `install-quick.sh`, `install.sh`, `gitmap/scripts/install.sh`, and `uninstall-quick.sh` with robust Python/find cleanup and EXIT traps.
- `05-unit-tests-and-ci-verification.md`: Author comprehensive unit tests in `gitmap/cmd/corrupted_dirs_test.go`, test shell scripts, and run CI quality gates.
