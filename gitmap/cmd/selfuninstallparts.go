package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/setup"
)

// selfDeployDir returns the directory the live binary is installed in.
// Falls back to "" if it cannot be resolved.
func selfDeployDir() string {
	self, err := os.Executable()
	if err != nil {
		return ""
	}
	resolved, err := filepath.EvalSymlinks(self)
	if err != nil {
		resolved = self
	}

	return filepath.Dir(resolved)
}

// selfDataDir returns the .gitmap data directory anchored to the binary.
func selfDataDir() string {
	deploy := selfDeployDir()
	if len(deploy) == 0 {
		return ""
	}

	return filepath.Join(deploy, constants.DefaultOutputFolder)
}

// removeDeployArtifacts deletes the gitmap binary plus any sibling
// artifacts (handoff temp copies, .old backups, completion files).
func removeDeployArtifacts(dir string) {
	if len(dir) == 0 {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgSelfUninstallSkipBin, err)

		return
	}
	for _, e := range entries {
		if !isGitmapArtifact(e.Name()) {
			continue
		}
		full := filepath.Join(dir, e.Name())
		removePathBestEffort(full)
	}
}

// isGitmapArtifact reports whether a filename belongs to the gitmap
// install (binary, handoff copies, completion outputs, .old backups).
func isGitmapArtifact(name string) bool {
	lower := strings.ToLower(name)
	if lower == "gitmap" || lower == "gitmap.exe" {
		return true
	}
	if strings.HasPrefix(lower, "gitmap-handoff-") {
		return true
	}
	if strings.HasSuffix(lower, ".old") && strings.HasPrefix(lower, "gitmap") {
		return true
	}
	if strings.HasPrefix(lower, "gitmap-completion") {
		return true
	}

	return false
}

// removeCompletionFiles wipes generated bash/zsh/fish completion files
// that gitmap may have written under the deploy dir.
func removeCompletionFiles(dir string) {
	if len(dir) == 0 {
		return
	}
	candidates := []string{"gitmap-completion.bash", "gitmap-completion.zsh", "gitmap-completion.fish"}
	for _, c := range candidates {
		removePathBestEffort(filepath.Join(dir, c))
	}
}

// removePathBestEffort deletes a file or directory and reports the
// outcome. Missing paths are silently ignored.
func removePathBestEffort(path string) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfUninstallRemove, path, err)

		return
	}
	if info.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfUninstallRemove, path, err)

		return
	}
	if info.IsDir() {
		fmt.Printf(constants.MsgSelfUninstallRemovedDir, path)

		return
	}
	fmt.Printf(constants.MsgSelfUninstallRemovedBin, path)
}

// removeProfileSnippet strips the gitmap shell-wrapper marker block
// from the user's shell profile, leaving the rest of the file intact.
func removeProfileSnippet(profile string) {
	if len(profile) == 0 {
		return
	}
	data, err := os.ReadFile(profile)
	if os.IsNotExist(err) {
		fmt.Printf(constants.MsgSelfUninstallSnippetMiss, profile)

		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfUninstallSnippetRead, profile, err)

		return
	}
	stripped, removed := stripMarkerBlock(string(data))
	if !removed {
		fmt.Printf(constants.MsgSelfUninstallSnippetMiss, profile)

		return
	}
	if writeErr := os.WriteFile(profile, []byte(stripped), 0o644); writeErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSelfUninstallSnippetWrite, profile, writeErr)

		return
	}
	fmt.Printf(constants.MsgSelfUninstallSnippetGone, profile)
}

// removeCompletionSourceLines strips the `# gitmap shell completion` comment
// and the `. '...completions.ps1'` / `source '...'` dot-source line from
// EVERY resolved profile. Without this the profile errors on the first prompt
// after uninstall because the script file is gone.
func removeCompletionSourceLines() {
	for _, p := range allProfilePaths() {
		removeCompletionFromProfile(p)
	}
}

// removeCompletionFromProfile strips gitmap completion lines from one profile.
func removeCompletionFromProfile(profile string) {
	data, err := os.ReadFile(profile)
	if err != nil {
		return
	}
	cleaned, changed := stripCompletionLines(string(data))
	if !changed {
		return
	}
	if err := os.WriteFile(profile, []byte(cleaned), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clean completion line from %s: %v\n", profile, err)

		return
	}
	fmt.Printf("  ✓ Removed completion source line from %s\n", profile)
}

// stripCompletionLines removes any line containing a gitmap completion
// source command and the preceding comment header.
func stripCompletionLines(content string) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var out strings.Builder
	removed := false
	for scanner.Scan() {
		line := scanner.Text()
		if isCompletionLine(line) {
			removed = true

			continue
		}
		out.WriteString(line)
		out.WriteString("\n")
	}
	result := out.String()
	if !strings.HasSuffix(content, "\n") {
		result = strings.TrimRight(result, "\n")
	}

	return result, removed
}

// isCompletionLine reports whether the line is a gitmap completion source
// command or the associated comment header.
func isCompletionLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "# gitmap shell completion" {
		return true
	}
	if strings.Contains(trimmed, "completions.ps1") && strings.HasPrefix(trimmed, ".") {
		return true
	}
	if strings.Contains(trimmed, "gitmap-completion") && strings.HasPrefix(trimmed, "source") {
		return true
	}

	return false
}

