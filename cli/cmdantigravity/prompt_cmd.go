package cmdantigravity

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

var (
	basePromptName   string
	basePromptTxt    string
	basePromptPrefix bool
	basePromptPf     bool
	basePromptSuffix bool
	basePromptSf     bool

	withNameTxt    string
	withNamePrefix bool
	withNamePf     bool
	withNameSuffix bool
	withNameSf     bool

	txtPromptPrefix bool
	txtPromptPf     bool
	txtPromptSuffix bool
	txtPromptSf     bool
	txtPromptTxt    string
)

// PromptCmd dispatches prompts to current repository in Antigravity.
var PromptCmd = &cobra.Command{
	Use:     "prompt",
	Aliases: []string{"pr"},
	Short:   "Send or enqueue a prompt for the current project in Antigravity",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBasePrompt(args)
	},
}

// PromptWithNameCmd dispatches a named prompt template with optional text.
var PromptWithNameCmd = &cobra.Command{
	Use:     "prompt-with-name <templateName>",
	Aliases: []string{"pwn"},
	Short:   "Send a named prompt template to current project in Antigravity",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPromptWithName(args)
	},
}

// PromptTxtCmd dispatches direct text to current project in Antigravity.
var PromptTxtCmd = &cobra.Command{
	Use:     "prompt-txt [text]",
	Aliases: []string{"pt"},
	Short:   "Send direct prompt text to current project in Antigravity",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPromptTxt(args)
	},
}

// PromptLsCmd lists all available prompt templates in a table.
var PromptLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List available prompt templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdprompt.RunPromptList(args)
	},
}

func init() {
	initBasePromptFlags(PromptCmd)
	initWithNameFlags(PromptWithNameCmd)
	initTxtFlags(PromptTxtCmd)
	PromptCmd.AddCommand(PromptLsCmd)
}

func initBasePromptFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&basePromptName, "name", "n", "", "Prompt template name")
	cmd.Flags().StringVarP(&basePromptTxt, "txt", "t", "", "Additional prompt text")
	cmd.Flags().BoolVar(&basePromptPrefix, "prefix", true, "Place template before text")
	cmd.Flags().BoolVar(&basePromptPf, "pf", true, "Place template before text (alias)")
	cmd.Flags().BoolVar(&basePromptSuffix, "suffix", false, "Place template after text")
	cmd.Flags().BoolVar(&basePromptSf, "sf", false, "Place template after text (alias)")
}

func initWithNameFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&withNameTxt, "txt", "t", "", "Additional prompt text")
	cmd.Flags().BoolVar(&withNamePrefix, "prefix", true, "Place template before text")
	cmd.Flags().BoolVar(&withNamePf, "pf", true, "Place template before text (alias)")
	cmd.Flags().BoolVar(&withNameSuffix, "suffix", false, "Place template after text")
	cmd.Flags().BoolVar(&withNameSf, "sf", false, "Place template after text (alias)")
}

func initTxtFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&txtPromptTxt, "txt", "t", "", "Prompt text")
	cmd.Flags().BoolVar(&txtPromptPrefix, "prefix", true, "Place template before text")
	cmd.Flags().BoolVar(&txtPromptPf, "pf", true, "Place template before text (alias)")
	cmd.Flags().BoolVar(&txtPromptSuffix, "suffix", false, "Place template after text")
	cmd.Flags().BoolVar(&txtPromptSf, "sf", false, "Place template after text (alias)")
}

func runBasePrompt(args []string) error {
	if IsHelpArg(args) {
		return PrintAgyPromptHelp()
	}

	repoRoot := resolveCurrentRepoRoot()
	templateName, templateContent, err := resolveBaseTemplate(args)
	if err != nil {
		return err
	}
	textVal := resolveBaseText(args, templateName)
	if len(templateContent) == 0 && len(textVal) == 0 {
		return apperror.NewSimple("prompt text or template name is required", "E9010")
	}
	isSuffix := resolveIsSuffix(basePromptSuffix, basePromptSf, basePromptPrefix, basePromptPf)
	assembled := AssemblePrompt(templateContent, textVal, isSuffix)

	return EnqueueWithDualQueuePolicy(repoRoot, templateName, assembled)
}

func resolveCurrentRepoRoot() string {
	root, err := gitutil.RepoRoot(".")
	hasRoot := err == nil && len(root) > 0
	if hasRoot {
		return root
	}
	cwd, _ := os.Getwd()

	return cwd
}

func resolveBaseTemplate(args []string) (string, string, error) {
	hasName := len(strings.TrimSpace(basePromptName)) > 0
	if hasName {
		content, err := resolveProjectTemplateContent(basePromptName)

		return basePromptName, content, err
	}
	hasArgs := len(args) > 0
	if hasArgs == false {
		return "", "", nil
	}

	tpl, hasTpl := cmdprompttemplate.FindTemplate(args[0])
	if hasTpl {
		return tpl.Name, tpl.Content, nil
	}

	return "", "", nil
}

func resolveBaseText(args []string, detectedTemplateName string) string {
	hasTxt := len(strings.TrimSpace(basePromptTxt)) > 0
	if hasTxt {
		return basePromptTxt
	}
	hasDetected := len(detectedTemplateName) > 0
	if hasDetected && len(args) > 1 {
		return strings.Join(args[1:], " ")
	}
	if hasDetected == false && len(args) > 0 {
		return strings.Join(args, " ")
	}

	return ""
}

func runPromptWithName(args []string) error {
	if IsHelpArg(args) {
		return PrintAgyPromptHelp()
	}

	repoRoot := resolveCurrentRepoRoot()
	tplName := args[0]
	tpl, hasTpl := cmdprompttemplate.FindTemplate(tplName)
	if hasTpl == false {
		return apperror.NewSimple("prompt template not found: "+tplName, "E9013")
	}
	textVal := resolveWithNameText(args)
	isSuffix := resolveIsSuffix(withNameSuffix, withNameSf, withNamePrefix, withNamePf)
	assembled := AssemblePrompt(tpl.Content, textVal, isSuffix)

	return EnqueueWithDualQueuePolicy(repoRoot, tpl.Name, assembled)
}

func resolveWithNameText(args []string) string {
	hasTxt := len(strings.TrimSpace(withNameTxt)) > 0
	if hasTxt {
		return withNameTxt
	}
	hasExtra := len(args) > 1
	if hasExtra {
		return strings.Join(args[1:], " ")
	}

	return ""
}

func runPromptTxt(args []string) error {
	if IsHelpArg(args) {
		return PrintAgyPromptHelp()
	}

	repoRoot := resolveCurrentRepoRoot()
	textVal := resolveTxtPromptText(args)
	if len(textVal) == 0 {
		return apperror.NewSimple("prompt text cannot be empty: gitmap agy prompt-txt <text>", "E9011")
	}
	isSuffix := resolveIsSuffix(txtPromptSuffix, txtPromptSf, txtPromptPrefix, txtPromptPf)
	assembled := formatTxtPrompt(textVal, isSuffix)

	return EnqueueWithDualQueuePolicy(repoRoot, "", assembled)
}

func formatTxtPrompt(textVal string, isSuffix bool) string {
	clean := strings.TrimSpace(textVal)
	if isSuffix {
		return clean + "\n\n"
	}

	return clean
}

func resolveTxtPromptText(args []string) string {
	hasTxt := len(strings.TrimSpace(txtPromptTxt)) > 0
	if hasTxt {
		return txtPromptTxt
	}
	hasArgs := len(args) > 0
	if hasArgs {
		return strings.Join(args, " ")
	}

	return ""
}
