package cmdagy

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
)

var (
	rerunPromptTpl    string
	rerunNoClipboard  bool
	rerunDryRun       bool
)

var agyRerunCmd = &cobra.Command{
	Use:   "rerun [last] [N]",
	Short: "Replay recent prompts with optional prefix template",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyRerun(args)
	},
}

func init() {
	agyRerunCmd.Flags().StringVarP(&rerunPromptTpl, "prompt", "p", cmdprompttemplate.DefaultTemplateID, "Prefix prompt template name or ID")
	agyRerunCmd.Flags().BoolVar(&rerunNoClipboard, "no-clipboard", false, "Do not copy constructed prompt to clipboard")
	agyRerunCmd.Flags().BoolVarP(&rerunDryRun, "dry-run", "d", false, "Preview constructed prompt without execution")
}

func runAgyRerun(args []string) error {
	count := parseRerunCount(args)
	tplContent := resolveRerunTemplate(rerunPromptTpl)
	cwd, _ := os.Getwd()
	prompts := resolvePromptsForRerun(cwd, count)
	if len(prompts) == 0 {
		fmt.Printf("  %s⚠ No recent prompts found to rerun.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}
	payload := buildRerunPayload(tplContent, prompts)
	renderRerunOutput(payload, len(prompts))
	handleRerunClipboard(payload)

	return nil
}

func parseRerunCount(args []string) int {
	cleanArgs := stripLeadingRerunTokens(args)
	if len(cleanArgs) == 0 {
		return 1
	}
	val, err := strconv.Atoi(cleanArgs[0])
	if err != nil || val < 1 {
		return 1
	}

	return val
}

func stripLeadingRerunTokens(args []string) []string {
	var out []string
	for _, a := range args {
		low := strings.ToLower(a)
		if low != "last" && low != "rerun" {
			out = append(out, a)
		}
	}

	return out
}

func resolveRerunTemplate(tplID string) string {
	if tplID == "" {
		tplID = cmdprompttemplate.DefaultTemplateID
	}
	tpl, found := cmdprompttemplate.FindTemplate(tplID)
	if found {
		return tpl.Content
	}

	return cmdprompttemplate.DefaultIsDoneContent
}
