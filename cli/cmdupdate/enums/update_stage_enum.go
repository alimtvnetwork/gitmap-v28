package enums

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// UpdateStageType represents a sequential execution phase in gitmap update.
type UpdateStageType int

// UpdateStage is an alias for UpdateStageType.
type UpdateStage = UpdateStageType

const (
	StageNone UpdateStageType = iota
	StageResolving
	StageDownloading
	StageVerifying
	StageSwapping
	StageComplete
)

var stageLabels = [...]string{"Preparing", "Resolving", "Downloading", "Verifying", "Swapping", "Complete"}

// IsValid reports whether the stage matches a known update stage.
func (s UpdateStageType) IsValid() bool {
	return s >= StageNone && s <= StageComplete
}

// String returns human-readable stage label.
func (s UpdateStageType) String() string {
	if int(s) > 0 && int(s) < len(stageLabels) {
		return stageLabels[s]
	}

	return stageLabels[0]
}

// StepNumber returns the 1-based step index.
func (s UpdateStageType) StepNumber() int {
	return int(s)
}

// TotalSteps returns total count of update stages.
func (s UpdateStageType) TotalSteps() int {
	return 5
}

// ParseUpdateStageType parses a raw string (label or number) into a validated UpdateStageType.
func ParseUpdateStageType(s string) (UpdateStageType, *apperror.AppError) {
	cleaned := strings.TrimSpace(s)
	if num, err := strconv.Atoi(cleaned); err == nil {
		val := UpdateStageType(num)
		if val.IsValid() {
			return val, nil
		}
	}

	for _, val := range AllUpdateStageTypes() {
		if strings.EqualFold(val.String(), cleaned) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseUpdateStageType", "E1000", fmt.Sprintf("invalid update stage: %q", s)).WithContext("value", s)
}

// AllUpdateStageTypes returns all valid UpdateStageType variants.
func AllUpdateStageTypes() []UpdateStageType {
	return []UpdateStageType{StageNone, StageResolving, StageDownloading, StageVerifying, StageSwapping, StageComplete}
}

// IsNone reports whether the stage is none.
func (s UpdateStageType) IsNone() bool { return s == StageNone }

// IsResolving reports whether the stage is resolving.
func (s UpdateStageType) IsResolving() bool { return s == StageResolving }

// IsDownloading reports whether the stage is downloading.
func (s UpdateStageType) IsDownloading() bool { return s == StageDownloading }

// IsVerifying reports whether the stage is verifying.
func (s UpdateStageType) IsVerifying() bool { return s == StageVerifying }

// IsSwapping reports whether the stage is swapping.
func (s UpdateStageType) IsSwapping() bool { return s == StageSwapping }

// IsComplete reports whether the stage is complete.
func (s UpdateStageType) IsComplete() bool { return s == StageComplete }
