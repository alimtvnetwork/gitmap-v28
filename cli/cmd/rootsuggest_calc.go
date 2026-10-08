package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

type commandScore struct {
	cmd  string
	dist int
}

func collectTopCommandCandidates() []string {
	seen := make(map[string]bool)
	var candidates []string

	candidates = appendCommandSlice(candidates, seen, primaryTopCommands)
	candidates = appendCommandSlice(candidates, seen, completion.AllCommands())
	candidates = appendCommandSlice(candidates, seen, []string{"run", "run-until"})
	candidates = appendMacroCandidates(candidates, seen)

	return candidates
}

func appendCommandSlice(candidates []string, seen map[string]bool, items []string) []string {
	for _, item := range items {
		clean := strings.TrimSpace(item)
		hasItem := len(clean) > 0
		hasSeen := seen[clean]

		if hasItem && !hasSeen {
			seen[clean] = true
			candidates = append(candidates, clean)
		}
	}

	return candidates
}

func appendMacroCandidates(candidates []string, seen map[string]bool) []string {
	macroList := macro.ListMacros()
	isSuccess := macroList.IsSuccess()

	if isSuccess {
		var names []string
		for _, m := range macroList.Data {
			names = append(names, m.Name)
		}

		return appendCommandSlice(candidates, seen, names)
	}

	return candidates
}

func suggestTopLevelCommands(command string) []string {
	norm := strings.ToLower(strings.TrimSpace(command))
	hasNorm := len(norm) > 0
	if !hasNorm {
		return nil
	}

	isPlease := norm == "pleae" || norm == "please"
	if isPlease {
		return []string{"pull", "release", "pull-release"}
	}

	candidates := collectTopCommandCandidates()
	scores := rankCandidateCommands(norm, candidates)
	best := selectBestSuggestions(scores)
	hasBest := len(best) > 0
	if hasBest {
		return best
	}

	return suggestByIntentKeywords(norm)
}

func suggestByIntentKeywords(norm string) []string {
	if strings.Contains(norm, "key") || strings.Contains(norm, "auth") {
		return []string{"deploy-keys-all", "ssh key add", "ssh copy-id"}
	}
	if strings.Contains(norm, "ssh") || strings.Contains(norm, "node") || strings.Contains(norm, "host") {
		return []string{"ssh ls", "ssh join", "ssh deploy-keys"}
	}
	if strings.Contains(norm, "clean") || strings.Contains(norm, "clear") || strings.Contains(norm, "cache") || strings.Contains(norm, "term") {
		return []string{"clear devtools", "clear terminal", "clean-dev"}
	}
	if strings.Contains(norm, "install") || strings.Contains(norm, "setup") {
		return []string{"install ls", "ssh install gitmap -t all", "deploy-bin all"}
	}
	if strings.Contains(norm, "fail") || strings.Contains(norm, "err") || strings.Contains(norm, "unknown") {
		return []string{"failed-commands", "errors", "pe"}
	}

	return nil
}

func rankCandidateCommands(input string, candidates []string) []commandScore {
	var scores []commandScore
	for _, c := range candidates {
		d := cmdclone.Levenshtein(input, c)
		isSub := strings.Contains(c, input)
		if isSub {
			d = 1
		}

		if d <= 3 {
			scores = append(scores, commandScore{cmd: c, dist: d})
		}
	}

	sort.SliceStable(scores, func(i, j int) bool {
		if scores[i].dist != scores[j].dist {
			return scores[i].dist < scores[j].dist
		}
		prefixI := strings.HasPrefix(scores[i].cmd, input)
		prefixJ := strings.HasPrefix(scores[j].cmd, input)
		if prefixI != prefixJ {
			return prefixI
		}
		diffI := candidateLenDiff(scores[i].cmd, input)
		diffJ := candidateLenDiff(scores[j].cmd, input)
		return diffI < diffJ
	})

	return scores
}

func candidateLenDiff(candidate, input string) int {
	diff := len(candidate) - len(input)
	if diff < 0 {
		return -diff
	}
	return diff
}

func selectBestSuggestions(scores []commandScore) []string {
	var result []string
	seen := make(map[string]bool)

	for _, sc := range scores {
		hasSeen := seen[sc.cmd]
		if !hasSeen {
			seen[sc.cmd] = true
			result = append(result, sc.cmd)
		}

		hasLimit := len(result) >= 3
		if hasLimit {
			break
		}
	}

	return result
}
