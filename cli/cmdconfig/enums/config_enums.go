package enums

import (
	"fmt"
	"strings"
)

// ConfigActionType represents supported configuration CLI actions.
type ConfigActionType string

const (
	ConfigActionExport      ConfigActionType = "export"
	ConfigActionImport      ConfigActionType = "import"
	ConfigActionWhatConfigs ConfigActionType = "what-configs"
	ConfigActionWhichFormat ConfigActionType = "which-format"
	ConfigActionGet         ConfigActionType = "get"
	ConfigActionSet         ConfigActionType = "set"
	ConfigActionList        ConfigActionType = "list"
	ConfigActionReset       ConfigActionType = "reset"
)

// IsValid reports whether the action is a supported configuration action.
func (a ConfigActionType) IsValid() bool {
	switch a {
	case ConfigActionExport, ConfigActionImport, ConfigActionWhatConfigs,
		ConfigActionWhichFormat, ConfigActionGet, ConfigActionSet,
		ConfigActionList, ConfigActionReset:
		return true
	default:
		return false
	}
}

// String returns the string representation of the action.
func (a ConfigActionType) String() string {
	return string(a)
}

// IsExport reports whether the action is export.
func (a ConfigActionType) IsExport() bool {
	return a == ConfigActionExport
}

// IsImport reports whether the action is import.
func (a ConfigActionType) IsImport() bool {
	return a == ConfigActionImport
}

// IsWhatConfigs reports whether the action is what-configs.
func (a ConfigActionType) IsWhatConfigs() bool {
	return a == ConfigActionWhatConfigs
}

// IsWhichFormat reports whether the action is which-format.
func (a ConfigActionType) IsWhichFormat() bool {
	return a == ConfigActionWhichFormat
}

// IsGet reports whether the action is get.
func (a ConfigActionType) IsGet() bool {
	return a == ConfigActionGet
}

// IsSet reports whether the action is set.
func (a ConfigActionType) IsSet() bool {
	return a == ConfigActionSet
}

// IsList reports whether the action is list.
func (a ConfigActionType) IsList() bool {
	return a == ConfigActionList
}

// IsReset reports whether the action is reset.
func (a ConfigActionType) IsReset() bool {
	return a == ConfigActionReset
}

// ParseConfigActionType parses a raw string into a validated ConfigActionType.
func ParseConfigActionType(s string) (ConfigActionType, error) {
	val := ConfigActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid config action: %q", s)
}

// AllConfigActionTypes returns all valid ConfigActionType variants.
func AllConfigActionTypes() []ConfigActionType {
	return []ConfigActionType{
		ConfigActionExport,
		ConfigActionImport,
		ConfigActionWhatConfigs,
		ConfigActionWhichFormat,
		ConfigActionGet,
		ConfigActionSet,
		ConfigActionList,
		ConfigActionReset,
	}
}

// ConfigOutputFormatType represents presentation output formats for config listings.
type ConfigOutputFormatType string

const (
	ConfigOutputFormatTable   ConfigOutputFormatType = "table"
	ConfigOutputFormatJSON    ConfigOutputFormatType = "json"
	ConfigOutputFormatYAML    ConfigOutputFormatType = "yaml"
	ConfigOutputFormatCompact ConfigOutputFormatType = "compact"
)

// IsValid reports whether the format matches a recognized output format.
func (o ConfigOutputFormatType) IsValid() bool {
	switch o {
	case ConfigOutputFormatTable, ConfigOutputFormatJSON,
		ConfigOutputFormatYAML, ConfigOutputFormatCompact:
		return true
	default:
		return false
	}
}

// String returns the string representation of the output format.
func (o ConfigOutputFormatType) String() string {
	return string(o)
}

// IsTable reports whether the format is table.
func (o ConfigOutputFormatType) IsTable() bool {
	return o == ConfigOutputFormatTable
}

// IsJSON reports whether the format is json.
func (o ConfigOutputFormatType) IsJSON() bool {
	return o == ConfigOutputFormatJSON
}

// IsYAML reports whether the format is yaml.
func (o ConfigOutputFormatType) IsYAML() bool {
	return o == ConfigOutputFormatYAML
}

// IsCompact reports whether the format is compact.
func (o ConfigOutputFormatType) IsCompact() bool {
	return o == ConfigOutputFormatCompact
}

