package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// MergeModeType picks the tag-merge strategy used by SyncMode / SyncAtMode.
// uint8 per the project-wide rule "smallest viable int type for ≤255 members".
type MergeModeType uint8

const (
	// MergeModeUnion is the v4.36.0 default — existing ∪ incoming, dedup'd.
	MergeModeUnion MergeModeType = 0
	// MergeModeReplace overwrites with incoming verbatim. Detector keeps brand.
	MergeModeReplace MergeModeType = 1
	// MergeModeIntersection keeps only tags present in BOTH sets, plus brand pin.
	MergeModeIntersection MergeModeType = 2
)

// IsValid reports whether the value matches a known variant.
func (m MergeModeType) IsValid() bool {
	switch m {
	case MergeModeUnion, MergeModeReplace, MergeModeIntersection:
		return true
	default:
		return false
	}
}

func (m MergeModeType) String() string {
	switch m {
	case MergeModeReplace:
		return constants.VSCodePMSyncModeReplace
	case MergeModeIntersection:
		return constants.VSCodePMSyncModeIntersection
	case MergeModeUnion:
		return constants.VSCodePMSyncModeUnion
	default:
		return constants.VSCodePMSyncModeUnion
	}
}

func ParseMergeMode(raw string) (MergeModeType, *apperror.AppError) {
	switch raw {
	case "", constants.VSCodePMSyncModeUnion:
		return MergeModeUnion, nil
	case constants.VSCodePMSyncModeReplace:
		return MergeModeReplace, nil
	case constants.VSCodePMSyncModeIntersection:
		return MergeModeIntersection, nil
	default:
		return MergeModeUnion, apperror.NewValidationError(fmt.Sprintf(
			constants.ErrVSCodePMSyncBadMode, raw,
			constants.VSCodePMSyncModeUnion,
			constants.VSCodePMSyncModeReplace,
			constants.VSCodePMSyncModeIntersection))
	}
}

// AllMergeModeTypes returns all valid variants.
func AllMergeModeTypes() []MergeModeType {
	return []MergeModeType{
		MergeModeUnion,
		MergeModeReplace,
		MergeModeIntersection,
	}
}

// IsUnion reports whether the value is MergeModeUnion.
func (m MergeModeType) IsUnion() bool {
	return m == MergeModeUnion
}

// IsReplace reports whether the value is MergeModeReplace.
func (m MergeModeType) IsReplace() bool {
	return m == MergeModeReplace
}

// IsIntersection reports whether the value is MergeModeIntersection.
func (m MergeModeType) IsIntersection() bool {
	return m == MergeModeIntersection
}

// HasMode reports whether the value is non-empty and valid.
func (m MergeModeType) HasMode() bool {
	return m.IsValid()
}
