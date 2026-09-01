// Package cmd — chrome_exec_test.go: unit tests for Chrome binary detection,
// launch options parsing, and multi-profile URL mappings.
package cmd

import (
	"testing"
)

func TestParseChromeLaunchArgsSingle(t *testing.T) {
	args := []string{"https://github.com", "--profile=Profile 1", "--incognito"}
	opts, err := parseChromeLaunchArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.IsIncog {
		t.Errorf("expected IsIncog to be true")
	}
	if len(opts.Targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(opts.Targets))
	}
	if opts.Targets[0].Profile != "Profile 1" {
		t.Errorf("expected Profile 1, got %s", opts.Targets[0].Profile)
	}
	if len(opts.Targets[0].URLs) != 1 || opts.Targets[0].URLs[0] != "https://github.com" {
		t.Errorf("expected https://github.com, got %v", opts.Targets[0].URLs)
	}
}

func TestParseChromeLaunchArgsMultiMapping(t *testing.T) {
	args := []string{"Profile 1=https://github.com,Profile 2=https://google.com"}
	opts, err := parseChromeLaunchArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(opts.Targets))
	}
	if opts.Targets[0].Profile != "Profile 1" || opts.Targets[0].URLs[0] != "https://github.com" {
		t.Errorf("unexpected target 0: %+v", opts.Targets[0])
	}
	if opts.Targets[1].Profile != "Profile 2" || opts.Targets[1].URLs[0] != "https://google.com" {
		t.Errorf("unexpected target 1: %+v", opts.Targets[1])
	}
}

func TestBuildChromeCmdArgs(t *testing.T) {
	opts := chromeLaunchOptions{
		IsIncog:  true,
		IsNewWin: true,
		AppURL:   "https://app.com",
	}
	args := buildChromeCmdArgs("Default", []string{"https://example.com"}, opts)
	hasProf := false
	hasIncog := false
	hasWin := false
	hasApp := false
	hasURL := false
	for _, a := range args {
		if a == "--profile-directory=Default" {
			hasProf = true
		}
		if a == "--incognito" {
			hasIncog = true
		}
		if a == "--new-window" {
			hasWin = true
		}
		if a == "--app=https://app.com" {
			hasApp = true
		}
		if a == "https://example.com" {
			hasURL = true
		}
	}
	if !hasProf || !hasIncog || !hasWin || !hasApp || !hasURL {
		t.Errorf("missing flags in built args: %v", args)
	}
}
