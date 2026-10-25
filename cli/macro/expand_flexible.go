package macro

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// CleanFlexiblePath normalizes and resolves flexible path inputs.
func CleanFlexiblePath(input string) string {
	cleaned := strings.TrimSpace(input)
	if len(cleaned) == 0 {
		return ""
	}

	cleaned = stripLeadingBullets(cleaned)
	cleaned = stripTrailingComment(cleaned)
	cleaned = stripSurroundingQuotes(cleaned)
	cleaned = stripAtBrackets(cleaned)
	cleaned = extractMarkdownURL(cleaned)
	cleaned = stripFileScheme(cleaned)
	cleaned = unescapeURLPath(cleaned)
	cleaned = expandAllTokens(cleaned)

	return resolveFlexiblePath(cleaned)
}

func stripLeadingBullets(s string) string {
	hasDash := strings.HasPrefix(s, "- ")
	hasStar := strings.HasPrefix(s, "* ")
	hasPlus := strings.HasPrefix(s, "+ ")
	if hasDash || hasStar || hasPlus {
		return strings.TrimSpace(s[2:])
	}

	return stripLeadingNumberBullet(s)
}

func stripLeadingNumberBullet(s string) string {
	idxDigitEnd := findLeadingDigitEnd(s)
	if idxDigitEnd == 0 || idxDigitEnd >= len(s) {
		return s
	}

	hasDotOrParen := s[idxDigitEnd] == '.' || s[idxDigitEnd] == ')'
	if !hasDotOrParen {
		return s
	}

	return stripLeadingBulletRest(s, idxDigitEnd)
}

func findLeadingDigitEnd(s string) int {
	idx := 0
	for idx < len(s) && isASCIIDigit(s[idx]) {
		idx++
	}

	return idx
}

func isASCIIDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func stripLeadingBulletRest(s string, idxDigitEnd int) string {
	after := s[idxDigitEnd+1:]
	rest := strings.TrimLeft(after, " \t")
	hasTrimmedSpace := len(rest) < len(after)
	if hasTrimmedSpace {
		return strings.TrimSpace(rest)
	}

	return s
}

func stripTrailingComment(s string) string {
	idxHash := findCommentHashIndex(s)
	if idxHash >= 0 {
		s = strings.TrimSpace(s[:idxHash])
	}

	idxSlash := findCommentSlashIndex(s)
	if idxSlash >= 0 {
		s = strings.TrimSpace(s[:idxSlash])
	}

	return s
}

func findCommentHashIndex(s string) int {
	depthBracket, depthParen := 0, 0
	for i := 0; i < len(s); i++ {
		depthBracket, depthParen = updateNestingDepth(s[i], depthBracket, depthParen)
		isZeroDepth := depthBracket == 0 && depthParen == 0
		isHash := s[i] == '#'
		isDelim := i == 0 || isCommentDelim(s[i-1])
		if isZeroDepth && isHash && isDelim {
			return i
		}
	}

	return -1
}

func isCommentDelim(b byte) bool {
	return b == ' ' || b == '\t' || b == ']' || b == ')'
}

func findCommentSlashIndex(s string) int {
	depthBracket, depthParen := 0, 0
	for i := 0; i < len(s)-1; i++ {
		depthBracket, depthParen = updateNestingDepth(s[i], depthBracket, depthParen)
		isZeroDepth := depthBracket == 0 && depthParen == 0
		isSlashSlash := s[i] == '/' && s[i+1] == '/'
		isSpecial := isURIPrefix(s, i) || isTempPrefix(s, i)
		isDelim := i == 0 || isCommentDelim(s[i-1])
		if isZeroDepth && isSlashSlash && !isSpecial && isDelim {
			return i
		}
	}

	return -1
}

func updateNestingDepth(b byte, depthBracket, depthParen int) (int, int) {
	if b == '[' {
		return depthBracket + 1, depthParen
	}
	if b == ']' && depthBracket > 0 {
		return depthBracket - 1, depthParen
	}
	if b == '(' {
		return depthBracket, depthParen + 1
	}
	if b == ')' && depthParen > 0 {
		return depthBracket, depthParen - 1
	}

	return depthBracket, depthParen
}

func isURIPrefix(s string, i int) bool {
	return i > 0 && s[i-1] == ':'
}