// ParseConfigOutputFormatType parses a raw string into a validated ConfigOutputFormatType.
func ParseConfigOutputFormatType(s string) (ConfigOutputFormatType, error) {
	val := ConfigOutputFormatType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid config output format: %q", s)
}

// AllConfigOutputFormatTypes returns all valid ConfigOutputFormatType variants.
func AllConfigOutputFormatTypes() []ConfigOutputFormatType {
	return []ConfigOutputFormatType{
		ConfigOutputFormatTable,
		ConfigOutputFormatJSON,
		ConfigOutputFormatYAML,
		ConfigOutputFormatCompact,
	}
}

// ConfigScopeType represents configuration inheritance scopes.
type ConfigScopeType string

const (
	ConfigScopeLocal   ConfigScopeType = "local"
	ConfigScopeGlobal  ConfigScopeType = "global"
	ConfigScopeSystem  ConfigScopeType = "system"
	ConfigScopeProject ConfigScopeType = "project"
)

// IsValid reports whether the scope matches a known configuration scope.
func (s ConfigScopeType) IsValid() bool {
	switch s {
	case ConfigScopeLocal, ConfigScopeGlobal, ConfigScopeSystem, ConfigScopeProject:
		return true
	default:
		return false
	}
}

// String returns the string representation of the scope.
func (s ConfigScopeType) String() string {
	return string(s)
}

// IsLocal reports whether the scope is local.
func (s ConfigScopeType) IsLocal() bool {
	return s == ConfigScopeLocal
}

// IsGlobal reports whether the scope is global.
func (s ConfigScopeType) IsGlobal() bool {
	return s == ConfigScopeGlobal
}

// IsSystem reports whether the scope is system.
func (s ConfigScopeType) IsSystem() bool {
	return s == ConfigScopeSystem
}

// IsProject reports whether the scope is project.
func (s ConfigScopeType) IsProject() bool {
	return s == ConfigScopeProject
}

// ParseConfigScopeType parses a raw string into a validated ConfigScopeType.
func ParseConfigScopeType(s string) (ConfigScopeType, error) {
	val := ConfigScopeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid config scope: %q", s)
}

// AllConfigScopeTypes returns all valid ConfigScopeType variants.
func AllConfigScopeTypes() []ConfigScopeType {
	return []ConfigScopeType{
		ConfigScopeLocal,
		ConfigScopeGlobal,
		ConfigScopeSystem,
		ConfigScopeProject,
	}
}

// ConfigKeyType represents recognized configuration key identifiers.
type ConfigKeyType string

const (
	ConfigKeyTargetDir     ConfigKeyType = "target_dir"
	ConfigKeyDefaultTool   ConfigKeyType = "default_tool"
	ConfigKeyExportDir     ConfigKeyType = "export_dir"
	ConfigKeySchemaVersion ConfigKeyType = "schema_version"
	ConfigKeyEncoding      ConfigKeyType = "encoding"
	ConfigKeyAutoImport    ConfigKeyType = "auto_import"
)

// IsValid reports whether the key matches a recognized configuration key.
func (k ConfigKeyType) IsValid() bool {
	switch k {
	case ConfigKeyTargetDir, ConfigKeyDefaultTool, ConfigKeyExportDir,
		ConfigKeySchemaVersion, ConfigKeyEncoding, ConfigKeyAutoImport:
		return true
	default:
		return false
	}
}

// String returns the string representation of the key.
func (k ConfigKeyType) String() string {
	return string(k)
}

// IsTargetDir reports whether the key is target_dir.
func (k ConfigKeyType) IsTargetDir() bool {
	return k == ConfigKeyTargetDir
}

// IsDefaultTool reports whether the key is default_tool.
func (k ConfigKeyType) IsDefaultTool() bool {
	return k == ConfigKeyDefaultTool
}

// IsExportDir reports whether the key is export_dir.
func (k ConfigKeyType) IsExportDir() bool {
	return k == ConfigKeyExportDir
}

// IsSchemaVersion reports whether the key is schema_version.
func (k ConfigKeyType) IsSchemaVersion() bool {
	return k == ConfigKeySchemaVersion
}

// IsEncoding reports whether the key is encoding.
func (k ConfigKeyType) IsEncoding() bool {
	return k == ConfigKeyEncoding
}

// IsAutoImport reports whether the key is auto_import.
func (k ConfigKeyType) IsAutoImport() bool {
	return k == ConfigKeyAutoImport
}

