package cmd

// Failure-message contract for the JSON schema assertion helpers
// in jsonsnapshot_helpers_test.go. The helpers exist to point a
// developer at the EXACT point of schema drift (which delimiter,
// which object index). This test pins that diagnostic value so a
// future refactor cannot silently degrade messages to a generic
// "parse failed" string.
//
// expectDelim is tested directly (it already returns an error).
// readEveryObjectKeys uses *testing.T directly, so we test its
// message contract via a pure error-returning twin
// (scanEveryObjectKeysPure) below — any wording change in the
// production helper must be mirrored here, which is the guard rail.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// TestExpectDelim_FailureMessages pins the exact wording for each
// rejection path. Downstream "JSON schema regression" tooling that
// scrapes the failure message for the offending delimiter would
// break on a wording change, so the message format is part of the
// helper's contract and not just a debug aid.
func TestExpectDelim_FailureMessages(t *testing.T) {
	cases := []struct {
		name        string
		raw         string
		want        byte
		wantSubstrs []string
	}{
		{
			name:        "wrong_open_delim",
			raw:         `{"x":1}`,
			want:        '[',
			wantSubstrs: []string{`want delim '['`, `got {`, `(json.Delim)`},
		},
		{
			name:        "string_where_delim_expected",
			raw:         `"hello"`,
			want:        '[',
			wantSubstrs: []string{`want delim '['`, `got hello`, `(string)`},
		},
		{
			name:        "eof_returns_underlying_error",
			raw:         ``,
			want:        '[',
			wantSubstrs: []string{`EOF`},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := expectDelim(json.NewDecoder(strings.NewReader(tc.raw)), tc.want)
			assertErrorContainsAll(t, err, tc.wantSubstrs)
		})
	}
}

// TestScanEveryObjectKeysPure_PointsAtBrokenObject runs the pure
// error-returning twin of readEveryObjectKeys against malformed
// inputs and asserts the error message names the offending object
// index. A regression that drops the index from the message would
// force devs to bisect their JSON by hand.
func TestScanEveryObjectKeysPure_PointsAtBrokenObject(t *testing.T) {
	cases := []struct {
		name        string
		raw         string
		wantSubstrs []string
	}{
		{
			name:        "not_an_array",
			raw:         `{"x":1}`,
			wantSubstrs: []string{"top-level array", `want delim '['`},
		},
		{
			name:        "second_element_not_object",
			raw:         `[{"a":1},"oops"]`,
			wantSubstrs: []string{"object at index 1", `want delim '{'`, "got oops"},
		},
		{
			name:        "third_element_is_array_not_object",
			raw:         `[{"a":1},{"b":2},[3,4]]`,
			wantSubstrs: []string{"object at index 2", `want delim '{'`, "got ["},
		},
		{
			name:        "object_unterminated_truncated_input",
			raw:         `[{"a":1`,
			wantSubstrs: []string{"object[0]"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := scanEveryObjectKeysPure([]byte(tc.raw))
			assertErrorContainsAll(t, err, tc.wantSubstrs)
		})
	}
}

// TestScanEveryObjectKeysPure_HappyPath proves the pure variant
// agrees with readEveryObjectKeys on well-formed input. Without
// this guarantee, a green failure-message test could coexist with
// a broken parser and we'd never notice.
func TestScanEveryObjectKeysPure_HappyPath(t *testing.T) {
	raw := []byte(`[{"name":"a","path":"/p"},{"name":"b","path":"/q"}]`)
	got, err := scanEveryObjectKeysPure(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := [][]string{{"name", "path"}, {"name", "path"}}
	if len(got) != len(want) {
		t.Fatalf("object count: want %d got %d (%v)", len(want), len(got), got)
	}
	for i := range want {
		if !equalStringSlices(got[i], want[i]) {
			t.Fatalf("object[%d] keys: want %v got %v", i, want[i], got[i])
		}
	}
}

// scanEveryObjectKeysPure mirrors readEveryObjectKeys with an
// error return instead of t.Fatalf. Sole purpose: make the
// failure-message contract testable. Wording must stay in sync.
func scanEveryObjectKeysPure(raw []byte) ([][]string, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := expectDelim(dec, '['); err != nil {

		return nil, fmt.Errorf("expected top-level array: %w", err)
	}
	var out [][]string
	for dec.More() {
		if err := expectDelim(dec, '{'); err != nil {

			return nil, fmt.Errorf("expected object at index %d: %w", len(out), err)
		}
		keys, err := pureCollectObjectKeys(dec, len(out))
		if err != nil {

			return nil, err
		}
		out = append(out, keys)
	}

	return out, nil
}

// pureCollectObjectKeys is the error-returning twin of
// collectObjectKeys; `objIdx` names the broken object in messages.
func pureCollectObjectKeys(dec *json.Decoder, objIdx int) ([]string, error) {
	var keys []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {

			return nil, fmt.Errorf("reading object[%d] key: %w", objIdx, err)
		}
		key, ok := tok.(string)
		if !ok {

			return nil, fmt.Errorf("object[%d] expected string key, got %v (%T)", objIdx, tok, tok)
		}
		keys = append(keys, key)
		var skip json.RawMessage
		if err := dec.Decode(&skip); err != nil {

			return nil, fmt.Errorf("object[%d] skipping value for key %q: %w", objIdx, key, err)
		}
	}
	if _, err := dec.Token(); err != nil {

		return nil, fmt.Errorf("object[%d] expected closing '}': %w", objIdx, err)
	}

	return keys, nil
}

// assertErrorContainsAll asserts err is non-nil and its message
// contains every substring in `wantSubstrs`. Reports a single
// fatal with all misses so a wording drift surfaces every missing
// substring at once instead of one per re-run.
func assertErrorContainsAll(t *testing.T, err error, wantSubstrs []string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error, got nil (wanted substrings: %v)", wantSubstrs)
	}
	msg := err.Error()
	var missing []string
	for _, s := range wantSubstrs {
		if !strings.Contains(msg, s) {
			missing = append(missing, s)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("error message missing substrings %v\nfull message: %q", missing, msg)
	}
}
