//go:build !windows

package ghtoken

func tokenFromSystemRegistry() (string, SourceType, bool) {
	return "", SourceNone, false
}
