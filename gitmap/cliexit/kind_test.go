package cliexit

// Tests for the Kind → exit-code mapping. Lock the numeric table so
// wrapper scripts (and CloneNowExitConfirmAborted=2 in particular)
// can't drift silently.

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// TestKindCode_Table pins every Kind to its documented exit code.
// If you change a value here, audit wrapper scripts and the
// cliexit_*_test.go assertions before merging.
func TestKindCode_Table(t *testing.T) {
	t.Parallel()
	cases := []struct {
		kind Kind
		want int
	}{
		{KindSuccess, 0},
		{KindExecutionFailed, 1},
		{KindUserCanceled, 2},
		{KindInvalidInput, 2},
		{KindVerifyFailed, 3},
		{KindPreconditionFailed, 4},
	}
	for _, tc := range cases {
		if got := KindCode(tc.kind); got != tc.want {
			t.Fatalf("KindCode(%v) = %d, want %d", tc.kind, got, tc.want)
		}
	}
}

// TestKindCode_UnknownDefaultsToOne guards against a future enum
// addition that forgets the table — we'd rather exit 1 than 0.
func TestKindCode_UnknownDefaultsToOne(t *testing.T) {
	t.Parallel()
	if got := KindCode(Kind(999)); got != 1 {
		t.Fatalf("unknown Kind should default to 1, got %d", got)
	}
}

// TestWithKindExtra_TagsContext verifies the kind label is added to
// Extras (so JSON consumers can branch on a string) without mutating
// the caller's original map.
func TestWithKindExtra_TagsContext(t *testing.T) {
	t.Parallel()
	original := map[string]string{"row": "3"}
	ctx := Context{
		Command: "clone-now",
		Op:      "git-clone",
		Extras:  original,
		Err:     errors.New("boom"),
	}
	tagged := withKindExtra(ctx, KindUserCanceled)
	if tagged.Extras["kind"] != "user-canceled" {
		t.Fatalf("kind label missing/wrong: %v", tagged.Extras)
	}
	if _, leaked := original["kind"]; leaked {
		t.Fatalf("withKindExtra mutated caller's map: %v", original)
	}
}

// TestFailKind_RendersKindInOutput drives the human renderer through
// withKindExtra to confirm the kind tag actually surfaces. We can't
// call FailKind directly (it os.Exits) so we exercise the rendering
// path it composes.
func TestFailKind_RendersKindInOutput(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	tagged := withKindExtra(Context{
		Command: "scan",
		Op:      "parse",
		Path:    "/tmp/x",
		Err:     errors.New("bad json"),
	}, KindInvalidInput)
	writeStructured(&buf, tagged, OutputHuman)
	if !strings.Contains(buf.String(), "kind=invalid-input") {
		t.Fatalf("expected kind=invalid-input in output:\n%s", buf.String())
	}
}
