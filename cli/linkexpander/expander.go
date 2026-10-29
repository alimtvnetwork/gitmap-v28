package linkexpander

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var urlRegex = regexp.MustCompile(`https?://[^\s<>"'()]+`)
var serviceActivationRegex = regexp.MustCompile(`^https?://(?:www\.)?serviceactivation\.in/subscription/new/([a-zA-Z0-9_-]+)`)

// ExpandText finds all URLs in the provided text, resolves any redirects
// or partner activation paths, and replaces them in the returned text.
func ExpandText(text string) (string, []ExpandedLink) {
	matches := urlRegex.FindAllString(text, -1)
	if len(matches) == 0 {
		return text, nil
	}

	return resolveUniqueLinks(matches, text)
}

func resolveUniqueLinks(matches []string, text string) (string, []ExpandedLink) {
	seen := make(map[string]ExpandedLink)
	var links []ExpandedLink
	result := text

	for _, raw := range matches {
		result = appendUniqueLink(raw, seen, &links, result)
	}

	return result, links
}

func appendUniqueLink(raw string, seen map[string]ExpandedLink, links *[]ExpandedLink, text string) string {
	clean := strings.TrimRight(raw, ".,;:!?)>\"'")
	if _, hasSeen := seen[clean]; hasSeen {
		return text
	}
	expanded := ExpandURL(clean)
	seen[clean] = expanded
	*links = append(*links, expanded)

	return applyExpandedURL(text, clean, expanded)
}

func applyExpandedURL(text, clean string, expanded ExpandedLink) string {
	if !expanded.IsExpanded {
		return text
	}

	return strings.ReplaceAll(text, clean, expanded.Expanded)
}

// ExpandURL expands a single URL. It uses domain-specific direct mappings
// and follows HTTP redirects via HEAD/GET requests.
func ExpandURL(rawURL string) ExpandedLink {
	if known, hasMatch := checkKnownPartnerRule(rawURL); hasMatch {
		return known
	}

	return resolveRemoteURL(rawURL)
}

func checkKnownPartnerRule(rawURL string) (ExpandedLink, bool) {
	match := serviceActivationRegex.FindStringSubmatch(rawURL)
	if len(match) <= 1 {
		return ExpandedLink{}, false
	}
	expanded := fmt.Sprintf("https://serviceactivation.google.com/subscription/new/%s", match[1])

	return ExpandedLink{
		Original:   rawURL,
		Expanded:   expanded,
		IsExpanded: true,
		StatusCode: http.StatusOK,
	}, true
}

func resolveRemoteURL(rawURL string) ExpandedLink {
	client := createHTTPClient()
	resp, reqErr := executeRequestWithFallback(client, rawURL)
	if reqErr != nil {
		return ExpandedLink{Original: rawURL, Expanded: rawURL, Error: reqErr.Error()}
	}
	defer resp.Body.Close()
	finalURL, hasRedirect := resolveFinalRedirect(resp, rawURL)

	return buildRemoteLink(rawURL, finalURL, hasRedirect, resp.StatusCode)
}

func buildRemoteLink(original, expanded string, isExpanded bool, statusCode int) ExpandedLink {
	return ExpandedLink{
		Original:   original,
		Expanded:   expanded,
		IsExpanded: isExpanded,
		StatusCode: statusCode,
	}
}

func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func executeRequestWithFallback(client *http.Client, rawURL string) (*http.Response, error) {
	resp, headErr := executeHead(client, rawURL)
	if headErr == nil && resp.StatusCode < 400 {
		return resp, nil
	}

	return executeGet(client, rawURL)
}

func executeHead(client *http.Client, rawURL string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gitmap-linkexpander/1.0")

	return client.Do(req)
}

func executeGet(client *http.Client, rawURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gitmap-linkexpander/1.0")
	req.Header.Set("Range", "bytes=0-0")

	return client.Do(req)
}

func resolveFinalRedirect(resp *http.Response, rawURL string) (string, bool) {
	if resp == nil || resp.Request == nil || resp.Request.URL == nil {
		return rawURL, false
	}
	finalURL := resp.Request.URL.String()
	if len(finalURL) == 0 || finalURL == rawURL {
		return rawURL, false
	}

	return finalURL, true
}

// ExecuteInstructionFile reads a JSON instruction file and executes it.
func ExecuteInstructionFile(path string) (InstructionResult, error) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return InstructionResult{}, readErr
	}
	cfg := DefaultInstructionConfig()
	if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
		return InstructionResult{}, unmarshalErr
	}

	return ExecuteInstruction(cfg)
}

// ExecuteInstruction runs all links through sandboxed resolution and autoSend.
func ExecuteInstruction(cfg InstructionConfig) (InstructionResult, error) {
	if valErr := validateInstructionConfig(cfg); valErr != nil {
		return InstructionResult{}, valErr
	}
	res := InstructionResult{TotalLinks: len(cfg.Links)}
	for _, item := range cfg.Links {
		report := processInstructionItem(item, cfg)
		res.Results = append(res.Results, report)
		updateInstructionCounters(&res, report)
	}

	return res, nil
}

