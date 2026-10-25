package cmdupdate

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// probeThrottleDuration specifies that sibling probing is run at most once per 24 hours.
const probeThrottleDuration = 24 * time.Hour

// probeWatermarkState stores the timestamp and slug from previous sibling probing.
type probeWatermarkState struct {
	LastProbeTime time.Time `json:"last_probe_time"`
	CachedSlug    string    `json:"cached_slug"`
	HasHit        bool      `json:"has_hit"`
}

var (
	probeWatermarkFileOverride string
	probeThrottleForceEnabled  bool
)

func probeWatermarkFile() string {
	if probeWatermarkFileOverride != "" {
		return probeWatermarkFileOverride
	}
	return filepath.Join(store.BinaryDataDir(), "probe-watermark.json")
}

func loadProbeWatermark() *probeWatermarkState {
	path := probeWatermarkFile()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var state probeWatermarkState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil
	}
	return &state
}

func saveProbeWatermark(cachedSlug string, hasHit bool) {
	state := probeWatermarkState{
		LastProbeTime: time.Now().UTC(),
		CachedSlug:    cachedSlug,
		HasHit:        hasHit,
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return
	}
	path := probeWatermarkFile()
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	_ = os.WriteFile(path, data, 0644)
}

func isTestEnv() bool {
	if flag.Lookup("test.v") != nil {
		return true
	}
	arg0 := strings.ToLower(os.Args[0])
	return strings.HasSuffix(arg0, ".test") || strings.HasSuffix(arg0, ".test.exe")
}

func checkProbeThrottle(base string, current int) (string, bool) {
	if !probeThrottleForceEnabled && isTestEnv() {
		return "", false
	}

	wm := loadProbeWatermark()
	if wm == nil {
		return "", false
	}

	if time.Since(wm.LastProbeTime) >= probeThrottleDuration {
		return "", false
	}

	if wm.HasHit && wm.CachedSlug != "" {
		return wm.CachedSlug, true
	}

	currentSlug := fmt.Sprintf("%s-v%d", base, current)

	return currentSlug, true
}

// repoSlugRe matches "<base>-v<N>" e.g. "gitmap-v28".
var repoSlugRe = lazyregex.New(`^([a-z][a-z0-9_-]*?)-v(\d+)$`)

// probeResult records a single sibling HEAD response.
type probeResult struct {
	Offset int
	Slug   string
	Status int
	IsHit  bool
}

// resolveLatestRepoSlug walks the probe → release → main fallback chain
// and returns the winning repo slug plus the source label that produced
// it. Spec: 02-spec/01-app/111-update-remote-probe.md.
func resolveLatestRepoSlug(httpClient *http.Client) (string, string, error) {
	base, currentN, err := parseCurrentRepoSlug(constants.UpdateCurrentRepoSlug)
	if err != nil {
		return "", "", err
	}

	if slug, hasSibling := checkSiblingCandidate(httpClient, base, currentN); hasSibling {
		fmt.Printf(constants.MsgUpdateProbeResolve, slug, constants.UpdateProbeSourceSibling)

		return slug, constants.UpdateProbeSourceSibling, nil
	}

	return resolveFallbackSlug(httpClient)
}

func checkSiblingCandidate(httpClient *http.Client, base string, currentN int) (string, bool) {
	slug, hasSibling := probeSiblings(httpClient, base, currentN, constants.UpdateProbeMaxSiblings)
	if !hasSibling {
		return "", false
	}

	if !hasNewerSiblingRelease(httpClient, slug, constants.Version) {
		return "", false
	}

	return slug, true
}

func hasNewerSiblingRelease(httpClient *http.Client, slug, currentVersion string) bool {
	tag, err := fetchRepoLatestReleaseTag(httpClient, slug)
	if err != nil || tag == "" {
		return false
	}

	return CompareSemVer(tag, currentVersion) > 0
}

func fetchRepoLatestReleaseTag(httpClient *http.Client, slug string) (string, error) {
	resp, err := getRepoReleaseResponse(httpClient, slug)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tag := decodeReleaseTagName(resp.Body)
	if tag == "" {
		return "", apperror.NewValidationError("empty release tag")
	}

	return tag, nil
}

func newReleaseGetRequest(ctx context.Context, slug string) (*http.Request, error) {
	url := fmt.Sprintf(constants.UpdateReleasesAPITmpl, constants.UpdateRepoOwner, slug)

	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

func getRepoReleaseResponse(httpClient *http.Client, slug string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(constants.UpdateProbeTimeoutSec)*time.Second)
	defer cancel()

	req, errReq := newReleaseGetRequest(ctx, slug)
	if errReq != nil {
		return nil, errReq
	}

	resp, errResp := httpClient.Do(req)
	if errResp != nil {
		return nil, errResp
	}

	return validateReleaseResponse(resp)
}

func validateReleaseResponse(resp *http.Response) (*http.Response, error) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()

		return nil, apperror.NewExecutionError(fmt.Sprintf("status %d", resp.StatusCode))
	}

	return resp, nil
}

func resolveFallbackSlug(httpClient *http.Client) (string, string, error) {
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

	return "", "", apperror.NewNotFoundError("no resolution")
}

// parseCurrentRepoSlug splits "gitmap-v28" into ("gitmap", 23).
func parseCurrentRepoSlug(slug string) (string, int, error) {
	m := repoSlugRe.FindStringSubmatch(slug)
	if m == nil {
		return "", 0, apperror.WrapSimple(apperror.NewExecutionError("no match"), fmt.Sprintf("could not parse current repo slug %q", slug))
	}

	n, err := strconv.Atoi(m[2])
	if err != nil {
		return "", 0, apperror.WrapSimple(err, fmt.Sprintf("could not parse current repo slug %q", slug))
	}

	return m[1], n, nil
}

// probeSiblings fires up to maxN parallel HEAD requests against
// "<base>-v<current+1..current+maxN>" and returns the highest-N hit.
// Probing is throttled to at most once per 24 hours. If probed in the
// last 24h, probing is skipped and the cached or current slug is returned.
func probeSiblings(httpClient *http.Client, base string, current, maxN int) (string, bool) {
	if cachedSlug, throttled := checkProbeThrottle(base, current); throttled {
		return cachedSlug, true
	}

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

	slug, isHit := pickMaxHit(results, base)
	recordProbeWatermark(slug, base, current, isHit)

	return slug, isHit
}

func recordProbeWatermark(slug, base string, current int, hasHit bool) {
	if isTestEnv() && !probeThrottleForceEnabled {
		return
	}

	if hasHit {
		saveProbeWatermark(slug, true)

		return
	}

	saveProbeWatermark(fmt.Sprintf("%s-v%d", base, current), false)
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
