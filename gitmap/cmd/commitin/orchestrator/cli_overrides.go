package orchestrator

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// buildCliOverrides projects the parser's RawArgs onto the
// profile.CliOverrides shape so profile.Resolve can apply layered
// precedence. nil-pointer / nil-slice = "user did not pass this flag".
func buildCliOverrides(raw *commitin.RawArgs) *profile.CliOverrides {
	c := &profile.CliOverrides{}
	addAuthor(c, raw)
	addConflictAndAffixes(c, raw)
	addOverridesAndIntel(c, raw)
	return c
}

func addAuthor(c *profile.CliOverrides, raw *commitin.RawArgs) {
	if raw.AuthorName != "" || raw.AuthorEmail != "" {
		c.Author = &profile.Author{Name: raw.AuthorName, Email: raw.AuthorEmail}
	}
}

func addConflictAndAffixes(c *profile.CliOverrides, raw *commitin.RawArgs) {
	if raw.ConflictMode != "" {
		s := raw.ConflictMode
		c.ConflictMode = &s
	}
	if raw.TitlePrefix != "" {
		s := raw.TitlePrefix
		c.TitlePrefix = &s
	}
	if raw.TitleSuffix != "" {
		s := raw.TitleSuffix
		c.TitleSuffix = &s
	}
	if len(raw.MessagePrefix) > 0 {
		c.MessagePrefix = raw.MessagePrefix
	}
	if len(raw.MessageSuffix) > 0 {
		c.MessageSuffix = raw.MessageSuffix
	}
}

func addOverridesAndIntel(c *profile.CliOverrides, raw *commitin.RawArgs) {
	if len(raw.OverrideMessages) > 0 {
		c.OverrideMessages = raw.OverrideMessages
	}
	if raw.OverrideOnlyWeak {
		b := true
		c.OverrideOnlyWeak = &b
	}
	if len(raw.WeakWords) > 0 {
		c.WeakWords = raw.WeakWords
	}
	if raw.FunctionIntel != "" {
		fi := profile.FunctionIntel{
			IsEnabled: raw.FunctionIntel == constants.CommitInFunctionIntelOn,
			Languages: raw.Languages,
		}
		c.FunctionIntel = &fi
	}
	if len(raw.MessageRules) > 0 {
		c.MessageRules = mapMessageRules(raw.MessageRules)
	}
}

func mapMessageRules(in []commitin.MessageRuleArg) []profile.MessageRule {
	out := make([]profile.MessageRule, 0, len(in))
	for _, r := range in {
		out = append(out, profile.MessageRule{Kind: r.Kind, Value: r.Value})
	}
	return out
}
