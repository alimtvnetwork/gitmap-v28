package jsonenvelope

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// WorkDirectoryConfig represents structured working directory configuration using root variables.
type WorkDirectoryConfig struct {
	Path        string `json:"path"`
	DefaultPath string `json:"defaultPath,omitempty"`
	IsApplied   bool   `json:"isApplied,omitempty"`
	IsEnforced  bool   `json:"isEnforced,omitempty"`
}

// EnvelopeAttributes holds metadata describing payload, origin, commands, and working directory.
type EnvelopeAttributes struct {
	Type                    string               `json:"type"`
	Source                  string               `json:"source,omitempty"`
	How                     string               `json:"how,omitempty"`
	Version                 string               `json:"version,omitempty"`
	GitmapVersion           string               `json:"gitmapVersion,omitempty"`
	ImportCommand           string               `json:"importCommand,omitempty"`
	ExportCommand           string               `json:"exportCommand,omitempty"`
	HelpCommand             string               `json:"helpCommand,omitempty"`
	Notes                   string               `json:"notes,omitempty"`
	Timestamp               string               `json:"timestamp,omitempty"`
	WorkDirectory           string               `json:"workDirectory,omitempty"`
	WorkDirectoryConfig     *WorkDirectoryConfig `json:"workDirectoryConfig,omitempty"`
	DefaultWorkDirectory    string               `json:"defaultWorkDirectory,omitempty"`
	IsWorkDirectoryApplied  bool                 `json:"isWorkDirectoryApplied,omitempty"`
	IsWorkDirectoryEnforced bool                 `json:"isWorkDirectoryEnforced,omitempty"`
}

