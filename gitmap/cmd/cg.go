package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

var cgCommand = "cg"

type cgOptions struct {
	All     bool
	Exclude string
	Repos   []string
	Action  string
}

func parseCGFlags(args []string) cgOptions {
	fs := flag.NewFlagSet(cgCommand, flag.ExitOnError)
	var opts cgOptions
	fs.BoolVar(&opts.All, "all", false, "Run on all workspaces")
	fs.StringVar(&opts.Exclude, "except", "", "Exclude repos (comma separated)")
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude repos (comma separated)")
	fs.Parse(reorderFlagsBeforeArgs(args))

	argsAfterParse := fs.Args()
	if len(argsAfterParse) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap cg install|update [--all] [--except repo1,repo2] [repo1 repo2...]\n")
		os.Exit(1)
	}

	opts.Action = argsAfterParse[0]
	if opts.Action == "all" {
		opts.All = true
		if len(argsAfterParse) > 1 {
			opts.Repos = argsAfterParse[1:]
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

func runCG(args []string) {
	opts := parseCGFlags(args)
	if opts.Action != "install" && opts.Action != "update" {
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", opts.Action)
		os.Exit(1)
	}

	repos := resolveCGRepos(opts)
	if len(repos) == 0 {
		fmt.Println("No repositories to process.")
		return
	}

	executeCGWorkers(repos)
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
	} else {
		targetRepos = opts.Repos
	}

	return applyCGExclusions(targetRepos, opts.Exclude)
}

func applyCGExclusions(repos []string, excludeCSV string) []string {
	if excludeCSV == "" {
		return repos
	}
	
	excludeList := strings.Split(excludeCSV, ",")
	var filtered []string
	for _, r := range repos {
		excluded := false
		for _, ex := range excludeList {
			if strings.TrimSpace(ex) == r {
				excluded = true
				break
			}
		}
		if !excluded {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
