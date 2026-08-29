package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

var cgCommand = "cg"

type cgOptions struct {
	All          bool
	DryRun       bool
	Exclude      string
	VersionValue string
	Repos        []string
	Action       string
}

func parseCGFlags(args []string) cgOptions {
	fs := flag.NewFlagSet(cgCommand, flag.ExitOnError)
	var opts cgOptions
	fs.BoolVar(&opts.All, "all", false, "Run on all workspaces")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Simulate execution without modifying files")
	fs.StringVar(&opts.Exclude, "except", "", "Exclude repos (comma separated)")
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude repos (comma separated)")
	fs.StringVar(&opts.VersionValue, "version", "", "Starting version value for version.json (default: 1.0.0)")
	fs.StringVar(&opts.VersionValue, "version-value", "", "Starting version value for version.json (default: 1.0.0)")
	fs.Parse(reorderFlagsBeforeArgs(args))

	argsAfterParse := fs.Args()
	if len(argsAfterParse) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap cg [version|status|install|update|install-version-json|install-prompts|update-prompts|prompts-status|prompts-version] [repo1 repo2...] [--all]\n")
		cliexit.HandleError(nil, 1)
	}

	opts.Action = argsAfterParse[0]
	if opts.Action == "repo" && len(argsAfterParse) > 1 {
		return parseCGRepoAction(opts, argsAfterParse)
	}
	return parseCGGeneralAction(opts, argsAfterParse)
}

func parseCGRepoAction(opts cgOptions, args []string) cgOptions {
	opts.Action = args[1]
	if len(args) > 2 {
		opts.Repos = args[2:]
	}
	return opts
}

func parseCGGeneralAction(opts cgOptions, args []string) cgOptions {
	if len(args) <= 1 {
		return opts
	}
	if args[1] == "all" && len(args) > 2 {
		opts.All = true
		opts.Repos = args[2:]
		return opts
	}
	if args[1] == "all" {
		opts.All = true
		return opts
	}
	opts.Repos = args[1:]
	return opts
}

func runCG(args []string) error {
	opts := parseCGFlags(args)

	repos := resolveCGRepos(opts)
	if len(repos) == 0 {
		fmt.Println("No repositories to process.")
		return nil
	}

	switch opts.Action {
	case "version", "ver", "-v":
		PrintCGVersion(repos)
	case "status", "stat", "ls":
		PrintCGStatus(repos)
	case "install-version-json", "install-version", "init-version", "version-init":
		runCGInstallVersionJSON(repos, opts.VersionValue, opts.DryRun)
	case "install-prompts", "install-prompt", "prompts-install", "update-prompts", "update-prompt", "prompts-update":
		runCGInstallPrompts(repos, opts.DryRun)
	case "prompts-status", "prompts-ls", "prompt-status":
		runPromptStatus(repos)
	case "prompts-version", "prompt-version":
		runPromptVersion(repos)
	case "update":
		runCGUpdateAction(repos)
	case "install":
		executeCGWorkers(repos)
	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", opts.Action)
		cliexit.HandleError(nil, 1)
	}
	return nil
}

func runCGUpdateAction(repos []string) {
	var toUpdate []string
	for _, r := range repos {
		if _, err := ReadCGMetadata(r); err == nil {
			toUpdate = append(toUpdate, r)
		} else {
			fmt.Printf("Skipping %s (coding guidelines not installed, run `gitmap cg install` first)\n", r)
		}
	}
	if len(toUpdate) > 0 {
		executeCGWorkers(toUpdate)
	}
}

func resolveCGRepos(opts cgOptions) []string {
	var targetRepos []string

	switch {
	case opts.All:
		targetRepos = resolveAllCGRepos()
	case len(opts.Repos) > 0:
		targetRepos = ResolveAllCGTargets(opts.Repos)
	default:
		targetRepos = resolveCurrentDirCGRepos()
	}

	return applyCGExclusions(targetRepos, opts.Exclude)
}

func resolveAllCGRepos() []string {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open DB: %v\n", err)
		return []string{}
	}
	defer db.Close()

	repos, err := db.ListRepos()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list repos: %v\n", err)
		return []string{}
	}

	var targetRepos []string
	for _, r := range repos {
		targetRepos = append(targetRepos, r.AbsolutePath)
	}
	return targetRepos
}

func resolveCurrentDirCGRepos() []string {
	cwd, _ := os.Getwd()
	gitDir := filepath.Join(cwd, ".git")
	if info, errStat := os.Stat(gitDir); errStat == nil && (info.IsDir() || !info.IsDir()) {
		return []string{cwd}
	}
	if childRepos, err := fsutil.DiscoverChildGitRepos(cwd); err == nil && len(childRepos) > 0 {
		return childRepos
	}
	return []string{cwd}
}

func applyCGExclusions(repos []string, excludeCSV string) []string {
	if excludeCSV == "" {
		return repos
	}

	excludes := strings.Split(excludeCSV, ",")
	excludeMap := make(map[string]bool)
	for _, e := range excludes {
		excludeMap[strings.TrimSpace(e)] = true
	}

	var filtered []string
	for _, r := range repos {
		base := strings.TrimSpace(r)
		name := base[strings.LastIndex(base, "/")+1:]
		if name == "" {
			name = base[strings.LastIndex(base, "\\")+1:]
		}
		if !excludeMap[name] && !excludeMap[base] {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
