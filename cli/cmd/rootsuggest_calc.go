package cmd

import (
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

type commandScore struct {
	cmd  string
	dist int
}

func collectTopCommandCandidates() []string {
	candidates := append([]string{}, primaryTopCommands...)
	candidates = append(candidates, "run", "run-until")
	macroList := macro.ListMacros()
	if macroList.IsSuccess() {
		for _, m := range macroList.Data {
			candidates = append(candidates, m.Name)
		}
	}

	return candidates
}

func suggestTopLevelCommands(command string) []string {
	norm := strings.ToLower(strings.TrimSpace(command))
	if norm == "" {
		return nil
	}

	if norm == "pleae" || norm == "please" {
		return []string{"pull", "release", "pull-release"}
	}

	candidates := collectTopCommandCandidates()
	scores := rankCandidateCommands(norm, candidates)
	best := selectBestSuggestions(scores)
	if len(best) > 0 {
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
		d := levenshtein(input, c)
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
		if !seen[sc.cmd] {
			seen[sc.cmd] = true
			result = append(result, sc.cmd)
		}

		if len(result) >= 3 {
			break
		}
	}

	return result
}
