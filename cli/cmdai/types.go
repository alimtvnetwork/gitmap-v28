// Package cmdai provides the native AI scripts catalog, discovery, execution, and autofix engine.
package cmdai

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ScriptCategoryType defines taxonomy categories for repository automation scripts.
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

// ScriptMetadata holds descriptive metadata, routing tokens, and capabilities for an AI script.
type ScriptMetadata struct {
	Number      string             `json:"number"`
	Filename    string             `json:"filename"`
	Slug        string             `json:"slug"`
	Aliases     []string           `json:"aliases"`
	Category    ScriptCategoryType `json:"category"`
	Description string             `json:"description"`
	HasFixMode  bool               `json:"hasFixMode"`
	DefaultArgs []string           `json:"defaultArgs"`
}

// PythonRuntime captures discovered Python executable environment metadata.
type PythonRuntime struct {
	ExecutablePath string `json:"executablePath"`
	Version        string `json:"version"`
	MajorVersion   int    `json:"majorVersion"`
	MinorVersion   int    `json:"minorVersion"`
	IsAvailable    bool   `json:"isAvailable"`
	IsPython3      bool   `json:"isPython3"`
}

// FixTargetDef maps a fix target type to its underlying script and configuration.
type FixTargetDef struct {
	Target      FixTargetType `json:"target"`
	ScriptToken string        `json:"scriptToken"`
	DefaultArgs []string      `json:"defaultArgs"`
	Description string        `json:"description"`
}

// ScriptListOptions encapsulates filter and rendering flags for script enumeration.
type ScriptListOptions struct {
	CategoryFilter string `json:"categoryFilter"`
	SearchQuery    string `json:"searchQuery"`
	IsFixOnly      bool   `json:"isFixOnly"`
	IsJsonFormat   bool   `json:"isJsonFormat"`
	IsVerbose      bool   `json:"isVerbose"`
}

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

// CreateScriptOptions configures the generation of a new AI script.
type CreateScriptOptions struct {
	Name        string             `json:"name"`
	Type        ScriptTemplateType `json:"type"`
	Description string             `json:"description"`
	IsParallel  bool               `json:"isParallel"`
	HasFixMode  bool               `json:"hasFixMode"`
	IsDryRun    bool               `json:"isDryRun"`
	IsForce     bool               `json:"isForce"`
}

type (
	// ScriptMetadataResult wraps a single ScriptMetadata outcome with an AppError.
	ScriptMetadataResult = result.Result[ScriptMetadata]

	// ScriptCatalogResult wraps a list of ScriptMetadata outcomes with an AppError.
	ScriptCatalogResult = result.Result[[]ScriptMetadata]

	// PythonRuntimeResult wraps a PythonRuntime outcome with an AppError.
	PythonRuntimeResult = result.Result[PythonRuntime]

	// FixTargetDefResult wraps a FixTargetDef outcome with an AppError.
	FixTargetDefResult = result.Result[FixTargetDef]

	// CreateScriptResult wraps the generated script path outcome with an AppError.
	CreateScriptResult = result.Result[string]
)
