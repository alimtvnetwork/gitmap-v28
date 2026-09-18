package cmdssh

import (
	"context"
	"testing"
)

func TestDispatchPrimarySSH_MatchedSubcommands(t *testing.T) {
	ctx := context.Background()
	matchedSubs := []string{
		"check", "health", "ping",
		"scan",
		"exec", "se",
		"join", "sj",
		"install", "i",
		"update", "u",
		"agy", "code",
		"compare", "matrix",
		"profiles", "profile", "p",
	}

	for _, sub := range matchedSubs {
		res := dispatchPrimarySSH(ctx, sub, []string{"--help"}, nil)
		if !res.IsMatched() {
			t.Errorf("dispatchPrimarySSH(%q) expected matched wrapper, got unmatched", sub)
		}
	}
}

func TestDispatchPrimarySSH_UnmatchedSubcommands(t *testing.T) {
	ctx := context.Background()
	unmatchedSubs := []string{
		"unknown-sub",
		"nonexistent",
		"foo-bar-xyz",
	}

	for _, sub := range unmatchedSubs {
		res := dispatchPrimarySSH(ctx, sub, []string{}, nil)
		if res.IsMatched() {
			t.Errorf("dispatchPrimarySSH(%q) expected unmatched wrapper, got matched", sub)
		}
	}
}

func TestRunSSHProfile_NilRunner(t *testing.T) {
	prev := ProfileRunner
	ProfileRunner = nil
	defer func() { ProfileRunner = prev }()

	if err := runSSHProfile([]string{}); err != nil {
		t.Errorf("expected nil error when ProfileRunner is nil, got: %v", err)
	}
}

func TestRunSSHProfile_CustomRunner(t *testing.T) {
	prev := ProfileRunner
	wasCalled := false
	ProfileRunner = func(args []string) error {
		wasCalled = true
		return nil
	}
	defer func() { ProfileRunner = prev }()

	if err := runSSHProfile([]string{"dev"}); err != nil || !wasCalled {
		t.Errorf("expected custom runner to execute, err=%v, called=%v", err, wasCalled)
	}
}
