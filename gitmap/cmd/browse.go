// Package cmd — browse.go: open URLs in default browser or Google Chrome.
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func runBrowseCmd(args []string) error {
	rawURL, isChrome := parseBrowseArgs(args)
	if rawURL == "" {
		printBrowseUsage()

		return apperror.NewValidationError("url required for browse command")
	}

	targetURL := normalizeBrowseURL(rawURL)
	if isChrome {
		return openURLInChrome(targetURL)
	}

	return openURLInDefaultBrowser(targetURL)
}

func parseBrowseArgs(args []string) (string, bool) {
	url := ""
	isChrome := false
	for _, a := range args {
		checkBrowseFlag(a, &url, &isChrome)
	}

	return url, isChrome
}

func checkBrowseFlag(a string, url *string, isChrome *bool) {
	if a == "--chrome" || a == "-c" {
		*isChrome = true

		return
	}

	if !strings.HasPrefix(a, "-") && *url == "" {
		*url = a
	}
}

func printBrowseUsage() {
	fmt.Println("Usage: gitmap open-url <url> [--chrome]")
	fmt.Println("       gitmap browse <url> [--chrome]")
	fmt.Println()
	fmt.Println("Opens web URLs in default browser or Google Chrome across platforms.")
	fmt.Println()
}

func normalizeBrowseURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	lower := strings.ToLower(trimmed)
	if hasURLProtocol(lower) {
		return trimmed
	}

	return "https://" + trimmed
}

func hasURLProtocol(lower string) bool {
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "chrome://") ||
		strings.HasPrefix(lower, "file://")
}

func openURLInChrome(url string) error {
	bin, err := findChromeBinaryPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(bin, url)
	configureDetachedProcess(cmd)
	if err := cmd.Start(); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("launch chrome with %q", url))
	}

	fmt.Printf("  %s🌐 Opened in Google Chrome: %s%s%s\n",
		constants.ColorCyan, constants.ColorWhite, url, constants.ColorReset)

	return nil
}

func openURLInDefaultBrowser(url string) error {
	if err := launchNativeOpener(url); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("launch default browser with %q", url))
	}

	fmt.Printf("  %s🌐 Opened in default browser: %s%s%s\n",
		constants.ColorCyan, constants.ColorWhite, url, constants.ColorReset)

	return nil
}
