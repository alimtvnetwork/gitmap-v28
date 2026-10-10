package cmdscan

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestScanRC_NormalizeRepoURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "https standard with .git",
			input:    "https://github.com/alimtvnetwork/gitmap-v28.git",
			expected: "github.com/alimtvnetwork/gitmap-v28",
		},
		{
			name:     "https standard with trailing slash",
			input:    "https://github.com/alimtvnetwork/gitmap-v28/",
			expected: "github.com/alimtvnetwork/gitmap-v28",
		},
		{
			name:     "ssh scp syntax with .git",
			input:    "git@github.com:alimtvnetwork/gitmap-v28.git",
			expected: "github.com/alimtvnetwork/gitmap-v28",
		},
		{
			name:     "ssh url scheme with trailing slash and .git",
			input:    "ssh://git@github.com/alimtvnetwork/gitmap-v28.git/",
			expected: "github.com/alimtvnetwork/gitmap-v28",
		},
		{
			name:     "uppercase and whitespace trimming",
			input:    "   HTTPS://GITHUB.COM/AlimTVNetwork/GitMap-v28.GIT/   ",
			expected: "github.com/alimtvnetwork/gitmap-v28",
		},
		{
			name:     "nested gitlab group ssh syntax",
			input:    "git@gitlab.com:enterprise/subgroup/platform.git",
			expected: "gitlab.com/enterprise/subgroup/platform",
		},
		{
			name:     "nested gitlab group https syntax",
			input:    "https://gitlab.com/enterprise/subgroup/platform.git",
			expected: "gitlab.com/enterprise/subgroup/platform",
		},
		{
			name:     "empty string",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NormalizeRepoURL(tt.input)
			if actual != tt.expected {
				t.Errorf("NormalizeRepoURL(%q) = %q, want %q", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestScanRC_FlagParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantRC       bool
		wantSeparate bool
	}{
		{
			name:         "no rc flags",
			args:         []string{"."},
			wantRC:       false,
			wantSeparate: false,
		},
		{
			name:         "rc canonical flag",
			args:         []string{".", "--rc"},
			wantRC:       true,
			wantSeparate: false,
		},
		{
			name:         "repo-cache alias flag",
			args:         []string{".", "--repo-cache"},
			wantRC:       true,
			wantSeparate: false,
		},
		{
			name:         "rc and separate canonical flag",
			args:         []string{".", "--rc", "--separate"},
			wantRC:       true,
			wantSeparate: true,
		},
		{
			name:         "rc and seprate typo alias",
			args:         []string{".", "--rc", "--seprate"},
			wantRC:       true,
			wantSeparate: true,
		},
		{
			name:         "rc and sep short alias",
			args:         []string{".", "--rc", "--sep"},
			wantRC:       true,
			wantSeparate: true,
		},
		{
			name:         "rc and s single letter alias",
			args:         []string{".", "--rc", "-s"},
			wantRC:       true,
			wantSeparate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, rc, separate, _, _, _ := ParseScanFlags(tt.args)
			if rc != tt.wantRC {
				t.Errorf("ParseScanFlags(%v) rc = %v, want %v", tt.args, rc, tt.wantRC)
			}
			if separate != tt.wantSeparate {
				t.Errorf("ParseScanFlags(%v) separate = %v, want %v", tt.args, separate, tt.wantSeparate)
			}
		})
	}
}

