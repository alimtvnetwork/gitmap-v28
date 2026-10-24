package cmdhooks

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	apperror "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks/enums"
)

// hookPlan is the per-repo install decision.
type hookPlan struct {
	repo   string
	action enums.HookPlanActionType
	detail string
}

type installOptions struct {
	isAll    bool
	workdir  string
	isDryRun bool
	hasYes   bool
}

type flagPointers struct {
	all     *bool
	workdir *string
	dryRun  *bool
	yes     *bool
}

type planStats struct {
	installed int
	skipped   int
	already   int
}

func hasHelpArg(args []string) bool {
	for _, a := range args {
		if a == constants.FlagArgHelp || a == constants.FlagArgHelpShort {
			return true
		}
	}

	return false
}

func setupInstallFlagSet() (*flag.FlagSet, flagPointers) {
	fs := flag.NewFlagSet(constants.FlagSetHooksInstall, flag.ContinueOnError)
	fps := flagPointers{
		all:     fs.Bool(constants.FlagAll, false, constants.FlagDescAll),
		workdir: fs.String(constants.FlagWorkdir, "", constants.FlagDescWorkdir),
		dryRun:  fs.Bool(constants.FlagDryRun, true, constants.FlagDescDryRun),
		yes:     fs.Bool(constants.FlagYes, false, constants.FlagDescYes),
	}

	return fs, fps
}

func parseInstallFlags(args []string) (installOptions, *appfault.AppError) {
	fs, fps := setupInstallFlagSet()
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(constants.HelpTextHooksInstall)

		return installOptions{}, apperror.WrapSimple(parseErr, constants.ErrWrapFlags)
	}
	if !*fps.all {
		fmt.Print(constants.HelpTextHooksInstall)

		return installOptions{}, appfault.NewValidationError(constants.ErrMissingAllFlag)
	}

	return installOptions{
		isAll:    *fps.all,
		workdir:  *fps.workdir,
		isDryRun: *fps.dryRun && !*fps.yes,
		hasYes:   *fps.yes,
	}, nil
}

func countPlanStats(plans []hookPlan) planStats {
	var stats planStats
	for _, p := range plans {
		switch p.action {
		case enums.HookPlanActionInstall:
			stats.installed++
		case enums.HookPlanActionAlready:
			stats.already++
		default:
			stats.skipped++
		}
	}

	return stats
}

func printInstallPlan(plans []hookPlan, stats planStats) {
	fmt.Printf(constants.MsgPlanHeaderFmt, len(plans), stats.installed, stats.already, stats.skipped)
	for _, p := range plans {
		fmt.Printf(constants.MsgPlanRowFmt, p.action, p.repo, p.detail)
	}
}

func buildHookPlans(workdir string, repos []string) []hookPlan {
	plans := make([]hookPlan, 0, len(repos))
	for _, repo := range repos {
		plans = append(plans, planHookInstall(workdir, repo))
	}

	return plans
}

// RunHooksInstall implements `gitmap hooks install --all` (alias: hkin).
func RunHooksInstall(args []string) *appfault.AppError {
	if hasHelpArg(args) {
		fmt.Print(constants.HelpTextHooksInstall)

		return nil
	}

	opts, err := parseInstallFlags(args)
	if err != nil {
		return err
	}

	return executeHookInstallFlow(opts)
}

func executeHookInstallFlow(opts installOptions) *appfault.AppError {
	workdir := resolveHooksWorkdir(opts.workdir)
	repos := listHookRepos(workdir)
	if len(repos) == 0 {
		fmt.Printf(constants.MsgNoReposFoundFmt, workdir)

		return nil
	}

	return runPlanExecution(workdir, repos, opts.isDryRun)
}

func runPlanExecution(workdir string, repos []string, isDryRun bool) *appfault.AppError {
	plans := buildHookPlans(workdir, repos)
	stats := countPlanStats(plans)
	printInstallPlan(plans, stats)

	if isDryRun {
		fmt.Println(constants.MsgDryRunNotice)

		return nil
	}

	return finalizeHookInstalls(plans, workdir, stats)
}

func finalizeHookInstalls(plans []hookPlan, workdir string, stats planStats) *appfault.AppError {
	written, writeErr := applyHookInstalls(plans, workdir)
	if writeErr != nil {
		return writeErr
	}

	fmt.Printf(constants.MsgInstallSummaryFmt, written, len(plans), stats.already, stats.skipped)

	return nil
}

// resolveHooksWorkdir mirrors the fleet workdir resolution.
func resolveHooksWorkdir(flagVal string) string {
	if strings.TrimSpace(flagVal) != "" {
		return flagVal
	}
	if envVal := strings.TrimSpace(os.Getenv(constants.EnvGitmapWorkdir)); envVal != "" {
		return envVal
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return constants.SentinelCurrentDir
	}

	return filepath.Join(home, constants.DefaultWorkspaceDir, constants.DefaultReposDir)
}

// listHookRepos enumerates workdir directories containing a .git dir.
func listHookRepos(workdir string) []string {
	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		return nil
	}

	return filterGitRepos(workdir, entries)
}

func filterGitRepos(workdir string, entries []os.DirEntry) []string {
	var repos []string
	for _, e := range entries {
		if hasRepoGitDir(workdir, e) {
			repos = append(repos, e.Name())
		}
	}

	sort.Strings(repos)

	return repos
}

func hasRepoGitDir(workdir string, e os.DirEntry) bool {
	if !e.IsDir() {
		return false
	}

	full := filepath.Join(workdir, e.Name(), constants.DirGit)
	_, statErr := os.Stat(full)

	return statErr == nil
}

// planHookInstall decides the action for one repo.
func planHookInstall(workdir, repo string) hookPlan {
	hookPath := filepath.Join(workdir, repo, constants.DirGit, constants.DirHooks, constants.FileCommitMsg)
	if isGuardInstalled(hookPath) {
		return hookPlan{
			repo:   repo,
			action: enums.HookPlanActionAlready,
			detail: constants.DetailAlreadyGuarded,
		}
	}

	return hookPlan{
		repo:   repo,
		action: enums.HookPlanActionInstall,
		detail: constants.PathHookCommitMsg,
	}
}

func installSingleHook(workdir, repo string) (bool, *appfault.AppError) {
	hookPath := filepath.Join(workdir, repo, constants.DirGit, constants.DirHooks, constants.FileCommitMsg)
	if isGuardInstalled(hookPath) {
		return false, nil
	}

	if writeErr := writeGuardHook(hookPath, repo); writeErr != nil {
		return false, writeErr
	}

	return true, nil
}

// applyHookInstalls writes the guard hook for planned repos.
func applyHookInstalls(plans []hookPlan, workdir string) (int, *appfault.AppError) {
	written := 0
	for _, p := range plans {
		if !p.action.IsInstall() {
			continue
		}

		hasWritten, err := installSingleHook(workdir, p.repo)
		if err != nil {
			return written, err
		}
		if hasWritten {
			written++
		}
	}

	return written, nil
}

func writeGuardHook(hookPath, repo string) *appfault.AppError {
	if writeErr := os.WriteFile(hookPath, []byte(constants.CommitMsgGuardTemplate), constants.HookFilePerm); writeErr != nil {
		return apperror.WrapSimple(writeErr, constants.ErrWrapWritePrefix+repo)
	}

	return nil
}

// isGuardInstalled reports whether our guard already occupies hookPath.
func isGuardInstalled(hookPath string) bool {
	data, readErr := os.ReadFile(hookPath)
	if readErr != nil {
		return false
	}

	return strings.Contains(string(data), constants.GuardMarker)
}
