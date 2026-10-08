// Package suggestion defines canonical domain models, multi-tier resolution engines,
// and Catppuccin terminal renderers for CLI suggestions and error remediation.
package suggestion

// SuggestionCategory classifies the underlying intent of a suggested action.
type SuggestionCategory string

const (
	// CategoryTypo indicates a mistyped command, flag, or token.
	CategoryTypo SuggestionCategory = "typo"
	// CategorySubcommand indicates a valid child subcommand recommendation.
	CategorySubcommand SuggestionCategory = "subcommand"
	// CategoryFlag indicates a supported flag recommendation matching input token.
	CategoryFlag SuggestionCategory = "flag"
	// CategoryRemediation indicates a prescriptive diagnostic step to fix a failure.
	CategoryRemediation SuggestionCategory = "remediation"
	// CategoryIntent indicates a semantic alias or alternative intent (e.g. docker -> cluster).
	CategoryIntent SuggestionCategory = "intent"
	// CategoryResource indicates a repository, branch, or node target suggestion.
	CategoryResource SuggestionCategory = "resource"
)

// ActionType defines the execution mode of the suggested command.
type ActionType string

const (
	// ActionAutoRun indicates the recommendation is safe to auto-execute.
	ActionAutoRun ActionType = "auto-run"
	// ActionManual indicates the recommendation requires user parameter adjustment.
	ActionManual ActionType = "manual"
	// ActionCopy indicates a copy-paste command candidate.
	ActionCopy ActionType = "copy"
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

// RenderFormat defines the rendering style for suggestions.
type RenderFormat string

const (
	// RenderFormatBox renders an elegant Catppuccin terminal box.
	RenderFormatBox RenderFormat = "box"
	// RenderFormatCompact renders a minimal single-line inline hint.
	RenderFormatCompact RenderFormat = "compact"
	// RenderFormatJSON renders machine-readable JSON for agents and tooling.
	RenderFormatJSON RenderFormat = "json"
)

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
