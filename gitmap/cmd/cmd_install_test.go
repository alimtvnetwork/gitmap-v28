package cmd

import (
	"reflect"
	"testing"
)

func TestParseCGFlags(t *testing.T) {
	opts := parseCGFlags([]string{"install", "--all"})
	if !opts.All || opts.Action != "install" {
		t.Errorf("Expected All=true, Action=install, got %+v", opts)
	}

	opts2 := parseCGFlags([]string{"--exclude", "repo3", "update", "repo1", "repo2"})
	if opts2.All || opts2.Action != "update" || opts2.Exclude != "repo3" || !reflect.DeepEqual(opts2.Repos, []string{"repo1", "repo2"}) {
		t.Errorf("Unexpected output: %+v", opts2)
	}
}

func TestParseSEFlags(t *testing.T) {
	opts := parseSEFlags([]string{"--exclude", "m1,m2", "ps", "echo test"})
	if opts.Exclude != "m1,m2" || len(opts.Args) != 2 || opts.Args[0] != "ps" {
		t.Errorf("Unexpected output: %+v", opts)
	}
}

func TestParseSJFlags(t *testing.T) {
	opts := parseSJFlags([]string{"--list"})
	if !opts.List {
		t.Errorf("Expected List=true, got %+v", opts)
	}

	opts2 := parseSJFlags([]string{"ls"})
	if !opts2.List {
		t.Errorf("Expected List=true for ls, got %+v", opts2)
	}
}
