// Package setup configures Git global settings from a JSON config file.
package setup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	gitConfigDiffTool            = "diff.tool"
	gitConfigDiffToolCmd         = "difftool.%s.cmd"
	gitConfigDiffToolPrompt      = "difftool.prompt"
	gitConfigDiffToolTrust       = "difftool.%s.trustExitCode"
	gitConfigMergeTool           = "merge.tool"
	gitConfigMergeToolCmd        = "mergetool.%s.cmd"
	gitConfigMergeToolPrompt     = "mergetool.prompt"
	gitConfigMergeToolKeepBackup = "mergetool.keepBackup"
	gitConfigMergeToolTrust      = "mergetool.%s.trustExitCode"
	gitConfigAlias               = "alias.%s"
	gitConfigCredentialHelper    = "credential.helper"
	gitValueTrue                 = "true"
	gitValueFalse                = "false"
)

// GitSetupConfig holds the full git-setup.json structure.
type GitSetupConfig struct {
	DiffTool         *ToolConfig       `json:"diffTool"`
	MergeTool        *ToolConfig       `json:"mergeTool"`
	Aliases          map[string]string `json:"aliases"`
	CredentialHelper string            `json:"credentialHelper"`
	Core             map[string]string `json:"core"`
}

// ToolConfig holds diff/merge tool configuration.
type ToolConfig struct {
	Name            string `json:"name"`
	Cmd             string `json:"cmd"`
	IsTrustExitCode bool   `json:"trustExitCode"`
}

// SetupResult tracks applied and failed settings.
type SetupResult struct {
	Applied int
	Skipped int
	Failed  int
	Errors  []string
}

// LoadConfig reads and parses the git-setup.json file.
func LoadConfig(path string) (GitSetupConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return GitSetupConfig{}, err
	}
	var cfg GitSetupConfig
	err = json.Unmarshal(data, &cfg)

	return cfg, err
}

// Apply applies the full git setup configuration.
func Apply(cfg GitSetupConfig, isDryRun bool) SetupResult {
	result := SetupResult{}

	if cfg.DiffTool != nil {
		applyDiffTool(cfg.DiffTool, isDryRun, &result)
	}
	if cfg.MergeTool != nil {
		applyMergeTool(cfg.MergeTool, isDryRun, &result)
	}
	if len(cfg.Aliases) > 0 {
		applyAliases(cfg.Aliases, isDryRun, &result)
	}
	if len(cfg.CredentialHelper) > 0 {
		applyCredentialHelper(cfg.CredentialHelper, isDryRun, &result)
	}
	if len(cfg.Core) > 0 {
		applyCoreSettings(cfg.Core, isDryRun, &result)
	}

	return result
}

// applyDiffTool configures git's global diff tool.
func applyDiffTool(tool *ToolConfig, isDryRun bool, r *SetupResult) {
	settings := []gitSetting{
		{gitConfigDiffTool, tool.Name},
		{fmt.Sprintf(gitConfigDiffToolCmd, tool.Name), tool.Cmd},
		{gitConfigDiffToolPrompt, gitValueFalse},
	}
	if tool.IsTrustExitCode {
		settings = append(settings, gitSetting{
			fmt.Sprintf(gitConfigDiffToolTrust, tool.Name), gitValueTrue,
		})
	}
	applySection(constants.SetupSectionDiff, settings, isDryRun, r)
}

// applyMergeTool configures git's global merge tool.
func applyMergeTool(tool *ToolConfig, isDryRun bool, r *SetupResult) {
	settings := []gitSetting{
		{gitConfigMergeTool, tool.Name},
		{fmt.Sprintf(gitConfigMergeToolCmd, tool.Name), tool.Cmd},
		{gitConfigMergeToolPrompt, gitValueFalse},
		{gitConfigMergeToolKeepBackup, gitValueFalse},
	}
	if tool.IsTrustExitCode {
		settings = append(settings, gitSetting{
			fmt.Sprintf(gitConfigMergeToolTrust, tool.Name), gitValueTrue,
		})
	}
	applySection(constants.SetupSectionMerge, settings, isDryRun, r)
}

// applyAliases configures git global aliases.
func applyAliases(aliases map[string]string, isDryRun bool, r *SetupResult) {
	settings := make([]gitSetting, 0, len(aliases))
	for name, value := range aliases {
		settings = append(settings, gitSetting{
			fmt.Sprintf(gitConfigAlias, name), value,
		})
	}
	applySection(constants.SetupSectionAlias, settings, isDryRun, r)
}

// applyCredentialHelper configures git's credential helper.
func applyCredentialHelper(helper string, isDryRun bool, r *SetupResult) {
	settings := []gitSetting{
		{gitConfigCredentialHelper, helper},
	}
	applySection(constants.SetupSectionCred, settings, isDryRun, r)
}

// applyCoreSettings configures git core settings.
func applyCoreSettings(core map[string]string, isDryRun bool, r *SetupResult) {
	settings := make([]gitSetting, 0, len(core))
	for key, value := range core {
		gitKey := mapCoreKey(key)
		settings = append(settings, gitSetting{gitKey, value})
	}
	applySection(constants.SetupSectionCore, settings, isDryRun, r)
}
