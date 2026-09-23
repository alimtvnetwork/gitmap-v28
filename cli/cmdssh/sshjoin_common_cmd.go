package cmdssh

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var (
	sjcPasswordFlag string
	sjcPortFlag     int
	sjcJSONFlag     bool
	sjcDryRunFlag   bool
)

// SSHJoinCommonCmd represents the 'gitmap ssh-join-common' command.
var SSHJoinCommonCmd = &cobra.Command{
	Use:     "ssh-join-common <username> <ip-tokens...> [flags]",
	Aliases: []string{"sjc", "join-common", "ssh-join-c"},
	Short:   "Batch join SSH machines using common credentials and shorthand octet notation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSSHJoinCommonCLI(args)
	},
}

func init() {
	SSHJoinCommonCmd.Flags().StringVarP(&sjcPasswordFlag, "pass", "p", "", "Common password for SSH targets")
	SSHJoinCommonCmd.Flags().IntVar(&sjcPortFlag, "port", 22, "Target SSH port (default 22)")
	SSHJoinCommonCmd.Flags().BoolVar(&sjcJSONFlag, "json", false, "Output results in JSON format")
	SSHJoinCommonCmd.Flags().BoolVarP(&sjcDryRunFlag, "dry-run", "d", false, "Preview expanded targets without connecting")
}

func parsePortArg(val string, fallback int) int {
	if p, err := strconv.Atoi(val); err == nil && p > 0 {
		return p
	}
	return fallback
}

func parseSJCFlags(args []string) (SSHCommonJoinOptions, error) {
	opts := SSHCommonJoinOptions{
		Port: 22,
	}
	var cleanArgs []string
	idx := 0
	for idx < len(args) {
		a := args[idx]
		if (a == "--pass" || a == "-p") && idx+1 < len(args) {
			opts.Password = args[idx+1]
			idx += 2
			continue
		}
		if a == "--port" && idx+1 < len(args) {
			opts.Port = parsePortArg(args[idx+1], opts.Port)
			idx += 2
			continue
		}
		if a == "--json" {
			opts.IsJSON = true
			idx++
			continue
		}
		if a == "--dry-run" || a == "-d" {
			opts.DryRun = true
			idx++
			continue
		}
		if a == "-h" || a == "--help" || a == "help" {
			opts.IsHelpRequested = true
			idx++
			continue
		}
		if strings.HasPrefix(a, "--pass=") {
			opts.Password = strings.TrimPrefix(a, "--pass=")
			idx++
			continue
		}
		cleanArgs = append(cleanArgs, a)
		idx++
	}

	if opts.IsHelpRequested {
		printSJCUsage()
		return opts, nil
	}

	if len(cleanArgs) < 2 {
		printSJCUsage()
		return opts, apperror.NewValidationError("usage: gitmap ssh-join-common <username> <ip-tokens...> [--pass <password>]")
	}

	opts.Username = cleanArgs[0]
	opts.RawIPs = cleanArgs[1:]
	return opts, nil
}

func printSJCUsage() {
	usage := `gitmap ssh-join-common (sjc) - Batch enroll multiple SSH nodes

Usage:
  gitmap ssh-join-common <username> <ip-tokens...> [flags]
  gitmap sjc <username> <ip-tokens...> [flags]

Flags:
  -p, --pass string      Common password for SSH targets
      --port int         Target SSH port (default 22)
  -d, --dry-run          Preview expanded targets without connecting
      --json             Output results in JSON format
  -h, --help             Show help for ssh-join-common

Examples:
  gitmap sjc administrator 192.168.1.3(w1),7(w2),12(w3) --pass Secret123
  gitmap sjc root 10.0.0.5(db1),6(db2),15(redis) -p Secret123 --dry-run
`
	fmt.Print(usage)
}

// RunSSHJoinCommonCLI executes batch SSH join from commandline arguments.
func RunSSHJoinCommonCLI(args []string) error {
	opts, err := parseSJCFlags(args)
	if err != nil {
		return err
	}
	if opts.IsHelpRequested {
		return nil
	}
	ctx := context.Background()
	_, appErr := ExecuteCommonJoin(ctx, os.Stdout, opts)
	if appErr != nil {
		return appErr
	}
	return nil
}
