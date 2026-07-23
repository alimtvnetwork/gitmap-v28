package orchestrator

import (
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestApplyExclusions_FiltersFolderAndFile(t *testing.T) {
	files := []string{
		"src/main.go",
		"vendor/lib/pkg.go",
		"node_modules/x/y.js",
		"README.md",
		"docs/secret.md",
	}
	rules := []profile.Exclusion{
		{Kind: constants.CommitInExclusionKindPathFolder, Value: "vendor"},
		{Kind: constants.CommitInExclusionKindPathFolder, Value: "node_modules"},
		{Kind: constants.CommitInExclusionKindPathFile, Value: "docs/secret.md"},
	}
	got := applyExclusions(files, rules)
	want := []string{"src/main.go", "README.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("applyExclusions: got %v want %v", got, want)
	}
}

func TestApplyExclusions_NoRulesReturnsInputUnchanged(t *testing.T) {
	files := []string{"a.go", "b.go"}
	got := applyExclusions(files, nil)
	if !reflect.DeepEqual(got, files) {
		t.Fatalf("expected pass-through, got %v", got)
	}
}

func TestApplyExclusions_FolderMatchesNestedSegment(t *testing.T) {
	rules := []profile.Exclusion{
		{Kind: constants.CommitInExclusionKindPathFolder, Value: "dist"},
	}
	got := applyExclusions([]string{"pkg/dist/out.js", "pkg/src/in.js"}, rules)
	want := []string{"pkg/src/in.js"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("nested folder match: got %v want %v", got, want)
	}
}

func TestApplyExclusions_FileExactMatchOnly(t *testing.T) {
	rules := []profile.Exclusion{
		{Kind: constants.CommitInExclusionKindPathFile, Value: "config.json"},
	}
	got := applyExclusions([]string{"config.json", "src/config.json"}, rules)
	want := []string{"src/config.json"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("file exact match: got %v want %v", got, want)
	}
}
