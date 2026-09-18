package cmdagy

import (
	"os"
	"path/filepath"
)

// LoadCanonicalRcaPrompt loads RCA bug fix or CI/CD fix prompt template.
func LoadCanonicalRcaPrompt(customPath string, isNoRelease bool) (string, string) {
	if prompt, path, hasCustom := readCustomPrompt(customPath); hasCustom {
		return prompt, path
	}

	for _, candidate := range selectRcaPromptCandidates(isNoRelease) {
		data, err := os.ReadFile(candidate)
		if err == nil && len(data) > 0 {
			return string(data), candidate
		}
	}

	return defaultRcaFixPromptFallback, "embedded-rca-fallback"
}

func selectRcaPromptCandidates(isNoRelease bool) []string {
	rootDir := resolveProjectRootDir()
	if isNoRelease {
		return []string{
			filepath.Join(rootDir, "01-prompts/07-bug-fix/01-fix-with-rca.md"),
			filepath.Join(rootDir, "01-prompts/16-ci-cd/01-ci-cd-fix.md"),
		}
	}

	return []string{
		filepath.Join(rootDir, "01-prompts/07-bug-fix/01-fix-with-rca.md"),
		filepath.Join(rootDir, "01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md"),
		filepath.Join(rootDir, "01-prompts/16-ci-cd/01-ci-cd-fix.md"),
	}
}

func selectPromptPath(isNoRelease bool) string {
	rootDir := resolveProjectRootDir()
	if isNoRelease {
		return filepath.Join(rootDir, "01-prompts/16-ci-cd/01-ci-cd-fix.md")
	}

	return filepath.Join(rootDir, "01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md")
}

func loadCicdFixPrompt(customPath string, isNoRelease bool) (string, string) {
	if prompt, path, hasCustom := readCustomPrompt(customPath); hasCustom {
		return prompt, path
	}

	targetPath := selectPromptPath(isNoRelease)
	data, err := os.ReadFile(targetPath)
	if err == nil {
		return string(data), targetPath
	}

	return defaultCicdFixWithReleasePromptFallback, "embedded-fallback"
}

func readCustomPrompt(customPath string) (string, string, bool) {
	if len(customPath) == 0 {
		return "", "", false
	}
	data, err := os.ReadFile(customPath)
	if err != nil {
		return "", "", false
	}

	return string(data), customPath, true
}
