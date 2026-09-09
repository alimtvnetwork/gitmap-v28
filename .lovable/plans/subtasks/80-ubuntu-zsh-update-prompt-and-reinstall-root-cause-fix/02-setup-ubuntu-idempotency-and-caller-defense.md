# Subtask 02: Setup Ubuntu Idempotency, Flag Integration & Caller Defense

## Objective
Refactor `setup.go` and `setup_ubuntu.go` to be 100% idempotent, non-blocking, and wire defense-in-depth callers across `selfinstall.go`, `setupverify.go`, `install.sh`, and `gitmap/scripts/install.sh`.

## Scope of Changes
1. **`gitmap/cmd/setup.go`**:
   - Add `--skip-zsh` flag in `parseSetupFlags(args)`.
   - Update `ensureZshUbuntuStep` call to pass `dryRun` and `isSkipZsh`.
2. **`gitmap/cmd/setup_ubuntu.go`**:
   - Refactor `ensureZshUbuntuStep(isDryRun, isSkipZsh bool)`.
   - If `isSkipZsh` or `isSkipZshEnv()` is true: return immediately.
   - If `zshPath, isZshInstalled := findZshBinary(); isZshInstalled`:
     - Print `✓ ZSH is already installed (<version>)`.
     - Check `isOhMyZshInstalled()`. If installed: print `✓ Oh-My-Zsh is already installed` and return immediately.
   - If missing, check `isStdinTerminal()`. If not a terminal, print notice and skip prompt.
   - Only prompt if missing and running in an interactive terminal.
3. **Caller Defense-in-Depth**:
   - `gitmap/cmd/selfinstall.go`: update `autoRunSetupAfterInstall` to pass `[]string{"--skip-zsh"}`.
   - `gitmap/cmd/setupverify.go`: update `autoRunSetupForCD` to pass `[]string{"--skip-zsh"}`.
   - `install.sh` & `gitmap/scripts/install.sh`: update line 1763/1764 to execute `GITMAP_SKIP_ZSH=1 "${bin_path}" setup --skip-zsh`.
4. **Docs**:
   - Update `gitmap/helptext/setup.md` to document `--skip-zsh`.

## Coding Guidelines Compliance
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested ifs.
- Affirmative booleans only.
