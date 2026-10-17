package cmdconflict

import (
	"fmt"
	"os"
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
		if path, err := exec.LookPath(ide.bin); err == nil && path != "" {
			return ide.bin, ide.name
		}
		if runtime.GOOS == "windows" {
			for _, ext := range []string{".cmd", ".exe"} {
				if path, err := exec.LookPath(ide.bin + ext); err == nil && path != "" {
					return ide.bin, ide.name
				}
			}
		}
	}

	return "", ""
}

// LaunchIDEForRepo launches the detected IDE for the given repository path and files.
// If files are provided, it passes them with the -r (reuse window) flag.
func LaunchIDEForRepo(repoPath string, files []string) error {
	bin, _ := DetectInstalledIDE()
	if bin == "" {
		return fmt.Errorf("no supported IDE found on PATH (checked: code, cursor, windsurf, codium)")
	}

	var args []string
	if len(files) > 0 {
		args = append(args, "-r")
		for _, f := range files {
			if !filepath.IsAbs(f) && repoPath != "" {
				args = append(args, filepath.Join(repoPath, f))
			} else {
				args = append(args, f)
			}
		}
	} else if repoPath != "" {
		args = append(args, repoPath)
	} else {
		args = append(args, ".")
	}

	cmd := exec.Command(bin, args...)
	if repoPath != "" {
		cmd.Dir = repoPath
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start()
}

// FormatMergeEditorGuide returns clear formatted instructions on how to use
// the VS Code 3-Way Merge Editor.
func FormatMergeEditorGuide(repoPath string, conflictedFiles []string) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("=================================================================\n")
	sb.WriteString("  VS Code 3-Way Merge Editor Walkthrough\n")
	sb.WriteString("=================================================================\n")
	if len(conflictedFiles) > 0 {
		sb.WriteString(fmt.Sprintf("  Conflicted Files (%d):\n", len(conflictedFiles)))
		for _, f := range conflictedFiles {
			sb.WriteString(fmt.Sprintf("    • %s\n", f))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("  1. Press Ctrl+Shift+G (Cmd+Shift+G on macOS) to open Source Control\n")
	sb.WriteString("  2. Click a conflicted file, then click 'Resolve in Merge Editor'\n")
	sb.WriteString("  3. Review the 3-Way split:\n")
	sb.WriteString("     - Incoming (Remote): Changes from incoming branch\n")
	sb.WriteString("     - Current (Local):   Your current branch changes\n")
	sb.WriteString("     - Result (Bottom):   Final merged output\n")
	sb.WriteString("  4. Click 'Accept Incoming', 'Accept Current', or combine changes\n")
	sb.WriteString("  5. Save the file (Ctrl+S / Cmd+S) — Git marks it resolved\n")
	sb.WriteString("  6. Stage & commit with GitMap:\n")
	sb.WriteString("     gitmap cpf \"resolve merge conflicts\"\n")
	sb.WriteString("=================================================================\n\n")

	return sb.String()
}
