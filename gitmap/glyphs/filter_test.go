package glyphs

import "testing"

// TestFilterSafeRewritesEmoji asserts representative glyphs collapse
// to their ASCII fallbacks under ModeSafe.
func TestFilterSafeRewritesEmoji(t *testing.T) {
	cases := map[string]string{
		"✓ done":   "v done",
		"→ next":   "-> next",
		"⚠️ warn":   "! warn",
		"📦 build": "[pkg] build",
		"plain":    "plain",
	}
	for in, want := range cases {
		got := string(Filter([]byte(in), ModeSafe))
		if got != want {
			t.Errorf("Filter(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestFilterRichPassthrough guards the zero-cost passthrough contract.
func TestFilterRichPassthrough(t *testing.T) {
	in := []byte("✓ → 📦 ⚠️")
	got := Filter(in, ModeRich)
	if string(got) != string(in) {
		t.Errorf("ModeRich must passthrough, got %q", got)
	}
}

// TestParse covers labels + auto fallback.
func TestParse(t *testing.T) {
	if Parse("rich") != ModeRich {
		t.Error("rich")
	}
	if Parse("safe") != ModeSafe {
		t.Error("safe")
	}
	if !IsValidLabel("auto") {
		t.Error("auto label invalid")
	}
}
