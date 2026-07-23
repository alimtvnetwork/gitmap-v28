package cmd

import "github.com/alimtvnetwork/gitmap-v27/gitmap/constants"

// CfrModifierFlags holds the parsed pre-URL modifier tokens accepted by
// `cfr` and `cfrp`. Tokens are order-independent and may appear anywhere
// before the first non-modifier positional (the repo URL or a flag).
//
// Modifiers currently recognised:
//
//	cg  → install alimtvnetwork coding-guidelines (v24) after clone/fix.
//	p   → promote to public (equivalent to using `cfrp`); no-op when the
//	      caller is already `cfrp`.
//
// Anything that is not a known modifier ends the scan — the remaining
// argv is handed to the existing flag/positional parser untouched. This
// keeps parsing conservative: an unknown leading positional (typo,
// future flag, URL that happens to start with a letter) is never
// silently swallowed.
type CfrModifierFlags struct {
	InstallCodingGuidelines bool
	PromotePublic           bool
	// NoCommit / NoPush are populated from `--no-commit` / `--no-push`
	// flags parsed downstream in parseCloneFixRepoArgs. They only take
	// effect when InstallCodingGuidelines is true; otherwise there is
	// no auto-commit step to skip.
	NoCommit bool
	NoPush   bool
}


// ParseCfrModifiers walks args from the front, consuming known modifier
// tokens (`cg`, `p`) in any order, and returns the parsed flags plus the
// remaining argv. Parsing stops at the first token that is either a flag
// (leading `-`) or an unrecognised positional so URLs and existing flags
// pass through untouched.
//
// Examples (rest shown as a Go slice literal):
//
//	ParseCfrModifiers([]string{"cg", "https://x/y"})       → {CG:true},         ["https://x/y"]
//	ParseCfrModifiers([]string{"p", "cg", "https://x/y"})  → {CG:true,Pub:true},["https://x/y"]
//	ParseCfrModifiers([]string{"cg", "p", "https://x/y"})  → {CG:true,Pub:true},["https://x/y"]
//	ParseCfrModifiers([]string{"--ssh", "cg", "url"})      → {},                ["--ssh","cg","url"]
//	ParseCfrModifiers([]string{"https://x/y", "cg"})       → {},                ["https://x/y","cg"]
func ParseCfrModifiers(args []string) (CfrModifierFlags, []string) {
	var flags CfrModifierFlags
	i := 0
	for i < len(args) {
		if !isCfrModifierToken(args[i]) {
			break
		}
		applyCfrModifier(&flags, args[i])
		i++
	}

	return flags, args[i:]
}

// isCfrModifierToken reports whether tok is one of the recognised
// pre-URL modifier tokens.
func isCfrModifierToken(tok string) bool {
	switch tok {
	case constants.CfrModifierCodingGuidelines, constants.CfrModifierPublic:
		return true
	}

	return false
}

// applyCfrModifier sets the matching field on flags. Idempotent: repeated
// tokens (`cg cg`) are accepted silently rather than erroring, matching
// the tolerant style of the surrounding CLI.
func applyCfrModifier(flags *CfrModifierFlags, tok string) {
	switch tok {
	case constants.CfrModifierCodingGuidelines:
		flags.InstallCodingGuidelines = true
	case constants.CfrModifierPublic:
		flags.PromotePublic = true
	}
}
