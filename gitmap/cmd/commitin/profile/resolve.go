package profile

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Resolved is the final flattened settings after applying the load
// order from spec §5.6 (CLI > --profile > --default > defaults).
type Resolved struct {
	ConflictMode     string
	Author           *Author
	Exclusions       []Exclusion
	MessageRules     []MessageRule
	MessagePrefix    []string
	MessageSuffix    []string
	TitlePrefix      string
	TitleSuffix      string
	OverrideMessages []string
	OverrideOnlyWeak bool
	WeakWords        []string
	FunctionIntel    FunctionIntel
}

// CliOverrides represents flag-level overrides; nil-pointer fields
// mean "user did not pass this flag" so the next layer wins.
type CliOverrides struct {
	ConflictMode     *string
	Author           *Author
	Exclusions       []Exclusion // empty slice = no override
	MessageRules     []MessageRule
	MessagePrefix    []string
	MessageSuffix    []string
	TitlePrefix      *string
	TitleSuffix      *string
	OverrideMessages []string
	OverrideOnlyWeak *bool
	WeakWords        []string
	FunctionIntel    *FunctionIntel
}

// Resolve applies the four-layer precedence: CLI > profile > defaults.
// `prof` may be nil (no --profile / --default). All defaults match
// spec §02 flag-default column.
func Resolve(cli *CliOverrides, prof *Profile) Resolved {
	r := defaultResolved()
	if prof != nil {
		applyProfile(&r, prof)
	}
	if cli != nil {
		applyCli(&r, cli)
	}
	return r
}

func defaultResolved() Resolved {
	return Resolved{
		ConflictMode: constants.CommitInDefaultConflictMode,
		WeakWords:    []string{"change", "update", "updates"},
		FunctionIntel: FunctionIntel{
			IsEnabled: false,
			Languages: []string{constants.CommitInLanguageGo},
		},
	}
}

func applyProfile(r *Resolved, p *Profile) {
	if p.ConflictMode != "" {
		r.ConflictMode = p.ConflictMode
	}
	if p.Author != nil {
		r.Author = p.Author
	}
	r.Exclusions = p.Exclusions
	r.MessageRules = p.MessageRules
	r.MessagePrefix = p.MessagePrefix
	r.MessageSuffix = p.MessageSuffix
	r.TitlePrefix = p.TitlePrefix
	r.TitleSuffix = p.TitleSuffix
	r.OverrideMessages = p.OverrideMessages
	r.OverrideOnlyWeak = p.OverrideOnlyWeak
	if len(p.WeakWords) > 0 {
		r.WeakWords = p.WeakWords
	}
	r.FunctionIntel = p.FunctionIntel
}

func applyCli(r *Resolved, c *CliOverrides) {
	if c.ConflictMode != nil {
		r.ConflictMode = *c.ConflictMode
	}
	if c.Author != nil {
		r.Author = c.Author
	}
	if c.Exclusions != nil {
		r.Exclusions = c.Exclusions
	}
	if c.MessageRules != nil {
		r.MessageRules = c.MessageRules
	}
	if c.MessagePrefix != nil {
		r.MessagePrefix = c.MessagePrefix
	}
	if c.MessageSuffix != nil {
		r.MessageSuffix = c.MessageSuffix
	}
	if c.TitlePrefix != nil {
		r.TitlePrefix = *c.TitlePrefix
	}
	if c.TitleSuffix != nil {
		r.TitleSuffix = *c.TitleSuffix
	}
	if c.OverrideMessages != nil {
		r.OverrideMessages = c.OverrideMessages
	}
	if c.OverrideOnlyWeak != nil {
		r.OverrideOnlyWeak = *c.OverrideOnlyWeak
	}
	if c.WeakWords != nil {
		r.WeakWords = c.WeakWords
	}
	if c.FunctionIntel != nil {
		r.FunctionIntel = *c.FunctionIntel
	}
}
