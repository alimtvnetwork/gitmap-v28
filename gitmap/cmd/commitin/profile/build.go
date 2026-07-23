package profile

// BuildArgs is the input to BuildFromResolved. Keeping it as a struct
// (rather than positional args) avoids a long argument list and makes
// future fields additive without breaking call sites.
type BuildArgs struct {
	Name           string
	SourceRepoPath string
	IsDefault      bool
	Resolved       Resolved
}

// BuildFromResolved materializes a Profile from a Resolved settings
// snapshot so the caller can hand it to SaveToDisk. The returned
// Profile is byte-stable when re-encoded (Encode applies canonical
// ordering + nil→empty-slice normalization).
func BuildFromResolved(args BuildArgs) *Profile {
	r := args.Resolved
	p := &Profile{
		Name:             args.Name,
		SchemaVersion:    CurrentSchemaVersion,
		SourceRepoPath:   args.SourceRepoPath,
		IsDefault:        args.IsDefault,
		ConflictMode:     r.ConflictMode,
		Author:           r.Author,
		Exclusions:       cloneExclusions(r.Exclusions),
		MessageRules:     cloneMessageRules(r.MessageRules),
		MessagePrefix:    cloneStrings(r.MessagePrefix),
		MessageSuffix:    cloneStrings(r.MessageSuffix),
		TitlePrefix:      r.TitlePrefix,
		TitleSuffix:      r.TitleSuffix,
		OverrideMessages: cloneStrings(r.OverrideMessages),
		OverrideOnlyWeak: r.OverrideOnlyWeak,
		WeakWords:        cloneStrings(r.WeakWords),
		FunctionIntel: FunctionIntel{
			IsEnabled: r.FunctionIntel.IsEnabled,
			Languages: cloneStrings(r.FunctionIntel.Languages),
		},
	}
	return p
}

func cloneStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func cloneExclusions(in []Exclusion) []Exclusion {
	if len(in) == 0 {
		return nil
	}
	out := make([]Exclusion, len(in))
	copy(out, in)
	return out
}

func cloneMessageRules(in []MessageRule) []MessageRule {
	if len(in) == 0 {
		return nil
	}
	out := make([]MessageRule, len(in))
	copy(out, in)
	return out
}