func TestScanRC_AutoMergeAndRefresh(t *testing.T) {
	tmpDir := t.TempDir()

	batch1 := []model.ScanRecord{
		{
			RepoName:     "repo-a",
			HTTPSUrl:     "https://github.com/example/repo-a.git",
			Branch:       "master",
			AbsolutePath: "/work/repo-a",
		},
		{
			RepoName:     "repo-b",
			HTTPSUrl:     "https://github.com/example/repo-b.git",
			Branch:       "main",
			AbsolutePath: "/work/repo-b",
		},
	}

	err := MergeScanRecordsWithExistingManifest(tmpDir, batch1, true)
	if err != nil {
		t.Fatalf("first MergeScanRecordsWithExistingManifest failed: %v", err)
	}

	manifestPath := filepath.Join(tmpDir, "01-gitmap", "gitmap.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest after batch 1 failed: %v", err)
	}

	var recs1 []exportRecord
	if err := json.Unmarshal(data, &recs1); err != nil {
		t.Fatalf("unmarshal batch 1 failed: %v", err)
	}
	if len(recs1) != 2 {
		t.Fatalf("batch 1 count = %d, want 2", len(recs1))
	}

	// Batch 2: Repo A with SSH format and updated branch + new Repo C
	batch2 := []model.ScanRecord{
		{
			RepoName:     "repo-a",
			SSHUrl:       "git@github.com:example/repo-a.git",
			Branch:       "feature/v2",
			AbsolutePath: "/work/repo-a-refreshed",
		},
		{
			RepoName:     "repo-c",
			HTTPSUrl:     "https://github.com/example/repo-c.git",
			Branch:       "develop",
			AbsolutePath: "/work/repo-c",
		},
	}

	err = MergeScanRecordsWithExistingManifest(tmpDir, batch2, true)
	if err != nil {
		t.Fatalf("second MergeScanRecordsWithExistingManifest failed: %v", err)
	}

	data2, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest after batch 2 failed: %v", err)
	}

	var recs2 []exportRecord
	if err := json.Unmarshal(data2, &recs2); err != nil {
		t.Fatalf("unmarshal batch 2 failed: %v", err)
	}

	if len(recs2) != 3 {
		t.Fatalf("merged count = %d, want 3", len(recs2))
	}

	// Check that Repo A metadata was updated
	foundRepoA := false
	for _, r := range recs2 {
		if NormalizeRepoURL(r.URL) == "github.com/example/repo-a" {
			foundRepoA = true
			if r.Branch != "feature/v2" {
				t.Errorf("repo-a branch = %q, want feature/v2", r.Branch)
			}
			if r.AbsolutePath != "/work/repo-a-refreshed" {
				t.Errorf("repo-a absolutePath = %q, want /work/repo-a-refreshed", r.AbsolutePath)
			}
		}
	}
	if !foundRepoA {
		t.Errorf("repo-a not found in merged manifest")
	}
}

func TestScanRC_SequentialAllocation(t *testing.T) {
	tmpDir := t.TempDir()

	// Seed 01-gitmap folder with gitmap.json
	primaryDir := filepath.Join(tmpDir, "01-gitmap")
	if err := os.MkdirAll(primaryDir, 0755); err != nil {
		t.Fatalf("mkdir primaryDir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(primaryDir, "gitmap.json"), []byte("[]"), 0644); err != nil {
		t.Fatalf("write primary gitmap.json: %v", err)
	}

	records := []model.ScanRecord{
		{
			RepoName: "repo-isolated",
			HTTPSUrl: "https://github.com/example/repo-isolated.git",
			Branch:   "main",
		},
	}

	// First separate allocation should allocate 02-gitmap.json
	err := AllocateNextSequentialFile(tmpDir, records, true)
	if err != nil {
		t.Fatalf("AllocateNextSequentialFile #1 failed: %v", err)
	}
	file02 := filepath.Join(tmpDir, "02-gitmap.json")
	if _, err := os.Stat(file02); err != nil {
		t.Fatalf("expected %s to exist: %v", file02, err)
	}

	// Second separate allocation should allocate 03-gitmap.json
	err = AllocateNextSequentialFile(tmpDir, records, true)
	if err != nil {
		t.Fatalf("AllocateNextSequentialFile #2 failed: %v", err)
	}
	file03 := filepath.Join(tmpDir, "03-gitmap.json")
	if _, err := os.Stat(file03); err != nil {
		t.Fatalf("expected %s to exist: %v", file03, err)
	}
}

func TestScanRC_ExportScanToRepoCache(t *testing.T) {
	tmpDir := t.TempDir()

	ResolveRepoCacheRootFn = func() string {
		return tmpDir
	}
	defer func() {
		ResolveRepoCacheRootFn = nil
	}()

	records := []model.ScanRecord{
		{
			RepoName: "repo-demo",
			HTTPSUrl: "https://github.com/example/repo-demo.git",
			Branch:   "main",
		},
	}

	// isSeparate = false -> writes to 01-gitmap/gitmap.json
	if err := ExportScanToRepoCache(records, false, true); err != nil {
		t.Fatalf("ExportScanToRepoCache master merge failed: %v", err)
	}
	mainPath := filepath.Join(tmpDir, "01-gitmap", "gitmap.json")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("expected %s to exist: %v", mainPath, err)
	}

	// isSeparate = true -> writes to 02-gitmap.json
	if err := ExportScanToRepoCache(records, true, true); err != nil {
		t.Fatalf("ExportScanToRepoCache sequential allocation failed: %v", err)
	}
	seqPath := filepath.Join(tmpDir, "02-gitmap.json")
	if _, err := os.Stat(seqPath); err != nil {
		t.Fatalf("expected %s to exist: %v", seqPath, err)
	}
}
