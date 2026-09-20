package cmdautomation

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// LocateOptions specifies parameters for finding developer tools.
type LocateOptions struct {
	Target string
	IsJSON bool
	IsCmd  bool
}

// LocateResult captures the output of locating developer tool binaries.
type LocateResult struct {
	Target     string `json:"target"`
	Path       string `json:"path"`
	HasFound   bool   `json:"hasFound"`
	DurationMs int64  `json:"durationMs"`
	Command    string `json:"command,omitempty"`
}

// RunLocate searches for developer tools with vswhere and known path fast paths.
func RunLocate(opts LocateOptions) (LocateResult, *apperror.AppError) {
	start := time.Now()
	target := resolveLocateTarget(opts.Target)
	resolvedPath, hasPath := findToolPath(target)

	durMs := time.Since(start).Milliseconds()
	res := buildLocateResult(target, resolvedPath, hasPath, durMs, opts.IsCmd)
	if !hasPath {
		return res, apperror.NewNotFound("tool_locate", "E_TOOL_NOT_FOUND", target)
	}

	return res, nil
}

func resolveLocateTarget(target string) string {
	clean := strings.TrimSpace(target)
	if len(clean) == 0 || clean == "vcvars" || clean == "vcvarsall" {
		return "vcvarsall.bat"
	}

	return clean
}

func buildLocateResult(target, path string, hasFound bool, dur int64, isCmd bool) LocateResult {
	cmdStr := ""
	if hasFound && isCmd {
		cmdStr = formatBatchInitCommand(path)
	}

	return LocateResult{
		Target:     target,
		Path:       path,
		HasFound:   hasFound,
		DurationMs: dur,
		Command:    cmdStr,
	}
}

func formatBatchInitCommand(path string) string {
	if strings.HasSuffix(strings.ToLower(path), "vcvarsall.bat") {
		return fmt.Sprintf("cmd.exe /c \"call \"%s\" x64 && set\"", path)
	}

	return path
}

func findToolPath(target string) (string, bool) {
	if isVcvarsTarget(target) {
		return locateVcvarsFastPath()
	}

	return locateGenericTool(target)
}

func isVcvarsTarget(target string) bool {
	lower := strings.ToLower(target)
	return lower == "vcvarsall.bat" || lower == "vcvars" || lower == "vcvars64.bat"
}

func locateVcvarsFastPath() (string, bool) {
	if runtime.GOOS != "windows" {
		return "", false
	}

	if path, ok := queryVsWhereForVcvars(); ok {
		return path, true
	}

	return scanKnownVcvarsPaths()
}

func queryVsWhereForVcvars() (string, bool) {
	vswherePath := findVsWhereExe()
	if len(vswherePath) == 0 {
		return "", false
	}

	args := []string{
		"-latest", "-products", "*",
		"-requires", "Microsoft.VisualStudio.Component.VC.Tools.x86.x64",
		"-find", "**/vcvarsall.bat",
	}
	out, err := exec.Command(vswherePath, args...).Output()
	if err != nil || len(out) == 0 {
		return "", false
	}

	firstLine := strings.TrimSpace(strings.Split(string(out), "\r\n")[0])
	if pathExists(firstLine) {
		return firstLine, true
	}

	return "", false
}

func findVsWhereExe() string {
	prog86 := os.Getenv("ProgramFiles(x86)")
	if len(prog86) == 0 {
		prog86 = `C:\Program Files (x86)`
	}

	candidate := filepath.Join(prog86, "Microsoft Visual Studio", "Installer", "vswhere.exe")
	if pathExists(candidate) {
		return candidate
	}

	prog := os.Getenv("ProgramFiles")
	if len(prog) == 0 {
		prog = `C:\Program Files`
	}

	candidate64 := filepath.Join(prog, "Microsoft Visual Studio", "Installer", "vswhere.exe")
	if pathExists(candidate64) {
		return candidate64
	}

	return ""
}

