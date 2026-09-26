package cmdprompttemplate

import (
	"crypto/rand"
	"math/big"
	"os"
	"strings"
	"time"
)

// Category constants for template classifications.
const (
	CategoryDefault = "default"
	CategoryUIUX    = "ui-ux"
	CategorySponsor = "sponsor"
	CategoryPR      = "pr-descriptions"
)

// GetAllCatalogTemplates aggregates all built-in catalog templates.
func GetAllCatalogTemplates() []PromptTemplate {
	var all []PromptTemplate
	all = append(all, GetDefaultCategoryTemplates()...)
	all = append(all, GetUIUXCategoryTemplates()...)
	all = append(all, GetSponsorCategoryTemplates()...)
	all = append(all, GetPRDescriptionTemplates()...)

	return all
}

// GetCategoryTemplates returns all templates for a specific category.
func GetCategoryTemplates(category string) []PromptTemplate {
	norm := normalizeCategoryName(category)
	switch norm {
	case CategoryDefault:
		return GetDefaultCategoryTemplates()
	case CategoryUIUX:
		return GetUIUXCategoryTemplates()
	case CategorySponsor:
		return GetSponsorCategoryTemplates()
	case CategoryPR:
		return GetPRDescriptionTemplates()
	default:
		return nil
	}
}

func normalizeCategoryName(raw string) string {
	low := strings.ToLower(strings.TrimSpace(raw))
	low = strings.TrimPrefix(low, "--")
	switch low {
	case "default", "def", "d":
		return CategoryDefault
	case "ui-ux", "uiux", "ui", "ux", "uu":
		return CategoryUIUX
	case "sponsor", "sp", "riseup", "riseup-asia":
		return CategorySponsor
	case "pr", "pr-desc", "pr-description", "pr-descriptions", "prs":
		return CategoryPR
	default:
		return low
	}
}

// ResolveTemplateContent resolves a template from category, slug, CSV list, or direct file.
func ResolveTemplateContent(token string, index int) string {
	trimmed := strings.TrimSpace(token)
	if len(trimmed) == 0 {
		return ""
	}

	if strings.Contains(trimmed, ",") {
		return pickFromCSV(trimmed, index)
	}

	if fileContent, isFile := tryReadFileContent(trimmed); isFile {
		return fileContent
	}

	templates := GetCategoryTemplates(trimmed)
	if len(templates) > 0 {
		idx := safeIndex(index, len(templates))

		return templates[idx].Content
	}

	for _, tmpl := range GetAllCatalogTemplates() {
		isMatch := strings.EqualFold(tmpl.Slug, trimmed) || strings.EqualFold(tmpl.ID, trimmed)
		if isMatch {
			return tmpl.Content
		}
	}

	return trimmed
}

// ResolveRandomTemplateContent randomly picks a template from category, CSV, or slug.
func ResolveRandomTemplateContent(token string) string {
	trimmed := strings.TrimSpace(token)
	if len(trimmed) == 0 {
		return ""
	}

	if strings.Contains(trimmed, ",") {
		parts := strings.Split(trimmed, ",")
		rIdx := secureRandomInt(len(parts))

		return strings.TrimSpace(parts[rIdx])
	}

	templates := GetCategoryTemplates(trimmed)
	if len(templates) > 0 {
		rIdx := secureRandomInt(len(templates))

		return templates[rIdx].Content
	}

	return ResolveTemplateContent(token, 0)
}

func pickFromCSV(csv string, index int) string {
	parts := strings.Split(csv, ",")
	cleanParts := make([]string, 0, len(parts))
	for _, p := range parts {
		c := strings.TrimSpace(p)
		if len(c) > 0 {
			cleanParts = append(cleanParts, c)
		}
	}
	if len(cleanParts) == 0 {
		return ""
	}
	idx := safeIndex(index, len(cleanParts))

	return cleanParts[idx]
}

func tryReadFileContent(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err == nil && len(data) > 0 {
		return strings.TrimSpace(string(data)), true
	}

	return "", false
}

func safeIndex(index, length int) int {
	if length <= 0 {
		return 0
	}
	if index < 0 {
		index = -index
	}

	return index % length
}

func secureRandomInt(max int) int {
	if max <= 1 {
		return 0
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return int(time.Now().UnixNano() % int64(max))
	}

	return int(nBig.Int64())
}
