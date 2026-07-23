package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestCommitTransferIncludeMergesDefault asserts the v6.0.0 breaking change:
// IncludeMerges defaults to true when no flag is passed.
func TestCommitTransferIncludeMergesDefault(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		constants.CmdCommitRight,
		constants.CmdCommitLeft,
		constants.CmdCommitBoth,
	} {
		spec, ok := commitTransferSpecFor(name)
		if !ok {
			t.Fatalf("unknown spec for %s", name)
		}

		opts, positional := parseCommitTransferArgs(spec, []string{"/tmp/left", "/tmp/right"})
		if len(positional) != 2 {
			t.Fatalf("%s: expected 2 positional args, got %d", name, len(positional))
		}
		if !opts.IncludeMerges {
			t.Errorf("%s: IncludeMerges = false, want true (v6.0.0 default)", name)
		}
	}
}

// TestCommitTransferIncludeMergesExplicit asserts that --include-merges and
// --no-include-merges override the default correctly.
func TestCommitTransferIncludeMergesExplicit(t *testing.T) {
	t.Parallel()

	spec, _ := commitTransferSpecFor(constants.CmdCommitRight)

	// --include-merges (redundant in v6, but must still work)
	optsOn, _ := parseCommitTransferArgs(spec, []string{
		"--include-merges", "/tmp/left", "/tmp/right",
	})
	if !optsOn.IncludeMerges {
		t.Error("--include-merges: IncludeMerges = false, want true")
	}

	// --no-include-merges (explicit opt-out)
	optsOff, _ := parseCommitTransferArgs(spec, []string{
		"--no-include-merges", "/tmp/left", "/tmp/right",
	})
	if optsOff.IncludeMerges {
		t.Error("--no-include-merges: IncludeMerges = true, want false")
	}
}
