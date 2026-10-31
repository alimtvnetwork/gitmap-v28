package cmdcg

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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

func isCGGeneralInstall(args []string) bool {
	return len(args) > 1 && (args[1] == "general" || args[1] == "general-install")
}

func parseCGInstallGeneral(opts cgOptions, args []string) cgOptions {
	opts.Action = "general"
	if len(args) > 2 {
		opts.Repos = args[2:]
	}

	return opts
}

// reposFromArgsAfter returns args[start:] when present, else the
// caller's existing repos. WHY: it keeps single-action parsing flat,
// so the nested-if linter never sees a depth-2 repos assignment.
func reposFromArgsAfter(existing []string, args []string, start int) []string {
	if len(args) > start {
		return args[start:]
	}

	return existing
}

func resolveCGSingleAction(opts cgOptions, args []string) cgOptions {
	if (opts.Action == "prompt" || opts.Action == "prompts") && len(args) > 1 && args[1] == "sync" {
		opts.Action = "prompt-sync"
		opts.Repos = reposFromArgsAfter(opts.Repos, args, 2)
		return opts
	}

	if (opts.Action == "skill" || opts.Action == "skills") && len(args) > 1 && args[1] == "sync" {
		opts.Action = "skill-sync"
		opts.Repos = reposFromArgsAfter(opts.Repos, args, 2)
		return opts
	}

	if opts.Action == "execute" && len(args) > 1 {
		opts.Action = "execute"
		opts.VersionValue = args[1]
		opts.Repos = reposFromArgsAfter(opts.Repos, args, 2)
		return opts
	}

	if opts.Action == "repo" && len(args) > 1 {
		return parseCGRepoAction(opts, args)
	}

	if opts.Action == "install" && isCGGeneralInstall(args) {
		return parseCGInstallGeneral(opts, args)
	}

	return parseCGGeneralAction(opts, args)
}

func resolveCGAction(opts cgOptions, args []string) cgOptions {
	if len(args) == 0 || isCGHelpArg(args[0]) {
		opts.Action = "help"

		return opts
	}

	opts.Action = args[0]

	return resolveCGSingleAction(opts, args)
}

func parseCGFlags(args []string) cgOptions {
	if len(args) == 0 || hasCGHelpArg(args) {
		return cgOptions{Action: "help"}
	}

	fs := flag.NewFlagSet(cgCommand, flag.ContinueOnError)
	var opts cgOptions
	initCGFlagSet(fs, &opts)
	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	return resolveCGAction(opts, fs.Args())
}

func parseCGRepoAction(opts cgOptions, args []string) cgOptions {
	opts.Action = args[1]
	if len(args) > 2 {
		opts.Repos = args[2:]
	}

	return opts
}

func parseCGAllAction(opts cgOptions, args []string) cgOptions {
	opts.All = true
	if len(args) > 2 {
		opts.Repos = args[2:]
	}

	return opts
}

func parseCGGeneralAction(opts cgOptions, args []string) cgOptions {
	if len(args) <= 1 {
		return opts
	}

	if args[1] == "all" {
		return parseCGAllAction(opts, args)
	}

	opts.Repos = args[1:]

	return opts
}

func executeCGAction(opts cgOptions, args []string) error {
	if opts.Action == "execute" {
		return RunCGExecutePrompt(opts.VersionValue, args)
	}

	repos := resolveCGRepos(opts)
	if len(repos) == 0 {
		fmt.Println("No repositories to process.")

		return nil
	}

	return dispatchCGAction(opts, repos, args)
}

func runCG(args []string) error {
	opts := parseCGFlags(args)
	if opts.Action == "help" {
		RenderCGHelp()

		return nil
	}

	return executeCGAction(opts, args)
}

func runCGUpdateAction(repos []string) {
	var toUpdate []string
	for _, r := range repos {
		if _, err := ReadCGMetadata(r); err == nil {
			toUpdate = append(toUpdate, r)
			continue
		}

		fmt.Printf("Skipping %s (coding guidelines not installed, run `gitmap cg install` first)\n", r)
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

func extractCGRepoPaths(repos []model.ScanRecord) []string {
	var targetRepos []string
	for _, r := range repos {
		targetRepos = append(targetRepos, r.AbsolutePath)
	}

	return targetRepos
}

func queryAllDBRepos(db *store.DB) []model.ScanRecord {
	repos, errList := db.ListRepos()
	if errList != nil {
		fmt.Fprintf(os.Stderr, "Failed to list repos: %v\n", errList)

		return []model.ScanRecord{}
	}

	return repos
}

func resolveAllCGRepos() []string {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open DB: %v\n", err)

		return []string{}
	}

	defer db.Close()

	return extractCGRepoPaths(queryAllDBRepos(db))
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

func buildCGExcludeMap(excludeCSV string) map[string]bool {
	excludes := strings.Split(excludeCSV, ",")
	excludeMap := make(map[string]bool)
	for _, e := range excludes {
		excludeMap[strings.TrimSpace(e)] = true
	}

	return excludeMap
}

func getCGRepoBaseName(path string) string {
	base := strings.TrimSpace(path)
	name := base[strings.LastIndex(base, "/")+1:]
	if name == "" {
		name = base[strings.LastIndex(base, "\\")+1:]
	}

	return name
}

func isCGRepoExcluded(r string, excludeMap map[string]bool) bool {
	base := strings.TrimSpace(r)
	name := getCGRepoBaseName(base)

	return excludeMap[name] || excludeMap[base]
}

func applyCGExclusions(repos []string, excludeCSV string) []string {
	if excludeCSV == "" {
		return repos
	}

	excludeMap := buildCGExcludeMap(excludeCSV)
	var filtered []string
	for _, r := range repos {
		if !isCGRepoExcluded(r, excludeMap) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}
