package store

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func GitProfilesPath() string {
	return filepath.Join(BinaryDataDir(), "git_profiles.json")
}

func LoadGitProfiles() (model.GitProfileConfig, error) {
	path := GitProfilesPath()
	data, err := os.ReadFile(path)
	if err != nil {
		cfg := discoverDefaultProfiles()
		saveErr := SaveGitProfiles(cfg)
		return cfg, saveErr
	}
	var cfg model.GitProfileConfig
	if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
		return cfg, apperror.WrapSimple(unmarshalErr, "unmarshal git profiles:")
	}
	return cfg, nil
}

func SaveGitProfiles(cfg model.GitProfileConfig) error {
	cfg.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal git profiles:")
	}
	if mkErr := os.MkdirAll(BinaryDataDir(), 0755); mkErr != nil {
		return apperror.WrapSimple(mkErr, "create data folder:")
	}
	return os.WriteFile(GitProfilesPath(), data, 0644)
}

func discoverDefaultProfiles() model.GitProfileConfig {
	cfg := model.GitProfileConfig{
		Profiles:  make([]model.GitProfile, 0),
		UpdatedAt: time.Now(),
	}
	user := detectGitHubUser()
	if user != "" {
		cfg.Profiles = append(cfg.Profiles, model.GitProfile{
			ID:         "prof_1",
			Name:       user,
			Provider:   "github",
			Type:       "user",
			AuthMethod: "gh-cli",
			IsDefault:  true,
			UsageCount: 1,
			LastUsedAt: time.Now(),
		})
		cfg.Active = user
		cfg.Default = user
	}
	discoverGitHubOrgs(&cfg)
	return cfg
}

func detectGitHubUser() string {
	cmd := exec.Command("gh", "api", "user", "--jq", ".login")
	out, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return strings.TrimSpace(string(out))
	}
	gitUserCmd := exec.Command("git", "config", "user.name")
	gitUserOut, gitErr := gitUserCmd.Output()
	if gitErr == nil {
		return strings.TrimSpace(string(gitUserOut))
	}
	return "default"
}

func discoverGitHubOrgs(cfg *model.GitProfileConfig) {
	cmd := exec.Command("gh", "api", "user/orgs", "--jq", ".[].login")
	out, err := cmd.Output()
	if err != nil {
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for idx, line := range lines {
		org := strings.TrimSpace(line)
		appendOrgProfile(cfg, org, idx+2)
	}
}

func appendOrgProfile(cfg *model.GitProfileConfig, org string, idx int) {
	if len(org) == 0 {
		return
	}
	cfg.Profiles = append(cfg.Profiles, model.GitProfile{
		ID:         org,
		Name:       org,
		Provider:   "github",
		Type:       "organization",
		AuthMethod: "gh-cli",
		UsageCount: 0,
		LastUsedAt: time.Now(),
	})
}
