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

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	imgURLRegex          = regexp.MustCompile(`(?i)<img[^>]+(?:id="screenshot-image"|class="screenshot__image")[^>]+src="([^"]+)"`)
	ogImgRegex           = regexp.MustCompile(`(?i)property="og:image"\s+content="([^"]+)"`)
	lightshotDirectRegex = regexp.MustCompile(`(?i)https://(?:image\.prntscr\.com|img\.lightshot\.app)/[^\s"'<>]+`)
)

func RunAssetPrnt(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs || args[0] == "-h" || args[0] == "--help" {
		printAssetPrntHelp()
		return nil
	}

	url := args[0]
	destPath := ""
	if len(args) > 1 {
		destPath = args[1]
	}

	return DownloadPrntAsset(url, destPath)
}

func printAssetPrntHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap asset download-prnt <url> [destination-path]")
	fmt.Println("  gitmap prnt <url> [destination-path]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  Downloads images from LightShot / prnt.sc URLs and saves them to local assets and brain.")
}

func DownloadPrntAsset(rawURL, destPath string) error {
	imgURL, err := resolveImageURL(rawURL)
	if err != nil {
		return err
	}

	data, err := fetchImageBytes(imgURL)
	if err != nil {
		return fmt.Errorf("failed to download image from %s: %w", imgURL, err)
	}

	slug := extractSlugFromURL(rawURL)
	finalPath := resolveDestinationPath(slug, destPath)
	_ = os.MkdirAll(filepath.Dir(finalPath), 0755)

	if writeErr := os.WriteFile(finalPath, data, 0644); writeErr != nil {
		return fmt.Errorf("failed to save image to %s: %w", finalPath, writeErr)
	}

	printAssetDownloadReport(rawURL, imgURL, finalPath, len(data))
	maybeSyncToActiveBrain(slug, data)
	return nil
}

func printAssetDownloadReport(rawURL, imgURL, finalPath string, size int) {
	fmt.Printf("\n%s✔ Successfully downloaded screenshot asset:%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  • Source URL:  %s%s%s\n", constants.ColorCyan, rawURL, constants.ColorReset)
	fmt.Printf("  • Direct URL:  %s\n", imgURL)
	fmt.Printf("  • Saved Path:  %s%s%s (%d bytes)\n", constants.ColorGreen, finalPath, constants.ColorReset, size)
}

func resolveImageURL(rawURL string) (string, error) {
	isDirect := lightshotDirectRegex.MatchString(rawURL)
	if isDirect {
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
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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
	if m := lightshotDirectRegex.FindString(html); len(m) > 0 {
		return m, nil
	}

	return "", fmt.Errorf("could not locate image source URL on %s", rawURL)
}

func fetchImageBytes(imgURL string) ([]byte, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func extractSlugFromURL(rawURL string) string {
	parts := strings.Split(strings.TrimRight(rawURL, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return fmt.Sprintf("screenshot_%d", time.Now().Unix())
}

func resolveDestinationPath(slug, customDest string) string {
	hasCustom := len(customDest) > 0
	if hasCustom {
		return customDest
	}

	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.png", slug, timestamp)

	return filepath.Join("assets", "screenshots", fileName)
}

func maybeSyncToActiveBrain(slug string, data []byte) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	brainBase := filepath.Join(home, ".gemini", "antigravity", "brain")
	entries, err := os.ReadDir(brainBase)
	hasEntries := err == nil && len(entries) > 0
	if !hasEntries {
		return
	}

	latestDir := findLatestBrainDir(brainBase, entries)
	hasLatestDir := len(latestDir) > 0
	if hasLatestDir {
		brainDest := filepath.Join(latestDir, fmt.Sprintf("screenshot_%s.png", slug))
		_ = os.WriteFile(brainDest, data, 0644)
		fmt.Printf("  • Brain Sync:  %s%s%s\n", constants.ColorGreen, brainDest, constants.ColorReset)
	}
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
