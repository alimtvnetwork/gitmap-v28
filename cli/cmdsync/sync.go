package cmdsync

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/enums"
)

// RunSync is the main entry point for the gitmap sync CLI command.
func RunSync(args []string) error {
	return RunSyncWithScope(enums.SyncScopeAll, args)
}

// RunSyncWithScope executes synchronization targeting specific asset domains.
func RunSyncWithScope(scope enums.SyncScopeType, args []string) error {
	opts, parseErr := parseSyncFlags(args)
	if parseErr != nil {
		return parseErr
	}
	opts.Scope = scope

	sourceRoot, srcErr := ResolveSourceRoot()
	if srcErr != nil {
		return apperror.WrapSimple(srcErr, "gitmap sync error")
	}

	targets, targetErr := ResolveTargetProjects(opts, sourceRoot)
	if targetErr != nil {
		return apperror.WrapSimple(targetErr, "failed to resolve target projects")
	}

	return dispatchSyncFleet(sourceRoot, targets, opts)
}

func dispatchSyncFleet(sourceRoot string, targets []ProjectConfig, opts SyncOptions) error {
	if len(targets) == 0 {
		fmt.Println("No target repositories found matching criteria.")

		return nil
	}

	printSyncBanner(sourceRoot, len(targets), opts)
	results := executeSyncFleet(sourceRoot, targets, opts)
	printSyncSummary(results)

	return nil
}

func parseSyncFlags(args []string) (SyncOptions, error) {
	var opts SyncOptions
	var repoFlag string
	fs := flag.NewFlagSet(constants.CmdSync, flag.ContinueOnError)
	bindSyncFlags(fs, &opts, &repoFlag)

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	populateRepoFilter(&opts, repoFlag)
	clampWorkerBounds(&opts)

	return opts, nil
}

func bindSyncFlags(fs *flag.FlagSet, opts *SyncOptions, repoFlag *string) {
	bindPathAndFilterFlags(fs, opts, repoFlag)
	bindConcurrencyAndExecutionFlags(fs, opts)
}

func bindPathAndFilterFlags(fs *flag.FlagSet, opts *SyncOptions, repoFlag *string) {
	fs.StringVar(&opts.Projects, constants.FlagProjects, "", "Path to projects JSON file or inline JSON string")
	fs.StringVar(&opts.Projects, constants.FlagProjectsShort, "", "Alias for --projects")
	fs.StringVar(repoFlag, constants.FlagRepo, "", "Filter to specific repository name or comma-separated names")
	fs.StringVar(repoFlag, constants.FlagRepoShort, "", "Alias for --repo")
	fs.StringVar(repoFlag, "repos", "", "Alias for --repo")
}

func bindConcurrencyAndExecutionFlags(fs *flag.FlagSet, opts *SyncOptions) {
	fs.IntVar(&opts.Workers, constants.FlagWorkers, constants.DefaultWorkers, "Parallel worker goroutine count")
	fs.IntVar(&opts.Workers, constants.FlagWorkersShort, constants.DefaultWorkers, "Alias for --workers")
	fs.BoolVar(&opts.DryRun, constants.FlagDryRun, false, "Preview sync without mutating git state or remotes")
	fs.BoolVar(&opts.DryRun, constants.FlagDryRunShort, false, "Alias for --dry-run")
	fs.BoolVar(&opts.NoPush, constants.FlagNoPush, false, "Commit locally without pushing to remote")
	fs.BoolVar(&opts.NoRelease, constants.FlagNoRelease, false, "Skip version bump release tagging")
	fs.BoolVar(&opts.Verbose, constants.FlagVerbose, false, "Verbose diagnostic output")
}

func populateRepoFilter(opts *SyncOptions, repoFlag string) {
	for _, p := range strings.Split(repoFlag, ",") {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			opts.Repos = append(opts.Repos, trimmed)
		}
	}
}

func clampWorkerBounds(opts *SyncOptions) {
	if opts.Workers <= 0 {
		opts.Workers = constants.MinWorkers
	}
	if opts.Workers > constants.MaxWorkers {
		opts.Workers = constants.MaxWorkers
	}
}

func printSyncBanner(sourceRoot string, count int, opts SyncOptions) {
	fmt.Printf("Source repository : %s\n", sourceRoot)
	fmt.Printf("Target count      : %d\n", count)
	fmt.Printf("Workers           : %d\n", opts.Workers)
	fmt.Printf("Dry run           : %v\n", opts.DryRun)
	fmt.Printf("No push           : %v\n", opts.NoPush)
	fmt.Println(strings.Repeat("=", 95))
}

