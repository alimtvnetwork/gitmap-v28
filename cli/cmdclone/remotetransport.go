package cmdclone

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// ApplyTransportFlag rewrites the `remote.origin.url` of dir to the
// requested transport and persists it via `git remote set-url`.
// Returns (changed, oldURL, newURL, err).
//
// When neither flag is set, returns immediately with changed=false.
// When both flags are set, --ssh wins and a one-line warning is
// printed to stderr (mirrors `gitmap clone` semantics).
//
// Unrecognized origin URLs fail-open: a warning is printed but no
// error is returned, so the caller can still run git push/pull.
//
// Spec: 02-spec/01-app/111-push-pull-transport-flags.md
func ApplyTransportFlag(dir string, useSSH, useHTTPS bool) (bool, string, string, error) {
	isSSH, isActionable := resolveRequestedTransport(useSSH, useHTTPS)
	if !isActionable {
		return false, "", "", nil
	}

	old, err := currentOriginURL(dir)
	if err != nil {
		return false, "", "", err
	}

	return rewriteTransportURL(dir, old, isSSH)
}

func resolveRequestedTransport(useSSH, useHTTPS bool) (bool, bool) {
	if !useSSH && !useHTTPS {
		return false, false
	}

	if useSSH && useHTTPS {
		fmt.Fprintln(os.Stderr, "⚠ both --ssh and --https set; --ssh wins")

		return true, true
	}

	return useSSH, true
}

func rewriteTransportURL(dir, old string, isSSH bool) (bool, string, string, error) {
	converted, isRecognized := convertURLForTransport(old, isSSH)
	if !isRecognized {
		fmt.Fprintf(os.Stderr, "⚠ remote.origin.url %q is not a recognized Git URL; skipping transport rewrite\n", old)

		return false, old, old, nil
	}

	hasChanged, err := applyConvertedURL(dir, old, converted)

	return hasChanged, old, converted, err
}

func convertURLForTransport(url string, isSSH bool) (string, bool) {
	if isSSH {
		return ConvertURLToSSH(url)
	}

	return ConvertURLToHTTPS(url)
}

func applyConvertedURL(dir, oldURL, convertedURL string) (bool, error) {
	if convertedURL == oldURL {
		return false, nil
	}

	if err := setOriginURL(dir, convertedURL); err != nil {
		return false, err
	}

	fmt.Printf("→ remote.origin.url: %s → %s\n", oldURL, convertedURL)

	return true, nil
}

// currentOriginURL reads the remote origin URL for dir using gitutil.
func currentOriginURL(dir string) (string, error) {
	url, err := gitutil.ReadRemoteOriginURL(dir)
	if err != nil {
		return "", fmt.Errorf("read remote.origin.url in %s: %w (is `git remote add origin <url>` set?)", dir, err)
	}

	return url, nil
}

// setOriginURL updates the remote origin URL for dir using gitutil.
func setOriginURL(dir, url string) error {
	return gitutil.SetRemoteOriginURL(dir, url)
}
