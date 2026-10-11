// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

type parsedFlagState struct {
	dest        string
	sources     []string
	configPath  string
	positionals []string
	dryRun      bool
	isHelp      bool
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		lower := strings.ToLower(strings.TrimSpace(a))
		if lower == "-h" || lower == "--help" || lower == "help" {
			return true
		}
	}
	return false
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
		dest := state.dest
		if dest == "" {
			dest = "icon-studio-v1"
		}
		destHint, resolved, skipped, err := parseSourceStringSlice(state.sources)
		if dest == "icon-studio-v1" && destHint != "" {
			dest = destHint
		}
		return dest, resolved, skipped, state.dryRun, err
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

func scanArgsIntoState(args []string) parsedFlagState {
	var s parsedFlagState
	i := 0
	for i < len(args) {
		a := strings.TrimSpace(args[i])
		lower := strings.ToLower(a)

		if lower == "-h" || lower == "--help" || lower == "help" {
			s.isHelp = true
			i++
			continue
		}

		if lower == "--dry-run" {
			s.dryRun = true
			i++
			continue
		}

		if lower == "-d" || lower == "--destination" || lower == "--dest" || lower == "--target" {
			if i+1 < len(args) {
				s.dest = strings.TrimSpace(args[i+1])
				i += 2
				continue
			}
		}

		if strings.HasPrefix(lower, "--destination=") || strings.HasPrefix(lower, "-d=") || strings.HasPrefix(lower, "--dest=") || strings.HasPrefix(lower, "--target=") {
			parts := strings.SplitN(a, "=", 2)
			if len(parts) == 2 {
				s.dest = strings.TrimSpace(parts[1])
			}
			i++
			continue
		}

		if lower == "-s" || lower == "--source" || lower == "-u" || lower == "--url" {
			if i+1 < len(args) {
				s.sources = append(s.sources, strings.TrimSpace(args[i+1]))
				i += 2
				continue
			}
		}

		if strings.HasPrefix(lower, "--source=") || strings.HasPrefix(lower, "-s=") || strings.HasPrefix(lower, "--url=") || strings.HasPrefix(lower, "-u=") {
			parts := strings.SplitN(a, "=", 2)
			if len(parts) == 2 {
				s.sources = append(s.sources, strings.TrimSpace(parts[1]))
			}
			i++
			continue
		}

		if lower == "-c" || lower == "--config" {
			if i+1 < len(args) {
				s.configPath = strings.TrimSpace(args[i+1])
				i += 2
				continue
			}
		}

		if strings.HasPrefix(lower, "--config=") || strings.HasPrefix(lower, "-c=") {
			parts := strings.SplitN(a, "=", 2)
			if len(parts) == 2 {
				s.configPath = strings.TrimSpace(parts[1])
			}
			i++
			continue
		}

		s.positionals = append(s.positionals, a)
		i++
	}

	return s
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
		return "", nil, nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg MergeAIConfigFile
	if errJSON := json.Unmarshal(data, &cfg); errJSON != nil {
		return "", nil, nil, fmt.Errorf("failed to parse JSON config %s: %w", configPath, errJSON)
	}

	destHint, resolved, skipped, errParse := ParseSourceManifestContent(strings.Join(cfg.Sources, "\n"))
	if errParse != nil {
		return "", nil, nil, errParse
	}

	dest := cfg.Destination
	if dest == "" {
		dest = destHint
	}
	if dest == "" {
		dest = "icon-studio-v1"
	}

	return dest, resolved, skipped, nil
}

func parseMultiMergeAIArgs(args []string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	// If first arg is a text file and subsequent args form destination name:
	if isTextFilePath(args[0]) {
		otherAreDest := true
		for _, a := range args[1:] {
			if isTextFilePath(a) || isFullURL(a) {
				otherAreDest = false
				break
			}
		}
		if otherAreDest {
			dest := strings.Join(args[1:], " ")
			_, resolved, skipped, err := ParseSourceFile(args[0])
			return dest, resolved, skipped, err
		}
	}

	// If last arg is a text file and preceding args form destination name:
	if isTextFilePath(args[len(args)-1]) {
		precedingAreDest := true
		for _, a := range args[:len(args)-1] {
			if isTextFilePath(a) || isFullURL(a) {
				precedingAreDest = false
				break
			}
		}
		if precedingAreDest {
			dest := strings.Join(args[:len(args)-1], " ")
			_, resolved, skipped, err := ParseSourceFile(args[len(args)-1])
			return dest, resolved, skipped, err
		}
	}

	dest := args[0]
	rawSources := args[1:]

	if len(rawSources) == 1 && isTextFilePath(rawSources[0]) {
		_, resolved, skipped, err := ParseSourceFile(rawSources[0])
		return dest, resolved, skipped, err
	}

	destHint, resolved, skipped, err := parseSourceStringSlice(rawSources)
	if dest == "" {
		dest = destHint
	}
	if dest == "" {
		dest = "icon-studio-v1"
	}
	return dest, resolved, skipped, err
}

func parseSourceStringSlice(rawSources []string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	var combinedLines []string
	for _, s := range rawSources {
		if isTextFilePath(s) {
			_, resFile, _, errFile := ParseSourceFile(s)
			if errFile == nil {
				for _, r := range resFile {
					combinedLines = append(combinedLines, r.RepoURL)
				}
				continue
			}
		}

		if strings.Contains(s, ",") {
			for _, part := range strings.Split(s, ",") {
				if clean := strings.TrimSpace(part); clean != "" {
					combinedLines = append(combinedLines, clean)
				}
			}
			continue
		}

		if clean := strings.TrimSpace(s); clean != "" {
			combinedLines = append(combinedLines, clean)
		}
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
