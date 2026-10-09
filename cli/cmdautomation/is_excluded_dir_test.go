package cmdautomation

import (
	"testing"
)

func TestIsExcludedDir_Git(t *testing.T) {
	isNonExcludedDir := !isExcludedDir(".git")
	if isNonExcludedDir {
		t.Error("expected .git to be excluded")
	}
}

func TestIsExcludedDir_Vendor(t *testing.T) {
	isNonExcludedDir := !isExcludedDir("vendor")
	if isNonExcludedDir {
		t.Error("expected vendor to be excluded")
	}
}

func TestIsExcludedDir_NodeModules(t *testing.T) {
	isNonExcludedDir := !isExcludedDir("node_modules")
	if isNonExcludedDir {
		t.Error("expected node_modules to be excluded")
	}
}

func TestIsExcludedDir_Regular(t *testing.T) {
	if isExcludedDir("src") {
		t.Error("expected src to not be excluded")
	}
}
