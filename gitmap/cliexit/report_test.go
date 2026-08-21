package cliexit

// Tests for the structured Report path: human mode shape, JSON mode
// shape, deterministic ordering, and the nil-err BUG guard.

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// TestWriteStructured_HumanMode pins the indented context tail and
// the leading canonical line so wrapper scripts that grep `gitmap
// <cmd>:` keep working unchanged.
func TestWriteStructured_HumanMode(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	writeStructured(&buf, Context{
		Command: "clone-now",
		Op:      "git-clone",
		Path:    "/tmp/repo",
		Args:    []string{"--execute", "manifest.json"},
		Mode:    "execute",
		Extras:  map[string]string{"row": "3", "url": "https://x/y.git"},
		Err:     errors.New("boom"),
	}, OutputHuman)

	got := buf.String()
	wantLead := "gitmap clone-now: git-clone on /tmp/repo failed: boom\n"
	isPrefixMissing := !strings.HasPrefix(got, wantLead)
	if isPrefixMissing == true {
		t.Fatalf("missing canonical lead line.\n got: %q\nwant prefix: %q", got, wantLead)
	}
	assertHumanOutputElements(t, got)
}

func assertHumanOutputElements(t *testing.T, got string) {
	t.Helper()
	for _, want := range []string{
		"  mode=execute",
		"  args=--execute manifest.json",
		"  row=3",
		"  url=https://x/y.git",
	} {
		isMissing := !strings.Contains(got, want)
		if isMissing == true {
			t.Fatalf("human output missing %q\nfull:\n%s", want, got)
		}
	}
}

// TestWriteStructured_JSONMode asserts the JSON payload is single-
// line, decodable, and carries every populated field with its
// declared type (args = []string, extras = map).
func TestWriteStructured_JSONMode(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	writeStructured(&buf, Context{
		Command: "scan",
		Op:      "walk",
		Path:    "/repos",
		Args:    []string{"--root", "/repos"},
		Mode:    "dry-run",
		Extras:  map[string]string{"depth": "2"},
		Err:     errors.New("io error"),
	}, OutputJSON)

	isSingleLine := strings.Count(buf.String(), "\n") == 1
	if isSingleLine == false {
		t.Fatalf("expected exactly one line, got: %q", buf.String())
	}
	assertJSONOutput(t, buf)
}

func assertJSONOutput(t *testing.T, buf bytes.Buffer) {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("output not valid JSON: %v\nraw=%q", err, buf.String())
	}
	for _, k := range []string{"command", "op", "path", "mode", "args", "extras", "error"} {
		_, ok := got[k]
		isMissingKey := !ok
		if isMissingKey == true {
			t.Fatalf("JSON missing key %q: %v", k, got)
		}
	}
	isErrorMismatch := got["error"] != "io error"
	if isErrorMismatch == true {
		t.Fatalf("error field wrong: %v", got["error"])
	}
}

// TestWriteStructured_OmitsEmpty verifies that empty optional fields
// don't leak into the output (no `path=`, no `args=`, no extras key).
func TestWriteStructured_OmitsEmpty(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	writeStructured(&buf, Context{
		Command: "amend",
		Op:      "parse",
		Err:     errors.New("bad input"),
	}, OutputHuman)

	got := buf.String()
	hasOnPath := strings.Contains(got, " on ")
	if hasOnPath == true {
		t.Fatalf("empty Path should elide ' on <subject>': %q", got)
	}
	hasMode := strings.Contains(got, "mode=")
	hasArgs := strings.Contains(got, "args=")
	hasModeOrArgs := hasMode == true || hasArgs == true
	if hasModeOrArgs == true {
		t.Fatalf("empty Mode/Args should not render: %q", got)
	}
}

// TestReport_NilErrSurfacesBug mirrors the Reportf guard so callers
// that misuse the structured entry-point also fail loudly.
func TestReport_NilErrSurfacesBug(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	writeStructured(&buf, Context{
		Command: "scan",
		Op:      "walk",
		Path:    "/repo",
	}, OutputHuman)
	isBugMissing := !strings.Contains(buf.String(), "BUG")
	if isBugMissing == true {
		t.Fatalf("expected BUG marker, got: %q", buf.String())
	}
}
