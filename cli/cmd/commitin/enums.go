// Package commitin contains the typed enums and shared types for the
// `gitmap commit-in` (cin) command. See spec/03-commit-in/.
//
// Enum design (per Core memory rules):
//   - Each enum is its own named uint8 type with a String() returning
//     the canonical PascalCase name from constants_commitin.go.
//   - Each enum has a single AllX() slice so DB seed code, parity
//     tests, and CLI validators iterate one source of truth.
//   - No magic strings — every literal is a constant.
package commitin

import "github.com/alimtvnetwork/gitmap-v28/cli/constants"

// ---- ConflictMode -------------------------------------------------
type ConflictModeType uint8

const (
	ConflictModeForceMerge ConflictModeType = iota + 1
	ConflictModePrompt
)

func (m ConflictModeType) String() string {
	switch m {
	case ConflictModeForceMerge:
		return constants.CommitInConflictModeForceMerge
	case ConflictModePrompt:
		return constants.CommitInConflictModePrompt
	}

	return ""
}

func AllConflictModes() []ConflictModeType {
	return []ConflictModeType{ConflictModeForceMerge, ConflictModePrompt}
}

// ---- InputKind ----------------------------------------------------
type InputKindType uint8

const (
	InputKindLocalFolder InputKindType = iota + 1
	InputKindGitUrl
	InputKindVersionedSibling
)

func (k InputKindType) String() string {
	switch k {
	case InputKindLocalFolder:
		return constants.CommitInInputKindLocalFolder
	case InputKindGitUrl:
		return constants.CommitInInputKindGitUrl
	case InputKindVersionedSibling:
		return constants.CommitInInputKindVersionedSibling
	}

	return ""
}

func AllInputKinds() []InputKindType {
	return []InputKindType{InputKindLocalFolder, InputKindGitUrl, InputKindVersionedSibling}
}

// ---- RunStatus ----------------------------------------------------
type RunStatusType uint8

const (
	RunStatusPending RunStatusType = iota + 1
	RunStatusRunning
	RunStatusCompleted
	RunStatusFailed
	RunStatusPartiallyFailed
)

func (s RunStatusType) String() string {
	switch s {
	case RunStatusPending:
		return constants.CommitInRunStatusPending
	case RunStatusRunning:
		return constants.CommitInRunStatusRunning
	case RunStatusCompleted:
		return constants.CommitInRunStatusCompleted
	case RunStatusFailed:
		return constants.CommitInRunStatusFailed
	case RunStatusPartiallyFailed:
		return constants.CommitInRunStatusPartiallyFailed
	}

	return ""
}

func AllRunStatuses() []RunStatusType {
	return []RunStatusType{RunStatusPending, RunStatusRunning, RunStatusCompleted, RunStatusFailed, RunStatusPartiallyFailed}
}

// ---- CommitOutcome ------------------------------------------------
type CommitOutcomeType uint8

const (
	CommitOutcomeCreated CommitOutcomeType = iota + 1
	CommitOutcomeSkipped
	CommitOutcomeFailed
)

func (o CommitOutcomeType) String() string {
	switch o {
	case CommitOutcomeCreated:
		return constants.CommitInOutcomeCreated
	case CommitOutcomeSkipped:
		return constants.CommitInOutcomeSkipped
	case CommitOutcomeFailed:
		return constants.CommitInOutcomeFailed
	}

	return ""
}

func AllCommitOutcomes() []CommitOutcomeType {
	return []CommitOutcomeType{CommitOutcomeCreated, CommitOutcomeSkipped, CommitOutcomeFailed}
}

// ---- SkipReason ---------------------------------------------------
type SkipReasonType uint8

const (
	SkipReasonDuplicateSourceSha SkipReasonType = iota + 1
	SkipReasonExcludedAllFiles
	SkipReasonEmptyAfterMessageRules
	SkipReasonDryRun
)

