package cmdconflict

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type ideCandidate struct {
	bin  string
	name string
}

var supportedIDEs = []ideCandidate{
	{bin: "code", name: "VS Code"},
	{bin: "cursor", name: "Cursor"},
	{bin: "windsurf", name: "Windsurf"},
	{bin: "codium", name: "VSCodium"},
}

// DetectInstalledIDE checks PATH for code, cursor, windsurf, or codium.
// Returns the binary name and friendly display name, or empty strings if none found.
func DetectInstalledIDE() (string, string) {
	for _, ide := range supportedIDEs {
		path, err := exec.LookPath(ide.bin)
		if err == nil && path != "" {
			return ide.bin, ide.name
		}
		if winPath := findWindowsIDEBinary(ide.bin); winPath != "" {
			return ide.bin, ide.name
		}
	}

	return "", ""
}

func findWindowsIDEBinary(bin string) string {
	if runtime.GOOS != "windows" {
		return ""
	}
	for _, ext := range []string{".cmd", ".exe"} {
		path, err := exec.LookPath(bin + ext)
		if err == nil && path != "" {
			return path
		}
	}
	return ""
}

// LaunchIDEForRepo launches the detected IDE for the given repository path and files.
// If files are provided, it passes them with the -r (reuse window) flag.
func LaunchIDEForRepo(repoPath string, files []string) error {
	bin, _ := DetectInstalledIDE()
	if bin == "" {
		return fmt.Errorf("no supported IDE found on PATH (checked: code, cursor, windsurf, codium)")
	}

	args := resolveIDEFileArgs(repoPath, files)
	cmd := exec.Command(bin, args...)
	if repoPath != "" {
		cmd.Dir = repoPath
	}

	return cmd.Start()
}

func resolveIDEFileArgs(repoPath string, files []string) []string {
	if len(files) == 0 {
		return resolveIDERepoArgs(repoPath)
	}
	args := []string{"-r"}
	for _, f := range files {
		args = append(args, resolveTargetFilePath(repoPath, f))
	}
	return args
}

func resolveTargetFilePath(repoPath, f string) string {
	if !filepath.IsAbs(f) && repoPath != "" {
		return filepath.Join(repoPath, f)
	}
	return f
}

func resolveIDERepoArgs(repoPath string) []string {
	if repoPath != "" {
		return []string{repoPath}
	}
	return []string{"."}
}

// FormatMergeEditorGuide returns a clean formatted walkthrough for the VS Code 3-Way Merge Editor.
func FormatMergeEditorGuide(repoPath string, conflictedFiles []string) string {
	var sb strings.Builder
	sb.WriteString("\n--------------------------------------------------------------------------------\n")
	sb.WriteString("  VS Code 3-Way Merge Editor — Resolution Guide\n")
	sb.WriteString("--------------------------------------------------------------------------------\n")
	sb.WriteString("  1. Press Ctrl+Shift+G (or Cmd+Shift+G on macOS) to open the Source Control view.\n")
	sb.WriteString("  2. Under 'Merge Changes', click on any conflicted file.\n")
	sb.WriteString("  3. Click 'Resolve in Merge Editor' in the bottom-right of the editor tab.\n")
	sb.WriteString("  4. Check 'Accept Incoming' (remote) or 'Accept Current' (local) for each conflict.\n")
	sb.WriteString("  5. Click 'Complete Merge' at the bottom-right and save the file (Ctrl+S).\n")
	sb.WriteString("  6. In GitMap, commit and push your resolution:\n")
	sb.WriteString("     gitmap cpf \"resolve merge conflicts\"\n")
	sb.WriteString("--------------------------------------------------------------------------------\n\n")

	return sb.String()
}