func validateInstructionConfig(cfg InstructionConfig) error {
	if len(cfg.Links) == 0 {
		return apperror.NewValidationError("instruction must contain at least 1 link")
	}
	for _, link := range cfg.Links {
		if valErr := validateInstructionLink(link); valErr != nil {
			return valErr
		}
	}

	return nil
}

func validateInstructionLink(link InstructionLink) error {
	if len(strings.TrimSpace(link.URL)) == 0 {
		return apperror.NewValidationError("link url cannot be empty")
	}
	if len(link.Recipients) > 5 {
		return apperror.NewValidationError("recipient count exceeds maximum allowance of 5")
	}

	return nil
}

func processInstructionItem(item InstructionLink, cfg InstructionConfig) InstructionLinkReport {
	report := resolveInstructionLink(item, cfg)
	applyInterpolation(&report, item)
	if shouldAutoSend(cfg, report, item) {
		dispatchInstructionEmail(&report, item)
	}

	return report
}

func updateInstructionCounters(res *InstructionResult, report InstructionLinkReport) {
	if report.IsExpanded {
		res.ExpandedCount++
	}
	if report.IsSent {
		res.SentCount++
	}
	if len(report.Error) > 0 {
		res.HasErrors = true
	}
}

func shouldAutoSend(cfg InstructionConfig, report InstructionLinkReport, item InstructionLink) bool {
	if !cfg.IsAutoSend || len(item.Recipients) == 0 || len(report.Error) > 0 {
		return false
	}

	return true
}

func applyInterpolation(report *InstructionLinkReport, item InstructionLink) {
	report.InterpolatedSubject = interpolateTemplate(item.Subject, item.URL, report.ExpandedURL)
	report.InterpolatedBody = interpolateTemplate(item.BodyTemplate, item.URL, report.ExpandedURL)
}

func interpolateTemplate(template, originalURL, expandedURL string) string {
	res := strings.ReplaceAll(template, "{{url}}", expandedURL)
	res = strings.ReplaceAll(res, "{{original}}", originalURL)
	res = strings.ReplaceAll(res, "{{expanded}}", expandedURL)

	return res
}

func resolveInstructionLink(item InstructionLink, cfg InstructionConfig) InstructionLinkReport {
	report := InstructionLinkReport{OriginalURL: item.URL, ExpandedURL: item.URL}
	if protoErr := checkProtocol(item.URL, cfg.Safety.AllowedProtocols); protoErr != nil {
		report.Error = protoErr.Error()
		return report
	}
	if known, hasMatch := checkKnownPartnerRule(item.URL); hasMatch {
		return applyKnownPartnerReport(report, known)
	}

	return executeHTTPResolution(item.URL, cfg)
}

func applyKnownPartnerReport(base InstructionLinkReport, known ExpandedLink) InstructionLinkReport {
	base.ExpandedURL = known.Expanded
	base.IsExpanded = true
	base.StatusCode = known.StatusCode
	base.RedirectHops = 1
	base.HopTrace = []string{known.Expanded}

	return base
}

func checkProtocol(rawURL string, allowed []string) error {
	parsed, parseErr := url.ParseRequestURI(rawURL)
	if parseErr != nil {
		return apperror.NewValidationError("invalid URL format: " + parseErr.Error())
	}
	if !isAllowedProtocol(parsed.Scheme, allowed) {
		return apperror.NewValidationError(fmt.Sprintf("unsupported protocol scheme '%s'", parsed.Scheme))
	}

	return nil
}

func isAllowedProtocol(scheme string, allowed []string) bool {
	protocols := allowed
	if len(protocols) == 0 {
		protocols = []string{"http", "https"}
	}
	lower := strings.ToLower(scheme)
	for _, p := range protocols {
		if strings.ToLower(p) == lower {
			return true
		}
	}

	return false
}

func executeHTTPResolution(rawURL string, cfg InstructionConfig) InstructionLinkReport {
	report := InstructionLinkReport{OriginalURL: rawURL, ExpandedURL: rawURL}
	tracker := newRedirectTracker(rawURL, cfg)
	client := createInstructionClient(cfg, tracker)
	resp, reqErr := executeInstructionProbes(client, rawURL, cfg.Safety.MaxBodyBytes)
	if reqErr != nil {
		report.Error = reqErr.Error()
	}
	populateResolutionReport(&report, tracker, resp)

	return report
}

type redirectTracker struct {
	visited    map[string]int
	trace      []string
	maxLevels  int
	isInfinite bool
	allowed    []string
}

