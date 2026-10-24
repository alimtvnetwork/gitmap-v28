// Package suggestion defines canonical domain models, multi-tier resolution engines,
// and Catppuccin terminal renderers for CLI suggestions and error remediation.
package diag

import "github.com/alimtvnetwork/gitmap-v28/cli/diag/enums"

// SuggestionCategoryType classifies the underlying intent of a suggested action.
// SuggestionCategoryType aliases the strongly-typed enum from the enums package.
type SuggestionCategoryType = enums.SuggestionCategoryType

const (
	CategoryTypo        = enums.CategoryTypo
	CategorySubcommand  = enums.CategorySubcommand
	CategoryFlag        = enums.CategoryFlag
	CategoryRemediation = enums.CategoryRemediation
	CategoryIntent      = enums.CategoryIntent
	CategoryResource    = enums.CategoryResource
)

// SuggestionCategory is a deprecated legacy alias for SuggestionCategoryType.
type SuggestionCategory = enums.SuggestionCategoryType

// ActionType defines the execution mode of the suggested command.
// ActionType aliases the strongly-typed enum from the enums package.
type ActionType = enums.ActionType

const (
	ActionAutoRun = enums.ActionAutoRun
	ActionManual  = enums.ActionManual
	ActionCopy    = enums.ActionCopy
)

// Suggestion represents an actionable CLI recommendation.
type Suggestion struct {
	Command     string             `json:"command"`
	Description string             `json:"description"`
	Category    SuggestionCategory `json:"category"`
	Confidence  float64            `json:"confidence"` // Normalized score: 0.0 to 1.0
	ActionType  ActionType         `json:"actionType"`
	Context     map[string]any     `json:"context,omitempty"`
}

// SuggestionGroup aggregates related suggestions under an explanatory banner.
type SuggestionGroup struct {
	Title       string       `json:"title"`
	Reason      string       `json:"reason"`
	InputToken  string       `json:"inputToken,omitempty"`
	Suggestions []Suggestion `json:"suggestions"`
}

// RenderFormatType defines the rendering style for suggestions.
// RenderFormatType aliases the strongly-typed enum from the enums package.
type RenderFormatType = enums.RenderFormatType

const (
	RenderFormatBox     = enums.RenderFormatBox
	RenderFormatCompact = enums.RenderFormatCompact
	RenderFormatJSON    = enums.RenderFormatJSON
)

// RenderFormat is a deprecated legacy alias for RenderFormatType.
type RenderFormat = enums.RenderFormatType

// Engine defines the contract for resolving suggestions.
type Engine interface {
	ResolveCommand(token string) SuggestionGroup
	ResolveFlag(cmdName, flagToken string) SuggestionGroup
	ResolveRemediation(err error, ctx map[string]any) SuggestionGroup
}

// HasSuggestions reports whether the group contains one or more suggestions.
func (g SuggestionGroup) HasSuggestions() bool {
	return len(g.Suggestions) > 0
}

// HasInputToken reports whether the group has an associated input token.
func (g SuggestionGroup) HasInputToken() bool {
	return len(g.InputToken) > 0
}

// BestSuggestion returns the top-confidence recommendation if available.
func (g SuggestionGroup) BestSuggestion() (Suggestion, bool) {
	if len(g.Suggestions) == 0 {
		return Suggestion{}, false
	}
	return g.Suggestions[0], true
}

// IsHighConfidence reports whether confidence is at or above 0.8.
func (s Suggestion) IsHighConfidence() bool {
	return s.Confidence >= 0.80
}

// IsAutoRunnable reports whether this suggestion can be automatically executed.
func (s Suggestion) IsAutoRunnable() bool {
	return s.ActionType == ActionAutoRun
}
