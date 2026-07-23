package commitin

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// validateConflictMode rejects anything outside the §2.5 enum. Empty
// is allowed (means "use default / profile / prompt") and short-
// circuits before this validator runs.
func validateConflictMode(value string) *ParseError {
	switch value {
	case "",
		constants.CommitInConflictModeForceMerge,
		constants.CommitInConflictModePrompt:
		return nil
	}
	return newBadArgs(constants.CommitInErrConflictMode, value)
}

// validateFunctionIntelToggle constrains --function-intel to on|off.
func validateFunctionIntelToggle(value string) *ParseError {
	switch value {
	case "", constants.CommitInFunctionIntelOn, constants.CommitInFunctionIntelOff:
		return nil
	}
	return newBadArgs(constants.CommitInErrFunctionIntelArg, value)
}

// validateLanguages rejects any token that is not a member of the
// FunctionIntelLanguage enum. Spec-aligned via AllLanguages().
func validateLanguages(values []string) *ParseError {
	if len(values) == 0 {
		return nil
	}
	known := languageSet()
	for _, v := range values {
		if !known[v] {
			return newBadArgs(constants.CommitInErrUnknownLanguage, v)
		}
	}
	return nil
}

// languageSet builds the lookup set from the typed enum so adding a
// language only requires editing enums.go.
func languageSet() map[string]bool {
	out := make(map[string]bool, len(AllLanguages()))
	for _, l := range AllLanguages() {
		out[l.String()] = true
	}
	return out
}

// validateAuthorPair enforces "both or neither" per §2.5.
func validateAuthorPair(name, email string) *ParseError {
	hasName := name != ""
	hasEmail := email != ""
	if hasName == hasEmail {
		return nil
	}
	return newBadArgs("%s", constants.CommitInErrAuthorPair)
}

// parseMessageRules turns CSV `Kind:Value` items into typed rules.
// Empty CSV is fine (no rules). Unknown Kind is a parse error.
func parseMessageRules(values []string) ([]MessageRuleArg, *ParseError) {
	if len(values) == 0 {
		return nil, nil
	}
	known := messageRuleKindSet()
	out := make([]MessageRuleArg, 0, len(values))
	for _, raw := range values {
		kind, value, ok := strings.Cut(raw, constants.CommitInMessageRuleKindSep)
		if !ok || !known[kind] || value == "" {
			return nil, newBadArgs(constants.CommitInErrMessageRuleShape, raw)
		}
		out = append(out, MessageRuleArg{Kind: kind, Value: value})
	}
	return out, nil
}

func messageRuleKindSet() map[string]bool {
	out := make(map[string]bool, len(AllMessageRuleKinds()))
	for _, k := range AllMessageRuleKinds() {
		out[k.String()] = true
	}
	return out
}

// requireSourceAndInputs enforces the §2.2 minimum: one <source> +
// at least one input (explicit or keyword).
func requireSourceAndInputs(source string, inputs []string, keyword string) *ParseError {
	if source == "" {
		return newBadArgs("%s", "missing <source>")
	}
	if len(inputs) == 0 && keyword == "" {
		return newBadArgs("%s", "missing <inputs...>")
	}
	return nil
}

// rejectMixedKeyword guards §2.4: a KEYWORD must appear alone.
func rejectMixedKeyword(keyword string, explicitInputs []string) *ParseError {
	if keyword == "" || len(explicitInputs) == 0 {
		return nil
	}
	return &ParseError{
		ExitCode: constants.CommitInExitBadArgs,
		Message:  fmt.Sprintf(constants.CommitInErrInputMixedKeyword, keyword),
	}
}
