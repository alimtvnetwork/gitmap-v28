package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// repoSlugRe matches "<base>-v<N>" e.g. "gitmap-v27".
var repoSlugRe = regexp.MustCompile(`^([a-z][a-z0-9_-]*?)-v(\d+)$`)

// probeResult records a single sibling HEAD response.
type probeResult struct {
	Offset int
	Slug   string
	Status int
	IsHit  bool
}

// resolveLatestRepoSlug walks the probe → release → main fallback chain
// and returns the winning repo slug plus the source label that produced
// it. Spec: spec/01-app/111-update-remote-probe.md.
func resolveLatestRepoSlug(httpClient *http.Client) (string, string, error) {
	base, currentN, err := parseCurrentRepoSlug(constants.UpdateCurrentRepoSlug)
	if err != nil {
		return "", "", err
	}

	if slug, ok := probeSiblings(httpClient, base, currentN, constants.UpdateProbeMaxSiblings); ok {
		fmt.Printf(constants.MsgUpdateProbeResolve, slug, constants.UpdateProbeSourceSibling)
		return slug, constants.UpdateProbeSourceSibling, nil
	}

	currentSlug := constants.UpdateCurrentRepoSlug
	if releaseFallbackOK(httpClient, currentSlug) {
		fmt.Printf(constants.MsgUpdateProbeResolve, currentSlug, constants.UpdateProbeSourceRelease)
		return currentSlug, constants.UpdateProbeSourceRelease, nil
	}

	if mainFallbackOK(httpClient, currentSlug) {
		fmt.Printf(constants.MsgUpdateProbeResolve, currentSlug, constants.UpdateProbeSourceMain)
		return currentSlug, constants.UpdateProbeSourceMain, nil
	}

	fmt.Fprint(os.Stderr, constants.ErrUpdateProbeNoResolve)
	return "", "", fmt.Errorf("no resolution")
}

// parseCurrentRepoSlug splits "gitmap-v27" into ("gitmap", 23).
func parseCurrentRepoSlug(slug string) (string, int, error) {
	m := repoSlugRe.FindStringSubmatch(slug)
	if m == nil {
		return "", 0, fmt.Errorf(constants.ErrUpdateProbeParseSlug, slug, fmt.Errorf("no match"))
	}
	n, err := strconv.Atoi(m[2])
	if err != nil {
		return "", 0, fmt.Errorf(constants.ErrUpdateProbeParseSlug, slug, err)
	}
	return m[1], n, nil
}

// probeSiblings fires up to maxN parallel HEAD requests against
// "<base>-v<current+1..current+maxN>" and returns the highest-N hit.
func probeSiblings(httpClient *http.Client, base string, current, maxN int) (string, bool) {
	fmt.Printf(constants.MsgUpdateProbeStart, base, current+1, base, current+maxN)

	results := make([]probeResult, maxN)
	var wg sync.WaitGroup
	for i := 1; i <= maxN; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			results[offset-1] = probeOne(httpClient, base, current+offset, offset)
		}(i)
	}
	wg.Wait()

	return pickMaxHit(results, base)
}

// probeOne runs a single HEAD request and returns the result.
func probeOne(httpClient *http.Client, base string, n, offset int) probeResult {
	slug := fmt.Sprintf("%s-v%d", base, n)
	url := fmt.Sprintf(constants.UpdateRepoHEADTmpl, constants.UpdateRepoOwner, slug)
	status := headStatus(httpClient, url)
	hit := status >= 200 && status < 300
	if hit {
		fmt.Printf(constants.MsgUpdateProbeHit, slug, status)
	}
	return probeResult{Offset: offset, Slug: slug, Status: status, IsHit: hit}
}

// pickMaxHit returns the highest-offset hit from results.
func pickMaxHit(results []probeResult, _ string) (string, bool) {
	for i := len(results) - 1; i >= 0; i-- {
		if results[i].IsHit {
			return results[i].Slug, true
		}
	}
	return "", false
}

// headStatus performs a HEAD request with a per-request timeout and
// returns the status code (0 on transport error).
func headStatus(httpClient *http.Client, url string) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.UpdateProbeTimeoutSec)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return 0
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

// releaseFallbackOK pings the GitHub releases API for the current repo.
func releaseFallbackOK(httpClient *http.Client, slug string) bool {
	url := fmt.Sprintf(constants.UpdateReleasesAPITmpl, constants.UpdateRepoOwner, slug)
	status := headStatus(httpClient, url)
	return status >= 200 && status < 300
}

// mainFallbackOK pings the repo root (which serves main HEAD).
func mainFallbackOK(httpClient *http.Client, slug string) bool {
	url := fmt.Sprintf(constants.UpdateRepoHEADTmpl, constants.UpdateRepoOwner, slug)
	status := headStatus(httpClient, url)
	return status >= 200 && status < 300
}

// newProbeClient returns the default HTTP client used by the probe.
func newProbeClient() *http.Client {
	return &http.Client{Timeout: time.Duration(constants.UpdateProbeTimeoutSec) * time.Second}
}
