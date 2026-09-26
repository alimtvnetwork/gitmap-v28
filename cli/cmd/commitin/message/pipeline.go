package message

import "strings"

// Build runs the §6.1 pipeline in canonical order.
func Build(in Inputs) Result {
	msg := stripRules(in.OriginalMessage, in.Resolved.MessageRules)
	msg = applyOverride(msg, in)
	chosenPrefix, chosenSuffix := pickAffixPair(in.Resolved.MessagePrefix, in.Resolved.MessageSuffix, in.PickIndex)
	seoTitle := resolveChosenTemplateTitle(chosenPrefix, chosenSuffix)
	msg = applyTitleReplacements(msg, in.Files, in.Resolved.TitleReplacements, seoTitle)
	msg = applyTitleAffix(msg, in.Resolved.TitlePrefix, in.Resolved.TitleSuffix)
	msg = expandTitleLineVars(msg, in.Files, seoTitle)
	msg = applyChosenBodyAffix(msg, chosenPrefix, chosenSuffix)
	msg = appendFunctionIntel(msg, in.FunctionIntel)

	return Result{Message: msg, IsEmpty: isEmpty(msg)}
}

func resolveChosenTemplateTitle(chosenPrefix, chosenSuffix string) string {
	if h := extractTemplateHeading(chosenSuffix); h != "" {
		return h
	}

	return extractTemplateHeading(chosenPrefix)
}

func expandTitleLineVars(msg string, files []string, seoTitle string) string {
	if !strings.Contains(msg, "$") {
		return msg
	}
	title, rest := splitTitleAndBody(msg)

	return expandFileAndTemplateVars(title, files, seoTitle) + rest
}

func applyOverride(msg string, in Inputs) string {
	pool := in.Resolved.OverrideMessages
	if len(pool) == 0 {
		return msg
	}

	if in.Resolved.OverrideOnlyWeak && !matchesWeak(msg, in.Resolved.WeakWords) {
		return msg
	}

	return pickOne(pool, in.PickIndex)
}

func appendFunctionIntel(msg, block string) string {
	if block == "" {
		return msg
	}

	if msg == "" {
		return block
	}

	return msg + "\n\n" + block
}

func isEmpty(msg string) bool { return strings.TrimSpace(msg) == "" }
