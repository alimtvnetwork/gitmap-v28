package cmdgit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CompareTool holds discovered executable information.
type CompareTool struct {
	Name string
	Path string
}

func findCompareTool() (CompareTool, error) {
	if tool, isFound := findBeyondCompare(); isFound {
		return tool, nil
	}

	if tool, isFound := findVSCode(); isFound {
		return tool, nil
	}

	if tool, isFound := findMeld(); isFound {
		return tool, nil
	}

	return findWinMerge()
}

func findBeyondCompare() (CompareTool, bool) {
	candidates := getBeyondCompareCandidates()
	for _, candidate := range candidates {
		if path, isFound := lookExecutable(candidate); isFound {
			return CompareTool{Name: constants.CompareToolBComp, Path: path}, true
		}
	}

	return CompareTool{}, false
}

func getBeyondCompareCandidates() []string {
	names := []string{"bcomp", "bcompare"}
	if runtime.GOOS == "windows" {
		return append(names, getWindowsBeyondComparePaths()...)
	}

	return append(names, "/usr/bin/bcompare", "/usr/local/bin/bcompare")
}

func getWindowsBeyondComparePaths() []string {
	return []string{
		`C:\Program Files\Beyond Compare 5\BComp.exe`,
		`C:\Program Files\Beyond Compare 5\BCompare.exe`,
		`C:\Program Files\Beyond Compare 4\BComp.exe`,
		`C:\Program Files\Beyond Compare 4\BCompare.exe`,
		`C:\Program Files (x86)\Beyond Compare 4\BComp.exe`,
		`C:\Program Files (x86)\Beyond Compare 4\BCompare.exe`,
	}
}

func findVSCode() (CompareTool, bool) {
	if path, isFound := lookExecutable("code"); isFound {
		return CompareTool{Name: constants.CompareToolCode, Path: path}, true
	}

	if runtime.GOOS == "windows" {
		return findWindowsVSCode()
	}

	return CompareTool{}, false
}

func findWindowsVSCode() (CompareTool, bool) {
	localAppData := os.Getenv("LOCALAPPDATA")
	path := filepath.Join(localAppData, "Programs", "Microsoft VS Code", "bin", "code.cmd")
	if isFileExist(path) {
		return CompareTool{Name: constants.CompareToolCode, Path: path}, true
	}

	return CompareTool{}, false
}

func findMeld() (CompareTool, bool) {
	if path, isFound := lookExecutable("meld"); isFound {
		return CompareTool{Name: constants.CompareToolMeld, Path: path}, true
	}

	return CompareTool{}, false
}

func findWinMerge() (CompareTool, error) {
	candidates := []string{"winmergeu", "winmerge"}
	if runtime.GOOS == "windows" {
		candidates = append(candidates, `C:\Program Files\WinMerge\WinMergeU.exe`, `C:\Program Files (x86)\WinMerge\WinMergeU.exe`)
	}

	for _, candidate := range candidates {
		if path, isFound := lookExecutable(candidate); isFound {
			return CompareTool{Name: constants.CompareToolWinMerge, Path: path}, nil
		}
	}

	return CompareTool{}, apperror.NewNotFound("findCompareTool", "E1004", "no compare tool found (Beyond Compare, VS Code, Meld, WinMerge)")
}

func lookExecutable(nameOrPath string) (string, bool) {
	if isFileExist(nameOrPath) {
		return nameOrPath, true
	}

	path, err := exec.LookPath(nameOrPath)

	return path, err == nil
}

func isFileExist(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)

	return err == nil && !info.IsDir()
}

func launchToolWithPaths(tool CompareTool, dirs []string) error {
	args := buildToolArguments(tool.Name, dirs)
	cmd := exec.Command(tool.Path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Launching %s to compare: %v\n", tool.Name, dirs)
	err := cmd.Run()
	if err != nil {
		return apperror.WrapSimple(err, "launch compare tool: "+tool.Name)
	}

	return nil
}

func buildToolArguments(toolName string, dirs []string) []string {
	if toolName == constants.CompareToolCode && len(dirs) == 2 {
		return []string{"--wait", "--diff", dirs[0], dirs[1]}
	}

	if toolName == constants.CompareToolCode {
		return append([]string{"--wait"}, dirs...)
	}

	return dirs
}
