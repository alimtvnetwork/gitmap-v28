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
