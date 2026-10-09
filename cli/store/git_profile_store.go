package store

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/jsonx"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func GitProfilesPath() string {
	return filepath.Join(BinaryDataDir(), "git_profiles.json")
}

func LoadGitProfiles() (model.GitProfileConfig, error) {
	path := GitProfilesPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return handleDefaultProfilesLoad()
	}

	cfg, unmarshalErr := unmarshalGitProfiles(data)
	if unmarshalErr != nil {
		return cfg, apperror.WrapSimple(unmarshalErr, "unmarshal git profiles:")
	}

	return cfg, nil
}

func handleDefaultProfilesLoad() (model.GitProfileConfig, error) {
	cfg := discoverDefaultProfiles()
	saveErr := SaveGitProfiles(cfg)

	return cfg, saveErr
}

func unmarshalGitProfiles(data []byte) (model.GitProfileConfig, error) {
	payload, _, extractErr := jsonx.ExtractPayload(data)
	if extractErr != nil {
		payload = data
	}

	var cfg model.GitProfileConfig
	if err := json.Unmarshal(payload, &cfg); err != nil {
		return cfg, err
	}

	if cfg.ProjectBindings == nil {
		cfg.ProjectBindings = make(map[string]model.ProjectBinding)
	}

	return cfg, nil
}

func SaveGitProfiles(cfg model.GitProfileConfig) error {
	if cfg.ProjectBindings == nil {
		cfg.ProjectBindings = make(map[string]model.ProjectBinding)
	}
	cfg.UpdatedAt = time.Now()
	envelope := jsonx.NewEnvelope(
		jsonx.TypeGitProfiles,
		GitProfilesPath(),
		"gitmap profile save",
		"1.0",
		cfg,
	)

	return writeProfilesEnvelope(envelope)
}

func writeProfilesEnvelope(envelope jsonx.Envelope[model.GitProfileConfig]) error {
	data, err := json.MarshalIndent(envelope, "", "  ")
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
		Profiles:        make([]model.GitProfile, 0),
		ProjectBindings: make(map[string]model.ProjectBinding),
		UpdatedAt:       time.Now(),
	}

	appendDefaultUserProfile(&cfg)
	discoverGitHubOrgs(&cfg)

	return cfg
}

func appendDefaultUserProfile(cfg *model.GitProfileConfig) {
	user := detectGitHubUser()
	if len(user) == 0 {
		return
	}

	cfg.Profiles = append(cfg.Profiles, newDefaultUserProfile(user))
	cfg.Active = user
	cfg.Default = user
}

func newDefaultUserProfile(user string) model.GitProfile {
	return model.GitProfile{
		ID:         "prof_1",
		Name:       user,
		Provider:   "github",
		Type:       "user",
		AuthMethod: "gh-cli",
		IsDefault:  true,
		UsageCount: 1,
		LastUsedAt: time.Now(),
	}
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
