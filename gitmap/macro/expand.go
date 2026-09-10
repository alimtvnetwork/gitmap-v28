package macro

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var winEnvRegex = regexp.MustCompile(`%([a-zA-Z0-9_]+)%`)

// ExpandPathAndEnv expands Windows %VAR%, Unix $VAR, and tilde ~ path tokens.
func ExpandPathAndEnv(input string) string {
	if len(input) == 0 {
		return input
	}
	expanded := expandWindowsEnv(input)
	expanded = os.ExpandEnv(expanded)
	expanded = expandTilde(expanded)

	return expanded
}

func expandWindowsEnv(input string) string {
	if !strings.Contains(input, "%") {
		return input
	}

	return winEnvRegex.ReplaceAllStringFunc(input, func(token string) string {
		varName := token[1 : len(token)-1]
		if val, hasVal := getEnvCaseInsensitive(varName); hasVal {
			return val
		}

		return token
	})
}

func getEnvCaseInsensitive(key string) (string, bool) {
	if val, hasKey := os.LookupEnv(key); hasKey {
		return val, true
	}
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], key) {
			return parts[1], true
		}
	}
	if strings.EqualFold(key, "temp") || strings.EqualFold(key, "tmp") {
		return os.TempDir(), true
	}

	return "", false
}

func expandTilde(input string) string {
	if !strings.Contains(input, "~") {
		return input
	}
	home, err := os.UserHomeDir()
	if err != nil || len(home) == 0 {
		return input
	}

	return replaceTildeTokens(input, home)
}

func replaceTildeTokens(input, home string) string {
	tokens := strings.Split(input, " ")
	for i, token := range tokens {
		tokens[i] = resolveSingleTildeToken(token, home)
	}

	return strings.Join(tokens, " ")
}

func resolveSingleTildeToken(token, home string) string {
	if token == "~" {
		return home
	}
	if strings.HasPrefix(token, "~/") {
		return filepath.Join(home, token[2:])
	}
	if runtime.GOOS == "windows" && strings.HasPrefix(token, `~\`) {
		return filepath.Join(home, token[2:])
	}

	return token
}

// NormalizeTargetPath resolves quotes, env variables (%VAR%, $VAR), tilde (~),
// cross-platform temp aliases (//temp, /temp, etc.), and relative directory paths.
func NormalizeTargetPath(target, currentDir string) string {
	stripped := stripWrappingQuotes(target)
	if len(stripped) == 0 || stripped == "-" {
		return stripped
	}
	expanded := ExpandPathAndEnv(stripped)
	if resolvedTemp, hasTemp := resolveTempAlias(expanded); hasTemp {
		return filepath.Clean(resolvedTemp)
	}
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded)
	}
	if len(currentDir) > 0 {
		return filepath.Clean(filepath.Join(currentDir, expanded))
	}

	return filepath.Clean(expanded)
}

func stripWrappingQuotes(input string) string {
	s := strings.TrimSpace(input)
	if len(s) < 2 {
		return s
	}
	hasDouble := s[0] == '"' && s[len(s)-1] == '"'
	hasSingle := s[0] == '\'' && s[len(s)-1] == '\''
	if hasDouble || hasSingle {
		return strings.TrimSpace(s[1 : len(s)-1])
	}

	return s
}

func resolveTempAlias(input string) (string, bool) {
	s := strings.TrimSpace(input)
	hasLeadingSlash := strings.HasPrefix(s, "/") || strings.HasPrefix(s, "\\")
	if !hasLeadingSlash {
		return input, false
	}
	stripped := strings.TrimLeft(s, "/\\")

	return matchTempPrefix(stripped)
}

func matchTempPrefix(stripped string) (string, bool) {
	lower := strings.ToLower(stripped)
	if lower == "temp" || lower == "tmp" {
		return os.TempDir(), true
	}
	if strings.HasPrefix(lower, "temp/") || strings.HasPrefix(lower, "temp\\") {
		return filepath.Join(os.TempDir(), stripped[5:]), true
	}
	if strings.HasPrefix(lower, "tmp/") || strings.HasPrefix(lower, "tmp\\") {
		return filepath.Join(os.TempDir(), stripped[4:]), true
	}

	return "", false
}
