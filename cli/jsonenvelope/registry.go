package jsonenvelope

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Known envelope type constants.
const (
	TypeSSHNodes         = "ssh-nodes"
	TypeMacro            = "macro"
	TypeMacroBundle      = "macro-bundle"
	TypeCommitPullConfig = "commit-pull-config"
	TypeUISettings       = "ui-settings"
	TypeTestInventory    = "test-inventory"
	TypePipelineConfig   = "pipeline-config"
	TypeGitProfiles      = "git-profiles"
	TypeTemplates        = "templates"
	TypeRepoManifest     = "repo-manifest"
)

// TypeDescriptor defines metadata, command recommendations, and detection logic for a data type.
type TypeDescriptor struct {
	Type                  string
	Name                  string
	Description           string
	SuggestedImportCmd    string
	SystemImpact          string
	LegacySignatureFilter func(raw []byte) bool
}

// Registry maps type identifiers to descriptors.
var registeredTypes = []TypeDescriptor{
	{
		Type:                  TypeSSHNodes,
		Name:                  "SSH Cluster Nodes",
		Description:           "Definitions and credentials for remote SSH cluster nodes.",
		SuggestedImportCmd:    "gitmap sj import %s -y",
		SystemImpact:          "Enrolls remote SSH machines into installation.db and credential vault.",
		LegacySignatureFilter: hasLegacySSHNodesSignature,
	},
	{
		Type:                  TypeMacro,
		Name:                  "Macro Automation Bundle",
		Description:           "Recorded automation macros, task pipelines, and step trees.",
		SuggestedImportCmd:    "gitmap macro import %s -y",
		SystemImpact:          "Registers macro definitions into SQLite repository macros table.",
		LegacySignatureFilter: hasLegacyMacroSignature,
	},
	{
		Type:                  TypeCommitPullConfig,
		Name:                  "Commit-In & Pull Configuration",
		Description:           "Declarative commit-in and pull synchronization preferences.",
		SuggestedImportCmd:    "gitmap commitin --config %s",
		SystemImpact:          "Configures commit-in execution flags, branch rules, and pull sync behaviors.",
		LegacySignatureFilter: hasLegacyCommitInSignature,
	},
	{
		Type:                  TypeUISettings,
		Name:                  "UI Layout & Graphics Preferences",
		Description:           "Terminal and Web UI graphics mode, layout, and browser auto-open preferences.",
		SuggestedImportCmd:    "gitmap cmdui import-settings %s",
		SystemImpact:          "Updates ~/.gitmap/ui_settings.json layout orientation and graphics settings.",
		LegacySignatureFilter: hasLegacyUISettingsSignature,
	},
	{
		Type:                  TypeTestInventory,
		Name:                  "Test Inventory Manifest",
		Description:           "Test execution timing records, historical pass rates, and inventory trees.",
		SuggestedImportCmd:    "gitmap test-inventory import %s",
		SystemImpact:          "Updates test execution historical inventory for smart test selection.",
		LegacySignatureFilter: hasLegacyTestInventorySignature,
	},
	{
		Type:                  TypePipelineConfig,
		Name:                  "Pipeline AI Configuration",
		Description:           "CI/CD pipeline monitoring rules, timeout thresholds, and runner configs.",
		SuggestedImportCmd:    "gitmap pipeline import %s",
		SystemImpact:          "Updates pipeline AI monitoring rules and dynamic wait calculation settings.",
		LegacySignatureFilter: hasLegacyPipelineConfigSignature,
	},
	{
		Type:                  TypeGitProfiles,
		Name:                  "Git Author Profiles",
		Description:           "Git author identities, GitHub accounts, and default commit signatures.",
		SuggestedImportCmd:    "gitmap backup restore --file %s",
		SystemImpact:          "Updates active Git user profile and credentials in git_profiles.json.",
		LegacySignatureFilter: hasLegacyGitProfilesSignature,
	},
	{
		Type:                  TypeTemplates,
		Name:                  "State Templates & Variables",
		Description:           "State template snippets, category sequences, and SEO/developer variables.",
		SuggestedImportCmd:    "gitmap templates import %s",
		SystemImpact:          "Enrolls template items and variables into templates.db split-database.",
		LegacySignatureFilter: hasLegacyTemplatesSignature,
	},
	{
		Type:                  TypeRepoManifest,
		Name:                  "Repository Scan Manifest",
		Description:           "Discovered Git repositories, remote URLs, branches, and clone metadata.",
		SuggestedImportCmd:    "gitmap clone %s",
		SystemImpact:          "Clones or synchronizes repositories registered in the scan manifest.",
		LegacySignatureFilter: hasLegacyRepoManifestSignature,
	},
}

// GetRegisteredTypes returns a copy of all registered format descriptors.
func GetRegisteredTypes() []TypeDescriptor {
	result := make([]TypeDescriptor, len(registeredTypes))
	copy(result, registeredTypes)
	return result
}

// FindDescriptorByType looks up a descriptor by its type identifier.
func FindDescriptorByType(typeName string) (TypeDescriptor, bool) {
	norm := strings.ToLower(strings.TrimSpace(typeName))
	if norm == TypeMacroBundle {
		norm = TypeMacro
	}
	for _, d := range registeredTypes {
		if strings.ToLower(d.Type) == norm {
			return d, true
		}
	}
	return TypeDescriptor{}, false
}

