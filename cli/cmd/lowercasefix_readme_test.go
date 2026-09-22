package cmd

import "testing"

func TestIsRootReadme(t *testing.T) {
	if !isRootReadme("README.md", "README.md") {
		t.Errorf("expected README.md in root to be root readme")
	}
	if !isRootReadme("Readme.md", "Readme.md") {
		t.Errorf("expected Readme.md in root to be root readme")
	}
	if !isRootReadme("README", "README") {
		t.Errorf("expected README in root to be root readme")
	}
	if isRootReadme("docs/README.md", "README.md") {
		t.Errorf("expected docs/README.md NOT to be root readme")
	}
	if isRootReadme("sub/dir/README.md", "README.md") {
		t.Errorf("expected sub/dir/README.md NOT to be root readme")
	}
}

func TestIsReadmeAlias(t *testing.T) {
	aliases := []string{"lowercase-readme", "lower-case-readme", "readme-lower", "readme-lowercase", "lcr"}
	for _, a := range aliases {
		if !isReadmeAlias(a) {
			t.Errorf("expected %q to be recognized as readme alias", a)
		}
	}
	if isReadmeAlias("lowercase") {
		t.Errorf("expected lowercase not to be readme alias")
	}
}

func TestResolveLcfPatterns_Readme(t *testing.T) {
	patterns, isReadme := resolveLcfPatterns([]string{"readme"})
	if !isReadme {
		t.Errorf("expected isReadme=true for 'readme'")
	}
	if len(patterns) != 1 || patterns[0] != "readme*" {
		t.Errorf("expected patterns=['readme*'], got %v", patterns)
	}
}
