package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// IsValid reports whether the conflict mode matches a known mode.
func (m ConflictModeType) IsValid() bool {
	switch m {
	case ConflictModeForceMerge, ConflictModePrompt:
		return true
	default:
		return false
	}
}

// ParseConflictModeType parses a canonical name into a validated ConflictModeType.
func ParseConflictModeType(s string) (ConflictModeType, *apperror.AppError) {
	for _, val := range AllConflictModes() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseConflictModeType", "E1000", fmt.Sprintf("invalid conflict mode: %q", s)).WithContext("value", s)
}

// IsValid reports whether the input kind matches a known kind.
func (k InputKindType) IsValid() bool {
	switch k {
	case InputKindLocalFolder, InputKindGitUrl, InputKindVersionedSibling:
		return true
	default:
		return false
	}
}

// ParseInputKindType parses a canonical name into a validated InputKindType.
func ParseInputKindType(s string) (InputKindType, *apperror.AppError) {
	for _, val := range AllInputKinds() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseInputKindType", "E1000", fmt.Sprintf("invalid input kind: %q", s)).WithContext("value", s)
}

// IsValid reports whether the run status matches a known status.
func (s RunStatusType) IsValid() bool {
	switch s {
	case RunStatusPending, RunStatusRunning, RunStatusCompleted, RunStatusFailed, RunStatusPartiallyFailed:
		return true
	default:
		return false
	}
}

// ParseRunStatusType parses a canonical name into a validated RunStatusType.
func ParseRunStatusType(s string) (RunStatusType, *apperror.AppError) {
	for _, val := range AllRunStatuses() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseRunStatusType", "E1000", fmt.Sprintf("invalid run status: %q", s)).WithContext("value", s)
}

// IsValid reports whether the commit outcome matches a known outcome.
func (o CommitOutcomeType) IsValid() bool {
	switch o {
	case CommitOutcomeCreated, CommitOutcomeSkipped, CommitOutcomeFailed:
		return true
	default:
		return false
	}
}

// ParseCommitOutcomeType parses a canonical name into a validated CommitOutcomeType.
func ParseCommitOutcomeType(s string) (CommitOutcomeType, *apperror.AppError) {
	for _, val := range AllCommitOutcomes() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseCommitOutcomeType", "E1000", fmt.Sprintf("invalid commit outcome: %q", s)).WithContext("value", s)
}

// IsValid reports whether the skip reason matches a known reason.
func (r SkipReasonType) IsValid() bool {
	switch r {
	case SkipReasonDuplicateSourceSha, SkipReasonExcludedAllFiles, SkipReasonEmptyAfterMessageRules, SkipReasonDryRun:
		return true
	default:
		return false
	}
}

// ParseSkipReasonType parses a canonical name into a validated SkipReasonType.
func ParseSkipReasonType(s string) (SkipReasonType, *apperror.AppError) {
	for _, val := range AllSkipReasons() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseSkipReasonType", "E1000", fmt.Sprintf("invalid skip reason: %q", s)).WithContext("value", s)
}

// IsValid reports whether the exclusion kind matches a known kind.
func (k ExclusionKindType) IsValid() bool {
	switch k {
	case ExclusionKindPathFolder, ExclusionKindPathFile:
		return true
	default:
		return false
	}
}

// ParseExclusionKindType parses a canonical name into a validated ExclusionKindType.
func ParseExclusionKindType(s string) (ExclusionKindType, *apperror.AppError) {
	for _, val := range AllExclusionKinds() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseExclusionKindType", "E1000", fmt.Sprintf("invalid exclusion kind: %q", s)).WithContext("value", s)
}

// IsValid reports whether the message rule kind matches a known kind.
func (k MessageRuleKindType) IsValid() bool {
	switch k {
	case MessageRuleKindStartsWith, MessageRuleKindEndsWith, MessageRuleKindContains:
		return true
	default:
		return false
	}
}

// ParseMessageRuleKindType parses a canonical name into a validated MessageRuleKindType.
func ParseMessageRuleKindType(s string) (MessageRuleKindType, *apperror.AppError) {
	for _, val := range AllMessageRuleKinds() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseMessageRuleKindType", "E1000", fmt.Sprintf("invalid message rule kind: %q", s)).WithContext("value", s)
}

// IsValid reports whether the function-intel language matches a known language.
func (l FunctionIntelLanguageType) IsValid() bool {
	switch l {
	case LanguageGo, LanguageJavaScript, LanguageTypeScript, LanguageRust, LanguagePython, LanguagePhp, LanguageJava, LanguageCSharp:
		return true
	default:
		return false
	}
}

// ParseFunctionIntelLanguageType parses a canonical name into a validated FunctionIntelLanguageType.
func ParseFunctionIntelLanguageType(s string) (FunctionIntelLanguageType, *apperror.AppError) {
	for _, val := range AllLanguages() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseFunctionIntelLanguageType", "E1000", fmt.Sprintf("invalid function intel language: %q", s)).WithContext("value", s)
}
