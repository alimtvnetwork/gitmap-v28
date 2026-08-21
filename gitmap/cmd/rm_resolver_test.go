package cmd

import (
	"testing"
)

func TestParseRmFlags(t *testing.T) {
	yes, dbOnly, args := parseRmFlags([]string{"-y", "--db-only", ".\\prompt-architect", "macro*"})
	if !yes {
		t.Errorf("parseRmFlags expected yes=true")
	}
	if !dbOnly {
		t.Errorf("parseRmFlags expected dbOnly=true")
	}
	if len(args) != 2 || args[0] != ".\\prompt-architect" || args[1] != "macro*" {
		t.Errorf("parseRmFlags args mismatch: %+v", args)
	}

	targets := expandRmTargets([]string{"foo,bar", "baz"})
	if len(targets) != 3 || targets[0] != "foo" || targets[1] != "bar" || targets[2] != "baz" {
		t.Errorf("expandRmTargets mismatch: %+v", targets)
	}
}
