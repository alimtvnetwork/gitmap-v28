// Package cmdagy provides Antigravity agent and AGM orchestration helpers.
package cmdagy

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

const (
	// DefaultAssetServerURL is the default local HTTP bridge endpoint for AGM and GitMap.
	DefaultAssetServerURL = "http://127.0.0.1:18188"

	// DefaultMemoryAssetsDir is the project-scoped asset storage directory.
	DefaultMemoryAssetsDir = ".ai-memory/assets"

	// DefaultGlobalAssetsDir is the machine-global fallback asset directory.
	DefaultGlobalAssetsDir = "~/.gitmap/assets"
)

var base64ImageRegex = lazyregex.New(`data:image/(png|jpeg|jpg|webp|gif);base64,([A-Za-z0-9+/=]+)`)

type imageMatchResult struct {
	MatchStart  int
	MatchEnd    int
	Replacement string
	SavedPath   string
}

// ExtractAndStoreBase64Images scans promptText for base64 data URIs, decodes and hashes them,
// saves binaries to targetDir, and rewrites the prompt with standalone URLs and asset comments.
func ExtractAndStoreBase64Images(promptText string, targetDir string) (string, []string, *appfault.AppError) {
	if strings.TrimSpace(promptText) == "" {
		return "", nil, nil
	}

	resolvedDir := resolveTargetDir(targetDir)
	if err := ensureDirExists(resolvedDir); err != nil {
		return "", nil, err
	}

	matches := base64ImageRegex.Compiled().FindAllStringSubmatchIndex(promptText, -1)

	return processImageMatches(promptText, resolvedDir, matches)
}

func processImageMatches(promptText string, targetDir string, matches [][]int) (string, []string, *appfault.AppError) {
	if len(matches) == 0 {
		return promptText, nil, nil
	}

	results, files, err := collectImageReplacements(promptText, targetDir, matches)
	if err != nil {
		return "", nil, err
	}

	rewritten := applyImageReplacements(promptText, results)

	return rewritten, files, nil
}

func collectImageReplacements(promptText string, targetDir string, matches [][]int) ([]imageMatchResult, []string, *appfault.AppError) {
	var results []imageMatchResult
	var files []string

	for _, matchIdx := range matches {
		res, err := processSingleImageMatch(promptText, targetDir, matchIdx)
		if err != nil {
			return nil, nil, err
		}

		results = append(results, res)
		files = appendIfNew(files, res.SavedPath)
	}

	return results, files, nil
}

func applyImageReplacements(promptText string, results []imageMatchResult) string {
	var sb strings.Builder
	lastIdx := 0

	for _, res := range results {
		sb.WriteString(promptText[lastIdx:res.MatchStart])
		sb.WriteString(res.Replacement)
		lastIdx = res.MatchEnd
	}

	sb.WriteString(promptText[lastIdx:])

	return sb.String()
}

func processSingleImageMatch(promptText string, targetDir string, matchIdx []int) (imageMatchResult, *appfault.AppError) {
	ext := strings.ToLower(promptText[matchIdx[2]:matchIdx[3]])
	b64Payload := promptText[matchIdx[4]:matchIdx[5]]

	data, err := base64.StdEncoding.DecodeString(b64Payload)
	if err != nil {
		return imageMatchResult{}, appfault.WrapSimple(err, "decodeBase64Image")
	}

	savedPath, fileName, saveErr := saveImagePayload(targetDir, ext, data)
	if saveErr != nil {
		return imageMatchResult{}, saveErr
	}

	return buildImageMatchResult(promptText, targetDir, fileName, savedPath, matchIdx), nil
}

func saveImagePayload(targetDir string, ext string, data []byte) (string, string, *appfault.AppError) {
	hash := sha256.Sum256(data)
	hashPrefix := hex.EncodeToString(hash[:])[:16]
	fileName := fmt.Sprintf("img-%s.%s", hashPrefix, ext)
	savedPath := filepath.Join(targetDir, fileName)

	if err := os.WriteFile(savedPath, data, 0644); err != nil {
		return "", "", appfault.WrapSimple(err, "saveImagePayload")
	}

	return savedPath, fileName, nil
}

func buildImageMatchResult(promptText, targetDir, fileName, savedPath string, matchIdx []int) imageMatchResult {
	assetRef := resolveAssetReference(targetDir, fileName)
	url := fmt.Sprintf("%s/assets/%s", DefaultAssetServerURL, fileName)
	comment := fmt.Sprintf("<!-- asset: %s -->", assetRef)

	replacement, end := formatReplacementWithBoundary(promptText, url, comment, matchIdx[1])

	return imageMatchResult{
		MatchStart:  matchIdx[0],
		MatchEnd:    end,
		Replacement: replacement,
		SavedPath:   savedPath,
	}
}

func formatReplacementWithBoundary(promptText, url, comment string, matchEnd int) (string, int) {
	if matchEnd < len(promptText) && promptText[matchEnd] == ')' {
		return fmt.Sprintf("%s)\n%s", url, comment), matchEnd + 1
	}

	return fmt.Sprintf("%s\n%s", url, comment), matchEnd
}

func resolveTargetDir(targetDir string) string {
	if strings.TrimSpace(targetDir) != "" {
		return targetDir
	}

	return DefaultMemoryAssetsDir
}

func ensureDirExists(dir string) *appfault.AppError {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return appfault.WrapSimple(err, "ensureDirExists")
	}

	return nil
}

func resolveAssetReference(targetDir, fileName string) string {
	combined := filepath.Join(targetDir, fileName)
	rel, err := filepath.Rel(".", combined)
	if err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}

	return filepath.ToSlash(combined)
}

func appendIfNew(items []string, item string) []string {
	for _, existing := range items {
		if existing == item {
			return items
		}
	}

	return append(items, item)
}

// FormatImageMarkdownStandalone produces the dual-reference markdown string for an image.
func FormatImageMarkdownStandalone(imagePath string, assetServerUrl string) string {
	if strings.TrimSpace(imagePath) == "" {
		return ""
	}

	serverURL := resolveServerURL(assetServerUrl)
	fileName := filepath.Base(imagePath)
	refPath := resolveAssetReference(filepath.Dir(imagePath), fileName)

	return fmt.Sprintf("![Screenshot](%s/assets/%s)\n<!-- asset: %s -->", serverURL, fileName, refPath)
}

func resolveServerURL(rawURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(rawURL), "/")
	if trimmed != "" {
		return trimmed
	}

	return DefaultAssetServerURL
}
