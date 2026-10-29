package linkexpander

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
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
		clean := strings.TrimRight(raw, ".,;:!?)>\"'")
		if _, hasSeen := seen[clean]; hasSeen {
			continue
		}
		expanded := ExpandURL(clean)
		seen[clean] = expanded
		links = append(links, expanded)
		result = applyExpandedURL(result, clean, expanded)
	}

	return result, links
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

	return ExpandedLink{
		Original:   rawURL,
		Expanded:   finalURL,
		IsExpanded: hasRedirect,
		StatusCode: resp.StatusCode,
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
