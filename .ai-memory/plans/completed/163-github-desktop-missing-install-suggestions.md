# Plan 163: GitHub Desktop Missing Install Suggestions & CLI Parity Suite

## Overview
Autonomously implement missing GitHub Desktop CLI detection, formatted installation suggestions with exact commands and multi-platform options, `--install` / `-i` flag support, and clean validation error handling to prevent raw execution stack traces across `gitmap github-desktop` (alias `gd`) and `gitmap desktop-sync` (alias `ds`).

## Key Goals
1. **Remediation Suggestion Engine (`cli/desktop/suggest.go`)**:
   - When `desktop.ResolveCLI()` returns empty, print formatted installation instructions with exact commands and options tailored to the runtime OS:
     - Option 1 (Recommended): `gitmap install github-desktop` (or `gitmap in gd`)
     - Option 2 (Native Package Managers):
       - Linux/Ubuntu: APT via shiftkey repo (`wget -qO - https://mirror.mwt.me/ghd/gpgkey | sudo tee ... && sudo apt update && sudo apt install -y github-desktop`) and Snap (`sudo snap install github-desktop --beta`)
       - Windows: `winget install --id GitHub.GitHubDesktop` and `choco install github-desktop`
       - macOS: `brew install --cask github-desktop`
     - Option 3 (Official GUI Installer): `https://desktop.github.com`
2. **Clean Validation Error & Zero Stack Traces**:
   - Replace `apperror.NewSimple(constants.MsgDesktopNotFound, "E9000")` with `apperror.NewValidationError(constants.MsgDesktopNotFound)`.
   - Because `appErr.Type == ErrorTypeValidation`, `cmd/root.go` emits `validation failed:` with zero raw execution stack traces.
3. **`--install` / `-i` Auto-Remediation Flag**:
   - Support `gitmap github-desktop --install` (and `gitmap gd -i` / `gitmap ds --install`).
   - When `--install` or `-i` is passed and the tool is missing, invoke `gitmap install github-desktop`, re-probe the CLI, and proceed.
4. **Tool Alias Parity in `cmdinstall` (`cli/cmdinstall/`)**:
   - Add aliases `"github-desktop"`, `"gh-desktop"`, and `"gd"` to `cmdinstall` so `gitmap install gd` and `gitmap in gh-desktop` resolve to `constants.ToolGitHubDesktop`.
5. **Documentation & Help Text Parity**:
   - Update `cli/helptext/github-desktop.md` and `cli/helptext/desktop-sync.md` documenting the `--install` flag and remediation options.
6. **Coding Guidelines Invariants**:
   - Functions strictly <= 15 lines (target <= 8 lines).
   - Affirmative booleans only (`is*`, `has*`).
   - Zero nested if statements (guard clauses and early returns).
   - Strict Unix LF line endings.

## Custom Rules
1. Functions strictly <= 15 lines (target <= 8 lines).
2. Affirmative booleans only (isRepo, hasInstall, isFound). No negative booleans like `isNonGitRepo`.
3. Universal AppError wrapping on all errors.
4. Absolute ban on `go test` and `go build` during routine turns.

> **Task Origin**: User reported `gitmap github-desktop` emitting `execute failed: [E9000:EXECUTION]` and raw stack trace when GitHub Desktop is missing. Requested suggesting install with exact commands and options similar to other commands.
> **Total Loops**: 1 continuous loop.

## Consolidated Subtasks Detail

# Subtask 163-01: Suggestion Engine, Alias Normalization & Clean Validation Error

## Target Files
- `cli/desktop/suggest.go` (NEW)
- `cli/desktop/suggest_test.go` (NEW)
- `cli/desktop/desktop.go`
- `cli/cmdinstall/install_aliases.go` (NEW)
- `cli/cmdinstall/install_packages.go`
- `cli/constants/constants_messages.go`

## Implemented Features
1. Created `cli/desktop/suggest.go` implementing `PrintInstallSuggestions()`, `FormatInstallSuggestions(targetOS string)`, `resolveNativeCommands(targetOS string)`, and `NewMissingCLIError()`.
2. Created `cli/desktop/suggest_test.go` with full unit test coverage across Windows, Linux, Darwin, and unknown OSes.
3. Updated `cli/desktop/desktop.go` `AddRepos` to invoke `PrintInstallSuggestions()` when the CLI is missing.
4. Created `cli/cmdinstall/install_aliases.go` registering aliases `"github-desktop"`, `"gh-desktop"`, and `"gd"` for `constants.ToolGitHubDesktop`.
5. Updated `cli/cmdinstall/install_packages.go` `resolveToolAlias` to check `desktopToolAliases`.
6. Added `FlagGHDesktopInstall = "--install"` and `FlagGHDesktopInstallShort = "-i"` to `cli/constants/constants_messages.go`.

# Subtask 163-02: Command Wiring, --install Flag Support & Help Text Parity

## Target Files
- `cli/cmd/githubdesktop.go`
- `cli/cmd/githubdesktop_test.go` (NEW)
- `cli/cmd/desktopsync.go`
- `cli/cmd/roottooling.go`
- `cli/helptext/github-desktop.md`
- `cli/helptext/desktop-sync.md`

## Implemented Features
1. Updated `cli/cmd/githubdesktop.go`:
   - Added `hasGHDesktopInstallFlag(args []string)` and `ensureGHDesktopInstalled(args []string)` to support `--install` / `-i`.
   - Refactored `registerGHDesktop(target string)` to call `desktop.PrintInstallSuggestions()` and return `desktop.NewMissingCLIError()`.
   - Renamed `isNonGitRepo` to affirmative `isRepo := isGitRepo(target)`.
2. Updated `cli/cmd/desktopsync.go`:
   - Updated `runDesktopSync(args []string)` to accept arguments and invoke `ensureGHDesktopInstalled(args)`.
   - Updated `syncToDesktop` when missing to call `desktop.PrintInstallSuggestions()` and return `desktop.NewMissingCLIError()`.
3. Updated `cli/cmd/roottooling.go` to pass `argsTail()` into `runDesktopSync`.
4. Created `cli/cmd/githubdesktop_test.go` verifying flag parsing, target resolution, subcommand matching, and validation error generation.
5. Updated `cli/helptext/github-desktop.md` and `cli/helptext/desktop-sync.md` documenting the `-i, --install` flag and missing CLI remediation.
