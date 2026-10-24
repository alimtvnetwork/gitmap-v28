package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// ConfigKeyType is a strongly-typed enum.
type ConfigKeyType string

const (
	KeyDefaultSplitSize   ConfigKeyType = "DownloaderConfig.DefaultSplitSize"
	KeyLargeFileSplitSize ConfigKeyType = "DownloaderConfig.LargeFileSplitSize"
	KeyLargeFileThreshold ConfigKeyType = "DownloaderConfig.LargeFileThreshold"
	KeyTinyFileThreshold  ConfigKeyType = "DownloaderConfig.TinyFileThreshold"
	KeyTinyFileSplitSize  ConfigKeyType = "DownloaderConfig.TinyFileSplitSize"
)

// IsValid reports whether the value matches a known variant.
func (c ConfigKeyType) IsValid() bool {
	switch c {
	case KeyDefaultSplitSize, KeyLargeFileSplitSize, KeyLargeFileThreshold, KeyTinyFileThreshold, KeyTinyFileSplitSize:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (c ConfigKeyType) String() string {
	return string(c)
}

// ParseConfigKeyType parses a raw string into a validated ConfigKeyType.
func ParseConfigKeyType(s string) (ConfigKeyType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllConfigKeyTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "keydefaultsplitsize":
		return KeyDefaultSplitSize, nil
	case "defaultsplitsize":
		return KeyDefaultSplitSize, nil
	case "keylargefilesplitsize":
		return KeyLargeFileSplitSize, nil
	case "largefilesplitsize":
		return KeyLargeFileSplitSize, nil
	case "keylargefilethreshold":
		return KeyLargeFileThreshold, nil
	case "largefilethreshold":
		return KeyLargeFileThreshold, nil
	case "keytinyfilethreshold":
		return KeyTinyFileThreshold, nil
	case "tinyfilethreshold":
		return KeyTinyFileThreshold, nil
	case "keytinyfilesplitsize":
		return KeyTinyFileSplitSize, nil
	case "tinyfilesplitsize":
		return KeyTinyFileSplitSize, nil
	}
	return "", apperror.NewValidation("ParseConfigKeyType", "E1000", fmt.Sprintf("invalid configkeytype: %q", s)).WithContext("value", s)
}

// AllConfigKeyTypes returns all valid variants.
func AllConfigKeyTypes() []ConfigKeyType {
	return []ConfigKeyType{
		KeyDefaultSplitSize,
		KeyLargeFileSplitSize,
		KeyLargeFileThreshold,
		KeyTinyFileThreshold,
		KeyTinyFileSplitSize,
	}
}

// IsDefaultSplitSize reports whether the value is KeyDefaultSplitSize.
func (c ConfigKeyType) IsDefaultSplitSize() bool {
	return c == KeyDefaultSplitSize
}

// IsLargeFileSplitSize reports whether the value is KeyLargeFileSplitSize.
func (c ConfigKeyType) IsLargeFileSplitSize() bool {
	return c == KeyLargeFileSplitSize
}

// IsLargeFileThreshold reports whether the value is KeyLargeFileThreshold.
func (c ConfigKeyType) IsLargeFileThreshold() bool {
	return c == KeyLargeFileThreshold
}

// IsTinyFileThreshold reports whether the value is KeyTinyFileThreshold.
func (c ConfigKeyType) IsTinyFileThreshold() bool {
	return c == KeyTinyFileThreshold
}

// IsTinyFileSplitSize reports whether the value is KeyTinyFileSplitSize.
func (c ConfigKeyType) IsTinyFileSplitSize() bool {
	return c == KeyTinyFileSplitSize
}

// HasKey reports whether the value is non-empty and valid.
func (c ConfigKeyType) HasKey() bool {
	return c.IsValid()
}
