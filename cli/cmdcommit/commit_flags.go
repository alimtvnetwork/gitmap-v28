package cmdcommit

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ExtractCommitFlags parses CLI arguments, strips metadata flags, and resolves ambient fallbacks.
func ExtractCommitFlags(args []string) ([]string, CommitMetadata) {
	var meta CommitMetadata
	clean := parseArgsFlags(args, &meta)
	applyAmbientFallbacks(&meta)

	return clean, meta
}

func parseArgsFlags(args []string, meta *CommitMetadata) []string {
	var clean []string
	isSkip := false
	for i := 0; i < len(args); i++ {
		if isSkip {
			isSkip = false
			continue
		}
		clean, isSkip = processSingleFlag(args, i, clean, meta)
	}

	return clean
}

func processSingleFlag(args []string, i int, clean []string, meta *CommitMetadata) ([]string, bool) {
	arg := args[i]
	if hasValuePrefix(arg, meta) {
		return clean, false
	}
	if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && isValuedFlag(arg, args[i+1], meta) {
		return clean, true
	}
	if isStandaloneFlag(arg, meta) {
		return clean, false
	}

	return append(clean, arg), false
}

func isStandaloneFlag(arg string, meta *CommitMetadata) bool {
	switch arg {
	case FlagPushLong, FlagPushShort:
		meta.HasPush = true
		return true
	case FlagDryRunLong, FlagDryRunShort:
		meta.IsDryRun = true
		return true
	case FlagAI:
		meta.IsAIMode = true
		return true
	case FlagSchedule, FlagScheduleShort:
		meta.IsTimeEnforced = true
		return true
	case FlagAutoWindow:
		meta.IsAutoWindow = true
		return true
	default:
		return false
	}
}

func hasValuePrefix(arg string, meta *CommitMetadata) bool {
	if hasProvenancePrefix(arg, meta) {
		return true
	}
	if hasAuthorOrTaskPrefix(arg, meta) {
		return true
	}

	return hasTimePrefix(arg, meta)
}

func hasProvenancePrefix(arg string, meta *CommitMetadata) bool {
	if hasAgentOrHeatmapPrefix(arg, meta) {
		return true
	}

	return hasAIOrSpecPrefix(arg, meta)
}

func hasAgentOrHeatmapPrefix(arg string, meta *CommitMetadata) bool {
	if matchPrefixTarget(arg, []string{FlagAgentID + "=", FlagAgent + "="}, &meta.AgentID) {
		meta.IsAIMode = true
		return true
	}

	return matchPrefixTarget(arg, []string{FlagHeatmapID + "=", FlagHeatmap + "="}, &meta.HeatmapID)
}

func hasAIOrSpecPrefix(arg string, meta *CommitMetadata) bool {
	if matchPrefixTarget(arg, []string{FlagAI + "="}, &meta.AIName) {
		meta.IsAIMode = true
		return true
	}

	return matchPrefixTarget(arg, []string{FlagSpec + "="}, &meta.SpecPath)
}

func hasAuthorOrTaskPrefix(arg string, meta *CommitMetadata) bool {
	if matchPrefixTarget(arg, []string{FlagAuthor + "=", FlagAuthorShort + "="}, &meta.Author) {
		return true
	}

	return matchPrefixTarget(arg, []string{FlagTask + "=", "-t="}, &meta.TaskSlug)
}

func hasTimePrefix(arg string, meta *CommitMetadata) bool {
	if matchPrefixTarget(arg, []string{FlagTimezone + "=", FlagTimezoneShort + "="}, &meta.Timezone) {
		return true
	}

	return hasScheduleOrAutoPrefix(arg, meta)
}

func hasScheduleOrAutoPrefix(arg string, meta *CommitMetadata) bool {
	if strings.HasPrefix(arg, FlagSchedule+"=") || strings.HasPrefix(arg, FlagScheduleShort+"=") {
		meta.IsTimeEnforced = true
		return true
	}
	if strings.HasPrefix(arg, FlagAutoWindow+"=") {
		meta.IsAutoWindow = true
		return true
	}

	return false
}

func matchPrefixTarget(arg string, prefixes []string, target *string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(arg, p) {
			*target = strings.TrimSpace(strings.TrimPrefix(arg, p))
			return true
		}
	}

	return false
}

