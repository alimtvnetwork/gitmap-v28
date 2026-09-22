// Package cmdantigravity — prompt_help.go renders help documentation for EDI prompt commands.
package cmdantigravity

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintAgyPromptHelp prints comprehensive CLI usage for the Antigravity prompt command suite.
func PrintAgyPromptHelp() error {
	printPromptHelpHeader()
	printPromptHelpSyntax()
	printPromptHelpBehavior()
	printAvailablePromptTemplates()
	printPromptHelpCommands()
	printPromptHelpFlags()
	printPromptHelpExamples()

	return nil
}

func printPromptHelpHeader() {
	fmt.Println()
	fmt.Printf("  %s💬 GitMap EDI — Antigravity Dynamic Prompt Engine%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func printPromptHelpSyntax() {
	fmt.Println("  SYNTAX:")
	fmt.Println("    gitmap agy prompt -name(n) <prompt-name> -txt(t) \"<text>\" [--prefix(pf)/suffix(sf)]")
	fmt.Println("    gitmap agy prompt-with-name(pwn) <prompt-name> -txt(t) \"<text>\" [--prefix(pf)/suffix(sf)]")
	fmt.Println("    gitmap agy prompt-txt(pt) \"<text>\" [--prefix(pf)/suffix(sf)]")
	fmt.Println("    gitmap agy prompt-project(p) <projectStartsWithName> -name(n) <name> -txt(t) \"<text>\"")
	fmt.Println("    gitmap agy prompt ls (or: gitmap agy list prompts)")
	fmt.Println()
}

func printPromptHelpBehavior() {
	fmt.Println("  DYNAMIC PROMPTING & CONVENTIONS:")
	fmt.Println("    • Single Prompt Injection: 'read-all' is NOT fixed or forced. It is one prompt template.")
	fmt.Println("      You can select and inject ANY prompt template dynamically by name as you wish.")
	fmt.Println("    • Git Repo Awareness: Running from inside a git repo automatically targets that project.")
	fmt.Println("    • Auto-Registration: Automatically registers the repository in Antigravity if missing.")
	fmt.Println("    • Flexible Assembly: By default, prefixes the template before text (--prefix / --pf).")
	fmt.Println("      Use --suffix / --sf to place custom text before the template.")
	fmt.Println()
}

func printAvailablePromptTemplates() {
	templates, err := cmdprompttemplate.LoadTemplates()
	if err != nil || len(templates) == 0 {
		return
	}

	fmt.Printf("  %sAVAILABLE PROMPT TEMPLATES & INJECTION EXAMPLES:%s\n",
		constants.ColorWhite, constants.ColorReset)
	for _, tpl := range templates {
		desc := tpl.Description
		if len(desc) == 0 {
			desc = "Reusable prompt template"
		}
		fmt.Printf("    • %s%-12s%s : %s\n", constants.ColorYellow, tpl.Name, constants.ColorReset, desc)
		fmt.Printf("      %sCommand:%s gitmap agy prompt -n %s -t \"<your custom instructions>\"\n",
			constants.ColorDim, constants.ColorReset, tpl.Name)
	}
	fmt.Println()
}

func printPromptHelpCommands() {
	fmt.Println("  COMMANDS:")
	fmt.Println("    prompt, pr           Dispatch single prompt to current repo (auto-registers if missing)")
	fmt.Println("    prompt-with-name     Dispatch named prompt template by argument (same as prompt -name)")
	fmt.Println("    prompt-txt, pt       Dispatch direct prompt text to current repo")
	fmt.Println("    prompt-project, p    Target specific Antigravity project by prefix and dispatch prompt")
	fmt.Println("    prompt ls            List all available prompt templates in a formatted table")
	fmt.Println()
}

func printPromptHelpFlags() {
	fmt.Println("  FLAGS:")
	fmt.Println("    -n, --name string    Prompt template name (e.g. read-all, is-done; any listed template)")
	fmt.Println("    -t, --txt string     Additional prompt instructions or task details")
	fmt.Println("    --prefix, --pf       Place template before text separated by 2 newlines (default)")
	fmt.Println("    --suffix, --sf       Place template after text separated by 2 newlines")
	fmt.Println("    -h, --help           Show this comprehensive prompt help")
	fmt.Println()
}

func printPromptHelpExamples() {
	fmt.Println("  EXAMPLES:")
	fmt.Println("    # Inject 'read-all' context ingestion prompt:")
	fmt.Println("    gitmap agy prompt -n read-all -t \"Defensively load memory, specs, and pending plans\"")
	fmt.Println()
	fmt.Println("    # Inject 'is-done' task completion verification prompt:")
	fmt.Println("    gitmap agy prompt -n is-done -t \"Verify all unit tests pass without regressions\"")
	fmt.Println()
	fmt.Println("    # Inject custom text prompt directly:")
	fmt.Println("    gitmap agy prompt -t \"Refactor database split-db to use PascalCase tables\"")
	fmt.Println()
	fmt.Println("    # Using prompt-with-name alias:")
	fmt.Println("    gitmap agy pwn is-done -t \"Check pipeline error logs and execute 4-part RCA\"")
	fmt.Println()
	fmt.Println("    # Target specific project:")
	fmt.Println("    gitmap agy prompt-project my-app -n is-done -t \"Verify release ceremony\"")
	fmt.Println()
	fmt.Println("    # List all available prompts:")
	fmt.Println("    gitmap agy prompt ls")
	fmt.Println()
}

// IsHelpArg checks whether CLI arguments request help.
func IsHelpArg(args []string) bool {
	hasArgs := len(args) > 0
	if hasArgs == false {
		return false
	}
	tok := strings.ToLower(strings.TrimSpace(args[0]))

	return tok == "help" || tok == "--help" || tok == "-h"
}