// ParseConfigKeyType parses a raw string into a validated ConfigKeyType.
func ParseConfigKeyType(s string) (ConfigKeyType, error) {
	val := ConfigKeyType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid config key: %q", s)
}

// AllConfigKeyTypes returns all valid ConfigKeyType variants.
func AllConfigKeyTypes() []ConfigKeyType {
	return []ConfigKeyType{
		ConfigKeyTargetDir,
		ConfigKeyDefaultTool,
		ConfigKeyExportDir,
		ConfigKeySchemaVersion,
		ConfigKeyEncoding,
		ConfigKeyAutoImport,
	}
}

// ConfigToolType represents tools supported for configuration export and import.
type ConfigToolType string

const (
	ConfigToolVSCode   ConfigToolType = "vscode"
	ConfigToolQTorrent ConfigToolType = "qtorrent"
	ConfigToolUTorrent ConfigToolType = "utorrent"
	ConfigToolAll      ConfigToolType = "all"
)

// IsValid reports whether the tool matches a supported configuration tool.
func (t ConfigToolType) IsValid() bool {
	switch t {
	case ConfigToolVSCode, ConfigToolQTorrent, ConfigToolUTorrent, ConfigToolAll:
		return true
	default:
		return false
	}
}

// String returns the string representation of the tool.
func (t ConfigToolType) String() string {
	return string(t)
}

// IsVSCode reports whether the tool is vscode.
func (t ConfigToolType) IsVSCode() bool {
	return t == ConfigToolVSCode
}

// IsQTorrent reports whether the tool is qtorrent.
func (t ConfigToolType) IsQTorrent() bool {
	return t == ConfigToolQTorrent
}

// IsUTorrent reports whether the tool is utorrent.
func (t ConfigToolType) IsUTorrent() bool {
	return t == ConfigToolUTorrent
}

// IsAll reports whether the tool selects all tools.
func (t ConfigToolType) IsAll() bool {
	return t == ConfigToolAll
}

// ParseConfigToolType parses a raw string or alias into a canonical ConfigToolType.
func ParseConfigToolType(s string) (ConfigToolType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "vscode", "code", "vcode":
		return ConfigToolVSCode, nil
	case "qtorrent", "qbittorrent", "qbit":
		return ConfigToolQTorrent, nil
	case "utorrent", "uttorrent", "u-torrent":
		return ConfigToolUTorrent, nil
	case "all", "*":
		return ConfigToolAll, nil
	default:
		return "", fmt.Errorf("invalid config tool: %q", s)
	}
}

// AllConfigToolTypes returns all valid ConfigToolType variants.
func AllConfigToolTypes() []ConfigToolType {
	return []ConfigToolType{
		ConfigToolVSCode,
		ConfigToolQTorrent,
		ConfigToolUTorrent,
		ConfigToolAll,
	}
}

// ConfigFileFormatType represents file formats detected by configuration inspection.
type ConfigFileFormatType string

const (
	ConfigFileFormatJSON   ConfigFileFormatType = "json"
	ConfigFileFormatYAML   ConfigFileFormatType = "yaml"
	ConfigFileFormatDB     ConfigFileFormatType = "db"
	ConfigFileFormatSQLite ConfigFileFormatType = "sqlite"
	ConfigFileFormatTOML   ConfigFileFormatType = "toml"
	ConfigFileFormatENV    ConfigFileFormatType = "env"
	ConfigFileFormatINI    ConfigFileFormatType = "ini"
	ConfigFileFormatCFG    ConfigFileFormatType = "cfg"
	ConfigFileFormatCONF   ConfigFileFormatType = "conf"
)

// IsValid reports whether the format matches a recognized configuration file format.
func (f ConfigFileFormatType) IsValid() bool {
	switch f {
	case ConfigFileFormatJSON, ConfigFileFormatYAML, ConfigFileFormatDB,
		ConfigFileFormatSQLite, ConfigFileFormatTOML, ConfigFileFormatENV,
		ConfigFileFormatINI, ConfigFileFormatCFG, ConfigFileFormatCONF:
		return true
	default:
		return false
	}
}

// String returns the string representation of the file format.
func (f ConfigFileFormatType) String() string {
	return string(f)
}

// IsJSON reports whether the file format is JSON.
func (f ConfigFileFormatType) IsJSON() bool {
	return f == ConfigFileFormatJSON
}

// IsYAML reports whether the file format is YAML.
func (f ConfigFileFormatType) IsYAML() bool {
	return f == ConfigFileFormatYAML
}