func isValuedFlag(arg, val string, meta *CommitMetadata) bool {
	if isAgentFlag(arg, val, meta) || isHeatmapFlag(arg, val, meta) {
		return true
	}
	if isAIFlag(arg, val, meta) || isSpecFlag(arg, val, meta) {
		return true
	}
	if isAuthorFlag(arg, val, meta) || isTaskFlag(arg, val, meta) {
		return true
	}

	return isTimezoneFlag(arg, val, meta)
}

func isTimezoneFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagTimezone && arg != FlagTimezoneShort {
		return false
	}
	meta.Timezone = strings.TrimSpace(val)

	return true
}

func isAgentFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagAgentID && arg != FlagAgent {
		return false
	}
	meta.AgentID = strings.TrimSpace(val)
	meta.IsAIMode = true

	return true
}

func isHeatmapFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagHeatmapID && arg != FlagHeatmap {
		return false
	}
	meta.HeatmapID = strings.TrimSpace(val)

	return true
}

func isAIFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagAI || strings.HasPrefix(val, "-") {
		return false
	}
	meta.AIName = strings.TrimSpace(val)
	meta.IsAIMode = true

	return true
}

func isSpecFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagSpec {
		return false
	}
	meta.SpecPath = strings.TrimSpace(val)

	return true
}

func isTaskFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagTask && arg != FlagTaskShort {
		return false
	}
	meta.TaskSlug = strings.TrimSpace(val)

	return true
}

func isAuthorFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagAuthor && arg != FlagAuthorShort {
		return false
	}
	meta.Author = strings.TrimSpace(val)

	return true
}

func applyAmbientFallbacks(meta *CommitMetadata) {
	fallbackString(&meta.AgentID, EnvGitmapAgentID, EnvAntigravityAgentID, EnvAgentID)
	fallbackString(&meta.HeatmapID, EnvGitmapHeatmapID, EnvGitmapHeatmapSessionID)
	fallbackString(&meta.TaskSlug, EnvGitmapTaskSlug, EnvGitmapTask)
	fallbackString(&meta.SpecPath, EnvGitmapSpecPath)
	resolveAuthorFallback(meta)
	resolveAIFallbacks(meta)
	resolveTimeFallbacks(meta)
}

func resolveTimeFallbacks(meta *CommitMetadata) {
	fallbackString(&meta.Timezone, EnvGitmapTimezone, "TZ")
	if !meta.IsTimeEnforced && os.Getenv(EnvGitmapTimeEnforced) == "true" {
		meta.IsTimeEnforced = true
	}
}

func resolveAuthorFallback(meta *CommitMetadata) {
	if meta.Author != "" {
		return
	}
	name, _ := execGitOutputCP("config", "user.name")
	email, _ := execGitOutputCP("config", "user.email")
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" && email == "" {
		return
	}
	meta.Author = formatAuthorIdentity(name, email)
}

func formatAuthorIdentity(name, email string) string {
	if name != "" && email != "" {
		return name + " <" + email + ">"
	}
	if name != "" {
		return name
	}

	return "<" + email + ">"
}

func resolveAIFallbacks(meta *CommitMetadata) {
	fallbackString(&meta.AIName, EnvGitmapAIName, EnvAIModel, EnvGeminiModel, EnvMuseModel)
	if hasAmbientAIEnv() || meta.AgentID != "" || meta.AIName != "" {
		meta.IsAIMode = true
	}
	if meta.IsAIMode && meta.AIName == "" {
		meta.AIName = store.DetectCurrentAIName()
	}
}

func hasAmbientAIEnv() bool {
	aiEnvs := []string{
		EnvGitmapAgentID, EnvAntigravityAgentID, EnvAgentID,
		EnvGitmapAIName, EnvAIModel, EnvGeminiModel, EnvMuseModel,
		"ANTIGRAVITY_INSTANCE_ID", "GEMINI_CLI_VERSION", "MUSE_RUNNER_ID",
	}

	return hasAnyEnvValue(aiEnvs)
}

func hasAnyEnvValue(keys []string) bool {
	for _, k := range keys {
		if strings.TrimSpace(os.Getenv(k)) != "" {
			return true
		}
	}

	return false
}

func fallbackString(target *string, envVars ...string) {
	if *target != "" {
		return
	}
	for _, env := range envVars {
		if val := strings.TrimSpace(os.Getenv(env)); val != "" {
			*target = val
			return
		}
	}
}
