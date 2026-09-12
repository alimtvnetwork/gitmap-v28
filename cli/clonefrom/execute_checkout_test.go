package clonefrom

// Tests for the per-row Checkout option:
//
//   - skip   → buildGitArgs appends --no-checkout AND the post-clone
//              hook is a no-op (no working tree materialized).
//   - auto   → legacy behavior, working tree present, no extra git
//              checkout call (covered indirectly by the existing
//              TestExecute_HappyPath which still passes byte-for-byte).
//   - force  → an explicit `git checkout <branch>` runs after clone.
//              Missing-branch case is surfaced as `failed` with a
//              clear MsgCloneFromBranchMissingFmt detail.
//
// Detached-HEAD case: when the row has NO Branch and Checkout=force,
// runPostCloneCheckout returns ("", true) without invoking git —
// asserted by TestPostCloneCheckout_NoBranchIsNoOp. This guards
// against accidentally trying to `git checkout ""` on a default-HEAD
// clone, which would error and break the row spuriously.

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// TestEffectiveCheckout_DefaultsToAuto pins the resolution rule:
// empty Checkout → "auto". A regression here would silently change
// the executor's behavior for every row that omits the field.
func TestEffectiveCheckout_DefaultsToAuto(t *testing.T) {
	if got := EffectiveCheckout(Row{}); got != constants.CloneFromCheckoutAuto {
		t.Fatalf("EffectiveCheckout(empty) = %q, want %q",
			got, constants.CloneFromCheckoutAuto)
	}

	if got := EffectiveCheckout(Row{Checkout: "force"}); got != "force" {
		t.Fatalf("EffectiveCheckout(force) = %q, want force", got)
	}
}

// TestBuildGitArgs_NoCheckoutOnlyForSkipMode asserts the executor
// passes --no-checkout EXACTLY when the resolved mode is "skip" and
// NEVER for auto/force. Critical: the depthflag golden + faithful-
// cmd verifier both depend on the default-row argv staying byte-
// identical to the pre-feature shape.
func TestBuildGitArgs_NoCheckoutOnlyForSkipMode(t *testing.T) {
	cases := []struct {
		name     string
		row      Row
		wantFlag bool
	}{
		{"empty=auto", Row{URL: "https://x/y.git"}, false},
		{"explicit auto", Row{URL: "https://x/y.git",
			Checkout: constants.CloneFromCheckoutAuto}, false},
		{"force", Row{URL: "https://x/y.git",
			Checkout: constants.CloneFromCheckoutForce}, false},
		{"skip", Row{URL: "https://x/y.git",
			Checkout: constants.CloneFromCheckoutSkip}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := buildGitArgs(tc.row, "out")
			has := containsTok(args, constants.CloneFromNoCheckoutFlag)
			if has != tc.wantFlag {
				t.Fatalf("--no-checkout present=%v, want %v\n argv: %v",
					has, tc.wantFlag, args)
			}
		})
	}
}

// TestPostCloneCheckout_NoBranchIsNoOp guards the detached-HEAD
// case: a force-mode row with no Branch must NOT try to run
// `git checkout ""` — that would error and break the row even
// though the user's intent ("just clone, default HEAD is fine")
// was satisfied. Returns ("", true) without invoking git.
func TestPostCloneCheckout_NoBranchIsNoOp(t *testing.T) {
	detail, ok := runPostCloneCheckout(
		Row{Checkout: constants.CloneFromCheckoutForce},
		"/nonexistent/path/should-not-be-touched",
		"/also/nonexistent",
	)
	if !ok {
		t.Fatalf("ok=false, want true (no-op for empty branch)")
	}

	if len(detail) != 0 {
		t.Errorf("detail = %q, want empty", detail)
	}
}

// TestPostCloneCheckout_AutoModeIsNoOp confirms the hot path stays
// branchless for the default mode. A failure here would mean every
// auto-mode row pays the cost of an extra exec.
func TestPostCloneCheckout_AutoModeIsNoOp(t *testing.T) {
	detail, ok := runPostCloneCheckout(
		Row{Branch: "main"}, // empty Checkout → auto
		"/nonexistent",
		"/also/nonexistent",
	)
	if !ok || len(detail) != 0 {
		t.Fatalf("auto-mode hook fired: detail=%q ok=%v", detail, ok)
	}
}

// TestValidateRow_RejectsBadCheckout pins parse-time validation:
// a `checkout` value other than "" / auto / skip / force errors out
// before any clone runs.
func TestValidateRow_RejectsBadCheckout(t *testing.T) {
	err := validateRow(Row{URL: "https://x/y.git", Checkout: "bogus"})
	if err == nil {
		t.Fatalf("validateRow accepted bogus checkout")
	}

	if !strings.Contains(err.Error(), "bogus") {
		t.Errorf("error %q does not mention bad value", err.Error())
	}
}
