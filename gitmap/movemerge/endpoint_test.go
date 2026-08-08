package movemerge

import (
	"path/filepath"
	"testing"
)

func TestClassifyEndpoint_FolderPaths(t *testing.T) {
	cases := []string{"./local", "/abs/path", "..\\rel", "plain-folder"}
	for _, raw := range cases {
		kind, _, _, _ := ClassifyEndpoint(raw)
		isFolder := kind == EndpointFolder
		if !isFolder {
			t.Errorf("ClassifyEndpoint(%q) kind = %v, want Folder", raw, kind)
		}
	}
}

func TestClassifyEndpoint_HTTPSWithBranch(t *testing.T) {
	kind, url, branch, _ := ClassifyEndpoint("https://github.com/owner/repo:develop")
	isURL := kind == EndpointURL
	if !isURL {
		t.Fatalf("kind = %v, want URL", kind)
	}
	isMatchURL := url == "https://github.com/owner/repo"
	if !isMatchURL {
		t.Errorf("url = %q", url)
	}
	isMatchBranch := branch == "develop"
	if !isMatchBranch {
		t.Errorf("branch = %q", branch)
	}
}

func TestClassifyEndpoint_HTTPSNoBranch(t *testing.T) {
	_, url, branch, _ := ClassifyEndpoint("https://github.com/owner/repo.git")
	isMatch := url == "https://github.com/owner/repo.git" && branch == ""
	if !isMatch {
		t.Errorf("got url=%q branch=%q", url, branch)
	}
}

func TestClassifyEndpoint_SCPGitAtForm(t *testing.T) {
	// git@host:user/repo has a colon but it is not a branch.
	_, url, branch, _ := ClassifyEndpoint("git@github.com:owner/repo.git")
	isMatch := url == "git@github.com:owner/repo.git" && branch == ""
	if !isMatch {
		t.Errorf("scp form: url=%q branch=%q", url, branch)
	}
}

func TestMapURLToFolder(t *testing.T) {
	got := MapURLToFolder("/tmp", "https://github.com/owner/my-repo.git")
	// Compare via filepath.ToSlash so the assertion is stable across
	// Windows (returns "\\tmp\\my-repo") and POSIX hosts. The
	// production behavior is correct on each OS — only the literal
	// separator differs, which is not what this test is checking.
	isMatch := filepath.ToSlash(got) == "/tmp/my-repo"
	if !isMatch {
		t.Errorf("MapURLToFolder = %q", got)
	}
}
