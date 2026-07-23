package release

import "testing"

func TestParseFullVersion(t *testing.T) {
	v, err := Parse("v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Errorf("expected 1.2.3, got %d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	if v.PreRelease != "" {
		t.Errorf("expected no pre-release, got %s", v.PreRelease)
	}
}

func TestParsePaddingMajorOnly(t *testing.T) {
	v, err := Parse("v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("expected 1.0.0, got %d.%d.%d", v.Major, v.Minor, v.Patch)
	}
}

func TestParsePaddingMajorMinor(t *testing.T) {
	v, err := Parse("v1.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 0 {
		t.Errorf("expected 1.2.0, got %d.%d.%d", v.Major, v.Minor, v.Patch)
	}
}

func TestParseNoPrefix(t *testing.T) {
	v, err := Parse("2.3.4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 2 || v.Minor != 3 || v.Patch != 4 {
		t.Errorf("expected 2.3.4, got %d.%d.%d", v.Major, v.Minor, v.Patch)
	}
}

func TestParsePreRelease(t *testing.T) {
	v, err := Parse("v1.0.0-rc.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.PreRelease != "rc.1" {
		t.Errorf("expected pre-release rc.1, got %s", v.PreRelease)
	}
	if v.IsPreRelease() == false {
		t.Error("expected IsPreRelease true")
	}
}

func TestParsePreReleaseBeta(t *testing.T) {
	v, err := Parse("v1.0.0-beta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.PreRelease != "beta" {
		t.Errorf("expected pre-release beta, got %s", v.PreRelease)
	}
}

func TestParseInvalidEmpty(t *testing.T) {
	_, err := Parse("v")
	if err == nil {
		t.Error("expected error for empty version")
	}
}

func TestParseInvalidText(t *testing.T) {
	_, err := Parse("abc")
	if err == nil {
		t.Error("expected error for non-numeric version")
	}
}

func TestParseInvalidTooManySegments(t *testing.T) {
	_, err := Parse("v1.2.3.4")
	if err == nil {
		t.Error("expected error for too many segments")
	}
}

func TestStringOutput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v1", "v1.0.0"},
		{"v1.2", "v1.2.0"},
		{"v1.2.3", "v1.2.3"},
		{"v1.0.0-rc.1", "v1.0.0-rc.1"},
	}
	for _, tt := range tests {
		v, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("Parse(%s) error: %v", tt.input, err)
		}
		if v.String() != tt.expected {
			t.Errorf("Parse(%s).String() = %s, want %s", tt.input, v.String(), tt.expected)
		}
	}
}

func TestCoreString(t *testing.T) {
	v, _ := Parse("v1.2.3")
	if v.CoreString() != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", v.CoreString())
	}

	v2, _ := Parse("v1.0.0-beta")
	if v2.CoreString() != "1.0.0-beta" {
		t.Errorf("expected 1.0.0-beta, got %s", v2.CoreString())
	}
}

func TestIsPreRelease(t *testing.T) {
	stable, _ := Parse("v1.0.0")
	if stable.IsPreRelease() {
		t.Error("stable version should not be pre-release")
	}

	pre, _ := Parse("v1.0.0-alpha.3")
	if pre.IsPreRelease() == false {
		t.Error("alpha version should be pre-release")
	}
}

func TestBumpPatch(t *testing.T) {
	v, _ := Parse("v1.2.3")
	bumped, err := Bump(v, "patch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bumped.String() != "v1.2.4" {
		t.Errorf("expected v1.2.4, got %s", bumped.String())
	}
}

func TestBumpMinor(t *testing.T) {
	v, _ := Parse("v1.2.3")
	bumped, err := Bump(v, "minor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bumped.String() != "v1.3.0" {
		t.Errorf("expected v1.3.0, got %s", bumped.String())
	}
}

func TestBumpMajor(t *testing.T) {
	v, _ := Parse("v0.9.1")
	bumped, err := Bump(v, "major")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bumped.String() != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", bumped.String())
	}
}

func TestBumpInvalidLevel(t *testing.T) {
	v, _ := Parse("v1.0.0")
	_, err := Bump(v, "invalid")
	if err == nil {
		t.Error("expected error for invalid bump level")
	}
}

func TestGreaterThanMajor(t *testing.T) {
	a, _ := Parse("v2.0.0")
	b, _ := Parse("v1.9.9")
	if a.GreaterThan(b) == false {
		t.Error("v2.0.0 should be greater than v1.9.9")
	}
	if b.GreaterThan(a) {
		t.Error("v1.9.9 should not be greater than v2.0.0")
	}
}

func TestGreaterThanMinor(t *testing.T) {
	a, _ := Parse("v1.3.0")
	b, _ := Parse("v1.2.9")
	if a.GreaterThan(b) == false {
		t.Error("v1.3.0 should be greater than v1.2.9")
	}
}

func TestGreaterThanPatch(t *testing.T) {
	a, _ := Parse("v1.2.4")
	b, _ := Parse("v1.2.3")
	if a.GreaterThan(b) == false {
		t.Error("v1.2.4 should be greater than v1.2.3")
	}
}

func TestGreaterThanEqual(t *testing.T) {
	a, _ := Parse("v1.2.3")
	b, _ := Parse("v1.2.3")
	if a.GreaterThan(b) {
		t.Error("equal versions should not be greater")
	}
}

func TestGreaterThanStableBeatsPreRelease(t *testing.T) {
	stable, _ := Parse("v1.0.0")
	pre, _ := Parse("v1.0.0-rc.1")
	if stable.GreaterThan(pre) == false {
		t.Error("v1.0.0 should be greater than v1.0.0-rc.1")
	}
	if pre.GreaterThan(stable) {
		t.Error("v1.0.0-rc.1 should not be greater than v1.0.0")
	}
}
