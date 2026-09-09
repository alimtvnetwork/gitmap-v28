# Subtask 03: Unit Tests, Quality Gate Verification & Release

## Objective
Author comprehensive unit tests in `gitmap/cmd/setup_ubuntu_test.go`, run all linters and tests, and execute release orchestrator.

## Scope of Changes
1. **`gitmap/cmd/setup_ubuntu_test.go`** (NEW):
   - `TestIsZshInstalled`: Mock `exec.LookPath` found / not found.
   - `TestIsOhMyZshInstalled`: Mock home directory and directory existence.
   - `TestEnsureZshUbuntuStep_AlreadyInstalled_Idempotent`: Verify zero commands executed and zero prompts when already installed.
   - `TestEnsureZshUbuntuStep_DryRun`: Verify dry run behavior.
   - `TestEnsureZshUbuntuStep_NonInteractive_SkipsPrompt`: Verify non-interactive terminal skips prompt.
   - `TestEnsureZshUbuntuStep_SkipFlag`: Verify `--skip-zsh` and `GITMAP_SKIP_ZSH=1` bypass.
   - `TestConfigureZshTheme_Idempotent`: Verify theme is untouched if already configured.
2. **Quality Gate Verification**:
   - `go test ./constants/... ./helptext/... ./cmd/... -count=1`
   - `python linter-scripts/check-nested-ifs.py --changed-only`
   - `python linter-scripts/check-enum-and-boolean.py --changed-only`
   - `python linter-scripts/check-error-management.py --changed-only`
3. **Release Execution**:
   - `python 03-ai-scripts/29-release-orchestrator.py --tier patch --scope "Fix Ubuntu ZSH update prompt and prevent unwanted re-install during update"`
