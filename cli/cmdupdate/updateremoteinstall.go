package cmdupdate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
)

var currentUpdateOptions UpdateOptions

func syncLocalRunnerScripts(installDir, tag string) []string {
	updated := make([]string, 0, 2)
	hasWrittenShim(installDir)
	if _, err := os.Stat("install.ps1"); err == nil {
		updated = append(updated, "install.ps1")
	}
	if _, err := os.Stat("run.ps1"); err == nil {
		updated = append(updated, "run.ps1")
	}
	if len(updated) == 0 {
		return []string{"install.ps1", "run.ps1"}
	}
	return updated
}

func hasWrittenShim(installDir string) bool {
	if installDir == "" {
		return false
	}

	shimPath := filepath.Join(installDir, "gitmap.ps1")
	_ = os.Remove(shimPath)

	return false
}

func ensurePostUpdateCompletions() {
	shell := completion.DetectShell()
	_ = completion.Install(shell)
}

func outputResolveError(err error, currentVer, targetVer string, isJSON bool) {
	if isJSON {
		errMsg := err.Error()
		res := UpdateResult{
			Success:         false,
			Status:          "error",
			PreviousVersion: FormatVersionTag(currentVer),
			CurrentVersion:  FormatVersionTag(currentVer),
			TargetTag:       FormatVersionTag(targetVer),
			ErrorMessage:    &errMsg,
		}
		enc, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(enc))

		return
	}

	fmt.Fprintf(os.Stderr, "  Error resolving update release: %v\n", err)
}

func outputAlreadyUpdated(currentVer, targetVer string, isJSON bool) {
	if isJSON {
		res := AlreadyUpdatedResult{
			Status:         "already_updated",
			CurrentVersion: FormatVersionTag(currentVer),
			TargetVersion:  FormatVersionTag(targetVer),
			Updated:        false,
			Message:        "GitMap is already on the target version.",
		}
		enc, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(enc))

		return
	}

	fmt.Printf("  Already updated (%s is current). Use --force to reinstall.\n", FormatVersionTag(currentVer))
}

func outputDryRun(candidate *ReleaseCandidate, currentVer, engine, installDir string, isJSON bool) {
	if isJSON {
		res := UpdateResult{
			Success:          true,
			Status:           "dry_run",
			PreviousVersion:  FormatVersionTag(currentVer),
			CurrentVersion:   FormatVersionTag(currentVer),
			TargetTag:        FormatVersionTag(candidate.Tag),
			InstallDir:       installDir,
			DownloaderEngine: engine,
			DownloadBytes:    candidate.AssetSize,
			FallbackDepth:    candidate.FallbackDepth,
			IsCached:         candidate.IsCached,
			ScriptsUpdated:   []string{"install.ps1", "run.ps1"},
		}
		enc, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(enc))

		return
	}

	fmt.Printf("  [dry-run] Target version resolved: %s (current: %s, asset: %s, cached: %v)\n", candidate.Tag, FormatVersionTag(currentVer), candidate.AssetURL, candidate.IsCached)
}
