package clonefrom

import (
	"testing"
)

func TestDetectLFSSmudgeError(t *testing.T) {
	output := `Error downloading object: assets/01-licensing.xmind (1929d9c): Smudge error: Error downloading
assets/01-licensing.xmind (1929d9c77f0c1c44704e0a2d4809f93a5bca60b48397cd6a3f5d3b3ee951bb56):
[1929d9c77f0c1c44704e0a2d4809f93a5bca60b48397cd6a3f5d3b3ee951bb56] Object does not exist on the server: [404] Object does not exist on the server

Errors logged to 'D:\work\lara-licensing\.git\lfs\logs\20260810T041448.594898.log'.

Use git lfs logs last to view the log.

error: external filter 'git-lfs filter-process' failed

fatal: assets/01-licensing.xmind: smudge filter lfs failed

warning: Clone succeeded, but checkout failed.

You can inspect what was checked out with 'git status'

and retry with 'git restore --source=HEAD :/'
`
	file, ok := detectLFSSmudgeError(output)
	if !ok {
		t.Fatalf("expected to detect LFS smudge error, got ok=false")
	}
	if file != "assets/01-licensing.xmind" {
		t.Errorf("expected file 'assets/01-licensing.xmind', got '%s'", file)
	}

	// Test fallback regex
	output2 := `Error downloading object: deep/path/file.bin (abc1234): Smudge error: [404] Object does not exist on the server`
	file2, ok2 := detectLFSSmudgeError(output2)
	if !ok2 {
		t.Fatalf("expected to detect LFS smudge error using fallback regex, got ok=false")
	}
	if file2 != "deep/path/file.bin" {
		t.Errorf("expected file 'deep/path/file.bin', got '%s'", file2)
	}

	// Test negative
	output3 := `fatal: repository 'https://github.com/missing/repo.git' not found`
	_, ok3 := detectLFSSmudgeError(output3)
	if ok3 {
		t.Fatalf("expected NOT to detect LFS smudge error for standard missing repo")
	}
}
