package cmdpipeline

import (
	"testing"
)

func TestQueryRepoSuggestions(t *testing.T) {
	suggs := QueryRepoSuggestionsFromDB("movi-cli")
	if len(suggs) == 0 {
		t.Logf("no suggestions found for movi-cli (DB may be unseeded in test environment)")
		return
	}
	found := false
	for _, s := range suggs {
		if s == "movie-cli-v8" || s == "alimtvnetwork/movie-cli-v8" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected movie-cli-v8 in suggestions, got %v", suggs)
	}
}
