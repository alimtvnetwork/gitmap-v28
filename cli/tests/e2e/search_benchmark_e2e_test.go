//go:build e2e

package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/searcher"
)

func TestSearchBenchmark_AUMVsPythonVsFind(t *testing.T) {
	rootPath, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("failed to resolve repository root: %v", err)
	}

	searchPattern := "SSHConnection"

	// 1. GitMap Native Searcher Benchmark
	targetDir := filepath.Join(rootPath, "cli", "cmdssh")
	files, _ := filepath.Glob(filepath.Join(targetDir, "*.go"))
	startAUM := time.Now()
	var aumMatches int
	for _, f := range files {
		content, rErr := os.ReadFile(f)
		if rErr == nil {
			res := searcher.SearchExact(string(content), searchPattern, f, filepath.Base(f))
			aumMatches += len(res)
		}
	}
	durationAUM := time.Since(startAUM)

	// 2. Python Search Benchmark (03-ai-scripts/12-fast-cached-grep.py)
	pythonScript := filepath.Join(rootPath, "03-ai-scripts", "12-fast-cached-grep.py")
	startPython := time.Now()
	cmd := exec.Command("python", pythonScript, "--pattern", searchPattern, "--limit", "50")
	cmd.Dir = rootPath
	pyOut, pyErr := cmd.Output()
	durationPython := time.Since(startPython)
	if pyErr != nil {
		t.Logf("Notice: python fast-grep output: %s, err: %v", string(pyOut), pyErr)
	}

	t.Logf("----------------------------------------------------------------------")
	t.Logf("Search Engine Benchmark Results (Target: %q)", searchPattern)
	t.Logf("  • GitMap Native AUM Searcher : %v (Matches: %d)", durationAUM, aumMatches)
	t.Logf("  • Python Fast Cached Grep    : %v", durationPython)
	t.Logf("----------------------------------------------------------------------")

	if durationAUM > 5*time.Second {
		t.Errorf("AUM search exceeded 5s threshold: %v", durationAUM)
	}
}
