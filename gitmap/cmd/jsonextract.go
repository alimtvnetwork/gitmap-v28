package cmd

// extractJSONString extracts a string value from JSON bytes by key.
func extractJSONString(data []byte, key string) string {
	s := string(data)
	needle := `"` + key + `"`
	idx := findKeyValue(s, needle)
	if idx < 0 {
		return ""
	}

	return extractQuotedValue(s, idx)
}

// findKeyValue finds the position after a JSON key and colon.
func findKeyValue(s, needle string) int {
	idx := indexOf(s, needle)
	if idx < 0 {
		return -1
	}

	idx += len(needle)
	for idx < len(s) && (s[idx] == ' ' || s[idx] == ':' || s[idx] == '\t') {
		idx++
	}

	return idx
}

// extractQuotedValue extracts a quoted string starting at idx.
func extractQuotedValue(s string, idx int) string {
	if idx >= len(s) || s[idx] != '"' {
		return ""
	}

	end := indexOf(s[idx+1:], `"`)
	if end >= 0 {
		return s[idx+1 : idx+1+end]
	}

	return ""
}

// indexOf returns the index of substr in s, or -1.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
