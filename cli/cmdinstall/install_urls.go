// Package cmdinstall — install_urls.go: `gitmap install urls <tag>`.
//
// Extracts every download URL from the release's installer scripts
// (install.ps1, install.sh, Homebrew cask when present in the repo) at the
// given tag and HEAD-checks each one. Read-only.
//
// Script source: raw.githubusercontent.com at the tag, falling back to the
// local checkout when the fetch fails (reported per script).
package cmdinstall

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

// isInstallURLsCommand reports whether args invoke `install urls`.
func isInstallURLsCommand(args []string) bool {
	return len(args) > 0 && strutil.EqualFoldAnyTrim(args[0], "urls")
}

// installURLPattern matches http(s) URLs in installer script bodies.
var installURLPattern = lazyregex.New(`https?://[^\s"'()<>\\]+`)

// installScriptSource is one installer script body plus where it came from.
type installScriptSource struct {
	name   string
	origin string
	body   string
}

// installURLCheck is the HEAD result for one extracted URL.
type installURLCheck struct {
	url        string
	source     string
	statusCode int
	isOK       bool
	note       string
}

// RunInstallURLs executes `gitmap install urls <tag>` — read-only.
// args is the tail after the `urls` token; args[0] is the tag.
func RunInstallURLs(args []string) *appfault.AppError {
	if isInstallHelpRequested(args) {
		fmt.Print(installURLsHelpText)
		cliexit.Exit(0)
	}

	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return appfault.NewValidationError("missing required argument: <tag>\n" + installURLsUsageLine +
			"example: gitmap install urls v6.558.0")
	}
	tag := strings.TrimSpace(args[0])

	sources, appErr := fetchInstallScripts(tag)
	if appErr != nil {
		return appErr
	}

	urls := extractInstallURLs(sources)
	if len(urls) == 0 {
		fmt.Println("no download URLs found in installer scripts")

		return nil
	}

	start := time.Now()
	checks := headCheckInstallURLs(urls)
	elapsed := time.Since(start).Round(time.Millisecond)

	printInstallURLSources(sources)
	printInstallURLChecks(checks, elapsed)

	return nil
}

// fetchInstallScripts loads install.sh / install.ps1 at the tag from GitHub,
// falling back to the local checkout per script. Also picks up a Homebrew
// cask (*.rb) when one is present in the repo.
func fetchInstallScripts(tag string) ([]installScriptSource, *appfault.AppError) {
	sources := make([]installScriptSource, 0, 3)
	for _, name := range []string{"install.sh", "install.ps1"} {
		source, appErr := fetchOneInstallScript(name, tag)
		if appErr != nil {
			return nil, appErr
		}
		sources = append(sources, source)
	}

	if cask, isFound := findHomebrewCask(); isFound {
		sources = append(sources, cask)
	}

	return sources, nil
}

// fetchOneInstallScript loads one script at the tag, else from the checkout.
func fetchOneInstallScript(name, tag string) (installScriptSource, *appfault.AppError) {
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
		constants.UpdateRepoOwner, constants.UpdateCurrentRepoSlug, tag, name)
	if body, appErr := httpGetBody(rawURL); appErr == nil {
		return installScriptSource{name: name, origin: "github " + tag, body: body}, nil
	}

	local := filepath.Join(constants.RepoPath, name)
	data, err := os.ReadFile(local)
	if err != nil {
		return installScriptSource{}, appfault.NewNotFoundError(
			fmt.Sprintf("installer script %s not found at tag %s and not in local checkout", name, tag))
	}

	return installScriptSource{name: name, origin: "local checkout", body: string(data)}, nil
}

// findHomebrewCask locates a Homebrew cask file in the repo, if present.
func findHomebrewCask() (installScriptSource, bool) {
	for _, dir := range []string{"Casks", "casks", "homebrew", ".github/homebrew"} {
		entries, err := os.ReadDir(filepath.Join(constants.RepoPath, dir))
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".rb") {
				path := filepath.Join(constants.RepoPath, dir, entry.Name())
				data, err := os.ReadFile(path)
				if err != nil {
					continue
				}

				return installScriptSource{name: entry.Name(), origin: "local checkout", body: string(data)}, true
			}
		}
	}

	return installScriptSource{}, false
}

// extractInstallURLs dedupes http(s) URLs across scripts, keeping the first
// script each URL was seen in.
func extractInstallURLs(sources []installScriptSource) []installURLCheck {
	seen := map[string]bool{}
	checks := []installURLCheck{}
	for _, source := range sources {
		for _, raw := range installURLPattern.FindAllString(source.body, -1) {
			url := strings.TrimRight(raw, ".,;:")
			if seen[url] {
				continue
			}
			seen[url] = true
			checks = append(checks, installURLCheck{url: url, source: source.name})
		}
	}
	sort.Slice(checks, func(i, j int) bool { return checks[i].url < checks[j].url })

	return checks
}

// headCheckInstallURLs HEAD-checks each URL (GET fallback on 405).
func headCheckInstallURLs(checks []installURLCheck) []installURLCheck {
	client := &http.Client{Timeout: 10 * time.Second}
	for i := range checks {
		code, note, isOK := probeInstallURL(client, checks[i].url)
		checks[i].statusCode = code
		checks[i].note = note
		checks[i].isOK = isOK
	}

	return checks
}

// probeInstallURL issues HEAD, falling back to GET when the server rejects HEAD.
func probeInstallURL(client *http.Client, url string) (int, string, bool) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return 0, err.Error(), false
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err.Error(), false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusMethodNotAllowed {
		return getInstallURL(client, url)
	}

	return resp.StatusCode, "", resp.StatusCode == http.StatusOK
}

// getInstallURL retries with GET for servers that reject HEAD.
func getInstallURL(client *http.Client, url string) (int, string, bool) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, err.Error(), false
	}
	defer resp.Body.Close()

	return resp.StatusCode, "", resp.StatusCode == http.StatusOK
}

// httpGetBody fetches a URL body, requiring HTTP 200.
func httpGetBody(url string) (string, *appfault.AppError) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", appfault.NewValidationError("fetch failed: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", appfault.NewValidationError(fmt.Sprintf("fetch %s: HTTP %d", url, resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", appfault.NewValidationError("read failed: " + err.Error())
	}

	return string(body), nil
}

// printInstallURLSources lists which script bodies were audited.
func printInstallURLSources(sources []installScriptSource) {
	for _, source := range sources {
		fmt.Printf("%s (%s)\n", source.name, source.origin)
	}
}

// printInstallURLChecks renders the per-URL verdicts and the summary line.
func printInstallURLChecks(checks []installURLCheck, elapsed time.Duration) {
	failed := 0
	for _, check := range checks {
		if check.isOK {
			continue
		}
		failed++
		detail := check.note
		if check.statusCode != 0 {
			detail = fmt.Sprintf("HTTP %d", check.statusCode)
		}
		fmt.Printf("✗ %s → %s (%s)\n", check.url, detail, check.source)
	}

	okCount := len(checks) - failed
	fmt.Printf("%d urls checked · %d × 200 ✓", len(checks), okCount)
	if failed > 0 {
		fmt.Printf(" · %d failed ✗", failed)
	}
	fmt.Printf(" (%s)\n", elapsed)
}
