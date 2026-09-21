// Package cmdantigravity — prompt_help.go renders help documentation for EDI prompt commands.
package cmdantigravity

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintAgyPromptHelp prints comprehensive CLI usage for the Antigravity prompt command suite.
func PrintAgyPromptHelp() error {
	printPromptHelpHeader()
	printPromptHelpSyntax()
	printPromptHelpBehavior()
	printPromptHelpCommands()
	printPromptHelpFlags()
	printPromptHelpExamples()

	return nil
}

func printPromptHelpHeader() {
	fmt.Println()
	fmt.Printf("  %s💬 GitMap EDI — Antigravity Prompt Engine%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func printPromptHelpSyntax() {
	fmt.Println("  SYNTAX:")
	fmt.Println("    gitmap agy prompt-project(p) <projectStartsWithName> -name(n) <name> -txt(t) \"<text>\" [--prefix(pf)/suffix(sf)]")
	fmt.Println("    gitmap agy prompt -name(n) <prompt-name> -txt(t) \"<text>\" [--prefix(pf)/suffix(sf)]")
	fmt.Println("    gitmap agy prompt-with-name(pwn) <prompt-name> -txt(t) \"<text>\" [--prefix(pf)/suffix(sf)]")
	fmt.Println("    gitmap agy prompt-txt(pt) \"<text>\" [--prefix(pf)/suffix(sf)]")
	fmt.Println("    gitmap agy prompt ls (or: gitmap agy list prompts)")
	fmt.Println()
}

func printPromptHelpBehavior() {
	fmt.Println("  BEHAVIOR & CONVENTIONS:")
	fmt.Println("    • Git Repo: Running from git repo does not need project name.")
	fmt.Println("    • Auto-Registration: Automatically registers repo in Antigravity if missing.")
	fmt.Println("    • Dual-Queue: Always adds default prompt to read all first, then queues current prompt.")
	fmt.Println("    • Default Name: If -name(n) is omitted, defaults to 'read-all' (same as prompt-with-name read-all).")
	fmt.Println("    • Formatting: By default, prefix template with 2 newlines (\\n\\n) before text. Use --suffix(sf) for post-text.")
	fmt.Println()
}

func printPromptHelpCommands() {
	fmt.Println("  COMMANDS:")
	fmt.Println("    prompt-project, p    Target Antigravity project by prefix name and dispatch prompt")
	fmt.Println("    prompt               Dispatch to current repo (auto-adds project, queues read-all first)")
	fmt.Println("    prompt-with-name     Dispatch named prompt template (same as prompt -name)")
	fmt.Println("    prompt-txt           Dispatch direct text (defaults to read-all prefix with 2 newlines)")
	fmt.Println("    prompt ls            List all available prompt templates in a formatted table")
	fmt.Println()
}

func printPromptHelpFlags() {
	fmt.Println("  FLAGS:")
	fmt.Println("    -n, --name string    Prompt template name (e.g. read-all, is-done; defaults to 'read-all')")
	fmt.Println("    -t, --txt string     Additional prompt text to include")
	fmt.Println("    --prefix, --pf       Place template before text separated by 2 newlines (default)")
	fmt.Println("    --suffix, --sf       Place template after text separated by 2 newlines")
	fmt.Println("    -h, --help           Show this prompt help")
	fmt.Println()
}

func printPromptHelpExamples() {
	fmt.Println("  EXAMPLES:")
	fmt.Println("    gitmap agy prompt-project my-app -name read-all -txt \"what we want to add here\" --prefix")
	fmt.Println("    gitmap agy prompt -name read-all -txt \"what we want to add here\" --prefix")
	fmt.Println("    gitmap agy prompt-with-name read-all -txt \"what we want to add here\" --prefix")
	fmt.Println("    gitmap agy prompt-txt \"what we want to add here\" --prefix")
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
