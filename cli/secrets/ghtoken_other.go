//go:build !windows

package secrets

func tokenFromSystemRegistry() (string, SourceType, bool) {
	return "", SourceNone, false
}
