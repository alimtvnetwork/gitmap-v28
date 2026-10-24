package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ToolchainEnvVarType represents environment variable names used by toolchains.
type ToolchainEnvVarType string

const (
	ToolchainEnvLibclangPath ToolchainEnvVarType = "LIBCLANG_PATH"
	ToolchainEnvPath         ToolchainEnvVarType = "PATH"
	ToolchainEnvCargoHome    ToolchainEnvVarType = "CARGO_HOME"
	ToolchainEnvRustupHome   ToolchainEnvVarType = "RUSTUP_HOME"
	ToolchainEnvSystemRoot   ToolchainEnvVarType = "SYSTEMROOT"
	ToolchainEnvComSpec      ToolchainEnvVarType = "COMSPEC"
	ToolchainEnvTmp          ToolchainEnvVarType = "TMP"
	ToolchainEnvTemp         ToolchainEnvVarType = "TEMP"
	ToolchainEnvUserProfile  ToolchainEnvVarType = "USERPROFILE"
)

// IsValid reports whether the environment variable name matches a known variant.
func (t ToolchainEnvVarType) IsValid() bool {
	switch t {
	case ToolchainEnvLibclangPath, ToolchainEnvPath, ToolchainEnvCargoHome,
		ToolchainEnvRustupHome, ToolchainEnvSystemRoot, ToolchainEnvComSpec,
		ToolchainEnvTmp, ToolchainEnvTemp, ToolchainEnvUserProfile:
		return true
	default:
		return false
	}
}

// String returns the string representation of the environment variable.
func (t ToolchainEnvVarType) String() string {
	return string(t)
}

// ParseToolchainEnvVarType parses a raw string into a validated ToolchainEnvVarType.
func ParseToolchainEnvVarType(s string) (ToolchainEnvVarType, *apperror.AppError) {
	val := ToolchainEnvVarType(strings.ToUpper(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseToolchainEnvVarType", "E1000", fmt.Sprintf("invalid toolchain env var: %q", s)).WithContext("value", s)
}

// WithValue returns the "NAME=value" assignment form of the variable,
// without string concatenation at the call site.
func (t ToolchainEnvVarType) WithValue(value string) string {
	return strings.Join([]string{t.String(), value}, "=")
}

// AssignmentPrefix returns the "NAME=" prefix used to match an existing
// entry in an environment list.
func (t ToolchainEnvVarType) AssignmentPrefix() string {
	return t.WithValue("")
}

// AllToolchainEnvVarTypes returns all valid ToolchainEnvVarType variants.
func AllToolchainEnvVarTypes() []ToolchainEnvVarType {
	return []ToolchainEnvVarType{
		ToolchainEnvLibclangPath,
		ToolchainEnvPath,
		ToolchainEnvCargoHome,
		ToolchainEnvRustupHome,
		ToolchainEnvSystemRoot,
		ToolchainEnvComSpec,
		ToolchainEnvTmp,
		ToolchainEnvTemp,
		ToolchainEnvUserProfile,
	}
}

// IsLibclangPath reports whether the variable is LIBCLANG_PATH.
func (t ToolchainEnvVarType) IsLibclangPath() bool {
	return t == ToolchainEnvLibclangPath
}

// IsPath reports whether the variable is PATH.
func (t ToolchainEnvVarType) IsPath() bool {
	return t == ToolchainEnvPath
}

// IsCargoHome reports whether the variable is CARGO_HOME.
func (t ToolchainEnvVarType) IsCargoHome() bool {
	return t == ToolchainEnvCargoHome
}

// IsRustupHome reports whether the variable is RUSTUP_HOME.
func (t ToolchainEnvVarType) IsRustupHome() bool {
	return t == ToolchainEnvRustupHome
}

// IsSystemRoot reports whether the variable is SYSTEMROOT.
func (t ToolchainEnvVarType) IsSystemRoot() bool {
	return t == ToolchainEnvSystemRoot
}

// IsComSpec reports whether the variable is COMSPEC.
func (t ToolchainEnvVarType) IsComSpec() bool {
	return t == ToolchainEnvComSpec
}

// IsTmp reports whether the variable is TMP.
func (t ToolchainEnvVarType) IsTmp() bool {
	return t == ToolchainEnvTmp
}

// IsTemp reports whether the variable is TEMP.
func (t ToolchainEnvVarType) IsTemp() bool {
	return t == ToolchainEnvTemp
}

// IsUserProfile reports whether the variable is USERPROFILE.
func (t ToolchainEnvVarType) IsUserProfile() bool {
	return t == ToolchainEnvUserProfile
}

// HasEnvVar reports whether the variable is non-empty and valid.
func (t ToolchainEnvVarType) HasEnvVar() bool {
	return t.IsValid()
}
