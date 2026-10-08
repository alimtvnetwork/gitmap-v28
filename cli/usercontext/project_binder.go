package usercontext

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var gitConfigSetter = func(dir, key, val string) error {
	cmd := exec.Command("git", "config", key, val)
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	return cmd.Run()
}

// NormalizeRepoKey converts a repo path or slug to a canonical lookup key.
func NormalizeRepoKey(repoDir string) string {
	trimmed := strings.TrimSpace(repoDir)
	if len(trimmed) == 0 || trimmed == "." {
		return resolveCurrentDirKey()
	}

	cleaned := filepath.Clean(trimmed)
	return strings.ToLower(filepath.Base(cleaned))
}

func resolveCurrentDirKey() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return strings.ToLower(filepath.Base(wd))
}

// ResolveProjectUser resolves the bound GitProfile for a given repository directory.
func ResolveProjectUser(repoDir string) (*model.GitProfile, bool) {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return nil, false
	}

	alias, hasBinding := lookupBindingAlias(cfg, repoDir)
	if !hasBinding {
		return nil, false
	}

	return findProfileByAlias(cfg.Profiles, alias)
}

func lookupBindingAlias(cfg model.GitProfileConfig, repoDir string) (string, bool) {
	if cfg.ProjectBindings == nil {
		return "", false
	}

	key := NormalizeRepoKey(repoDir)
	if binding, ok := cfg.ProjectBindings[key]; ok {
		return binding.ProfileAlias, true
	}

	if binding, ok := cfg.ProjectBindings[repoDir]; ok {
		return binding.ProfileAlias, true
	}

	return "", false
}

func findProfileByAlias(profiles []model.GitProfile, alias string) (*model.GitProfile, bool) {
	for idx := range profiles {
		if isProfileMatch(profiles[idx], alias) {
			return &profiles[idx], true
		}
	}
	return nil, false
}

func isProfileMatch(p model.GitProfile, alias string) bool {
	if p.ID == alias {
		return true
	}
	return strings.EqualFold(p.Name, alias)
}

// BindProject associates a repository directory with a Git profile alias.
func BindProject(repoDir, profileAlias string) error {
	alias := strings.TrimSpace(profileAlias)
	if len(alias) == 0 {
		return errors.New("profile alias cannot be empty")
	}

	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return err
	}

	return applyProjectBinding(&cfg, repoDir, alias)
}

func applyProjectBinding(cfg *model.GitProfileConfig, repoDir, alias string) error {
	if cfg.ProjectBindings == nil {
		cfg.ProjectBindings = make(map[string]model.ProjectBinding)
	}

	key := NormalizeRepoKey(repoDir)
	cfg.ProjectBindings[key] = model.ProjectBinding{
		RepoSlug:     key,
		ProfileAlias: alias,
		BoundAt:      time.Now(),
	}

	return store.SaveGitProfiles(*cfg)
}

// UnbindProject removes the profile binding for a repository directory.
func UnbindProject(repoDir string) error {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return err
	}

	return removeProjectBinding(&cfg, repoDir)
}

func removeProjectBinding(cfg *model.GitProfileConfig, repoDir string) error {
	if cfg.ProjectBindings == nil {
		return nil
	}

	key := NormalizeRepoKey(repoDir)
	delete(cfg.ProjectBindings, key)
	delete(cfg.ProjectBindings, repoDir)

	return store.SaveGitProfiles(*cfg)
}

// SyncProjectUser applies the bound profile's user identity to local git config.
func SyncProjectUser(repoDir string) error {
	profile, hasProfile := ResolveProjectUser(repoDir)
	if !hasProfile || profile == nil {
		return nil
	}

	dir := resolveTargetDir(repoDir)
	if err := applyGitName(dir, profile.Name); err != nil {
		return err
	}

	return applyGitEmail(dir, profile.Email)
}

func resolveTargetDir(repoDir string) string {
	cleaned := strings.TrimSpace(repoDir)
	if len(cleaned) == 0 || cleaned == "." {
		return ""
	}
	return cleaned
}

func applyGitName(dir, name string) error {
	if len(name) == 0 {
		return nil
	}
	return gitConfigSetter(dir, "user.name", name)
}

func applyGitEmail(dir, email string) error {
	if len(email) == 0 {
		return nil
	}
	return gitConfigSetter(dir, "user.email", email)
}
