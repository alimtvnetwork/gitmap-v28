package cmdsummary

import (
	"strings"
	"testing"
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
