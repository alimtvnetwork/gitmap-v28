package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func executeCreateRepo(args []string, defaultLocal bool) error {
	params, parseErr := parseCreateParams(args, defaultLocal)
	if parseErr != nil {
		return parseErr
	}

	if initErr := initLocalRepo(params); initErr != nil {
		return initErr
	}

	remoteURL, pushErr := pushRemoteRepo(params)
	if pushErr != nil {
		return pushErr
	}

	absDir, _ := filepath.Abs(params.LocalDir)
	if params.IsNoDesktop {
		workspacesync.SyncWithoutDesktop(absDir, params.Name)
	} else {
		workspacesync.SyncAll(absDir, params.Name)
	}
	recordProfileUsage(params.Profile)
	if params.IsCD {
		WriteShellHandoff(absDir)
	}

	return reportCreatedRepo(params, remoteURL)
}

func isRemoteGitURL(raw string) bool {
	prefixes := []string{"http://", "https://", "git@", "ssh://"}
	for _, p := range prefixes {
		if strings.HasPrefix(raw, p) {
			return true
		}
	}

	return false
}

// EnsureOrProvisionDestinationRepo ensures a destination exists; if not, provisions it.
func EnsureOrProvisionDestinationRepo(rawTarget string, isLocal bool) (string, error) {
	if isRemoteGitURL(rawTarget) {
		return rawTarget, nil
	}

	abs, err := filepath.Abs(rawTarget)
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve destination:")
	}

	info, statErr := os.Stat(abs)
	isFound := statErr == nil && info.IsDir()
	if isFound {
		return abs, nil
	}

	return provisionMissingDestination(rawTarget, abs, isLocal)
}

func provisionMissingDestination(rawTarget, abs string, isLocal bool) (string, error) {
	name := filepath.Base(rawTarget)
	slug := SlugifyRepoName(name)
	targetDir := abs
	if !isPathLike(rawTarget) {
		targetDir = filepath.Join(".", slug)
	}
	absTarget, _ := filepath.Abs(targetDir)

	if info, err := os.Stat(absTarget); err == nil && info.IsDir() {
		return absTarget, nil
	}

	params := createRepoParams{
		Name: name, Slug: slug, LocalDir: absTarget,
		Description:  "Provisioned by GitMap replay engine",
		IsSkipRemote: isLocal,
	}

	if err := initLocalRepo(params); err != nil {
		return "", err
	}

	tryPushRemoteProvisioned(params)
	workspacesync.SyncAll(absTarget, name)

	return absTarget, nil
}

func tryPushRemoteProvisioned(p createRepoParams) {
	if p.IsSkipRemote {
		return
	}

	remoteURL, err := pushRemoteRepo(p)
	if err != nil {
		fmt.Printf("  %sNotice: remote repository creation via gh skipped (%v); proceeding locally.%s\n",
			constants.ColorYellow, err, constants.ColorReset)

		return
	}

	fmt.Printf("  %s✓ Provisioned remote repository: %s%s\n", constants.ColorGreen, remoteURL, constants.ColorReset)
}
