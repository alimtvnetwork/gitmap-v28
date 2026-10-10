package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"strings"
	"time"
)

// RunMuseCLI routes the top-level `gitmap muse` commands.
func RunMuseCLI(args []string) error {
	if len(args) == 0 {
		printMuseUsage()
		return nil
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "install", "in", "i":
		return handleMuseInstallCommand(args[1:])
	case "status", "stat", "check":
		return handleMuseStatusCommand()
	case "version", "-v", "--version":
		return handleMuseVersionCommand()
	default:
		return handleMuseDirectOrUsage(args)
	}
}

func handleMuseDirectOrUsage(args []string) error {
	first := args[0]
	if strings.HasPrefix(first, "-") {
		return handleMuseInstallCommand(args)
	}

	printMuseUsage()
	return nil
}

func handleMuseInstallCommand(args []string) error {
	opts := parseMuseCliFlags(args)
	_, err := cmdinstall.RunMuseInstaller(opts)
	return err
}

func handleMuseStatusCommand() error {
	path, ver, isInstalled := cmdinstall.CheckExistingMuse()
	if !isInstalled {
		fmt.Printf("%s○ Meta Muse is not installed.%s\n", constants.ColorYellow, constants.ColorReset)
		fmt.Println("  Run 'gitmap muse install' to install Meta Muse.")
		return nil
	}

	fmt.Printf("%s● Meta Muse is installed%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Binary:  %s\n", path)
	fmt.Printf("  Version: %s\n", ver)
	return nil
}

func handleMuseVersionCommand() error {
	_, ver, isInstalled := cmdinstall.CheckExistingMuse()
	if !isInstalled {
		fmt.Println("Meta Muse is not installed.")
		return nil
	}

	fmt.Printf("Meta Muse version %s\n", ver)
	return nil
}

func parseMuseCliFlags(args []string) cmdinstall.MuseInstallOptions {
	opts := cmdinstall.MuseInstallOptions{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--force" || a == "-f" || a == "-y" || a == "--yes" {
			opts.Force = true
			continue
		}
		if a == "--dry-run" {
			opts.DryRun = true
			continue
		}
		if a == "--verify" {
			opts.Verify = true
			continue
		}
		if a == "--platform" && i+1 < len(args) {
			opts.Platform = cmdinstall.ResolveMusePlatform(args[i+1])
			i++
			continue
		}
		if a == "--timeout" && i+1 < len(args) {
			opts.Timeout = parseMuseTimeoutDuration(args[i+1])
			i++
			continue
		}
	}
	return opts
}

func parseMuseTimeoutDuration(val string) time.Duration {
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0
	}
	return d
}

func printMuseUsage() {
	fmt.Println("Usage: gitmap muse <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  install    Install Meta Muse AI developer agent for your platform")
	fmt.Println("  status     Check if Meta Muse binary is installed and registered in PATH")
	fmt.Println("  version    Display current Meta Muse version")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --platform <windows|linux|darwin>  Explicitly target platform installer")
	fmt.Println("  --force, -f                        Force reinstallation even if already present")
	fmt.Println("  --dry-run                          Display planned install command without executing")
	fmt.Println("  --verify                           Verify installation with 'muse --version'")
	fmt.Println("  --timeout <duration>               Set custom download/execution timeout")
}
