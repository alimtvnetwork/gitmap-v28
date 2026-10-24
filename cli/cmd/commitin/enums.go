// Package commitin contains the shared types for the
// `gitmap commit-in` (cin) command. See 02-spec/03-commit-in/.
//
// The typed enums themselves live in cli/cmd/commitin/enums; this file
// only re-exports them so existing in-package and subpackage call sites
// keep compiling unchanged.
package commitin

import "github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/enums"

// ---- ConflictMode -------------------------------------------------
type ConflictModeType = enums.ConflictModeType

const (
	ConflictModeForceMerge = enums.ConflictModeForceMerge
	ConflictModePrompt     = enums.ConflictModePrompt
)

func AllConflictModes() []ConflictModeType {
	return enums.AllConflictModes()
}

// ---- InputKind ----------------------------------------------------
type InputKindType = enums.InputKindType

const (
	InputKindLocalFolder      = enums.InputKindLocalFolder
	InputKindGitUrl           = enums.InputKindGitUrl
	InputKindVersionedSibling = enums.InputKindVersionedSibling
)

func AllInputKinds() []InputKindType {
	return enums.AllInputKinds()
}

// ---- RunStatus ----------------------------------------------------
type RunStatusType = enums.RunStatusType

const (
	RunStatusPending         = enums.RunStatusPending
	RunStatusRunning         = enums.RunStatusRunning
	RunStatusCompleted       = enums.RunStatusCompleted
	RunStatusFailed          = enums.RunStatusFailed
	RunStatusPartiallyFailed = enums.RunStatusPartiallyFailed
)

func AllRunStatuses() []RunStatusType {
	return enums.AllRunStatuses()
}

// ---- CommitOutcome ------------------------------------------------
type CommitOutcomeType = enums.CommitOutcomeType

const (
	CommitOutcomeCreated = enums.CommitOutcomeCreated
	CommitOutcomeSkipped = enums.CommitOutcomeSkipped
	CommitOutcomeFailed  = enums.CommitOutcomeFailed
)

func AllCommitOutcomes() []CommitOutcomeType {
	return enums.AllCommitOutcomes()
}

// ---- SkipReason ---------------------------------------------------
type SkipReasonType = enums.SkipReasonType

const (
	SkipReasonDuplicateSourceSha     = enums.SkipReasonDuplicateSourceSha
	SkipReasonExcludedAllFiles       = enums.SkipReasonExcludedAllFiles
	SkipReasonEmptyAfterMessageRules = enums.SkipReasonEmptyAfterMessageRules
	SkipReasonDryRun                 = enums.SkipReasonDryRun
)

func AllSkipReasons() []SkipReasonType {
	return enums.AllSkipReasons()
}

// ---- ExclusionKind ------------------------------------------------
type ExclusionKindType = enums.ExclusionKindType

const (
	ExclusionKindPathFolder = enums.ExclusionKindPathFolder
	ExclusionKindPathFile   = enums.ExclusionKindPathFile
)

func AllExclusionKinds() []ExclusionKindType {
	return enums.AllExclusionKinds()
}

// ---- MessageRuleKind ----------------------------------------------
type MessageRuleKindType = enums.MessageRuleKindType

const (
	MessageRuleKindStartsWith = enums.MessageRuleKindStartsWith
	MessageRuleKindEndsWith   = enums.MessageRuleKindEndsWith
	MessageRuleKindContains   = enums.MessageRuleKindContains
)

func AllMessageRuleKinds() []MessageRuleKindType {
	return enums.AllMessageRuleKinds()
}

// ---- FunctionIntelLanguage ----------------------------------------
type FunctionIntelLanguageType = enums.FunctionIntelLanguageType

const (
	LanguageGo         = enums.LanguageGo
	LanguageJavaScript = enums.LanguageJavaScript
	LanguageTypeScript = enums.LanguageTypeScript
	LanguageRust       = enums.LanguageRust
	LanguagePython     = enums.LanguagePython
	LanguagePhp        = enums.LanguagePhp
	LanguageJava       = enums.LanguageJava
	LanguageCSharp     = enums.LanguageCSharp
)

func AllLanguages() []FunctionIntelLanguageType {
	return enums.AllLanguages()
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
