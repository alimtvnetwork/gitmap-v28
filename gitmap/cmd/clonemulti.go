package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// urlListSeparators are the characters that split a single positional
// arg into multiple URLs. Comma is the documented form; semicolon is
// accepted because (a) bash users reach for it naturally and (b) when
// PowerShell sees an unquoted `;` it terminates the statement, so by
// the time we receive a `;`-bearing token the user almost certainly
// meant a list (single quoted form: `clone 'a;b;c'`).
const urlListSeparators = ",;"

// flattenURLArgs splits each positional arg on commas/semicolons,
// sanitizes each segment (BOM, smart quotes, zero-width chars,
// surrounding whitespace), drops empties, and returns the ordered,
// deduplicated list of URLs. Both space- and list-separated forms
// are accepted, mixable:
//
//	gitmap clone a b c          → [a, b, c]
//	gitmap clone a,b,c          → [a, b, c]
//	gitmap clone "a;b;c"        → [a, b, c]   (PowerShell quoted)
//	gitmap clone a,b c d;e      → [a, b, c, d, e]
//
// Dedup is case-insensitive with trailing ".git" normalised; first-seen wins.
// See: spec/01-app/104-clone-multi.md and mem://features/clone-multi.
func flattenURLArgs(args []string) []string {
	out := make([]string, 0, len(args))
	seen := make(map[string]struct{}, len(args))

	for _, raw := range args {
		for _, part := range splitOnURLSeparators(raw) {
			cleaned := sanitizeURLToken(part)
			if cleaned == "" {
				continue
			}

			key := normaliseURLKey(cleaned)
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, cleaned)
		}
	}

	return out
}

// splitOnURLSeparators splits on every rune in urlListSeparators.
// We deliberately do not use strings.Split because we want both
// `,` and `;` to act as boundaries simultaneously.
func splitOnURLSeparators(raw string) []string {
	return strings.FieldsFunc(raw, func(r rune) bool {
		return strings.ContainsRune(urlListSeparators, r)
	})
}

// sanitizeURLToken removes characters that survive copy-paste from
// browsers, docs, and terminals but break URL parsing downstream:
//   - U+FEFF (BOM) — Windows clipboard frequently injects this
//   - U+200B…U+200D (zero-width spaces) — copied from rich-text sources
//   - Smart quotes (U+2018, U+2019, U+201C, U+201D) — Word/Slack auto-fix
//   - Surrounding ASCII whitespace and matched ASCII quotes/backticks
//
// A token that is nothing but separators / quotes / whitespace
// returns "" so the caller drops it instead of producing a bogus
// "invalid URL" warning for what was really a typo.
func sanitizeURLToken(s string) string {
	cleaned := stripInvisibleRunes(s)
	cleaned = replaceSmartQuotes(cleaned)
	cleaned = strings.TrimSpace(cleaned)
	cleaned = trimMatchingWrappers(cleaned)
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.Trim(cleaned, urlListSeparators)

	return strings.TrimSpace(cleaned)
}

// stripInvisibleRunes drops BOM and zero-width characters anywhere
// in the string. These never belong in a URL.
func stripInvisibleRunes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '\uFEFF', '\u200B', '\u200C', '\u200D':
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}

// replaceSmartQuotes folds curly quotes back to ASCII so that
// trimMatchingWrappers can strip them in one pass.
func replaceSmartQuotes(s string) string {
	r := strings.NewReplacer(
		"\u2018", "'", "\u2019", "'",
		"\u201C", "\"", "\u201D", "\"",
	)

	return r.Replace(s)
}

// trimMatchingWrappers strips one matched pair of `'`, `"`, or backticks
// surrounding the token. Only matched pairs are stripped — a stray
// trailing quote stays, so the caller still sees a recognizably broken
// URL rather than a silently "fixed" one.
func trimMatchingWrappers(s string) string {
	if len(s) < 2 {
		return s
	}
	first, last := s[0], s[len(s)-1]
	if first == last && (first == '\'' || first == '"' || first == '`') {
		return s[1 : len(s)-1]
	}

	return s
}

// normaliseURLKey lowercases and trims a trailing ".git" so that
// "https://x/y" and "HTTPS://X/Y.git" dedupe to the same key.
func normaliseURLKey(url string) string {
	lower := strings.ToLower(url)
	lower = strings.TrimSuffix(lower, ".git")

	return lower
}

// classifyURLs partitions a flattened arg list into valid URLs and
// invalid (non-URL, non-empty) entries. Invalid entries are reported
// to stderr but do NOT abort the batch — the caller decides exit code.
func classifyURLs(flat []string) (valid, invalid []string) {
	for _, candidate := range flat {
		if isDirectURL(candidate) {
			valid = append(valid, candidate)

			continue
		}

		invalid = append(invalid, candidate)
		fmt.Fprintf(os.Stderr, constants.MsgCloneInvalidURLFmt, candidate)
	}

	return valid, invalid
}

// executeDirectCloneOne is the non-fatal sibling of executeDirectClone:
// it clones a single URL, persists to the DB, optionally registers with
// GitHub Desktop, and returns any error instead of calling os.Exit.
// Folder name is derived verbatim from the URL (see resolveCloneFolder).
func executeDirectCloneOne(url, folderName string, ghDesktopFlag, noReplace bool) error {
	repoName := repoNameFromURL(url)
	folderName = resolveCloneFolder(repoName, folderName)

	absPath, err := filepath.Abs(folderName)
	if err != nil {
		return fmt.Errorf("resolve abs path for %s: %w", folderName, err)
	}

	url = coerceURLToStoredTransport(url)

	errClone := executeCloneStrategy(url, absPath, noReplace)
	if errClone != nil {
		return errClone
	}

	persistRecloneTransport(url)
	upsertDirectClone(url, repoName, folderName, absPath)

	if ghDesktopFlag {
		registerSingleDesktop(repoName, absPath)
		fmt.Printf(constants.MsgCloneRegisteredInline, repoName)
	}

	return nil
}

// resolveCloneFolder returns the destination folder name for a clone.
//
// Base `clone` NEVER rewrites the folder name. An explicit folder
// argument wins; otherwise the folder is the repo name exactly as it
// appears in the URL, including any trailing `-vN` suffix
// (wp-onboarding-v13 stays wp-onboarding-v13/). Version flattening and
// version bumping belong to `gitmap clone-next`, not here. This mirrors
// the single-URL path in executeDirectClone so one URL and many URLs
// always produce the same on-disk name.
func resolveCloneFolder(repoName, folderName string) string {
	if len(folderName) > 0 {
		return folderName
	}

	return repoName
}

func executeCloneStrategy(url string, absPath string, noReplace bool) error {
	if noReplace {
		return handleNoReplaceClone(url, absPath)
	}
	return handleReplaceClone(url, absPath)
}

func handleNoReplaceClone(url string, absPath string) error {
	_, statErr := os.Stat(absPath)
	if statErr == nil {
		return fmt.Errorf("target exists: %s (use without --no-replace to replace)", absPath)
	}
	cloneErr := runCloneCommand(url, absPath)
	if cloneErr != nil {
		return fmt.Errorf("git clone: %w", cloneErr)
	}
	return nil
}

func handleReplaceClone(url string, absPath string) error {
	_, replaceErr := cloneReplacing(url, absPath)
	if replaceErr != nil {
		return fmt.Errorf("clone-replace: %w", replaceErr)
	}
	return nil
}
