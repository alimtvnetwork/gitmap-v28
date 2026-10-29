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
	instructionPath, isJSON, cleanArgs, parseErr := parseEmailExpandArgs(args)
	if parseErr != nil {
		return parseErr
	}
	if len(instructionPath) > 0 {
		return executeInstructionExpand(instructionPath, isJSON)
	}
	if len(cleanArgs) == 0 {
		return apperror.NewValidationError("usage: gitmap space email expand <url-or-text> [--instruction <path>] [--json]")
	}

	return runStandardEmailExpand(cleanArgs, isJSON)
}

func runStandardEmailExpand(cleanArgs []string, isJSON bool) error {
	input := strings.Join(cleanArgs, " ")
	expandedText, links := linkexpander.ExpandText(input)
	if isJSON {
		return printExpandResultJSON(input, expandedText, links)
	}

	return printExpandResultText(expandedText, links)
}

func executeInstructionExpand(path string, isJSON bool) error {
	res, execErr := linkexpander.ExecuteInstructionFile(path)
	if execErr != nil {
		return execErr
	}
	if isJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(res)
	}

	return printInstructionSummaryText(res)
}

func parseEmailExpandArgs(args []string) (string, bool, []string, error) {
	var instructionPath string
	isJSON := false
	var cleanArgs []string

	for i := 0; i < len(args); i++ {
		path, jsonFlag, consumed, err := parseOneExpandArg(args, i)
		if err != nil {
			return "", false, nil, err
		}
		applyParsedArg(path, jsonFlag, args[i], consumed, &instructionPath, &isJSON, &cleanArgs)
		i += consumed
	}

	return instructionPath, isJSON, cleanArgs, nil
}

func applyParsedArg(path string, jsonFlag bool, raw string, consumed int, outPath *string, outJSON *bool, clean *[]string) {
	if len(path) > 0 {
		*outPath = path
	}
	if jsonFlag {
		*outJSON = true
	}
	if consumed == 0 && !jsonFlag && !isInstructionArg(raw) {
		*clean = append(*clean, raw)
	}
}

func isInstructionArg(arg string) bool {
	return strings.HasPrefix(arg, "--instruction=") || strings.HasPrefix(arg, "-i=")
}

func parseOneExpandArg(args []string, i int) (string, bool, int, error) {
	a := args[i]
	if a == "--json" {
		return "", true, 0, nil
	}
	if a == "--instruction" || a == "-i" {
		return parseInstructionValue(args, i, a)
	}
	if strings.HasPrefix(a, "--instruction=") || strings.HasPrefix(a, "-i=") {
		return extractInlineInstruction(a), false, 0, nil
	}

	return "", false, 0, nil
}

func parseInstructionValue(args []string, i int, flagName string) (string, bool, int, error) {
	if i+1 >= len(args) {
		return "", false, 0, apperror.NewValidationError("flag " + flagName + " requires a value")
	}

	return args[i+1], false, 1, nil
}

func extractInlineInstruction(arg string) string {
	_, val, _ := strings.Cut(arg, "=")

	return val
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
	fmt.Print(`Usage: gitmap space email expand <url-or-text> [flags]

Expands URLs and partner redirect links (such as serviceactivation.in -> serviceactivation.google.com).

Flags:
  --instruction, -i <path>   Execute link expansion from JSON instruction file
  --json                     Output structured JSON result
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
