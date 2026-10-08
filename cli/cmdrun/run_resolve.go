package cmdrun

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// CanResolveTarget returns true if input resolves to an executable script file.
func CanResolveTarget(input string) bool {
	target, err := ResolveRunTarget(input)

	return err == nil && target != nil
}

// ResolveRunTarget resolves a script path, inspecting direct files, candidate extensions, and shebangs.
func ResolveRunTarget(input string) (*RunTarget, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, apperror.NewValidationError("script target path cannot be empty")
	}

	normPath := normalizeToRelative(trimmed)
	if target, isFound := probeDirectTarget(trimmed, normPath); isFound {
		return target, nil
	}

	if target, isFound := probeCandidateExtensions(trimmed, normPath); isFound {
		return target, nil
	}

	return nil, apperror.NewValidationError(fmt.Sprintf("target script '%s' not found", input))
}

func normalizeToRelative(input string) string {
	clean := filepath.Clean(input)
	if !filepath.IsAbs(clean) {
		return clean
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return clean
	}

	rel, relErr := filepath.Rel(cwd, clean)
	if relErr == nil && !strings.HasPrefix(rel, "..") {
		return rel
	}

	return clean
}

func probeDirectTarget(rawInput, path string) (*RunTarget, bool) {
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		return nil, false
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		ext = inspectShebangExtension(path)
	}

	interpreter := detectInterpreterForExt(ext)

	return &RunTarget{
		RawInput:      rawInput,
		ResolvedPath:  filepath.ToSlash(path),
		Extension:     ext,
		Interpreter:   interpreter,
		IsDirectMatch: true,
	}, true
}

func probeCandidateExtensions(rawInput, basePath string) (*RunTarget, bool) {
	for _, ext := range SupportedRunExtensions {
		candidate := basePath + ext
		fi, err := os.Stat(candidate)
		if err != nil || fi.IsDir() {
			continue
		}

		interpreter := detectInterpreterForExt(ext)

		return &RunTarget{
			RawInput:      rawInput,
			ResolvedPath:  filepath.ToSlash(candidate),
			Extension:     ext,
			Interpreter:   interpreter,
			IsDirectMatch: false,
		}, true
	}

	return nil, false
}

func inspectShebangExtension(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return ""
	}

	line := strings.TrimSpace(scanner.Text())
	if !strings.HasPrefix(line, "#!") {
		return ""
	}

	return mapShebangToExt(strings.ToLower(line))
}

func mapShebangToExt(line string) string {
	if strings.Contains(line, "python") {
		return ".py"
	}
	if strings.Contains(line, "bash") || strings.Contains(line, "/sh") {
		return ".sh"
	}
	if strings.Contains(line, "node") || strings.Contains(line, "bun") {
		return ".js"
	}
	if strings.Contains(line, "pwsh") || strings.Contains(line, "powershell") {
		return ".ps1"
	}

	return ""
}

func detectInterpreterForExt(ext string) string {
	switch ext {
	case ".py":
		return "python"
	case ".ps1":
		return "powershell"
	case ".sh", ".bash":
		return "bash"
	case ".js":
		return "node"
	case ".ts":
		return "bun"
	case ".go":
		return "go"
	default:
		return "unknown"
	}
}
