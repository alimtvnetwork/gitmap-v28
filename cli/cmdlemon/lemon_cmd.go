package cmdlemon

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// GitmapDaemonCmd represents the cobra command entry point for gitmap gitmap-daemon / daemon.
var GitmapDaemonCmd = &cobra.Command{
	Use:     "gitmap-daemon",
	Aliases: []string{"daemon", "lemon", "lm", "lemon-daemon"},
	Short:   "Background gitmap-daemon and email command service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunLemonCLI(args)
	},
}

// LemonCmd aliases GitmapDaemonCmd for backward compatibility.
var LemonCmd = GitmapDaemonCmd

// DaemonCmd aliases GitmapDaemonCmd for concise shorthand CLI registration.
var DaemonCmd = GitmapDaemonCmd

// RunLemonCLI routes subcommands for gitmap daemon and legacy lemon.
func RunLemonCLI(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printLemonHelp()
		return nil
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "run", "start-foreground":
		return runLemonDaemonCmd(args[1:])
	case "status", "st", "info":
		return runLemonStatusCmd(args[1:])
	case "poll", "once":
		return runLemonPollOnceCmd(args[1:])
	case "service", "svc":
		return runLemonServiceCmd(args[1:])
	case "enable", "on":
		return runDaemonEnableCmd()
	case "disable", "off":
		return runDaemonDisableCmd()
	default:
		return apperror.NewValidationError(fmt.Sprintf("unknown lemon subcommand %q (gitmap-daemon); see 'gitmap daemon help'", args[0]))
	}
}

func printLemonHelp() {
	fmt.Println("🤖 GitMap Daemon (gitmap-daemon / Lemon) & Email Command Service")
	fmt.Println("Usage:")
	fmt.Println("  gitmap daemon run [--interval 8s] [--once] [--verbose] [--force]")
	fmt.Println("  gitmap daemon status [--json]")
	fmt.Println("  gitmap daemon poll")
	fmt.Println("  gitmap daemon enable")
	fmt.Println("  gitmap daemon disable")
	fmt.Println("  gitmap daemon service <install|uninstall|start|stop|restart|status>")
	fmt.Println("Aliases:")
	fmt.Println("  gitmap daemon, gitmap gitmap-daemon, gitmap lemon, gitmap lm")
}

func runLemonDaemonCmd(args []string) error {
	interval := parseIntervalFlag(args)
	isOnce := hasFlag(args, "--once")
	isVerbose := hasFlag(args, "--verbose", "-v")
	isForce := hasFlag(args, "--force", "-f")

	if !isForce && !IsDaemonEnabled() {
		fmt.Println("Notice: gitmap-daemon is disabled in settings (daemon.enabled=false). Run 'gitmap daemon enable' or 'gitmap settings set daemon.enabled true' to enable.")
		return nil
	}

	cfg := DaemonConfig{
		Interval:  interval,
		IsOnce:    isOnce,
		IsVerbose: isVerbose,
		PollFunc: func(ctx context.Context) error {
			if isVerbose {
				fmt.Printf("[%s] 🤖 gitmap-daemon polled cycle completed\n", time.Now().Format(time.RFC3339))
			}
			return nil
		},
	}

	daemon := NewLemonDaemon(cfg)
	if !isOnce {
		fmt.Printf("🤖 Starting gitmap-daemon (interval=%v, min-guard=5s). Press Ctrl+C to stop.\n", daemon.cfg.Interval)
	}

	return daemon.Run(context.Background())
}

func runLemonPollOnceCmd(args []string) error {
	return runLemonDaemonCmd(append([]string{"--once"}, args...))
}

func runDaemonEnableCmd() error {
	if err := SetDaemonEnabled(true); err != nil {
		return apperror.WrapSimple(err, "failed to enable daemon in settings")
	}
	fmt.Println("✔ gitmap-daemon is now ENABLED in settings (daemon.enabled=true)")
	return nil
}

func runDaemonDisableCmd() error {
	if err := SetDaemonEnabled(false); err != nil {
		return apperror.WrapSimple(err, "failed to disable daemon in settings")
	}
	fmt.Println("✔ gitmap-daemon is now DISABLED in settings (daemon.enabled=false)")
	return nil
}
