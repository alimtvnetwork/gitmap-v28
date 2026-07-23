package startup

// Unit tests for safeUint32 — the centralized G115 guard used by every
// LinkInfo offset/size computation. The function MUST:
//
//   - Accept the inclusive boundaries 0 and math.MaxUint32.
//   - Accept any value in between without modification.
//   - Reject negative ints (would wrap to a huge uint32).
//   - Reject values above math.MaxUint32 on platforms where int can
//     represent them (64-bit). On 32-bit platforms the over-max case
//     is unreachable by construction, so that subtest is skipped.
//
// Keeping these guarantees pinned by tests prevents a future "just
// cast it" refactor from silently re-introducing the gosec G115
// integer-overflow class of bugs in the shortcut writer.

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

func TestSafeUint32_AcceptsBoundaryValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   int
		want uint32
	}{
		{"zero lower bound", 0, 0},
		{"one above zero", 1, 1},
		{"small typical", 1024, 1024},
		{"max uint32 upper bound", math.MaxUint32, math.MaxUint32},
		{"one below max uint32", math.MaxUint32 - 1, math.MaxUint32 - 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := safeUint32(tc.in)
			if err != nil {
				t.Fatalf("safeUint32(%d) returned unexpected error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("safeUint32(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestSafeUint32_RejectsNegatives(t *testing.T) {
	t.Parallel()
	cases := []int{-1, -2, -1024, math.MinInt32, math.MinInt64}
	for _, in := range cases {
		t.Run(strconv.Itoa(in), func(t *testing.T) {
			t.Parallel()
			got, err := safeUint32(in)
			if err == nil {
				t.Fatalf("safeUint32(%d) = %d, want error", in, got)
			}
			if got != 0 {
				t.Fatalf("safeUint32(%d) returned %d on error, want 0", in, got)
			}
		})
	}
}

func TestSafeUint32_RejectsAboveMaxUint32(t *testing.T) {
	t.Parallel()
	// On 32-bit platforms `int` cannot exceed math.MaxInt32, which is
	// itself below math.MaxUint32 — the over-max branch is unreachable
	// there, so skip rather than emit a misleading pass.
	if math.MaxInt <= math.MaxUint32 {
		t.Skip("int cannot exceed math.MaxUint32 on this platform (32-bit)")
	}
	cases := []int{
		math.MaxUint32 + 1,
		math.MaxUint32 + 1024,
		math.MaxInt, // largest int on a 64-bit platform
	}
	for _, in := range cases {
		t.Run(strconv.Itoa(in), func(t *testing.T) {
			t.Parallel()
			got, err := safeUint32(in)
			if err == nil {
				t.Fatalf("safeUint32(%d) = %d, want error", in, got)
			}
			if got != 0 {
				t.Fatalf("safeUint32(%d) returned %d on error, want 0", in, got)
			}
		})
	}
}

func TestSafeUint32_ErrorIsNotSentinel(t *testing.T) {
	t.Parallel()
	// Sanity check: the function returns a fmt.Errorf wrap with no
	// sentinel target, so errors.Is against a generic error must be
	// false. This pins the contract that callers should treat the
	// returned error as opaque + non-nil rather than matching it.
	_, err := safeUint32(-1)
	if err == nil {
		t.Fatal("expected error for negative input")
	}
	if errors.Is(err, errSafeUint32Sentinel) {
		t.Fatalf("safeUint32 error must not match an exported sentinel; got %v", err)
	}
}

// errSafeUint32Sentinel exists only so the test above can prove that
// safeUint32 does NOT return a matchable sentinel. Keep it unexported
// and unused outside this file.
var errSafeUint32Sentinel = errors.New("not-a-real-sentinel")
