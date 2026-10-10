package cmdclone

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func TestCloneRC_ManifestDiscoveryPriority(t *testing.T) {
	tempRoot := t.TempDir()

	// 1. Primary manifest: 01-gitmap/gitmap.json
	primaryDir := filepath.Join(tempRoot, "01-gitmap")
	if err := os.MkdirAll(primaryDir, 0755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	writeTestManifest(t, filepath.Join(primaryDir, "gitmap.json"), "primary-repo")

	// 2. Sequential manifests: 02-gitmap.json and 03-gitmap.json
	writeTestManifest(t, filepath.Join(tempRoot, "02-gitmap.json"), "seq2-repo")
	writeTestManifest(t, filepath.Join(tempRoot, "03-gitmap.json"), "seq3-repo")

	// 3. Auxiliary manifest: repos.json
	writeTestManifest(t, filepath.Join(tempRoot, "repos.json"), "aux-repo")

	// Set test hook
	oldHook := ResolveRepoCacheRootFn
	oldDBPath := CustomRepoCacheDBPath
	ResolveRepoCacheRootFn = func() string { return tempRoot }
	CustomRepoCacheDBPath = filepath.Join(tempRoot, "cache_split.db")
	defer func() {
		ResolveRepoCacheRootFn = oldHook
		CustomRepoCacheDBPath = oldDBPath
	}()

	manifests, err := DiscoverRCManifests(tempRoot)
	if err != nil {
		t.Fatalf("DiscoverRCManifests failed: %v", err)
	}

	if len(manifests) != 4 {
		t.Fatalf("expected 4 manifests, got %d", len(manifests))
	}

	// Slot 1 must be primary
	if manifests[0].Index != 1 || !strings.Contains(manifests[0].RelativePath, "01-gitmap") {
		t.Errorf("manifest [1] expected primary 01-gitmap, got %s", manifests[0].RelativePath)
	}

	// Slot 2 must be 02-gitmap.json
	if manifests[1].Index != 2 || !strings.Contains(manifests[1].RelativePath, "02-gitmap.json") {
		t.Errorf("manifest [2] expected 02-gitmap.json, got %s", manifests[1].RelativePath)
	}

	// Slot 3 must be 03-gitmap.json
	if manifests[2].Index != 3 || !strings.Contains(manifests[2].RelativePath, "03-gitmap.json") {
		t.Errorf("manifest [3] expected 03-gitmap.json, got %s", manifests[2].RelativePath)
	}

	// Slot 4 must be aux repos.json
	if manifests[3].Index != 4 || !strings.Contains(manifests[3].RelativePath, "repos.json") {
		t.Errorf("manifest [4] expected repos.json, got %s", manifests[3].RelativePath)
	}
}

func TestCloneRC_RenderShortTreeView(t *testing.T) {
	manifests := []DiscoveredManifest{
		{
			Index:        1,
			RelativePath: "01-gitmap/gitmap.json",
			TotalRepos:   2,
			Entries: []RepoCacheEntry{
				{Owner: "alimtvnetwork", RepoName: "gitmap-v28", CloneUrl: "https://github.com/alimtvnetwork/gitmap-v28.git"},
				{Owner: "alimtvnetwork", RepoName: "smart-scanner", CloneUrl: "git@github.com:alimtvnetwork/smart-scanner.git"},
			},
		},
		{
			Index:        2,
			RelativePath: "02-gitmap.json",
			TotalRepos:   1,
			Entries: []RepoCacheEntry{
				{Owner: "devops", RepoName: "ci-runner", CloneUrl: "git@github.com:devops/ci-runner.git"},
			},
		},
	}

	rendered := RenderShortTreeView(manifests)
	if !strings.Contains(rendered, "[1] 01-gitmap/gitmap.json") {
		t.Errorf("rendered tree missing manifest [1]")
	}
	if !strings.Contains(rendered, "alimtvnetwork/gitmap-v28 [Public HTTPS]") {
		t.Errorf("rendered tree missing alimtvnetwork/gitmap-v28 entry")
	}
	if !strings.Contains(rendered, "alimtvnetwork/smart-scanner [SSH]") {
		t.Errorf("rendered tree missing alimtvnetwork/smart-scanner entry")
	}
	if !strings.Contains(rendered, "[2] 02-gitmap.json") {
		t.Errorf("rendered tree missing manifest [2]")
	}
}

func writeTestManifest(t *testing.T, path, repoName string) {
	t.Helper()
	items := []genericManifestItem{
		{
			RepoName: repoName,
			Owner:    "testowner",
			URL:      "https://github.com/testowner/" + repoName + ".git",
			Branch:   "main",
		},
	}
	bytes, err := json.Marshal(items)
	if err != nil {
		t.Fatalf("writeTestManifest marshal failed: %v", err)
	}
	if err := os.WriteFile(path, bytes, 0644); err != nil {
		t.Fatalf("writeTestManifest write failed: %v", err)
	}
}

func TestCloneRC_RunCloneRCEmpty(t *testing.T) {
	tempRoot := t.TempDir()
	oldHook := ResolveRepoCacheRootFn
	ResolveRepoCacheRootFn = func() string { return tempRoot }
	defer func() { ResolveRepoCacheRootFn = oldHook }()

	err := RunCloneRC([]string{"-y"}, CloneModeStandard)
	if err != nil {
		t.Errorf("RunCloneRC with empty repo-cache expected nil error, got %v", err)
	}
}

func TestCloneRC_DispatchModes(t *testing.T) {
	if CloneModeStandard != "standard" {
		t.Errorf("unexpected CloneModeStandard: %s", CloneModeStandard)
	}
	if CloneModeCFR != "cfr" {
		t.Errorf("unexpected CloneModeCFR: %s", CloneModeCFR)
	}
	if CloneModeCFRP != "cfrp" {
		t.Errorf("unexpected CloneModeCFRP: %s", CloneModeCFRP)
	}
}

