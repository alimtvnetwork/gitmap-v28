// Package cmdai provides the native AI scripts catalog, discovery, execution, and autofix engine.
package cmdai

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ScriptCategoryType is the canonical script category type, defined in enums.
type ScriptCategoryType = enums.ScriptCategoryType

const (
	CategoryCore        = enums.CategoryCore
	CategoryFormatting  = enums.CategoryFormatting
	CategoryGuideline   = enums.CategoryGuideline
	CategoryCicd        = enums.CategoryCicd
	CategoryPaths       = enums.CategoryPaths
	CategoryAudit       = enums.CategoryAudit
	CategoryDiscovery   = enums.CategoryDiscovery
	CategoryMaintenance = enums.CategoryMaintenance
	CategoryPlans       = enums.CategoryPlans
	CategoryRelease     = enums.CategoryRelease
	CategoryDatabase    = enums.CategoryDatabase
	CategoryUnknown     = enums.CategoryUnknown
)

// FixTargetType is the canonical fix target type, defined in enums.
type FixTargetType = enums.FixTargetType

const (
	FixTargetGuidelines = enums.FixTargetGuidelines
	FixTargetNewlines   = enums.FixTargetNewlines
	FixTargetPaths      = enums.FixTargetPaths
	FixTargetNaming     = enums.FixTargetNaming
	FixTargetEncoding   = enums.FixTargetEncoding
	FixTargetSpelling   = enums.FixTargetSpelling
	FixTargetPlans      = enums.FixTargetPlans
	FixTargetAll        = enums.FixTargetAll
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
	CategoryFilter  string `json:"categoryFilter"`
	SearchQuery     string `json:"searchQuery"`
	IsFixOnly       bool   `json:"isFixOnly"`
	IsJsonFormat    bool   `json:"isJsonFormat"`
	IsVerbose       bool   `json:"isVerbose"`
	IsCopyClipboard bool   `json:"isCopyClipboard"`
	IsFrequentOnly  bool   `json:"isFrequentOnly"`
}

// ScriptTemplateType is the canonical script template type, defined in enums.
type ScriptTemplateType = enums.ScriptTemplateType

const (
	TemplateLinter    = enums.TemplateLinter
	TemplateFixer     = enums.TemplateFixer
	TemplateAuditor   = enums.TemplateAuditor
	TemplateChecker   = enums.TemplateChecker
	TemplateGenerator = enums.TemplateGenerator
	TemplateUtil      = enums.TemplateUtil
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
	ScriptMetadataResult = result.Wrap[ScriptMetadata]

	// ScriptCatalogResult wraps a list of ScriptMetadata outcomes with an AppError.
	ScriptCatalogResult = result.WrapSlice[ScriptMetadata]

	// PythonRuntimeResult wraps a PythonRuntime outcome with an AppError.
	PythonRuntimeResult = result.Wrap[PythonRuntime]

	// FixTargetDefResult wraps a FixTargetDef outcome with an AppError.
	FixTargetDefResult = result.Wrap[FixTargetDef]

	// CreateScriptResult wraps the generated script path outcome with an AppError.
	CreateScriptResult = result.Wrap[string]
)
