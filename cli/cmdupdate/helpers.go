package cmdupdate

import "strings"

// Delegate hooks
var (
	ResolveDeployedAndConfigPathsFn func() (string, string)
	RunPostUpdateMigrateFn          func() error
	RequireOnlineFn                 func()
	PrintGitmapIdentityBlockLongFn  func()
)

func ResolveDeployedAndConfigPaths() (string, string) {
	if ResolveDeployedAndConfigPathsFn != nil {
		return ResolveDeployedAndConfigPathsFn()
	}
	return "", ""
}

func RunPostUpdateMigrate() error {
	if RunPostUpdateMigrateFn != nil {
		return RunPostUpdateMigrateFn()
	}
	return nil
}

func requireOnline() {
	if RequireOnlineFn != nil {
		RequireOnlineFn()
	}
}

func extractJSONString(data []byte, key string) string {
	s := string(data)
	needle := `"` + key + `"`
	idx := findKeyValue(s, needle)
	if idx < 0 {
		return ""
	}

	return extractQuotedValue(s, idx)
}

func findKeyValue(s, needle string) int {
	idx := strings.Index(s, needle)
	if idx < 0 {
		return -1
	}

	idx += len(needle)
	for idx < len(s) && (s[idx] == ' ' || s[idx] == ':' || s[idx] == '	') {
		idx++
	}

	return idx
}

func extractQuotedValue(s string, idx int) string {
	if idx >= len(s) || s[idx] != '"' {
		return ""
	}

	end := strings.Index(s[idx+1:], `"`)
	if end >= 0 {
		return s[idx+1 : idx+1+end]
	}

	return ""
}
