package cmdagy

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolvePromptReadConvID(args []string) (string, error) {
	if len(args) > 0 && len(args[0]) > 0 {
		return resolveTargetConvID(args[0])
	}

	return resolveActiveOrLatestConvID()
}

func resolveActiveOrLatestConvID() (string, error) {
	cwd, _ := os.Getwd()
	activeID := findMatchingActiveConvID(cwd)
	if len(activeID) > 0 {
		return activeID, nil
	}

	latestID, _ := findLatestModifiedTranscript()
	if len(latestID) > 0 {
		return latestID, nil
	}

	return "", apperror.NewSimple("no active or recent conversation found in Antigravity brain", "E9021")
}

func resolveTargetConvID(target string) (string, error) {
	id, isFound := searchConvTarget(target)
	if isFound {
		return id, nil
	}

	return "", apperror.NewSimple(fmt.Sprintf("conversation not found for target: %s", target), "E9020")
}

func searchConvTarget(target string) (string, bool) {
	exactID, isExact := matchDirectConvDir(target)
	if isExact {
		return exactID, true
	}
	prefixID, hasPrefixMatch := matchConvPrefix(target)
	if hasPrefixMatch {
		return prefixID, true
	}

	return matchConvByWorkspace(target)
}

func matchDirectConvDir(target string) (string, bool) {
	tPath := resolveTranscriptPathForConv(target)
	hasTranscript := checkFileExists(tPath)

	return target, hasTranscript
}

func matchConvPrefix(target string) (string, bool) {
	brainDir, err := GetBrainLogsDirPath()
	if err != nil {
		return "", false
	}

	entries, err := os.ReadDir(brainDir)
	if err != nil {
		return "", false
	}

	return findPrefixMatchingDir(entries, target)
}

func findPrefixMatchingDir(entries []os.DirEntry, target string) (string, bool) {
	targetLow := strings.ToLower(target)
	for _, e := range entries {
		isMatch := e.IsDir() && strings.HasPrefix(strings.ToLower(e.Name()), targetLow)
		if isMatch {
			return e.Name(), true
		}
	}

	return "", false
}

func matchConvByWorkspace(target string) (string, bool) {
	convID := findMatchingActiveConvID(target)
	hasConv := len(convID) > 0

	return convID, hasConv
}
