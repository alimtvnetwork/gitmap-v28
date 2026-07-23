package cmd

import (
	"reflect"
	"testing"
)

// TestShouldRewriteToCloneCoversReportedInvocations pins the three exact
// invocations the user reported as failing with "Unknown command", plus
// the leading-flag and SSH variants that share the same code path. Any
// regression here would re-break the bare-URL shortcut.
func TestShouldRewriteToCloneCoversReportedInvocations(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{
			name: "single comma-glued URL list (PowerShell paste)",
			args: []string{"https://github.com/a/b,https://github.com/a/c,https://github.com/a/d"},
			want: true,
		},
		{
			name: "comma-then-space split across argv (bash paste)",
			args: []string{"https://github.com/a/b,", "https://github.com/a/c,", "https://github.com/a/d"},
			want: true,
		},
		{
			name: "mixed comma + space separators across argv",
			args: []string{"https://github.com/a/b,", "https://github.com/a/c", "https://github.com/a/d"},
			want: true,
		},
		{
			name: "single bare URL",
			args: []string{"https://github.com/a/b"},
			want: true,
		},
		{
			name: "leading flag then URL",
			args: []string{"--verbose", "https://github.com/a/b"},
			want: true,
		},
		{
			name: "SSH shorthand",
			args: []string{"git@github.com:a/b.git"},
			want: true,
		},
		{
			name: "GitLab URL",
			args: []string{"https://gitlab.com/a/b"},
			want: true,
		},
		{
			name: "no args",
			args: nil,
			want: false,
		},
		{
			name: "known subcommand, not a URL",
			args: []string{"scan"},
			want: false,
		},
		{
			name: "folder path, not a URL",
			args: []string{"./my-repo"},
			want: false,
		},
	}

	for _, tc := range cases {
		got := shouldRewriteToClone(tc.args)
		if got != tc.want {
			t.Errorf("%s: shouldRewriteToClone(%q) = %v, want %v", tc.name, tc.args, got, tc.want)
		}
	}
}

// TestSplitOnCommaHandlesPasteArtifacts pins the splitter against the
// real-world paste shapes (trailing commas, double spaces, empty
// pieces) that produced the original "Unknown command" reports.
func TestSplitOnCommaHandlesPasteArtifacts(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{"a, b, c", []string{"a", "b", "c"}},
		{"a,,b", []string{"a", "b"}},
		{"a,", []string{"a"}},
		{",a,b,", []string{"a", "b"}},
		{"  a  ,  b  ", []string{"a", "b"}},
		{"", []string{}},
	}

	for _, tc := range cases {
		got := splitOnComma(tc.in)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("splitOnComma(%q) = %#v, want %#v", tc.in, got, tc.want)
		}
	}
}

// TestLooksLikeURLTokenAcceptsAllSupportedShapes guards the
// shape-detection helper against accidental tightening — the user's
// invocations all flow through this single predicate.
func TestLooksLikeURLTokenAcceptsAllSupportedShapes(t *testing.T) {
	good := []string{
		"https://github.com/a/b",
		"http://example.com/a/b",
		"ssh://git@github.com/a/b",
		"git@github.com:a/b.git",
		"https://gitlab.com/a/b",
		"https://github.com/a/b,https://github.com/c/d",
		"  https://github.com/a/b  ",
	}
	for _, s := range good {
		if !looksLikeURLToken(s) {
			t.Errorf("looksLikeURLToken(%q) = false, want true", s)
		}
	}

	bad := []string{"", "scan", "./my-repo", "version", "--verbose"}
	for _, s := range bad {
		if looksLikeURLToken(s) {
			t.Errorf("looksLikeURLToken(%q) = true, want false", s)
		}
	}
}
