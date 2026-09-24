package cmdui

import (
	"testing"
)

func TestDetectSyntaxLanguage(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"config.json", "json"},
		{"readme.md", "markdown"},
		{"script.py", "python"},
		{"index.html", "html"},
		{"styles.css", "css"},
		{"run.sh", "shell"},
		{"unknown.xyz", "plaintext"},
	}

	for _, c := range cases {
		actual := detectSyntaxLanguage(c.path)
		if actual != c.expected {
			t.Errorf("path %s: expected %s, got %s", c.path, c.expected, actual)
		}
	}
}
