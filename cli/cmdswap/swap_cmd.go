package cmdswap

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
)

// Run is the main entry point for gitmap swap-binary command.
func Run(args []string) error {
	if isHelpRequested(args) {
		helpdoc.Print("swap-binary")

		return nil
	}

	opts, parseErr := parseSwapFlags(args)
	if parseErr != nil {
		return parseErr
	}

	return executeSwapOperation(opts)
}

func isHelpRequested(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" {
			return true
		}
	}

	return false
}

func parseSwapFlags(args []string) (SwapOptions, error) {
	opts := SwapOptions{}

	for i := 0; i < len(args); i++ {
		consumed, err := handleSwapFlag(args, i, &opts)
		if err != nil {
			return opts, err
		}

		i += consumed
	}

	return opts, nil
}

func handleSwapFlag(args []string, idx int, opts *SwapOptions) (int, error) {
	arg := args[idx]

	if isBooleanSwapFlag(arg, opts) {
		return 0, nil
	}

	return parseValuedSwapFlag(args, idx, opts)
}

func isBooleanSwapFlag(arg string, opts *SwapOptions) bool {
	switch arg {
	case "--kill", "-k":
		opts.HasKill = true
		return true
	case "--reopen", "--restart":
		opts.HasRestart = true
		return true
	case "--force", "-f":
		opts.HasForce = true
		return true
	case "--script-only":
		opts.IsScriptOnly = true
		return true
	case "--json":
		opts.IsJSONOutput = true
		return true
	default:
		return false
	}
}

func parseValuedSwapFlag(args []string, idx int, opts *SwapOptions) (int, error) {
	arg := args[idx]
	if idx+1 >= len(args) {
		return 0, apperror.NewValidationError(fmt.Sprintf("flag %s requires an argument", arg))
	}

	val := args[idx+1]
	switch arg {
	case "--from":
		opts.FromPath = val
		return 1, nil
	case "--to", "--target", "-t":
		opts.ToPath = val
		return 1, nil
	case "--pid":
		return parsePIDFlag(val, opts)
	default:
		return 0, nil
	}
}

func parsePIDFlag(val string, opts *SwapOptions) (int, error) {
	pid, err := strconv.Atoi(val)
	if err != nil || pid <= 0 {
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid pid value: %s", val))
	}

	opts.TargetPID = pid

	return 1, nil
}

func executeSwapOperation(opts SwapOptions) error {
	result, execErr := ExecuteSwap(opts)
	if execErr != nil {
		return execErr
	}

	output, renderErr := RenderResultOutput(result, opts.IsJSONOutput)
	if renderErr != nil {
		return renderErr
	}

	fmt.Println(output)

	return nil
}