func (r SkipReasonType) String() string {
	switch r {
	case SkipReasonDuplicateSourceSha:
		return constants.CommitInSkipReasonDuplicateSourceSha
	case SkipReasonExcludedAllFiles:
		return constants.CommitInSkipReasonExcludedAllFiles
	case SkipReasonEmptyAfterMessageRules:
		return constants.CommitInSkipReasonEmptyAfterMessageRules
	case SkipReasonDryRun:
		return constants.CommitInSkipReasonDryRun
	}

	return ""
}

func AllSkipReasons() []SkipReasonType {
	return []SkipReasonType{SkipReasonDuplicateSourceSha, SkipReasonExcludedAllFiles, SkipReasonEmptyAfterMessageRules, SkipReasonDryRun}
}

// ---- ExclusionKind ------------------------------------------------
type ExclusionKindType uint8

const (
	ExclusionKindPathFolder ExclusionKindType = iota + 1
	ExclusionKindPathFile
)

func (k ExclusionKindType) String() string {
	switch k {
	case ExclusionKindPathFolder:
		return constants.CommitInExclusionKindPathFolder
	case ExclusionKindPathFile:
		return constants.CommitInExclusionKindPathFile
	}

	return ""
}

func AllExclusionKinds() []ExclusionKindType {
	return []ExclusionKindType{ExclusionKindPathFolder, ExclusionKindPathFile}
}

// ---- MessageRuleKind ----------------------------------------------
type MessageRuleKindType uint8

const (
	MessageRuleKindStartsWith MessageRuleKindType = iota + 1
	MessageRuleKindEndsWith
	MessageRuleKindContains
)

func (k MessageRuleKindType) String() string {
	switch k {
	case MessageRuleKindStartsWith:
		return constants.CommitInMessageRuleKindStartsWith
	case MessageRuleKindEndsWith:
		return constants.CommitInMessageRuleKindEndsWith
	case MessageRuleKindContains:
		return constants.CommitInMessageRuleKindContains
	}

	return ""
}

func AllMessageRuleKinds() []MessageRuleKindType {
	return []MessageRuleKindType{MessageRuleKindStartsWith, MessageRuleKindEndsWith, MessageRuleKindContains}
}

// ---- FunctionIntelLanguage ----------------------------------------
type FunctionIntelLanguageType uint8

const (
	LanguageGo FunctionIntelLanguageType = iota + 1
	LanguageJavaScript
	LanguageTypeScript
	LanguageRust
	LanguagePython
	LanguagePhp
	LanguageJava
	LanguageCSharp
)

func (l FunctionIntelLanguageType) String() string {
	switch l {
	case LanguageGo:
		return constants.CommitInLanguageGo
	case LanguageJavaScript:
		return constants.CommitInLanguageJavaScript
	case LanguageTypeScript:
		return constants.CommitInLanguageTypeScript
	case LanguageRust:
		return constants.CommitInLanguageRust
	case LanguagePython:
		return constants.CommitInLanguagePython
	case LanguagePhp:
		return constants.CommitInLanguagePhp
	case LanguageJava:
		return constants.CommitInLanguageJava
	case LanguageCSharp:
		return constants.CommitInLanguageCSharp
	}

	return ""
}

func AllLanguages() []FunctionIntelLanguageType {
	return []FunctionIntelLanguageType{LanguageGo, LanguageJavaScript, LanguageTypeScript, LanguageRust, LanguagePython, LanguagePhp, LanguageJava, LanguageCSharp}
}

// Backward-compatible type aliases.
type (
	ConflictMode          = ConflictModeType
	InputKind             = InputKindType
	RunStatus             = RunStatusType
	CommitOutcome         = CommitOutcomeType
	SkipReason            = SkipReasonType
	ExclusionKind         = ExclusionKindType
	MessageRuleKind       = MessageRuleKindType
	FunctionIntelLanguage = FunctionIntelLanguageType
)
