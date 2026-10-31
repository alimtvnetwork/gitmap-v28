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

const (
	AutobotProfileID = "autobot"
	AutobotName      = "GitMap automation"
	AutobotEmail     = "info@rgitmap.com"
	AutobotType      = "automation agent"
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

	return parseAndSeedProfiles(data)
}

func parseAndSeedProfiles(data []byte) (model.GitProfileConfig, error) {
	cfg, err := unmarshalGitProfiles(data)
	if err != nil {
		return cfg, apperror.WrapSimple(err, "unmarshal git profiles:")
	}

	SeedAutobotProfile(&cfg)

	return cfg, nil
}

func handleDefaultProfilesLoad() (model.GitProfileConfig, error) {
	cfg := discoverDefaultProfiles()
	saveErr := SaveGitProfiles(cfg)

	return cfg, saveErr
}

func unmarshalGitProfiles(data []byte) (model.GitProfileConfig, error) {
	payload := extractProfilePayload(data)
	var cfg model.GitProfileConfig
	if err := json.Unmarshal(payload, &cfg); err != nil {
		return cfg, err
	}

	ensureProjectBindings(&cfg)

	return cfg, nil
}

func extractProfilePayload(data []byte) []byte {
	payload, _, err := jsonx.ExtractPayload(data)
	if err != nil {
		return data
	}

	return payload
}

func ensureProjectBindings(cfg *model.GitProfileConfig) {
	if cfg.ProjectBindings == nil {
		cfg.ProjectBindings = make(map[string]model.ProjectBinding)
	}
}

func SaveGitProfiles(cfg model.GitProfileConfig) error {
	ensureProjectBindings(&cfg)
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

// NewAutobotProfile constructs the standard automation agent Git profile.
func NewAutobotProfile() model.GitProfile {
	return model.GitProfile{
		ID:         AutobotProfileID,
		Name:       AutobotName,
		Email:      AutobotEmail,
		Provider:   "local",
		Type:       AutobotType,
		AuthMethod: "none",
		IsDefault:  false,
		UsageCount: 0,
		LastUsedAt: time.Now(),
	}
}

// SeedAutobotProfile ensures the autobot profile exists in the configuration.
func SeedAutobotProfile(cfg *model.GitProfileConfig) {
	if hasAutobotProfile(cfg.Profiles) {
		return
	}

	cfg.Profiles = append(cfg.Profiles, NewAutobotProfile())
}

func hasAutobotProfile(profiles []model.GitProfile) bool {
	for _, p := range profiles {
		if p.ID == AutobotProfileID {
			return true
		}
	}

	return false
}

func discoverDefaultProfiles() model.GitProfileConfig {
	cfg := model.GitProfileConfig{
		Profiles:        make([]model.GitProfile, 0),
		ProjectBindings: make(map[string]model.ProjectBinding),
		UpdatedAt:       time.Now(),
	}

	appendDefaultUserProfile(&cfg)
	discoverGitHubOrgs(&cfg)
	SeedAutobotProfile(&cfg)

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

	gitUserCmd := exec.Command("git", "config", "--global", "user.name")
	gitUserOut, gitErr := gitUserCmd.Output()
	if gitErr == nil && len(strings.TrimSpace(string(gitUserOut))) > 0 {
		return strings.TrimSpace(string(gitUserOut))
	}

	return ""
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
