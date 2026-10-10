package cmdcreate

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrepo"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func IsRemoteGitURL(raw string) bool {
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
	if IsRemoteGitURL(rawTarget) {
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

	return ProvisionMissingDestination(rawTarget, abs, isLocal)
}

func ProvisionMissingDestination(rawTarget, abs string, isLocal bool) (string, error) {
	name := filepath.Base(rawTarget)
	slug := cmdrepo.SlugifyRepoName(name)
	targetDir := abs
	if !cmdrepo.IsPathLike(rawTarget) {
		targetDir = filepath.Join(".", slug)
	}
	absTarget, _ := filepath.Abs(targetDir)

	if info, err := os.Stat(absTarget); err == nil && info.IsDir() {
		return absTarget, nil
	}

	params := cmdrepo.CreateRepoParams{
		Name: name, Slug: slug, LocalDir: absTarget,
		Description:  "Provisioned by GitMap replay engine",
		IsSkipRemote: isLocal,
	}

	if err := cmdrepo.InitLocalRepo(params); err != nil {
		return "", err
	}

	tryPushRemoteProvisioned(params)
	if !workspacesync.IsTempOrTestPath(absTarget) {
		workspacesync.SyncAll(absTarget, name)
	}

	return absTarget, nil
}

func tryPushRemoteProvisioned(p cmdrepo.CreateRepoParams) {
	if p.IsSkipRemote {
		return
	}

	remoteURL, err := cmdrepo.PushRemoteRepo(p)
	if err != nil {
		fmt.Printf("  %sNotice: remote repository creation via gh skipped (%v); proceeding locally.%s\n",
			constants.ColorYellow, err, constants.ColorReset)

		return
	}

	fmt.Printf("  %s✓ Provisioned remote repository: %s%s\n", constants.ColorGreen, remoteURL, constants.ColorReset)
}

func applyWorkspaceSync(params cmdrepo.CreateRepoParams, absDir string) {
	if params.IsNoSync {
		workspacesync.SyncAgyOnly(absDir, params.Name)
		return
	}
	if params.IsNoDesktop {
		workspacesync.SyncWithoutDesktop(absDir, params.Name)
		return
	}
	workspacesync.SyncAll(absDir, params.Name)
}
