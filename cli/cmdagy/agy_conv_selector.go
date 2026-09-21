package cmdagy

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// FindMatchingConversations returns all conversations matching repoRoot sorted by activity.
func FindMatchingConversations(repoRoot string) ([]AgyConvInfo, error) {
	pClean := cleanProjectWorkspace(repoRoot)
	pName := filepath.Base(pClean)
	summaries, err := scanConversationsFromSummaries(pClean, pName)
	hasSummaries := err == nil && len(summaries) > 0
	if hasSummaries {
		return summaries, nil
	}

	return findMatchingConversationsFromDBDir(pClean)
}

func findMatchingConversationsFromDBDir(pClean string) ([]AgyConvInfo, error) {
	convs, err := scanAllConversations()
	hasErr := err != nil
	if hasErr {
		return nil, apperror.WrapSimple(err, "scan conversations")
	}

	matches := filterMatchingConversations(pClean, convs)
	sortConversations(matches)

	return matches, nil
}

func filterMatchingConversations(pClean string, convs []AgyConvInfo) []AgyConvInfo {
	var matches []AgyConvInfo
	for _, c := range convs {
		isMatch := isConvPathMatch(pClean, c.CleanPath)
		if isMatch {
			matches = append(matches, c)
		}
	}

	return matches
}

func sortConversations(convs []AgyConvInfo) {
	sort.Slice(convs, func(i, j int) bool {
		return compareConvInfo(convs[i], convs[j])
	})
}

func compareConvInfo(a, b AgyConvInfo) bool {
	hasDiffUserSteps := a.UserSteps != b.UserSteps
	if hasDiffUserSteps {
		return a.UserSteps > b.UserSteps
	}

	hasDiffSteps := a.StepCount != b.StepCount
	if hasDiffSteps {
		return a.StepCount > b.StepCount
	}

	return a.ID < b.ID
}

// SelectMatchingConversation automatically picks the best matching conversation by project name and path without prompting.
func SelectMatchingConversation(repoRoot string) (AgyConvInfo, error) {
	convs, err := FindMatchingConversations(repoRoot)
	hasErr := err != nil
	if hasErr {
		return AgyConvInfo{}, err
	}
	isZeroMatches := len(convs) == 0
	if isZeroMatches {
		return AgyConvInfo{}, apperror.NewSimple("no matching conversation found for workspace", "E9020")
	}

	return convs[0], nil
}

// SelectMatchingConversationWithIO selects matching conversation with specified IO streams.
func SelectMatchingConversationWithIO(repoRoot string, r io.Reader, w io.Writer, isInteractive bool) (AgyConvInfo, error) {
	convs, err := FindMatchingConversations(repoRoot)
	hasErr := err != nil
	if hasErr {
		return AgyConvInfo{}, err
	}

	return PromptSelectConversation(convs, r, w, isInteractive)
}

// PromptSelectConversation handles conversation selection with 0, 1, or multiple matches.
func PromptSelectConversation(convs []AgyConvInfo, r io.Reader, w io.Writer, isInteractive bool) (AgyConvInfo, error) {
	isZeroMatches := len(convs) == 0
	if isZeroMatches {
		return AgyConvInfo{}, apperror.NewSimple("no matching conversation found for workspace", "E9020")
	}

	sortConversations(convs)
	hasMultiple := len(convs) > 1
	isPromptNeeded := isInteractive && hasMultiple
	if isPromptNeeded {
		return promptUserChoice(convs, r, w)
	}

	return convs[0], nil
}

func promptUserChoice(convs []AgyConvInfo, r io.Reader, w io.Writer) (AgyConvInfo, error) {
	printConvList(convs, w)
	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	hasErr := err != nil && len(input) == 0
	if hasErr {
		return convs[0], nil
	}

	return resolveSelectedConv(convs, input), nil
}

func printConvList(convs []AgyConvInfo, w io.Writer) {
	fmt.Fprintf(w, "Multiple active conversations found:\n")
	for i, c := range convs {
		printConvOption(w, i, c)
	}

	fmt.Fprintf(w, "Select conversation [1-%d] (default 1): ", len(convs))
}

func printConvOption(w io.Writer, idx int, c AgyConvInfo) {
	isDefault := idx == 0
	if isDefault {
		fmt.Fprintf(w, "  [%d] %s (user steps: %d, steps: %d) [Default]\n", idx+1, c.ID, c.UserSteps, c.StepCount)

		return
	}

	fmt.Fprintf(w, "  [%d] %s (user steps: %d, steps: %d)\n", idx+1, c.ID, c.UserSteps, c.StepCount)
}

func resolveSelectedConv(convs []AgyConvInfo, input string) AgyConvInfo {
	idx := parseConvIndex(input, len(convs))

	return convs[idx]
}

func parseConvIndex(input string, max int) int {
	trimmed := strings.TrimSpace(input)
	idx, err := strconv.Atoi(trimmed)
	hasErr := err != nil
	if hasErr {
		return 0
	}

	isValid := idx >= 1 && idx <= max
	if isValid {
		return idx - 1
	}

	return 0
}
