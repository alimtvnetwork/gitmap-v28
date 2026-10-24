package cmdconflict

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict/enums"
)

// DetectInstalledIDE checks PATH for code, cursor, windsurf, or codium.
// Returns the binary name and friendly display name, or empty strings if none found.
func DetectInstalledIDE() (string, string) {
	for _, tool := range enums.AllConflictToolTypes() {
		if bin, name, hasFound := checkIDETool(tool); hasFound {
			return bin, name
		}
	}

	return "", ""
}

func checkIDETool(tool enums.ConflictToolType) (string, string, bool) {
	bin := tool.BinaryName()
	name := tool.DisplayName()
	path, err := exec.LookPath(bin)
	if err == nil && path != "" {
		return bin, name, true
	}

	if winPath := findWindowsIDEBinary(bin); winPath != "" {
		return bin, name, true
	}

	return "", "", false
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
		return apperror.NewNotFoundError(constants.ErrNoSupportedIDE)
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

	args := []string{constants.FlagReuseWindow}

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
	writeGuideHeader(&sb)
	writeGuideSteps(&sb)
	sb.WriteString(constants.GuideDivider + "\n")

	return sb.String()
}

func writeGuideHeader(sb *strings.Builder) {
	sb.WriteString("\n" + constants.GuideDivider)
	sb.WriteString(constants.GuideTitle)
	sb.WriteString(constants.GuideDivider)
}

func writeGuideSteps(sb *strings.Builder) {
	sb.WriteString(constants.GuideStep1)
	sb.WriteString(constants.GuideStep2)
	sb.WriteString(constants.GuideStep3)
	sb.WriteString(constants.GuideStep4)
	sb.WriteString(constants.GuideStep5)
	sb.WriteString(constants.GuideStep6)
	sb.WriteString(constants.GuideStep6Sub)
}
