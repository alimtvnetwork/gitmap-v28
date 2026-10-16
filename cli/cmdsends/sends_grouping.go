// Package cmdsends provides semantic change grouping and commit message synthesis.
package cmdsends

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// FileCategory defines semantic classification categories for changed files.
type FileCategory string

const (
	// CategoryStyle classifies CSS and styling assets.
	CategoryStyle FileCategory = "style"
	// CategoryCore classifies core programming code logic.
	CategoryCore FileCategory = "core"
	// CategoryBoolean classifies boolean conventions and predicate files.
	CategoryBoolean FileCategory = "boolean"
	// CategoryAsset classifies static visual assets, icons, and logos.
	CategoryAsset FileCategory = "asset"
	// CategorySpec classifies documentation and architectural specifications.
	CategorySpec FileCategory = "spec"
	// CategoryOther classifies auxiliary configuration and general files.
	CategoryOther FileCategory = "other"
)

// ChangedFile represents an individual changed file with its semantic classification.
type ChangedFile struct {
	RelativePath string       `json:"relativePath"`
	Category     FileCategory `json:"category"`
}

// CategorySummary aggregates changed files under a semantic category heading.
type CategorySummary struct {
	Category      FileCategory  `json:"category"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	Files         []ChangedFile `json:"files"`
	OverflowCount int           `json:"overflowCount"`
}

// CommitDescriptionPayload encapsulates title and structured categories.
type CommitDescriptionPayload struct {
	Title      string            `json:"title"`
	Categories []CategorySummary `json:"categories"`
}

// MaxFilesInCategory defines the maximum files listed per category before truncation.
const MaxFilesInCategory = 8

// ClassifyFilePath determines the semantic FileCategory for a given relative path.
func ClassifyFilePath(filePath string) FileCategory {
	normalized := filepath.ToSlash(filepath.Clean(filePath))
	lower := strings.ToLower(normalized)
	base := path.Base(lower)
	ext := path.Ext(lower)

	if isSpecPath(lower, ext) {
		return CategorySpec
	}
	if isAssetPath(lower, base, ext) {
		return CategoryAsset
	}
	if isBooleanPath(base) {
		return CategoryBoolean
	}
	if isStylePath(base, ext) {
		return CategoryStyle
	}
	if isCoreCodePath(ext) {
		return CategoryCore
	}

	return CategoryOther
}

func isSpecPath(lower, ext string) bool {
	hasSpecPrefix := strings.HasPrefix(lower, "02-spec/") || strings.Contains(lower, "/02-spec/")
	hasDocsPrefix := strings.HasPrefix(lower, "docs/") || strings.Contains(lower, "/docs/")
	isDocExt := ext == ".md" || ext == ".rst" || ext == ".txt"
	return hasSpecPrefix || hasDocsPrefix || isDocExt
}

func isAssetPath(lower, base, ext string) bool {
	hasAssetPrefix := strings.HasPrefix(lower, "assets/") || strings.Contains(lower, "/assets/")
	hasVisualKeyword := strings.Contains(base, "avatar") || strings.Contains(base, "logo") || strings.Contains(base, "icon")
	isVisualExt := ext == ".png" || ext == ".svg" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".webp" || ext == ".ico"
	return hasAssetPrefix || hasVisualKeyword || isVisualExt
}

func isBooleanPath(base string) bool {
	hasBoolKeyword := strings.Contains(base, "bool")
	isPredicate := isPredicateFileName(base)
	return hasBoolKeyword || isPredicate
}

func isPredicateFileName(base string) bool {
	prefixes := []string{"is_", "has_", "can_", "should_", "is-", "has-", "can-", "should-"}
	for _, pfx := range prefixes {
		if strings.HasPrefix(base, pfx) {
			return true
		}
	}
	return false
}

func isStylePath(base, ext string) bool {
	isStyleExt := ext == ".css" || ext == ".scss" || ext == ".less" || ext == ".pcss"
	hasTailwind := strings.Contains(base, "tailwind")
	return isStyleExt || hasTailwind
}

func isCoreCodePath(ext string) bool {
	switch ext {
	case ".go", ".ts", ".tsx", ".rs", ".py", ".php", ".c", ".cpp", ".h", ".hpp":
		return true
	default:
		return false
	}
}

// CategoryMetadata holds title and narrative descriptions for a category.
type CategoryMetadata struct {
	Title       string
	Description string
}

func resolveCategoryMetadata(cat FileCategory) CategoryMetadata {
	switch cat {
	case CategoryStyle:
		return CategoryMetadata{
			Title:       "Styling & CSS Updates",
			Description: "Modernized responsive layout tokens and component stylesheets.",
		}
	case CategoryCore:
		return CategoryMetadata{
			Title:       "Core Code Logic Changes",
			Description: "Updated core business logic, handlers, and routing interfaces.",
		}
	case CategoryBoolean:
		return CategoryMetadata{
			Title:       "Boolean & Guard Updates",
			Description: "Normalized conditional branching, predicate helpers, and guard returns.",
		}
	case CategoryAsset:
		return CategoryMetadata{
			Title:       "Asset & Image Updates",
			Description: "Synchronized visual iconography, branding assets, and static media.",
		}
	case CategorySpec:
		return CategoryMetadata{
			Title:       "Specification & Documentation",
			Description: "Updated architectural specifications, documentation, and task records.",
		}
	default:
		return CategoryMetadata{
			Title:       "General Project Updates",
			Description: "Updated repository configuration and auxiliary files.",
		}
	}
}

// CategorizeChangedFiles parses porcelain output and groups changed files into categories.
func CategorizeChangedFiles(porcelainOutput string) []CategorySummary {
	lines := strings.Split(porcelainOutput, "\n")
	seen := make(map[string]bool)
	categoryBuckets := make(map[FileCategory][]ChangedFile)

	for _, line := range lines {
		cleanLine := strings.TrimRight(line, "\r")
		if len(cleanLine) < 3 {
			continue
		}

		rawPath := strings.TrimSpace(cleanLine[2:])
		hasArrow := strings.Contains(rawPath, " -> ")
		if hasArrow {
			parts := strings.Split(rawPath, " -> ")
			rawPath = strings.TrimSpace(parts[len(parts)-1])
		}
		rawPath = strings.Trim(rawPath, "\"")

		cleanPath := filepath.ToSlash(filepath.Clean(rawPath))
		cleanPath = strings.TrimPrefix(cleanPath, "./")
		cleanPath = strings.TrimPrefix(cleanPath, "/")

		if cleanPath == "" || cleanPath == "." {
			continue
		}

		key := strings.ToLower(cleanPath)
		if seen[key] {
			continue
		}
		seen[key] = true

		cat := ClassifyFilePath(cleanPath)
		categoryBuckets[cat] = append(categoryBuckets[cat], ChangedFile{
			RelativePath: cleanPath,
			Category:     cat,
		})
	}

	orderedCategories := []FileCategory{
		CategoryStyle,
		CategoryCore,
		CategoryBoolean,
		CategoryAsset,
		CategorySpec,
		CategoryOther,
	}

	var summaries []CategorySummary
	for _, cat := range orderedCategories {
		files := categoryBuckets[cat]
		hasFiles := len(files) > 0
		if hasFiles {
			meta := resolveCategoryMetadata(cat)
			summary := CategorySummary{
				Category:    cat,
				Title:       meta.Title,
				Description: meta.Description,
			}

			if len(files) <= MaxFilesInCategory {
				summary.Files = files
			} else {
				summary.Files = files[:MaxFilesInCategory]
				summary.OverflowCount = len(files) - MaxFilesInCategory
			}

			summaries = append(summaries, summary)
		}
	}

	return summaries
}

// FormatMarkdownCommitBody renders a conventional multi-line commit message with relative markdown links.
func FormatMarkdownCommitBody(title string, categories []CategorySummary) string {
	cleanTitle := strings.TrimSpace(title)
	hasCategories := len(categories) > 0
	if hasCategories {
		var b strings.Builder
		b.WriteString(cleanTitle)
		b.WriteString("\n\n")

		for i, cat := range categories {
			if i > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString("### ")
			b.WriteString(cat.Title)
			b.WriteString("\n")
			b.WriteString(cat.Description)

			for _, f := range cat.Files {
				b.WriteString("\n- [")
				b.WriteString(f.RelativePath)
				b.WriteString("](")
				b.WriteString(f.RelativePath)
				b.WriteString(")")
			}

			hasOverflow := cat.OverflowCount > 0
			if hasOverflow {
				b.WriteString(fmt.Sprintf("\n- ... and %d additional files", cat.OverflowCount))
			}
		}

		return b.String()
	}

	return cleanTitle
}

// BuildGroupedCommitMessage parses porcelain and synthesizes the final commit message.
func BuildGroupedCommitMessage(title, porcelainOutput string) string {
	categories := CategorizeChangedFiles(porcelainOutput)
	return FormatMarkdownCommitBody(title, categories)
}
