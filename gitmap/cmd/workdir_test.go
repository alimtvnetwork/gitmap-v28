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

	optsHelp := parseWorkDirFlags([]string{"help"})
	if optsHelp.Action != "help" {
		t.Fatalf("expected action help, got %s", optsHelp.Action)
	}

	optsDef := parseWorkDirFlags([]string{"default"})
	if optsDef.Action != "default" {
		t.Fatalf("expected action default, got %s", optsDef.Action)
	}
}

func TestWorkDirHelpAndUsage(t *testing.T) {
	optsHelp := workDirOptions{Action: "help"}
	err := dispatchWorkDirAction(optsHelp)
	if err != nil {
		t.Fatalf("dispatchWorkDirAction help returned error: %v", err)
	}

	optsDashH := workDirOptions{Action: "-h"}
	errDashH := dispatchWorkDirAction(optsDashH)
	if errDashH != nil {
		t.Fatalf("dispatchWorkDirAction -h returned error: %v", errDashH)
	}
}

func TestCDWorkDirKeyword(t *testing.T) {
	if !isWorkDirKeyword("work") {
		t.Error("expected 'work' to be a workdir keyword")
	}

	if !isWorkDirKeyword("default") {
		t.Error("expected 'default' to be a workdir keyword")
	}

	if !isWorkDirKeyword("workdir") {
		t.Error("expected 'workdir' to be a workdir keyword")
	}

	if isWorkDirKeyword("unknown-repo") {
		t.Error("did not expect 'unknown-repo' to be a workdir keyword")
	}
}
