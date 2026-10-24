package cmdconflict

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

func parseResolveFlags(args []string) (string, string) {
	target := ""
	strategy := ""
	for _, arg := range args {
		target, strategy = processResolveFlagArg(arg, target, strategy)
	}

	return target, strategy
}

func processResolveFlagArg(arg, target, strategy string) (string, string) {
	if act, err := enums.ParseConflictActionType(arg); err == nil {
		return target, act.String()
	}

	if !strings.HasPrefix(arg, "-") && target == "" {
		return arg, strategy
	}

	return target, strategy
}

// ResolveConflictAction performs a targeted resolution action on repoPath.
// Returns true if the action was recognized and executed.
func ResolveConflictAction(repoPath string, action string) (bool, error) {
	act, err := enums.ParseConflictActionType(action)
	if err != nil {
		return false, nil
	}

	return true, dispatchConflictAction(repoPath, act)
}

func dispatchConflictAction(repoPath string, act enums.ConflictActionType) error {
	switch act {
	case enums.ConflictActionTheirs:
		return applyTheirs(repoPath)
	case enums.ConflictActionOurs:
		return applyOurs(repoPath)
	default:
		return dispatchExtendedAction(repoPath, act)
	}
}

func dispatchExtendedAction(repoPath string, act enums.ConflictActionType) error {
	switch act {
	case enums.ConflictActionIDE:
		return applyIDE(repoPath)
	case enums.ConflictActionAbort:
		return applyAbort(repoPath)
	case enums.ConflictActionRerere:
		return applyRerere(repoPath)
	default:
		return nil
	}
}

func applyTheirs(repoPath string) error {
	return applySideResolution(repoPath, constants.GitArgTheirs, constants.MsgAcceptedTheirsFmt)
}

func applyOurs(repoPath string) error {
	return applySideResolution(repoPath, constants.GitArgOurs, constants.MsgAcceptedOursFmt)
}

func applySideResolution(repoPath, sideFlag, msgFmt string) error {
	files, err := GetConflictedFiles(repoPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Printf(constants.MsgNoConflictedFilesDetectedFmt,
			constants.ColorYellow+constants.GlyphWarn+constants.ColorReset, repoPath)

		return nil
	}

	return finishSideResolution(repoPath, sideFlag, msgFmt, files)
}

func finishSideResolution(repoPath, sideFlag, msgFmt string, files []string) error {
	checkoutAndStageFiles(repoPath, sideFlag, files)
	fmt.Printf(msgFmt, constants.ColorGreen+constants.GlyphCheck+constants.ColorReset, len(files), repoPath)
	fmt.Println(constants.MsgNextCommitInstruction)

	return nil
}

func checkoutAndStageFiles(repoPath, sideFlag string, files []string) {
	for _, f := range files {
		_ = exec.Command(constants.GitBin, constants.GitFlagDir, repoPath, constants.GitSubCmdCheckout, sideFlag, constants.GitArgDoubleDash, f).Run()
		_ = exec.Command(constants.GitBin, constants.GitFlagDir, repoPath, constants.GitArgAdd, constants.GitArgDoubleDash, f).Run()
	}
}

func applyIDE(repoPath string) error {
	files, _ := GetConflictedFiles(repoPath)
	ideErr := LaunchIDEForRepo(repoPath, files)
	if ideErr != nil {
		fmt.Printf(constants.MsgIDELaunchFailedFmt, constants.ColorYellow+constants.GlyphWarn+constants.ColorReset, ideErr)
	} else {
		_, name := DetectInstalledIDE()
		fmt.Printf(constants.MsgIDELaunchedFmt, constants.ColorGreen+constants.GlyphCheck+constants.ColorReset, name, repoPath)
	}

	fmt.Print(FormatMergeEditorGuide(repoPath, files))

	return nil
}

func applyAbort(repoPath string) error {
	if hasAbortedSubCmd(repoPath, constants.GitSubCmdMerge, constants.MsgMergeAbortedFmt) {
		return nil
	}

	if hasAbortedSubCmd(repoPath, constants.GitSubCmdRebase, constants.MsgRebaseAbortedFmt) {
		return nil
	}

	if hasAbortedSubCmd(repoPath, constants.GitSubCmdCherryPick, constants.MsgCherryPickAbortedFmt) {
		return nil
	}

	return apperror.NewExecutionError(fmt.Sprintf(constants.ErrNoAbortInProgressFmt, repoPath))
}

func hasAbortedSubCmd(repoPath, subCmd, msgFmt string) bool {
	cmd := exec.Command(constants.GitBin, constants.GitFlagDir, repoPath, subCmd, constants.GitArgAbort)
	if cmd.Run() != nil {
		return false
	}

	fmt.Printf(msgFmt, constants.ColorGreen+constants.GlyphCheck+constants.ColorReset, repoPath)

	return true
}

func applyRerere(repoPath string) error {
	_ = exec.Command(constants.GitBin, constants.GitFlagDir, repoPath, constants.GitSubCmdConfig, constants.GitConfigRerereEnabled, constants.GitTrue).Run()
	_ = exec.Command(constants.GitBin, constants.GitFlagDir, repoPath, constants.GitSubCmdConfig, constants.GitConfigRerereAutoupdate, constants.GitTrue).Run()
	fmt.Printf(constants.MsgRerereEnabledFmt,
		constants.ColorGreen+constants.GlyphCheck+constants.ColorReset, repoPath)

	return nil
}

// RunResolveCmd implements 'gitmap resolve <repo> [--theirs | --ours | --ide | --abort]'.
func RunResolveCmd(args []string) error {
	if hasHelpFlag(args) {
		helpdoc.Print(constants.CmdResolve)

		return nil
	}

	target, strategy := parseResolveFlags(args)

	return executeResolve(target, strategy)
}

func executeResolve(target, strategy string) error {
	repoPath := ResolveRepoPath(target)
	if repoPath == "" {
		return apperror.NewValidationError(constants.ErrResolveUsage)
	}

	if strategy != "" {
		_, err := ResolveConflictAction(repoPath, strategy)

		return err
	}

	return resolveDefaultWorkflow(target, repoPath)
}

func resolveDefaultWorkflow(target, repoPath string) error {
	files, err := GetConflictedFiles(repoPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Printf(constants.MsgNoConflictsDetectedFmt, constants.ColorGreen+constants.GlyphCheck+constants.ColorReset, repoPath)

		return nil
	}

	return presentConflictAndLaunch(target, repoPath, files)
}

func presentConflictAndLaunch(target, repoPath string, files []string) error {
	displayName := target
	if displayName == "" {
		displayName = filepath.Base(getAbsPath(repoPath))
	}

	printConflictTree(displayName, repoPath, files)
	printQuickActions(displayName)
	autoLaunchIDEIfAvailable(repoPath, files)

	return nil
}

func autoLaunchIDEIfAvailable(repoPath string, files []string) {
	_, name := DetectInstalledIDE()
	if name == "" {
		return
	}

	fmt.Printf(constants.MsgAutoLaunchingIDEFmt, name)
	_ = LaunchIDEForRepo(repoPath, files)
	fmt.Print(FormatMergeEditorGuide(repoPath, files))
}
