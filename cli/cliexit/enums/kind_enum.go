package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// KindType classifies a CLI failure into one of the documented buckets.
// Use FailKind(ctx, kind, mode) at the cmd entry-point error path so
// the (message, exit-code) pair stays atomic.
type KindType int

const (
	// KindSuccess is the zero-value sentinel. Never pass it to
	// FailKind — it's only here so `var k KindType` defaults to a
	// meaningful name in debuggers.
	KindSuccess KindType = 0
	// KindExecutionFailed covers any operational failure that
	// happened while doing real work: git clone died, a file
	// couldn't be written, a network call timed out. Exit 1.
	KindExecutionFailed KindType = 1
	// KindUserCanceled is "the user said no" — interactive
	// abort, non-TTY without --yes, Ctrl-C at a confirm prompt.
	// Exit 2. Distinct semantics from KindInvalidInput but
	// shares the code (see file header).
	KindUserCanceled KindType = 2
	// KindInvalidInput is bad CLI input: unknown flag, missing
	// positional, malformed manifest, unparseable URL. Exit 2.
	KindInvalidInput KindType = 3
	// KindVerifyFailed is a read-only check that found a
	// divergence (clone-audit mismatch, verify-cmd-faithful
	// mismatch). Exit 3 — distinguishable from operational
	// failures so CI can branch on "drift vs broken".
	KindVerifyFailed KindType = 4
	// KindPreconditionFailed means the environment isn't ready:
	// `git` not on PATH, DB locked by another process, required
	// secret missing. Exit 4 — separate from input errors so
	// retry logic can wait-and-retry instead of giving up.
	KindPreconditionFailed KindType = 5
)

// IsValid reports whether the value matches a known variant.
func (k KindType) IsValid() bool {
	switch k {
	case KindSuccess, KindExecutionFailed, KindUserCanceled, KindInvalidInput, KindVerifyFailed, KindPreconditionFailed:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (k KindType) String() string {
	return strconv.Itoa(int(k))
}

// ParseKindType parses a raw string into a validated KindType.
func ParseKindType(s string) (KindType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := KindType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllKindTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "kindsuccess":
		return KindSuccess, nil
	case "success":
		return KindSuccess, nil
	case "kindexecutionfailed":
		return KindExecutionFailed, nil
	case "executionfailed":
		return KindExecutionFailed, nil
	case "kindusercanceled":
		return KindUserCanceled, nil
	case "usercanceled":
		return KindUserCanceled, nil
	case "kindinvalidinput":
		return KindInvalidInput, nil
	case "invalidinput":
		return KindInvalidInput, nil
	case "kindverifyfailed":
		return KindVerifyFailed, nil
	case "verifyfailed":
		return KindVerifyFailed, nil
	case "kindpreconditionfailed":
		return KindPreconditionFailed, nil
	case "preconditionfailed":
		return KindPreconditionFailed, nil
	}
	return 0, apperror.NewValidation("ParseKindType", "E1000", fmt.Sprintf("invalid kindtype: %q", s)).WithContext("value", s)
}

// AllKindTypes returns all valid variants.
func AllKindTypes() []KindType {
	return []KindType{
		KindSuccess,
		KindExecutionFailed,
		KindUserCanceled,
		KindInvalidInput,
		KindVerifyFailed,
		KindPreconditionFailed,
	}
}

// IsSuccess reports whether the value is KindSuccess.
func (k KindType) IsSuccess() bool {
	return k == KindSuccess
}

// IsExecutionFailed reports whether the value is KindExecutionFailed.
func (k KindType) IsExecutionFailed() bool {
	return k == KindExecutionFailed
}

// IsUserCanceled reports whether the value is KindUserCanceled.
func (k KindType) IsUserCanceled() bool {
	return k == KindUserCanceled
}

// IsInvalidInput reports whether the value is KindInvalidInput.
func (k KindType) IsInvalidInput() bool {
	return k == KindInvalidInput
}

// IsVerifyFailed reports whether the value is KindVerifyFailed.
func (k KindType) IsVerifyFailed() bool {
	return k == KindVerifyFailed
}

// IsPreconditionFailed reports whether the value is KindPreconditionFailed.
func (k KindType) IsPreconditionFailed() bool {
	return k == KindPreconditionFailed
}

// HasKind reports whether the value is non-empty and valid.
func (k KindType) HasKind() bool {
	return k.IsValid()
}
