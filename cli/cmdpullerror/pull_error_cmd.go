// Package cmdpullerror provides the CLI command implementation for inspecting repository pull errors.
package cmdpullerror

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunPullErrorCLI is the main entry point for 'gitmap pull-error' and its aliases.
func RunPullErrorCLI(args []string) error {
	if len(args) > 0 && isPullErrorAlias(args[0]) {
		args = args[1:]
	}

	opts := parsePullErrorOptions(args)
	if opts.IsHelp {
		renderPullErrorHelp()
		return nil
	}

	if opts.IsSSH {
		return runPullErrorSSH(opts)
	}

	return executeLocalPullErrors(opts)
}

func isPullErrorAlias(val string) bool {
	switch strings.ToLower(val) {
	case "pull-error", "pull-errors", "pulle", "pull-e":
		return true
	default:
		return false
	}
}

func parsePullErrorOptions(args []string) PullErrorOptions {
	opts := PullErrorOptions{Limit: 50}
	var positional []string

	for i := 0; i < len(args); i++ {
		consumed := parseOptionFlag(args, i, &opts)
		if consumed > 0 {
			i += consumed - 1
			continue
		}
		if !strings.HasPrefix(args[i], "-") {
			positional = append(positional, args[i])
		}
	}

	opts.RepoSlug = resolveTargetRepo(positional)
	return opts
}

func parseOptionFlag(args []string, idx int, opts *PullErrorOptions) int {
	arg := strings.ToLower(args[idx])
	switch arg {
	case "--json", "-j":
		opts.IsJSON = true
		return 1
	case "--ssh":
		opts.IsSSH = true
		return 1
	case "--help", "-h", "help":
		opts.IsHelp = true
		return 1
	case "--limit", "-l":
		return parseLimitOption(args, idx, opts)
	default:
		return 0
	}
}

func parseLimitOption(args []string, idx int, opts *PullErrorOptions) int {
	if idx+1 >= len(args) {
		return 1
	}
	if parsed, err := strconv.Atoi(args[idx+1]); err == nil && parsed > 0 {
		opts.Limit = parsed
		return 2
	}
	return 1
}

func resolveTargetRepo(positional []string) string {
	if len(positional) > 0 {
		return positional[0]
	}
	return resolveDefaultTarget()
}

func resolveDefaultTarget() string {
	if isGitRepoCWD() {
		topLevel := gitRepoTopLevel()
		if topLevel != "" {
			return filepath.Base(topLevel)
		}
		if cwd, err := os.Getwd(); err == nil && cwd != "" {
			return filepath.Base(cwd)
		}
	}
	return "all"
}

func isGitRepoCWD() bool {
	if hasLocalDotGit() {
		return true
	}
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func hasLocalDotGit() bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}
	info, statErr := os.Stat(filepath.Join(cwd, ".git"))
	return statErr == nil && info != nil
}

func gitRepoTopLevel() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func runPullErrorSSH(opts PullErrorOptions) error {
	remoteCmd := buildRemotePullErrorCmd(opts)
	return cmdssh.RunFleetPASCommand("pull-errors", remoteCmd, func() error {
		return executeLocalPullErrors(opts)
	})
}

func buildRemotePullErrorCmd(opts PullErrorOptions) string {
	cmd := "gitmap pull-errors"
	if opts.RepoSlug != "" && opts.RepoSlug != "all" {
		cmd += " " + opts.RepoSlug
	}
	if opts.IsJSON {
		cmd += " --json"
	}
	if opts.Limit > 0 && opts.Limit != 50 {
		cmd += fmt.Sprintf(" --limit %d", opts.Limit)
	}
	return cmd
}

func executeLocalPullErrors(opts PullErrorOptions) error {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return fmt.Errorf("open pull split db: %w", err)
	}
	defer db.Close()

	records, err := db.QueryLatestPullErrors(opts.RepoSlug, opts.Limit)
	if err != nil {
		return fmt.Errorf("query pull errors: %w", err)
	}

	return renderPullErrorsOutput(records, opts)
}

func renderPullErrorsOutput(records []store.PullErrorRecord, opts PullErrorOptions) error {
	if opts.IsJSON {
		return emitPullErrorsJSON(records)
	}
	if len(records) == 0 {
		renderNoPullErrors(opts.RepoSlug)
		return nil
	}
	renderPullErrorsList(records, opts.RepoSlug)
	return nil
}

func emitPullErrorsJSON(records []store.PullErrorRecord) error {
	if records == nil {
		records = []store.PullErrorRecord{}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
