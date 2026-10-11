package completion

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// Install writes the completion script and adds a source line to the profile.
func Install(shell string) error {
	script, err := Generate(shell)
	if err != nil {
		return err
	}

	scriptPath, profilePath := resolvePaths(shell)

	return writeAndSource(script, scriptPath, profilePath, shell)
}

// resolvePaths returns the script file path and profile path for the shell.
func resolvePaths(shell string) (string, string) {
	switch shell {
	case constants.ShellPowerShell:
		return resolvePowerShellPaths()
	case constants.ShellBash:
		return resolveBashPaths()
	default:
		return resolveZshPaths()
	}
}

// resolveBashPaths returns paths for Bash completion.
func resolveBashPaths() (string, string) {
	home, _ := os.UserHomeDir()
	scriptPath := filepath.Join(home, ".local", "share", constants.CompDirName, constants.CompFileBash)
	profile := filepath.Join(home, ".bashrc")

	return scriptPath, profile
}

// resolveZshPaths returns paths for Zsh completion.
func resolveZshPaths() (string, string) {
	home, _ := os.UserHomeDir()
	scriptPath := filepath.Join(home, ".local", "share", constants.CompDirName, constants.CompFileZsh)
	profile := filepath.Join(home, ".zshrc")

	return scriptPath, profile
}

// writeAndSource writes the script file and adds a source line to the profile.
func writeAndSource(script, scriptPath, profilePath, shell string) error {
	if err := writeScriptFile(scriptPath, script); err != nil {
		return err
	}

	profilePaths := []string{profilePath}
	if shell == constants.ShellPowerShell {
		profilePaths = resolvePowerShellProfilePaths()
	}

	return addSourceLines(scriptPath, profilePaths, shell)
}

// addSourceLines appends the source command to every resolved profile.
func addSourceLines(scriptPath string, profilePaths []string, shell string) error {
	for _, profilePath := range profilePaths {
		if err := addSourceLine(scriptPath, profilePath, shell); err != nil {
			return err
		}
	}

	return nil
}

// writeScriptFile creates directories and writes the completion script.
func writeScriptFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

const (
	compMarkerStart     = "# >>> gitmap shell completion & predictive suggestions >>>"
	compMarkerEnd       = "# <<< gitmap shell completion & predictive suggestions <<<"
	compLegacyMarker    = "# gitmap shell completion"
	compUnixMarkerStart = "# >>> gitmap shell completion >>>"
	compUnixMarkerEnd   = "# <<< gitmap shell completion <<<"
)

// addSourceLine updates or appends the shell completion and predictive suggestions snippet.
func addSourceLine(scriptPath, profilePath, shell string) error {
	snippet := buildSourceLine(scriptPath, shell)
	if err := os.MkdirAll(filepath.Dir(profilePath), 0o755); err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}

	existingBytes, err := os.ReadFile(profilePath)
	existing := ""
	if err == nil {
		existing = string(existingBytes)
	}

	updated, hasChanged := reconcileProfileSnippet(existing, snippet, shell)
	if !hasChanged {
		fmt.Fprintf(os.Stderr, constants.MsgCompAlreadyDone, shell)

		return nil
	}

	if err := os.WriteFile(profilePath, []byte(updated), 0o644); err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCompProfileWrite, profilePath)

	return nil
}

// buildSourceLine returns the shell-appropriate source command.
func buildSourceLine(scriptPath, shell string) string {
	if shell == constants.ShellPowerShell {
		return toNativeEOL(buildPowerShellSnippet(scriptPath), resolveProfileEOL("", shell))
	}

	return buildUnixSnippet(scriptPath)
}

func buildPowerShellSnippet(scriptPath string) string {
	return fmt.Sprintf(`%s
if (-not $global:__gitmap_suggestions_configured) {
    $global:__gitmap_suggestions_configured = $true
    Import-Module PSReadLine -MinimumVersion 2.2.0 -ErrorAction SilentlyContinue
    if (Get-Command Set-PSReadLineOption -ErrorAction SilentlyContinue) {
        try {
            Set-PSReadLineOption -PredictionSource HistoryAndPlugin -ErrorAction SilentlyContinue
        } catch {
            try {
                Set-PSReadLineOption -PredictionSource History -ErrorAction SilentlyContinue
            } catch {}
        }
        try {
            Set-PSReadLineOption -PredictionViewStyle ListView -ErrorAction SilentlyContinue
        } catch {}
    }
}
if (Test-Path -LiteralPath '%s') {
    . '%s'
}
%s`, compMarkerStart, scriptPath, scriptPath, compMarkerEnd)
}

func buildUnixSnippet(scriptPath string) string {
	return fmt.Sprintf("%s\n[[ -f '%s' ]] && source '%s'\n%s",
		compUnixMarkerStart, scriptPath, scriptPath, compUnixMarkerEnd)
}

// reconcileProfileSnippet reconciles existing profile content with the latest snippet.
func reconcileProfileSnippet(existing, snippet, shell string) (string, bool) {
	eol := resolveProfileEOL(existing, shell)
	normExisting := strings.ReplaceAll(existing, "\r\n", "\n")
	normSnippet := strings.ReplaceAll(snippet, "\r\n", "\n")

	if shell == constants.ShellPowerShell {
		return reconcilePowerShellProfile(normExisting, normSnippet, existing, eol)
	}

	return reconcileUnixProfile(normExisting, normSnippet, existing, eol)
}

func resolveProfileEOL(existing, shell string) string {
	if strings.Contains(existing, "\r\n") || (runtime.GOOS == "windows" && shell == constants.ShellPowerShell) {
		return "\r\n"
	}

	return "\n"
}

func reconcilePowerShellProfile(normExisting, normSnippet, rawExisting, eol string) (string, bool) {
	if strings.Contains(normExisting, compMarkerStart) && strings.Contains(normExisting, compMarkerEnd) {
		replaced := replaceExistingBlock(normExisting, normSnippet, compMarkerStart, compMarkerEnd)
		cleaned := collapseBlankLines(stripLegacyPowerShellBlocks(replaced))
		final := toNativeEOL(cleaned, eol)

		return final, final != rawExisting
	}

	cleaned := collapseBlankLines(stripLegacyPowerShellBlocks(normExisting))
	appended := appendSnippetToProfile(cleaned, normSnippet)
	final := toNativeEOL(appended, eol)

	return final, final != rawExisting
}
