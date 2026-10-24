package enums

import "testing"

func TestCommitVerbType_Properties(t *testing.T) {
	verbs := []CommitVerbType{
		CommitVerbFeature,
		CommitVerbBug,
		CommitVerbRefactor,
		CommitVerbChore,
		CommitVerbRelease,
		CommitVerbFlat,
		CommitVerbPullPush,
	}

	for _, v := range verbs {
		if !v.IsValid() {
			t.Errorf("expected %s to be valid", v)
		}

		if v.String() != string(v) {
			t.Errorf("expected String() == string(v) for %s", v)
		}
	}

	if !CommitVerbFeature.IsFeature() {
		t.Errorf("expected CommitVerbFeature.IsFeature() to be true")
	}

	if CommitVerbBug.IsFeature() {
		t.Errorf("expected CommitVerbBug.IsFeature() to be false")
	}

	if !CommitVerbBug.IsBug() {
		t.Errorf("expected CommitVerbBug.IsBug() to be true")
	}

	if CommitVerbFeature.Prefix() != "Feature: " {
		t.Errorf("expected 'Feature: ', got %s", CommitVerbFeature.Prefix())
	}

	if CommitVerbBug.Prefix() != "Bug: " {
		t.Errorf("expected 'Bug: ', got %s", CommitVerbBug.Prefix())
	}

	if CommitVerbFeature.UppercaseToken() != "FEAT" {
		t.Errorf("expected 'FEAT', got %s", CommitVerbFeature.UppercaseToken())
	}

	if CommitVerbBug.UppercaseToken() != "FIX" {
		t.Errorf("expected 'FIX', got %s", CommitVerbBug.UppercaseToken())
	}

	invalid := CommitVerbType("unknown")
	if invalid.IsValid() {
		t.Errorf("expected invalid verb to return false")
	}
}

func TestConventionalVerbType_Properties(t *testing.T) {
	verbs := AllConventionalVerbs()
	if len(verbs) != 5 {
		t.Fatalf("expected 5 conventional verbs, got %d", len(verbs))
	}

	for _, v := range verbs {
		if !v.IsValid() {
			t.Errorf("expected %s to be valid", v)
		}

		if v.String() == "" {
			t.Errorf("expected non-empty string for %s", v)
		}
	}

	invalid := ConventionalVerbType("INVALID")
	if invalid.IsValid() {
		t.Errorf("expected invalid conventional verb to return false")
	}
}

func TestCommitFormatType_Properties(t *testing.T) {
	formats := []CommitFormatType{
		FormatStandardBracket,
		FormatHyphenSeparated,
		FormatConventional,
	}

	for _, f := range formats {
		if !f.IsValid() {
			t.Errorf("expected format %d to be valid", f)
		}

		if f.String() == "Unknown" {
			t.Errorf("expected known string for format %d", f)
		}
	}

	if !FormatHyphenSeparated.IsHyphenSeparated() {
		t.Errorf("expected FormatHyphenSeparated.IsHyphenSeparated() to be true")
	}

	if !FormatConventional.IsConventional() {
		t.Errorf("expected FormatConventional.IsConventional() to be true")
	}

	invalid := CommitFormatType(99)
	if invalid.IsValid() {
		t.Errorf("expected invalid format to return false")
	}
}