func scanKnownVcvarsPaths() (string, bool) {
	candidates := buildKnownVcvarsCandidates()
	for _, c := range candidates {
		if pathExists(c) {
			return c, true
		}
	}

	return searchBoundedVcvars()
}

func buildKnownVcvarsCandidates() []string {
	return []string{
		`C:\Program Files\Microsoft Visual Studio\2022\Community\VC\Auxiliary\Build\vcvarsall.bat`,
		`C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC\Auxiliary\Build\vcvarsall.bat`,
		`C:\Program Files\Microsoft Visual Studio\2022\Professional\VC\Auxiliary\Build\vcvarsall.bat`,
		`C:\Program Files\Microsoft Visual Studio\2022\BuildTools\VC\Auxiliary\Build\vcvarsall.bat`,
		`C:\Program Files (x86)\Microsoft Visual Studio\2019\Community\VC\Auxiliary\Build\vcvarsall.bat`,
		`C:\Program Files (x86)\Microsoft Visual Studio\2019\Enterprise\VC\Auxiliary\Build\vcvarsall.bat`,
		`C:\Program Files (x86)\Microsoft Visual Studio\2019\Professional\VC\Auxiliary\Build\vcvarsall.bat`,
		`C:\Program Files (x86)\Microsoft Visual Studio\2019\BuildTools\VC\Auxiliary\Build\vcvarsall.bat`,
		`C:\BuildTools\VC\Auxiliary\Build\vcvarsall.bat`,
	}
}

func searchBoundedVcvars() (string, bool) {
	roots := []string{
		`C:\Program Files\Microsoft Visual Studio`,
		`C:\Program Files (x86)\Microsoft Visual Studio`,
		`C:\BuildTools`,
	}
	for _, root := range roots {
		if path, ok := scanRootForFile(root, "vcvarsall.bat", 4); ok {
			return path, true
		}
	}

	return "", false
}

func scanRootForFile(root, filename string, maxDepth int) (string, bool) {
	if !dirExists(root) {
		return "", false
	}

	var found string
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || len(found) > 0 {
			return filepath.SkipDir
		}
		if info.IsDir() && computeDepth(root, p) > maxDepth {
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.EqualFold(info.Name(), filename) {
			found = p
			return filepath.SkipDir
		}
		return nil
	})

	return found, len(found) > 0
}

func computeDepth(root, p string) int {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return 0
	}

	return len(strings.Split(rel, string(filepath.Separator)))
}

func locateGenericTool(target string) (string, bool) {
	if path, err := exec.LookPath(target); err == nil {
		return filepath.Clean(path), true
	}

	return checkStandardToolDirs(target)
}

func checkStandardToolDirs(target string) (string, bool) {
	dirs := resolveStandardToolDirs()
	for _, d := range dirs {
		candidate := filepath.Join(d, target)
		if pathExists(candidate) {
			return candidate, true
		}
	}

	return "", false
}

func resolveStandardToolDirs() []string {
	if runtime.GOOS == "windows" {
		return []string{
			`C:\BuildTools\bin`,
			`C:\Program Files\Git\bin`,
			`C:\Program Files\Git\usr\bin`,
			`C:\Program Files\LLVM\bin`,
		}
	}

	return []string{
		"/usr/local/bin",
		"/usr/bin",
		"/opt/homebrew/bin",
		"/bin",
	}
}

func pathExists(p string) bool {
	if len(p) == 0 {
		return false
	}
	info, err := os.Stat(p)

	return err == nil && !info.IsDir()
}

func dirExists(p string) bool {
	if len(p) == 0 {
		return false
	}
	info, err := os.Stat(p)

	return err == nil && info.IsDir()
}

func renderLocateResult(res LocateResult, isJSON bool) {
	if isJSON {
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
		return
	}

	if res.HasFound {
		fmt.Println(res.Path)
		return
	}

	fmt.Fprintf(os.Stderr, "Tool '%s' not found.\n", res.Target)
}
