package cmdos

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runOSFixLink inspects and repairs broken symlinks and shared directories.
func runOSFixLink(args []string) error {
	checkHelp("fix-link", args)
	opts, paths := parseFixLinkArgs(args)

	if len(paths) == 0 {
		paths = resolveDefaultFixLinkPaths()
	}

	results, err := executeFixLinkRuns(paths, opts)
	if err != nil {
		return err
	}

	return renderLinkResults(results, opts)
}

func parseFixLinkArgs(args []string) (FixLinkOptions, []string) {
	opts := FixLinkOptions{}
	var paths []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isTargetFlag(arg) && i+1 < len(args) {
			opts.TargetOverride = expandHome(args[i+1])
			i++

			continue
		}

		if applyFlagIfMatch(arg, &opts) {
			continue
		}

		if !strings.HasPrefix(arg, "-") {
			paths = append(paths, expandHome(arg))
		}
	}

	return opts, paths
}

func isTargetFlag(arg string) bool {
	return arg == "--target" || arg == "-t"
}

func applyFlagIfMatch(arg string, opts *FixLinkOptions) bool {
	switch arg {
	case "--force", "-f":
		opts.IsForce = true

		return true
	case "--dry-run", "-n":
		opts.IsDryRun = true

		return true
	case "--recursive", "-r":
		opts.IsRecursive = true

		return true
	case constants.FlagJSON:
		opts.IsJSON = true

		return true
	default:
		return false
	}
}

func executeFixLinkRuns(paths []string, opts FixLinkOptions) ([]LinkResult, error) {
	var allResults []LinkResult

	for _, p := range paths {
		items, err := inspectAndRepairPath(p, opts)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, items...)
	}

	return allResults, nil
}
