package cmd

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func resolveCreationProfile(args []string) (model.GitProfile, error) {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return model.GitProfile{}, apperror.WrapSimple(err, "load profiles:")
	}

	req := extractFlagVal(args, "--profile")
	if req == "" {
		req = extractFlagVal(args, "--org")
	}

	if req != "" {
		_, p, findErr := pickProfileBySequenceOrName(cfg.Profiles, req)

		return p, findErr
	}

	return pickDefaultProfile(cfg), nil
}

func pickDefaultProfile(cfg model.GitProfileConfig) model.GitProfile {
	for _, p := range cfg.Profiles {
		if p.IsDefault || p.Name == cfg.Default {
			return p
		}
	}

	if len(cfg.Profiles) > 0 {
		return cfg.Profiles[0]
	}

	return model.GitProfile{Name: "default", Provider: "github", Type: "user"}
}

func recordProfileUsage(prof model.GitProfile) {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return
	}

	for i := range cfg.Profiles {
		if cfg.Profiles[i].Name == prof.Name {
			cfg.Profiles[i].UsageCount++
			cfg.Profiles[i].LastUsedAt = time.Now()
			break
		}
	}

	_ = store.SaveGitProfiles(cfg)
}
