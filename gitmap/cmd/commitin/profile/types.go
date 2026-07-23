package profile

// Author identity bound to a profile (both fields must be set together).
type Author struct {
	Name  string
	Email string
}

// Exclusion represents one PathFolder/PathFile rule.
type Exclusion struct {
	Kind  string
	Value string
}

// MessageRule represents one StartsWith/EndsWith/Contains line-strip rule.
type MessageRule struct {
	Kind  string
	Value string
}

// FunctionIntel block of a profile.
type FunctionIntel struct {
	IsEnabled bool
	Languages []string
}

// Profile is the canonical in-memory representation of a commit-in
// profile (matches spec §5.2 byte-for-byte when re-serialized).
type Profile struct {
	Name             string
	SchemaVersion    int
	SourceRepoPath   string
	IsDefault        bool
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

// CurrentSchemaVersion is the only SchemaVersion v1 loaders accept.
const CurrentSchemaVersion = 1

// LoadError signals a profile load/decode failure with stable .Reason
// strings the caller can map to exit codes.
type LoadError struct {
	Path   string
	Reason string
	Cause  error
}

func (e *LoadError) Error() string {
	if e.Cause == nil {
		return e.Reason
	}
	return e.Reason + ": " + e.Cause.Error()
}

func (e *LoadError) Unwrap() error { return e.Cause }
