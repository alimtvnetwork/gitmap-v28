package cliexit

// Unit tests for the cliexit formatter contract. Locks the byte-exact
// shape of the user-facing failure line so future refactors can't
// silently drift the format that wrapper scripts (and our own CI
// annotations) grep against.

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// TestFormatLine_Shape pins the documented format for both the
// with-subject and elided-subject cases.
func TestFormatLine_Shape(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		command string
		op      string
		subject string
		err     error
		want    string
	}{
		{
			name:    "with_subject",
			command: "clone-from",
			op:      "parse",
			subject: "/tmp/manifest.json",
			err:     errors.New("invalid json"),
			want:    "gitmap clone-from: parse on /tmp/manifest.json failed: invalid json",
		},
		{
			name:    "empty_subject_elided",
			command: "scan",
			op:      "config-load",
			subject: "",
			err:     errors.New("permission denied"),
			want:    "gitmap scan: config-load failed: permission denied",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := formatLine(tc.command, tc.op, tc.subject, tc.err)
			if got != tc.want {
				t.Fatalf("formatLine mismatch:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// TestWriteReport_NilErrSurfacesBug guarantees the logic-bug guard
// fires loudly instead of producing a half-formed line. This is the
// contract that lets call sites stop defending against nil-err at
// every site themselves.
func TestWriteReport_NilErrSurfacesBug(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	writeReport(&buf, "scan", "walk", "/repo", nil)
	out := buf.String()
	if !strings.Contains(out, "BUG") {
		t.Fatalf("expected BUG marker, got: %q", out)
	}
	if !strings.Contains(out, "op=walk") || !strings.Contains(out, "subject=/repo") {
		t.Fatalf("expected op + subject in BUG line, got: %q", out)
	}
}

// TestWriteReport_TrailingNewline asserts every line ends with \n so
// concatenated stderr output stays line-delimited (important for
// log scrapers and the CI annotation regex).
func TestWriteReport_TrailingNewline(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	writeReport(&buf, "clone", "git-clone", "https://x/y.git", errors.New("boom"))
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Fatalf("missing trailing newline: %q", buf.String())
	}
}
