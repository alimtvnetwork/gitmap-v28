package cmdasset

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	imgURLRegex      = regexp.MustCompile(`(?i)<img[^>]+(?:id="screenshot-image"|class="screenshot__image")[^>]+src="([^"]+)"`)
	ogImgRegex       = regexp.MustCompile(`(?i)property="og:image"\s+content="([^"]+)"`)
	directImageRegex = regexp.MustCompile(`(?i)https?://(?:image\.prntscr\.com|img\.lightshot\.app|api\.prntscr\.com|i\.imgur\.com)/[^\s"'<>]+`)
	fallbackImgRegex = regexp.MustCompile(`(?i)https?://(?:image\.prntscr\.com|api\.prntscr\.com|i\.imgur\.com)/[^\s"'<>]+\.(?:png|jpg|jpeg)`)
)

func resolveImageURL(rawURL string) (string, error) {
	if directImageRegex.MatchString(rawURL) {
		return rawURL, nil
	}

	html, err := fetchPageHTML(rawURL)
	if err != nil {
		return "", err
	}

	return parseImageURLFromHTML(rawURL, html)
}

func fetchPageHTML(targetURL string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error %d fetching %s", resp.StatusCode, targetURL)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func parseImageURLFromHTML(rawURL, html string) (string, error) {
	if m := imgURLRegex.FindStringSubmatch(html); len(m) > 1 {
		return m[1], nil
	}
	if m := ogImgRegex.FindStringSubmatch(html); len(m) > 1 {
		return m[1], nil
	}
	if m := directImageRegex.FindString(html); len(m) > 0 {
		return m, nil
	}
	if m := fallbackImgRegex.FindString(html); len(m) > 0 {
		return m, nil
	}

	return "", fmt.Errorf("could not locate image source URL on %s", rawURL)
}

func fetchImageBytes(imgURL string) ([]byte, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodGet, imgURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d downloading image from %s", resp.StatusCode, imgURL)
	}

	return io.ReadAll(resp.Body)
}

func extractSlugFromURL(rawURL string) string {
	parts := strings.Split(strings.TrimRight(rawURL, "/"), "/")
	if len(parts) > 0 {
		clean := parts[len(parts)-1]
		if idx := strings.Index(clean, "?"); idx != -1 {
			clean = clean[:idx]
		}
		if clean != "" {
			return clean
		}
	}

	return fmt.Sprintf("screenshot_%d", time.Now().Unix())
}

func resolveDestinationPath(slug, customDest, outputDir string) string {
	if len(customDest) > 0 {
		return customDest
	}

	baseDir := outputDir
	if baseDir == "" {
		baseDir = filepath.Join("assets", "screenshots")
	}

	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.png", slug, timestamp)

	return filepath.Join(baseDir, fileName)
}

func maybeSyncToActiveBrain(slug string, data []byte) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	brainBase := filepath.Join(home, ".gemini", "antigravity", "brain")
	entries, err := os.ReadDir(brainBase)
	if err != nil || len(entries) == 0 {
		return ""
	}

	latestDir := findLatestBrainDir(brainBase, entries)
	if len(latestDir) == 0 {
		return ""
	}

	brainDest := filepath.Join(latestDir, fmt.Sprintf("screenshot_%s.png", slug))
	if writeErr := os.WriteFile(brainDest, data, 0644); writeErr != nil {
		return ""
	}

	return brainDest
}

func findLatestBrainDir(brainBase string, entries []os.DirEntry) string {
	var latestDir string
	var latestMod time.Time

	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		p := filepath.Join(brainBase, e.Name())
		fi, fiErr := os.Stat(p)
		if fiErr == nil && fi.ModTime().After(latestMod) {
			latestMod = fi.ModTime()
			latestDir = p
		}
	}

	return latestDir
}
