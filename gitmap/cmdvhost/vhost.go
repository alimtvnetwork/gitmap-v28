package cmdvhost

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func dispatchVHostOpsSub(sub string, args []string) (bool, error) {
	switch sub {
	case "rm", "remove", "delete", "del":
		return true, runNginxRm(args)
	case "enable", "en":
		return true, runVHostEnable(args)
	case "disable", "dis":
		return true, runVHostDisable(args)
	case "test", "check", "t":
		return true, runVHostTest(args)
	case "reload", "restart", "r":
		return true, runVHostReload(args)
	default:
		return false, nil
	}
}

func dispatchVHostSub(sub string, args []string) error {
	switch sub {
	case "list", "ls":
		return runVHostList(args)
	case "create", "add", "new":
		return runVHostCreate(args)
	}

	isHandled, err := dispatchVHostOpsSub(sub, args)
	if isHandled {
		return err
	}

	return apperror.NewValidationError("unknown vhost subcommand: " + sub)
}

// runVHost handles CLI routing for the 'vhost' top-level command.
func runVHost(args []string) error {
	checkHelp("vhost", args)
	isNoArgs := len(args) == 0
	if isNoArgs {
		return runVHostList([]string{})
	}

	sub := args[0]
	rest := args[1:]

	return dispatchVHostSub(sub, rest)
}

func printVHostRow(info VHostInfo) {
	status := "disabled"
	if info.IsEnabled {
		status = "enabled"
	}

	fmt.Printf("  %-25s %-12s %-10s %s\n", info.Domain, info.SiteType, status, info.ConfigPath)
}

func printVHostTable(vhosts []VHostInfo) {
	isEmpty := len(vhosts) == 0
	if isEmpty {
		fmt.Printf("No virtual hosts found.\n")

		return
	}

	fmt.Printf("\n  %-25s %-12s %-10s %s\n", "DOMAIN", "TYPE", "STATUS", "CONFIG PATH")
	fmt.Printf("  %-25s %-12s %-10s %s\n", "------", "----", "------", "-----------")
	for _, vh := range vhosts {
		printVHostRow(vh)
	}

	fmt.Println()
}

func runVHostList(args []string) error {
	opts := VHostOptions{}
	vhosts, err := ListVHosts(opts)
	if err != nil {
		return err
	}

	printVHostTable(vhosts)

	return nil
}

func runVHostEnable(args []string) error {
	isEmpty := len(args) < 1
	if isEmpty {
		return apperror.NewValidationError("domain name required: gitmap vhost enable <domain>")
	}

	domain := args[0]
	opts := VHostOptions{}
	err := EnableVHost(domain, opts)
	if err != nil {
		return err
	}

	fmt.Printf("%sVirtual host enabled: %s%s\n", constants.ColorGreen, domain, constants.ColorReset)

	return nil
}

func runVHostDisable(args []string) error {
	isEmpty := len(args) < 1
	if isEmpty {
		return apperror.NewValidationError("domain name required: gitmap vhost disable <domain>")
	}

	domain := args[0]
	opts := VHostOptions{}
	err := DisableVHost(domain, opts)
	if err != nil {
		return err
	}

	fmt.Printf("%sVirtual host disabled: %s%s\n", constants.ColorGreen, domain, constants.ColorReset)

	return nil
}

func isDryRunArg(arg string) bool {
	return arg == "--dry-run" || arg == "-n"
}

func parseDryRunOption(args []string) bool {
	for _, a := range args {
		if isDryRunArg(a) {
			return true
		}
	}

	return false
}

func runVHostTest(args []string) error {
	opts := VHostOptions{IsDryRun: parseDryRunOption(args)}
	out, err := TestNginxConfig(opts)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}

func runVHostReload(args []string) error {
	opts := VHostOptions{IsDryRun: parseDryRunOption(args)}
	err := ReloadNginx(opts)
	if err != nil {
		return err
	}

	fmt.Printf("%sNginx configuration reloaded successfully.%s\n", constants.ColorGreen, constants.ColorReset)

	return nil
}
