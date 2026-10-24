package enums

import (
	"fmt"
	"strings"
)

// SuggestionCategoryType classifies the underlying intent of a suggested action.
type SuggestionCategoryType string

const (
	// CategoryTypo indicates a mistyped command, flag, or token.
	CategoryTypo SuggestionCategoryType = "typo"
	// CategorySubcommand indicates a valid child subcommand recommendation.
	CategorySubcommand SuggestionCategoryType = "subcommand"
	// CategoryFlag indicates a supported flag recommendation matching input token.
	CategoryFlag SuggestionCategoryType = "flag"
	// CategoryRemediation indicates a prescriptive diagnostic step to fix a failure.
	CategoryRemediation SuggestionCategoryType = "remediation"
	// CategoryIntent indicates a semantic alias or alternative intent (e.g. docker -> cluster).
	CategoryIntent SuggestionCategoryType = "intent"
	// CategoryResource indicates a repository, branch, or node target suggestion.
	CategoryResource SuggestionCategoryType = "resource"
)

// IsValid reports whether the value matches a known variant.
func (s SuggestionCategoryType) IsValid() bool {
	switch s {
	case CategoryTypo, CategorySubcommand, CategoryFlag, CategoryRemediation, CategoryIntent, CategoryResource:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (s SuggestionCategoryType) String() string {
	return string(s)
}

// ParseSuggestionCategoryType parses a raw string into a validated SuggestionCategoryType.
// ParseSuggestionCategoryType returns a plain error (not
// *apperror.AppError) by recorded exemption: apperror imports cli/diag,
// so an AppError Parse here would close an import cycle
// (diag -> diag/enums -> apperror -> diag).
func ParseSuggestionCategoryType(s string) (SuggestionCategoryType, error) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllSuggestionCategoryTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "categorytypo":
		return CategoryTypo, nil
	case "typo":
		return CategoryTypo, nil
	case "categorysubcommand":
		return CategorySubcommand, nil
	case "subcommand":
		return CategorySubcommand, nil
	case "categoryflag":
		return CategoryFlag, nil
	case "flag":
		return CategoryFlag, nil
	case "categoryremediation":
		return CategoryRemediation, nil
	case "remediation":
		return CategoryRemediation, nil
	case "categoryintent":
		return CategoryIntent, nil
	case "intent":
		return CategoryIntent, nil
	case "categoryresource":
		return CategoryResource, nil
	case "resource":
		return CategoryResource, nil
	}
	return "", fmt.Errorf("invalid suggestioncategorytype: %q", s)
}

// AllSuggestionCategoryTypes returns all valid variants.
func AllSuggestionCategoryTypes() []SuggestionCategoryType {
	return []SuggestionCategoryType{
		CategoryTypo,
		CategorySubcommand,
		CategoryFlag,
		CategoryRemediation,
		CategoryIntent,
		CategoryResource,
	}
}

// IsTypo reports whether the value is CategoryTypo.
func (s SuggestionCategoryType) IsTypo() bool {
	return s == CategoryTypo
}

// IsSubcommand reports whether the value is CategorySubcommand.
func (s SuggestionCategoryType) IsSubcommand() bool {
	return s == CategorySubcommand
}

// IsFlag reports whether the value is CategoryFlag.
func (s SuggestionCategoryType) IsFlag() bool {
	return s == CategoryFlag
}

// IsRemediation reports whether the value is CategoryRemediation.
func (s SuggestionCategoryType) IsRemediation() bool {
	return s == CategoryRemediation
}

// IsIntent reports whether the value is CategoryIntent.
func (s SuggestionCategoryType) IsIntent() bool {
	return s == CategoryIntent
}

// IsResource reports whether the value is CategoryResource.
func (s SuggestionCategoryType) IsResource() bool {
	return s == CategoryResource
}

// HasCategory reports whether the value is non-empty and valid.
func (s SuggestionCategoryType) HasCategory() bool {
	return s.IsValid()
}
