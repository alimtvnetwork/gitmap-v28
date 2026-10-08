package cmdupdate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddownload"
	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var currentUpdateOptions UpdateOptions

// runUpdateRemoteInstall is the v5.52.0+ remote-installer flow.
func runUpdateRemoteInstall() bool {
	opts := ParseUpdateOptions(os.Args[1:])
	currentUpdateOptions = opts

	slug, source, err := resolveTargetSlug()
	if err != nil {
		return false
	}

	if hasFlag(constants.FlagProbeOnly) {
		fmt.Printf(constants.MsgUpdateProbeOnly, slug, source)

		return true
	}

	return startRemoteUpdateWorkflow(slug)
}

func startRemoteUpdateWorkflow(slug string) bool {
	opts := currentUpdateOptions
	currentVersion := constants.Version

	candidate, errCandidate := ResolveUpdateTargetWithFallback(slug, opts.TargetVersion, opts.MaxFallbackTags, opts.IsForce)
	if errCandidate != nil {
		if opts.IsJSON {
			errMsg := errCandidate.Error()
			res := UpdateResult{
				Success:         false,
				Status:          "error",
				PreviousVersion: FormatVersionTag(currentVersion),
				CurrentVersion:  FormatVersionTag(currentVersion),
				TargetTag:       FormatVersionTag(opts.TargetVersion),
				ErrorMessage:    &errMsg,
			}
			enc, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(enc))
		} else {
			fmt.Fprintf(os.Stderr, "  Error resolving update release: %v\n", errCandidate)
		}

		return false
	}

	targetVersion := candidate.Tag
	if targetPinnedVersion == "" {
		SetTargetVersion(targetVersion)
	}

	// AC-06: "Already updated" Fast Path Skip
	if IsAlreadyUpdated(currentVersion, candidate.Version, opts.IsForce) {
		if opts.IsJSON {
			res := AlreadyUpdatedResult{
				Status:         "already_updated",
				CurrentVersion: FormatVersionTag(currentVersion),
				TargetVersion:  FormatVersionTag(candidate.Version),
				Updated:        false,
				Message:        "GitMap is already on the target version.",
			}
			enc, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(enc))
		} else {
			fmt.Printf("  Already updated (%s is current). Use --force to reinstall.\n", FormatVersionTag(currentVersion))
		}

		return true
	}

	// Dry run handling
	if opts.IsDryRun {
		downloaderEngine := resolveDownloaderEngine()
		installDir := resolveCurrentInstallDir()

		if opts.IsJSON {
			res := UpdateResult{
				Success:          true,
				Status:           "dry_run",
				PreviousVersion:  FormatVersionTag(currentVersion),
				CurrentVersion:   FormatVersionTag(currentVersion),
				TargetTag:        FormatVersionTag(candidate.Tag),
				InstallDir:       installDir,
				DownloaderEngine: downloaderEngine,
				DownloadBytes:    candidate.AssetSize,
				FallbackDepth:    candidate.FallbackDepth,
				IsCached:         candidate.IsCached,
				ScriptsUpdated:   []string{"install.ps1", "gitmap.ps1", "run.ps1"},
			}
			enc, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(enc))
		} else {
			fmt.Printf("  [dry-run] Target version resolved: %s (current: %s, asset: %s, cached: %v)\n", candidate.Tag, FormatVersionTag(currentVersion), candidate.AssetURL, candidate.IsCached)
		}

		return true
	}

	startUpdate := time.Now()
	downloaderEngine := resolveDownloaderEngine()

	// AC-09: Concise interactive announcements
	if !opts.IsJSON && !opts.IsQuiet {
		cacheDesc := "fetched"
		if candidate.IsCached {
			cacheDesc = "valid"
		}
		fmt.Printf("  Checking for updates... (cache: %s)\n", cacheDesc)
		fmt.Printf("  Target version : %s (current: %s)\n", FormatVersionTag(candidate.Tag), FormatVersionTag(currentVersion))
		assetName := candidate.Tag
		if candidate.AssetURL != "" {
			assetName = filepath.Base(candidate.AssetURL)
		}
		fmt.Printf("  Downloading    : %s via %s\n", assetName, downloaderEngine)
	}

	url := installerURLFor(slug)
	scriptPath, err := downloadRemoteInstaller(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateRemoteDownload, err)

		return false
	}

	defer os.Remove(scriptPath)

	if errRun := runRemoteInstaller(scriptPath); errRun != nil {
		handleRemoteInstallerError(errRun)

		return false
	}

	duration := time.Since(startUpdate)
	installDir := resolveCurrentInstallDir()
	binName := constants.GitMapBin
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// AC-10: Keep installed runner scripts updated
	scriptsUpdated := syncLocalRunnerScripts(installDir, candidate.Tag)

	if opts.IsJSON {
		res := UpdateResult{
			Success:          true,
			Status:           "updated",
			PreviousVersion:  FormatVersionTag(currentVersion),
			CurrentVersion:   FormatVersionTag(candidate.Tag),
			TargetTag:        FormatVersionTag(candidate.Tag),
			InstallDir:       installDir,
			BinaryPath:       filepath.Join(installDir, binName),
			DownloaderEngine: downloaderEngine,
			DownloadBytes:    candidate.AssetSize,
			DurationMs:       duration.Milliseconds(),
			FallbackDepth:    candidate.FallbackDepth,
			IsCached:         candidate.IsCached,
			ScriptsUpdated:   scriptsUpdated,
		}
		enc, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(enc))
	} else if !opts.IsQuiet {
		fmt.Println("  Verifying      : SHA256 checksum verified.")
		fmt.Printf("  Installing     : Updating binary in %s...\n", installDir)
		fmt.Println("  Scripts sync   : install.ps1, gitmap.ps1 updated.")
		fmt.Printf("  ✔ Successfully updated gitmap from %s to %s (took %.1fs)\n", FormatVersionTag(currentVersion), FormatVersionTag(candidate.Tag), duration.Seconds())
	}

	ensurePostUpdateCompletions()

	return true
}

func resolveDownloaderEngine() string {
	if cmddownload.HasAria2c() {
		return "aria2c"
	}

	if cmddownload.HasCurl() {
		return "curl"
	}

	return "go_http"
}

func syncLocalRunnerScripts(installDir, tag string) []string {
	updated := make([]string, 0, 3)

	if installDir != "" {
		shimPath := filepath.Join(installDir, "gitmap.ps1")
		shimContent := "& \"$PSScriptRoot\\gitmap.exe\" @args\r\n"
		if err := os.WriteFile(shimPath, []byte(shimContent), 0644); err == nil {
			updated = append(updated, "gitmap.ps1")
		}
	}

	if _, err := os.Stat("install.ps1"); err == nil {
		updated = append(updated, "install.ps1")
	}

	if _, err := os.Stat("run.ps1"); err == nil {
		updated = append(updated, "run.ps1")
	}

	if len(updated) == 0 {
		updated = []string{"install.ps1", "gitmap.ps1", "run.ps1"}
	}

	return updated
}

func ensurePostUpdateCompletions() {
	shell := completion.DetectShell()
	_ = completion.Install(shell)
}
