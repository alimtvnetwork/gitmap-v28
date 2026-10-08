package cmdsync

import (
	"flag"
	"fmt"
	"strings"
	"sync"
)

// RunSync is the main entry point for the gitmap sync CLI command.
func RunSync(args []string) error {
	opts, parseErr := parseSyncFlags(args)
	if parseErr != nil {
		return parseErr
	}

	sourceRoot, srcErr := ResolveSourceRoot()
	if srcErr != nil {
		return fmt.Errorf("gitmap sync error: %w", srcErr)
	}

	targets, targetErr := ResolveTargetProjects(opts, sourceRoot)
	if targetErr != nil {
		return fmt.Errorf("failed to resolve target projects: %w", targetErr)
	}

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
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)

	var repoFlag string
	fs.StringVar(&opts.Projects, "projects", "", "Path to projects JSON file or inline JSON string")
	fs.StringVar(&opts.Projects, "p", "", "Alias for --projects")
	fs.StringVar(&repoFlag, "repo", "", "Filter to specific repository name or comma-separated names")
	fs.StringVar(&repoFlag, "r", "", "Alias for --repo")
	fs.IntVar(&opts.Workers, "workers", 8, "Parallel worker goroutine count")
	fs.IntVar(&opts.Workers, "w", 8, "Alias for --workers")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Preview sync without mutating git state or remotes")
	fs.BoolVar(&opts.DryRun, "n", false, "Alias for --dry-run")
	fs.BoolVar(&opts.NoPush, "no-push", false, "Commit locally without pushing to remote")
	fs.BoolVar(&opts.NoRelease, "no-release", false, "Skip version bump release tagging")
	fs.BoolVar(&opts.Verbose, "verbose", false, "Verbose diagnostic output")

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	if repoFlag != "" {
		parts := strings.Split(repoFlag, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				opts.Repos = append(opts.Repos, trimmed)
			}
		}
	}

	if opts.Workers <= 0 {
		opts.Workers = 4
	}
	if opts.Workers > 16 {
		opts.Workers = 16
	}

	return opts, nil
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
	numWorkers := opts.Workers
	if numWorkers > len(targets) {
		numWorkers = len(targets)
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for proj := range jobs {
				res := SyncSingleRepo(sourceRoot, proj, opts)
				statusTag := fmt.Sprintf("[%s]", res.Status)
				fmt.Printf("%-6s %-32s | Pre: %-8s | Post: %-8s | +%d/~%d/-%d | %s\n",
					statusTag, res.Repo, res.PreTag, res.PostTag, res.FilesAdded, res.FilesUpdated, res.FilesRemoved, res.Action)
				results <- res
			}
		}()
	}

	for _, t := range targets {
		jobs <- t
	}
	close(jobs)

	wg.Wait()
	close(results)

	var list []RepoSyncResult
	for r := range results {
		list = append(list, r)
	}

	return list
}

func printSyncSummary(results []RepoSyncResult) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 95))
	fmt.Println("MULTI-REPOSITORY SYNCHRONIZATION SUMMARY")
	fmt.Println(strings.Repeat("=", 95))
	fmt.Printf("%-34s | %-8s | %-10s | %-10s | %-12s | %s\n",
		"Repository", "Status", "Pre-Tag", "Post-Tag", "Files", "Action")
	fmt.Println(strings.Repeat("-", 95))

	for _, r := range results {
		filesStr := fmt.Sprintf("+%d/~%d", r.FilesAdded, r.FilesUpdated)
		fmt.Printf("%-34s | %-8s | %-10s | %-10s | %-12s | %s\n",
			r.Repo, r.Status, r.PreTag, r.PostTag, filesStr, r.Action)
	}
	fmt.Println()
}
