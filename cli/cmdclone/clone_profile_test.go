package cmdclone

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestResolveDefaultProfile_MatchesConfigDefault(t *testing.T) {
	cfg := model.GitProfileConfig{
		Default: "work",
		Profiles: []model.GitProfile{
			{Name: "personal", ID: "p1", IsDefault: false},
			{Name: "work", ID: "p2", IsDefault: false},
		},
	}
	resolved := ResolveDefaultProfile(cfg)
	if resolved == nil || resolved.Name != "work" {
		t.Fatalf("expected work profile, got %v", resolved)
	}
}

func TestResolveDefaultProfile_FallbackToIsDefaultFlag(t *testing.T) {
	cfg := model.GitProfileConfig{
		Default: "",
		Profiles: []model.GitProfile{
			{Name: "personal", ID: "p1", IsDefault: false},
			{Name: "work", ID: "p2", IsDefault: true},
		},
	}
	resolved := ResolveDefaultProfile(cfg)
	if resolved == nil || resolved.Name != "work" {
		t.Fatalf("expected work profile via IsDefault, got %v", resolved)
	}
}

func TestCollectSecondaryProfiles_FiltersOutDefault(t *testing.T) {
	defProf := &model.GitProfile{Name: "work", ID: "p2"}
	cfg := model.GitProfileConfig{
		Profiles: []model.GitProfile{
			{Name: "personal", ID: "p1"},
			{Name: "work", ID: "p2"},
			{Name: "open-source", ID: "p3"},
		},
	}
	secondaries := CollectSecondaryProfiles(cfg, defProf)
	if len(secondaries) != 2 {
		t.Fatalf("expected 2 secondary profiles, got %d", len(secondaries))
	}
}

func TestProbeSecondaryProfiles_FallbackWinnerOn404(t *testing.T) {
	oldProber := gitLsRemoteProber
	defer func() { gitLsRemoteProber = oldProber }()

	gitLsRemoteProber = func(ctx context.Context, prof model.GitProfile, url string) (bool, error) {
		isSecondProfile := prof.Name == "secondary-2"
		if isSecondProfile {
			return true, nil
		}
		return false, errors.New("404: Repository not found")
	}

	profiles := []model.GitProfile{
		{Name: "secondary-1", ID: "s1"},
		{Name: "secondary-2", ID: "s2"},
	}
	winner, err := ProbeSecondaryProfiles(context.Background(), profiles, "https://github.com/org/repo")
	assertWinningProfile(t, winner, err, "secondary-2")
}

func assertWinningProfile(t *testing.T, winner *model.GitProfile, err error, expected string) {
	if err != nil {
		t.Fatalf("expected probe success, got err: %v", err)
	}
	if winner == nil || winner.Name != expected {
		t.Fatalf("expected %s as winner, got %v", expected, winner)
	}
}

func TestProbeSecondaryProfiles_DiagnosticFailureWhenAllFail(t *testing.T) {
	oldProber := gitLsRemoteProber
	defer func() { gitLsRemoteProber = oldProber }()

	gitLsRemoteProber = func(ctx context.Context, prof model.GitProfile, url string) (bool, error) {
		return false, errors.New("401 Unauthorized")
	}

	profiles := []model.GitProfile{
		{Name: "sec-a", ID: "sa"},
		{Name: "sec-b", ID: "sb"},
	}
	winner, err := ProbeSecondaryProfiles(context.Background(), profiles, "https://github.com/org/private-repo")
	assertDiagnosticFailure(t, winner, err)
}

func assertDiagnosticFailure(t *testing.T, winner *model.GitProfile, err error) {
	if winner != nil {
		t.Fatalf("expected nil winner, got %v", winner)
	}
	if err == nil {
		t.Fatalf("expected diagnostic error when all profiles fail")
	}
	errMsg := err.Error()
	hasSecA := strings.Contains(errMsg, "sec-a")
	hasSecB := strings.Contains(errMsg, "sec-b")
	if !hasSecA || !hasSecB {
		t.Fatalf("expected diagnostic error to list failed profiles, got: %s", errMsg)
	}
}

func TestIsAuthOrNotFoundError(t *testing.T) {
	assertAuthErrorMatches(t, errors.New("fatal: repository not found"), true)
	assertAuthErrorMatches(t, errors.New("HTTP 401 Unauthorized"), true)
	assertAuthErrorMatches(t, errors.New("Permission denied (publickey)"), true)
	assertAuthErrorMatches(t, errors.New("could not read Username"), true)
	assertAuthErrorMatches(t, errors.New("no space left on device"), false)
}

func assertAuthErrorMatches(t *testing.T, err error, expected bool) {
	isAuthErr := isAuthOrNotFoundError(err, "")
	if isAuthErr != expected {
		t.Fatalf("expected %v for error %v, got %v", expected, err, isAuthErr)
	}
}
