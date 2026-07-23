package clonefrom

// Parser tests dedicated to the per-row `checkout` field added in
// the checkout-mode feature. Kept in a separate file so the existing
// parse_test.go stays focused on the original column set and the
// new tests are easy to delete if the feature is ever rolled back.

import (
	"strings"
	"testing"
)

// TestParseFile_JSON_CheckoutField confirms the JSON parser reads
// the optional `checkout` key, lower-cases it, and rejects bogus
// values at row level (so the user gets a row-pointing error
// instead of a confusing "git clone failed" later).
func TestParseFile_JSON_CheckoutField(t *testing.T) {
	body := `[
  {"url": "https://x/y.git", "checkout": "Skip"},
  {"url": "https://x/z.git", "checkout": "force"}
]`
	plan, err := ParseFile(writeTemp(t, "co.json", body))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if plan.Rows[0].Checkout != "skip" {
		t.Errorf("row0 checkout = %q, want skip (lower-cased)",
			plan.Rows[0].Checkout)
	}
	if plan.Rows[1].Checkout != "force" {
		t.Errorf("row1 checkout = %q, want force", plan.Rows[1].Checkout)
	}
}

// TestParseFile_JSON_CheckoutFieldRejectsBogus pins the parse-time
// failure path: a typo in the manifest fails the whole load with a
// row-pointing error message.
func TestParseFile_JSON_CheckoutFieldRejectsBogus(t *testing.T) {
	body := `[{"url": "https://x/y.git", "checkout": "yolo"}]`
	_, err := ParseFile(writeTemp(t, "bad.json", body))
	if err == nil {
		t.Fatalf("ParseFile accepted bogus checkout")
	}
	if !strings.Contains(err.Error(), "yolo") {
		t.Errorf("error %q does not mention bad value", err.Error())
	}
}

// TestParseFile_CSV_CheckoutColumn confirms the CSV header recognizer
// picks up the new column at any position (here: between dest and
// branch) and lower-cases the value.
func TestParseFile_CSV_CheckoutColumn(t *testing.T) {
	body := "url,dest,Checkout,branch,depth\n" +
		"https://x/y.git,,Skip,,\n" +
		"https://x/z.git,custom,FORCE,main,1\n"
	plan, err := ParseFile(writeTemp(t, "co.csv", body))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if plan.Rows[0].Checkout != "skip" {
		t.Errorf("row0 checkout = %q, want skip", plan.Rows[0].Checkout)
	}
	if plan.Rows[1].Checkout != "force" || plan.Rows[1].Branch != "main" {
		t.Errorf("row1 = %+v", plan.Rows[1])
	}
}

// TestMergeRows_CheckoutLaterWins covers the dedup-merge contract
// for the new field: a later row with an explicit checkout
// overrides an earlier row's empty/old value, matching how
// branch/depth already behave.
func TestMergeRows_CheckoutLaterWins(t *testing.T) {
	first := Row{URL: "u", Checkout: "auto"}
	later := Row{URL: "u", Checkout: "skip"}
	out := mergeRows(first, later)
	if out.Checkout != "skip" {
		t.Errorf("merged checkout = %q, want skip", out.Checkout)
	}
	// Empty later does NOT clobber a non-empty first.
	out2 := mergeRows(Row{URL: "u", Checkout: "force"}, Row{URL: "u"})
	if out2.Checkout != "force" {
		t.Errorf("empty later clobbered: got %q, want force", out2.Checkout)
	}
}
