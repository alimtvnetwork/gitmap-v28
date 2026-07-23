package completion

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// InstallCDFunction writes the gitmap/gcd shell wrapper to user profiles.
func InstallCDFunction(shell string) error {
	snippet := cdSnippet(shell)
	if len(snippet) == 0 {
		return fmt.Errorf(constants.ErrCompUnknownShell, shell)
	}
	if shell == constants.ShellPowerShell {
		return installPowerShellCDFunction(snippet)
	}

	return appendCDFunctions(snippet, cdProfilePaths(shell))
}

func installPowerShellCDFunction(snippet string) error {
	if err := appendCDFunctions(snippet, cdProfilePaths(constants.ShellPowerShell)); err != nil {
		return err
	}
	if runtime.GOOS != constants.OSWindows {
		return nil
	}

	return installPowerShellCommandShim()
}

func installPowerShellCommandShim() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(exe)
	body := renderPowerShellCommandShim(dir)

	return os.WriteFile(filepath.Join(dir, constants.PowerShellShimFile), []byte(body), 0o755)
}

func renderPowerShellCommandShim(dir string) string {
	escaped := strings.ReplaceAll(dir,
		constants.PowerShellSingleQuote,
		constants.PowerShellEscapedQuote)

	return fmt.Sprintf(constants.PowerShellShimTemplateFmt, escaped)
}

// cdSnippet returns the gcd function body for the given shell.
func cdSnippet(shell string) string {
	switch shell {
	case constants.ShellPowerShell:
		return constants.CDFuncPowerShell
	case constants.ShellBash:
		return constants.CDFuncBash
	case constants.ShellZsh:
		return constants.CDFuncZsh
	default:
		return ""
	}
}

// cdProfilePaths returns all profile paths to write the cd function to.
func cdProfilePaths(shell string) []string {
	switch shell {
	case constants.ShellPowerShell:
		return resolvePowerShellProfilePaths()
	case constants.ShellBash:
		home, _ := os.UserHomeDir()
		return []string{filepath.Join(home, ".bashrc")}
	default:
		home, _ := os.UserHomeDir()
		return []string{filepath.Join(home, ".zshrc")}
	}
}

// appendCDFunctions appends the managed wrapper to every resolved profile.
func appendCDFunctions(snippet string, profilePaths []string) error {
	for _, profilePath := range profilePaths {
		if err := appendCDFunction(snippet, profilePath); err != nil {
			return err
		}
	}

	return nil
}

// appendCDFunction appends the gitmap command wrapper to the profile if not present.
func appendCDFunction(snippet, profilePath string) error {
	if err := os.MkdirAll(filepath.Dir(profilePath), 0o755); err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}

	existing, err := os.ReadFile(profilePath)
	if err == nil {
		text := string(existing)
		if hasCurrentCDFunction(text) {
			next := reconcileCDFunction(text, snippet)
			if next == text {
				fmt.Fprintf(os.Stderr, constants.MsgCDFuncAlready)

				return nil
			}

			if writeErr := os.WriteFile(profilePath, []byte(next), 0o644); writeErr != nil {
				return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, writeErr)
			}
			fmt.Fprintf(os.Stderr, constants.MsgCDFuncInstalled)

			return nil
		}
	}

	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\n%s\n%s\n", constants.CDFuncMarker, snippet)
	if err != nil {
		return fmt.Errorf(constants.ErrCompProfileWrite, profilePath, err)
	}

	fmt.Fprintf(os.Stderr, constants.MsgCDFuncInstalled)

	return nil
}

func hasCurrentCDFunction(text string) bool {
	return strings.Contains(text, constants.CDFuncMarker) &&
		hasCDFunctionEnd(text)
}

func hasCDFunctionEnd(text string) bool {
	return strings.Contains(text, constants.CDFuncMarkerEnd) ||
		strings.Contains(text, constants.CDFuncMarkerEndLegacy)
}

func replaceCDFunction(text, snippet string) string {
	start := strings.Index(text, constants.CDFuncMarker)
	if start < 0 {
		return text
	}
	end, marker := findCDFunctionEnd(text[start:])
	if end < 0 {
		return text
	}
	end += start + len(marker)
	replacement := constants.CDFuncMarker + "\n" + snippet

	return text[:start] + replacement + text[end:]
}

func reconcileCDFunction(text, snippet string) string {
	start := strings.Index(text, constants.CDFuncMarker)
	if start < 0 {
		return text
	}
	endRel, marker := findCDFunctionEnd(text[start:])
	if endRel < 0 {
		return text
	}
	blockEnd := start + endRel + len(marker)

	desired := constants.CDFuncMarker + "\n" + snippet
	existing := text[start:blockEnd]
	after := text[blockEnd:]
	trailingSignificant := strings.TrimSpace(after) != ""

	if existing == desired && !trailingSignificant {
		return text
	}
	if !trailingSignificant {
		return text[:start] + desired + after
	}

	return appendBlock(text[:start]+after, desired)
}

func moveCDFunctionToEnd(text, snippet string) string {
	without := removeCDFunction(text)
	block := constants.CDFuncMarker + "\n" + snippet

	return appendBlock(without, block)
}

func removeCDFunction(text string) string {
	start := strings.Index(text, constants.CDFuncMarker)
	if start < 0 {
		return text
	}
	end, marker := findCDFunctionEnd(text[start:])
	if end < 0 {
		return text
	}

	return text[:start] + text[start+end+len(marker):]
}

func appendBlock(text, block string) string {
	trimmed := strings.TrimRight(text, "\r\n")
	if len(trimmed) == 0 {
		return block + "\n"
	}

	return trimmed + "\n\n" + block + "\n"
}

func findCDFunctionEnd(text string) (int, string) {
	current := strings.Index(text, constants.CDFuncMarkerEnd)
	legacy := strings.Index(text, constants.CDFuncMarkerEndLegacy)
	if current >= 0 && (legacy < 0 || current < legacy) {
		return current, constants.CDFuncMarkerEnd
	}
	return legacy, constants.CDFuncMarkerEndLegacy
}
