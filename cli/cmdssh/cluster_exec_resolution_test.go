package cmdssh

import (
	"testing"
)

func TestParseClusterExecArgs_GitmapMultiToken(t *testing.T) {
	args := []string{"status --json"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.target != "all" || opts.command != "gitmap status --json" {
		t.Fatalf("unexpected result: %+v, err: %v", opts, err)
	}
}

func TestParseClusterExecArgs_ExplicitTargetGitmapCompound(t *testing.T) {
	args := []string{"devbox", "gitmap status && gitmap pipeline"}
	opts, err := parseClusterExecArgs(args)
	expectedCmd := "gitmap status && gitmap pipeline"
	if err != nil || opts.target != "devbox" || opts.command != expectedCmd {
		t.Fatalf("unexpected result: %+v, err: %v", opts, err)
	}
}

func TestParseClusterExecArgs_ExplicitTargetGitmapSubcommand(t *testing.T) {
	args := []string{"devbox", "status --json"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.target != "devbox" || opts.command != "gitmap status --json" {
		t.Fatalf("unexpected result: %+v, err: %v", opts, err)
	}
}

func TestParseClusterExecArgs_MultiTargetComma(t *testing.T) {
	args := []string{"devbox,worker-1", "uname -a"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.target != "devbox,worker-1" || opts.command != "uname -a" {
		t.Fatalf("unexpected result: %+v, err: %v", opts, err)
	}
}

func TestParseClusterExecArgs_NewGitmapCommands(t *testing.T) {
	args := []string{"prompts-template ls"}
	opts, err := parseClusterExecArgs(args)
	if err != nil || opts.target != "all" || opts.command != "gitmap prompts-template ls" {
		t.Fatalf("unexpected result: %+v, err: %v", opts, err)
	}
}
