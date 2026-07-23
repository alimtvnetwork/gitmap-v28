package commitin

import "github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

// RawArgs is the in-memory shape of a successful `commit-in` parse.
// All resolved-but-unresolved-defaults values stay zero so the caller
// can layer profile + interactive prompts on top per spec §5.6.
//
// Pure value type: zero git, zero filesystem, zero DB. Trivially
// printable for golden-test diffing.
type RawArgs struct {
	Source string   // <source> argv token, verbatim
	Inputs []string // expanded input list (post separator/quote split)
	// Keyword captures the special-mode token when present. Empty
	// string when the user passed explicit inputs.
	Keyword     string // "all" | "-N" (verbatim) | ""
	KeywordTail int    // N when Keyword == "-N"; 0 otherwise

	// Behavior flags (defaults match spec §2.5 column "Default").
	UseDefaultProfile    bool
	ProfileName          string
	SaveProfileName      string
	SaveProfileOverwrite bool
	SetDefault           bool

	AuthorName  string
	AuthorEmail string

	ConflictMode  string // ConflictMode enum literal; empty = unset
	Exclude       []string
	MessageRules  []MessageRuleArg
	MessagePrefix []string
	MessageSuffix []string
	TitlePrefix   string
	TitleSuffix   string

	OverrideMessages []string
	OverrideOnlyWeak bool
	WeakWords        []string

	FunctionIntel string   // "on" | "off" | "" (unset)
	Languages     []string // FunctionIntelLanguage literals

	IsNoPrompt        bool
	IsDryRun          bool
	IsKeepTemp        bool
	IsNoReleaseBranch bool // --no-release-branch; default false (branches ON)
}

// MessageRuleArg is the parsed shape of one `--message-exclude` entry.
// `Kind` is one of the MessageRuleKind enum literals.
type MessageRuleArg struct {
	Kind  string
	Value string
}

// ParseError carries the exit code that the caller should propagate to
// `os.Exit` along with a stderr-ready message. Per zero-swallow rule,
// every error path through the parser returns exactly one ParseError.
type ParseError struct {
	ExitCode int
	Message  string
}

func (e *ParseError) Error() string { return e.Message }

func newBadArgs(format string, args ...any) *ParseError {
	return &ParseError{
		ExitCode: constants.CommitInExitBadArgs,
		Message:  sprintf(format, args...),
	}
}
