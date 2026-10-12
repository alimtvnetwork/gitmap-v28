package cmd

import (
	"strings"
	"testing"
)

func TestIsAntigravityIdeToken(t *testing.T) {
	validTokens := []string{
		"antigravity",
		"antigravity-ide",
		"ide",
		"agy-ide",
		"ANTIGRAVITY",
		"IDE",
		"Agy-Ide",
	}
	for _, tok := range validTokens {
		if !isAntigravityIdeToken(tok) {
			t.Errorf("expected isAntigravityIdeToken(%q) = true, got false", tok)
		}
	}

	invalidTokens := []string{
		"agm",
		"agy",
		"gitmap",
		"all",
		"unknown",
		"",
	}
	for _, tok := range invalidTokens {
		if isAntigravityIdeToken(tok) {
			t.Errorf("expected isAntigravityIdeToken(%q) = false, got true", tok)
		}
	}
}

func TestIsAntigravityIdeUpdateTarget(t *testing.T) {
	tests := []struct {
		args []string
		want bool
	}{
		{[]string{"antigravity"}, true},
		{[]string{"ide"}, true},
		{[]string{"antigravity-ide"}, true},
		{[]string{"agy-ide"}, true},
		{[]string{"--dry-run", "ide"}, true},
		{[]string{"--version=2.19.1", "antigravity"}, true},
		{[]string{"-n", "agy-ide"}, true},
		{[]string{"agm"}, false},
		{[]string{"all"}, false},
		{[]string{"gitmap"}, false},
		{[]string{"--dry-run"}, false},
		{[]string{}, false},
	}

	for _, tc := range tests {
		if got := isAntigravityIdeUpdateTarget(tc.args); got != tc.want {
			t.Errorf("isAntigravityIdeUpdateTarget(%v) = %v; want %v", tc.args, got, tc.want)
		}
	}
}

func TestSuggestUpdateTargetIde(t *testing.T) {
	tests := []struct {
		token string
		want  string
	}{
		{"antigravity", "antigravity"},
		{"antigravity-ide", "antigravity"},
		{"ide", "ide"},
		{"agy-ide", "ide"},
	}

	for _, tc := range tests {
		if got := suggestUpdateTarget(tc.token); got != tc.want {
			t.Errorf("suggestUpdateTarget(%q) = %q; want %q", tc.token, got, tc.want)
		}
	}
}

func TestIsKnownUpdateTargetOrFlagIde(t *testing.T) {
	tokens := []string{"antigravity", "antigravity-ide", "ide", "agy-ide"}
	for _, tok := range tokens {
		if !isKnownUpdateTargetOrFlag(tok) {
			t.Errorf("expected isKnownUpdateTargetOrFlag(%q) = true, got false", tok)
		}
	}
}

func TestUpdateUsageHelpTextContainsIde(t *testing.T) {
	if !strings.Contains(updateUsageHelpText, "gitmap update antigravity") {
		t.Errorf("expected updateUsageHelpText to contain 'gitmap update antigravity'")
	}
	if !strings.Contains(updateUsageHelpText, "gitmap update ide") {
		t.Errorf("expected updateUsageHelpText to contain 'gitmap update ide'")
	}
}

func TestResolveUpdatePackageAntigravity(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"ide"}, "antigravity"},
		{[]string{"antigravity"}, "antigravity"},
		{[]string{"antigravity-ide"}, "antigravity"},
		{[]string{"agy-ide"}, "antigravity"},
		{[]string{"agm"}, "agm"},
		{[]string{"all"}, "gitmap"},
		{[]string{}, "gitmap"},
	}

	for _, tc := range tests {
		if got := resolveUpdatePackage(tc.args); got != tc.want {
			t.Errorf("resolveUpdatePackage(%v) = %q; want %q", tc.args, got, tc.want)
		}
	}
}

func TestRunUpdateAntigravityIdeTargetDryRun(t *testing.T) {
	args := []string{"ide", "--dry-run"}
	err := runUpdateAntigravityIdeTarget(args)
	if err != nil {
		t.Errorf("expected nil error for runUpdateAntigravityIdeTarget dry-run, got %v", err)
	}

	args2 := []string{"antigravity", "--dry-run", "--version=2.19.1"}
	err2 := runUpdateAntigravityIdeTarget(args2)
	if err2 != nil {
		t.Errorf("expected nil error for runUpdateAntigravityIdeTarget dry-run with version, got %v", err2)
	}
}
