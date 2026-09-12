package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type serverCmdOptions struct {
	Target   string
	Command  string
	IsSudo   bool
	IsScript bool
	Config   string
	Exclude  string
	Parallel int
	DryRun   bool
}

// runServerCmd dispatches remote server commands across cluster and SSH nodes.
func runServerCmd(args []string) error {
	if len(args) == 0 || hasHelpFlag(args) {
		printServerCmdUsage()

		return nil
	}

	opts, err := parseServerCmdFlags(args)
	if err != nil {
		return err
	}

	return executeServerCmd(opts)
}

func parseServerCmdFlags(args []string) (serverCmdOptions, error) {
	fs := flag.NewFlagSet("server-cmd", flag.ContinueOnError)
	var opts serverCmdOptions
	fs.BoolVar(&opts.IsSudo, "sudo", false, "Execute remote command with sudo escalation")
	fs.BoolVar(&opts.IsScript, "script", false, "Stage and execute on-the-fly script (/tmp/on-the-fly-cmd)")
	fs.StringVar(&opts.Config, "config", "", "Path to custom cluster topology config JSON")
	fs.StringVar(&opts.Exclude, "exclude", "", "Comma-separated list of nodes to exclude")
	fs.IntVar(&opts.Parallel, "parallel", 10, "Maximum concurrent remote worker routines")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Show target execution plan without connecting")

	if err := fs.Parse(args); err != nil {
		return opts, apperror.WrapSimple(err, "cmd.servercmd.parseFlags")
	}

	positional := fs.Args()
	if len(positional) < 2 {
		return opts, missingServerCmdArgsError(positional)
	}

	opts.Target = strings.ToLower(strings.TrimSpace(positional[0]))
	opts.Command = strings.Join(positional[1:], " ")

	return opts, nil
}

func missingServerCmdArgsError(args []string) error {
	msg := fmt.Sprintf("target and command required, got %d arguments; see 'gitmap server-cmd --help'", len(args))

	return apperror.NewWithDetails(
		"cmd.servercmd",
		"E4010",
		msg,
		"cmd.servercmd",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		nil,
	)
}

func printServerCmdUsage() {
	fmt.Println("Usage: gitmap server-cmd <target> \"<command>\" [flags]")
	fmt.Println()
	fmt.Println("Target Selectors:")
	fmt.Println("  all                    Execute across all cluster nodes (control plane and workers)")
	fmt.Println("  control, servers       Execute on control plane server nodes only")
	fmt.Println("  workers, clients       Execute on worker client nodes only")
	fmt.Println("  <node-id-or-ip>        Target a specific node by ID, display ID, or hostname")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --sudo                 Run command with non-interactive sudo escalation")
	fmt.Println("  --script               Stage and execute as /tmp/on-the-fly-cmd/script-<ts>.sh")
	fmt.Println("  --config <path>        Custom cluster topology JSON file (fallback to DB)")
	fmt.Println("  --exclude <nodes>      Comma-separated list of nodes to exclude")
	fmt.Println("  --parallel <N>         Maximum parallel connections (default: 10)")
	fmt.Println("  --dry-run              Display node plan without executing")
	fmt.Println("  -h, --help             Show this help message")
}
