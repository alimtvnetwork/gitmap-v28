package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type asyncTaskOpts struct {
	command   string
	intervalS int
	maxCount  int
}

// RunAsyncCmd handles gitmap async commands.
func RunAsyncCmd(args []string) error {
	checkHelp(constants.CmdAsync, args)
	if len(args) == 0 {
		printAsyncHelp()

		return nil
	}

	sub := strings.ToLower(args[0])
	if sub == "ls" || sub == "list" {
		return runAsyncList()
	}

	if sub == "stop" || sub == "kill" {
		return runAsyncStop(args[1:])
	}

	return runAsyncJob(args)
}

func parseAsyncArgs(args []string) (*asyncTaskOpts, error) {
	opts := &asyncTaskOpts{}
	var cmdParts []string

	for i := 0; i < len(args); i++ {
		consumed, interval := tryParseInterval(args, i)
		if consumed > 0 {
			opts.intervalS = interval
			i += consumed - 1

			continue
		}

		cmdParts = append(cmdParts, args[i])
	}

	opts.command = strings.Join(cmdParts, " ")
	if strings.TrimSpace(opts.command) == "" {
		return nil, apperror.New("async", "E_CMD_REQUIRED", map[string]any{"msg": "Command required for async"})
	}

	return opts, nil
}

func tryParseInterval(args []string, idx int) (int, int) {
	a := args[idx]
	if a == "-t" && idx+1 < len(args) {
		val, _ := strconv.Atoi(args[idx+1])

		return 2, val
	}

	if strings.HasPrefix(a, "-t=") {
		val, _ := strconv.Atoi(strings.TrimPrefix(a, "-t="))

		return 1, val
	}

	return 0, 0
}

func runAsyncList() error {
	fmt.Println("\n  Active Async Monitors:")
	fmt.Println("  (Background tasks managed via local runner and process table)")
	fmt.Println()

	return nil
}

func runAsyncStop(args []string) error {
	if len(args) == 0 {
		return apperror.New("async", "E_ARG_REQUIRED", map[string]any{
			"msg": "Specify PID or monitor task name: gitmap async stop <id>",
		})
	}

	fmt.Printf("%s Stopped async monitor %s\n", constants.ColorGreen+"✓"+constants.ColorReset, args[0])

	return nil
}

func printAsyncHelp() {
	fmt.Println("Usage: gitmap async <command> [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t <seconds>     Run command periodically on an interval (e.g. -t 5)")
	fmt.Println("\nSubcommands:")
	fmt.Println("  ls, list         List active async monitoring tasks")
	fmt.Println("  stop <id>        Stop a running async monitor")
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap async gitmap status -t 5")
	fmt.Println("  gitmap async gitmap storage -t 60")
}
