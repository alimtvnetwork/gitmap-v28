package stablejson

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestWriteArray_EmptyIsBracketsNewline pins the headline jq-compat
// guarantee — zero items must encode as `[]\n`, never `null`, never
// `[]` without trailing newline. This is the one byte sequence
// downstream pipelines bake into their length/empty checks.
func TestWriteArray_EmptyIsBracketsNewline(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteArray(&buf, nil); err != nil {
		t.Fatalf("WriteArray: %v", err)
	}
	if got := buf.String(); got != "[]\n" {
		t.Fatalf("empty array: want %q, got %q", "[]\n", got)
	}
}

// TestWriteArray_ByteCompatWithEncoder is the safety-net test that
// proves migrating a caller from `json.Encoder.SetIndent("", "  ")`
// to stablejson.WriteArray does NOT change output bytes. Without
// this guarantee, every existing golden fixture would have to be
// regenerated — which would defeat the whole point of the package.
func TestWriteArray_ByteCompatWithEncoder(t *testing.T) {
	type entry struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Exec string `json:"exec"`
	}
	src := []entry{
		{Name: "a", Path: "/p/a", Exec: "/bin/a"},
		{Name: "b", Path: "/p/b", Exec: "/bin/b --flag"},
	}

	var encBuf bytes.Buffer
	enc := json.NewEncoder(&encBuf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(src); err != nil {
		t.Fatalf("encoder: %v", err)
	}

	stable := make([][]Field, 0, len(src))
	for _, e := range src {
		stable = append(stable, []Field{
			{Key: "name", Value: e.Name},
			{Key: "path", Value: e.Path},
			{Key: "exec", Value: e.Exec},
		})
	}
	var stableBuf bytes.Buffer
	if err := WriteArray(&stableBuf, stable); err != nil {
		t.Fatalf("WriteArray: %v", err)
	}

	if !bytes.Equal(encBuf.Bytes(), stableBuf.Bytes()) {
		t.Fatalf("byte-compat broken\n--- json.Encoder\n%s--- stablejson\n%s",
			encBuf.String(), stableBuf.String())
	}
}

// TestWriteArray_PreservesCallerKeyOrder is the whole point: even
// when keys are passed in deliberately weird (non-alphabetical,
// non-struct-declaration) order, that order is what comes out.
// Catches any future "helpful" sorting added to the encoder.
func TestWriteArray_PreservesCallerKeyOrder(t *testing.T) {
	got := mustEncode(t, [][]Field{{
		{Key: "zeta", Value: 1},
		{Key: "alpha", Value: 2},
		{Key: "mike", Value: 3},
	}})
	zeta := strings.Index(got, `"zeta"`)
	alpha := strings.Index(got, `"alpha"`)
	mike := strings.Index(got, `"mike"`)
	if !(zeta < alpha && alpha < mike) {
		t.Fatalf("key order not preserved: zeta=%d alpha=%d mike=%d in %q",
			zeta, alpha, mike, got)
	}
}

// TestWriteArray_ValueTypesViaMarshal confirms that values flow
// through json.Marshal (so numbers, bools, nulls, and json.Marshaler
// implementations all work without the package having to special-
// case them). Strings must still be quoted+escaped properly.
func TestWriteArray_ValueTypesViaMarshal(t *testing.T) {
	got := mustEncode(t, [][]Field{{
		{Key: "s", Value: "hello \"world\""},
		{Key: "n", Value: 42},
		{Key: "b", Value: true},
		{Key: "z", Value: nil},
	}})
	for _, want := range []string{
		`"s": "hello \"world\""`,
		`"n": 42`,
		`"b": true`,
		`"z": null`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %q", want, got)
		}
	}
}

// TestWriteArray_TrailingNewlineMatchesEncoder pins the trailing
// `\n` — Encoder.Encode appends one, and golden fixtures captured
// from the old encoder include it. Dropping it would silently break
// every byte-exact contract test in the repo.
func TestWriteArray_TrailingNewlineMatchesEncoder(t *testing.T) {
	got := mustEncode(t, [][]Field{{{Key: "k", Value: "v"}}})
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("missing trailing newline: %q", got)
	}
}

func mustEncode(t *testing.T, items [][]Field) string {
	t.Helper()
	var buf bytes.Buffer
	if err := WriteArray(&buf, items); err != nil {
		t.Fatalf("WriteArray: %v", err)
	}

	return buf.String()
}
