// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolveFinalDest(dest, destHint string) string {
	if dest != "" {
		return dest
	}
	if destHint != "" {
		return destHint
	}
	return "icon-studio-v1"
}

func parseNamedFlagSources(state parsedFlagState) (string, []ResolvedSourceItem, []SkippedSourceItem, bool, error) {
	destHint, resolved, skipped, err := parseSourceStringSlice(state.sources)
	dest := resolveFinalDest(state.dest, destHint)
	return dest, resolved, skipped, state.dryRun, err
}

func parseMergeAIArgs(args []string) (string, []ResolvedSourceItem, []SkippedSourceItem, bool, error) {
	if len(args) == 0 || hasHelpFlag(args) {
		RenderMergeAIHelp()
		return "", nil, nil, false, nil
	}

	state := scanArgsIntoState(args)
	if state.isHelp {
		RenderMergeAIHelp()
		return "", nil, nil, false, nil
	}

	if state.configPath != "" {
		dest, resolved, skipped, err := parseJSONConfigFile(state.configPath)
		return dest, resolved, skipped, state.dryRun, err
	}

	// Case A: Named flags were used for sources or dest
	if len(state.sources) > 0 || (state.dest != "" && len(state.positionals) == 0) {
		return parseNamedFlagSources(state)
	}

	// Case B: Pure positional arguments
	if len(state.positionals) == 1 {
		dest, resolved, skipped, err := parseSingleMergeAIArg(state.positionals[0])
		return dest, resolved, skipped, state.dryRun, err
	}

	dest, resolved, skipped, err := parseMultiMergeAIArgs(state.positionals)
	if state.dest != "" {
		dest = state.dest
	}
	return dest, resolved, skipped, state.dryRun, err
}

func parseSingleMergeAIArg(arg string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	lower := strings.ToLower(arg)
	if strings.HasSuffix(lower, ".json") {
		return parseJSONConfigFile(arg)
	}

	destHint, resolved, skipped, err := ParseSourceFile(arg)
	if err != nil {
		return "", nil, nil, err
	}

	dest := destHint
	if dest == "" {
		dest = "icon-studio-v1"
		pterm.Info.Printf("[GitMap Merge-AI] No destination target specified; defaulting to: %s\n", dest)
	}

	return dest, resolved, skipped, nil
}

func parseJSONConfigFile(configPath string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	resolvedPath, errResolve := ResolveSourceFilePath(configPath)
	if errResolve != nil {
		resolvedPath = configPath
	}

	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return "", nil, nil, apperror.WrapSimple(err, fmt.Sprintf("failed to read config file %s", configPath))
	}

	var cfg MergeAIConfigFile
	if errJSON := json.Unmarshal(data, &cfg); errJSON != nil {
		return "", nil, nil, apperror.WrapSimple(errJSON, fmt.Sprintf("failed to parse JSON config %s", configPath))
	}

	destHint, resolved, skipped, errParse := ParseSourceManifestContent(strings.Join(cfg.Sources, "\n"))
	if errParse != nil {
		return "", nil, nil, errParse
	}

	dest := resolveFinalDest(cfg.Destination, destHint)
	return dest, resolved, skipped, nil
}

func isAllDestinationArgs(args []string) bool {
	for _, a := range args {
		if isTextFilePath(a) || isFullURL(a) {
			return false
		}
	}
	return true
}

func parseMultiMergeAIArgs(args []string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	// If first arg is a text file and subsequent args form destination name:
	if isTextFilePath(args[0]) && isAllDestinationArgs(args[1:]) {
		dest := strings.Join(args[1:], " ")
		_, resolved, skipped, err := ParseSourceFile(args[0])
		return dest, resolved, skipped, err
	}

	// If last arg is a text file and preceding args form destination name:
	if isTextFilePath(args[len(args)-1]) && isAllDestinationArgs(args[:len(args)-1]) {
		dest := strings.Join(args[:len(args)-1], " ")
		_, resolved, skipped, err := ParseSourceFile(args[len(args)-1])
		return dest, resolved, skipped, err
	}

	dest := args[0]
	rawSources := args[1:]

	if len(rawSources) == 1 && isTextFilePath(rawSources[0]) {
		_, resolved, skipped, err := ParseSourceFile(rawSources[0])
		return dest, resolved, skipped, err
	}

	destHint, resolved, skipped, err := parseSourceStringSlice(rawSources)
	dest = resolveFinalDest(dest, destHint)
	return dest, resolved, skipped, err
}

func expandTextFileSource(path string) []string {
	_, resFile, _, errFile := ParseSourceFile(path)
	if errFile != nil {
		return nil
	}
	var lines []string
	for _, r := range resFile {
		lines = append(lines, r.RepoURL)
	}
	return lines
}

func expandCommaSeparatedSources(raw string) []string {
	var lines []string
	for _, part := range strings.Split(raw, ",") {
		clean := strings.TrimSpace(part)
		if clean != "" {
			lines = append(lines, clean)
		}
	}
	return lines
}

func expandSourceArg(s string) []string {
	if isTextFilePath(s) {
		return expandTextFileSource(s)
	}
	if strings.Contains(s, ",") {
		return expandCommaSeparatedSources(s)
	}
	clean := strings.TrimSpace(s)
	if clean != "" {
		return []string{clean}
	}
	return nil
}

func parseSourceStringSlice(rawSources []string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	var combinedLines []string
	for _, s := range rawSources {
		combinedLines = append(combinedLines, expandSourceArg(s)...)
	}

	return ParseSourceManifestContent(strings.Join(combinedLines, "\n"))
}

func isTextFilePath(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	if strings.HasSuffix(lower, ".txt") || strings.HasSuffix(lower, ".manifest") {
		return true
	}
	if resolved, err := ResolveSourceFilePath(path); err == nil && resolved != "" {
		return true
	}
	return false
}