// allProfilePaths returns every profile file the completion installer may
// have written to. Covers PowerShell (Core + Legacy) on Windows and
// bash/zsh on Unix.
func allProfilePaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	if isWindows() {
		docs := filepath.Join(home, "Documents")

		return []string{
			filepath.Join(docs, "PowerShell", "profile.ps1"),
			filepath.Join(docs, "PowerShell", "Microsoft.PowerShell_profile.ps1"),
			filepath.Join(docs, "WindowsPowerShell", "profile.ps1"),
			filepath.Join(docs, "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1"),
		}
	}

	return []string{
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".config", "powershell", "profile.ps1"),
		filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"),
	}
}

// stripMarkerBlock removes any line range delimited by the gitmap
// shell-wrapper marker open/close lines, regardless of manager string.
func stripMarkerBlock(content string) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var out strings.Builder
	skip := false
	removed := false
	for scanner.Scan() {
		line := scanner.Text()
		if !skip && isMarkerOpen(line) {
			skip = true
			removed = true

			continue
		}
		if skip && line == setup.MarkerClose() {
			skip = false

			continue
		}
		if !skip {
			out.WriteString(line)
			out.WriteString("\n")
		}
	}
	if !strings.HasSuffix(content, "\n") {
		return strings.TrimRight(out.String(), "\n"), removed
	}

	return out.String(), removed
}

// isMarkerOpen matches the marker open line for any manager string.
func isMarkerOpen(line string) bool {
	return strings.HasPrefix(line, "# gitmap shell wrapper v") &&
		strings.Contains(line, " - managed by ")
}

// resolveProfilesForShellMode returns the deterministic list of profile
// files self-uninstall should strip the PATH snippet from for the given
// --shell-mode value. Mirrors install.sh's should_write_profile gating
// so uninstall touches exactly the same files self-install touched.
//
// Behavior:
//   - "auto" / "both": every known profile across every family (safest
//     for full removal — clears any snippet self-install may have left).
//   - singleton ("zsh"|"bash"|"pwsh"|"fish"): only that family's files.
//   - combo ("zsh+pwsh", etc.): strict union of the listed families;
//     skips ~/.profile and any unlisted family. Same contract as
//     self-install combos.
//
// Output is de-duplicated and order-stable so the printed target list
// matches what executeSelfUninstall actually touches.
func resolveProfilesForShellMode(mode string) []string {
	families := shellModeFamilies(mode)
	seen := map[string]bool{}
	var out []string
	for _, fam := range families {
		for _, p := range profilesForFamily(fam) {
			if len(p) == 0 || seen[p] {
				continue
			}
			seen[p] = true
			out = append(out, p)
		}
	}

	return out
}

// shellModeFamilies expands a --shell-mode value into the concrete shell
// families it covers. `auto` and `both` expand to every supported family;
// singletons return themselves; combos split on ShellModeComboSep.
func shellModeFamilies(mode string) []string {
	allFamilies := []string{
		constants.ShellModeZsh,
		constants.ShellModeBash,
		constants.ShellModePwsh,
		constants.ShellModeFish,
	}
	if mode == constants.ShellModeAuto || mode == constants.ShellModeBoth || len(mode) == 0 {
		return allFamilies
	}
	if strings.Contains(mode, constants.ShellModeComboSep) {
		return strings.Split(mode, constants.ShellModeComboSep)
	}

	return []string{mode}
}

// profilesForFamily returns the conventional profile files for one
// shell family on the current OS. Returns nil for unknown families
// (caller dedups, so this is safe).
func profilesForFamily(family string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	if isWindows() {
		return windowsProfilesForFamily(home, family)
	}

	return unixProfilesForFamily(home, family)
}

// unixProfilesForFamily lists the rc files for one shell family on Unix.
// pwsh covers both the legacy and Core profile names under
// ~/.config/powershell to mirror allProfilePaths().
func unixProfilesForFamily(home, family string) []string {
	switch family {
	case constants.ShellModeZsh:
		return []string{
			filepath.Join(home, ".zshrc"),
			filepath.Join(home, ".zprofile"),
		}
	case constants.ShellModeBash:
		return []string{
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".bash_profile"),
		}
	case constants.ShellModePwsh:
		return []string{
			filepath.Join(home, ".config", "powershell", "profile.ps1"),
			filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"),
		}
	case constants.ShellModeFish:
		return []string{filepath.Join(home, ".config", "fish", "config.fish")}
	}

	return nil
}

// windowsProfilesForFamily lists the rc files for one shell family on
// Windows. Only pwsh is meaningful on Windows; other families resolve
// to nil so a Linux-style `--shell-mode zsh+pwsh` invocation on Windows
// still cleans the pwsh profiles without erroring.
func windowsProfilesForFamily(home, family string) []string {
	if family != constants.ShellModePwsh {
		return nil
	}
	docs := filepath.Join(home, "Documents")

	return []string{
		filepath.Join(docs, "PowerShell", "profile.ps1"),
		filepath.Join(docs, "PowerShell", "Microsoft.PowerShell_profile.ps1"),
		filepath.Join(docs, "WindowsPowerShell", "profile.ps1"),
		filepath.Join(docs, "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1"),
	}
}
