package cmd

// Regression coverage for the bare-base rewrite branch added in
// v5.8.0: when v1 is in the target set, `applyAllTargets` must ALSO
// substitute standalone `{base}` occurrences (the pre-versioned
// remote name) — bounded by word boundaries so it never clobbers
// `{base}-vN`, `{base}.js`, or substrings of unrelated identifiers.

import "testing"

func TestApplyAllTargets_BareBase_V1To2(t *testing.T) {
	body := "url=https://github.com/x/img-pdf go to img-pdf/main\nalso img-pdf-v1 and img-pdf-v2 stay correct"
	got, count := applyAllTargets(body, "img-pdf", 2, []int{1})
	want := "url=https://github.com/x/img-pdf-v2 go to img-pdf-v2/main\nalso img-pdf-v2 and img-pdf-v2 stay correct"
	if got != want {
		t.Fatalf("rewrite mismatch.\n got:  %q\n want: %q", got, want)
	}
	// 2 bare-base hits + 1 -v1 hit = 3
	if count != 3 {
		t.Fatalf("expected 3 replacements, got %d", count)
	}
}

func TestApplyBareBase_WordBoundaryGuards(t *testing.T) {
	// None of these should be rewritten — each fails a boundary check.
	cases := []string{
		"myimg-pdf trailing-word", // prev byte is letter
		"img-pdfx trailing-word",  // next byte is letter
		"img-pdf.js extension",    // next byte is '.'
		"img-pdf_alt underscore",  // next byte is '_'
		"img-pdf-v2 already-versioned",
		"img-pdf-tools dashed",
	}
	for _, in := range cases {
		got, n := applyBareBase(in, "img-pdf", 2)
		if got != in || n != 0 {
			t.Errorf("boundary violation for %q -> %q (n=%d)", in, got, n)
		}
	}
}

func TestApplyBareBase_SkippedWhenV1NotInTargets(t *testing.T) {
	body := "img-pdf bare and img-pdf-v2 versioned"
	got, _ := applyAllTargets(body, "img-pdf", 3, []int{2})
	// v1 not in targets -> bare-base sweep MUST NOT run.
	want := "img-pdf bare and img-pdf-v3 versioned"
	if got != want {
		t.Fatalf("bare-base ran without v1 target.\n got:  %q\n want: %q", got, want)
	}
}

// v5.38.0: bare-base sweep is now restricted to the v1→v2 transition.
// At v3+ the bare token (e.g. `gitmap`) is almost always an unrelated
// identifier (binary name, package, brand) and must be left alone.
func TestApplyAllTargets_BareBase_SkippedAtV3Plus(t *testing.T) {
	body := "url=https://github.com/x/acme and acme-v1 plus acme-v2"
	got, _ := applyAllTargets(body, "acme", 3, []int{1, 2})
	want := "url=https://github.com/x/acme and acme-v3 plus acme-v3"
	if got != want {
		t.Fatalf("bare-base ran at v3 (must only run at v1→v2).\n got:  %q\n want: %q", got, want)
	}
}

func TestApplyAllTargets_BareBase_SkippedAtV4WithV1InTargets(t *testing.T) {
	body := "acme and acme-v1 and acme-v2 and acme-v3"
	got, _ := applyAllTargets(body, "acme", 4, []int{1, 2, 3})
	want := "acme and acme-v4 and acme-v4 and acme-v4"
	if got != want {
		t.Fatalf("bare-base must be skipped when current>2 even if v1 targeted.\n got:  %q\n want: %q", got, want)
	}
}
