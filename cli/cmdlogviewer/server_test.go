package cmdlogviewer

import (
	"testing"
)

func TestResolveTargetRepoAndPath(t *testing.T) {
	repoRoot, relPath, err := ResolveTargetRepoAndPath([]string{"."})
	if err != nil {
		t.Fatalf("ResolveTargetRepoAndPath failed: %v", err)
	}

	if repoRoot == "" {
		t.Error("expected non-empty repoRoot")
	}

	if relPath == "" {
		t.Error("expected non-empty relPath")
	}

	rootRepo, rootRel, err := ResolveTargetRepoAndPath([]string{repoRoot})
	if err != nil {
		t.Fatalf("ResolveTargetRepoAndPath(repoRoot) failed: %v", err)
	}

	if rootRel != "." {
		t.Errorf("expected rootRel '.', got %q", rootRel)
	}
	if rootRepo != repoRoot {
		t.Errorf("expected rootRepo %q, got %q", repoRoot, rootRepo)
	}
}

func TestComputeRelativePath(t *testing.T) {
	tests := []struct {
		repoRoot string
		absPath  string
		expected string
	}{
		{"/repo", "/repo", "."},
		{"/repo", "/repo/sub/file.go", "sub/file.go"},
	}

	for _, tt := range tests {
		got := computeRelativePath(tt.repoRoot, tt.absPath)
		if got != tt.expected {
			t.Errorf("computeRelativePath(%q, %q) = %q, want %q", tt.repoRoot, tt.absPath, got, tt.expected)
		}
	}
}

func TestBindFreePort(t *testing.T) {
	ln, port, err := bindFreePort(9876)
	if err != nil {
		t.Fatalf("bindFreePort failed: %v", err)
	}
	defer ln.Close()

	if port < 9876 {
		t.Errorf("expected port >= 9876, got %d", port)
	}
}

func TestFormatTargetURL(t *testing.T) {
	urlNormal := formatTargetURL(8090, "foo/bar.go", false)
	if urlNormal != "http://127.0.0.1:8090/?path=foo%2Fbar.go" {
		t.Errorf("unexpected normal url: %s", urlNormal)
	}

	urlCompare := formatTargetURL(8090, "foo/bar.go", true)
	if urlCompare != "http://127.0.0.1:8090/?path=foo%2Fbar.go&compare=1" {
		t.Errorf("unexpected compare url: %s", urlCompare)
	}
}
