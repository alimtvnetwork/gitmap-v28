package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func TestIsClusterHelpToken(t *testing.T) {
	if !isClusterHelpToken("--help") {
		t.Error("expected --help to be true")
	}

	if !isClusterHelpToken("-h") {
		t.Error("expected -h to be true")
	}

	if !isClusterHelpToken("help") {
		t.Error("expected help to be true")
	}

	if isClusterHelpToken("other") {
		t.Error("expected other to be false")
	}
}

func TestRunCluster_Empty(t *testing.T) {
	err := runCluster([]string{})
	if err != nil {
		t.Fatalf("expected nil on empty args, got: %v", err)
	}
}

func TestRunCluster_HelpToken(t *testing.T) {
	tokens := []string{"help", "--help", "-h"}
	for _, tok := range tokens {
		err := runCluster([]string{tok})
		if err != nil {
			t.Fatalf("expected nil on %q, got: %v", tok, err)
		}
	}
}

func TestRunCluster_Unknown(t *testing.T) {
	err := runCluster([]string{"unknown-cluster-xyz"})
	if err == nil {
		t.Fatal("expected error on unknown cluster subcommand, got nil")
	}
}

func TestDispatchClusterSubcommand_Matches(t *testing.T) {
	prev := cliexit.SetExitFunc(func(code int) {})
	defer cliexit.SetExitFunc(prev)

	commands := []string{"add", "join", "ping", "nodes", "ls", "remove", "rm"}
	for _, cmd := range commands {
		res := dispatchClusterSubcommand(cmd, []string{})
		if !res.Data {
			t.Errorf("expected subcommand %q to be matched", cmd)
		}
	}
}

func TestDispatchClusterSubcommand_Unknown(t *testing.T) {
	res := dispatchClusterSubcommand("unknown-cmd-xyz", []string{})
	if res.Data {
		t.Error("expected unknown subcommand to not be matched")
	}
}

func TestRunCluster_InvertedHelp(t *testing.T) {
	prev := cliexit.SetExitFunc(func(code int) {})
	defer cliexit.SetExitFunc(prev)

	err := runCluster([]string{"help", "unknown-xyz"})
	if err == nil {
		t.Fatal("expected error for unknown inverted help, got nil")
	}

	_ = runCluster([]string{"help", "nodes"})
}

func TestDispatchInvertedClusterHelp_Branches(t *testing.T) {
	prev := cliexit.SetExitFunc(func(code int) {})
	defer cliexit.SetExitFunc(prev)

	resUnmatched := dispatchInvertedClusterHelp([]string{"nodes"})
	if resUnmatched.Data {
		t.Error("expected non-inverted help to be unmatched")
	}

	resMatched := dispatchInvertedClusterHelp([]string{"help", "nodes"})
	if !resMatched.Data {
		t.Error("expected help nodes to be matched")
	}

	resUnknown := dispatchInvertedClusterHelp([]string{"help", "nonexistent"})
	if !resUnknown.Data || resUnknown.AppError() == nil {
		t.Error("expected unknown command under help to be matched with AppError")
	}
}
