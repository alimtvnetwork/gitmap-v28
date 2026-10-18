// Package cmdclone — clone_profile_fallback.go handles multi-login fallback during direct cloning.
package cmdclone

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DirectCloneParams encapsulates options for cloning a single repository.
type DirectCloneParams struct {
	URL          string
	FolderName   string
	GHDesktop    bool
	NoReplace    bool
	Output       string
	NoVSCodeSync bool
	IsClean      bool
	IsJSON       bool
	UseSSH       bool
	Profile      string
}

func runDirectCloneWithProfileFallback(params DirectCloneParams, url, absPath string) (*model.GitProfile, error) {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return runFallbackWithoutProfiles(params, url, absPath)
	}

	if len(params.Profile) > 0 {
		return executeCloneWithSpecificProfile(params, url, absPath, cfg, params.Profile)
	}

	defaultProf := ResolveDefaultProfile(cfg)
	return executeCloneWithDefaultAndFallback(params, url, absPath, cfg, defaultProf)
}

func runFallbackWithoutProfiles(params DirectCloneParams, url, absPath string) (*model.GitProfile, error) {
	err := executeSingleCloneCommand(url, absPath, params.NoReplace)
	return nil, err
}

func executeCloneWithSpecificProfile(
	params DirectCloneParams,
	url, absPath string,
	cfg model.GitProfileConfig,
	alias string,
) (*model.GitProfile, error) {
	prof := findProfileByName(cfg.Profiles, alias)
	if prof == nil {
		return nil, apperror.NewNotFoundError(fmt.Sprintf("profile not found: %s", alias))
	}
	cloneURL, restore := ApplyProfileCredentials(*prof, url)
	defer restore()
	err := executeSingleCloneCommand(cloneURL, absPath, params.NoReplace)
	return prof, err
}

func executeCloneWithDefaultAndFallback(
	params DirectCloneParams,
	url, absPath string,
	cfg model.GitProfileConfig,
	defaultProf *model.GitProfile,
) (*model.GitProfile, error) {
	cloneURL := url
	var restore func()
	if defaultProf != nil {
		cloneURL, restore = ApplyProfileCredentials(*defaultProf, url)
		defer restore()
	}
	err := executeSingleCloneCommand(cloneURL, absPath, params.NoReplace)
	if err == nil {
		return defaultProf, nil
	}
	if !isAuthOrNotFoundError(err, "") {
		return nil, err
	}
	return attemptSecondaryProfilesFallback(params, url, absPath, cfg, defaultProf)
}

func attemptSecondaryProfilesFallback(
	params DirectCloneParams,
	url, absPath string,
	cfg model.GitProfileConfig,
	defaultProf *model.GitProfile,
) (*model.GitProfile, error) {
	secondaries := CollectSecondaryProfiles(cfg, defaultProf)
	if len(secondaries) == 0 {
		return nil, apperror.NewNotFoundError(fmt.Sprintf("primary profile failed and no secondary profiles configured for %s", url))
	}
	winner, probeErr := ProbeSecondaryProfiles(context.Background(), secondaries, url)
	if probeErr != nil {
		return nil, probeErr
	}
	cleanDirectCloneTarget(absPath, true)
	fmt.Printf("✓ Fallback: verified access via profile %q, cloning...\n", winner.Name)
	winURL, winRestore := ApplyProfileCredentials(*winner, url)
	defer winRestore()
	cloneErr := executeSingleCloneCommand(winURL, absPath, params.NoReplace)
	return winner, cloneErr
}

func executeSingleCloneCommand(url, absPath string, noReplace bool) error {
	if noReplace {
		return runCloneCommand(url, absPath)
	}
	_, err := cloneReplacing(url, absPath)
	return err
}