// UnmarshalJSON supports workDirectory as either a plain string or a WorkDirectoryConfig object.
func (attrs *EnvelopeAttributes) UnmarshalJSON(data []byte) error {
	type Alias EnvelopeAttributes
	var raw struct {
		Alias
		RawWorkDir json.RawMessage `json:"workDirectory"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*attrs = EnvelopeAttributes(raw.Alias)
	return attrs.applyRawWorkDirectory(raw.RawWorkDir)
}

func (attrs *EnvelopeAttributes) applyRawWorkDirectory(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	var strVal string
	if err := json.Unmarshal(raw, &strVal); err == nil {
		attrs.WorkDirectory = strVal
		return nil
	}
	var cfg WorkDirectoryConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return err
	}
	attrs.applyWorkDirConfig(&cfg)
	return nil
}

func (attrs *EnvelopeAttributes) applyWorkDirConfig(cfg *WorkDirectoryConfig) {
	attrs.WorkDirectoryConfig = cfg
	attrs.WorkDirectory = cfg.Path
	if cfg.DefaultPath != "" && attrs.DefaultWorkDirectory == "" {
		attrs.DefaultWorkDirectory = cfg.DefaultPath
	}
	if cfg.IsApplied {
		attrs.IsWorkDirectoryApplied = true
	}
	if cfg.IsEnforced {
		attrs.IsWorkDirectoryEnforced = true
	}
}

// MarshalJSON serializes workDirectory as an object if WorkDirectoryConfig is populated.
func (attrs EnvelopeAttributes) MarshalJSON() ([]byte, error) {
	if attrs.WorkDirectoryConfig != nil {
		type Alias EnvelopeAttributes
		return json.Marshal(&struct {
			Alias
			WorkDirectory *WorkDirectoryConfig `json:"workDirectory"`
		}{
			Alias:         Alias(attrs),
			WorkDirectory: attrs.WorkDirectoryConfig,
		})
	}
	type Alias EnvelopeAttributes
	return json.Marshal(Alias(attrs))
}

// Envelope is a generic envelope wrapping attributes, optional variables, and typed data.
type Envelope[T any] struct {
	Attributes EnvelopeAttributes `json:"attributes"`
	Variables  map[string]any     `json:"variables,omitempty"`
	Data       T                  `json:"data"`
}

// RawEnvelope wraps attributes and optional variables with unparsed json.RawMessage data.
type RawEnvelope struct {
	Attributes EnvelopeAttributes `json:"attributes"`
	Variables  map[string]any     `json:"variables,omitempty"`
	Data       json.RawMessage    `json:"data"`
}

// RelativeBaseName extracts the clean file basename for local directory execution commands.
func RelativeBaseName(source string) string {
	trimmed := strings.TrimSpace(source)
	if trimmed == "" {
		return ""
	}
	norm := strings.ReplaceAll(trimmed, "\\", "/")
	base := filepath.Base(norm)
	if base == "." || base == "/" || base == "\\" {
		return trimmed
	}
	return base
}

// ResolveRelativeJSONPath resolves a relative JSON filename against CWD and common relative subdirectories.
func ResolveRelativeJSONPath(inputPath string) string {
	trimmed := strings.TrimSpace(inputPath)
	if trimmed == "" {
		return inputPath
	}
	if _, err := os.Stat(trimmed); err == nil {
		return trimmed
	}
	base := RelativeBaseName(trimmed)
	searchDirs := []string{
		".",
		"01-gitmap",
		filepath.Join("01-gitmap", "summaries"),
		filepath.Join("01-gitmap", "temp"),
		"vault",
		filepath.Join("02-antigravity-manager", "vault"),
		filepath.Join("02-antigravity-and-event-manager", "vault"),
	}
	for _, dir := range searchDirs {
		candidate := filepath.Join(dir, base)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return inputPath
}

// NewEnvelope constructs a typed envelope with standard v2.0 metadata, workDir/repoDir variables, and payload.
func NewEnvelope[T any](
	dataType string,
	source string,
	how string,
	version string,
	payload T,
) Envelope[T] {
	curTime := time.Now().UTC().Format(time.RFC3339)
	ver := version
	if ver == "" {
		ver = "2.0"
	}
	baseSource := RelativeBaseName(source)
	workDirCfg := &WorkDirectoryConfig{
		Path:        "${workDir}",
		DefaultPath: `D:\work`,
		IsApplied:   true,
		IsEnforced:  false,
	}
	attrs := EnvelopeAttributes{
		Type:                   dataType,
		Source:                 baseSource,
		How:                    how,
		Version:                ver,
		GitmapVersion:          constants.Version,
		Timestamp:              curTime,
		ExportCommand:          how,
		WorkDirectory:          "${workDir}",
		WorkDirectoryConfig:    workDirCfg,
		DefaultWorkDirectory:   `D:\work`,
		IsWorkDirectoryApplied: true,
	}
	populateDescriptorDefaults(&attrs, dataType, baseSource)
	return Envelope[T]{
		Attributes: attrs,
		Variables: map[string]any{
			"workDir":    `D:\work`,
			"repoDir":    `${workDir}\gitmap`,
			"secretsDir": ".",
		},
		Data: payload,
	}
}

func populateDescriptorDefaults(attrs *EnvelopeAttributes, dataType, source string) {
	desc, ok := FindDescriptorByType(dataType)
	if !ok {
		return
	}
	baseSource := RelativeBaseName(source)
	if desc.SuggestedImportCmd != "" && baseSource != "" {
		attrs.ImportCommand = fmt.Sprintf(desc.SuggestedImportCmd, baseSource)
	}
	if baseSource != "" {
		attrs.HelpCommand = fmt.Sprintf("gitmap which-format %s", baseSource)
	}
	attrs.Notes = desc.Description
}

// NewEnvelopeWithAttributes constructs a typed envelope with custom attributes, variables, and payload.
func NewEnvelopeWithAttributes[T any](
	attrs EnvelopeAttributes,
	variables map[string]any,
	payload T,
) Envelope[T] {
	if attrs.Timestamp == "" {
		attrs.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if attrs.GitmapVersion == "" {
		attrs.GitmapVersion = constants.Version
	}
	if attrs.Version == "" {
		attrs.Version = "2.0"
	}
	return Envelope[T]{
		Attributes: attrs,
		Variables:  variables,
		Data:       payload,
	}
}

// IsEnvelope returns true if the raw JSON conforms to the {attributes, data} envelope schema.
func IsEnvelope(raw []byte) bool {
	var probe struct {
		Attributes *struct {
			Type string `json:"type"`
		} `json:"attributes"`
		Data json.RawMessage `json:"data"`
	}
	err := json.Unmarshal(raw, &probe)
	if err != nil {
		return false
	}
	return probe.Attributes != nil && len(probe.Data) > 0
}

// ExtractEnvelope parses the complete raw envelope including attributes and variables.
func ExtractEnvelope(raw []byte) (RawEnvelope, error) {
	var env RawEnvelope
	err := json.Unmarshal(raw, &env)
	if err != nil {
		return RawEnvelope{}, apperror.WrapSimple(err, "ExtractEnvelope.Unmarshal")
	}
	return env, nil
}

// ExtractPayload extracts the core payload bytes and attributes from either an envelope or legacy JSON.
func ExtractPayload(raw []byte) ([]byte, EnvelopeAttributes, error) {
	clean := strings.TrimSpace(string(raw))
	if len(clean) == 0 {
		return nil, EnvelopeAttributes{}, apperror.NewValidationError("ExtractPayload: empty input JSON")
	}

	if !IsEnvelope(raw) {
		return extractLegacyOrFlatPayload(raw)
	}

	return extractEnvelopePayload(raw)
}

func extractLegacyOrFlatPayload(raw []byte) ([]byte, EnvelopeAttributes, error) {
	env, err := ExtractEnvelope(raw)
	if err == nil && hasEnvelopeVariablesOrWorkDir(env) {
		allVars := MergeVariables(env.Variables, env.Attributes.WorkDirectoryConfig)
		return ExpandVariables(raw, allVars), env.Attributes, nil
	}
	attrs := EnvelopeAttributes{
		Version: "legacy",
	}
	return raw, attrs, nil
}

func hasEnvelopeVariablesOrWorkDir(env RawEnvelope) bool {
	return len(env.Variables) > 0 || env.Attributes.WorkDirectoryConfig != nil
}

func extractEnvelopePayload(raw []byte) ([]byte, EnvelopeAttributes, error) {
	env, err := ExtractEnvelope(raw)
	if err != nil {
		return nil, EnvelopeAttributes{}, err
	}
	allVars := MergeVariables(env.Variables, env.Attributes.WorkDirectoryConfig)
	data := env.Data
	if len(allVars) > 0 {
		data = ExpandVariables(data, allVars)
		env.Attributes.WorkDirectory = ResolveStringVariable(env.Attributes.WorkDirectory, allVars)
		if env.Attributes.WorkDirectoryConfig != nil {
			env.Attributes.WorkDirectoryConfig.Path = ResolveStringVariable(env.Attributes.WorkDirectoryConfig.Path, allVars)
		}
	}
	return data, env.Attributes, nil
}
