package cmdsummary

import (
	"strings"
	"testing"
	"time"
)

func TestParseSummaryArgs(t *testing.T) {
	dir, limit := parseSummaryArgs([]string{"my-repo", "5"})
	if dir != "my-repo" {
		t.Fatalf("expected dir 'my-repo', got '%s'", dir)
	}
	if limit != 5 {
		t.Fatalf("expected limit 5, got %d", limit)
	}

	dirDefault, limitDefault := parseSummaryArgs([]string{})
	if dirDefault != "." {
		t.Fatalf("expected default dir '.', got '%s'", dirDefault)
	}
	if limitDefault != 8 {
		t.Fatalf("expected default limit 8, got %d", limitDefault)
	}
}

func TestParseFullSummaryOptions(t *testing.T) {
	opts := parseFullSummaryOptions([]string{"5", "--json", "+pe", "--force-all"})
	if opts.ReleasesCount != 5 {
		t.Fatalf("expected ReleasesCount 5, got %d", opts.ReleasesCount)
	}
	if !opts.IsJSON {
		t.Fatalf("expected IsJSON true")
	}
	if !opts.WithPE {
		t.Fatalf("expected WithPE true")
	}
	if !opts.ForceAll {
		t.Fatalf("expected ForceAll true")
	}
}

func TestSynthesizeReleaseGist(t *testing.T) {
	commitLog := "a1b2c3d Fix token refresh race condition\ne5f6g7h Add Casbin RBAC schema support\n9z8y7x6 Update Go dependencies"
	gist := synthesizeReleaseGist(commitLog)
	if !strings.Contains(gist, "Fix token refresh race condition") {
		t.Fatalf("expected gist to contain token fix, got: %s", gist)
	}
	if !strings.Contains(gist, "Add Casbin RBAC schema support") {
		t.Fatalf("expected gist to contain Casbin support, got: %s", gist)
	}

	words := strings.Fields(gist)
	if len(words) > 200 {
		t.Fatalf("expected words <= 200, got %d", len(words))
	}
}

func TestExtractHeatedFilesFromDiff(t *testing.T) {
	numstat := "142\t18\tinternal/auth/jwt.go\n35\t12\tpkg/db/sqlite.go\n4\t1\treadme.md\n"
	metrics := extractHeatedFilesFromDiff(numstat)
	if len(metrics) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(metrics))
	}
	if metrics[0].Path != "internal/auth/jwt.go" {
		t.Fatalf("expected top heated file to be jwt.go, got: %s", metrics[0].Path)
	}
	if metrics[0].ChangesCount != 160 {
		t.Fatalf("expected total changes 160, got %d", metrics[0].ChangesCount)
	}
}

func TestSplitDBCacheRoundTrip(t *testing.T) {
	repoURL := "https://github.com/mock/test-summary-repo.git"
	slug := "test-summary-repo"
	localPath := "repos/test-summary-repo"

	rec := ReleaseSummaryRecord{
		TagName:       "v1.0.0",
		TagCommitHash: "abcdef123456",
		ReleaseDate:   "2026-10-10",
		SummaryGist:   "Initial release with split-db caching and heated file metrics.",
		WordCount:     9,
		HeatedFiles: []HeatedFileMetric{
			{
				Path:         "main.go",
				ChangesCount: 50,
				Insertions:   40,
				Deletions:    10,
				Description:  "Feature expansion",
			},
		},
	}

	if err := SaveCachedRelease(repoURL, slug, localPath, rec); err != nil {
		t.Fatalf("SaveCachedRelease failed: %v", err)
	}

	cached, isHit := GetCachedRelease(repoURL, "v1.0.0", "abcdef123456")
	if !isHit {
		t.Fatalf("expected cache hit for v1.0.0, but got miss")
	}
	if cached.SummaryGist != rec.SummaryGist {
		t.Fatalf("expected gist '%s', got '%s'", rec.SummaryGist, cached.SummaryGist)
	}
	if len(cached.HeatedFiles) != 1 || cached.HeatedFiles[0].Path != "main.go" {
		t.Fatalf("expected 1 heated file 'main.go', got: %+v", cached.HeatedFiles)
	}
}

func TestPipelineErrorTruncation(t *testing.T) {
	// Construct 35 lines of stack trace
	var rawLines []string
	for i := 1; i <= 35; i++ {
		rawLines = append(rawLines, strings.Repeat("x", 10)+"-line-"+strings.Repeat("0", 2)+string(rune('A'+(i%26))))
	}

	traceLines := rawLines
	if len(traceLines) > 25 {
		traceLines = traceLines[len(traceLines)-25:]
	}

	if len(traceLines) != 25 {
		t.Fatalf("expected truncated lines to be 25, got %d", len(traceLines))
	}
	if traceLines[0] != rawLines[10] {
		t.Fatalf("expected first line of truncated trace to match index 10, got '%s'", traceLines[0])
	}
	if traceLines[24] != rawLines[34] {
		t.Fatalf("expected last line of truncated trace to match index 34, got '%s'", traceLines[24])
	}
}

