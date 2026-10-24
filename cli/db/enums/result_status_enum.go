package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// ResultStatusType is a strongly-typed enum.
type ResultStatusType int

const (
	ResultStatusPending      ResultStatusType = 1
	ResultStatusSucceeded    ResultStatusType = 2
	ResultStatusFailed       ResultStatusType = 3
	ResultStatusSkipped      ResultStatusType = 4
	ResultStatusDeferred     ResultStatusType = 5
	ResultStatusRequiresAuth ResultStatusType = 6
)

// IsValid reports whether the value matches a known variant.
func (r ResultStatusType) IsValid() bool {
	switch r {
	case ResultStatusPending, ResultStatusSucceeded, ResultStatusFailed, ResultStatusSkipped, ResultStatusDeferred, ResultStatusRequiresAuth:
		return true
	default:
		return false
	}
}

func (s ResultStatusType) String() string {
	switch s {
	case ResultStatusPending:
		return "Pending"
	case ResultStatusSucceeded:
		return "Succeeded"
	case ResultStatusFailed:
		return "Failed"
	case ResultStatusSkipped:
		return "Skipped"
	case ResultStatusDeferred:
		return "Deferred"
	case ResultStatusRequiresAuth:
		return "RequiresAuth"
	default:
		return fmt.Sprintf("ResultStatusType(%d)", int(s))
	}
}

func ParseResultStatus(s string) (ResultStatusType, *apperror.AppError) {
	switch strings.ToLower(s) {
	case "pending":
		return ResultStatusPending, nil
	case "succeeded":
		return ResultStatusSucceeded, nil
	case "failed":
		return ResultStatusFailed, nil
	case "skipped":
		return ResultStatusSkipped, nil
	case "deferred":
		return ResultStatusDeferred, nil
	case "requiresauth":
		return ResultStatusRequiresAuth, nil
	default:
		return 0, apperror.NewValidationError(fmt.Sprintf("invalid ResultStatus: %s", s))
	}
}

// AllResultStatusTypes returns all valid variants.
func AllResultStatusTypes() []ResultStatusType {
	return []ResultStatusType{
		ResultStatusPending,
		ResultStatusSucceeded,
		ResultStatusFailed,
		ResultStatusSkipped,
		ResultStatusDeferred,
		ResultStatusRequiresAuth,
	}
}

// IsPending reports whether the value is ResultStatusPending.
func (r ResultStatusType) IsPending() bool {
	return r == ResultStatusPending
}

// IsSucceeded reports whether the value is ResultStatusSucceeded.
func (r ResultStatusType) IsSucceeded() bool {
	return r == ResultStatusSucceeded
}

// IsFailed reports whether the value is ResultStatusFailed.
func (r ResultStatusType) IsFailed() bool {
	return r == ResultStatusFailed
}

// IsSkipped reports whether the value is ResultStatusSkipped.
func (r ResultStatusType) IsSkipped() bool {
	return r == ResultStatusSkipped
}

// IsDeferred reports whether the value is ResultStatusDeferred.
func (r ResultStatusType) IsDeferred() bool {
	return r == ResultStatusDeferred
}

// IsRequiresAuth reports whether the value is ResultStatusRequiresAuth.
func (r ResultStatusType) IsRequiresAuth() bool {
	return r == ResultStatusRequiresAuth
}

// HasStatus reports whether the value is non-empty and valid.
func (r ResultStatusType) HasStatus() bool {
	return r.IsValid()
}
