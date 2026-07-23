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

import "github.com/alimtvnetwork/gitmap-v27/gitmap/constants"

// ---- ConflictMode -------------------------------------------------
type ConflictMode uint8

const (
	ConflictModeForceMerge ConflictMode = iota + 1
	ConflictModePrompt
)

func (m ConflictMode) String() string {
	switch m {
	case ConflictModeForceMerge:
		return constants.CommitInConflictModeForceMerge
	case ConflictModePrompt:
		return constants.CommitInConflictModePrompt
	}
	return ""
}

func AllConflictModes() []ConflictMode {
	return []ConflictMode{ConflictModeForceMerge, ConflictModePrompt}
}

// ---- InputKind ----------------------------------------------------
type InputKind uint8

const (
	InputKindLocalFolder InputKind = iota + 1
	InputKindGitUrl
	InputKindVersionedSibling
)

func (k InputKind) String() string {
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

func AllInputKinds() []InputKind {
	return []InputKind{InputKindLocalFolder, InputKindGitUrl, InputKindVersionedSibling}
}

// ---- RunStatus ----------------------------------------------------
type RunStatus uint8

const (
	RunStatusPending RunStatus = iota + 1
	RunStatusRunning
	RunStatusCompleted
	RunStatusFailed
	RunStatusPartiallyFailed
)

func (s RunStatus) String() string {
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

func AllRunStatuses() []RunStatus {
	return []RunStatus{RunStatusPending, RunStatusRunning, RunStatusCompleted, RunStatusFailed, RunStatusPartiallyFailed}
}

// ---- CommitOutcome ------------------------------------------------
type CommitOutcome uint8

const (
	CommitOutcomeCreated CommitOutcome = iota + 1
	CommitOutcomeSkipped
	CommitOutcomeFailed
)

func (o CommitOutcome) String() string {
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

func AllCommitOutcomes() []CommitOutcome {
	return []CommitOutcome{CommitOutcomeCreated, CommitOutcomeSkipped, CommitOutcomeFailed}
}

// ---- SkipReason ---------------------------------------------------
type SkipReason uint8

const (
	SkipReasonDuplicateSourceSha SkipReason = iota + 1
	SkipReasonExcludedAllFiles
	SkipReasonEmptyAfterMessageRules
	SkipReasonDryRun
)

func (r SkipReason) String() string {
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

func AllSkipReasons() []SkipReason {
	return []SkipReason{SkipReasonDuplicateSourceSha, SkipReasonExcludedAllFiles, SkipReasonEmptyAfterMessageRules, SkipReasonDryRun}
}

// ---- ExclusionKind ------------------------------------------------
type ExclusionKind uint8

const (
	ExclusionKindPathFolder ExclusionKind = iota + 1
	ExclusionKindPathFile
)

func (k ExclusionKind) String() string {
	switch k {
	case ExclusionKindPathFolder:
		return constants.CommitInExclusionKindPathFolder
	case ExclusionKindPathFile:
		return constants.CommitInExclusionKindPathFile
	}
	return ""
}

func AllExclusionKinds() []ExclusionKind {
	return []ExclusionKind{ExclusionKindPathFolder, ExclusionKindPathFile}
}

// ---- MessageRuleKind ----------------------------------------------
type MessageRuleKind uint8

const (
	MessageRuleKindStartsWith MessageRuleKind = iota + 1
	MessageRuleKindEndsWith
	MessageRuleKindContains
)

func (k MessageRuleKind) String() string {
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

func AllMessageRuleKinds() []MessageRuleKind {
	return []MessageRuleKind{MessageRuleKindStartsWith, MessageRuleKindEndsWith, MessageRuleKindContains}
}

// ---- FunctionIntelLanguage ----------------------------------------
type FunctionIntelLanguage uint8

const (
	LanguageGo FunctionIntelLanguage = iota + 1
	LanguageJavaScript
	LanguageTypeScript
	LanguageRust
	LanguagePython
	LanguagePhp
	LanguageJava
	LanguageCSharp
)

func (l FunctionIntelLanguage) String() string {
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

func AllLanguages() []FunctionIntelLanguage {
	return []FunctionIntelLanguage{LanguageGo, LanguageJavaScript, LanguageTypeScript, LanguageRust, LanguagePython, LanguagePhp, LanguageJava, LanguageCSharp}
}
