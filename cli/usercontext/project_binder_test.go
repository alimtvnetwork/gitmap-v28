package usercontext

import (
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestBindAndResolveProject(t *testing.T) {
	tempDir := t.TempDir()
	store.SetBinaryDataDirForTesting(tempDir)
	defer store.SetBinaryDataDirForTesting("")

	seedProfilesConfig(t)

	if err := BindProject("/projects/my-service", "work-profile"); err != nil {
		t.Fatalf("unexpected bind error: %v", err)
	}

	profile, hasProfile := ResolveProjectUser("/projects/my-service")
	if !hasProfile || profile == nil {
		t.Fatalf("expected profile to resolve")
	}
	if profile.Name != "Work User" {
		t.Fatalf("expected Work User, got %s", profile.Name)
	}

	slugProfile, hasSlug := ResolveProjectUser("my-service")
	if !hasSlug || slugProfile == nil {
		t.Fatalf("expected profile to resolve by slug")
	}
}

func TestUnbindProject(t *testing.T) {
	tempDir := t.TempDir()
	store.SetBinaryDataDirForTesting(tempDir)
	defer store.SetBinaryDataDirForTesting("")

	seedProfilesConfig(t)
	_ = BindProject("repo-alpha", "work-profile")

	if err := UnbindProject("repo-alpha"); err != nil {
		t.Fatalf("unexpected unbind error: %v", err)
	}

	profile, hasProfile := ResolveProjectUser("repo-alpha")
	if hasProfile || profile != nil {
		t.Fatalf("expected profile to be unbound")
	}
}

func TestSyncProjectUser(t *testing.T) {
	tempDir := t.TempDir()
	store.SetBinaryDataDirForTesting(tempDir)
	defer store.SetBinaryDataDirForTesting("")

	seedProfilesConfig(t)
	_ = BindProject("sync-repo", "work-profile")

	var appliedKey, appliedVal string
	origSetter := gitConfigSetter
	gitConfigSetter = func(dir, key, val string) error {
		appliedKey = key
		appliedVal = val
		return nil
	}
	defer func() { gitConfigSetter = origSetter }()

	if err := SyncProjectUser("sync-repo"); err != nil {
		t.Fatalf("unexpected sync error: %v", err)
	}
	if appliedKey != "user.email" || appliedVal != "work@company.com" {
		t.Fatalf("expected user.email set to work@company.com, got %s=%s", appliedKey, appliedVal)
	}
}

func seedProfilesConfig(t *testing.T) {
	cfg := model.GitProfileConfig{
		Profiles: []model.GitProfile{
			{
				ID:         "work-profile",
				Name:       "Work User",
				Email:      "work@company.com",
				Provider:   "github",
				Type:       "user",
				AuthMethod: "gh-cli",
				LastUsedAt: time.Now(),
			},
		},
		Active:    "work-profile",
		UpdatedAt: time.Now(),
	}

	if err := store.SaveGitProfiles(cfg); err != nil {
		t.Fatalf("failed to seed profiles: %v", err)
	}
}
