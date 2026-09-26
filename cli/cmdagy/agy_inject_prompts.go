package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

var (
	injectPromptsPrefix  string
	injectPromptsSuffix  string
	injectPromptsRerun   int
	isInjectPromptsWatch bool
	isInjectPromptsUIUX  bool
	injectPromptsSSH     string
)

// AgyInjectPromptsCmd provides iterative prompt injection from folders or files.
var AgyInjectPromptsCmd = &cobra.Command{
	Use:     "inject-prompts <folder-or-text>",
	Aliases: []string{"ip", "inject-prompts-ssh", "ip-ssh"},
	Short:   "Inject prompts iteratively with templates, prefix, suffix, and watch mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: gitmap agy inject-prompts <folder-or-text> [--watch] [--rerun N] [--prefix ...] [--suffix ...]")
		}
		appErr := runAgyInjectPrompts(args[0])
		if appErr != nil {
			return appErr
		}
		return nil
	},
}

func init() {
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "prefix", "", "Prefix template from catalog or CSV")
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsPrefix, "pfx", "", "Prefix template (alias)")
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsSuffix, "suffix", "", "Suffix template from catalog or CSV")
	AgyInjectPromptsCmd.Flags().BoolVar(&isInjectPromptsUIUX, "ui-ux", false, "Use random UI/UX template as prefix")
	AgyInjectPromptsCmd.Flags().BoolVar(&isInjectPromptsUIUX, "uu", false, "Use random UI/UX template (alias)")
	AgyInjectPromptsCmd.Flags().IntVar(&injectPromptsRerun, "rerun", 0, "Rerun loop count")
	AgyInjectPromptsCmd.Flags().BoolVar(&isInjectPromptsWatch, "watch", false, "Watch until each prompt completes before continuing")
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsSSH, "nodes", "", "Specific SSH nodes to target")
	AgyInjectPromptsCmd.Flags().StringVar(&injectPromptsSSH, "ssh-only-node", "", "SSH node targeting")
	AgyCmd.AddCommand(AgyInjectPromptsCmd)
}

func runAgyInjectPrompts(target string) *apperror.AppError {
	cwd, _ := os.Getwd()
	prefixText, suffixText := resolveAppliedTemplates(injectPromptsPrefix, injectPromptsSuffix, isInjectPromptsUIUX)

	fi, err := os.Stat(target)
	if err == nil && fi.IsDir() {
		return injectFromFolder(target, cwd, prefixText, suffixText)
	}

	return injectSingleTarget(target, cwd, prefixText, suffixText)
}

func injectFromFolder(dir, projectPath, prefix, suffix string) *apperror.AppError {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return apperror.WrapSimple(err, "read inject folder")
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".md" || ext == ".txt" {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(files)

	if len(files) == 0 {
		return apperror.NewSimple("no .md or .txt prompt files found in "+dir, "E_NO_PROMPTS")
	}

	fmt.Printf("\n  %s Discovered %d prompt files in %s\n\n", constants.ColorCyan+"▶"+constants.ColorReset, len(files), dir)

	for i, f := range files {
		data, readErr := os.ReadFile(f)
		if readErr != nil {
			continue
		}
		title := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		content := decoratePromptBody(string(data), prefix, suffix)

		fmt.Printf("  [%d/%d] Injecting: %s%s%s\n", i+1, len(files), constants.ColorYellow, title, constants.ColorReset)
		executePromptWithRerun(projectPath, content, title, injectPromptsRerun, isInjectPromptsWatch)
	}

	return nil
}

func injectSingleTarget(target, projectPath, prefix, suffix string) *apperror.AppError {
	var body, title string
	if fi, err := os.Stat(target); err == nil && !fi.IsDir() {
		data, _ := os.ReadFile(target)
		body = string(data)
		title = strings.TrimSuffix(filepath.Base(target), filepath.Ext(target))
	} else {
		body = target
		title = "Injected Custom Prompt"
	}

	content := decoratePromptBody(body, prefix, suffix)
	executePromptWithRerun(projectPath, content, title, injectPromptsRerun, isInjectPromptsWatch)
	return nil
}

func executePromptWithRerun(projectPath, content, title string, rerunCount int, isWatch bool) {
	loops := 1
	if rerunCount > 0 {
		loops += rerunCount
	}

	for l := 1; l <= loops; l++ {
		loopTitle := title
		if loops > 1 {
			loopTitle = fmt.Sprintf("%s (Run %d/%d)", title, l, loops)
		}

		res := InjectAgyPrompt(projectPath, content, loopTitle, "user_inject", false)
		if res.IsSuccess {
			fmt.Printf("    %s Prompt injected successfully.\n", constants.ColorGreen+"✔"+constants.ColorReset)
		} else {
			fmt.Printf("    %s Injection warning: %s\n", constants.ColorYellow+"⚠"+constants.ColorReset, res.Message)
		}

		if isWatch {
			waitForPromptCompletion(projectPath)
		}
	}
}

func waitForPromptCompletion(projectPath string) {
	fmt.Printf("    %s Waiting for prompt execution to finish...", constants.ColorDim)
	for {
		time.Sleep(3 * time.Second)
		activeConvs, _ := FetchActiveRunningConversations()
		hasRunning := false
		for _, c := range activeConvs {
			if strings.EqualFold(c.WorkspacePath, projectPath) || strings.Contains(projectPath, c.ProjectID) {
				hasRunning = true
				break
			}
		}
		if !hasRunning {
			break
		}
	}
	fmt.Printf(" %sDone.%s\n", constants.ColorGreen, constants.ColorReset)
}

func resolveAppliedTemplates(prefixSpec, suffixSpec string, isUIUX bool) (string, string) {
	var prefixText, suffixText string

	if isUIUX {
		prefixText = cmdprompttemplate.ResolveRandomTemplateContent("ui-ux")
	} else if prefixSpec != "" {
		prefixText = cmdprompttemplate.ResolveRandomTemplateContent(prefixSpec)
	}

	if suffixSpec != "" {
		suffixText = cmdprompttemplate.ResolveRandomTemplateContent(suffixSpec)
	}

	return prefixText, suffixText
}

func decoratePromptBody(body, prefix, suffix string) string {
	var parts []string
	if strings.TrimSpace(prefix) != "" {
		parts = append(parts, strings.TrimSpace(prefix))
	}
	parts = append(parts, strings.TrimSpace(body))
	if strings.TrimSpace(suffix) != "" {
		parts = append(parts, strings.TrimSpace(suffix))
	}

	return strings.Join(parts, "\n\n")
}
