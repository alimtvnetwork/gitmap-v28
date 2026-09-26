package commitin

import (
	"os"
	"strings"
)

// Parse converts an argv slice (already stripped of the leading
// `commit-in` / `cin` token) into a fully-validated RawArgs.
func Parse(args []string) (*RawArgs, *ParseError) {
	fs, raw, csv := newFlagSet()
	if err := fs.Parse(reorder(args)); err != nil {
		return nil, newBadArgs("%v", err)
	}

	if perr := finalizeFlagFanout(raw, csv); perr != nil {
		return nil, perr
	}

	if perr := splitPositionalOrConfig(raw, fs.Args()); perr != nil {
		return nil, perr
	}

	if perr := applyConfigFileIfPresent(raw); perr != nil {
		return nil, perr
	}

	if perr := validateAll(raw); perr != nil {
		return nil, perr
	}

	return raw, nil
}

// finalizeFlagFanout splits every CSV holder into its typed slice and
// runs the message-rule shape validator.
func finalizeFlagFanout(raw *RawArgs, csv *csvHolder) *ParseError {
	raw.Exclude = splitCSV(csv.exclude)
	raw.MessagePrefix = splitCSV(csv.messagePrefix)
	raw.MessageSuffix = splitCSV(csv.messageSuffix)
	raw.OverrideMessages = splitCSV(csv.overrideMessages)
	raw.WeakWords = splitCSV(csv.weakWords)
	raw.Languages = splitCSV(csv.languages)
	rules, perr := parseMessageRules(splitCSV(csv.messageExclude))
	if perr != nil {
		return perr
	}

	raw.MessageRules = rules
	maybeLoadStateSEOTemplates(raw)

	return nil
}

func maybeLoadStateSEOTemplates(raw *RawArgs) {
	cat := raw.SEOTemplate
	if cat == "" && raw.IsSponsor {
		cat = "seo"
	}
	if cat == "" || PrecompileStateCategoryHook == nil {
		return
	}
	if compiled := PrecompileStateCategoryHook(cat, nil); len(compiled) > 0 {
		raw.MessageSuffix = append(raw.MessageSuffix, compiled...)
	}
}

func splitPositionalOrConfig(raw *RawArgs, positional []string) *ParseError {
	if len(positional) == 0 {
		if raw.ConfigPath != "" {
			return nil
		}
		return newBadArgs("%s", "missing <source>")
	}
	if len(positional) == 1 && raw.ConfigPath == "" && isExistingJSONConfigFile(positional[0]) {
		raw.ConfigPath = positional[0]
		return nil
	}

	return splitPositional(raw, positional)
}

func isExistingJSONConfigFile(token string) bool {
	trimmed := strings.TrimSpace(token)
	if !strings.HasSuffix(strings.ToLower(trimmed), ".json") {
		return false
	}
	info, err := os.Stat(trimmed)

	return err == nil && !info.IsDir()
}

// splitPositional consumes the leftover argv: first token is <source>,
// remainder is the input list.
func splitPositional(raw *RawArgs, positional []string) *ParseError {
	if len(positional) == 0 {
		return newBadArgs("%s", "missing <source>")
	}

	raw.Source = positional[0]
	rest := positional[1:]
	if len(rest) == 1 {
		return applyKeywordArg(raw, rest)
	}

	raw.Inputs = splitInputs(rest)

	return nil
}

func applyKeywordArg(raw *RawArgs, rest []string) *ParseError {
	kw, tail, isKw, perr := classifyKeyword(rest[0])
	if perr != nil {
		return perr
	}

	if !isKw {
		raw.Inputs = splitInputs(rest)

		return nil
	}

	raw.Keyword = kw
	raw.KeywordTail = tail

	return nil
}

func validateAll(raw *RawArgs) *ParseError {
	if perr := requireSourceAndInputs(raw.Source, raw.Inputs, raw.Keyword); perr != nil {
		return perr
	}

	if perr := rejectMixedKeyword(raw.Keyword, raw.Inputs); perr != nil {
		return perr
	}

	if perr := validateAuthorPair(raw.AuthorName, raw.AuthorEmail); perr != nil {
		return perr
	}

	if perr := validateConflictMode(raw.ConflictMode); perr != nil {
		return perr
	}

	if perr := validateFunctionIntelToggle(raw.FunctionIntel); perr != nil {
		return perr
	}

	return validateLanguages(raw.Languages)
}

func reorder(args []string) []string {
	flags, positional := splitFlagsAndPositional(args)
	out := make([]string, 0, len(args))
	out = append(out, flags...)
	out = append(out, positional...)

	return out
}

func splitFlagsAndPositional(args []string) ([]string, []string) {
	bools := boolFlagSet()
	flags := make([]string, 0, len(args))
	positional := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		tok := args[i]
		if !strings.HasPrefix(tok, "-") || tok == "-" || isTailKeywordToken(tok) {
			positional = append(positional, tok)
			continue
		}

		flags = append(flags, tok)
		if needsValue(tok, bools) && i+1 < len(args) {
			flags = append(flags, args[i+1])
			i++
		}
	}

	return flags, positional
}

func isTailKeywordToken(tok string) bool {
	if len(tok) < 2 || tok[0] != '-' {
		return false
	}

	for i := 1; i < len(tok); i++ {
		if tok[i] < '0' || tok[i] > '9' {
			return false
		}
	}

	return true
}

func needsValue(tok string, bools map[string]bool) bool {
	if strings.Contains(tok, "=") {
		return false
	}

	name := strings.TrimLeft(tok, "-")

	return !bools[name]
}
