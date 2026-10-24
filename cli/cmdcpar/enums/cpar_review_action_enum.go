package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// CPARReviewActionType represents user choice in review flow.
type CPARReviewActionType string

const (
	ActionCommitAllFeature CPARReviewActionType = "feature"
	ActionCommitAllBug     CPARReviewActionType = "bug"
	ActionSingleRepo       CPARReviewActionType = "single"
	ActionCommitAllChore   CPARReviewActionType = "chore"
	ActionAbort            CPARReviewActionType = "abort"
)

// IsValid reports whether the value matches a known variant.
func (c CPARReviewActionType) IsValid() bool {
	switch c {
	case ActionCommitAllFeature, ActionCommitAllBug, ActionSingleRepo, ActionCommitAllChore, ActionAbort:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (c CPARReviewActionType) String() string {
	return string(c)
}

// ParseCPARReviewActionType parses a raw string into a validated CPARReviewActionType.
func ParseCPARReviewActionType(s string) (CPARReviewActionType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllCPARReviewActionTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "actioncommitallfeature":
		return ActionCommitAllFeature, nil
	case "commitallfeature":
		return ActionCommitAllFeature, nil
	case "actioncommitallbug":
		return ActionCommitAllBug, nil
	case "commitallbug":
		return ActionCommitAllBug, nil
	case "actionsinglerepo":
		return ActionSingleRepo, nil
	case "singlerepo":
		return ActionSingleRepo, nil
	case "actioncommitallchore":
		return ActionCommitAllChore, nil
	case "commitallchore":
		return ActionCommitAllChore, nil
	case "actionabort":
		return ActionAbort, nil
	case "abort":
		return ActionAbort, nil
	}
	return "", apperror.NewValidation("ParseCPARReviewActionType", "E1000", fmt.Sprintf("invalid cparreviewactiontype: %q", s)).WithContext("value", s)
}

// AllCPARReviewActionTypes returns all valid variants.
func AllCPARReviewActionTypes() []CPARReviewActionType {
	return []CPARReviewActionType{
		ActionCommitAllFeature,
		ActionCommitAllBug,
		ActionSingleRepo,
		ActionCommitAllChore,
		ActionAbort,
	}
}

// IsCommitAllFeature reports whether the value is ActionCommitAllFeature.
func (c CPARReviewActionType) IsCommitAllFeature() bool {
	return c == ActionCommitAllFeature
}

// IsCommitAllBug reports whether the value is ActionCommitAllBug.
func (c CPARReviewActionType) IsCommitAllBug() bool {
	return c == ActionCommitAllBug
}

// IsSingleRepo reports whether the value is ActionSingleRepo.
func (c CPARReviewActionType) IsSingleRepo() bool {
	return c == ActionSingleRepo
}

// IsCommitAllChore reports whether the value is ActionCommitAllChore.
func (c CPARReviewActionType) IsCommitAllChore() bool {
	return c == ActionCommitAllChore
}

// IsAbort reports whether the value is ActionAbort.
func (c CPARReviewActionType) IsAbort() bool {
	return c == ActionAbort
}

// HasReviewAction reports whether the value is non-empty and valid.
func (c CPARReviewActionType) HasReviewAction() bool {
	return c.IsValid()
}
