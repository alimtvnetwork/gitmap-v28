package clonenext

import "testing"

func TestParseRepoName(t *testing.T) {
	tests := []struct {
		input   string
		base    string
		version int
		hasVer  bool
	}{
		{"macro-ahk-v11", "macro-ahk", 11, true},
		{"macro-ahk-v1", "macro-ahk", 1, true},
		{"macro-ahk", "macro-ahk", 1, false},
		{"my-app-v100", "my-app", 100, true},
	}
	for _, tt := range tests {
		p := ParseRepoName(tt.input)
		if p.BaseName != tt.base || p.CurrentVersion != tt.version || p.HasVersion != tt.hasVer {
			t.Errorf("ParseRepoName(%q) = %+v, want base=%q ver=%d has=%v",
				tt.input, p, tt.base, tt.version, tt.hasVer)
		}
	}
}

func TestResolveTarget(t *testing.T) {
	p := ParsedRepo{BaseName: "macro-ahk", CurrentVersion: 11, HasVersion: true}

	// v++ increments by 1
	v, err := ResolveTarget(p, "v++")
	if err != nil || v != 12 {
		t.Errorf("v++ = %d, %v; want 12, nil", v, err)
	}

	// v+1 is an alias for v++
	v, err = ResolveTarget(p, "v+1")
	if err != nil || v != 12 {
		t.Errorf("v+1 = %d, %v; want 12, nil", v, err)
	}

	// Explicit version
	v, err = ResolveTarget(p, "v15")
	if err != nil || v != 15 {
		t.Errorf("v15 = %d, %v; want 15, nil", v, err)
	}

	// Invalid: not a version arg
	_, err = ResolveTarget(p, "xyz")
	if err == nil {
		t.Error("expected error for invalid arg 'xyz'")
	}

	// Invalid: v0
	_, err = ResolveTarget(p, "v0")
	if err == nil {
		t.Error("expected error for v0")
	}

	// Invalid: negative
	_, err = ResolveTarget(p, "v-3")
	if err == nil {
		t.Error("expected error for v-3")
	}

	// Invalid: malformed
	_, err = ResolveTarget(p, "vabc")
	if err == nil {
		t.Error("expected error for vabc")
	}

	// Case insensitive
	v, err = ResolveTarget(p, "V++")
	if err != nil || v != 12 {
		t.Errorf("V++ = %d, %v; want 12, nil", v, err)
	}

	v, err = ResolveTarget(p, "V+1")
	if err != nil || v != 12 {
		t.Errorf("V+1 = %d, %v; want 12, nil", v, err)
	}
}

func TestTargetRepoName(t *testing.T) {
	got := TargetRepoName("macro-ahk", 12)
	if got != "macro-ahk-v12" {
		t.Errorf("got %q, want macro-ahk-v12", got)
	}
}

func TestReplaceRepoInURL(t *testing.T) {
	got := ReplaceRepoInURL("https://github.com/user/macro-ahk-v11.git", "macro-ahk-v11", "macro-ahk-v12")
	if got != "https://github.com/user/macro-ahk-v12.git" {
		t.Errorf("got %q", got)
	}

	got = ReplaceRepoInURL("git@github.com:user/macro-ahk-v11.git", "macro-ahk-v11", "macro-ahk-v12")
	if got != "git@github.com:user/macro-ahk-v12.git" {
		t.Errorf("got %q", got)
	}
}
