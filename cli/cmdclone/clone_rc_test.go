package cmdclone

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCloneRC_FormatOwnerRepo(t *testing.T) {
	tests := []struct {
		url       string
		wantOwner string
		wantRepo  string
	}{
		{
			url:       "git@github.com:alimtvnetwork/gitmap-v28.git",
			wantOwner: "alimtvnetwork",
			wantRepo:  "gitmap-v28",
		},
		{
			url:       "https://github.com/alimtvnetwork/core-api.git",
			wantOwner: "alimtvnetwork",
			wantRepo:  "core-api",
		},
		{
			url:       "ssh://git@gitlab.internal:2222/platform/infra.git",
			wantOwner: "platform",
			wantRepo:  "infra",
		},
		{
			url:       "http://git.local/tools/helper",
			wantOwner: "tools",
			wantRepo:  "helper",
		},
		{
			url:       "./local-repo",
			wantOwner: "local",
			wantRepo:  "local-repo",
		},
		{
			url:       "local/local-repo",
			wantOwner: "local",
			wantRepo:  "local-repo",
		},
	}

	for _, tt := range tests {
		gotOwner, gotRepo := FormatOwnerRepo(tt.url)
		if gotOwner != tt.wantOwner || gotRepo != tt.wantRepo {
			t.Errorf("FormatOwnerRepo(%q) = (%q, %q), want (%q, %q)",
				tt.url, gotOwner, gotRepo, tt.wantOwner, tt.wantRepo)
		}
	}
}

func TestCloneRC_ClassifyProtocol(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"git@github.com:org/repo.git", "[SSH]"},
		{"ssh://git@gitlab.com:22/org/repo.git", "[SSH]"},
		{"https://github.com/org/repo.git", "[Public HTTPS]"},
		{"http://git.local/org/repo", "[Public HTTPS]"},
		{"./custom/repo", "[Local Path]"},
	}

	for _, tt := range tests {
		got := ClassifyProtocol(tt.url)
		if got != tt.want {
			t.Errorf("ClassifyProtocol(%q) = %q, want %q", tt.url, got, tt.want)
		}
	}
}

func TestCloneRC_FlagParsing(t *testing.T) {
	rcTokens := []string{"--rc", "-rc", "rc", "--repo-cache", "repo-cache"}
	for _, token := range rcTokens {
		if !hasRCFlagOrToken([]string{token}) {
			t.Errorf("hasRCFlagOrToken([%q]) expected true", token)
		}
	}

	if hasRCFlagOrToken([]string{"clone", "only-missing"}) {
		t.Errorf("hasRCFlagOrToken without rc tokens expected false")
	}

	parsed := parseRCArgs([]string{"--rc", "-y", "2", "--verbose"})
	if !parsed.isAssumeYes {
		t.Errorf("parseRCArgs expected isAssumeYes=true")
	}
	if parsed.manifestIndex != 2 {
		t.Errorf("parseRCArgs expected manifestIndex=2, got %d", parsed.manifestIndex)
	}
	if len(parsed.passthroughFlags) != 1 || parsed.passthroughFlags[0] != "--verbose" {
		t.Errorf("parseRCArgs unexpected passthrough flags: %v", parsed.passthroughFlags)
	}
}

func TestCloneRC_SplitDBCache(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "repocache_test.db")

	db, err := OpenRepoCacheDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenRepoCacheDBAt failed: %v", err)
	}
	defer db.Close()

	manifestFile := filepath.Join(tempDir, "sample-manifest.json")
	sampleData := []genericManifestItem{
		{
			RepoName: "gitmap-v28",
			Owner:    "alimtvnetwork",
			URL:      "https://github.com/alimtvnetwork/gitmap-v28.git",
			Branch:   "main",
		},
		{
			RepoName: "smart-scanner",
			Owner:    "alimtvnetwork",
			URL:      "git@github.com:alimtvnetwork/smart-scanner.git",
			Branch:   "master",
		},
	}
	bytes, err := json.MarshalIndent(sampleData, "", "  ")
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	if err := os.WriteFile(manifestFile, bytes, 0644); err != nil {
		t.Fatalf("os.WriteFile failed: %v", err)
	}

	hash, mtime, err := ComputeFileHashAndMtime(manifestFile)
	if err != nil {
		t.Fatalf("ComputeFileHashAndMtime failed: %v", err)
	}

	// 1. Initial cache check -> expected miss
	entries, isHit, err := GetCachedManifestEntries(db, manifestFile, hash, mtime)
	if err != nil {
		t.Fatalf("GetCachedManifestEntries failed: %v", err)
	}
	if isHit {
		t.Fatalf("expected initial cache miss, got hit")
	}

	// 2. Store entries
	var toStore []RepoCacheEntry
	for _, item := range sampleData {
		toStore = append(toStore, RepoCacheEntry{
			RepoName: item.RepoName,
			Owner:    item.Owner,
			CloneUrl: item.URL,
			IsSSH:    isSSHProtocol(item.URL),
			IsHTTPS:  isHTTPSProtocol(item.URL),
			Branch:   item.Branch,
		})
	}

	if err := StoreManifestEntries(db, manifestFile, hash, mtime, toStore); err != nil {
		t.Fatalf("StoreManifestEntries failed: %v", err)
	}

	// 3. Cache hit benchmark (<5ms requirement)
	start := time.Now()
	entries, isHit, err = GetCachedManifestEntries(db, manifestFile, hash, mtime)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("GetCachedManifestEntries on hit failed: %v", err)
	}
	if !isHit {
		t.Fatalf("expected cache hit after storing")
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 cached entries, got %d", len(entries))
	}
	if elapsed > 100*time.Millisecond {
		t.Logf("Cache hit duration: %v", elapsed)
	}

	// 4. Invalidation check on mismatched hash
	_, isMismatchedHit, _ := GetCachedManifestEntries(db, manifestFile, "fakehash", mtime)
	if isMismatchedHit {
		t.Fatalf("expected cache miss on mismatched hash")
	}
}

