package cmdide

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcursor"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

func runIDERemove(args []string) error {
	opts, remaining := parseIDEOptions(args)
	if len(remaining) == 0 {
		return apperror.NewValidation("cmdide.runIDERemove", "E1000", "missing repository path to remove")
	}
	absPath, err := filepath.Abs(remaining[0])
	if err != nil {
		return apperror.WrapSimple(err, "cmdide.runIDERemove: abs path")
	}
	res := removeRepoFromIDEs(absPath, opts)
	printActionResult(res, opts.IsJSON, opts.IsQuiet)
	return nil
}

func removeRepoFromIDEs(absPath string, opts IDEOptions) IDEActionResult {
	res := IDEActionResult{Path: absPath, Action: "remove"}
	removeVSCodeAndCursor(absPath, opts, &res)
	removeAntigravityAndDesktop(absPath, opts, &res)
	return res
}

func performIDERemoveAction(isTargeted, isDryRun bool, applyFn func() bool) bool {
	if !isTargeted {
		return false
	}
	if isDryRun {
		return true
	}
	return applyFn()
}

func removeVSCodeAndCursor(absPath string, opts IDEOptions, res *IDEActionResult) {
	res.IsVSCodeAffected = performIDERemoveAction(opts.IsVSCodeTargeted, opts.IsDryRun, func() bool {
		return vscodepm.RemoveEntry(absPath) == nil
	})
	res.IsCursorAffected = performIDERemoveAction(opts.IsCursorTargeted, opts.IsDryRun, func() bool {
		cursorPath, _ := cmdcursor.GetCursorProjectsJSONPath()
		return vscodepm.RemoveEntryAt(cursorPath, absPath) == nil
	})
}

func removeAntigravityAndDesktop(absPath string, opts IDEOptions, res *IDEActionResult) {
	res.IsAntigravityAffected = performIDERemoveAction(opts.IsAntigravityTargeted, opts.IsDryRun, func() bool {
		return removeAgyProjectByPath(absPath)
	})
	res.IsDesktopAffected = performIDERemoveAction(opts.IsDesktopTargeted, opts.IsDryRun, func() bool {
		return desktop.RemoveRepo(absPath) == nil
	})
}

func removeAgyProjectByPath(absPath string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	configDir := filepath.Join(home, ".gemini", "config", "projects")
	entries, rErr := os.ReadDir(configDir)
	if rErr != nil {
		return false
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" && matchAndRemoveAgyFile(configDir, e.Name(), absPath) {
			return true
		}
	}
	return false
}

func matchAndRemoveAgyFile(configDir, fileName, absPath string) bool {
	fullPath := filepath.Join(configDir, fileName)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return false
	}
	cleanRepo := filepath.ToSlash(filepath.Clean(absPath))
	if strings.Contains(string(data), cleanRepo) {
		return os.Remove(fullPath) == nil
	}
	return false
}
