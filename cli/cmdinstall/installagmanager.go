package cmdinstall

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runInstallAgManagerWithOpts(opts installOptions) error {
	ver, isFound := isAgManagerInstalled()
	if isFound {
		fmt.Printf("  ✓ Antigravity Manager is already installed (%s)\n", ver)

		return nil
	}

	fmt.Println("Fetching release for Antigravity-Manager...")
	assetURL, relVer, err := resolveAgManagerAssetURL(opts.Version)
	if err != nil {
		reportVerificationFailure(constants.ToolAgManager, "ag-manager")

		return nil
	}

	if opts.DryRun {
		fmt.Printf("  [dry-run] Would download %s (version: %s) and execute installer\n", assetURL, relVer)

		return nil
	}

	performAgManagerDownloadAndInstall(assetURL, relVer)

	return nil
}

func performAgManagerDownloadAndInstall(assetURL, ver string) {
	fmt.Printf("Downloading %s (version: %s)...\n", assetURL, ver)
	tmpPath, err := downloadAgManagerFile(assetURL)
	if err != nil {
		reportVerificationFailure(constants.ToolAgManager, "ag-manager")

		return
	}

	fmt.Printf("Installing %s...\n", filepath.Base(tmpPath))
	if err := executeAgManagerInstaller(tmpPath); err != nil {
		reportVerificationFailure(constants.ToolAgManager, "ag-manager")

		return
	}

	recordAgManagerInstalled(ver)
	fmt.Println(constants.ColorGreen + "✓" + constants.ColorReset + " Antigravity Manager installed successfully.")
}

func recordAgManagerInstalled(ver string) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return
	}

	defer splitDB.Close()
	_ = splitDB.SaveInstalledTool("ag-manager", ver, "github-release")
}

func downloadAgManagerFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]
	tmpPath := filepath.Join(os.TempDir(), name)
	out, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}

	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	return tmpPath, err
}
