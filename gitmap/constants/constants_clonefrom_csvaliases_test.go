package constants

import "testing"

// TestCanonicalCSVColumn pins the alias map: every accepted spelling
// must resolve to the documented canonical column, and unknown
// headers must return "" (so the parser can skip them silently
// rather than misclassify them as a known column).
func TestCanonicalCSVColumn(t *testing.T) {
	cases := []struct {
		raw, want string
	}{
		// canonical names round-trip
		{"url", CSVColumnURL},
		{"dest", CSVColumnDest},
		{"branch", CSVColumnBranch},
		{"depth", CSVColumnDepth},
		{"checkout", CSVColumnCheckout},

		// case + whitespace + separator variants
		{"URL", CSVColumnURL},
		{"  Url  ", CSVColumnURL},
		{"HTTPS_URL", CSVColumnURL},
		{"https-url", CSVColumnURL},
		{"Https URL", CSVColumnURL},
		{"httpsURL", CSVColumnURL},

		// dest aliases
		{"relpath", CSVColumnDest},
		{"RelPath", CSVColumnDest},
		{"path", CSVColumnDest},
		{"folder", CSVColumnDest},

		// branch / depth / checkout aliases
		{"ref", CSVColumnBranch},
		{"clone-depth", CSVColumnDepth},
		{"CheckoutMode", CSVColumnCheckout},

		// unknown → empty so parser ignores extra columns
		{"", ""},
		{"notes", ""},
		{"sha", ""},
		{"comment", ""},
	}
	for _, tc := range cases {
		if got := CanonicalCSVColumn(tc.raw); got != tc.want {
			t.Errorf("CanonicalCSVColumn(%q) = %q, want %q", tc.raw, got, tc.want)
		}
	}
}

// TestCSVColumnAliasesNoConflict asserts every value in the alias
// map is one of the five canonical column constants. Guards against
// typos that would silently route a header to a non-existent column.
func TestCSVColumnAliasesNoConflict(t *testing.T) {
	valid := map[string]bool{
		CSVColumnURL: true, CSVColumnDest: true, CSVColumnBranch: true,
		CSVColumnDepth: true, CSVColumnCheckout: true,
	}
	for alias, canonical := range CSVColumnAliases {
		if !valid[canonical] {
			t.Errorf("alias %q maps to unknown canonical %q", alias, canonical)
		}
	}
}