// DetectFormat inspects the raw JSON bytes and returns the matching descriptor and attributes.
func DetectFormat(raw []byte) (TypeDescriptor, EnvelopeAttributes, bool) {
	if IsEnvelope(raw) {
		return detectFromEnvelope(raw)
	}
	return detectFromLegacy(raw)
}

func detectFromEnvelope(raw []byte) (TypeDescriptor, EnvelopeAttributes, bool) {
	var env RawEnvelope
	err := json.Unmarshal(raw, &env)
	if err != nil {
		return TypeDescriptor{}, EnvelopeAttributes{}, false
	}

	desc, ok := FindDescriptorByType(env.Attributes.Type)
	if !ok {
		// Custom/unknown envelope type
		customDesc := TypeDescriptor{
			Type:               env.Attributes.Type,
			Name:               fmt.Sprintf("Custom Type (%s)", env.Attributes.Type),
			Description:        "User-defined typed envelope JSON format.",
			SuggestedImportCmd: fmt.Sprintf("gitmap import --type %s %%s -y", env.Attributes.Type),
			SystemImpact:       fmt.Sprintf("Imports data payload for custom type '%s'.", env.Attributes.Type),
		}
		return customDesc, env.Attributes, true
	}
	return desc, env.Attributes, true
}

func detectFromLegacy(raw []byte) (TypeDescriptor, EnvelopeAttributes, bool) {
	for _, d := range registeredTypes {
		if d.LegacySignatureFilter != nil && d.LegacySignatureFilter(raw) {
			inferredAttrs := EnvelopeAttributes{
				Type:    d.Type,
				Source:  "legacy-unwrapped",
				Version: "legacy",
			}
			return d, inferredAttrs, true
		}
	}
	return TypeDescriptor{}, EnvelopeAttributes{}, false
}

func hasLegacySSHNodesSignature(raw []byte) bool {
	str := string(raw)
	hasNodesArray := strings.Contains(str, `"nodes"`) && (strings.Contains(str, `"ip_address"`) || strings.Contains(str, `"ipAddress"`))
	if hasNodesArray {
		return true
	}
	hasDirectWorkers := (strings.Contains(str, `"worker_id"`) || strings.Contains(str, `"workerId"`)) &&
		(strings.Contains(str, `"ip_address"`) || strings.Contains(str, `"ipAddress"`))
	return hasDirectWorkers
}

func hasLegacyMacroSignature(raw []byte) bool {
	str := string(raw)
	hasMacroFields := strings.Contains(str, `"steps"`) && strings.Contains(str, `"action"`)
	if hasMacroFields {
		return true
	}
	hasCommandSteps := strings.Contains(str, `"command"`) && strings.Contains(str, `"is_background"`)
	return hasCommandSteps
}

func hasLegacyCommitInSignature(raw []byte) bool {
	str := string(raw)
	hasTreeOrSync := strings.Contains(str, `"isApplyTree"`) || strings.Contains(str, `"isApplyFinalSync"`)
	if hasTreeOrSync {
		return true
	}
	hasLegacyFlags := strings.Contains(str, `"isPushImmediate"`) || strings.Contains(str, `"pullDirection"`)
	return hasLegacyFlags
}

func hasLegacyUISettingsSignature(raw []byte) bool {
	str := string(raw)
	hasGraphics := strings.Contains(str, `"graphicsMode"`) || strings.Contains(str, `"autoOpenBrowser"`)
	if hasGraphics {
		return true
	}
	hasLayout := strings.Contains(str, `"commitInLayout"`) || strings.Contains(str, `"prReplayMode"`)
	return hasLayout
}

func hasLegacyTestInventorySignature(raw []byte) bool {
	str := string(raw)
	return strings.Contains(str, `"test_inventory"`) || (strings.Contains(str, `"total_tests"`) && strings.Contains(str, `"packages"`))
}

func hasLegacyPipelineConfigSignature(raw []byte) bool {
	str := string(raw)
	return strings.Contains(str, `"pipeline_ai"`) || (strings.Contains(str, `"eta_threshold"`) && strings.Contains(str, `"workflows"`))
}

func hasLegacyGitProfilesSignature(raw []byte) bool {
	str := string(raw)
	return strings.Contains(str, `"profiles"`) && (strings.Contains(str, `"authMethod"`) || strings.Contains(str, `"provider"`))
}

func hasLegacyTemplatesSignature(raw []byte) bool {
	str := string(raw)
	return strings.Contains(str, `"exportId"`) && (strings.Contains(str, `"templates"`) || strings.Contains(str, `"variables"`))
}

func hasLegacyRepoManifestSignature(raw []byte) bool {
	str := string(raw)
	hasRepoKeys := strings.Contains(str, `"cloneInstruction"`) || strings.Contains(str, `"repoId"`)
	if hasRepoKeys {
		return true
	}
	return strings.Contains(str, `"discoveredUrl"`) || (strings.Contains(str, `"sshUrl"`) && strings.Contains(str, `"relativePath"`))
}
