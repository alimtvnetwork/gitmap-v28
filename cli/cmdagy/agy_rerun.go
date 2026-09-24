package cmdagy

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	rerunPromptTpl   string
	rerunNoClipboard bool
	rerunDryRun      bool
	rerunRestartFlag bool
	rerunNoRestart   bool
	rerunProjectFlag string
	rerunConvFlag    string
)

var agyRerunCmd = &cobra.Command{
	Use:   "rerun [1|2|3|4|project] [flags]",
	Short: "Rerun last prompt for project with Antigravity IDE restart and media attachments",
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgyRerun(args)
		if appErr != nil {
			return appErr
		}

		return nil
	},
}

func init() {
	agyRerunCmd.Flags().StringVarP(&rerunPromptTpl, "prompt", "p", cmdprompttemplate.DefaultTemplateID, "Prefix prompt template name or ID")
	agyRerunCmd.Flags().BoolVar(&rerunNoClipboard, "no-clipboard", false, "Do not copy constructed prompt to clipboard")
	agyRerunCmd.Flags().BoolVarP(&rerunDryRun, "dry-run", "d", false, "Preview constructed prompt without execution")
	agyRerunCmd.Flags().BoolVarP(&rerunRestartFlag, "restart", "r", true, "Restart Antigravity IDE and replay prompt")
	agyRerunCmd.Flags().BoolVar(&rerunNoRestart, "no-restart", false, "Do not restart IDE, inject directly")
	agyRerunCmd.Flags().StringVarP(&rerunProjectFlag, "project", "P", "", "Target project by index (1, 2, 3...) or name")
	agyRerunCmd.Flags().StringVarP(&rerunConvFlag, "conversation", "c", "", "Target conversation ID")
}

// RunRerunTopLevelCLI executes agy rerun from top-level gitmap aliases.
func RunRerunTopLevelCLI(args []string) error {
	runArgs := append([]string{"rerun"}, args...)
	AgyCmd.SetArgs(runArgs)

	return AgyCmd.Execute()
}

func runAgyRerun(args []string) *apperror.AppError {
	if isLegacyLastInvocation(args) {
		return runLegacyAgyRerun(args)
	}

	target := resolveRerunProjectTarget(args)
	isRestart := rerunRestartFlag && !rerunNoRestart
	tplName := extractRerunTemplateName(args)

	err := RestartAndRerunProject(target, isRestart, rerunDryRun, tplName)
	if err != nil {
		return apperror.WrapSimple(err, "agy rerun")
	}

	return nil
}

func isLegacyLastInvocation(args []string) bool {
	if len(args) == 0 {
		return false
	}

	return strings.EqualFold(args[0], "last")
}

func resolveRerunProjectTarget(args []string) string {
	if rerunProjectFlag != "" {
		return rerunProjectFlag
	}

	for _, a := range args {
		trimmed := strings.TrimSpace(a)
		if trimmed != "" && !strings.HasPrefix(trimmed, "-") {
			return trimmed
		}
	}

	return ""
}

func runLegacyAgyRerun(args []string) *apperror.AppError {
	count := parseRerunCount(args)
	tplName := extractRerunTemplateName(args)
	tplContent := resolveRerunTemplate(tplName)
	cwd, _ := os.Getwd()
	prompts := resolvePromptsForRerun(cwd, count)
	if len(prompts) == 0 {
		printNoPromptsWarning()

		return nil
	}

	executeRerunPayload(tplContent, prompts)

	return nil
}

func printNoPromptsWarning() {
	fmt.Printf("  %s⚠ No recent prompts found to rerun.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func executeRerunPayload(tplContent string, prompts []AgyPromptEntry) {
	payload := buildRerunPayload(tplContent, prompts)
	renderRerunOutput(payload, len(prompts))
	handleRerunClipboard(payload)
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
	for i := 0; i < len(args); i++ {
		tok := args[i]
		if isPromptFlag(tok) {
			i++
			continue
		}

		if isSkippableRerunToken(tok) {
			continue
		}

		out = append(out, tok)
	}

	return out
}

func isPromptFlag(tok string) bool {
	low := strings.ToLower(tok)

	return low == "-p" || low == "--prompt"
}

func isSkippableRerunToken(tok string) bool {
	low := strings.ToLower(tok)
	if low == "last" || low == "rerun" || low == "rr" {
		return true
	}

	return strings.HasPrefix(tok, "-")
}

func extractRerunTemplateName(args []string) string {
	for i := 0; i < len(args); i++ {
		tpl := matchPromptFlagAt(args, i)
		if tpl != "" {
			return tpl
		}
	}

	return rerunPromptTpl
}

func matchPromptFlagAt(args []string, i int) string {
	if isPromptFlag(args[i]) && i+1 < len(args) {
		return args[i+1]
	}

	if strings.HasPrefix(args[i], "-p=") {
		return strings.TrimPrefix(args[i], "-p=")
	}

	if strings.HasPrefix(args[i], "--prompt=") {
		return strings.TrimPrefix(args[i], "--prompt=")
	}

	return ""
}

func resolveRerunTemplate(tplID string) string {
	if tplID == "" {
		tplID = cmdprompttemplate.DefaultTemplateID
	}

	tpl, isFound := cmdprompttemplate.FindTemplate(tplID)
	if isFound {
		return tpl.Content
	}

	return cmdprompttemplate.DefaultIsDoneContent
}
