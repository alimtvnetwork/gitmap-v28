package cmdfix

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func TestIsNeedsMergeError(t *testing.T) {
	cases := []struct {
		output string
		want   bool
	}{
		{"error: could not write index\nspec/23-app-db/00-overview.md: needs merge", true},
		{"fatal: git-write-tree: error building trees", true},
		{"fatal: Pulling is not possible because you have unmerged files.", true},
		{"You have unmerged paths.", true},
		{"working tree clean", false},
	}

	for _, tc := range cases {
		got := isNeedsMergeError(tc.output)
		if got != tc.want {
			t.Errorf("isNeedsMergeError(%q) = %v, want %v", tc.output, got, tc.want)
		}
	}
}

func TestIsConflictError(t *testing.T) {
	cases := []struct {
		output string
		want   bool
	}{
		{"CONFLICT (content): Merge conflict in file.go", true},
		{"Automatic merge failed; fix conflicts and then commit the result.", true},
		{"failed to merge in the changes.", true},
		{"error: patch failed: file.go:10", false},
	}

	for _, tc := range cases {
		got := isConflictError(tc.output)
		if got != tc.want {
			t.Errorf("isConflictError(%q) = %v, want %v", tc.output, got, tc.want)
		}
	}
}

func TestStripConflictMarkers(t *testing.T) {
	raw := `| # | File | Purpose |
|---|------|---------|
<<<<<<< Updated upstream
| 00 | 00-overview.md | Overview |
=======
| 00 | 00-overview.md | Overview |
>>>>>>> Stashed changes
| 01 | 01-erd.md | ERD |`

	resolved := stripConflictMarkers(raw)
	if strings.Contains(resolved, "<<<<<<<") || strings.Contains(resolved, ">>>>>>>") {
		t.Errorf("expected conflict markers to be removed, got:\n%s", resolved)
	}

	occurrences := strings.Count(resolved, "| 00 | 00-overview.md | Overview |")
	if occurrences != 1 {
		t.Errorf("expected duplicate table row deduplicated to 1, got %d", occurrences)
	}
}

func TestAutoResolveRepoConflicts(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.name", "gitmap").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.email", "gitmap@test.com").Run()

	testFile := filepath.Join(tempDir, "readme.md")
	_ = os.WriteFile(testFile, []byte("# Initial\n"), 0644)
	_ = exec.Command("git", "-C", tempDir, "add", "-A").Run()
	_ = exec.Command("git", "-C", tempDir, "commit", "-m", "initial").Run()

	conflictedContent := `# Initial
<<<<<<< HEAD
Our version
=======
Their version
>>>>>>> branch
`
	_ = os.WriteFile(testFile, []byte(conflictedContent), 0644)

	isResolved := autoResolveRepoConflicts(tempDir)
	if !isResolved {
		t.Fatalf("expected autoResolveRepoConflicts to succeed")
	}

	afterBytes, _ := os.ReadFile(testFile)
	afterStr := string(afterBytes)
	if strings.Contains(afterStr, "<<<<<<<") || strings.Contains(afterStr, ">>>>>>>") {
		t.Errorf("expected conflict markers stripped, got:\n%s", afterStr)
	}
}

func TestRecoverNeedsMerge(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.name", "gitmap").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.email", "gitmap@test.com").Run()

	testFile := filepath.Join(tempDir, "readme.md")
	_ = os.WriteFile(testFile, []byte("# Initial\n"), 0644)
	_ = exec.Command("git", "-C", tempDir, "add", "-A").Run()
	_ = exec.Command("git", "-C", tempDir, "commit", "-m", "initial").Run()

	conflictedContent := `# Initial
<<<<<<< HEAD
Our version
=======
Their version
>>>>>>> branch
`
	_ = os.WriteFile(testFile, []byte(conflictedContent), 0644)

	step := gitutil.RemediationStep{
		Name: "git",
		Args: []string{"-C", tempDir, "stash", "-u"},
	}
	outStr := "error: could not write index\nreadme.md: needs merge\nfatal: git-write-tree: error building trees"

	isRecovered, msg := recoverNeedsMerge(step, outStr)
	if !isRecovered {
		t.Fatalf("expected recoverNeedsMerge to succeed, msg: %s", msg)
	}

	afterBytes, _ := os.ReadFile(testFile)
	if strings.Contains(string(afterBytes), "<<<<<<<") {
		t.Errorf("expected conflict markers resolved, got:\n%s", string(afterBytes))
	}
}

func TestRecoverStashPopConflict(t *testing.T) {
	tempDir := t.TempDir()
	_ = exec.Command("git", "init", tempDir).Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.name", "gitmap").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.email", "gitmap@test.com").Run()

	testFile := filepath.Join(tempDir, "readme.md")
	_ = os.WriteFile(testFile, []byte("# Initial\n"), 0644)
	_ = exec.Command("git", "-C", tempDir, "add", "-A").Run()
	_ = exec.Command("git", "-C", tempDir, "commit", "-m", "initial").Run()

	_ = os.WriteFile(testFile, []byte("# Stashed\n"), 0644)
	_ = exec.Command("git", "-C", tempDir, "stash", "-u").Run()

	step := gitutil.RemediationStep{
		Name: "git",
		Args: []string{"-C", tempDir, "stash", "pop"},
	}
	outStr := "CONFLICT (content): Merge conflict in readme.md\nfailed to merge in the changes"

	isRecovered, _ := recoverStashPopConflict(step, outStr)
	if !isRecovered {
		t.Fatalf("expected recoverStashPopConflict to succeed")
	}

	stashList, _ := exec.Command("git", "-C", tempDir, "stash", "list").CombinedOutput()
	if strings.TrimSpace(string(stashList)) != "" {
		t.Errorf("expected stash to be dropped, got: %s", string(stashList))
	}
}
