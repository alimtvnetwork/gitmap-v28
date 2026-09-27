// Package cmdagy — agy_rerun_help.go renders rich guidance and documentation for prompt rerun and replay workflows.
package cmdagy

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var agyRerunHelpCmd = &cobra.Command{
	Use:     "help",
	Aliases: []string{"man", "info", "usage"},
	Short:   "Show comprehensive help for rerun command",
	Run: func(cmd *cobra.Command, args []string) {
		RenderAgyRerunHelp()
	},
}

// IsRerunHelpRequested reports whether any token in args requests rerun help.
func IsRerunHelpRequested(args []string) bool {
	for _, arg := range args {
		if IsRerunHelpToken(arg) {
			return true
		}
	}

	return false
}

// IsRerunHelpToken reports whether tok is a help-seeking keyword or flag.
func IsRerunHelpToken(tok string) bool {
	clean := strings.ToLower(strings.TrimSpace(tok))

	return clean == "help" ||
		clean == "-h" ||
		clean == "--help" ||
		clean == "-help" ||
		clean == "--h" ||
		clean == "/?" ||
		clean == "-?" ||
		clean == "--?" ||
		clean == "?" ||
		clean == "/h" ||
		clean == "/help" ||
		clean == "man" ||
		clean == "info" ||
		clean == "usage" ||
		clean == "-usage" ||
		clean == "--usage" ||
		clean == "/usage"
}

func isRerunHelpRequested(args []string) bool {
	return IsRerunHelpRequested(args)
}

func isRerunHelpToken(tok string) bool {
	return IsRerunHelpToken(tok)
}

// RenderAgyRerunHelp prints the rich interactive Antigravity rerun guide.
func RenderAgyRerunHelp() {
	renderRerunHelpHeader()
	renderRerunHelpCommands()
	renderRerunHelpFlagsPrimary()
	renderRerunHelpFlagsSecondary()
	renderRerunHelpExamples()
}

func renderRerunHelpHeader() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║        ANTIGRAVITY RERUN & PROMPT REPLAY GUIDE           ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Replay active or recent prompts into Antigravity workspaces with optional\n")
	fmt.Printf("  IDE restart, clipboard sync, model selection, and queued prefix guards.\n\n")
}

func renderRerunHelpCommands() {
	fmt.Printf("  %sCOMMANDS:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • %sgitmap rerun%s (or %srr%s)             Replay prompt into current workspace\n", constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap rerun 1%s (or %srr 1%s)         Replay prompt into project sequence #1\n", constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap rerun all%s (or %srra%s)        Replay active prompts across all running projects\n", constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap rerun queue%s (or %srrq%s)      Re-inject queued prompts with completion prefix\n", constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap rerun help%s                    Show this comprehensive help guide\n\n", constants.ColorCyan, constants.ColorReset)
}

func renderRerunHelpFlagsPrimary() {
	fmt.Printf("  %sFLAGS & OPTIONS:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • %s-p, --prompt <tpl>%s        Prefix prompt template name or ID\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s-d, --dry-run%s             Preview constructed prompt without execution\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s-r, --restart%s             Restart Antigravity IDE before replaying prompt\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s--no-restart%s              Direct prompt injection without closing IDE\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s-n, --new-conversation%s    Create new conversation for prompt replay (default: true)\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s-m, --model <model>%s       Model for new conversation (flash_lite, flash, pro)\n", constants.ColorYellow, constants.ColorReset)
}

func renderRerunHelpFlagsSecondary() {
	fmt.Printf("    • %s-P, --project <id|seq>%s    Target project by sequence number, ID, or slug\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s-a, --all%s                 Rerun active prompts across all active projects\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s-q, --queue%s               Re-inject queued prompts with completion prefix\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s--prefix <prefix>%s         Custom prefix added to queued prompts\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("    • %s--no-clipboard%s            Do not copy constructed prompt to system clipboard\n\n", constants.ColorYellow, constants.ColorReset)
}

func renderRerunHelpExamples() {
	fmt.Printf("  %sEXAMPLES:%s\n", constants.ColorWhite, constants.ColorReset)
	fmt.Printf("    • %sgitmap rerun%s                      Replay last prompt in active project\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap rr 2 --dry-run%s             Preview prompt replay for project #2\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap rerun gitmap -r%s            Restart IDE and replay into 'gitmap'\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • %sgitmap rerun all --model pro%s      Replay all projects using Pro model\n\n", constants.ColorCyan, constants.ColorReset)
}
