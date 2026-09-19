package cmd

import (
	"testing"
)

func TestIsSyncHelp(t *testing.T) {
	if !isSyncHelp("help") {
		t.Errorf("expected isSyncHelp(\"help\") to be true")
	}

	if !isSyncHelp("-h") {
		t.Errorf("expected isSyncHelp(\"-h\") to be true")
	}

	if !isSyncHelp("--help") {
		t.Errorf("expected isSyncHelp(\"--help\") to be true")
	}

	if isSyncHelp("ignore") {
		t.Errorf("expected isSyncHelp(\"ignore\") to be false")
	}
}
