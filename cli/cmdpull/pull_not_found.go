package cmdpull

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func handlePullTargetNotFound(opts pullOptions) {
	if len(opts.slug) > 0 {
		handlePullSlugNotFound(opts.slug)

		return
	}
	if len(opts.group) > 0 {
		handlePullGroupNotFound(opts.group)

		return
	}
	if opts.all {
		handlePullNoReposFound()

		return
	}

	cliexit.HandleError(nil, int(cliexit.ExitCodeNotFound))
}

func handlePullSlugNotFound(slug string) {
	arrow := resolveSubArrow()
	fmt.Fprintf(os.Stderr, "    %s%s%s %serror: repository '%s' not found%s\n",
		constants.ColorRed, arrow, constants.ColorReset,
		constants.ColorBold, slug, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "      No tracked repository matching '%s' was found in the database or scan records.\n\n", slug)
	printPullSuggestions(slug)

	cliexit.HandleError(nil, int(cliexit.ExitCodeNotFound))
}

func printPullSuggestions(slug string) {
	fmt.Fprintf(os.Stderr, "    %sSuggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	if hint := resolveCommandHint(slug); len(hint) > 0 {
		fmt.Fprintf(os.Stderr, "      • Did you mean '%s'?\n", hint)
	}
	fmt.Fprintf(os.Stderr, "      • Run 'gitmap search %s' or 'gitmap ls' to search tracked repositories\n", slug)
	printCwdPullSuggestion()
	fmt.Fprintf(os.Stderr, "      • Run 'gitmap pull --all' to pull all tracked repositories\n\n")
}

func printCwdPullSuggestion() {
	if isGitRepoCWD() {
		fmt.Fprintf(os.Stderr, "      • Run 'gitmap pull' without arguments to pull current repository\n")

		return
	}

	fmt.Fprintf(os.Stderr, "      • Run 'gitmap pull' inside a git repository directory\n")
}

func handlePullGroupNotFound(group string) {
	arrow := resolveSubArrow()
	fmt.Fprintf(os.Stderr, "    %s%s%s %serror: group '%s' not found%s\n",
		constants.ColorRed, arrow, constants.ColorReset,
		constants.ColorBold, group, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "      No tracked repositories found in group '%s'.\n\n", group)
	fmt.Fprintf(os.Stderr, "    %sSuggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "      • Run 'gitmap group ls' to list available groups\n")
	fmt.Fprintf(os.Stderr, "      • Run 'gitmap pull --all' to pull all tracked repositories\n\n")

	cliexit.HandleError(nil, int(cliexit.ExitCodeNotFound))
}

func handlePullNoReposFound() {
	arrow := resolveSubArrow()
	fmt.Fprintf(os.Stderr, "    %s%s%s %serror: no tracked repositories found%s\n",
		constants.ColorRed, arrow, constants.ColorReset,
		constants.ColorBold, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "      The repository database is empty.\n\n")
	fmt.Fprintf(os.Stderr, "    %sSuggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "      • Run 'gitmap scan' to discover and register repositories on this machine\n\n")

	cliexit.HandleError(nil, int(cliexit.ExitCodeNotFound))
}

func handlePushTargetNotFound(opts pushOptions) {
	if len(opts.slug) > 0 {
		handlePushSlugNotFound(opts.slug)

		return
	}
	if len(opts.group) > 0 {
		handlePullGroupNotFound(opts.group)

		return
	}
	if opts.all {
		handlePullNoReposFound()

		return
	}

	cliexit.HandleError(nil, int(cliexit.ExitCodeNotFound))
}

func handlePushSlugNotFound(slug string) {
	arrow := resolveSubArrow()
	fmt.Fprintf(os.Stderr, "    %s%s%s %serror: repository '%s' not found%s\n",
		constants.ColorRed, arrow, constants.ColorReset,
		constants.ColorBold, slug, constants.ColorReset)
	fmt.Fprintf(os.Stderr, "      No tracked repository matching '%s' was found in the database or scan records.\n\n", slug)
	printPushSuggestions(slug)

	cliexit.HandleError(nil, int(cliexit.ExitCodeNotFound))
}

func printPushSuggestions(slug string) {
	fmt.Fprintf(os.Stderr, "    %sSuggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	if hint := resolveCommandHint(slug); len(hint) > 0 {
		fmt.Fprintf(os.Stderr, "      • Did you mean '%s'?\n", hint)
	}
	fmt.Fprintf(os.Stderr, "      • Run 'gitmap search %s' or 'gitmap ls' to search tracked repositories\n", slug)
	printCwdPushSuggestion()
	fmt.Fprintf(os.Stderr, "      • Run 'gitmap push --all' to push all tracked repositories\n\n")
}

func printCwdPushSuggestion() {
	if isGitRepoCWD() {
		fmt.Fprintf(os.Stderr, "      • Run 'gitmap push' without arguments to push current repository\n")

		return
	}

	fmt.Fprintf(os.Stderr, "      • Run 'gitmap push' inside a git repository directory\n")
}

func resolveCommandHint(slug string) string {
	lower := strings.ToLower(slug)
	if hint := resolvePipelineCommandHint(lower); len(hint) > 0 {
		return hint
	}

	return resolveGeneralCommandHint(lower)
}

func resolvePipelineCommandHint(cmd string) string {
	switch cmd {
	case "pe":
		return "gitmap pe (view pipeline error logs)"
	case "pd":
		return "gitmap pd (view pipeline details and runner timings)"
	case "pipeline":
		return "gitmap pipeline (CI/CD pipeline diagnostics)"
	case "hist", "history":
		return "gitmap history (view commit/pipeline history)"
	case "eta", "waittime":
		return "gitmap eta (view remaining wait time)"
	default:
		return ""
	}
}

func resolveGeneralCommandHint(cmd string) string {
	switch cmd {
	case "st", "status":
		return "gitmap status (view repository status)"
	case "pa", "pull-all":
		return "gitmap pull-all (pull all tracked repositories)"
	case "aef", "fix":
		return "gitmap aef (feed errors to Antigravity IDE)"
	case "agy":
		return "gitmap agy (Google Antigravity integration)"
	default:
		return ""
	}
}
