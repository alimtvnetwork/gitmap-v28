package cmd

import (
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestPickProfileBySequenceOrName(t *testing.T) {
	profiles := []model.GitProfile{
		{ID: "prof_1", Name: "alimtvnetwork", Provider: "github", Type: "user", IsDefault: true, UsageCount: 5},
		{ID: "prof_2", Name: "riseup-asia", Provider: "github", Type: "organization", IsDefault: false, UsageCount: 2},
	}

	// Test by sequence number "1"
	idx1, prof1, err1 := pickProfileBySequenceOrName(profiles, "1")
	if err1 != nil || idx1 != 0 || prof1.Name != "alimtvnetwork" {
		t.Fatalf("expected profile 1 'alimtvnetwork', got: %v, %v, %v", idx1, prof1, err1)
	}

	// Test by sequence number "2"
	idx2, prof2, err2 := pickProfileBySequenceOrName(profiles, "2")
	if err2 != nil || idx2 != 1 || prof2.Name != "riseup-asia" {
		t.Fatalf("expected profile 2 'riseup-asia', got: %v, %v, %v", idx2, prof2, err2)
	}

	// Test by name
	idxName, profName, errName := pickProfileBySequenceOrName(profiles, "riseup-asia")
	if errName != nil || idxName != 1 || profName.Type != "organization" {
		t.Fatalf("expected profile 'riseup-asia', got: %v, %v, %v", idxName, profName, errName)
	}

	// Test not found
	_, _, errNotFound := pickProfileBySequenceOrName(profiles, "nonexistent")
	if errNotFound == nil {
		t.Fatalf("expected error for nonexistent profile")
	}
}

func TestNormalizeCreateArgs(t *testing.T) {
	argsWithRepo := []string{"repo", "my-test-service", "--private"}
	norm1 := normalizeCreateArgs(argsWithRepo)
	if len(norm1) != 2 || norm1[0] != "my-test-service" {
		t.Fatalf("expected ['my-test-service', '--private'], got: %v", norm1)
	}

	argsWithoutRepo := []string{"direct-service", "--public"}
	norm2 := normalizeCreateArgs(argsWithoutRepo)
	if len(norm2) != 2 || norm2[0] != "direct-service" {
		t.Fatalf("expected ['direct-service', '--public'], got: %v", norm2)
	}
}

func TestResolveBackupRepoSlug(t *testing.T) {
	prof := model.GitProfile{Name: "alimtvnetwork", Provider: "github"}

	// Default backup slug
	slug1 := resolveBackupRepoSlug(prof, nil)
	if slug1 != "alimtvnetwork/gitmap-cloud-backup" {
		t.Fatalf("expected 'alimtvnetwork/gitmap-cloud-backup', got: %s", slug1)
	}

	// Custom repo slug override
	slug2 := resolveBackupRepoSlug(prof, []string{"--repo", "custom-org/backup-vault"})
	if slug2 != "custom-org/backup-vault" {
		t.Fatalf("expected 'custom-org/backup-vault', got: %s", slug2)
	}
}

func TestApplyDefaultProfile(t *testing.T) {
	cfg := model.GitProfileConfig{
		Profiles: []model.GitProfile{
			{Name: "user1", IsDefault: true, LastUsedAt: time.Now()},
			{Name: "org1", IsDefault: false, LastUsedAt: time.Now()},
		},
		Default: "user1",
		Active:  "user1",
	}

	target := cfg.Profiles[1]
	for i := range cfg.Profiles {
		cfg.Profiles[i].IsDefault = (cfg.Profiles[i].Name == target.Name)
	}
	cfg.Default = target.Name
	cfg.Active = target.Name

	if cfg.Default != "org1" || !cfg.Profiles[1].IsDefault || cfg.Profiles[0].IsDefault {
		t.Fatalf("expected org1 as default, got: %+v", cfg)
	}
}
