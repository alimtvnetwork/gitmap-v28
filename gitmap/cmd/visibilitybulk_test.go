// Package cmd — visibilitybulk_test.go
//
// Unit tests for bulk-visibility parsing helpers.
package cmd

import (
	"fmt"
	"testing"
)

func TestExtractBaseAndVersionFromArg_URL(t *testing.T) {
	// Per fix-repo digit-capture rule (mem://constraints, v4.12.0): derive the
	// expected version from the same int used to format the input URL so a
	// future `gitmap-vN` bump rewrites both sides atomically.
	const wantVer = 25
	input := fmt.Sprintf("https://github.com/alimtvnetwork/gitmap-v%d", wantVer)
	base, ver := extractBaseAndVersionFromArg(input)
	if base != "gitmap" || ver != wantVer {
		t.Fatalf("expected (gitmap, %d), got (%s, %d)", wantVer, base, ver)
	}
}

func TestExtractBaseAndVersionFromArg_Slug(t *testing.T) {
	base, ver := extractBaseAndVersionFromArg("alimtvnetwork/gitmap-v40")
	if base != "gitmap" || ver != 40 {
		t.Fatalf("expected (gitmap, 40), got (%s, %d)", base, ver)
	}
}

func TestExtractBaseAndVersionFromArg_BareNameVersioned(t *testing.T) {
	base, ver := extractBaseAndVersionFromArg("macro-ahk-v10")
	if base != "macro-ahk" || ver != 10 {
		t.Fatalf("expected (macro-ahk, 10), got (%s, %d)", base, ver)
	}
}

func TestExtractBaseAndVersionFromArg_UnversionedDefaultsTo1(t *testing.T) {
	base, ver := extractBaseAndVersionFromArg("gitmap")
	if base != "gitmap" || ver != 1 {
		t.Fatalf("expected (gitmap, 1), got (%s, %d)", base, ver)
	}
}

func TestExtractBaseAndVersionFromArg_UnversionedURL(t *testing.T) {
	base, ver := extractBaseAndVersionFromArg("https://github.com/alimtvnetwork/gitmap")
	if base != "gitmap" || ver != 1 {
		t.Fatalf("expected (gitmap, 1), got (%s, %d)", base, ver)
	}
}

func TestParseBulkRequest_TwoArgValid(t *testing.T) {
	// Derive both the input slug AND expected StartVer from the same int so
	// a future `gitmap-vN` bump (fix-repo rewriter) updates them atomically.
	// Hardcoding either side leaks the digit-capture desync bug.
	const inputVer = 27
	slug := fmt.Sprintf("gitmap-v%d", inputVer)
	wantStartVer := inputVer - 1
	req, ok := parseBulkRequest([]string{slug, "3"})
	if !ok {
		t.Fatal("expected ok=true for valid two-arg request")
	}
	if req.BaseRepo != "gitmap" {
		t.Fatalf("expected BaseRepo=gitmap, got %s", req.BaseRepo)
	}
	if req.StartVer != wantStartVer {
		t.Fatalf("expected StartVer=%d, got %d", wantStartVer, req.StartVer)
	}
	if req.Count != 3 {
		t.Fatalf("expected Count=3, got %d", req.Count)
	}
}


func TestParseBulkRequest_Empty(t *testing.T) {
	_, ok := parseBulkRequest([]string{})
	if ok {
		t.Fatal("expected ok=false for empty positional")
	}
}

func TestParseBulkRequest_SingleNonInt(t *testing.T) {
	_, ok := parseBulkRequest([]string{"gitmap"})
	if ok {
		t.Fatal("expected ok=false for single non-integer arg")
	}
}

func TestParseBulkRequest_TooManyArgs(t *testing.T) {
	_, ok := parseBulkRequest([]string{"a", "b", "c"})
	if ok {
		t.Fatal("expected ok=false for three positional args")
	}
}