func isTempPrefix(s string, i int) bool {
	rest := strings.ToLower(s[i+2:])
	hasTemp := strings.HasPrefix(rest, "temp")
	hasTmp := strings.HasPrefix(rest, "tmp")

	return hasTemp || hasTmp
}

func stripSurroundingQuotes(s string) string {
	trimmed := strings.TrimSpace(s)
	for len(trimmed) >= 2 {
		isQuoted := isMatchingQuote(trimmed)
		if !isQuoted {
			break
		}

		trimmed = strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	}

	return trimmed
}

func isMatchingQuote(s string) bool {
	isDouble := s[0] == '"' && s[len(s)-1] == '"'
	isSingle := s[0] == '\'' && s[len(s)-1] == '\''
	isBacktick := s[0] == '`' && s[len(s)-1] == '`'

	return isDouble || isSingle || isBacktick
}

func stripAtBrackets(s string) string {
	trimmed := strings.TrimSpace(s)
	hasPrefix := strings.HasPrefix(trimmed, "@[")
	hasSuffix := strings.HasSuffix(trimmed, "]")
	if hasPrefix && hasSuffix {
		return stripSurroundingQuotes(trimmed[2 : len(trimmed)-1])
	}

	return trimmed
}

func extractMarkdownURL(s string) string {
	trimmed := strings.TrimSpace(s)
	isCandidate := hasMarkdownLinkPattern(trimmed)
	if !isCandidate {
		return trimmed
	}

	return extractTargetBetweenParens(trimmed)
}

func hasMarkdownLinkPattern(s string) bool {
	hasParenSuffix := strings.HasSuffix(s, ")")
	isImage := strings.HasPrefix(s, "![")
	isLink := strings.HasPrefix(s, "[")

	return hasParenSuffix && (isImage || isLink)
}

func extractTargetBetweenParens(s string) string {
	idxOpenParen := findMarkdownURLOpenParen(s)
	idxCloseParen := strings.LastIndex(s, ")")
	hasValidBounds := idxOpenParen >= 0 && idxCloseParen > idxOpenParen
	if !hasValidBounds {
		return s
	}

	target := s[idxOpenParen+1 : idxCloseParen]

	return strings.Trim(strings.TrimSpace(target), "<>")
}

func findMarkdownURLOpenParen(s string) int {
	idx := strings.Index(s, "](")
	if idx >= 0 {
		return idx + 1
	}

	idxSpace := strings.Index(s, "] (")
	if idxSpace >= 0 {
		return idxSpace + 2
	}

	return -1
}

func unescapeURLPath(s string) string {
	unescaped, err := url.PathUnescape(s)
	if err == nil {
		return unescaped
	}

	return s
}

func stripFileScheme(s string) string {
	stripped, hasScheme := trimFileSchemePrefix(s)
	if !hasScheme {
		return s
	}

	normalized := normalizeWindowsDriveSlash(stripped)

	return filepath.FromSlash(normalized)
}

func trimFileSchemePrefix(s string) (string, bool) {
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "file:///") {
		return s[8:], true
	}

	if strings.HasPrefix(lower, "file://") {
		return s[7:], true
	}

	return s, false
}

func normalizeWindowsDriveSlash(s string) string {
	hasSlashDrive := len(s) >= 3 && isSlashChar(s[0]) && isDriveLetter(s[1]) && s[2] == ':'
	if hasSlashDrive {
		return s[1:]
	}

	return s
}

func isSlashChar(b byte) bool {
	return b == '/' || b == '\\'
}

func isDriveLetter(b byte) bool {
	isLower := b >= 'a' && b <= 'z'
	isUpper := b >= 'A' && b <= 'Z'

	return isLower || isUpper
}

func expandAllTokens(s string) string {
	expanded := expandWindowsEnv(s)
	expanded = expandUnixEnv(expanded)
	expanded = expandTilde(expanded)

	return expandMacroTemp(expanded)
}

func isCommandLine(s string) bool {
	hasCmdPrefix := startsWithShellCommand(s)
	hasSpace := strings.Contains(s, " ")

	return hasCmdPrefix && hasSpace
}

func startsWithShellCommand(s string) bool {
	commands := []string{"echo ", "cat ", "cd ", "mkdir ", "rm ", "touch "}
	for _, cmd := range commands {
		if strings.HasPrefix(s, cmd) {
			return true
		}
	}

	return false
}

func isWindowsDrivePath(s string) bool {
	hasDrive := len(s) >= 2 && isDriveLetter(s[0]) && s[1] == ':'

	return hasDrive
}

