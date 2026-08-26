// Package model — version_manifest.go defines the multi-component SSOT version schema.
package model

// ComponentVersion represents a sub-component version configuration.
type ComponentVersion struct {
	Version     string `json:"version"` // "inherit" or explicit SemVer string (e.g. "1.0.0")
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
}

// RepositoryVersionManifest is the canonical schema for version.json.
type RepositoryVersionManifest struct {
	SchemaDescription string                    `json:"$schema_description"`
	Documentation     string                    `json:"$documentation"`
	InheritanceRules  string                    `json:"$inheritance_rules,omitempty"`
	Version           string                    `json:"version"`
	Backend           *ComponentVersion         `json:"backend,omitempty"`
	Frontend          *ComponentVersion         `json:"frontend,omitempty"`
	CLI               *ComponentVersion         `json:"cli,omitempty"`
	CodingGuidelines  *PromptArchitectMetadata  `json:"coding-guidelines,omitempty"`
	PromptArchitect   *PromptArchitectMetadata  `json:"promptArchitectByRiseupAsia,omitempty"`
}

// ResolveVersion returns the resolved version for a given component, resolving "inherit" to the root version.
func (m *RepositoryVersionManifest) ResolveVersion(component string) string {
	var comp *ComponentVersion
	switch component {
	case "backend":
		comp = m.Backend
	case "frontend":
		comp = m.Frontend
	case "cli":
		comp = m.CLI
	}

	if comp != nil && comp.Version != "" && comp.Version != "inherit" {
		return comp.Version
	}
	return m.Version
}
