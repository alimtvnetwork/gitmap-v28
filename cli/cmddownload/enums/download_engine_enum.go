package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// DownloadEngineType represents the active downloader backend.
type DownloadEngineType string

const (
	EngineAria2c DownloadEngineType = "aria2c"
	EngineCurl   DownloadEngineType = "curl"
	EngineGoHTTP DownloadEngineType = "go_http"
	EngineAuto   DownloadEngineType = "auto"
)

// IsValid reports whether the value matches a known variant.
func (d DownloadEngineType) IsValid() bool {
	switch d {
	case EngineAria2c, EngineCurl, EngineGoHTTP, EngineAuto:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (d DownloadEngineType) String() string {
	return string(d)
}

// ParseDownloadEngineType parses a raw string into a validated DownloadEngineType.
func ParseDownloadEngineType(s string) (DownloadEngineType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllDownloadEngineTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "enginearia2c":
		return EngineAria2c, nil
	case "aria2c":
		return EngineAria2c, nil
	case "enginecurl":
		return EngineCurl, nil
	case "curl":
		return EngineCurl, nil
	case "enginegohttp":
		return EngineGoHTTP, nil
	case "gohttp":
		return EngineGoHTTP, nil
	case "engineauto":
		return EngineAuto, nil
	case "auto":
		return EngineAuto, nil
	}
	return "", apperror.NewValidation("ParseDownloadEngineType", "E1000", fmt.Sprintf("invalid downloadenginetype: %q", s)).WithContext("value", s)
}

// AllDownloadEngineTypes returns all valid variants.
func AllDownloadEngineTypes() []DownloadEngineType {
	return []DownloadEngineType{
		EngineAria2c,
		EngineCurl,
		EngineGoHTTP,
		EngineAuto,
	}
}

// IsAria2c reports whether the value is EngineAria2c.
func (d DownloadEngineType) IsAria2c() bool {
	return d == EngineAria2c
}

// IsCurl reports whether the value is EngineCurl.
func (d DownloadEngineType) IsCurl() bool {
	return d == EngineCurl
}

// IsGoHTTP reports whether the value is EngineGoHTTP.
func (d DownloadEngineType) IsGoHTTP() bool {
	return d == EngineGoHTTP
}

// IsAuto reports whether the value is EngineAuto.
func (d DownloadEngineType) IsAuto() bool {
	return d == EngineAuto
}

// HasEngine reports whether the value is non-empty and valid.
func (d DownloadEngineType) HasEngine() bool {
	return d.IsValid()
}