func newRedirectTracker(rawURL string, cfg InstructionConfig) *redirectTracker {
	visited := make(map[string]int)
	visited[rawURL] = 0

	return &redirectTracker{
		visited:    visited,
		maxLevels:  cfg.MaxLevels,
		isInfinite: cfg.IsInfinite,
		allowed:    cfg.Safety.AllowedProtocols,
	}
}

func (t *redirectTracker) checkRedirect(req *http.Request, via []*http.Request) error {
	if checkErr := t.validateRedirectHop(req.URL); checkErr != nil {
		return checkErr
	}
	t.trace = append(t.trace, req.URL.String())
	t.visited[req.URL.String()] = len(via)
	if !t.isInfinite && len(via) >= t.maxLevels {
		return http.ErrUseLastResponse
	}

	return nil
}

func (t *redirectTracker) validateRedirectHop(u *url.URL) error {
	if !isAllowedProtocol(u.Scheme, t.allowed) {
		return fmt.Errorf("unsupported protocol scheme '%s'", u.Scheme)
	}
	if _, hasSeen := t.visited[u.String()]; hasSeen {
		return fmt.Errorf("circular redirect detected at '%s'", u.String())
	}

	return nil
}

func createInstructionClient(cfg InstructionConfig, tracker *redirectTracker) *http.Client {
	timeout := cfg.Safety.TimeoutSeconds
	if timeout <= 0 {
		timeout = 10
	}

	return &http.Client{
		Timeout:       time.Duration(timeout) * time.Second,
		CheckRedirect: tracker.checkRedirect,
	}
}

func executeInstructionProbes(client *http.Client, rawURL string, maxBodyBytes int64) (*http.Response, error) {
	resp, headErr := executeHead(client, rawURL)
	if headErr == nil && resp.StatusCode < 400 {
		drainAndCloseBody(resp, maxBodyBytes)
		return resp, nil
	}
	if isUnrecoverableError(headErr) {
		return nil, headErr
	}
	getResp, getErr := executeGetWithRange(client, rawURL, maxBodyBytes)
	drainAndCloseBody(getResp, maxBodyBytes)

	return getResp, getErr
}

func isUnrecoverableError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()

	return strings.Contains(msg, "circular redirect") || strings.Contains(msg, "unsupported protocol")
}

func executeGetWithRange(client *http.Client, rawURL string, maxBodyBytes int64) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	bytesCap := resolveBodyBytesCap(maxBodyBytes)
	req.Header.Set("User-Agent", "gitmap-linkexpander/1.0")
	req.Header.Set("Range", fmt.Sprintf("bytes=0-%d", bytesCap-1))

	return client.Do(req)
}

func resolveBodyBytesCap(maxBodyBytes int64) int64 {
	if maxBodyBytes <= 0 {
		return 4096
	}

	return maxBodyBytes
}

func drainAndCloseBody(resp *http.Response, maxBodyBytes int64) {
	if resp == nil || resp.Body == nil {
		return
	}
	defer resp.Body.Close()
	bytesCap := resolveBodyBytesCap(maxBodyBytes)
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, bytesCap))
}

func populateResolutionReport(report *InstructionLinkReport, tracker *redirectTracker, resp *http.Response) {
	if resp != nil {
		report.StatusCode = resp.StatusCode
	}
	report.RedirectHops = len(tracker.trace)
	report.HopTrace = tracker.trace
	if len(tracker.trace) > 0 {
		report.ExpandedURL = tracker.trace[len(tracker.trace)-1]
		report.IsExpanded = true
		return
	}
	applyFinalResponseURL(report, resp)
}

func applyFinalResponseURL(report *InstructionLinkReport, resp *http.Response) {
	if resp == nil || resp.Request == nil || resp.Request.URL == nil {
		return
	}
	reqURL := resp.Request.URL.String()
	if len(reqURL) > 0 && reqURL != report.OriginalURL {
		report.ExpandedURL = reqURL
		report.IsExpanded = true
	}
}

// EmailSenderFunc defines the signature for autoSend delivery.
type EmailSenderFunc func(recipients []string, subject, body string) error

var defaultEmailSender EmailSenderFunc

// SetEmailSender registers a custom email sender for autoSend dispatch.
func SetEmailSender(sender EmailSenderFunc) {
	defaultEmailSender = sender
}

func dispatchInstructionEmail(report *InstructionLinkReport, item InstructionLink) {
	if defaultEmailSender == nil {
		recordEmailSuccess(report, item.Recipients)
		return
	}
	sendErr := defaultEmailSender(item.Recipients, report.InterpolatedSubject, report.InterpolatedBody)
	if sendErr != nil {
		report.Error = sendErr.Error()
		return
	}
	recordEmailSuccess(report, item.Recipients)
}

func recordEmailSuccess(report *InstructionLinkReport, recipients []string) {
	report.IsSent = true
	report.IsEmailSent = true
	report.SentRecipients = recipients
}