func executeSyncFleet(sourceRoot string, targets []ProjectConfig, opts SyncOptions) []RepoSyncResult {
	jobs := make(chan ProjectConfig, len(targets))
	results := make(chan RepoSyncResult, len(targets))

	var wg sync.WaitGroup
	numWorkers := calculateWorkerCount(opts.Workers, len(targets))

	spawnSyncWorkers(&wg, numWorkers, sourceRoot, opts, jobs, results)
	enqueueSyncJobs(jobs, targets)

	wg.Wait()
	close(results)

	return collectSyncResults(results)
}

func calculateWorkerCount(configuredWorkers, targetCount int) int {
	if configuredWorkers > targetCount {
		return targetCount
	}

	return configuredWorkers
}

func spawnSyncWorkers(wg *sync.WaitGroup, count int, sourceRoot string, opts SyncOptions, jobs <-chan ProjectConfig, results chan<- RepoSyncResult) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runSingleWorkerLoop(sourceRoot, opts, jobs, results)
		}()
	}
}

func runSingleWorkerLoop(sourceRoot string, opts SyncOptions, jobs <-chan ProjectConfig, results chan<- RepoSyncResult) {
	for proj := range jobs {
		res := SyncSingleRepo(sourceRoot, proj, opts)
		printWorkerResult(res)
		results <- res
	}
}

func printWorkerResult(res RepoSyncResult) {
	statusTag := fmt.Sprintf("[%s]", res.Status)
	preservedInfo := ""
	if res.FilesPreserved > 0 || res.FilesConflicts > 0 {
		preservedInfo = fmt.Sprintf(" [preserved %d human edits, %d sidecars]", res.FilesPreserved, res.FilesConflicts)
	}
	fmt.Printf("%-6s %-32s | Pre: %-8s | Post: %-8s | +%d/~%d/-%d | %s%s\n",
		statusTag, res.Repo, res.PreTag, res.PostTag, res.FilesAdded, res.FilesUpdated, res.FilesRemoved, res.Action, preservedInfo)
}

func enqueueSyncJobs(jobs chan<- ProjectConfig, targets []ProjectConfig) {
	for _, t := range targets {
		jobs <- t
	}
	close(jobs)
}

func collectSyncResults(results <-chan RepoSyncResult) []RepoSyncResult {
	var list []RepoSyncResult
	for r := range results {
		list = append(list, r)
	}

	return list
}

func printSyncSummary(results []RepoSyncResult) {
	printSummaryHeader()
	for _, r := range results {
		printSummaryRow(r)
	}
	printHumanPreservedSummary(results)
	fmt.Println()
}

func printHumanPreservedSummary(results []RepoSyncResult) {
	var countPreserved, countConflicts int
	for _, r := range results {
		countPreserved += r.FilesPreserved
		countConflicts += r.FilesConflicts
	}
	if countPreserved == 0 && countConflicts == 0 {
		return
	}

	fmt.Printf("\n● Human Edits Detected & Preserved (Zero Data Loss Invariant):\n")
	for _, r := range results {
		if r.FilesPreserved > 0 || r.FilesConflicts > 0 {
			fmt.Printf("  • %-30s : %d local edit(s) preserved untouched, %d .upstream sidecar(s) generated\n", r.Repo, r.FilesPreserved, r.FilesConflicts)
		}
	}
	fmt.Printf("Notice: Local changes were left intact. Upstream differences are safely isolated in '<file>.upstream' sidecars.\n")
}

func printSummaryHeader() {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 95))
	fmt.Println(constants.SummaryHeader)
	fmt.Println(strings.Repeat("=", 95))
	fmt.Printf("%-34s | %-8s | %-10s | %-10s | %-12s | %s\n",
		constants.SummaryColRepo, constants.SummaryColStatus, constants.SummaryColPreTag,
		constants.SummaryColPostTag, constants.SummaryColFiles, constants.SummaryColAction)
	fmt.Println(strings.Repeat("-", 95))
}

func printSummaryRow(r RepoSyncResult) {
	filesStr := fmt.Sprintf("+%d/~%d", r.FilesAdded, r.FilesUpdated)
	fmt.Printf("%-34s | %-8s | %-10s | %-10s | %-12s | %s\n",
		r.Repo, r.Status, r.PreTag, r.PostTag, filesStr, r.Action)
}