// IsDB reports whether the file format is DB.
func (f ConfigFileFormatType) IsDB() bool {
	return f == ConfigFileFormatDB
}

// IsSQLite reports whether the file format is SQLite.
func (f ConfigFileFormatType) IsSQLite() bool {
	return f == ConfigFileFormatSQLite
}

// IsTOML reports whether the file format is TOML.
func (f ConfigFileFormatType) IsTOML() bool {
	return f == ConfigFileFormatTOML
}

// IsENV reports whether the file format is ENV.
func (f ConfigFileFormatType) IsENV() bool {
	return f == ConfigFileFormatENV
}

// IsINI reports whether the file format is INI.
func (f ConfigFileFormatType) IsINI() bool {
	return f == ConfigFileFormatINI
}

// IsCFG reports whether the file format is CFG.
func (f ConfigFileFormatType) IsCFG() bool {
	return f == ConfigFileFormatCFG
}

// IsCONF reports whether the file format is CONF.
func (f ConfigFileFormatType) IsCONF() bool {
	return f == ConfigFileFormatCONF
}

var configFileFormatMap = map[string]ConfigFileFormatType{
	"json":    ConfigFileFormatJSON,
	"yaml":    ConfigFileFormatYAML,
	"yml":     ConfigFileFormatYAML,
	"db":      ConfigFileFormatDB,
	"sqlite":  ConfigFileFormatSQLite,
	"sqlite3": ConfigFileFormatSQLite,
	"toml":    ConfigFileFormatTOML,
	"env":     ConfigFileFormatENV,
	"ini":     ConfigFileFormatINI,
	"cfg":     ConfigFileFormatCFG,
	"conf":    ConfigFileFormatCONF,
}

// ParseConfigFileFormatType parses a raw string or extension into a ConfigFileFormatType.
func ParseConfigFileFormatType(s string) (ConfigFileFormatType, error) {
	lower := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(s)), ".")
	if val, ok := configFileFormatMap[lower]; ok {
		return val, nil
	}

	return "", fmt.Errorf("invalid config file format: %q", s)
}

// AllConfigFileFormatTypes returns all valid ConfigFileFormatType variants.
func AllConfigFileFormatTypes() []ConfigFileFormatType {
	return []ConfigFileFormatType{
		ConfigFileFormatJSON,
		ConfigFileFormatYAML,
		ConfigFileFormatDB,
		ConfigFileFormatSQLite,
		ConfigFileFormatTOML,
		ConfigFileFormatENV,
		ConfigFileFormatINI,
		ConfigFileFormatCFG,
		ConfigFileFormatCONF,
	}
}

// ConfigFileEncodingType represents supported content encoding types for payloads.
type ConfigFileEncodingType string

const (
	ConfigFileEncodingUTF8   ConfigFileEncodingType = "utf-8"
	ConfigFileEncodingBase64 ConfigFileEncodingType = "base64"
)

// IsValid reports whether the encoding matches a recognized encoding type.
func (e ConfigFileEncodingType) IsValid() bool {
	switch e {
	case ConfigFileEncodingUTF8, ConfigFileEncodingBase64:
		return true
	default:
		return false
	}
}

// String returns the string representation of the encoding type.
func (e ConfigFileEncodingType) String() string {
	return string(e)
}

// IsUTF8 reports whether the encoding is utf-8.
func (e ConfigFileEncodingType) IsUTF8() bool {
	return e == ConfigFileEncodingUTF8
}

// IsBase64 reports whether the encoding is base64.
func (e ConfigFileEncodingType) IsBase64() bool {
	return e == ConfigFileEncodingBase64
}

// ParseConfigFileEncodingType parses a raw string into a validated ConfigFileEncodingType.
func ParseConfigFileEncodingType(s string) (ConfigFileEncodingType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "utf-8", "utf8":
		return ConfigFileEncodingUTF8, nil
	case "base64", "b64":
		return ConfigFileEncodingBase64, nil
	default:
		return "", fmt.Errorf("invalid config file encoding: %q", s)
	}
}

// AllConfigFileEncodingTypes returns all valid ConfigFileEncodingType variants.
func AllConfigFileEncodingTypes() []ConfigFileEncodingType {
	return []ConfigFileEncodingType{
		ConfigFileEncodingUTF8,
		ConfigFileEncodingBase64,
	}
}
