package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	rootDir := "d:/wp-work/riseup-asia/gitmap"

	magicStringRe := regexp.MustCompile(`(==|!=)\s*("[a-zA-Z0-9_\-\s]+"|[0-9]{2,})`)
	emptyCatchRe := regexp.MustCompile(`catch\s*\([^)]*\)\s*\{\s*\}`)
	bannedShortRe := regexp.MustCompile(`\b(arr|cb|fn|el|msg|ctx|obj|val)\b`)
	missingTypeSuffixReGo := regexp.MustCompile(`type\s+([A-Z][a-zA-Z0-9_]*)\s+(struct|string|int|uint)`)
	missingTypeSuffixReTS := regexp.MustCompile(`(type|interface|enum)\s+([A-Z][a-zA-Z0-9_]*)`)

	var magicMatches []string
	var swallowedMatches []string
	var shortMatches []string
	var suffixMatches []string

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		
		ext := filepath.Ext(path)
		if ext != ".go" && ext != ".ts" && ext != ".tsx" {
			return nil
		}
		
		if strings.Contains(path, "vendor") || strings.Contains(path, "node_modules") || strings.Contains(path, ".git") || strings.Contains(path, ".lovable") {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(contentBytes)
		lines := strings.Split(content, "\n")

		relPath, _ := filepath.Rel(rootDir, path)
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		for i, line := range lines {
			lineNum := i + 1

			if len(magicMatches) < 100 && magicStringRe.MatchString(line) {
				magicMatches = append(magicMatches, fmt.Sprintf("- [ ] `%s:%d` - %s", relPath, lineNum, strings.TrimSpace(line)))
			}

			if len(shortMatches) < 100 && bannedShortRe.MatchString(line) {
				shortMatches = append(shortMatches, fmt.Sprintf("- [ ] `%s:%d` - %s", relPath, lineNum, strings.TrimSpace(line)))
			}

			if ext == ".go" {
				if len(suffixMatches) < 100 {
					m := missingTypeSuffixReGo.FindStringSubmatch(line)
					if len(m) > 1 {
						name := m[1]
						if !strings.HasSuffix(name, "Type") && !strings.HasSuffix(name, "Enum") {
							suffixMatches = append(suffixMatches, fmt.Sprintf("- [ ] `%s:%d` - %s", relPath, lineNum, strings.TrimSpace(line)))
						}
					}
				}
			} else if ext == ".ts" || ext == ".tsx" {
				if len(suffixMatches) < 100 {
					m := missingTypeSuffixReTS.FindStringSubmatch(line)
					if len(m) > 2 {
						name := m[2]
						if !strings.HasSuffix(name, "Type") && !strings.HasSuffix(name, "Enum") {
							suffixMatches = append(suffixMatches, fmt.Sprintf("- [ ] `%s:%d` - %s", relPath, lineNum, strings.TrimSpace(line)))
						}
					}
				}
			}
			
			if (ext == ".ts" || ext == ".tsx") && emptyCatchRe.MatchString(line) {
				if len(swallowedMatches) < 100 {
					swallowedMatches = append(swallowedMatches, fmt.Sprintf("- [ ] `%s:%d` - Empty catch block", relPath, lineNum))
				}
			}
			
			if ext == ".go" {
				matched, _ := regexp.MatchString(`^\s*return\s+([a-zA-Z0-9_]+,\s*)*err\s*$`, line)
				if matched && len(swallowedMatches) < 100 {
					swallowedMatches = append(swallowedMatches, fmt.Sprintf("- [ ] `%s:%d` - Missing error wrapping: %s", relPath, lineNum, strings.TrimSpace(line)))
				}
			}
		}

		return nil
	})

	os.MkdirAll("d:/wp-work/riseup-asia/gitmap/.lovable/plans/subtasks/01-coding-guideline-fixes", 0755)

	writeTask("06-magic-strings.md", "Magic Strings and Numbers", magicMatches)
	writeTask("07-swallowed-errors.md", "Swallowed Errors and Missing Context", swallowedMatches)
	writeTask("08-banned-short-identifiers.md", "Banned Short Identifiers", shortMatches)
	writeTask("09-missing-type-suffixes.md", "Missing Enum or Type Suffixes", suffixMatches)

	fmt.Println("Audit complete. Created subtasks.")
}

func writeTask(filename, title string, matches []string) {
	if len(matches) > 50 {
		matches = matches[:50]
	}
	content := fmt.Sprintf("# Subtask: %s\n\nFind and fix instances of %s.\n\n", title, strings.ToLower(title))
	if len(matches) == 0 {
		content += "No instances found.\n"
	} else {
		content += strings.Join(matches, "\n") + "\n"
	}

	path := filepath.Join("d:/wp-work/riseup-asia/gitmap/.lovable/plans/subtasks/01-coding-guideline-fixes", filename)
	os.WriteFile(path, []byte(content), 0644)
}
