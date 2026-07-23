package release

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/ghtoken"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/verbose"
)

// uploadToGitHub creates a GitHub release and uploads assets.
func uploadToGitHub(v Version, assets []string, opts Options) {
	token, source, err := ghtoken.Resolve()
	if err != nil {
		if len(assets) > 0 {
			fmt.Fprint(os.Stderr, constants.ErrAssetNoToken)
		}

		return
	}
	fmt.Printf(constants.MsgTokenFromSource, source)

	owner, repo, err := ParseRemoteOrigin()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAssetRemoteParse, err)

		return
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("github: creating release %s on %s/%s (%d asset(s))", v.String(), owner, repo, len(assets))
	}

	name := constants.ReleaseTagPrefix + v.String()
	if len(opts.Notes) > 0 {
		name = opts.Notes
	}

	// Release body: only the gitmap source repo gets the auto-detected
	// CHANGELOG.md notes + pinned install snippet. For any other repo,
	// the body stays empty — the user wanted "just create the tag,
	// don't dump gitmap's changelog into their release notes".
	var body string
	if ShouldPrintInstallHint(getRemoteURL()) {
		body = DetectChangelog()
		body = AppendPinnedInstallSnippet(body, v.String())
	}
	ghRelease, err := CreateGitHubRelease(owner, repo, v.String(), name, body, token, opts.IsDraft)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ GitHub release creation failed: %v\n", err)

		return
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("github: release created, id=%d", ghRelease.ID)
	}

	if len(assets) > 0 {
		fmt.Printf(constants.MsgAssetUploadStart, len(assets))
		UploadAllAssets(owner, repo, ghRelease.ID, assets, token)
	}
}

// buildGoAssetsIfApplicable cross-compiles Go binaries when --bin is passed.
func buildGoAssetsIfApplicable(v Version, opts Options) []string {
	if !opts.Bin {
		return nil
	}

	if !DetectGoProject() {
		return nil
	}

	modName, err := ReadModuleName()
	if err != nil {
		return nil
	}

	fmt.Printf(constants.MsgAssetDetected, modName)

	packages := FindMainPackages()
	if len(packages) == 0 {
		fmt.Print(constants.MsgAssetNoMain)

		return nil
	}

	targets, err := ResolveTargets(opts.Targets, opts.ConfigTargets)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Invalid targets: %v\n", err)

		return nil
	}

	stagingDir, err := EnsureStagingDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Staging dir: %v\n", err)

		return nil
	}

	fmt.Printf(constants.MsgAssetCrossCompile, len(targets)*len(packages))

	results := CrossCompile(v.String(), targets, packages, stagingDir)
	successful := CollectSuccessfulBuilds(results)

	fmt.Printf(constants.MsgAssetBuildSummary, len(successful), len(results))

	return successful
}
