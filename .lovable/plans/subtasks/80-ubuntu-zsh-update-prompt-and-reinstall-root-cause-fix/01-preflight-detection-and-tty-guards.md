# Subtask 01: Pre-Flight ZSH/OMZ Detection Engine & Probe Registration

## Objective
Implement modular pre-flight detection for ZSH and Oh-My-Zsh in `gitmap/cmd/setup_ubuntu_detect.go` and explicitly register `constants.ToolZsh` in `gitmap/cmd/installprobe.go`.

## Scope of Changes
1. **`gitmap/cmd/setup_ubuntu_detect.go`** (NEW):
   - `findZshBinary() (string, bool)`: Checks `exec.LookPath("zsh")` and `/bin/zsh`, `/usr/bin/zsh`, `/usr/local/bin/zsh`.
   - `getZshVersion(zshPath string) string`: Runs `zsh --version` and parses clean version string (e.g. `5.8.1`, `5.9`).
   - `isOhMyZshInstalled() bool`: Checks `$ZSH` env and `~/.oh-my-zsh` directory existence.
   - `isStdinTerminal() bool`: Checks `os.Stdin.Stat()` for `os.ModeCharDevice`.
   - `isSkipZshEnv() bool`: Checks if `os.Getenv(constants.EnvGitmapSkipZsh) == "1"`.
   - Provide package-level test seam variables for mock testing without root/Ubuntu.
2. **`gitmap/constants/constants_cli.go`**:
   - Add `FlagSkipZsh = "skip-zsh"`
   - Add `FlagDescSkipZsh = "Skip ZSH and Oh-My-Zsh setup on Ubuntu"`
   - Add `EnvGitmapSkipZsh = "GITMAP_SKIP_ZSH"`
3. **`gitmap/cmd/installprobe.go`**:
   - Register `constants.ToolZsh: {bins: []string{"zsh"}, args: []string{"--version"}}` in `toolProbeMap`.

## Coding Guidelines Compliance
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested ifs.
- Affirmative booleans only.
- Strict Unix LF line endings.