func TestGistWordCeilingAndFileListSuppression(t *testing.T) {
	// Test file path suppression
	commitLogWithFiles := "abc1234 Updated internal/auth/token.go, pkg/db/sqlite.go with new schema\ndef5678 Fixed session cache in pkg/cache/redis.go"
	gist := synthesizeReleaseGist(commitLogWithFiles)

	if strings.Contains(gist, "internal/auth/token.go") {
		t.Fatalf("expected raw file path internal/auth/token.go to be suppressed from gist, got: %s", gist)
	}
	if strings.Contains(gist, "pkg/cache/redis.go") {
		t.Fatalf("expected raw file path pkg/cache/redis.go to be suppressed from gist, got: %s", gist)
	}

	// Test hard word ceiling <= 200 words
	var longLines []string
	for i := 0; i < 50; i++ {
		longLines = append(longLines, "abc0000 Feature enhancement adding additional capability for modular subsystem components")
	}
	longGist := synthesizeReleaseGist(strings.Join(longLines, "\n"))
	words := strings.Fields(longGist)
	if len(words) > 200 {
		t.Fatalf("expected word ceiling <= 200, got %d words", len(words))
	}
	if !strings.HasSuffix(longGist, "...") {
		t.Fatalf("expected truncated long gist to end with ellipsis, got: %s", longGist)
	}
}

func TestSummarizeFileChurnExtensions(t *testing.T) {
	cases := []struct {
		path string
		ins  int
		del  int
		desc string
	}{
		{"cmd/main.go", 100, 10, "Feature expansion and new logic implementation"},
		{"cmd/main.go", 10, 100, "Refactoring, pruning, and dead code elimination"},
		{"cmd/main.go", 20, 20, "Iterative feature enhancements and logic adjustments"},
		{"docs/arch.md", 5, 2, "Documentation and architecture updates"},
		{"config.json", 10, 1, "Configuration, schema, and dependency updates"},
		{"schema.sql", 50, 0, "Database schema and migration updates"},
		{"service.proto", 15, 2, "API protocol and RPC schema updates"},
		{"ui/styles.css", 30, 5, "User interface styling and markup updates"},
		{"scripts/build.sh", 8, 2, "Automation scripting and build workflow updates"},
	}

	for _, c := range cases {
		got := summarizeFileChurn(c.path, c.ins, c.del)
		if got != c.desc {
			t.Errorf("summarizeFileChurn(%s) = '%s', want '%s'", c.path, got, c.desc)
		}
	}
}

func TestPorcelainStatusParsing(t *testing.T) {
	sample := "?? new_file.txt\nM  staged_and_modified.go\n M unstaged.go\nA  added.go\n"
	untracked, modified, staged, pending := parsePorcelainStatusLines(sample)

	if untracked != 1 {
		t.Fatalf("expected 1 untracked, got %d", untracked)
	}
	if staged != 2 { // "M " and "A "
		t.Fatalf("expected 2 staged, got %d", staged)
	}
	if modified != 1 { // " M"
		t.Fatalf("expected 1 modified, got %d", modified)
	}
	if len(pending) != 4 {
		t.Fatalf("expected 4 pending files, got %d", len(pending))
	}
}

func TestUpdateAndGetCachedHeadState(t *testing.T) {
	repoURL := "https://github.com/mock/test-head-state.git"
	slug := "test-head-state"
	localPath := "repos/test-head-state"
	headHash := "feedbeef12345678"

	if err := UpdateRepoHeadState(repoURL, slug, localPath, headHash, true, 3, time.Now()); err != nil {
		t.Fatalf("UpdateRepoHeadState failed: %v", err)
	}

	headState, found := GetCachedHeadState(repoURL)
	if !found {
		t.Fatalf("expected GetCachedHeadState to find record for %s", repoURL)
	}
	if headState.HeadCommitHash != headHash {
		t.Fatalf("expected headHash '%s', got '%s'", headHash, headState.HeadCommitHash)
	}
	if !headState.IsDirty {
		t.Fatalf("expected isDirty true")
	}
	if headState.DirtyFilesCount != 3 {
		t.Fatalf("expected dirtyFilesCount 3, got %d", headState.DirtyFilesCount)
	}
}

func BenchmarkSplitDBCacheHit(b *testing.B) {
	repoURL := "https://github.com/mock/bench-summary-repo.git"
	slug := "bench-summary-repo"
	localPath := "repos/bench-summary-repo"
	tag := "v2.0.0"
	tagHash := "112233445566"

	rec := ReleaseSummaryRecord{
		TagName:       tag,
		TagCommitHash: tagHash,
		ReleaseDate:   "2026-10-10",
		SummaryGist:   "Benchmark cached gist for high performance sub-15ms retrieval.",
		WordCount:     9,
		HeatedFiles: []HeatedFileMetric{
			{
				Path:         "bench.go",
				ChangesCount: 20,
				Insertions:   15,
				Deletions:    5,
				Description:  "Feature expansion",
			},
		},
	}

	_ = SaveCachedRelease(repoURL, slug, localPath, rec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cached, hit := GetCachedRelease(repoURL, tag, tagHash)
		if !hit || cached == nil {
			b.Fatalf("cache miss during benchmark")
		}
	}
}
