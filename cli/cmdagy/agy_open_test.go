package cmdagy

import (
	"os"
	"testing"
)

func TestResolveTargetPath(t *testing.T) {
	cwd, _ := os.Getwd()
	if path := resolveTargetPath("."); path != cwd {
		t.Errorf("expected %s, got %s", cwd, path)
	}
	if path := resolveTargetPath(""); path != cwd {
		t.Errorf("expected %s, got %s", cwd, path)
	}
}

func TestIsAgyOpenPathArg(t *testing.T) {
	if !isAgyOpenPathArg(".") {
		t.Errorf("expected '.' to be recognized as path")
	}
	if !isAgyOpenPathArg("..") {
		t.Errorf("expected '..' to be recognized as path")
	}
	if isAgyOpenPathArg("optimize-projects") {
		t.Errorf("expected 'optimize-projects' not to be path")
	}
}
