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
	name := resolveRepoDisplayName(absPath, opts.Name)
	folderURI := formatFolderURI(absPath)
	res.IsVSCodeAffected = performIDEAddAction(opts.IsVSCodeTargeted, opts.IsDryRun, func() bool {
		tags := buildProjectTags(opts.IsPinned)
		pair := vscodepm.Pair{RootPath: absPath, Name: name, Paths: []string{absPath}, Tags: tags}
		_, errSync := vscodepm.SyncMode([]vscodepm.Pair{pair}, vscodepm.MergeModeUnion)
		errDb := UpdateStateVscdbForVSCode(folderURI, name)

		return errSync == nil || errDb == nil
	})
	res.IsCursorAffected = performIDEAddAction(opts.IsCursorTargeted, opts.IsDryRun, func() bool {
		_, _, errSync := cmdcursor.SyncCursorProjects([]string{absPath})
		errDb := UpdateStateVscdbForCursor(folderURI, name)

		return errSync == nil || errDb == nil
	})
}

func buildProjectTags(hasPinned bool) []string {
	if hasPinned {
		return []string{"gitmap", "pinned"}
	}

	return []string{"gitmap"}
}

func resolveRepoDisplayName(absPath, customName string) string {
	if customName != "" {
		return customName
	}

	return filepath.Base(absPath)
}

func formatFolderURI(absPath string) string {
	clean := filepath.ToSlash(absPath)
	if len(clean) > 1 && clean[1] == ':' {
		return "file:///" + string(clean[0]) + "%3A" + clean[2:]
	}
	if strings.HasPrefix(clean, "/") {
		return "file://" + clean
	}

	return "file:///" + clean
}

func addAntigravityAndDesktop(absPath string, opts IDEOptions, res *IDEActionResult) {
	name := resolveRepoDisplayName(absPath, opts.Name)
	res.IsAntigravityAffected = performIDEAddAction(opts.IsAntigravityTargeted, opts.IsDryRun, func() bool {
		isSynced := workspacesync.SyncAntigravity(absPath, name)
		if opts.IsPinned {
			_ = SyncAntigravityPinned(absPath, name)
		}

		return isSynced
	})
	if !opts.IsDesktopTargeted {
		return
	}
	if opts.IsDryRun {
		res.IsDesktopAffected = true

		return
	}
	executeDesktopAdd(absPath, res)
}

func executeDesktopAdd(absPath string, res *IDEActionResult) {
	initErrorDetails(res)
	cli := desktop.ResolveCLI()
	if cli == "" {
		res.ErrorDetails["desktop"] = "GitHub Desktop CLI not found"
		res.HasErrors = true
		return
	}
	output, err := exec.Command(cli, "open", absPath).CombinedOutput()
	if err == nil {
		res.IsDesktopAffected = true
		return
	}
	res.IsDesktopAffected = false
	res.HasErrors = true
	res.ErrorDetails["desktop"] = formatDesktopError(err, output)
}

func initErrorDetails(res *IDEActionResult) {
	if res.ErrorDetails == nil {
		res.ErrorDetails = make(map[string]string)
	}
}

func formatDesktopError(err error, output []byte) string {
	trimmed := strings.TrimSpace(string(output))
	if trimmed != "" {
		return err.Error() + ": " + trimmed
	}
	return err.Error()
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
	if isQuiet {
		return
	}
	printActionSummaryLines(res)
}

func printActionSummaryLines(res IDEActionResult) {
	fmt.Printf("%s✔ IDE %s completed for:%s %s\n", constants.ColorGreen, titleWord(res.Action), constants.ColorReset, res.Path)
	fmt.Printf("  • VS Code:     %s\n", formatIDEStatus(res.IsVSCodeAffected, "vscode", res.ErrorDetails))
	fmt.Printf("  • Cursor:      %s\n", formatIDEStatus(res.IsCursorAffected, "cursor", res.ErrorDetails))
	fmt.Printf("  • Antigravity: %s\n", formatIDEStatus(res.IsAntigravityAffected, "antigravity", res.ErrorDetails))
	fmt.Printf("  • Desktop:     %s\n", formatIDEStatus(res.IsDesktopAffected, "desktop", res.ErrorDetails))
	tokens := formatRegisteredTokens(res)
	if tokens != "" {
		fmt.Printf("\nREGISTERED:%s\n", tokens)
	}
}

func formatRegisteredTokens(res IDEActionResult) string {
	var tokens []string
	if res.IsVSCodeAffected {
		tokens = append(tokens, "vscode")
	}
	if res.IsCursorAffected {
		tokens = append(tokens, "cursor")
	}
	if res.IsAntigravityAffected {
		tokens = append(tokens, "antigravity")
	}
	if res.IsDesktopAffected {
		tokens = append(tokens, "github-desktop")
	}

	return strings.Join(tokens, ",")
}

func formatIDEStatus(isAffected bool, key string, errDetails map[string]string) string {
	if isAffected {
		return constants.ColorGreen + "✔ updated" + constants.ColorReset
	}
	if errMsg, hasErr := errDetails[key]; hasErr && errMsg != "" {
		return constants.ColorRed + "✗ failed (" + errMsg + ")" + constants.ColorReset
	}
	return constants.ColorDim + "— skipped" + constants.ColorReset
}
