package cmdhooks

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"
	"sort"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const hooksInstallHelp = `hooks install — install the commit-msg guard hook across repos.

Description:
  Installs the house commit-msg guard (no colons in the subject line)
  into every repo in the workdir. Safe by default: --dry-run prints the
  plan; a real install only happens with --yes.

Usage:
  gitmap hooks install --all [--workdir <dir>] [--dry-run] [--yes]
  gitmap hkin --all [--yes]

Flags:
  --all           required: target all repos in the workdir
  --workdir <dir> repo workdir (default: $GITMAP_WORKDIR or ~/workspace/repos)
  --dry-run       print the install plan without writing (default: true)
  --yes           actually install the hooks
  --help, -h      show this help

Example:
  gitmap hooks install --all --dry-run
`

// hookPlan is the per-repo install decision.
type hookPlan struct {
	repo   string
	action string
	detail string
}

// RunHooksInstall implements `gitmap hooks install --all` (alias: hkin).
func RunHooksInstall(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(hooksInstallHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("hooks install", flag.ContinueOnError)
	allFlag := fs.Bool("all", false, "target all repos in the workdir")
	workdirFlag := fs.String("workdir", "", "repo workdir")
	dryRun := fs.Bool("dry-run", true, "print the plan without installing")
	yes := fs.Bool("yes", false, "actually install the hooks")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(hooksInstallHelp)
		return apperror.WrapSimple(parseErr, "hooks.install.flags")
	}
	if !*allFlag {
		fmt.Print(hooksInstallHelp)
		return appfault.NewValidationError("missing required flag: --all")
	}

	workdir := resolveHooksWorkdir(*workdirFlag)
	repos := listHookRepos(workdir)
	if len(repos) == 0 {
		fmt.Printf("hooks install: no git repos found in %s\n", workdir)
		return nil
	}

	plans := make([]hookPlan, 0, len(repos))
	for _, repo := range repos {
		plans = append(plans, planHookInstall(workdir, repo))
	}

	isDryRun := *dryRun && !*yes
	installed := 0
	skipped := 0
	already := 0
	for _, p := range plans {
		switch p.action {
		case "install":
			installed++
		case "already":
			already++
		default:
			skipped++
		}
	}
	fmt.Printf("hooks install plan: %d repos (%d would install, %d already guarded, %d skipped)\n",
		len(plans), installed, already, skipped)
	for _, p := range plans {
		fmt.Printf("  %-9s %-32s %s\n", p.action, p.repo, p.detail)
	}

	if isDryRun {
		fmt.Println("dry-run: no hooks were installed. re-run with --yes to install.")
		return nil
	}

	written, writeErr := applyHookInstalls(plans, workdir)
	if writeErr != nil {
		return writeErr
	}
	fmt.Printf("hooks install: installed commit-msg guard in %d/%d repos (%d already had it, %d skipped).\n",
		written, len(plans), already, skipped)

	return nil
}

// resolveHooksWorkdir mirrors the fleet workdir resolution.
func resolveHooksWorkdir(flagVal string) string {
	if strings.TrimSpace(flagVal) != "" {
		return flagVal
	}
	if envVal := strings.TrimSpace(os.Getenv("GITMAP_WORKDIR")); envVal != "" {
		return envVal
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "."
	}

	return filepath.Join(home, "workspace", "repos")
}

// listHookRepos enumerates workdir directories containing a .git dir.
func listHookRepos(workdir string) []string {
	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		return nil
	}
	var repos []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		full := filepath.Join(workdir, e.Name())
		if _, statErr := os.Stat(filepath.Join(full, ".git")); statErr != nil {
			continue
		}
		repos = append(repos, e.Name())
	}
	sort.Strings(repos)

	return repos
}

// planHookInstall decides the action for one repo.
func planHookInstall(workdir, repo string) hookPlan {
	hookPath := filepath.Join(workdir, repo, ".git", "hooks", "commit-msg")
	if isGuardInstalled(hookPath) {
		return hookPlan{repo: repo, action: "already", detail: "guard already installed"}
	}

	return hookPlan{repo: repo, action: "install", detail: ".git/hooks/commit-msg"}
}

// applyHookInstalls writes the guard hook for planned repos.
func applyHookInstalls(plans []hookPlan, workdir string) (int, *appfault.AppError) {
	written := 0
	for _, p := range plans {
		if p.action != "install" {
			continue
		}
		hookPath := filepath.Join(workdir, p.repo, ".git", "hooks", "commit-msg")
		if isGuardInstalled(hookPath) {
			continue
		}
		if writeErr := os.WriteFile(hookPath, []byte(commitMsgGuard), 0755); writeErr != nil {
			return written, apperror.WrapSimple(writeErr, "hooks.install.write: "+p.repo)
		}
		written++
	}

	return written, nil
}

// isGuardInstalled reports whether our guard already occupies hookPath.
func isGuardInstalled(hookPath string) bool {
	data, readErr := os.ReadFile(hookPath)
	if readErr != nil {
		return false
	}

	return strings.Contains(string(data), guardMarker)
}
