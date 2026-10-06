package cmdide

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcursor"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func runIDEAdd(args []string) error {
	opts, remaining := parseIDEOptions(args)
	if len(remaining) == 0 {
		return apperror.NewValidation("cmdide.runIDEAdd", "E1000", "missing repository path to add")
	}
	absPath, err := validateRepoPath(remaining[0])
	if err != nil {
		return err
	}
	res := addRepoToIDEs(absPath, opts)
	printActionResult(res, opts.IsJSON, opts.IsQuiet)
	return nil
}

func validateRepoPath(rawPath string) (string, error) {
	abs, err := filepath.Abs(rawPath)
	if err != nil {
		return "", apperror.WrapSimple(err, "cmdide.validateRepoPath: abs path")
	}
	gitPath := filepath.Join(abs, ".git")
	info, statErr := os.Stat(gitPath)
	if statErr != nil || (!info.IsDir() && !strings.HasSuffix(rawPath, ".git")) {
		return "", apperror.NewValidation("cmdide.validateRepoPath", "E1000", "path is not a git repository: "+rawPath)
	}
	return abs, nil
}

func addRepoToIDEs(absPath string, opts IDEOptions) IDEActionResult {
	res := IDEActionResult{Path: absPath, Action: "add"}
	addVSCodeAndCursor(absPath, opts, &res)
	addAntigravityAndDesktop(absPath, opts, &res)
	return res
}

func performIDEAddAction(isTargeted, isDryRun bool, applyFn func() bool) bool {
	if !isTargeted {
		return false
	}
	if isDryRun {
		return true
	}
	return applyFn()
}

func addVSCodeAndCursor(absPath string, opts IDEOptions, res *IDEActionResult) {
	res.IsVSCodeAffected = performIDEAddAction(opts.IsVSCodeTargeted, opts.IsDryRun, func() bool {
		pair := vscodepm.Pair{RootPath: absPath, Name: filepath.Base(absPath), Tags: []string{"gitmap"}}
		_, err := vscodepm.SyncMode([]vscodepm.Pair{pair}, vscodepm.MergeModeUnion)
		return err == nil
	})
	res.IsCursorAffected = performIDEAddAction(opts.IsCursorTargeted, opts.IsDryRun, func() bool {
		_, _, err := cmdcursor.SyncCursorProjects([]string{absPath})
		return err == nil
	})
}

func addAntigravityAndDesktop(absPath string, opts IDEOptions, res *IDEActionResult) {
	res.IsAntigravityAffected = performIDEAddAction(opts.IsAntigravityTargeted, opts.IsDryRun, func() bool {
		return workspacesync.SyncAntigravity(absPath, filepath.Base(absPath))
	})
	cli := desktop.ResolveCLI()
	res.IsDesktopAffected = performIDEAddAction(opts.IsDesktopTargeted && cli != "", opts.IsDryRun, func() bool {
		return exec.Command(cli, absPath).Run() == nil
	})
}

func titleWord(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func printActionResult(res IDEActionResult, isJSON, isQuiet bool) {
	if isJSON {
		b, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(b))
		return
	}
	if !isQuiet {
		fmt.Printf("%s✔ IDE %s completed for:%s %s\n", constants.ColorGreen, titleWord(res.Action), constants.ColorReset, res.Path)
		fmt.Printf("  • VS Code:     %s\n", formatStatusBool(res.IsVSCodeAffected))
		fmt.Printf("  • Cursor:      %s\n", formatStatusBool(res.IsCursorAffected))
		fmt.Printf("  • Antigravity: %s\n", formatStatusBool(res.IsAntigravityAffected))
		fmt.Printf("  • Desktop:     %s\n", formatStatusBool(res.IsDesktopAffected))
	}
}

func formatStatusBool(val bool) string {
	if val {
		return constants.ColorGreen + "✔ updated" + constants.ColorReset
	}
	return constants.ColorDim + "— skipped" + constants.ColorReset
}