func resolveAbsPath(s string) string {
	if fileExists(s) {
		return s
	}

	return filepath.Clean(s)
}

func resolveFlexiblePath(s string) string {
	if isCommandLine(s) {
		return s
	}

	isAbs := filepath.IsAbs(s) || isWindowsDrivePath(s)
	if isAbs {
		return resolveAbsPath(s)
	}

	return resolveRelativePath(s)
}

func resolveRelativePath(s string) string {
	cwd, err := os.Getwd()
	hasCWD := err == nil
	if !hasCWD {
		return filepath.Clean(s)
	}

	resolved, isFound := tryResolveLocations(s, cwd)
	if isFound {
		return resolved
	}

	return filepath.Clean(s)
}

func tryResolveLocations(s, cwd string) (string, bool) {
	rcResolved, hasRC := checkRepoCacheCandidate(s, cwd)
	if hasRC {
		return rcResolved, true
	}

	cwdResolved, hasCWD := resolveFromCWD(s, cwd)
	if hasCWD {
		return cwdResolved, true
	}

	return resolveFromRepoRoot(s, cwd)
}

func checkRepoCacheCandidate(s, cwd string) (string, bool) {
	if !isRepoCacheRelative(s) {
		return "", false
	}

	return resolveRepoCacheCandidates(s, cwd)
}

func isRepoCacheRelative(s string) bool {
	norm := filepath.ToSlash(s)
	hasPrefix := strings.HasPrefix(norm, "../repo-cache/")

	return hasPrefix
}

func resolveRepoCacheCandidates(s, cwd string) (string, bool) {
	sub := extractRepoCacheSubpath(s)
	candidates := buildRepoCacheCandidates(cwd, sub)
	for _, cand := range candidates {
		if fileExists(cand) {
			return cand, true
		}
	}

	return "", false
}

func extractRepoCacheSubpath(s string) string {
	norm := filepath.ToSlash(s)
	sub := strings.TrimPrefix(norm, "../repo-cache/")

	return filepath.FromSlash(sub)
}

func buildRepoCacheCandidates(cwd, sub string) []string {
	cwdParent := filepath.Dir(cwd)
	repoRoot := resolveRepoRoot(cwd)
	repoRootParent := filepath.Dir(repoRoot)

	return []string{
		filepath.Join(cwdParent, "repo-cache", sub),
		filepath.Join(cwd, "..", "repo-cache", sub),
		filepath.Join(repoRoot, "..", "repo-cache", sub),
		filepath.Join(repoRootParent, "repo-cache", sub),
		filepath.Join(filepath.Dir(cwdParent), "repo-cache", sub),
	}
}

func resolveFromCWD(s, cwd string) (string, bool) {
	candidate := filepath.Join(cwd, s)
	if fileExists(candidate) {
		return candidate, true
	}

	return "", false
}

func resolveFromRepoRoot(s, cwd string) (string, bool) {
	root := resolveRepoRoot(cwd)
	hasRoot := len(root) > 0
	if !hasRoot {
		return "", false
	}

	candidate := filepath.Join(root, s)
	if fileExists(candidate) {
		return candidate, true
	}

	return "", false
}

func resolveRepoRoot(cwd string) string {
	res := gitutil.RepoRoot(cwd)
	hasRoot := res.IsSuccess() && len(res.Data) > 0
	if hasRoot {
		return res.Data
	}

	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func hasFlexibleFormat(input string) bool {
	hasAtBracket := strings.Contains(input, "@[")
	hasMarkdown := strings.HasPrefix(input, "![") || strings.Contains(input, "](")
	hasFileScheme := strings.Contains(strings.ToLower(input), "file://")
	hasBullet := hasLeadingBulletToken(input)
	hasRCRel := strings.HasPrefix(filepath.ToSlash(input), "../repo-cache/")

	return hasAtBracket || hasMarkdown || hasFileScheme || hasBullet || hasRCRel
}

func hasLeadingBulletToken(s string) bool {
	trimmed := strings.TrimSpace(s)
	hasDash := strings.HasPrefix(trimmed, "- ")
	hasStar := strings.HasPrefix(trimmed, "* ")
	hasPlus := strings.HasPrefix(trimmed, "+ ")
	hasNumber := stripLeadingNumberBullet(trimmed) != trimmed

	return hasDash || hasStar || hasPlus || hasNumber
}
