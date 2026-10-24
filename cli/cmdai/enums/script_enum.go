package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ScriptCategoryType defines taxonomy categories for repository automation fspath.
type ScriptCategoryType string

const (
	CategoryCore        ScriptCategoryType = "core"
	CategoryFormatting  ScriptCategoryType = "formatting"
	CategoryGuideline   ScriptCategoryType = "guideline"
	CategoryCicd        ScriptCategoryType = "cicd"
	CategoryPaths       ScriptCategoryType = "paths"
	CategoryAudit       ScriptCategoryType = "audit"
	CategoryDiscovery   ScriptCategoryType = "discovery"
	CategoryMaintenance ScriptCategoryType = "maintenance"
	CategoryPlans       ScriptCategoryType = "plans"
	CategoryRelease     ScriptCategoryType = "release"
	CategoryDatabase    ScriptCategoryType = "database"
	CategoryUnknown     ScriptCategoryType = "unknown"
)

// FixTargetType defines standard auto-remediation targets across repository domains.
type FixTargetType string

const (
	FixTargetGuidelines FixTargetType = "guidelines"
	FixTargetNewlines   FixTargetType = "newlines"
	FixTargetPaths      FixTargetType = "paths"
	FixTargetNaming     FixTargetType = "naming"
	FixTargetEncoding   FixTargetType = "encoding"
	FixTargetSpelling   FixTargetType = "spelling"
	FixTargetPlans      FixTargetType = "plans"
	FixTargetAll        FixTargetType = "all"
)

// ScriptTemplateType defines the archetype of script to scaffold.
type ScriptTemplateType string

const (
	TemplateLinter    ScriptTemplateType = "linter"
	TemplateFixer     ScriptTemplateType = "fixer"
	TemplateAuditor   ScriptTemplateType = "auditor"
	TemplateChecker   ScriptTemplateType = "checker"
	TemplateGenerator ScriptTemplateType = "generator"
	TemplateUtil      ScriptTemplateType = "util"
)

// IsValid reports whether the script category matches a known category.
func (c ScriptCategoryType) IsValid() bool {
	switch c {
	case CategoryCore, CategoryFormatting, CategoryGuideline, CategoryCicd,
		CategoryPaths, CategoryAudit, CategoryDiscovery, CategoryMaintenance,
		CategoryPlans, CategoryRelease, CategoryDatabase, CategoryUnknown:
		return true
	default:
		return false
	}
}

// String returns the string representation of the script category.
func (c ScriptCategoryType) String() string {
	return string(c)
}

// ParseScriptCategoryType parses a raw string into a validated ScriptCategoryType.
func ParseScriptCategoryType(s string) (ScriptCategoryType, *apperror.AppError) {
	val := ScriptCategoryType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseScriptCategoryType", "E1000", fmt.Sprintf("invalid script category: %q", s)).WithContext("value", s)
}

// AllScriptCategoryTypes returns all valid ScriptCategoryType variants.
func AllScriptCategoryTypes() []ScriptCategoryType {
	return []ScriptCategoryType{
		CategoryCore, CategoryFormatting, CategoryGuideline, CategoryCicd,
		CategoryPaths, CategoryAudit, CategoryDiscovery, CategoryMaintenance,
		CategoryPlans, CategoryRelease, CategoryDatabase, CategoryUnknown,
	}
}

// IsValid reports whether the fix target matches a known target.
func (f FixTargetType) IsValid() bool {
	switch f {
	case FixTargetGuidelines, FixTargetNewlines, FixTargetPaths, FixTargetNaming,
		FixTargetEncoding, FixTargetSpelling, FixTargetPlans, FixTargetAll:
		return true
	default:
		return false
	}
}

// String returns the string representation of the fix target.
func (f FixTargetType) String() string {
	return string(f)
}

// ParseFixTargetType parses a raw string into a validated FixTargetType.
func ParseFixTargetType(s string) (FixTargetType, *apperror.AppError) {
	val := FixTargetType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseFixTargetType", "E1000", fmt.Sprintf("invalid fix target: %q", s)).WithContext("value", s)
}

// AllFixTargetTypes returns all valid FixTargetType variants.
func AllFixTargetTypes() []FixTargetType {
	return []FixTargetType{
		FixTargetGuidelines, FixTargetNewlines, FixTargetPaths, FixTargetNaming,
		FixTargetEncoding, FixTargetSpelling, FixTargetPlans, FixTargetAll,
	}
}

// IsValid reports whether the script template matches a known template.
func (t ScriptTemplateType) IsValid() bool {
	switch t {
	case TemplateLinter, TemplateFixer, TemplateAuditor, TemplateChecker, TemplateGenerator, TemplateUtil:
		return true
	default:
		return false
	}
}

// String returns the string representation of the script template.
func (t ScriptTemplateType) String() string {
	return string(t)
}

// ParseScriptTemplateType parses a raw string into a validated ScriptTemplateType.
func ParseScriptTemplateType(s string) (ScriptTemplateType, *apperror.AppError) {
	val := ScriptTemplateType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseScriptTemplateType", "E1000", fmt.Sprintf("invalid script template: %q", s)).WithContext("value", s)
}

// AllScriptTemplateTypes returns all valid ScriptTemplateType variants.
func AllScriptTemplateTypes() []ScriptTemplateType {
	return []ScriptTemplateType{TemplateLinter, TemplateFixer, TemplateAuditor, TemplateChecker, TemplateGenerator, TemplateUtil}
}
