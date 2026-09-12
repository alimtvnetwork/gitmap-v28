package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"testing"
)

// TestParseHygieneFormat exercises the format flag parser.
func TestParseHygieneFormat(t *testing.T) {
	cases := map[string]hygieneFormat{
		"":      hygieneFormatTable,
		"table": hygieneFormatTable,
		"json":  hygieneFormatJSON,
		"csv":   hygieneFormatCSV,
	}

	for in, want := range cases {
		got, err := parseHygieneFormat(in)
		if err != nil || got != want {
			t.Fatalf("parseHygieneFormat(%q) = %q,%v want %q", in, got, err, want)
		}
	}

	if _, err := parseHygieneFormat("xml"); err == nil {
		t.Fatalf("expected error for invalid format")
	}
}

// TestEmitJSONAndCSV captures stdout and decodes the emitted payload
// to confirm the schema used by stale/dedupe/size/orphans is parseable.
func TestEmitJSONAndCSV(t *testing.T) {
	type row struct {
		Path string `json:"path"`
	}

	withStdout(t, func() { emitJSON([]row{{Path: "a"}, {Path: "b"}}) }, func(buf []byte) {
		var got []row
		if err := json.Unmarshal(buf, &got); err != nil {
			t.Fatalf("json: %v\n%s", err, buf)
		}

		if len(got) != 2 || got[0].Path != "a" {
			t.Fatalf("unexpected json: %+v", got)
		}
	})
	withStdout(t, func() { emitCSV([]string{"path"}, [][]string{{"a"}, {"b"}}) }, func(buf []byte) {
		r := csv.NewReader(bytes.NewReader(buf))
		recs, err := r.ReadAll()
		if err != nil {
			t.Fatalf("csv: %v", err)
		}

		if len(recs) != 3 || recs[0][0] != "path" || recs[2][0] != "b" {
			t.Fatalf("unexpected csv: %v", recs)
		}
	})
}

// withStdout redirects os.Stdout for the duration of fn and feeds the
// captured bytes to inspect.
func withStdout(t *testing.T, fn func(), inspect func([]byte)) {
	t.Helper()
	out := captureStdoutForTest(t, fn)
	inspect([]byte(out))
}
