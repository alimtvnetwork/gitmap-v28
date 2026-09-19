package cmdautomation

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func findPathViolationsInText(file, content string) []PathViolation {
	var vios []PathViolation
	lines := strings.Split(content, "\n")
	for idx, line := range lines {
		vios = append(vios, checkLinePathViolations(file, idx+1, line)...)
	}
	return vios
}

func checkLinePathViolations(file string, lineNum int, line string) []PathViolation {
	var vios []PathViolation
	matches := reWinDrive.FindAllString(line, -1)
	matches = append(matches, reFileUri.FindAllString(line, -1)...)
	for _, m := range matches {
		if isIgnoredPathMarker(m) {
			continue
		}
		vios = append(vios, PathViolation{
			File:       file,
			LineNumber: lineNum,
			RawPath:    m,
		})
	}
	return vios
}

func isIgnoredPathMarker(m string) bool {
	return strings.Contains(m, `\\.\`) ||
		strings.HasSuffix(m, ".") ||
		strings.Contains(m, "http://") ||
		strings.Contains(m, "https://")
}

func applyPathFixIfEnabled(file, content string, isFix bool) int {
	if !isFix {
		return 0
	}
	cleaned, count := sanitizePathsInText(content)
	if count == 0 {
		return 0
	}
	if err := os.WriteFile(file, []byte(cleaned), 0o644); err != nil {
		return 0
	}
	return count
}

func sanitizePathsInText(content string) (string, int) {
	root := findRepoRoot()
	cleanRoot := filepath.ToSlash(root)
	patterns := []string{
		"file:///" + cleanRoot + "/",
		cleanRoot + "/",
		root + "\\",
	}
	modified := content
	count := 0
	for _, pat := range patterns {
		if strings.Contains(modified, pat) {
			modified = strings.ReplaceAll(modified, pat, "")
			count++
		}
	}
	return modified, count
}

func collectTextFiles(root string, extensions []string) ([]string, *apperror.AppError) {
	extMap := buildNormalizedExtMap(extensions)
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() && isExcludedDir(d.Name()) {
			return filepath.SkipDir
		}
		if !d.IsDir() && isTargetTextFile(path, extMap) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, apperror.WrapSimple(err, "walk text files")
	}
	return files, nil
}

func buildNormalizedExtMap(extensions []string) map[string]bool {
	m := make(map[string]bool)
	for _, ext := range extensions {
		norm := strings.ToLower(strings.TrimSpace(ext))
		if !strings.HasPrefix(norm, ".") {
			norm = "." + norm
		}
		m[norm] = true
	}
	return m
}

func isTargetTextFile(path string, extMap map[string]bool) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if len(extMap) > 0 {
		return extMap[ext]
	}
	return ext == ".md" || ext == ".txt" || ext == ".go" || ext == ".ts" ||
		ext == ".py" || ext == ".json" || ext == ".yml" || ext == ".yaml"
}
