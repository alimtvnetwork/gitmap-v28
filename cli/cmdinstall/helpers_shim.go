package cmdinstall

import "strings"

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
	for idx < len(s) && (s[idx] == ' ' || s[idx] == ':' || s[idx] == '\t') {
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
