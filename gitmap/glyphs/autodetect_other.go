//go:build !windows

package glyphs

import (
	"os"
	"strings"
)

func init() {
	isLegacyUnixHost = legacyUnixHost
}

func legacyUnixHost() bool {
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return true
	}

	return !isUTF8Locale(resolveLocale())
}

func resolveLocale() string {
	locale := os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LC_CTYPE")
	}

	if locale == "" {
		locale = os.Getenv("LANG")
	}

	return strings.ToLower(locale)
}

func isUTF8Locale(loc string) bool {
	if loc == "" {
		return true
	}

	return strings.Contains(loc, "utf-8") || strings.Contains(loc, "utf8")
}
