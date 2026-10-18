package cmdspace

import (
	"fmt"
)

func runSpaceEmail(args []string) error {
	if len(args) == 0 {
		printSendEmailUIHelp()

		return nil
	}

	if args[0] == "ui" {
		return runSpaceSendEmailUI(args[1:])
	}

	if hasHelpFlag(args) {
		printSendEmailUIHelp()

		return nil
	}

	printSendEmailUIHelp()

	return fmt.Errorf("unknown space email subcommand: %s", args[0])
}

func runSpaceSendEmailUI(args []string) error {
	if hasHelpFlag(args) {
		printSendEmailUIHelp()

		return nil
	}

	return RunSpaceSendEmailUI(args)
}

func hasUIFlag(args []string) bool {
	for _, a := range args {
		if a == "--ui" || a == "-ui" {
			return true
		}
	}

	return false
}

func filterUIFlag(args []string) []string {
	filtered := make([]string, 0, len(args))
	for _, a := range args {
		if a != "--ui" && a != "-ui" {
			filtered = append(filtered, a)
		}
	}

	return filtered
}
