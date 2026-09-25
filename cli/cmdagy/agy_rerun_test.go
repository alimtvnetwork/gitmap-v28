package cmdagy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func makeTestAgyProject(id, name, uri, updatedAt string) AgyProject {
	p := AgyProject{
		ID:        id,
		Name:      name,
		UpdatedAt: updatedAt,
	}
	if uri != "" {
		p.ProjectResources = &AgyProjectResources{
			Resources: []AgyResource{
				{GitFolder: &AgyGitFolder{FolderURI: uri}},
			},
		}
	}
	return p
}

func TestResolveClosestActiveProject_FuzzyMatch(t *testing.T) {
	fileScheme := "file:" + "//" + "/"
	projects := []AgyProject{
		makeTestAgyProject("proj-1", "gitmap", fileScheme+"virtual/gitmap", "2026-09-24T00:00:00Z"),
		makeTestAgyProject("proj-2", "wp-exam", fileScheme+"virtual/wp-exam", "2026-09-21T00:00:00Z"),
		makeTestAgyProject("proj-3", "movie-cli", fileScheme+"virtual/movie-cli", "2026-09-20T00:00:00Z"),
	}

	testCases := []struct {
		input       string
		expectedID  string
		expectedSeq int
	}{
		{"wp-exam", "proj-2", 2},
		{"wp-xampp", "proj-2", 2},
		{"wpexam", "proj-2", 2},
		{"exam", "proj-2", 2},
		{"movie", "proj-3", 3},
		{"gitmap", "proj-1", 1},
		{"1", "proj-1", 1},
		{"2", "proj-2", 2},
		{"3", "proj-3", 3},
	}

	for _, tc := range testCases {
		p, seq := resolveClosestActiveProject(projects, tc.input)
		if p.ID != tc.expectedID {
			t.Errorf("for target %q: expected ID %q, got %q", tc.input, tc.expectedID, p.ID)
		}
		if seq != tc.expectedSeq {
			t.Errorf("for target %q: expected seq %d, got %d", tc.input, tc.expectedSeq, seq)
		}
	}
}

func TestResolveClosestActiveProject_CwdMatch(t *testing.T) {
	tempDir := t.TempDir()
	pDir := filepath.Join(tempDir, "sample-project")
	_ = os.MkdirAll(pDir, 0755)

	fileScheme := "file:" + "//" + "/"
	cleanPDir := strings.TrimPrefix(filepath.ToSlash(pDir), "/")
	projects := []AgyProject{
		makeTestAgyProject("other-1", "other-one", fileScheme+"virtual/other", "2026-08-01T00:00:00Z"),
		makeTestAgyProject("match-1", "sample-project", fileScheme+cleanPDir, "2026-08-01T00:00:00Z"),
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(pDir)
	defer func() { _ = os.Chdir(origWd) }()

	p, seq := resolveClosestActiveProject(projects, "")
	if p.ID != "match-1" {
		t.Errorf("expected cwd match 'match-1', got %q", p.ID)
	}
	if seq != 2 {
		t.Errorf("expected sequence 2, got %d", seq)
	}

	p1, seq1 := resolveClosestActiveProject(projects, "1")
	if p1.ID != "match-1" || seq1 != 2 {
		t.Errorf("expected target '1' to match CWD project 'match-1', got %q seq %d", p1.ID, seq1)
	}
}

func TestSortProjectsByActivityAndPins_Ordering(t *testing.T) {
	projects := []AgyProject{
		{ID: "proj-old", Name: "old-project", UpdatedAt: "2026-08-01T00:00:00Z"},
		{ID: "proj-act", Name: "act-project", UpdatedAt: "2026-08-05T00:00:00Z"},
	}

	entries := []projectSortEntry{
		{project: projects[0], pinnedRank: 999999, lastActTime: "2026-08-01T00:00:00Z"},
		{project: projects[1], pinnedRank: 999999, lastActTime: "2026-09-24T10:00:00Z"},
	}

	isSecondNewer := compareProjectEntries(entries[1], entries[0])
	if isSecondNewer == false {
		t.Errorf("expected act-project with recent activity to rank before old-project")
	}

	pinnedOldEntry := projectSortEntry{project: projects[0], pinnedRank: 1, lastActTime: "2026-08-01T00:00:00Z"}
	isRecentFirst := compareProjectEntries(entries[1], pinnedOldEntry)
	if isRecentFirst == false {
		t.Errorf("expected recent activity to rank before idle pinned project")
	}

	pinnedEqualEntry := projectSortEntry{project: projects[0], pinnedRank: 1, lastActTime: "2026-08-01T00:00:00Z"}
	unpinnedEqualEntry := projectSortEntry{project: projects[1], pinnedRank: 999999, lastActTime: "2026-08-01T00:00:00Z"}
	isPinnedTieBreak := compareProjectEntries(pinnedEqualEntry, unpinnedEqualEntry)
	if isPinnedTieBreak == false {
		t.Errorf("expected pinned entry to break tie when activity timestamps are equal")
	}

	runningEntry := projectSortEntry{project: projects[0], isRunning: true, lastActTime: "2026-08-01T00:00:00Z"}
	idleEntry := projectSortEntry{project: projects[1], isRunning: false, lastActTime: "2026-09-24T10:00:00Z"}
	isRunningFirst := compareProjectEntries(runningEntry, idleEntry)
	if isRunningFirst == false {
		t.Errorf("expected actively running project to rank before idle project")
	}

	cwdEntry := projectSortEntry{project: projects[1], isCwdMatch: true, isRunning: false}
	isCwdFirst := compareProjectEntries(cwdEntry, runningEntry)
	if isCwdFirst == false {
		t.Errorf("expected CWD match to rank first")
	}
}

func TestRerunDefaultFlags_NonDestructive(t *testing.T) {
	flag := agyRerunCmd.Flag("restart")
	if flag == nil {
		t.Fatalf("expected --restart flag to exist")
	}
	if flag.DefValue != "false" {
		t.Errorf("expected --restart flag default to be 'false', got %q", flag.DefValue)
	}

	newConvFlag := agyRerunCmd.Flag("new-conversation")
	if newConvFlag == nil {
		t.Fatalf("expected --new-conversation flag to exist")
	}
	if newConvFlag.DefValue != "true" {
		t.Errorf("expected --new-conversation flag default to be 'true', got %q", newConvFlag.DefValue)
	}
}

func TestFormatPromptWithMedia_BothLinkAndMarkdown(t *testing.T) {
	content := "Test prompt with screenshot"
	imgURI := filepath.ToSlash(filepath.Join(t.TempDir(), "test.png"))
	media := []rawTranscriptMedia{
		{MimeType: "image/png", URI: imgURI},
	}

	res := FormatPromptWithMedia(content, media)
	if !containsString(res, "Attached Picture 1: "+imgURI) {
		t.Errorf("expected text list item, got:\n%s", res)
	}
	expectedMD := "![" + "Picture 1](file:" + "//" + "/" + imgURI + ")"
	if !containsString(res, expectedMD) {
		t.Errorf("expected markdown image embed, got:\n%s", res)
	}
}

func containsString(haystack, needle string) bool {
	return filepath.FromSlash(haystack) != "" && filepath.FromSlash(needle) != "" && (len(haystack) >= len(needle)) && (haystack == needle || (len(needle) > 0 && (needle == haystack || (len(haystack) > len(needle) && (haystack[:len(needle)] == needle || haystack[len(haystack)-len(needle):] == needle || containsSubstring(haystack, needle))))))
}

func containsSubstring(h, n string) bool {
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return true
		}
	}
	return false
}
