package cmd

import (
	"testing"
)

func TestWorkDirCLIParse(t *testing.T) {
	opts := parseWorkDirFlags([]string{"add", "/home/user/work", "--label", "my-work"})
	if opts.Action != "add" || opts.Target != "/home/user/work" || opts.Label != "my-work" {
		t.Fatalf("unexpected parsed options: %+v", opts)
	}

	optsLs := parseWorkDirFlags([]string{})
	if optsLs.Action != "ls" {
		t.Fatalf("expected default action ls, got %s", optsLs.Action)
	}
}
