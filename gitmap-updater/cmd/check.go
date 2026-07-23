package cmd

import (
	"fmt"
	"os"
)

// RunCheck checks GitHub for the latest release and compares with installed version.
func RunCheck() {
	fmt.Print(MsgChecking)

	current, err := getInstalledVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrGetVersion, err)
		os.Exit(1)
	}

	fmt.Printf(MsgCurrentVer, current)

	latest, err := fetchLatestTag()
	if err != nil {
		fmt.Fprintf(os.Stderr, ErrFetchRelease, err)
		os.Exit(1)
	}

	fmt.Printf(MsgLatestVer, latest)

	if normalizeVersion(current) == normalizeVersion(latest) {
		fmt.Printf(MsgUpToDate, current)
		return
	}

	fmt.Printf(MsgUpdateAvail, current, latest)
	fmt.Println("\n  → Run 'gitmap-updater run' to update.")
	fmt.Println()
}
