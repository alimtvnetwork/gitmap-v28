package commitin

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Each enumExpect locks the spec's authoritative member list (in order)
// to the constants tokens. Adding/removing/reordering a member must
// update the spec, the constants_commitin.go enum block, the typed
// enum in enums.go, AND this table — in the same commit.
type enumExpect struct {
	name string
	got  []string
	want []string
}

func TestCommitInEnumsMatchSpec(t *testing.T) {
	cases := []enumExpect{
		{
			name: "ConflictMode",
			got:  toStrings(AllConflictModes()),
			want: []string{
				constants.CommitInConflictModeForceMerge,
				constants.CommitInConflictModePrompt,
			},
		},
		{
			name: "InputKind",
			got:  toStrings(AllInputKinds()),
			want: []string{
				constants.CommitInInputKindLocalFolder,
				constants.CommitInInputKindGitUrl,
				constants.CommitInInputKindVersionedSibling,
			},
		},
		{
			name: "RunStatus",
			got:  toStrings(AllRunStatuses()),
			want: []string{
				constants.CommitInRunStatusPending,
				constants.CommitInRunStatusRunning,
				constants.CommitInRunStatusCompleted,
				constants.CommitInRunStatusFailed,
				constants.CommitInRunStatusPartiallyFailed,
			},
		},
		{
			name: "CommitOutcome",
			got:  toStrings(AllCommitOutcomes()),
			want: []string{
				constants.CommitInOutcomeCreated,
				constants.CommitInOutcomeSkipped,
				constants.CommitInOutcomeFailed,
			},
		},
		{
			name: "SkipReason",
			got:  toStrings(AllSkipReasons()),
			want: []string{
				constants.CommitInSkipReasonDuplicateSourceSha,
				constants.CommitInSkipReasonExcludedAllFiles,
				constants.CommitInSkipReasonEmptyAfterMessageRules,
				constants.CommitInSkipReasonDryRun,
			},
		},
		{
			name: "ExclusionKind",
			got:  toStrings(AllExclusionKinds()),
			want: []string{
				constants.CommitInExclusionKindPathFolder,
				constants.CommitInExclusionKindPathFile,
			},
		},
		{
			name: "MessageRuleKind",
			got:  toStrings(AllMessageRuleKinds()),
			want: []string{
				constants.CommitInMessageRuleKindStartsWith,
				constants.CommitInMessageRuleKindEndsWith,
				constants.CommitInMessageRuleKindContains,
			},
		},
		{
			name: "FunctionIntelLanguage",
			got:  toStrings(AllLanguages()),
			want: []string{
				constants.CommitInLanguageGo,
				constants.CommitInLanguageJavaScript,
				constants.CommitInLanguageTypeScript,
				constants.CommitInLanguageRust,
				constants.CommitInLanguagePython,
				constants.CommitInLanguagePhp,
				constants.CommitInLanguageJava,
				constants.CommitInLanguageCSharp,
			},
		},
	}
	for _, tc := range cases {
		if !equalSlices(tc.got, tc.want) {
			t.Errorf("enum %s mismatch:\n  got : %v\n  want: %v\n  fix : update spec/03-commit-in/04-database-schema.md, constants_commitin.go, and enums.go together",
				tc.name, tc.got, tc.want)
		}
		// PascalCase guard: every member must be non-empty and start
		// with an uppercase ASCII letter (matches DB-schema rule that
		// enum-mirror Names are PascalCase).
		for _, m := range tc.got {
			if m == "" || m[0] < 'A' || m[0] > 'Z' {
				t.Errorf("enum %s member %q must be PascalCase", tc.name, m)
			}
		}
	}
}

// Exit codes must be unique to keep error attribution unambiguous. Any
// duplicate would silently collapse two distinct failure modes into one.
func TestCommitInExitCodesUnique(t *testing.T) {
	codes := map[string]int{
		"Ok":              constants.CommitInExitOk,
		"PartiallyFailed": constants.CommitInExitPartiallyFailed,
		"BadArgs":         constants.CommitInExitBadArgs,
		"SourceUnusable":  constants.CommitInExitSourceUnusable,
		"InputUnusable":   constants.CommitInExitInputUnusable,
		"DbFailed":        constants.CommitInExitDbFailed,
		"ProfileMissing":  constants.CommitInExitProfileMissing,
		"MissingAnswer":   constants.CommitInExitMissingAnswer,
		"ConflictAborted": constants.CommitInExitConflictAborted,
		"LockBusy":        constants.CommitInExitLockBusy,
		"FunctionIntel":   constants.CommitInExitFunctionIntel,
	}
	seen := map[int]string{}
	for name, code := range codes {
		if other, dup := seen[code]; dup {
			t.Errorf("commit-in exit code %d collides: %s and %s", code, other, name)
		}
		seen[code] = name
	}
}

// Flag long names must be unique and kebab-case (no underscores, no
// uppercase). Spec §2.5 lists every flag once.
func TestCommitInFlagNamesShape(t *testing.T) {
	flags := []string{
		constants.CommitInFlagDefault,
		constants.CommitInFlagProfile,
		constants.CommitInFlagSaveProfile,
		constants.CommitInFlagSaveProfileOverwrite,
		constants.CommitInFlagSetDefault,
		constants.CommitInFlagAuthorName,
		constants.CommitInFlagAuthorEmail,
		constants.CommitInFlagConflict,
		constants.CommitInFlagExclude,
		constants.CommitInFlagMessageExclude,
		constants.CommitInFlagMessagePrefix,
		constants.CommitInFlagMessageSuffix,
		constants.CommitInFlagTitlePrefix,
		constants.CommitInFlagTitleSuffix,
		constants.CommitInFlagOverrideMessages,
		constants.CommitInFlagOverrideOnlyWeak,
		constants.CommitInFlagWeakWords,
		constants.CommitInFlagFunctionIntel,
		constants.CommitInFlagLanguages,
		constants.CommitInFlagNoPrompt,
		constants.CommitInFlagDryRun,
		constants.CommitInFlagKeepTemp,
	}
	seen := map[string]bool{}
	for _, f := range flags {
		if f == "" {
			t.Error("commit-in flag name must not be empty")
			continue
		}
		if strings.ContainsAny(f, "_ ") || strings.ToLower(f) != f {
			t.Errorf("commit-in flag %q must be lowercase kebab-case", f)
		}
		if seen[f] {
			t.Errorf("commit-in flag %q is duplicated", f)
		}
		seen[f] = true
	}
}

// ---- helpers ------------------------------------------------------

type stringer interface{ String() string }

func toStrings[T stringer](items []T) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.String()
	}
	return out
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
