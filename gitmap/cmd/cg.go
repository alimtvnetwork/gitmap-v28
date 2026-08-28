package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
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
		os.Exit(1)
	}

	opts.Action = argsAfterParse[0]
	if opts.Action == "repo" && len(argsAfterParse) > 1 {
		opts.Action = argsAfterParse[1]
		if len(argsAfterParse) > 2 {
			opts.Repos = argsAfterParse[2:]
		}
	} else if len(argsAfterParse) > 1 {
		if argsAfterParse[1] == "all" {
			opts.All = true
			if len(argsAfterParse) > 2 {
				opts.Repos = argsAfterParse[2:]
			}
		} else {
			opts.Repos = argsAfterParse[1:]
		}
	}

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
	case "install-prompts", "install-prompt", "prompts-install":
		runCGInstallPrompts(repos, opts.DryRun)
	case "update-prompts", "update-prompt", "prompts-update":
		runCGInstallPrompts(repos, opts.DryRun)
	case "prompts-status", "prompts-ls", "prompt-status":
		runPromptStatus(repos)
	case "prompts-version", "prompt-version":
		runPromptVersion(repos)
	case "update":
		// Selective update: only update repos with coding-guidelines in version.json
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
	case "install":
		executeCGWorkers(repos)
	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", opts.Action)
		os.Exit(1)
	}
	return nil
}

func resolveCGRepos(opts cgOptions) []string {
	var targetRepos []string

	if opts.All {
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

		for _, r := range repos {
			targetRepos = append(targetRepos, r.AbsolutePath)
		}
	} else if len(opts.Repos) > 0 {
		targetRepos = ResolveAllCGTargets(opts.Repos)
	} else {
		cwd, _ := os.Getwd()
		gitDir := filepath.Join(cwd, ".git")
		if info, errStat := os.Stat(gitDir); errStat == nil && (info.IsDir() || !info.IsDir()) {
			targetRepos = []string{cwd}
		} else if childRepos, err := fsutil.DiscoverChildGitRepos(cwd); err == nil && len(childRepos) > 0 {
			targetRepos = childRepos
		} else {
			targetRepos = []string{cwd}
		}
	}

	return applyCGExclusions(targetRepos, opts.Exclude)
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
