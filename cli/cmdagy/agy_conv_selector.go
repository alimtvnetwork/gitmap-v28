package cmdagy

import (
	"io"
	"path/filepath"
	"sort"

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

// SelectMatchingConversationWithIO selects matching conversation automatically without prompting.
func SelectMatchingConversationWithIO(repoRoot string, r io.Reader, w io.Writer, isInteractive bool) (AgyConvInfo, error) {
	return SelectMatchingConversation(repoRoot)
}

// PromptSelectConversation handles conversation selection automatically without prompting the user.
func PromptSelectConversation(convs []AgyConvInfo, r io.Reader, w io.Writer, isInteractive bool) (AgyConvInfo, error) {
	isZeroMatches := len(convs) == 0
	if isZeroMatches {
		return AgyConvInfo{}, apperror.NewSimple("no matching conversation found for workspace", "E9020")
	}

	sortConversations(convs)

	return convs[0], nil
}
