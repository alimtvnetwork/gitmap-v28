package cmdpull

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func TestFormatRepoGitCmd_ForwardSlash(t *testing.T) {
	windowsPath := `D:\work\gitmap`
	cmd := formatRepoGitCmd(windowsPath, "stash")

	expected := `git -C "D:/work/gitmap" stash`
	if cmd != expected {
		t.Fatalf("expected %q, got %q", expected, cmd)
	}
	if strings.Contains(cmd, `\\`) || strings.Contains(cmd, `\`) {
		t.Fatalf("command contains backslashes: %s", cmd)
	}
}

func TestResolveDirtyTreeDualHints_PullBeforeChanges(t *testing.T) {
	windowsPath := `D:\work\gitmap`
	l1, c1, l2, c2 := resolveDirtyTreeDualHints(windowsPath)

	if l1 != "Commit WIP" || l2 != "Stash Changes" {
		t.Fatalf("unexpected labels: %s, %s", l1, l2)
	}
	if !strings.Contains(c1, "pull --rebase") {
		t.Fatalf("commitCmd should pull before finishing: %s", c1)
	}
	if !strings.Contains(c2, "pull") {
		t.Fatalf("stashCmd should pull before stash pop: %s", c2)
	}
	if strings.Contains(c1, `\`) || strings.Contains(c2, `\`) {
		t.Fatalf("commands contain backslashes: c1=%s, c2=%s", c1, c2)
	}
	expectedStash := `git -C "D:/work/gitmap" stash -u && git -C "D:/work/gitmap" pull && git -C "D:/work/gitmap" stash pop`
	if c2 != expectedStash {
		t.Fatalf("expected stashCmd %q, got %q", expectedStash, c2)
	}
}

func TestRenderItemizedDirtyFiles_ForwardSlash(t *testing.T) {
	diag := gitutil.DirtyDiagnosis{
		ModifiedFiles:  []string{`cli\cmd\test.go`, `pkg\util\helper.go`},
		UntrackedFiles: []string{`new_dir\file.txt`},
	}
	var buf bytes.Buffer
	renderItemizedDirtyFiles(&buf, diag)
	out := buf.String()

	if strings.Contains(out, `\`) {
		t.Fatalf("rendered output contains backslashes: %s", out)
	}
	if !strings.Contains(out, "cli/cmd/test.go") {
		t.Fatalf("expected forward-slash path cli/cmd/test.go, got: %s", out)
	}
	if !strings.Contains(out, "new_dir/file.txt") {
		t.Fatalf("expected forward-slash path new_dir/file.txt, got: %s", out)
	}
}
