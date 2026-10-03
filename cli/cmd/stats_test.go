package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func verifyCandidateMatch(t *testing.T, args []string, isExpected bool) {
	t.Helper()
	isMatch := isSSHForwardCandidate(args)

	if isMatch != isExpected {
		t.Fatalf("isSSHForwardCandidate(%v) = %v; want %v", args, isMatch, isExpected)
	}
}

func TestIsSSHForwardCandidate_ValidCommands(t *testing.T) {
	candidates := []string{
		"fix-auth", "login", "host", "hosts", "nodes",
		"enable", "port", "troubleshoot", "key", "copy-id", "auth-key",
	}

	for _, cmd := range candidates {
		verifyCandidateMatch(t, []string{cmd}, true)
		verifyCandidateMatch(t, []string{cmd, "extra-arg"}, true)
	}
}

func TestIsSSHForwardCandidate_CaseInsensitive(t *testing.T) {
	cases := []string{"FIX-AUTH", "Login", "HOSTS", "TroubleShoot"}

	for _, cmd := range cases {
		verifyCandidateMatch(t, []string{cmd}, true)
	}
}

func TestIsSSHForwardCandidate_NonMatching(t *testing.T) {
	verifyCandidateMatch(t, nil, false)
	verifyCandidateMatch(t, []string{}, false)
	verifyCandidateMatch(t, []string{""}, false)
	verifyCandidateMatch(t, []string{"stats"}, false)
	verifyCandidateMatch(t, []string{"--json"}, false)
	verifyCandidateMatch(t, []string{"status"}, false)
}

func verifyRoundDuration(t *testing.T, input float64, expected int64) {
	t.Helper()
	actual := store.RoundAvgDuration(input)

	if actual != expected {
		t.Fatalf("RoundAvgDuration(%f) = %d; want %d", input, actual, expected)
	}
}

func TestStatsAvgDuration_FloatConversion(t *testing.T) {
	verifyRoundDuration(t, 26040.96023391813, 26041)
	verifyRoundDuration(t, 0.0, 0)
	verifyRoundDuration(t, 1234.4, 1234)
	verifyRoundDuration(t, 1234.5, 1235)
	verifyRoundDuration(t, 1234.6, 1235)
}
