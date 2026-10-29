package cmdspace

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/linkexpander"
)

func runSpaceEmail(args []string) error {
	if len(args) == 0 {
		printSendEmailUIHelp()

		return nil
	}

	if args[0] == "ui" {
		return runSpaceSendEmailUI(args[1:])
	}

	if args[0] == "send" {
		return RunSendEmail(args[1:])
	}

	if args[0] == "expand" {
		return runSpaceEmailExpand(args[1:])
	}

	if hasHelpFlag(args) {
		printSendEmailUIHelp()

		return nil
	}

	printSendEmailUIHelp()

	return apperror.NewValidationError(fmt.Sprintf("unknown space email subcommand: %s", args[0]))
}

func runSpaceEmailExpand(args []string) error {
	if len(args) == 0 || hasHelpFlag(args) {
		printEmailExpandHelp()

		return nil
	}

	isJSON, cleanArgs := filterJSONFlag(args)
	if len(cleanArgs) == 0 {
		return apperror.NewValidationError("usage: gitmap space email expand <url-or-text> [--json]")
	}

	input := strings.Join(cleanArgs, " ")
	expandedText, links := linkexpander.ExpandText(input)
	if isJSON {
		return printExpandResultJSON(input, expandedText, links)
	}

	return printExpandResultText(expandedText, links)
}

func filterJSONFlag(args []string) (bool, []string) {
	isJSON := false
	filtered := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--json" {
			isJSON = true
			continue
		}
		filtered = append(filtered, a)
	}

	return isJSON, filtered
}

func printExpandResultJSON(input, expandedText string, links []linkexpander.ExpandedLink) error {
	res := linkexpander.ExpandResult{
		OriginalText: input,
		ExpandedText: expandedText,
		HasChanges:   expandedText != input,
		Count:        len(links),
		Links:        links,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(res)
}

func printExpandResultText(expandedText string, links []linkexpander.ExpandedLink) error {
	if len(links) == 0 {
		fmt.Println(expandedText)

		return nil
	}
	for _, l := range links {
		if l.IsExpanded {
			fmt.Printf("%s -> %s\n", l.Original, l.Expanded)
			continue
		}
		fmt.Printf("%s (unchanged)\n", l.Original)
	}

	return nil
}

func printEmailExpandHelp() {
	fmt.Print(`Usage: gitmap space email expand <url-or-text> [--json]

Expands URLs and partner redirect links (such as serviceactivation.in -> serviceactivation.google.com).

Flags:
  --json    Output structured JSON result
`)
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
