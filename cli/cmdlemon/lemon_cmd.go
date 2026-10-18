package cmdlemon

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// LemonCmd represents the cobra command entry point for gitmap lemon.
var LemonCmd = &cobra.Command{
	Use:     "lemon",
	Aliases: []string{"lm", "lemon-daemon"},
	Short:   "Background lemon daemon and email command service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunLemonCLI(args)
	},
}

// RunLemonCLI routes subcommands for gitmap lemon.
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
	default:
		return fmt.Errorf("unknown lemon subcommand %q; see 'gitmap lemon help'", args[0])
	}
}

func printLemonHelp() {
	fmt.Println("🍋 GitMap Lemon Daemon & Email Command Service")
	fmt.Println("Usage:")
	fmt.Println("  gitmap lemon run [--interval 8s] [--once] [--verbose]")
	fmt.Println("  gitmap lemon status [--json]")
	fmt.Println("  gitmap lemon poll")
	fmt.Println("  gitmap lemon service <install|uninstall|start|stop|restart|status>")
}

func runLemonDaemonCmd(args []string) error {
	interval := parseIntervalFlag(args)
	isOnce := hasFlag(args, "--once")
	isVerbose := hasFlag(args, "--verbose", "-v")

	cfg := DaemonConfig{
		Interval:  interval,
		IsOnce:    isOnce,
		IsVerbose: isVerbose,
		PollFunc: func(ctx context.Context) error {
			if isVerbose {
				fmt.Printf("[%s] 🍋 Lemon daemon polled cycle completed\n", time.Now().Format(time.RFC3339))
			}
			return nil
		},
	}

	daemon := NewLemonDaemon(cfg)
	if !isOnce {
		fmt.Printf("🍋 Starting Lemon daemon (interval=%v, min-guard=5s). Press Ctrl+C to stop.\n", daemon.cfg.Interval)
	}
	return daemon.Run(context.Background())
}

func runLemonStatusCmd(args []string) error {
	isJSON := hasFlag(args, "--json")
	var state DaemonState
	data, err := os.ReadFile(StateFilePath)
	if err == nil {
		_ = json.Unmarshal(data, &state)
	}

	if isJSON {
		out, _ := json.MarshalIndent(state, "", "  ")
		fmt.Println(string(out))
		return nil
	}

	fmt.Println("🍋 GitMap Lemon Daemon Status")
	fmt.Printf("  Running:      %v\n", state.IsRunning)
	fmt.Printf("  PID:          %d\n", state.Pid)
	fmt.Printf("  Poll Count:   %d\n", state.PollCount)
	fmt.Printf("  Last Polled:  %s\n", state.LastPollTime.Format(time.RFC3339))
	return nil
}

func runLemonPollOnceCmd(args []string) error {
	return runLemonDaemonCmd(append([]string{"--once"}, args...))
}

func runLemonServiceCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing service action: install, uninstall, start, stop, restart, status")
	}
	mgr := NewServiceManager()
	details, err := mgr.Execute(ServiceActionType(args[0]))
	if err != nil {
		return err
	}
	if details != nil {
		fmt.Printf("✔ Service action %q applied (name=%s, status=%s, platform=%s)\n",
			args[0], details.Name, details.StatusType, details.Platform)
	}
	return nil
}

func parseIntervalFlag(args []string) time.Duration {
	for i, arg := range args {
		if (arg == "--interval" || arg == "-i") && i+1 < len(args) {
			d, err := time.ParseDuration(args[i+1])
			if err == nil {
				return d
			}
		}
		if strings.HasPrefix(arg, "--interval=") {
			val := strings.TrimPrefix(arg, "--interval=")
			d, err := time.ParseDuration(val)
			if err == nil {
				return d
			}
		}
	}
	return DefaultPollingInterval
}

func hasFlag(args []string, flags ...string) bool {
	for _, arg := range args {
		for _, f := range flags {
			if arg == f {
				return true
			}
		}
	}
	return false
}
