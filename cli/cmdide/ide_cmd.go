package cmdide

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunIDE dispatches subcommands for the gitmap ide command group.
func RunIDE(args []string) error {
	if len(args) == 0 {
		return RunHelp()
	}
	return dispatchSubcommand(args[0], args[1:])
}

func dispatchSubcommand(subcmd string, args []string) error {
	switch subcmd {
	case "add":
		return runIDEAdd(args)
	case "sync":
		return runIDESync(args)
	case "remove", "rm":
		return runIDERemove(args)
	case "list", "ls":
		return runIDEList(args)
	case "status":
		return runIDEStatus(args)
	case "help", "-h", "--help":
		return RunHelp()
	default:
		return apperror.NewValidation("cmdide.RunIDE", "E1000", "unknown ide subcommand: "+subcmd)
	}
}

func defaultIDEOptions() IDEOptions {
	return IDEOptions{IsVSCodeTargeted: true, IsCursorTargeted: true, IsAntigravityTargeted: true, IsDesktopTargeted: true, IsForceCreate: true}
}

func parseIDEOptions(args []string) (IDEOptions, []string) {
	opts, remaining := defaultIDEOptions(), []string{}
	for i := 0; i < len(args); i++ {
		step, isFlag := parseSingleFlag(&opts, args, i)
		if isFlag {
			i += step
			continue
		}
		remaining = append(remaining, args[i])
	}
	return opts, remaining
}

func parseSingleFlag(opts *IDEOptions, args []string, idx int) (int, bool) {
	arg, hasNext := args[idx], idx+1 < len(args)
	if (arg == "-t" || arg == "--target" || arg == "-i" || arg == "--ide") && hasNext {
		applyTargetFilter(opts, args[idx+1])
		return 1, true
	}
	if (arg == "-e" || arg == "--exclude" || arg == "--skip-sync" || arg == "--exclude-sync") && hasNext {
		applyExcludeFilter(opts, args[idx+1])
		return 1, true
	}
	if (arg == "-d" || arg == "--dir") && hasNext {
		opts.TargetDirectory = args[idx+1]
		return 1, true
	}
	return checkBoolFlags(opts, arg)
}

func checkBoolFlags(opts *IDEOptions, arg string) (int, bool) {
	switch arg {
	case "-n", "--dry-run":
		opts.IsDryRun = true
	case "-j", "--json":
		opts.IsJSON = true
	case "-q", "--quiet":
		opts.IsQuiet = true
	case "--force-create":
		opts.IsForceCreate = true
	default:
		return 0, false
	}
	return 0, true
}

func applyTargetFilter(opts *IDEOptions, target string) {
	low := strings.ToLower(target)
	if low != TargetAll {
		opts.IsVSCodeTargeted = strings.Contains(low, TargetVSCode) || strings.Contains(low, "code")
		opts.IsCursorTargeted = strings.Contains(low, TargetCursor)
		opts.IsAntigravityTargeted = strings.Contains(low, TargetAntigravity) || strings.Contains(low, "agy")
		opts.IsDesktopTargeted = strings.Contains(low, TargetDesktop)
	}
}

func applyExcludeFilter(opts *IDEOptions, excluded string) {
	low := strings.ToLower(excluded)
	opts.IsVSCodeTargeted = opts.IsVSCodeTargeted && !strings.Contains(low, TargetVSCode) && !strings.Contains(low, "code")
	opts.IsCursorTargeted = opts.IsCursorTargeted && !strings.Contains(low, TargetCursor)
	opts.IsAntigravityTargeted = opts.IsAntigravityTargeted && !strings.Contains(low, TargetAntigravity) && !strings.Contains(low, "agy")
	opts.IsDesktopTargeted = opts.IsDesktopTargeted && !strings.Contains(low, TargetDesktop)
}
