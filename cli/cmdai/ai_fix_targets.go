package cmdai

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var fixTargetDefs = []FixTargetDef{
	{
		Target:      FixTargetEncoding,
		ScriptToken: "10",
		DefaultArgs: []string{"--fix"},
		Description: "Normalizes files to strict UTF-8 with UNIX LF line endings",
	},
	{
		Target:      FixTargetNewlines,
		ScriptToken: "04",
		DefaultArgs: []string{"--fix"},
		Description: "Fixes trailing whitespace and missing final newlines",
	},
	{
		Target:      FixTargetNaming,
		ScriptToken: "08",
		DefaultArgs: []string{"--fix"},
		Description: "Enforces lowercase filenames and affirmative boolean naming",
	},
	{
		Target:      FixTargetGuidelines,
		ScriptToken: "05",
		DefaultArgs: []string{"--fix"},
		Description: "Composite autofixer for guideline and boolean conventions",
	},
	{
		Target:      FixTargetPaths,
		ScriptToken: "07",
		DefaultArgs: []string{"--fix"},
		Description: "Detects and fixes absolute paths and file URIs",
	},
	{
		Target:      FixTargetSpelling,
		ScriptToken: "27b",
		DefaultArgs: []string{"--fix"},
		Description: "Audits and auto-fixes British to American English spelling",
	},
	{
		Target:      FixTargetPlans,
		ScriptToken: "20",
		DefaultArgs: []string{"--fix"},
		Description: "Consolidates and normalizes plan files and subtasks",
	},
}

// AllFixTargetDefs returns all registered individual fix targets in execution order.
func AllFixTargetDefs() []FixTargetDef {
	return fixTargetDefs
}

// ResolveFixTarget matches a target name or alias to its FixTargetDef.
func ResolveFixTarget(token string) (FixTargetDef, *apperror.AppError) {
	norm := normalizeFixToken(token)
	isAll := norm == string(FixTargetAll)
	if isAll {
		return compositeFixTargetDef(), nil
	}

	for _, def := range fixTargetDefs {
		isMatch := string(def.Target) == norm
		if isMatch {
			return def, nil
		}
	}

	ctx := map[string]any{"target": token, "normalized": norm}

	return FixTargetDef{}, apperror.New("resolve_target", "E_UNKNOWN_FIX_TARGET", ctx)
}

func compositeFixTargetDef() FixTargetDef {
	return FixTargetDef{
		Target:      FixTargetAll,
		ScriptToken: "all",
		DefaultArgs: []string{"--fix"},
		Description: "Composite repository autofix pipeline",
	}
}

func normalizeFixToken(token string) string {
	clean := strings.ToLower(strings.TrimSpace(token))
	isAll := clean == "all" || clean == "everything" || clean == ""
	if isAll {
		return string(FixTargetAll)
	}

	return resolveTokenAlias(clean)
}

func resolveTokenAlias(clean string) string {
	isGuidelines := clean == "guidelines" || clean == "guideline" || clean == "cg"
	if isGuidelines {
		return string(FixTargetGuidelines)
	}

	isNewlines := clean == "newlines" || clean == "newline" || clean == "whitespace" || clean == "ws"
	if isNewlines {
		return string(FixTargetNewlines)
	}

	isPaths := clean == "paths" || clean == "path" || clean == "relative-paths"
	if isPaths {
		return string(FixTargetPaths)
	}

	isNaming := clean == "naming" || clean == "name" || clean == "booleans" || clean == "bool"
	if isNaming {
		return string(FixTargetNaming)
	}

	return resolveTokenAliasExtended(clean)
}

func resolveTokenAliasExtended(clean string) string {
	isEncoding := clean == "encoding" || clean == "utf8" || clean == "lf"
	if isEncoding {
		return string(FixTargetEncoding)
	}

	isSpelling := clean == "spelling" || clean == "spell" || clean == "misspell"
	if isSpelling {
		return string(FixTargetSpelling)
	}

	isPlans := clean == "plans" || clean == "plan"
	if isPlans {
		return string(FixTargetPlans)
	}

	return clean
}

// MergeFixArgs merges target default arguments with extra user-supplied arguments.
func MergeFixArgs(def FixTargetDef, extraArgs []string) []string {
	merged := append([]string{}, def.DefaultArgs...)
	for _, arg := range extraArgs {
		isFixArg := arg == "--fix"
		if isFixArg {
			continue
		}

		merged = append(merged, arg)
	}

	return merged
}
