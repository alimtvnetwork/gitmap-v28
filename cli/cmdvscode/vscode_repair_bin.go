package cmdvscode

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type binaryRepairResult struct {
	ActiveInstall string
	VersionInfo   string
	IsHealthy     bool
	SyncedFolders int
}

var commitDirRegex = regexp.MustCompile(`^[0-9a-f]{10}$`)

func discoverInstallationCandidates() []string {
	var candidates []string
	addIfValid := func(dir string) {
		if dir != "" && isDir(dir) && !contains(candidates, dir) {
			candidates = append(candidates, dir)
		}
	}

	addIfValid(findCandidateFromPath())
	if runtime.GOOS == "windows" {
		addIfValid(filepath.Join(os.Getenv("ProgramFiles"), "Microsoft VS Code"))
		addIfValid(filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft VS Code"))
		addIfValid(filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Microsoft VS Code"))
	}

	return candidates
}

func findCandidateFromPath() string {
	codePath, err := exec.LookPath("code.cmd")
	if err != nil {
		codePath, err = exec.LookPath("code")
	}
	if err != nil || codePath == "" {
		return ""
	}

	binDir := filepath.Dir(codePath)

	return filepath.Dir(binDir)
}

func testCodeExecutable(installDir string) (string, bool) {
	cmdName := resolveCodeLauncher(installDir)
	if cmdName == "" {
		return "", false
	}

	cmd := exec.Command(cmdName, "--version")
	out, err := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(out))
	if err == nil && !strings.Contains(outStr, "icu_util") {
		return outStr, true
	}

	return outStr, false
}

func resolveCodeLauncher(installDir string) string {
	if runtime.GOOS == "windows" {
		binCmd := filepath.Join(installDir, "bin", "code.cmd")
		if isFile(binCmd) {
			return binCmd
		}
		exe := filepath.Join(installDir, "Code.exe")
		if isFile(exe) {
			return exe
		}
	}

	stdBin := filepath.Join(installDir, "bin", "code")
	if isFile(stdBin) {
		return stdBin
	}

	return ""
}

func syncCommitFoldersBetween(candidates []string) int {
	totalSynced := 0
	for _, src := range candidates {
		commitDirs := findCommitDirs(src)
		for _, dst := range candidates {
			if src == dst {
				continue
			}
			totalSynced += copyMissingCommitDirs(commitDirs, src, dst)
		}
	}

	return totalSynced
}

func findCommitDirs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var commitDirs []string
	for _, e := range entries {
		if e.IsDir() && commitDirRegex.MatchString(e.Name()) {
			commitDirs = append(commitDirs, e.Name())
		}
	}

	return commitDirs
}

func copyMissingCommitDirs(dirs []string, srcParent, dstParent string) int {
	synced := 0
	for _, d := range dirs {
		dstPath := filepath.Join(dstParent, d)
		if !isDir(dstPath) {
			srcPath := filepath.Join(srcParent, d)
			if copyDir(srcPath, dstPath) == nil {
				synced++
			}
		}
	}

	return synced
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}

func runWingetRepair() (string, error) {
	wingetPath, err := exec.LookPath("winget")
	if err != nil {
		return "", fmt.Errorf("winget is not available on this machine")
	}

	args := []string{"install", "--id", "Microsoft.VisualStudioCode", "--force", "--accept-source-agreements", "--accept-package-agreements"}
	cmd := exec.Command(wingetPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func isDir(p string) bool {
	fi, err := os.Stat(p)

	return err == nil && fi.IsDir()
}

func isFile(p string) bool {
	fi, err := os.Stat(p)

	return err == nil && !fi.IsDir()
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return true
		}
	}

	return false
}
