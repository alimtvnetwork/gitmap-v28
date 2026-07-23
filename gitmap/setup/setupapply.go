package setup

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// gitSetting is a key-value pair for git config.
type gitSetting struct {
	Key   string
	Value string
}

// applySection applies a group of settings and prints results.
func applySection(section string, settings []gitSetting, dryRun bool, r *SetupResult) {
	fmt.Printf("\n  %s■ %s%s\n", constants.ColorYellow, section, constants.ColorReset)
	for _, s := range settings {
		applyOne(s, dryRun, r)
	}
}

// applyOne applies a single git config --global setting.
func applyOne(s gitSetting, dryRun bool, r *SetupResult) {
	if dryRun {
		printDryRunSetting(s, r)

		return
	}

	current := getCurrentValue(s.Key)
	if current == s.Value {
		printSkippedSetting(s, r)

		return
	}

	executeSetting(s, r)
}

// printDryRunSetting logs a dry-run setting.
func printDryRunSetting(s gitSetting, r *SetupResult) {
	fmt.Printf("  %s[dry-run]%s git config --global %s %q\n",
		constants.ColorDim, constants.ColorReset, s.Key, s.Value)
	r.Skipped++
}

// printSkippedSetting logs an already-set setting.
func printSkippedSetting(s gitSetting, r *SetupResult) {
	fmt.Printf("  %s⊘ %s%s = %s (already set)\n",
		constants.ColorDim, s.Key, constants.ColorReset, s.Value)
	r.Skipped++
}

// executeSetting runs git config --global for a single key-value.
func executeSetting(s gitSetting, r *SetupResult) {
	cmd := exec.Command(constants.GitBin, constants.GitConfigCmd, constants.SetupGlobalFlag, s.Key, s.Value)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("%s: %s — %s", s.Key, err, strings.TrimSpace(string(out)))
		fmt.Printf("  %s✗ %s%s\n", constants.ColorYellow, errMsg, constants.ColorReset)
		r.Failed++
		r.Errors = append(r.Errors, errMsg)

		return
	}

	fmt.Printf("  %s✓ %s%s = %s\n", constants.ColorGreen, s.Key, constants.ColorReset, s.Value)
	r.Applied++
}

// getCurrentValue reads the current git config value for a key.
func getCurrentValue(key string) string {
	cmd := exec.Command(constants.GitBin, constants.GitConfigCmd, constants.SetupGlobalFlag, constants.GitGetFlag, key)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// mapCoreKey maps JSON-friendly keys to git config keys.
func mapCoreKey(jsonKey string) string {
	coreMap := map[string]string{
		"autocrlf":      "core.autocrlf",
		"longpaths":     "core.longpaths",
		"editor":        "core.editor",
		"safecrlf":      "core.safecrlf",
		"defaultBranch": "init.defaultBranch",
	}
	if mapped, ok := coreMap[jsonKey]; ok {
		return mapped
	}

	return "core." + jsonKey
}
