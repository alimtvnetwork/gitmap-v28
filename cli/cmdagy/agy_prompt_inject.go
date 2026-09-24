package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	injectProjectFlag string
	injectConvFlag    string
	injectTitleFlag   string
	injectSkipFlag    bool
)

func resolveInjectProjectTarget(args []string) (string, error) {
	target := injectProjectFlag
	if target == "" && len(args) >= 2 {
		target = args[1]
	}
	if target == "" {
		target = "."
	}

	absPath, err := filepath.Abs(target)
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve project path")
	}

	return absPath, nil
}

func loadInjectPromptPayload(slugOrPath string) (string, string, error) {
	if _, err := os.Stat(slugOrPath); err != nil {
		return findPromptInSearchPaths(slugOrPath)
	}

	content, readErr := os.ReadFile(slugOrPath)
	if readErr != nil {
		return "", "", apperror.WrapSimple(readErr, "read prompt file")
	}
	title := strings.TrimSuffix(filepath.Base(slugOrPath), filepath.Ext(slugOrPath))

	return string(content), title, nil
}

func findPromptInSearchPaths(slug string) (string, string, error) {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(".ai-memory", "prompts", slug+".md"),
		filepath.Join("01-prompts", slug+".md"),
		filepath.Join(home, ".gitmap", "prompts", slug+".md"),
	}

	for _, cand := range candidates {
		if data, err := os.ReadFile(cand); err == nil {
			return string(data), slug, nil
		}
	}

	return slug, "Custom Prompt", nil
}

// RunAgyPromptInject executes prompt injection for Antigravity projects and conversations.
func RunAgyPromptInject(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: gitmap agy prompt inject <slug-or-file> [target-project] [--conversation <id>]")
	}

	projectPath, err := resolveInjectProjectTarget(args)
	if err != nil {
		return err
	}

	content, title, loadErr := loadInjectPromptPayload(args[0])
	if loadErr != nil {
		return loadErr
	}

	if injectTitleFlag != "" {
		title = injectTitleFlag
	}

	return executeInjectPayload(projectPath, content, title)
}

func executeInjectPayload(projectPath, content, title string) error {
	convID := resolveInjectConvID(projectPath)
	fmt.Printf("\n%s Injecting prompt \033[1m%q\033[0m into target:\n", constants.ColorCyan+"▶"+constants.ColorReset, title)
	fmt.Printf("    • Project:      %s\n", projectPath)
	fmt.Printf("    • Conversation: %s\n", convID)

	res := InjectAgyPrompt(projectPath, content, title, "user_inject", injectSkipFlag)
	if res.IsSuccess {
		fmt.Printf("  %s %s\n\n", constants.ColorGreen+"✓"+constants.ColorReset, res.Message)

		return nil
	}

	fmt.Printf("  %s %s\n\n", constants.ColorYellow+"!"+constants.ColorReset, res.Message)

	return nil
}

func resolveInjectConvID(projectPath string) string {
	if injectConvFlag != "" {
		return injectConvFlag
	}

	conv, err := SelectMatchingConversation(projectPath)
	if err == nil && conv.ID != "" {
		return conv.ID
	}

	return "default (active)"
}
