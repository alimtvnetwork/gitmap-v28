package cmdpull

import (
	"strings"
	"testing"
)

func TestFormatStepBadgeSafe(t *testing.T) {
	badge := FormatStepBadge(PullStepTypeUpToDate, true)
	isExpected := badge == "[=]"
	if !isExpected {
		t.Fatalf("expected [=], got %q", badge)
	}

	errBadge := FormatStepBadge(PullStepTypeError, true)
	isErrExpected := errBadge == "[X]"
	if !isErrExpected {
		t.Fatalf("expected [X], got %q", errBadge)
	}
}

func TestFormatStepBadgeRich(t *testing.T) {
	badge := FormatStepBadge(PullStepTypeUpToDate, false)
	isExpected := badge == "✔"
	if !isExpected {
		t.Fatalf("expected ✔, got %q", badge)
	}

	ffBadge := FormatStepBadge(PullStepTypeFastForward, false)
	isFFExpected := ffBadge == "⚡"
	if !isFFExpected {
		t.Fatalf("expected ⚡, got %q", ffBadge)
	}
}

func TestParseGitPullOutputUpToDate(t *testing.T) {
	output := "Already up to date.\n"
	step, rng, changes := ParseGitPullOutput(output, "abc1234", "abc1234")
	isUpToDate := step == PullStepTypeUpToDate
	if !isUpToDate {
		t.Fatalf("expected UP_TO_DATE, got %s", step)
	}

	isRangeExpected := rng == "abc1234"
	if !isRangeExpected {
		t.Fatalf("expected abc1234, got %q", rng)
	}

	isChangesExpected := changes == "up-to-date"
	if !isChangesExpected {
		t.Fatalf("expected up-to-date, got %q", changes)
	}
}

func TestParseGitPullOutputFastForward(t *testing.T) {
	output := "Updating abc1234..def5678\nFast-forward\n 2 files changed, 10 insertions(+), 2 deletions(-)\n"
	step, rng, changes := ParseGitPullOutput(output, "abc1234567", "def5678901")
	isFF := step == PullStepTypeFastForward
	if !isFF {
		t.Fatalf("expected FAST_FORWARD, got %s", step)
	}

	isRangeExpected := rng == "abc1234..def5678"
	if !isRangeExpected {
		t.Fatalf("expected abc1234..def5678, got %q", rng)
	}

	hasPlus := strings.Contains(changes, "+10")
	if !hasPlus {
		t.Fatalf("expected +10 in changes, got %q", changes)
	}
}

func TestParseGitPullOutputConflict(t *testing.T) {
	output := "CONFLICT (content): Merge conflict in file.go\nAutomatic merge failed\n"
	step, rng, changes := ParseGitPullOutput(output, "abc1234", "def5678")
	isConflict := step == PullStepTypeConflict
	if !isConflict {
		t.Fatalf("expected CONFLICT, got %s", step)
	}

	isRangeEmpty := rng == ""
	if !isRangeEmpty {
		t.Fatalf("expected empty range on conflict, got %q", rng)
	}

	isConflictWord := changes == "conflict"
	if !isConflictWord {
		t.Fatalf("expected conflict changes, got %q", changes)
	}
}

func TestParseGitPullOutputError(t *testing.T) {
	output := "fatal: refusing to merge unrelated histories\n"
	step, _, changes := ParseGitPullOutput(output, "abc1234", "def5678")
	isErr := step == PullStepTypeError
	if !isErr {
		t.Fatalf("expected ERROR, got %s", step)
	}

	isErrWord := changes == "error"
	if !isErrWord {
		t.Fatalf("expected error changes, got %q", changes)
	}
}
