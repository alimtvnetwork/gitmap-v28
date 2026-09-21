package cmdantigravity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

var (
	projPromptName   string
	projPromptTxt    string
	projPromptPrefix bool
	projPromptPf     bool
	projPromptSuffix bool
	projPromptSf     bool
)

// PromptProjectCmd targets a project by prefix and dispatches a templated/text prompt.
var PromptProjectCmd = &cobra.Command{
	Use:     "prompt-project <projectStartsWithName>",
	Aliases: []string{"p", "prompt-p"},
	Short:   "Send a prompt to an Antigravity project matching prefix",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPromptProject(args)
	},
}

func init() {
	initPromptProjectFlags(PromptProjectCmd)
}

func initPromptProjectFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&projPromptName, "name", "n", "", "Prompt template name")
	cmd.Flags().StringVarP(&projPromptTxt, "txt", "t", "", "Additional prompt text")
	cmd.Flags().BoolVar(&projPromptPrefix, "prefix", true, "Place template before text")
	cmd.Flags().BoolVar(&projPromptPf, "pf", true, "Place template before text (alias)")
	cmd.Flags().BoolVar(&projPromptSuffix, "suffix", false, "Place template after text")
	cmd.Flags().BoolVar(&projPromptSf, "sf", false, "Place template after text (alias)")
}

func runPromptProject(args []string) error {
	targetPath, findErr := FindProjectByPrefix(args[0])
	if findErr != nil {
		return findErr
	}
	templateContent, tplErr := resolveProjectTemplateContent(projPromptName)
	if tplErr != nil {
		return tplErr
	}
	textVal := resolveProjectText(args)
	if len(templateContent) == 0 && len(textVal) == 0 {
		return apperror.NewSimple("prompt text or template name is required", "E9010")
	}
	isSuffix := resolveIsSuffix(projPromptSuffix, projPromptSf, projPromptPrefix, projPromptPf)
	assembled := AssemblePrompt(templateContent, textVal, isSuffix)

	return EnqueueWithDualQueuePolicy(targetPath, projPromptName, assembled)
}

func resolveProjectTemplateContent(name string) (string, error) {
	hasName := len(strings.TrimSpace(name)) > 0
	if hasName == false {
		return "", nil
	}
	tpl, hasTpl := cmdprompttemplate.FindTemplate(name)
	if hasTpl {
		return tpl.Content, nil
	}

	return "", apperror.NewSimple("prompt template not found: "+name, "E9013")
}

func resolveProjectText(args []string) string {
	hasTxt := len(strings.TrimSpace(projPromptTxt)) > 0
	if hasTxt {
		return projPromptTxt
	}
	hasExtraArgs := len(args) > 1
	if hasExtraArgs {
		return strings.Join(args[1:], " ")
	}

	return ""
}

func resolveIsSuffix(suffix, sf, prefix, pf bool) bool {
	isSuffixExplicit := suffix || sf
	if isSuffixExplicit {
		return true
	}

	return false
}

// FindProjectByPrefix locates a project directory by name, ID, or path prefix.
func FindProjectByPrefix(prefix string) (string, error) {
	cleanPrefix := strings.ToLower(strings.TrimSpace(prefix))
	configDir := resolveAgyProjectsConfigDir()
	path, isFound := searchProjectsDirByPrefix(configDir, cleanPrefix)
	if isFound {
		return path, nil
	}

	return searchLocalDirectoryByPrefix(prefix)
}

func searchProjectsDirByPrefix(configDir, prefix string) (string, bool) {
	entries, err := os.ReadDir(configDir)
	hasErr := err != nil
	if hasErr {
		return "", false
	}
	for _, entry := range entries {
		isJson := filepath.Ext(entry.Name()) == ".json"
		if isJson == false {
			continue
		}

		path, isMatch := checkProjectFilePrefix(filepath.Join(configDir, entry.Name()), prefix)
		if isMatch {
			return path, true
		}
	}

	return "", false
}

func checkProjectFilePrefix(filePath, prefix string) (string, bool) {
	data, err := os.ReadFile(filePath)
	hasData := err == nil
	if hasData == false {
		return "", false
	}
	var cfg rawProjectConfig
	hasCfg := json.Unmarshal(data, &cfg) == nil
	if hasCfg && isProjectConfigMatch(cfg, prefix) {
		return resolveProjectConfigPath(cfg), true
	}

	return "", false
}

func isProjectConfigMatch(cfg rawProjectConfig, prefix string) bool {
	lowName := strings.ToLower(cfg.Name)
	lowID := strings.ToLower(cfg.ID)
	if strings.HasPrefix(lowName, prefix) || strings.HasPrefix(lowID, prefix) {
		return true
	}

	return false
}

func resolveProjectConfigPath(cfg rawProjectConfig) string {
	hasResources := cfg.ProjectResources != nil && len(cfg.ProjectResources.Resources) > 0
	if hasResources && cfg.ProjectResources.Resources[0].GitFolder != nil {
		return parseFolderURI(cfg.ProjectResources.Resources[0].GitFolder.FolderURI)
	}

	return ""
}

func searchLocalDirectoryByPrefix(prefix string) (string, error) {
	info, err := os.Stat(prefix)
	isDir := err == nil && info.IsDir()
	if isDir {
		return resolveExistingDirAbs(prefix)
	}

	root, rootErr := gitutil.RepoRoot(prefix)
	hasRoot := rootErr == nil && len(root) > 0
	if hasRoot {
		return root, nil
	}

	return "", apperror.NewSimple("no Antigravity project found matching prefix: "+prefix, "E9012")
}

func resolveExistingDirAbs(prefix string) (string, error) {
	abs, absErr := filepath.Abs(prefix)
	hasAbs := absErr == nil
	if hasAbs {
		return abs, nil
	}

	return prefix, nil
}
